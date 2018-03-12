package models

type Endpoint struct {
	EndpointID int
	Url        string
	Headers    []Header
}

type Header struct {
	Key   string
	Value string
}
