package main

import (
	"log"

	"github.com/speix/cue/app/http"
)

func main() {

	server := http.NewServer()

	log.Fatal(server.ListenAndServe())

	/*cue := StartCue() // Start the Cue (Queues, Dispatchers, Workers, Listeners)

	server := &http.Server{
		Addr: ":" + os.Getenv("CUE_SERVER_PORT"),
	}

	http.Handle("/", validate(cue, cue.Payload)) // Validate each task request and serve it

	log.Fatal(server.ListenAndServe())*/
}

/*func validate(h http.Handler, filter PayloadFilter) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := filter.Validate(w, r)
		if err != nil {
			return
		}

		h.ServeHTTP(w, r)
	})
}*/
