package domain

const (
	// node states
	// Ready等待审批人审批
	Ready = "ready"
	// Complete审批人完成审批
	Complete = "complete"
	// Revoke 审批人撤回审批结果，中间状态
	// Revoke = "revoke"
	ReadyCN = "待审批"
	// Complete审批人完成审批
	CompleteCN = "已审批"
	// Revoke 审批人撤回审批结果，中间状态
	// RevokeCN = "撤回审批"

	// approve state
	Pass   = "pass"
	Reject = "reject"

	// flow instance state
	Created      = "created"
	Initializing = "initializing"
	Pending      = "pending"
	Finish       = "finish"

	CreatedCN      = "待审批"
	InitializingCN = "初始化中"
	PendingCN      = "审批中"
	FinishCN       = "审批完成"

	QuickApprover = "emergency"
	AutoPass      = "auto"
	AutoInit      = "auto_init"
)

var en2cn = map[string]string{
	Created:      CreatedCN,
	Initializing: InitializingCN,
	Pending:      PendingCN,
	Finish:       FinishCN,
}
