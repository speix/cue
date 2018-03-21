package models

import (
	"fmt"
	"time"
)

type Dispatcher struct {
	WorkerPool chan chan Task // pool of worker's channels that are registered with the dispatcher
	nWorkers   int            // number of worker for each pool
}

const (
	retryDelay = 5 * time.Second
)

type Result struct {
	Error   error
	task    *Task
	worker  *Worker
	message string
}

var results = make(chan Result, 100)

func CreateDispatcher(nWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Task),
		nWorkers:   nWorkers,
	}
}

func (d *Dispatcher) Start(queue *Queue) {

	fmt.Println("Workers assignment for Queue:", queue.Name)

	for i := 1; i <= d.nWorkers; i++ {
		worker := CreateWorker(d.WorkerPool, queue)
		fmt.Println("Spawned worker", worker.ID)
		worker.Start()
	}

	go d.dispatch(queue)
}

func (d *Dispatcher) dispatch(queue *Queue) {
	for {
		select {
		case task := <-queue.Tasks: // a task has been received

			go func(task Task) {

				// obtain an available worker Task channel
				// block until a worker is available
				taskChannel := <-d.WorkerPool

				// dispatch the task to the worker's Task channel
				taskChannel <- task

			}(task)
		}
	}
}

func (d *Dispatcher) Listen() {

	// Listen for results in results channel
	go func() {
		for result := range results {

			if result.Error != nil {

				endpointRetries := result.worker.queue.Endpoint.Retries
				taskRetries := result.task.Retries

				fmt.Println("Finished", result.task.Name)
				fmt.Println(result.message)

				if endpointRetries != 0 && endpointRetries != taskRetries { // Retry working the task

					fmt.Println("Retrying task", result.task.Name)

					result.task.Delay = retryDelay // Set up retry delay

					result.task.Retries += 1

					result.worker.queue.Tasks <- *result.task // Send the task back to the queue for pick up

				}

			}

		}
	}()

}
