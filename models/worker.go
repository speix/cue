package models

import (
	"fmt"
	"time"
)

type Worker struct {
	WorkerPool    chan chan Task
	WorkerChannel chan Task
	quit          chan bool
}

func CreateWorker(workerPool chan chan Task) *Worker {
	return &Worker{
		WorkerPool:    workerPool,
		WorkerChannel: make(chan Task),
		quit:          make(chan bool),
	}
}

// Start runs a loop for the worker execution
// listening for a quit signal in case we need to stop it
func (w Worker) Start() {
	go func() {
		for {

			// register current worker to the Worker Pool
			w.WorkerPool <- w.WorkerChannel

			select {

			case task := <-w.WorkerChannel:
				// received a work request, do some work
				fmt.Println("Pulled", task.Name)
				time.Sleep(task.Delay)

				// give back the response to the results channel
				results <- Result{task: &task, message: "Finished processing: " + task.Name}

			case <-w.quit:
				fmt.Println("quitting the channel")
				return

			}

		}
	}()
}

// Stop signals the worker to stop listening for work requests
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
