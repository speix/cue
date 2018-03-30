package app

import (
	"bytes"
	"encoding/json"
	"time"
)

type Task struct {
	Name     string          // Name of the task
	Messages []*bytes.Buffer // Messages to be forwarded to the next peer
	Delay    time.Duration   // Run after x number of seconds
	Retries  int             // Number of retries for each task on run time
}

type TaskService interface {
	CreateTask(name string, payload json.RawMessage, delay time.Duration) (*Task, error)
}
