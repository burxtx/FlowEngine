package core

type Job interface {
	Name() string
	Execute(...interface{})
}
