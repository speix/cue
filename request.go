package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Env struct {
	db Storage
}

type TaskRequestHandler struct {
	Payload *Payload
	Pool    QueuesPool
}

func StartCue() *TaskRequestHandler {

	queues := loadQueues() // Load queues from database

	handler := &TaskRequestHandler{
		Payload: &Payload{
			QMapper: make(map[string]bool),
		},
		Pool: QueuesPool{},
	}

	for i := range queues {

		handler.Pool.Add(queues[i]) // Add queue to the pool of queues

		handler.Payload.QMap(queues[i].Name) // Add available queue names to the Payload as reference

		dispatcher := CreateDispatcher(queues[i].Workers) // Create a dispatcher for each queue

		dispatcher.Start(queues[i]) // Start workers running on each queue

		dispatcher.Listen() // Listen for tasks
	}

	return handler
}

func (h *TaskRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	task, err := CreateTask(h.Payload.TaskName, h.Payload.Messages, time.Duration(h.Payload.Delay)*time.Second)
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

	h.Pool[h.Payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}

func loadQueues() []*Queue {

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("CUE_DB_HOST"), os.Getenv("CUE_DB_USER"), os.Getenv("CUE_DB_PASS"), os.Getenv("CUE_DB_NAME"))

	db, err := NewDB(dataSource)
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
