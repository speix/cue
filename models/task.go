package models

import (
	"fmt"
	"time"
)

type Task struct {
	Name   string
	Delay  time.Duration
	Result string
}

func CreateTask(name string, delay time.Duration) *Task {
	return &Task{
		Name:  name,
		Delay: delay,
	}
}

func (t *Task) Process() {
	fmt.Println("Processing:", t.Name)
	fmt.Println("===============================")
	time.Sleep(t.Delay)
	fmt.Println("Finished:", t.Name)
}
