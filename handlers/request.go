package handlers

import (
	"fmt"
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

	// TODO: introduce number of retries on each queue
	// TODO: Unit test the code

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("CUE_DB_HOST"), os.Getenv("CUE_DB_USER"), os.Getenv("CUE_DB_PASS"), os.Getenv("CUE_DB_NAME"))

	db, err := models.NewDB(dataSource)
	if err != nil {
		fmt.Println(err)
	}

	env := &Env{db}

	dbQueues, err := env.db.GetQueues()
	if err != nil {
		fmt.Println(err)
	}

	for i := range dbQueues {

		fmt.Println("Creating queue:", dbQueues[i].Name)
		queue := models.CreateQueue(dbQueues[i])

		fmt.Println("Adding", queue.Name, "queue to the Pool of queues")
		h.Pool.Add(queue.Name, queue)

		fmt.Printf("Creating dispatcher with %v workers\n", queue.Workers)
		dispatcher := models.CreateDispatcher(queue.Workers)

		fmt.Println("Starting the dispatcher")
		dispatcher.Start(queue)
		dispatcher.Listen()
	}

}

func (h TaskRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	task := models.CreateTask(h.Payload.TaskName, h.Payload.Content, time.Duration(h.Payload.Delay)*time.Second)

	fmt.Println("Received", task.Name, "with delay", task.Delay)

	h.Pool[h.Payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}
