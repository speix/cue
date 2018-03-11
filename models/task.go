package models

import (
	"bytes"
	"time"
)

type Task struct {
	Name    string
	Payload *bytes.Buffer
	Delay   time.Duration
}

func CreateTask(name, payload string, delay time.Duration) *Task {
	return &Task{
		Name:    name,
		Payload: bytes.NewBufferString(payload),
		Delay:   delay,
	}
}
