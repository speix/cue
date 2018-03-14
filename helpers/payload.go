package helpers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/speix/cue/models"
)

type Payload struct {
	QueueName string `json:"queue"`
	TaskName  string `json:"task"`
	Content   string `json:"payload"`
	Delay     int    `json:"delay"`
}

type ServiceResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (p *Payload) Validate(w http.ResponseWriter, r *http.Request, pool models.Queues) error {

	w.Header().Set("Content-Type", "application/json")
	response := ServiceResponse{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New("")
	}

	// check if there is a QueueName request
	if len(p.QueueName) == 0 {
		response.Error = true
		response.Message = "Queue is not set"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New("")
	}

	// check if there is a Payload request
	if len(p.Content) == 0 {
		response.Error = true
		response.Message = "Payload is not set"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New("")
	}

	// check if delay is between 0 to 30 minutes
	if !inBetween(p.Delay, 0, 1800) {

		response.Error = true
		response.Message = "Delay must be in seconds between 0 and 1800."
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New("")
	}

	// check if payload attribute is a json field
	if !isJSON(p.Content) {
		response.Error = true
		response.Message = "Payload format is not json"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New("")
	}

	if _, ok := pool[p.QueueName]; !ok {
		response.Error = true
		response.Message = "Queue " + p.QueueName + " not found"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(404)
		w.Write(responseJson)

		return errors.New("")
	}

	return nil
}

func inBetween(number, min, max int) bool {

	if (number >= min) && (number <= max) {
		return true
	}

	return false
}

func isJSON(str string) bool {
	var js json.RawMessage

	return json.Unmarshal([]byte(str), &js) == nil
}
