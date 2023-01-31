package domain

type FlowInstance struct {
	ID          int64
	name        string
	state       string
	nodes       []*Node
	size        int
	currentNode *Node
	HeadNode    *Node
	TailNode    *Node
	Result      Result
	biz         BizInfo
}

func NewFlowInstance() *FlowInstance {
	res := Result{Name: Initializing}
	return &FlowInstance{
		Result: res,
		size:   0,
	}
}

func (f *FlowInstance) SetCurrentNode(v *Node) {
	f.currentNode = v
}

func (f *FlowInstance) SetNodes(v []*Node) {
	f.nodes = v
}

func (f *FlowInstance) GetNodes() []*Node {
	return f.nodes
}

func (f *FlowInstance) SetTailNode(v *Node) {
	f.TailNode = v
}

func (f *FlowInstance) SetHeadNode(v *Node) {
	f.HeadNode = v
}

func (f *FlowInstance) GetCurrentNode() *Node {
	return f.currentNode
}
