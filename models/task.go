package models

import (
	"time"
)

type Task struct {
	Name    string
	Payload string
	Delay   time.Duration
}

func CreateTask(name, payload string, delay time.Duration) *Task {
	return &Task{
		Name:    name,
		Payload: payload,
		Delay:   delay,
	}
}
