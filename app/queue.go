package app

type Queue struct {
	QueueID  int `db:"queue_id"`
	Name     string
	Mode     string
	Workers  int
	Endpoint Endpoint
	Tasks    chan Task
}

type QueuesPool map[string]*Queue

type QueueService interface {
	CreateQueues() ([]*Queue, error)
	Add(queue *Queue) QueuesPool
}
