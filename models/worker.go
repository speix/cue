package models

import (
	"crypto/tls"
	"encoding/json"
	"errors"
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
				fmt.Println("Pulled", task.Name, "by worker", w.ID, "with delay", task.Delay)

				time.Sleep(task.Delay)

				result := Result{task: &task, worker: &w}

				client := &http.Client{
					Timeout: time.Duration(w.queue.Endpoint.Timeout * time.Second),
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}

				request, err := http.NewRequest("POST", w.queue.Endpoint.Url, task.Payload)
				if err != nil {
					result.Error = err
					result.message = "Failed to prepare request: " + err.Error()
					results <- result
					break
				}

				for h := range w.queue.Endpoint.Headers {
					request.Header.Add(w.queue.Endpoint.Headers[h].Key, w.queue.Endpoint.Headers[h].Value)
				}

				response, err := client.Do(request)
				if err != nil {
					result.Error = err
					result.message = "Failed to execute request: " + err.Error()
					results <- result
					break
				}

				if response.StatusCode != 200 {

					body := &responseBody{}

					err = json.NewDecoder(response.Body).Decode(&body)
					if err != nil {
						result.Error = err
						result.message = "Response error: " + response.Status + " " + body.Message
						results <- result
						break
					}

					result.Error = errors.New("")
					result.message = "Unable to connect: " + response.Status
					results <- result
					break
				}

				result.message = "Finished processing: " + task.Name + " response " + response.Status
				results <- result

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
