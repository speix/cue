package main

import (
	"log"
	"net/http"
	"os"

	"github.com/speix/cue/helpers"

	"github.com/speix/cue/handlers"
)

func main() {

	cue := handlers.StartCue() // Start the Cue (Queues, Dispatchers, Workers, Listeners)

	server := &http.Server{
		Addr: ":" + os.Getenv("CUE_SERVER_PORT"),
	}

	http.Handle("/", validate(cue, cue.Payload)) // Validate each task request and serve it

	log.Fatal(server.ListenAndServe())
}

func validate(h http.Handler, filter helpers.RequestResponseFilter) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := filter.Validate(w, r)
		if err != nil {
			return
		}

		h.ServeHTTP(w, r)
	})
}
