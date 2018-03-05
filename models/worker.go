package models

import "fmt"

type Worker struct{}

func SpawnWorkers(queue *Queue, nWorkers int) {

	for i := 0; i < nWorkers; i++ {
		fmt.Println("Swawned workers:", i+1)
		go func() {
			for task := range queue.Tasks {
				fmt.Println("Extracting:", task.Name)
				task.Process()
			}
		}()

	}
}
