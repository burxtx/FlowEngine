package domain

type Candidate struct {
	ID       int64
	Approver string
	Node     *Node
	PID      *FlowInstance
}
