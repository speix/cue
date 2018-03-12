package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/twinj/uuid"
)

type responseBody struct {
	Message string
}

type Worker struct {
	ID            uuid.Uuid
	WorkerPool    chan chan Task
	WorkerChannel chan Task
	queue         *Queue
	quit          chan bool
}

func CreateWorker(workerPool chan chan Task, queue *Queue) *Worker {
	return &Worker{
		ID:            uuid.NewV4(),
		WorkerPool:    workerPool,
		WorkerChannel: make(chan Task),
		queue:         queue,
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
				fmt.Println("Pulled", task.Name, "by worker", w.ID)
				time.Sleep(task.Delay)

				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}

				request, err := http.NewRequest("POST", "", task.Payload)

				for h := range w.queue.Headers {
					fmt.Println(h)
				}

				if err != nil {
					results <- Result{task: &task, message: "Failed to prepare request: " + err.Error()}
					break
				}

				response, err := client.Do(request)
				if err != nil {
					results <- Result{task: &task, message: "Failed to execute request: " + err.Error()}
					break
				}

				if response.StatusCode != 200 {
					body := &responseBody{}
					err = json.NewDecoder(response.Body).Decode(&body)

					if err != nil {
						results <- Result{task: &task, message: body.Message}
						break
					}

					results <- Result{task: &task, message: "Unable to connect: " + response.Status}
					break
				}

				// give back the response to the results channel
				results <- Result{task: &task, message: "Finished processing: " + task.Name + " response " + response.Status}

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
