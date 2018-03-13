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

func CreateTask(name, content string, delay time.Duration) *Task {
	return &Task{
		Name:    name,
		Payload: bytes.NewBufferString(content),
		Delay:   delay,
	}
}
