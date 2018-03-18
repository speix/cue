package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/speix/cue/helpers"

	"github.com/speix/cue/models"
)

type Env struct {
	db models.Storage
}

type TaskRequestHandler struct {
	Payload *helpers.Payload
	Pool    models.QueuesPool
}

func StartCue() *TaskRequestHandler {

	queues := loadQueues() // Load queues from database

	handler := &TaskRequestHandler{
		Payload: &helpers.Payload{
			QMapper: make(map[string]bool),
		},
		Pool: models.QueuesPool{},
	}

	for i := range queues {

		handler.Pool.Add(queues[i]) // Add queue to the pool of queues

		handler.Payload.QMap(queues[i].Name) // Add available queue names to the Payload as reference

		dispatcher := models.CreateDispatcher(queues[i].Workers) // Create a dispatcher for each queue

		dispatcher.Start(queues[i]) // Start workers running on each queue

		dispatcher.Listen() // Listen for tasks
	}

	return handler
}

func (h *TaskRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	task, err := models.CreateTask(h.Payload.TaskName, h.Payload.Messages, time.Duration(h.Payload.Delay)*time.Second)
	if err != nil {
		response := helpers.ServiceResponse{
			Error:   true,
			Message: err.Error(),
		}
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)
		return
	}

	fmt.Println("Received", task.Name, "with delay", task.Delay)

	h.Pool[h.Payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}

func loadQueues() []*models.Queue {

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("CUE_DB_HOST"), os.Getenv("CUE_DB_USER"), os.Getenv("CUE_DB_PASS"), os.Getenv("CUE_DB_NAME"))

	db, err := models.NewDB(dataSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	env := &Env{db}

	queues, err := env.db.CreateQueues()
	if err != nil {
		log.Fatal(err.Error())
	}

	return queues
}
