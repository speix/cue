package app

import "time"

type Endpoint struct {
	EndpointID int `db:"queue_endpoint_id"`
	QueueID    int
	Url        string        // Request URL
	Timeout    time.Duration // Request Timeout
	Retries    int           // Number of retries on a failed request
	Method     string        // HTTP Method to be used for the request
	Headers    []Header      // Extra HTTP Headers for the request
}

type Header struct {
	Key   string
	Value string
}
