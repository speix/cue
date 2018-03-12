package models

import (
	"errors"
)

type Queue struct {
	QueueID   int
	Name      string
	Mode      string
	Workers   int
	Endpoints []Endpoint
	Headers   string
	Tasks     chan Task
}

type Queues map[string]*Queue

func (pool Queues) Add(reference string, queue *Queue) Queues {
	pool[reference] = queue
	return pool
}

func CreateQueue(name, mode, headers string, workers int) *Queue {
	return &Queue{
		Name:    name,
		Mode:    mode,
		Workers: workers,
		Headers: headers,
		Tasks:   make(chan Task, 100),
	}
}

func (db *DB) GetQueues() ([]*Queue, error) {

	sql := "select q.queue_id, q.name, q.mode, q.workers from queue q"

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	queues := make([]*Queue, 0)
	for rows.Next() {

		q := new(Queue)
		err := rows.Scan(&q.QueueID, &q.Name, &q.Mode, &q.Workers)
		if err != nil {
			return nil, err
		}

		queues = append(queues, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(queues) == 0 {
		return nil, errors.New("No queues found in database ")
	}

	return queues, nil
}
