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
	Payload   string `json:"payload"`
	Delay     int    `json:"delay"`
}

type ServiceResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type Env struct {
	db models.Storage
}

var pool = make(models.Queues)

func init() {

	// TODO: validate queues against stored ones.
	// TODO: valid delay input (is a number in seconds between 1 and 1800.
	// TODO: extract request validation sequence to a separate method.
	// TODO: Unit test the code

	db, err := models.NewDB("host=192.168.10.70 user=et_psql password= dbname=etable sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	env := &Env{db}

	queues, err := env.db.GetQueues()
	if err != nil {
		fmt.Println(err)
	}

	for i := range queues {

		fmt.Println("Creating queue:", queues[i].Name)
		queue := models.CreateQueue(queues[i].Name, queues[i].Mode, queues[i].Workers)

		fmt.Println("Adding", queue.Name, "queue to the Pool of queues")
		pool.Add(queue.Name, queue)

		fmt.Printf("Creating dispatcher with %v workers\n", queue.Workers)
		dispatcher := models.CreateDispatcher(queue.Workers)

		fmt.Println("Starting the dispatcher")
		dispatcher.Start(queue)
		dispatcher.Listen()
	}

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

	// check if there is a QueueName request attribute
	if len(payload.QueueName) == 0 {
		response.Error = true
		response.Message = "Queue is empty"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return
	}

	// check if there is a Payload request attribute
	if len(payload.Payload) == 0 {
		response.Error = true
		response.Message = "Payload is empty"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return
	}

	// check if payload attribute is a json field
	if !IsJSON(payload.Payload) {
		response.Error = true
		response.Message = "Payload format is not json"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return
	}

	// check if pool of queues contains the requested QueueName
	if _, ok := pool[payload.QueueName]; !ok {
		response.Error = true
		response.Message = "Queue " + payload.QueueName + " not found"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(404)
		w.Write(responseJson)

		return
	}

	task := models.CreateTask(payload.TaskName, payload.Payload, time.Duration(payload.Delay)*time.Second)

	fmt.Println("Received", task.Name, "with delay", task.Delay)

	pool[payload.QueueName].Tasks <- *task

	w.WriteHeader(http.StatusCreated)
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
