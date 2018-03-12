package models

type Endpoint struct {
	EndpointID int `db:"queue_endpoint_id"`
	QueueID    int
	Url        string
	Headers    []Header
}

type Header struct {
	Key   string
	Value string
}
