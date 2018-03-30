package services

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/speix/cue/app"
	"github.com/twinj/uuid"
)

type WorkerService struct {
	worker app.Worker
}

type workerServiceResponse struct {
	Message string
}

func (service *WorkerService) CreateWorker(workerPool chan chan app.Task, queue *app.Queue) *app.Worker {
	return &app.Worker{
		ID:            uuid.NewV4(),
		WorkerPool:    workerPool,
		WorkerChannel: make(chan app.Task),
		Queue:         queue,
		Quit:          make(chan bool),
	}
}

// Start runs a loop for the worker execution
// listening for a quit signal in case we need to stop it
func (service *WorkerService) Start() {

	go func() {
		for {

			// register current worker to the Worker Pool
			service.worker.WorkerPool <- service.worker.WorkerChannel

			select {

			case task := <-service.worker.WorkerChannel:

				// received a work request, do some work
				time.Sleep(task.Delay)

				fmt.Printf("Running %v by worker %v\n", task.Name, service.worker.ID)

				app.Results <- service.execute(&task)

			case <-service.worker.Quit:
				fmt.Println("quitting the channel")
				return

			}

		}
	}()

}

// Stop signals the worker to stop listening for work requests
func (service *WorkerService) Stop() {
	go func() {
		service.worker.Quit <- true
	}()
}

// execute forwards each Message of the Task to the Endpoint Url
func (service *WorkerService) execute(task *app.Task) app.Result {

	result := app.Result{Task: task, Worker: &service.worker}

	client := &http.Client{
		Timeout: time.Duration(service.worker.Queue.Endpoint.Timeout * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	for i := range task.Messages {

		request, err := http.NewRequest(service.worker.Queue.Endpoint.Method, service.worker.Queue.Endpoint.Url, task.Messages[i])
		if err != nil {
			result.Error = err
			result.Message = "Failed to prepare request: " + err.Error()
			return result
		}

		for h := range service.worker.Queue.Endpoint.Headers {
			request.Header.Add(service.worker.Queue.Endpoint.Headers[h].Key, service.worker.Queue.Endpoint.Headers[h].Value)
		}

		response, err := client.Do(request)
		if err != nil {
			result.Error = err
			result.Message = "Failed to execute request: " + err.Error()
			return result
		}

		if response.StatusCode != 200 {

			body := &workerServiceResponse{}

			err = json.NewDecoder(response.Body).Decode(&body)
			if err != nil {
				result.Error = err
				result.Message = "Response error: " + response.Status + " " + body.Message
				return result
			}

			result.Error = errors.New("")
			result.Message = "Unable to connect: " + response.Status
			return result
		}

		result.Message = "Finished: " + task.Name + " response " + response.Status
		response.Body.Close()
	}

	return result
}
