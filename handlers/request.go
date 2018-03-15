package handlers

import (
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
	Pool    models.Queues
}

func (h TaskRequestHandler) StartCue() {

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("CUE_DB_HOST"), os.Getenv("CUE_DB_USER"), os.Getenv("CUE_DB_PASS"), os.Getenv("CUE_DB_NAME"))

	db, err := models.NewDB(dataSource)
	if err != nil {
		fmt.Println(err)
	}

	env := &Env{db}

	queues, err := env.db.CreateQueues()
	if err != nil {
		log.Fatal(err)
	}

	h.Payload.QMapper = make(map[string]bool)

	for i := range queues {

		h.Pool.Add(queues[i])

		h.Payload.QMapper[queues[i].Name] = true // Add available queue names to the Payload as reference

		dispatcher := models.CreateDispatcher(queues[i].Workers)

		dispatcher.Start(queues[i])

		dispatcher.Listen()
	}

}

func (h TaskRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	task := models.CreateTask(h.Payload.TaskName, h.Payload.Payload, time.Duration(h.Payload.Delay)*time.Second)

	fmt.Println("Received", task.Name, "with delay", task.Delay)

	h.Pool[h.Payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}
