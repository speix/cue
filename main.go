package main

import (
	"log"
	"net/http"

	"github.com/speix/cue/handlers"

	"github.com/speix/cue/models"
)

func main() {

	queue := models.CreateQueue("mail", "push")
	queue.SpawnWorkers(2)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RequestHandler(queue, w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
