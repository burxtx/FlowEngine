package approval

const (
	ErrApproveNotPermitted   = "当前无审批权限"
	ErrFlowInstantNotCreated = "审批流未初始化"
)

type Config struct {
	StartPoint               int
	ReApproveSameOneRestrict bool
}
