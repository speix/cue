package app

import (
	"encoding/json"
)

type Payload struct {
	QueueName string          `json:"queue"`
	TaskName  string          `json:"task"`
	Messages  json.RawMessage `json:"messages"`
	Delay     int             `json:"delay"`
	QMapper   map[string]bool
}

type PayloadService interface {
	QMap(queueName string) map[string]bool
}
