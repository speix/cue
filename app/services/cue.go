package services

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/speix/cue/app"
)

type CueService struct {
	cue          app.Cue
	queueService QueueService
	payload      PayloadService
	dispatcher   app.DispatcherService
}

func (s *CueService) StartCue() {

	queues := s.LoadQueues() // Load queues from database

	/*handler := &TaskRequestHandler{
		Payload: &Payload{
			QMapper: make(map[string]bool),
		},
		Pool: QueuesPool{},
	}*/

	for i := range queues {

		s.queueService.Add(queues[i]) // Add queue to the pool of queues

		s.payload.QMap(queues[i].Name) // Add available queue names to the Payload as reference

		dispatcher := s.dispatcher.CreateDispatcher(queues[i].Workers) // Create a dispatcher for each queue

		s.dispatcher.Start()

		dispatcher.Start(queues[i]) // Start workers running on each queue

		dispatcher.Listen() // Listen for tasks
	}

	return handler

}

func (s *CueService) LoadQueues() []*app.Queue {

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("CUE_DB_HOST"), os.Getenv("CUE_DB_USER"), os.Getenv("CUE_DB_PASS"), os.Getenv("CUE_DB_NAME"))

	db, err := NewDB(dataSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	s.queueService.DB = db.DB

	queues, err := s.queueService.CreateQueues()
	if err != nil {
		log.Fatal(err.Error())
	}

	return queues
}

func NewDB(dataSource string) (*app.Cue, error) {
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	return &app.Cue{DB: db}, nil
}
