package models

import (
	"bytes"
	"encoding/json"
	"time"
)

type Task struct {
	Name    string        // Name of the task
	Payload *bytes.Buffer // Payload to be forwarded to the next peer
	Delay   time.Duration // Run after x number of seconds
	Retries int           // Number of retries for each task on run time
}

func CreateTask(name string, payload json.RawMessage, delay time.Duration) *Task {

	return &Task{
		Name:    name,
		Payload: bytes.NewBufferString(string(payload)),
		Delay:   delay,
		Retries: 0,
	}
}
