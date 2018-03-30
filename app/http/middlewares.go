package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/speix/cue/app/helpers"

	"github.com/speix/cue/app"
)

type ServiceResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func validatePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Hit the validation middleware")

		w.Header().Set("Content-Type", "application/json")
		response := ServiceResponse{}
		p := app.Payload{}

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			response.Error = true
			response.Message = err.Error()
			responseJson, _ := json.Marshal(response)

			w.WriteHeader(400)
			w.Write(responseJson)

			return
		}

		// check if there is a QueueName request
		if len(p.QueueName) == 0 {
			response.Error = true
			response.Message = "Queue is not set"
			responseJson, _ := json.Marshal(response)

			w.WriteHeader(400)
			w.Write(responseJson)

			return
		}

		// check if there is a Payload request
		if len(p.Messages) == 0 {
			response.Error = true
			response.Message = "Messages is not set"
			responseJson, _ := json.Marshal(response)

			w.WriteHeader(400)
			w.Write(responseJson)

			return
		}

		// check if delay is between 0 to 30 minutes
		if !helpers.InBetween(p.Delay, 0, 1800) {

			response.Error = true
			response.Message = "Delay must be in seconds between 0 and 1800."
			responseJson, _ := json.Marshal(response)

			w.WriteHeader(400)
			w.Write(responseJson)

			return
		}

		// check if payload attribute is a json field
		if !helpers.IsJSON(p.Messages) {
			response.Error = true
			response.Message = "Message format is not json"
			responseJson, _ := json.Marshal(response)

			w.WriteHeader(400)
			w.Write(responseJson)

			return
		}

		// check if requested queue name is available
		if !p.QMapper[p.QueueName] {
			response.Error = true
			response.Message = "Queue " + p.QueueName + " not found"
			responseJson, _ := json.Marshal(response)

			w.WriteHeader(404)
			w.Write(responseJson)

			return
		}

		next.ServeHTTP(w, r)
	})
}
