package app

type Dispatcher struct {
	WorkerPool chan chan Task // pool of worker's channels that are registered with the dispatcher
	NWorkers   int            // number of worker for each pool
}

type Result struct {
	Error   error
	Task    *Task
	Worker  *Worker
	Message string
}

type DispatcherService interface {
	CreateDispatcher(n int) *Dispatcher
	Start(q *Queue)
	Listen()
	Dispatch(q *Queue)
}

var Results = make(chan Result, 100)
