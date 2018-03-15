package models

import (
	"bytes"
	"encoding/json"
	"time"
)

type Task struct {
	Name    string
	Payload *bytes.Buffer
	Delay   time.Duration
	Retries int
}

func CreateTask(name string, payload json.RawMessage, delay time.Duration) *Task {

	return &Task{
		Name:    name,
		Payload: bytes.NewBufferString(string(payload)),
		Delay:   delay,
		Retries: 0,
	}
}
