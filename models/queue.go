package models

import "fmt"

type Queue struct {
	Name  string
	Mode  string
	Tasks chan Task
}

func CreateQueue(name string, mode string) *Queue {
	return &Queue{
		Name:  name,
		Mode:  mode,
		Tasks: make(chan Task, 100),
	}
}

func (q *Queue) SpawnWorkers(nWorkers int) {
	for i := 0; i < nWorkers; i++ {
		fmt.Println("Spawned workers:", i+1)
		go func() {
			for task := range q.Tasks {
				fmt.Println("Extracting:", task.Name)
				task.Process()
			}
		}()

	}
}
