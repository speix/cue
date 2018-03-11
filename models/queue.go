package models

import (
	"errors"
)

type Queue struct {
	Name    string
	Mode    string
	Workers int
	Url     string
	Headers string
	Tasks   chan Task
}

type Queues map[string]*Queue

func (pool Queues) Add(reference string, queue *Queue) Queues {
	pool[reference] = queue
	return pool
}

func CreateQueue(name, mode, url, headers string, workers int) *Queue {
	return &Queue{
		Name:    name,
		Mode:    mode,
		Workers: workers,
		Url:     url,
		Headers: headers,
		Tasks:   make(chan Task, 100),
	}
}

func (db *DB) GetQueues() ([]*Queue, error) {

	sql := `select
				q.name,
				q.mode,
				q.workers,
				qe.url,
				jsonb_agg(
					jsonb_build_object(
						'key', qeh.key,
						'value', qeh.value
					)
				) as headers
			from queue q
				inner join queue_endpoint qe on qe.queue_id = q.queue_id
				inner join queue_endpoint_header qeh on qeh.queue_endpoint_id = qe.queue_endpoint_id
			group by q.name, q.mode, q.workers, qe.url`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	queues := make([]*Queue, 0)
	for rows.Next() {

		q := new(Queue)
		err := rows.Scan(&q.Name, &q.Mode, &q.Workers, &q.Url, &q.Headers)
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
