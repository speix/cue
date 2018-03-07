package models

import (
	"time"
)

type Task struct {
	Name  string
	Delay time.Duration
}

func CreateTask(name string, delay time.Duration) *Task {
	return &Task{
		Name:  name,
		Delay: delay,
	}
}
