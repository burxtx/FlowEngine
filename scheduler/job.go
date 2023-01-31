package scheduler

type Job interface {
	Execute()
	Name() string
}
