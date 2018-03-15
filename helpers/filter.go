package helpers

import (
	"net/http"
)

type RequestResponseFilter interface {
	Validate(w http.ResponseWriter, r *http.Request) error
}
