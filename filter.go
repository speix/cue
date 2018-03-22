package main

import "net/http"

type PayloadFilter interface {
	Validate(w http.ResponseWriter, r *http.Request) error
	QMap(queueName string) map[string]bool
}
