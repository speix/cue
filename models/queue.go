package models

type Queue struct {
	Name        string
	Mode        string
	Workers     int
	Subscribers Endpoints
	Tasks       chan Task
}

type Queues map[string]*Queue

func (pool Queues) Add(reference string, queue *Queue) Queues {
	pool[reference] = queue
	return pool
}

func CreateQueue(name, mode string, workers int) *Queue {
	return &Queue{
		Name:    name,
		Mode:    mode,
		Workers: workers,
		Tasks:   make(chan Task, 100),
	}
}

func (db *DB) GetQueues() ([]*Queue, error) {
	rows, err := db.Query("select q.name, q.mode, q.workers from queue q")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	queues := make([]*Queue, 0)
	for rows.Next() {

		q := new(Queue)
		err := rows.Scan(&q.Name, &q.Mode, &q.Workers)
		if err != nil {
			return nil, err
		}

		queues = append(queues, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return queues, nil
}
