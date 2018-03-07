package models

type Queue struct {
	Name        string
	Mode        string
	Subscribers Endpoints
	Tasks       chan Task
}

type Queues map[string]*Queue

func (pool Queues) Add(reference string, queue *Queue) Queues {
	pool[reference] = queue
	return pool
}

func CreateQueue(name string, mode string) *Queue {
	return &Queue{
		Name:  name,
		Mode:  mode,
		Tasks: make(chan Task, 100),
	}
}
