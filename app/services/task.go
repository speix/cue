package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/speix/cue/app"
)

type TaskService struct {
	task app.Task
}

func (s *TaskService) CreateTask(name string, payload json.RawMessage, delay time.Duration) (*app.Task, error) {
	task := &app.Task{
		Name:    name,
		Delay:   delay,
		Retries: 0,
	}

	err := s.convertMessages(payload)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) convertMessages(payload json.RawMessage) error {
	var pMessages []interface{}

	err := json.Unmarshal(payload, &pMessages)
	if err != nil {
		return errors.New("messages must be a valid json array")
	}

	for m := range pMessages {
		jsonString, err := json.Marshal(pMessages[m])
		if err != nil {
			return err
		}
		s.task.Messages = append(s.task.Messages, bytes.NewBufferString(string(jsonString)))
	}

	return nil
}
