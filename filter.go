package main

import "net/http"

type RequestResponseFilter interface {
	Validate(w http.ResponseWriter, r *http.Request) error
	QMap(queueName string) map[string]bool
}
