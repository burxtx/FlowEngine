package approval

import (
	"context"
	"sync"

	"github.com/burxtx/FlowEngine/approval/domain"
	"github.com/burxtx/FlowEngine/persistent"
	"github.com/burxtx/FlowEngine/persistent/candidate"
	"github.com/burxtx/FlowEngine/persistent/node"
	"github.com/burxtx/FlowEngine/persistent/process"

	"gorm.io/gorm"
)

type FlowEngine struct {
	mu  sync.Mutex
	db  *gorm.DB
	cfg Config
}

// Create 构造节点，候选人，并持久化
func (f *FlowEngine) Create(ctx context.Context, approvers [][]string,
	nodesName []string, submitter, processName string) (flowInstance *domain.FlowInstance, err error) {
	tx := persistent.GetDBFromCtx(ctx)
	commit := false
	if tx == nil {
		tx = f.db.Begin()
		commit = true
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if commit {
			err = tx.Commit().Error
		}
	}()
	nodePst := f.NodeRepo(tx)
	ns := NewNode(nodePst)

	cdtPst := f.CandidateRepo(tx)
	cdt := NewNodeCandidate(cdtPst)

	pPst := f.ProcessRepo(tx)
	finst := NewInstance(pPst)

	// 初始化审批流实例
	flowInstance, err = finst.Init(processName)
	if err != nil {
		return nil, err
	}

	// 持久化节点
	// TODO 处理头尾节点
	flowNodes := make([]*domain.Node, 0)
	// 绑定节点和审批人
	users := make([][]string, 0)
	users = append(users, []string{submitter})
	users = append(users, approvers...)
	users = append(users, []string{"system"})

	for i := 0; i < len(nodesName); i++ {
		node, err := ns.Prepare(nodesName[i], flowInstance.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		for j := 0; j < len(users[i]); j++ {
			// if i == 0 {
			// 	continue
			// }
			err := cdt.Prepare(users[i][j], node.ID, flowInstance.ID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		flowNodes = append(flowNodes, node)
	}
	// 定义节点流转
	// 普通流程下的节点各有一个父子节点
	for i := 0; i < len(flowNodes); i++ {
		var nextID int64
		if i == len(flowNodes)-1 {
			nextID = 0
			flowNodes[i].Children = nil
		} else {
			nextID = flowNodes[i+1].ID
			flowNodes[i].Children = flowNodes[i+1]
		}
		// 持久化
		err := ns.Connect(flowNodes[i].ID, nextID, flowInstance.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// 初始化当前节点
	// 第一个节点是提交人,从审批人节点开始
	curNode, err := ns.Load(flowNodes[f.cfg.StartPoint].ID)
	if err != nil {
		return flowInstance, err
	}
	err = finst.SetCurrentNode(flowInstance, curNode)
	if err != nil {
		return flowInstance, err
	}

	err = finst.SetHeadNode(flowInstance, flowNodes[0])
	if err != nil {
		return flowInstance, err
	}

	err = finst.SetTailNode(flowInstance, flowNodes[len(flowNodes)-1])
	if err != nil {
		return flowInstance, err
	}
	// init status
	err = ns.SetStatus(flowNodes[f.cfg.StartPoint], domain.Result{Name: domain.Ready})
	if err != nil {
		return flowInstance, err
	}

	if f.cfg.StartPoint == 1 {
		err = ns.SetStatus(flowNodes[0], domain.Result{Name: domain.Complete})
		if err != nil {
			return flowInstance, err
		}
	}

	err = finst.SetResult(flowInstance, domain.Result{Name: domain.Created})
	if err != nil {
		return nil, err
	}
	flowInstance.SetNodes(flowNodes)
	return
}

func (f *FlowEngine) Load(ctx context.Context, pid int64) (*domain.FlowInstance, error) {
	nodePst := f.NodeRepo(f.db)
	n := NewNode(nodePst)

	pPst := f.ProcessRepo(f.db)
	finst := NewInstance(pPst)

	procDO, err := finst.Get(pid)
	if err != nil {
		return nil, err
	}
	curNode, err := n.Get(procDO.CurrentNode)
	if err != nil {
		return nil, err
	}
	headNode, err := n.Get(procDO.HeadNode)
	if err != nil {
		return nil, err
	}
	tailNode, err := n.Get(procDO.TailNode)
	if err != nil {
		return nil, err
	}
	nodes, err := n.List(pid)
	if err != nil {
		return nil, err
	}
	fiDOM := domain.NewFlowInstance()
	fiDOM.SetCurrentNode(curNode)
	fiDOM.SetNodes(nodes)
	fiDOM.SetHeadNode(headNode)
	fiDOM.SetTailNode(tailNode)
	fiDOM.ID = pid
	return fiDOM, nil
}

func (f *FlowEngine) ListSuccessors(ctx context.Context, pid int64) ([]string, error) {
	nodePst := f.NodeRepo(f.db)
	n := NewNode(nodePst)

	pPst := f.ProcessRepo(f.db)
	finst := NewInstance(pPst)

	cdtPst := f.CandidateRepo(f.db)
	cdt := NewNodeCandidate(cdtPst)

	procDO, err := finst.Get(pid)
	if err != nil {
		return nil, err
	}
	curNode, err := n.Get(procDO.CurrentNode)
	if err != nil {
		return nil, err
	}
	if f.cfg.ReApproveSameOneRestrict && curNode.Status.Name == domain.Reject {
		return []string{curNode.User}, nil
	} else {
		successors, err := cdt.ListNodeUsers(curNode)
		if err != nil {
			return nil, err
		}
		return successors, nil
	}
}

func (f *FlowEngine) SetDB(tx *gorm.DB) {
	f.db = tx
}

func NewInstance(repo process.ProcessRepository) *FlowInstance {
	return newInstance(repo)
}
func NewEngine(db *gorm.DB, cfg Config) *FlowEngine {
	return &FlowEngine{
		db:  db,
		cfg: cfg,
	}
}

func (f *FlowEngine) NodeRepo(tx *gorm.DB) node.NodeRepository {
	r := &node.Repository{tx}
	return r
}

func (f *FlowEngine) ProcessRepo(tx *gorm.DB) process.ProcessRepository {
	r := &process.Repository{tx}
	return r
}

func (f *FlowEngine) CandidateRepo(tx *gorm.DB) candidate.CandidateRepository {
	r := &candidate.Repository{tx}
	return r
}

func (f *FlowEngine) ListPendingInstances(ctx context.Context, approver string) ([]int64, error) {
	tx := f.db.Begin()

	nodePst := f.NodeRepo(tx)
	ns := NewNode(nodePst)

	cdtPst := f.CandidateRepo(tx)
	cdt := NewNodeCandidate(cdtPst)

	res := make([]int64, 0)

	userNodes, err := cdt.ListUserNodes(approver)
	if err != nil {
		return res, err
	}
	tmpMap := make(map[int64]bool, 0)
	for _, userNode := range userNodes {
		domNode, err := ns.Load(userNode.NodeID)
		if err != nil {
			return res, err
		}
		if domNode.Status.Name == domain.Ready ||
			(f.cfg.ReApproveSameOneRestrict &&
				domNode.Status.Name == domain.Reject &&
				domNode.User == approver) {
			if _, exist := tmpMap[domNode.PID]; !exist {
				tmpMap[domNode.PID] = true
			}
		}
	}
	for k := range tmpMap {
		res = append(res, k)
	}
	return res, nil
}

func (f *FlowEngine) ListCompletedInstances(ctx context.Context, approver string) ([]int64, error) {
	tx := f.db.Begin()

	nodePst := f.NodeRepo(tx)
	ns := NewNode(nodePst)

	res := make([]int64, 0)
	tmpMap := make(map[int64]bool, 0)
	domNodes, err := ns.UserQuery(approver, domain.Complete)
	if err != nil {
		return res, err
	}
	domNodesRejected, err := ns.UserQuery(approver, domain.Reject)
	if err != nil {
		return res, err
	}
	domNodes = append(domNodes, domNodesRejected...)
	for _, node := range domNodes {
		if node.Memo == domain.AutoInit {
			continue
		}
		if _, exist := tmpMap[node.PID]; !exist {
			tmpMap[node.PID] = true
		}
	}
	for k := range tmpMap {
		res = append(res, k)
	}

	return res, nil
}

func (f *FlowEngine) ResetFlow(ctx context.Context, pid int64) (err error) {
	nodePst := f.NodeRepo(f.db)
	ns := NewNode(nodePst)

	pPst := f.ProcessRepo(f.db)
	finst := NewInstance(pPst)

	fiDOM, err := f.Load(ctx, pid)
	if err != nil {
		return
	}
	nodes := fiDOM.GetNodes()
	for i, v := range nodes {
		if i == 0 {
			err = ns.SetStatus(v, domain.Result{Name: domain.Complete})
			if err != nil {
				return
			}
			continue
		}
		if i == 1 {
			err = ns.SetStatus(v, domain.Result{Name: domain.Ready})
			if err != nil {
				return
			}
			continue
		}
		// reset other nodes
		err = ns.SetStatus(v, domain.Result{Name: ""})
		if err != nil {
			return
		}
	}
	curNode, err := ns.Load(nodes[1].ID)
	if err != nil {
		return
	}
	err = finst.SetCurrentNode(fiDOM, curNode)
	if err != nil {
		return
	}
	err = finst.SetResult(fiDOM, domain.Result{Name: domain.Created})
	if err != nil {
		return
	}
	return nil
}
