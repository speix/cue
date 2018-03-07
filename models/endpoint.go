package models

type Endpoint struct {
	Url     string
	Headers Headers
}

type Endpoints []Endpoint

type Header struct {
	Key   string
	Value string
}

type Headers []Header
