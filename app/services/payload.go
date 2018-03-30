package services

import (
	"github.com/speix/cue/app"
)

type PayloadService struct {
	payload app.Payload
}

func (s *PayloadService) QMap(queueName string) map[string]bool {
	s.payload.QMapper[queueName] = true
	return s.payload.QMapper
}
