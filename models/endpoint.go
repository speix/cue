package models

import "time"

type Endpoint struct {
	EndpointID int `db:"queue_endpoint_id"`
	QueueID    int
	Url        string
	Timeout    time.Duration
	Retries    int
	Headers    []Header
}

type Header struct {
	Key   string
	Value string
}
