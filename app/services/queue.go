package services

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/speix/cue/app"
)

type QueueService struct {
	queue app.Queue
	pool  app.QueuesPool
	DB    *sqlx.DB
}

func (s *QueueService) Add(queue *app.Queue) app.QueuesPool {
	s.pool[queue.Name] = queue
	return s.pool
}

func (s *QueueService) CreateQueues() ([]*app.Queue, error) {

	queues := make([]*app.Queue, 0)

	err := s.DB.Select(&queues, "select * from queue")
	if err != nil {
		return nil, err
	}

	for q := range queues {

		endpoint := app.Endpoint{}
		err = s.DB.Get(&endpoint, "select queue_endpoint_id, url, timeout, retries, method from queue_endpoint where queue_id=$1", queues[q].QueueID)
		if err != nil {
			return nil, err
		}

		headers := make([]app.Header, 0)
		err = s.DB.Select(&headers, "select key, value from queue_endpoint_header where queue_endpoint_id=$1", endpoint.EndpointID)
		if err != nil {
			return nil, err
		}

		endpoint.Headers = headers
		queues[q].Endpoint = endpoint
		queues[q].Tasks = make(chan app.Task, 100)
	}

	defer s.DB.Close()

	if len(queues) == 0 {
		return nil, errors.New("No queues found in database ")
	}

	return queues, nil

}
