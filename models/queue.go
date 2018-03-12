package models

import (
	"errors"
)

type Queue struct {
	QueueID  int `db:"queue_id"`
	Name     string
	Mode     string
	Workers  int
	Endpoint Endpoint
	Tasks    chan Task
}

type Queues map[string]*Queue

func (pool Queues) Add(reference string, queue *Queue) Queues {
	pool[reference] = queue
	return pool
}

func CreateQueue(queue *Queue) *Queue {
	return &Queue{
		Name:     queue.Name,
		Mode:     queue.Mode,
		Workers:  queue.Workers,
		Endpoint: queue.Endpoint,
		Tasks:    make(chan Task, 100),
	}
}

func (db *DB) GetQueues() ([]*Queue, error) {

	queues := make([]*Queue, 0)
	err := db.Select(&queues, "select * from queue")
	if err != nil {
		return nil, err
	}

	for q := range queues {

		endpoint := Endpoint{}
		err = db.Get(&endpoint, "select queue_endpoint_id, url from queue_endpoint where queue_id=$1", queues[q].QueueID)
		if err != nil {
			return nil, err
		}

		headers := make([]Header, 0)
		err = db.Select(&headers, "select key, value from queue_endpoint_header where queue_endpoint_id=$1", endpoint.EndpointID)
		if err != nil {
			return nil, err
		}

		endpoint.Headers = headers
		queues[q].Endpoint = endpoint
	}

	defer db.Close()

	if len(queues) == 0 {
		return nil, errors.New("No queues found in database ")
	}

	return queues, nil
}
