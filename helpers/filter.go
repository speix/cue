package helpers

import (
	"net/http"

	"github.com/speix/cue/models"
)

type RequestResponseFilter interface {
	Validate(w http.ResponseWriter, r *http.Request, queues models.Queues) error
}
