package domain

import "fmt"

type Node struct {
	ID         int64
	PID        int64
	Name       string
	Children   *Node
	Parents    *Node
	Status     Result
	User       string
	UpdateTime int64
	Memo       string
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) Execute(callback func()) {
	fmt.Printf("executing node %s %d...", n.Name, n.ID)
}
