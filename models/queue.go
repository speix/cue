package models

type Queue struct {
	Name        string
	Mode        string
	Subscribers Endpoints
	Tasks       chan Task
}

type Queues map[string]*Queue

func (pool Queues) Add(reference string, queue *Queue) Queues {
	pool[reference] = queue
	return pool
}

func CreateQueue(name string, mode string) *Queue {
	return &Queue{
		Name:  name,
		Mode:  mode,
		Tasks: make(chan Task, 100),
	}
}

func (db *DB) GetQueues() ([]*Queue, error) {
	rows, err := db.Query("select q.name, q.mode from queue q")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	queues := make([]*Queue, 0)
	for rows.Next() {

		q := new(Queue)
		err := rows.Scan(&q.Name, &q.Mode)
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
