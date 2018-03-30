package app

import "github.com/twinj/uuid"

type Worker struct {
	ID            uuid.Uuid
	WorkerPool    chan chan Task
	WorkerChannel chan Task
	Queue         *Queue
	Quit          chan bool
}

type WorkerService interface {
	CreateWorker(workerPool chan chan Task, queue *Queue) *Worker
	Start()
	Stop()
}
