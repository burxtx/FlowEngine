package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/burxtx/FlowEngine/approval"
	"github.com/burxtx/FlowEngine/persistent"
)

var (
	MySQLInstance   = "USER:PASSWD@tcp(localhost:3306)/flowengine"
	MaxIdleConns    = 2
	MaxOpenConns    = 0
	ConnMaxLifeTime = 3
)

func MockGetFreeAmount(approver string) float64 {
	return 1
}

func InitEngine() *approval.FlowEngine {
	cfg := persistent.Config{
		DSN:             MySQLInstance,
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 3,
	}
	persistent.NewDB(cfg)
	c := approval.Config{StartPoint: 0}
	return approval.NewEngine(persistent.GetDB(), c)
}

func InitScheduler(user string) *ApprovalScheduler {
	cfg := persistent.Config{
		DSN:             MySQLInstance,
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 3,
	}
	persistent.NewDB(cfg)
	r := NewRule(100, MockGetFreeAmount, true, true, true, false, user)
	return NewApprovalScheduler(persistent.GetDB(), r)
}

func TestSingleApprover(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	approvers := [][]string{
		{"Thomas", "Sarah"},
	}
	nodes := []string{
		"创建工单",
		"经理审批",
		"结束",
	}
	submitter := "chris"
	processName := "单人审批流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)
	// 审批
	err = scheduler.Trigger(ctx, "Thomas", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
}

func TestSelfApprover(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	approvers := [][]string{
		{"Thomas", "Sarah", "chris"},
	}
	nodes := []string{
		"创建工单",
		"经理审批",
		"结束",
	}
	submitter := "chris"
	processName := "自己申请自己审批流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)
	// 审批
	err = scheduler.Trigger(ctx, "chris", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
}

func Test3ApproversContinuousApprove(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
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
	processName := "同一人自动审批流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)

	// 审批
	err = scheduler.Trigger(ctx, "Thomas", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
}

func TestRejectSchedule(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	approvers := [][]string{
		{"Thomas", "Sarah"},
		{"Daniel", "Susan"},
		{"William", "Lisa"},
	}
	nodes := []string{
		"创建工单",
		"开发审批",
		"经理审批",
		"部门审批",
		"结束",
	}
	submitter := "chris"
	processName := "拒绝流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)

	// 审批
	err = scheduler.Trigger(ctx, "Thomas", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)
	err = scheduler.Trigger(ctx, "Susan", "reject", "memo1", fi)
	if err != nil {
		t.Error(err)
	}
}

func TestRejectThenPassApprove(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	approvers := [][]string{
		{"Thomas", "Sarah"},
		{"Daniel", "Susan"},
		{"William", "Lisa"},
	}
	nodes := []string{
		"创建工单",
		"开发审批",
		"经理审批",
		"部门审批",
		"结束",
	}
	submitter := "chris"
	processName := "先拒绝后通过流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)

	// 审批
	err = scheduler.Trigger(ctx, "Thomas", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)
	err = scheduler.Trigger(ctx, "Susan", "reject", "memo1", fi)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	err = scheduler.Trigger(ctx, "Daniel", "pass", "memo2", fi)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	err = scheduler.Trigger(ctx, "Lisa", "pass", "memo3", fi)
	if err != nil {
		t.Error(err)
	}
}

func TestNormalPassApprove(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	approvers := [][]string{
		{"Thomas", "Sarah"},
		{"Daniel", "Susan"},
		{"William", "Lisa"},
	}
	nodes := []string{
		"创建工单",
		"开发审批",
		"经理审批",
		"部门审批",
		"结束",
	}
	submitter := "chris"
	processName := "正常审批通过流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)

	// 审批
	err = scheduler.Trigger(ctx, "Thomas", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)
	err = scheduler.Trigger(ctx, "Susan", "pass", "memo1", fi)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	err = scheduler.Trigger(ctx, "William", "pass", "memo2", fi)
	if err != nil {
		t.Error(err)
	}
}

func TestQuickPassApprove(t *testing.T) {
	f := InitEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	approvers := [][]string{
		{"Thomas", "Sarah"},
		{"Daniel", "Susan"},
		{"William", "Lisa"},
	}
	nodes := []string{
		"创建工单",
		"开发审批",
		"经理审批",
		"部门审批",
		"结束",
	}
	submitter := "chris"
	processName := "紧急审批通过流程"
	fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", fi)

	scheduler := InitScheduler(submitter)

	// 审批
	err = scheduler.Trigger(ctx, "emergency", "pass", "memo", fi)
	if err != nil {
		t.Error(err)
	}
}
