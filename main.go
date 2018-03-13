package main

import (
	"log"
	"net/http"

	"github.com/speix/cue/models"

	"github.com/speix/cue/helpers"

	"github.com/speix/cue/handlers"
)

func main() {

	// task request payload representation
	payload := &helpers.Payload{}

	// pool of different queues to be available by the system
	pool := models.Queues{}

	// package payload/request inside the handler
	request := handlers.TaskRequestHandler{
		Payload: payload,
		Pool:    pool,
	}

	// load available queues from database and fill up the pool, start the dispatcher with workers
	request.StartCue()

	server := &http.Server{
		Addr: ":8000",
	}

	// handle request on selected path but first function-chain the validation of each task request
	http.Handle("/", validate(request, payload, pool))

	log.Fatal(server.ListenAndServe())
}

func validate(h http.Handler, payload *helpers.Payload, pool models.Queues) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := payload.Validate(w, r, pool)
		if err != nil {
			return
		}

		h.ServeHTTP(w, r)
	})
}
