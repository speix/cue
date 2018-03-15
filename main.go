package main

import (
	"log"
	"net/http"
	"os"

	"github.com/speix/cue/models"

	"github.com/speix/cue/helpers"

	"github.com/speix/cue/handlers"
)

func main() {

	// package payload/request inside the handler
	request := handlers.TaskRequestHandler{
		Payload: &helpers.Payload{},
		Pool:    models.Queues{},
	}

	// load available queues from database and fill up the pool, start the dispatcher with workers
	request.StartCue()

	server := &http.Server{
		Addr: ":" + os.Getenv("CUE_SERVER_PORT"),
	}

	// handle request on selected path but first function-chain the validation of each task request
	http.Handle("/", validate(request, request.Payload))

	log.Fatal(server.ListenAndServe())
}

func validate(h http.Handler, payload *helpers.Payload) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := payload.Validate(w, r)
		if err != nil {
			return
		}

		h.ServeHTTP(w, r)
	})
}
