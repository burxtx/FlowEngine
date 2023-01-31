package approval

import (
	"errors"

	"github.com/burxtx/FlowEngine/approval/domain"
	"github.com/burxtx/FlowEngine/persistent/node"
	"gorm.io/gorm"
)

type Node struct {
	// ID       int64
	// name     string
	// Children *Node
	// Parents  *Node
	// Status   Result
	node.NodeRepository
}

// 节点权限判断
func (n *Node) Permitted(nodeID int64, user string) bool {
	return false
}

// 获取节点
func (n *Node) Take(nodeID int64, user string) error {
	return nil
}

// 释放节点
func (n *Node) Release(nodeID int64, user string) error {
	return nil
}

// 完成节点
func (n *Node) Complete(nodeID int64, user string) error {
	return nil
}

func (n *Node) Get(nodeID int64) (*domain.Node, error) {
	node := domain.NewNode()
	nodeDO, err := n.NodeRepository.Get(nodeID)
	if err != nil {
		return node, err
	}
	node.ID = nodeDO.ID
	node.Name = nodeDO.Name
	node.PID = nodeDO.ProcessID
	node.UpdateTime = nodeDO.UpdateTime
	node.Memo = nodeDO.Memo
	node.User = nodeDO.UserName
	node.Status = domain.Result{Name: nodeDO.Status}
	return node, nil
}

// 修改状态
func (n *Node) SetStatus(domNode *domain.Node, domResult domain.Result) error {
	newNode := node.Node{
		ID:       domNode.ID,
		Status:   domResult.Name,
		UserName: domNode.User,
		Memo:     domNode.Memo,
	}
	err := n.NodeRepository.Update(newNode)
	if err != nil {
		return err
	}
	domNode.Status = domResult
	return nil
}

func (n *Node) List(pID int64) ([]*domain.Node, error) {
	query := node.Node{
		ProcessID: pID,
	}
	nodesDO, err := n.NodeRepository.List(query)
	res := make([]*domain.Node, 0)
	for _, v := range nodesDO {
		nodeDOM, err := n.Load(v.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, nodeDOM)
	}

	return res, err
}

func (n *Node) UserQuery(approver, state string) ([]*domain.Node, error) {
	query := node.Node{
		UserName: approver,
		Status:   state,
	}
	nodesDO, err := n.NodeRepository.List(query)
	res := make([]*domain.Node, 0)
	for _, v := range nodesDO {
		nodeDOM, err := n.Load(v.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, nodeDOM)
	}

	return res, err
}

// Prepare
func (n *Node) Prepare(name string, pID int64) (*domain.Node, error) {
	ni := domain.NewNode()
	// 序列化node
	newNode := node.NewNode(name, pID)
	nodeID, err := n.NodeRepository.Save(newNode)
	if err != nil {
		return nil, err
	}
	ni.Name = name
	ni.ID = nodeID
	return ni, nil
}

func (n *Node) Connect(pre, next, pID int64) error {
	// 序列化node
	newNodeLine := node.NewNodeLine(pre, next, pID)
	_, err := n.NodeRepository.AddLine(newNodeLine)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) LoadPreNode(next int64) (*domain.Node, error) {
	query := node.NodeLine{
		Next: next,
	}
	v, err := n.GetLine(query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	parent, err := n.Get(v.Pre)
	if err != nil {
		return nil, err
	}
	return parent, nil
}

func (n *Node) LoadNextNode(pre int64) (*domain.Node, error) {
	query := node.NodeLine{Pre: pre}
	v, err := n.GetLine(query)
	if err != nil {
		return nil, err
	}
	if v.Next == 0 {
		return nil, nil
	}
	next, err := n.Get(v.Next)
	if err != nil {
		return nil, err
	}
	return next, nil
}

func (n *Node) Load(id int64) (*domain.Node, error) {
	v, err := n.Get(id)
	if err != nil {
		return nil, err
	}
	pre, err := n.LoadPreNode(id)
	if err != nil {
		return nil, err
	}
	next, err := n.LoadNextNode(id)
	if err != nil {
		return nil, err
	}
	v.Parents = pre
	v.Children = next

	return v, nil
}

func NewNode(repo node.NodeRepository) *Node {
	return &Node{
		NodeRepository: repo,
	}
}
