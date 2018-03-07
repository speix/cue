package models

import "fmt"

type Dispatcher struct {
	WorkerPool chan chan Task // pool of worker's channels that are registered with the dispatcher
	nWorkers   int            // number of worker for each pool
}

type Result struct {
	task    *Task
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
			// Log the results in file or db
			fmt.Println(result.message)
		}
	}()
}
