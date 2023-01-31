package main

import "fmt"

type Flow struct {
	head *Node
}

type Node struct {
	name string
}

func initFlow(v *Node) *Flow {
	return &Flow{
		head: v,
	}
}

func changeHead(fi *Flow, newHead *Node) {
	fi.head = newHead
}

func main() {
	node1 := new(Node)
	node1.name = "a"
	node2 := new(Node)
	node2.name = "b"
	flow := initFlow(node1)
	fmt.Printf("%#v\n", flow)
	changeHead(flow, node2)
	fmt.Printf("%#v\n", flow)
}
