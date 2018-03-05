package models

import (
	"fmt"
	"time"
)

type Task struct {
	Name  string
	Delay time.Duration
}

func CreateTask(name string, delay time.Duration) Task {
	return Task{
		Name:  name,
		Delay: delay,
	}
}

func (t *Task) Process() {
	fmt.Println("Received task:", t.Name)
	time.Sleep(t.Delay)
	fmt.Println("Finished task", t.Name)
}
