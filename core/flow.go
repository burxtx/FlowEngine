package core

type Flow interface {
	AddFirst()
	AddLast()
	AddBefore()
	AddAfter()
	Remove()
	First()
	Last()
	Run()
}
