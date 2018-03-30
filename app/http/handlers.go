package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/speix/cue/app"
)

type Handler struct {
	taskService app.TaskService
	payload     app.Payload
	pool        app.QueuesPool
}

func (h Handler) TaskRequestHandler(w http.ResponseWriter, r *http.Request) {

	task, err := h.taskService.CreateTask(h.payload.TaskName, h.payload.Messages, time.Duration(h.payload.Delay)*time.Second)
	if err != nil {
		response := ServiceResponse{
			Error:   true,
			Message: err.Error(),
		}
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)
		return
	}

	fmt.Printf("Received %v with delay %v\n", task.Name, task.Delay)

	h.pool[h.payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}
