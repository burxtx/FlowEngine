package scheduler

type Scheduler interface {
	// Start starts the scheduler
	Start()
	Schedule(job Job) error
	Stop()
}
