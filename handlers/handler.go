package handlers

import (
	"encoding/json"
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
	Message string `json:"message"`
}

func RequestHandler(queue *models.Queue, w http.ResponseWriter, r *http.Request) {

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

	task := models.CreateTask(payload.TaskName, time.Duration(payload.Delay)*time.Second)

	fmt.Printf("Added: %s Delay: %s\n", task.Name, task.Delay)

	queue.Tasks <- *task
}
