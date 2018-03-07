package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/speix/cue/models"
)

type AddTaskRequestContainer struct {
	QueueName string `json:"queue"`
	TaskName  string `json:"task"`
	Delay     int    `json:"delay"`
}

type ServiceResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

var pool = make(models.Queues)

func init() {

	fmt.Println("Starting up Pool of Queues")
	queueEmail := models.CreateQueue("email", "push")
	queueSms := models.CreateQueue("sms", "push")

	fmt.Println("Adding available system queues")
	pool.Add("email", queueEmail)
	pool.Add("sms", queueSms)

	fmt.Println("Dispatching workers to each queue")

	dispatcherEmail := models.CreateDispatcher(1)
	//dispatcherSms := models.CreateDispatcher(10)

	dispatcherEmail.Start(pool["email"])
	//dispatcherSms.Start(pool["sms"])

	dispatcherEmail.Listen()
	//dispatcherSms.Listen()
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	payload := AddTaskRequestContainer{}
	response := ServiceResponse{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return
	}

	if len(payload.QueueName) == 0 {
		response.Error = true
		response.Message = errors.New("queue is empty").Error()
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return
	}

	if payload.QueueName != "email" && payload.QueueName != "sms" {
		response.Error = true
		response.Message = errors.New("unknown queue").Error()
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(404)
		w.Write(responseJson)

		return
	}

	task := models.CreateTask(payload.TaskName, time.Duration(payload.Delay)*time.Second)

	fmt.Println("Received", task.Name, "with delay", task.Delay)

	pool[payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}
