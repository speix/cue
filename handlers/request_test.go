package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/speix/cue/helpers"
	"github.com/speix/cue/models"
)

func TestTaskRequestHandler_ServeHTTP(t *testing.T) {

	expected := 201
	queue := &models.Queue{
		Name:  "myQueue",
		Tasks: make(chan models.Task, 100),
	}

	h := &TaskRequestHandler{
		Payload: &helpers.Payload{
			QueueName: "myQueue",
			TaskName:  "myTask",
			Messages:  []byte(`[{"somekey": "somevalue"}]`),
			Delay:     10,
		},
		Pool: models.QueuesPool{},
	}

	h.Pool.Add(queue)

	handler := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	if expected != resp.StatusCode {
		t.Errorf("Expected status code %v got %v", expected, resp.StatusCode)
	}
}
