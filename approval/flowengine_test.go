package approval

import (
	"context"
	"testing"

	"github.com/burxtx/FlowEngine/persistent"
)

var (
	MySQLInstance   = "USER:PASSWD@tcp(localhost:3306)/flowengine"
	MaxIdleConns    = 2
	MaxOpenConns    = 0
	ConnMaxLifeTime = 3
)

func InitEngine() *FlowEngine {
	cfg := persistent.Config{
		DSN:             MySQLInstance,
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 3,
	}
	persistent.NewDB(cfg)
	c := Config{StartPoint: 0}
	return NewEngine(persistent.GetDB(), c)
}

func TestCreateEngine(t *testing.T) {
	f := InitEngine()
	ctx := context.Background()
	approvers := [][]string{
		{"Thomas", "Sarah"},
		{"Thomas", "Sarah"},
		{"Thomas", "Sarah"},
	}
	nodes := []string{
		"创建工单",
		"开发审批",
		"经理审批",
		"部门审批",
		"结束",
	}
	submitter := "chris"
	processName := "提测流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Log(fi)
}

func TestList(t *testing.T) {
	f := InitEngine()
	ctx := context.Background()
	flowIDs, err := f.ListPendingInstances(ctx, "Thomas")
	if err != nil {
		t.Error(err)
	}
	t.Log(flowIDs)

	flowIDs, err = f.ListCompletedInstances(ctx, "Susan")
	if err != nil {
		t.Error(err)
	}
	t.Log(flowIDs)
}
