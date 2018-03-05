package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/speix/cue/handlers"

	"github.com/speix/cue/models"
)

func main() {

	tasks := make(chan models.Task, 100)

	i := 1
	for i = 1; i <= 20; i++ {
		go func(i int) {
			for task := range tasks {
				task.Process()
			}
		}(i)
	}
	fmt.Println("Spawned", i, "workers")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RequestHandler(tasks, w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
