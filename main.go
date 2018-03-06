package main

import (
	"log"
	"net/http"

	"github.com/speix/cue/handlers"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RequestHandler(w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
