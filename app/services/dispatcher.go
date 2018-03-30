package services

import (
	"fmt"
	"time"

	"github.com/speix/cue/app"
)

const (
	retryDelay = 5 * time.Second
)

type DispatcherService struct {
	dispatcher app.Dispatcher
	worker     app.WorkerService
}

func (s *DispatcherService) CreateDispatcher(nWorkers int) *app.Dispatcher {
	return &app.Dispatcher{
		WorkerPool: make(chan chan app.Task),
		NWorkers:   nWorkers,
	}
}

func (s *DispatcherService) Start(queue *app.Queue) {

	fmt.Println("Workers assignment for Queue:", queue.Name)

	for i := 1; i <= s.dispatcher.NWorkers; i++ {
		worker := s.worker.CreateWorker(s.dispatcher.WorkerPool, queue)
		fmt.Println("Spawned worker", worker.ID)
		s.worker.Start()
	}

	go s.Dispatch(queue)
}

func (s *DispatcherService) Dispatch(queue *app.Queue) {
	for {
		select {
		case task := <-queue.Tasks: // a task has been received

			go func(task app.Task) {

				// obtain an available worker Task channel
				// block until a worker is available
				taskChannel := <-s.dispatcher.WorkerPool

				// dispatch the task to the worker's Task channel
				taskChannel <- task

			}(task)
		}
	}
}

func (s *DispatcherService) Listen() {

	// Listen for results in results channel
	go func() {
		for result := range app.Results {

			fmt.Println(result.Message)

			if result.Error != nil {

				endpointRetries := result.Worker.Queue.Endpoint.Retries
				taskRetries := result.Task.Retries

				if endpointRetries != 0 && endpointRetries != taskRetries { // Retry working the task

					fmt.Printf("Retrying %v with delay %v\n", result.Task.Name, retryDelay)

					result.Task.Delay = retryDelay // Set up retry delay

					result.Task.Retries += 1

					result.Worker.Queue.Tasks <- *result.Task // Send the task back to the queue for pick up

				} else {
					fmt.Println("Discarding task:", result.Task.Name)
				}

			}

		}
	}()

}
