package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrQueueNotSet    = errors.New("queue is not set")
	ErrMessagesNotSet = errors.New("messages are not set")
	ErrInvalidDelay   = errors.New("delay must be in seconds between 0 and 1800")
	ErrInvalidJson    = errors.New("message format is not json")
	ErrQueueNotFound  = errors.New("queue not found")
)

type Payload struct {
	QueueName string          `json:"queue"`
	TaskName  string          `json:"task"`
	Messages  json.RawMessage `json:"messages"`
	Delay     int             `json:"delay"`
	QMapper   map[string]bool
}

type ServiceResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (p *Payload) QMap(queueName string) map[string]bool {
	p.QMapper[queueName] = true
	return p.QMapper
}

func (p *Payload) Validate(w http.ResponseWriter, r *http.Request) error {

	w.Header().Set("Content-Type", "application/json")
	response := ServiceResponse{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New(response.Message)
	}

	// check if there is a QueueName request
	if len(p.QueueName) == 0 {
		response.Error = true
		response.Message = "Queue is not set"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New(response.Message)
	}

	// check if there is a Payload request
	if len(p.Messages) == 0 {
		response.Error = true
		response.Message = "Messages is not set"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New(response.Message)
	}

	// check if delay is between 0 to 30 minutes
	if !inBetween(p.Delay, 0, 1800) {

		response.Error = true
		response.Message = "Delay must be in seconds between 0 and 1800."
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New(response.Message)
	}

	// check if payload attribute is a json field
	if !isJSON(p.Messages) {
		response.Error = true
		response.Message = "Message format is not json"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(400)
		w.Write(responseJson)

		return errors.New(response.Message)
	}

	// check if requested queue name is available
	if !p.QMapper[p.QueueName] {
		response.Error = true
		response.Message = "Queue " + p.QueueName + " not found"
		responseJson, _ := json.Marshal(response)

		w.WriteHeader(404)
		w.Write(responseJson)

		return errors.New(response.Message)
	}

	return nil
}

// inBetween checks if a given number is between two others
func inBetween(number, min, max int) bool {

	if (number >= min) && (number <= max) {
		return true
	}

	return false
}

// isJSON checks if a given string is a proper JSON object
func isJSON(str json.RawMessage) bool {
	var js json.RawMessage

	return json.Unmarshal([]byte(string(str)), &js) == nil
}
