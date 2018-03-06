package models

import (
	"time"
)

type Task struct {
	Name  string
	Delay time.Duration
	Queue string
}

func CreateTask(name string, delay time.Duration, queue string) *Task {
	return &Task{
		Name:  name,
		Delay: delay,
		Queue: queue,
	}
}
