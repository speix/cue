package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PseudoRequestBody struct {
	Name string
}

func TestPayload_Validate_MissingRequestBody(t *testing.T) {
	expected := "EOF"
	responseContainer := &ServiceResponse{}

	payload := &Payload{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		payload.Validate(w, r)
	}

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	json.NewDecoder(resp.Body).Decode(responseContainer)

	if expected != responseContainer.Message {
		t.Errorf("Expected %v got %v", expected, responseContainer.Message)
	}
}

func TestPayload_Validate_QueueIsSet(t *testing.T) {
	expected := "Queue is not set"
	responseContainer := &ServiceResponse{}
	b, _ := json.Marshal(PseudoRequestBody{"Test"})

	payload := &Payload{
		TaskName: "task",
		Payload:  json.RawMessage{},
		Delay:    10,
		QMapper:  make(map[string]bool),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		payload.Validate(w, r)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	json.NewDecoder(resp.Body).Decode(responseContainer)

	if expected != responseContainer.Message {
		t.Errorf("Expected %v got %v", expected, responseContainer.Message)
	}
}

func TestPayload_Validate_PayloadIsSet(t *testing.T) {

	expected := "Payload is not set"
	responseContainer := &ServiceResponse{}
	b, _ := json.Marshal(PseudoRequestBody{"Test"})

	payload := &Payload{
		QueueName: "queue",
		TaskName:  "task",
		Payload:   json.RawMessage{},
		Delay:     10,
		QMapper:   make(map[string]bool),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		payload.Validate(w, r)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	json.NewDecoder(resp.Body).Decode(responseContainer)

	if expected != responseContainer.Message {
		t.Errorf("Expected %v got %v", expected, responseContainer.Message)
	}
}

func TestPayload_Validate_QueueExists(t *testing.T) {

	responseContainer := &ServiceResponse{}
	b, _ := json.Marshal(PseudoRequestBody{"Test"})

	jsonData := []byte(`{"somekey": "somevalue"}`)
	availableQueue := "email_queue"
	expected := "Queue " + availableQueue + " not found"

	payload := &Payload{
		QueueName: "email_queue",
		TaskName:  "task",
		Payload:   (json.RawMessage)(jsonData),
		Delay:     10,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		payload.Validate(w, r)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	json.NewDecoder(resp.Body).Decode(responseContainer)

	if expected != responseContainer.Message {
		t.Errorf("Expected %v got %v", expected, responseContainer.Message)
	}
}
