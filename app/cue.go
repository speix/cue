package app

import "github.com/jmoiron/sqlx"

type Cue struct {
	*sqlx.DB
}

type CueService interface {
	LoadQueues() []*Queue
	StartCue()
}
