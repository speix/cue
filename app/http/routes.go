package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	MiddleWares []mux.MiddlewareFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Add new task",
		"POST",
		"/",
		Handler{}.TaskRequestHandler,
		[]mux.MiddlewareFunc{validatePayload},
	},
}
