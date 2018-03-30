package http

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func NewServer() *http.Server {

	return &http.Server{
		Handler:      newRouter(),
		Addr:         ":" + os.Getenv("CUE_SERVER_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func newRouter() *mux.Router {

	router := mux.NewRouter()

	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc

		if route.MiddleWares != nil {
			for m := range route.MiddleWares {
				handler = route.MiddleWares[m](handler)
			}
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
