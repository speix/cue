package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/speix/cue/models"
)

type AddTaskRequestContainer struct {
	QueueName string `json:"queuename"`
	TaskName  string `json:"taskname"`
	Delay     int    `json:"delay"`
}

type ServiceResponse struct {
	Message string `json:"message"`
}

func RequestHandler(tasks chan models.Task, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	response := ServiceResponse{}
	payload := AddTaskRequestContainer{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response.Message = err.Error()
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return
	}

	delay := time.Duration(payload.Delay) * time.Second
	task := models.CreateTask(payload.TaskName, delay)

	go func() {
		fmt.Printf("Added: %s Delay: %s\n", task.Name, task.Delay)
		tasks <- task
	}()
}
