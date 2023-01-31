package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/burxtx/FlowEngine/approval"
	"github.com/burxtx/FlowEngine/approval/domain"
	"github.com/burxtx/FlowEngine/persistent"
	"github.com/burxtx/FlowEngine/persistent/candidate"
	"github.com/burxtx/FlowEngine/persistent/node"
	"github.com/burxtx/FlowEngine/persistent/process"
	persistentScheduler "github.com/burxtx/FlowEngine/persistent/scheduler"

	"gorm.io/gorm"
)

type ApprovalScheduler struct {
	db    *gorm.DB
	rule  Rule
	queue persistentScheduler.TaskQueue
	stop  chan int
}

type Rule struct {
	Amount                   bool
	OrderAmount              float64
	GetFreeAmountsFunc       func(string) float64
	Continuous               bool
	QuickApprove             bool
	ReApproveSameOneRestrict bool
	SkipMyCreated            string
}

// 阻塞方法，需要调用方自行处理error
func (as *ApprovalScheduler) Schedule(ctx context.Context, fi *domain.FlowInstance, cur *domain.Node) error {
	if cur == nil {
		return nil
	}

	// 遍历节点
	tail := fi.TailNode
	// for {
	// tx := as.db.Begin()
	tx := persistent.GetDBFromCtx(ctx)
	nodePst := as.NodeRepo(tx)
	ns := approval.NewNode(nodePst)

	pPst := as.ProcessRepo(tx)
	finst := approval.NewInstance(pPst)

	// cur, err := as.queue.Pop()
	// if err != nil {
	// 	return err
	// }
	//if cur == nil {
	//time.Sleep(500 * time.Millisecond)
	//	return nil
	//}
	// select {
	// // 当人为执行审批后执行
	// case :
	approver := cur.User
	memo := cur.Memo
	for {
		if cur.ID != tail.ID {
			if cur.Status.Name != domain.Ready && cur.Status.Name != domain.Reject {
				break
			}
			cur.User = approver
			cur.Memo = memo
			// 当前节点标记为完成
			err := ns.SetStatus(cur, domain.Result{Name: domain.Complete})
			if err != nil {
				tx.Rollback()
				return err
			}
			// 找到子节点, 之后cur代表子节点
			cur, err = ns.Load(cur.Children.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
			// 此时cur.User为空，需重新赋值approver
			cur.User = approver
			err = ns.SetStatus(cur, domain.Result{Name: domain.Ready})
			if err != nil {
				tx.Rollback()
				return err
			}
			err = finst.SetCurrentNode(fi, cur)
			if err != nil {
				tx.Rollback()
				return err
			}
			err = finst.SetResult(fi, domain.Result{Name: domain.Pending})
			if err != nil {
				tx.Rollback()
				return err
			}
			// fmt.Printf("curc: %#v\n", cur)

			if cur.ID == tail.ID {
				continue
			} else {
				autoPass, passer, err := as.AutoPassNext(tx, cur)
				if err != nil {
					tx.Rollback()
					return err
				}
				if !autoPass {
					break
				} else {
					approver = passer.approver
					memo = passer.memo
				}
			}
		} else {
			err := ns.SetStatus(cur, domain.Result{Name: domain.Complete})
			if err != nil {
				tx.Rollback()
				return err
			}
			err = finst.SetResult(fi, domain.Result{Name: domain.Finish})
			if err != nil {
				tx.Rollback()
				return err
			}
			return tx.Commit().Error
		}
	}
	// case <-as.stop:
	// 	return nil
	// case <-ctx.Done():
	// 	return nil
	// }
	// }
	return tx.Commit().Error
}

func (as *ApprovalScheduler) Trigger(ctx context.Context, approver, action, remark string, fi *domain.FlowInstance) error {
	permitted := false
	tx := as.db.Begin()
	txCtx := persistent.CtxWithTransaction(ctx, tx)
	nodePst := as.NodeRepo(tx)
	ns := approval.NewNode(nodePst)

	cdtPst := as.CandidateRepo(tx)
	cdt := approval.NewNodeCandidate(cdtPst)

	pPst := as.ProcessRepo(tx)
	finst := approval.NewInstance(pPst)

	if approver == domain.QuickApprover && as.rule.QuickApprove {
		err := finst.SetCurrentNode(fi, fi.TailNode)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = finst.SetResult(fi, domain.Result{Name: domain.Finish})
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}

	userNodes, err := cdt.ListUserNodes(approver, fi.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	curNodeDO, err := finst.GetCurrentNode(fi)
	if err != nil {
		tx.Rollback()
		return err
	}
	curNodeDOM, err := ns.Load(curNodeDO)
	if err != nil {
		tx.Rollback()
		return err
	}
	// fmt.Printf("current node: %#v\n", curNodeDOM)
	if as.rule.ReApproveSameOneRestrict && curNodeDOM.Status.Name == domain.Reject {
		if approver == curNodeDOM.User {
			permitted = true
		} else {
			return fmt.Errorf("只有驳回人有权继续审批")
		}
	} else {
		for i := 0; i < len(userNodes); i++ {
			if curNodeDOM.ID == userNodes[i].NodeID &&
				(curNodeDOM.Status.Name == domain.Ready) {
				permitted = true
				break
			}
		}
	}

	if permitted {
		curNodeDOM.User = approver
		curNodeDOM.Memo = remark

		var sq persistentScheduler.ScheduleQueue
		sq.NodeID = curNodeDOM.ID
		sq.User = curNodeDOM.User
		sq.State = action
		sq.Memo = curNodeDOM.Memo
		sq.Name = curNodeDOM.Name
		sq.ProcessID = curNodeDOM.PID
		serialized, err := json.Marshal(curNodeDOM)
		if err != nil {
			return err
		}
		sq.DomainNode = string(serialized)
		err = as.queue.Push(&sq)
		if err != nil {
			return err
		}

		if action == domain.Pass {
			err = as.Schedule(txCtx, fi, curNodeDOM)
			if err != nil {
				return err
			}
			return nil
		} else if action == domain.Reject {
			err := ns.SetStatus(curNodeDOM, domain.Result{Name: domain.Reject})
			if err != nil {
				tx.Rollback()
				return err
			}
			err = finst.SetResult(fi, domain.Result{Name: domain.Pending})
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		return fmt.Errorf("当前没有审批权限")
	}

	return tx.Commit().Error
}

func (as *ApprovalScheduler) Stop() {
	as.stop <- 1
}

func (as *ApprovalScheduler) ListTriggerHistory(ctx context.Context, pid int64) ([]*persistentScheduler.ScheduleQueue, error) {
	query := persistentScheduler.ScheduleQueue{
		ProcessID: pid,
	}
	return as.queue.List(query)
}

type AutoPasser struct {
	approver, memo string
}

func (as *ApprovalScheduler) AutoPassNext(tx *gorm.DB, cur *domain.Node) (bool, AutoPasser, error) {
	res := false
	AutoPasser := AutoPasser{}
	cdtPst := as.CandidateRepo(tx)
	cdt := approval.NewNodeCandidate(cdtPst)

	nextApprovers, err := cdt.ListNodeUsers(cur)
	if err != nil {
		return res, AutoPasser, err
	}

	amt := false
	if as.rule.Amount {
		canContinue, passedBy, err := as.ApplyAmountRule(nextApprovers)
		if err != nil {
			return res, AutoPasser, err
		}
		if canContinue {
			amt = true
			AutoPasser.approver = passedBy
			AutoPasser.memo = domain.AutoPass
		}
	}

	continuous := false
	if as.rule.Continuous {
		// fmt.Printf("curUser: %s, nextApprovers: %s \n", cur.User, nextApprovers)
		canContinue, err := as.ApplyContinuousRule(cur.User, nextApprovers)
		if err != nil {
			return res, AutoPasser, err
		}
		if canContinue {
			continuous = true
			AutoPasser.approver = cur.User
			AutoPasser.memo = domain.AutoPass
		}
	}

	skipMyCreated := false
	if len(as.rule.SkipMyCreated) > 0 {
		canSkipMyCreated, err := as.ApplySkipMyCreatedRule(nextApprovers)
		if err != nil {
			return res, AutoPasser, err
		}
		if canSkipMyCreated {
			skipMyCreated = true
			AutoPasser.approver = as.rule.SkipMyCreated
			AutoPasser.memo = domain.AutoPass
		}
	}

	if skipMyCreated || continuous || amt {
		return true, AutoPasser, nil
	} else {
		return false, AutoPasser, nil
	}
}

func (as *ApprovalScheduler) ApplySkipMyCreatedRule(nextApprovers []string) (bool, error) {
	autoPass := false
	for _, v := range nextApprovers {
		if v == as.rule.SkipMyCreated {
			autoPass = true
			break
		}
	}
	return autoPass, nil
}

func (as *ApprovalScheduler) ApplyContinuousRule(curApprover string, nextApprovers []string) (bool, error) {
	autoPass := false
	for _, approver := range nextApprovers {
		if curApprover == approver {
			autoPass = true
		}
	}
	return autoPass, nil
}

func (as *ApprovalScheduler) ApplyAmountRule(nextApprovers []string) (bool, string, error) {
	autoPass := false
	passedBy := ""
	for _, approver := range nextApprovers {
		amount := as.rule.GetFreeAmountsFunc(approver)
		if amount >= as.rule.OrderAmount {
			autoPass = true
			passedBy = approver
		}
	}
	return autoPass, passedBy, nil
}

func NewRule(orderAmount float64, fn func(string) float64,
	amount, continous, quickApprove, re bool, skipMy string) Rule {
	return Rule{
		OrderAmount:              orderAmount,  //自动过审金额
		GetFreeAmountsFunc:       fn,           //获取工单金额的方法
		Amount:                   amount,       //根据金额自动审批
		Continuous:               continous,    //连续审批
		QuickApprove:             quickApprove, //快速结单
		ReApproveSameOneRestrict: re,           //审批人驳回后，只能继续由驳回人审批
		SkipMyCreated:            skipMy,       //申请人和审批人为同一人时，自动过审
	}
}

func NewApprovalScheduler(db *gorm.DB, rule Rule) *ApprovalScheduler {
	taskQueue := persistentScheduler.NewTaskQueueInstance(db)
	return &ApprovalScheduler{
		db:    db,
		stop:  make(chan int),
		queue: taskQueue,
		rule:  rule,
	}
}

func (as *ApprovalScheduler) NodeRepo(tx *gorm.DB) node.NodeRepository {
	r := &node.Repository{tx}
	return r
}

func (as *ApprovalScheduler) ProcessRepo(tx *gorm.DB) process.ProcessRepository {
	r := &process.Repository{tx}
	return r
}

func (as *ApprovalScheduler) CandidateRepo(tx *gorm.DB) candidate.CandidateRepository {
	r := &candidate.Repository{tx}
	return r
}
