package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"
)

type Task struct {
	Name     string          // Name of the task
	Messages []*bytes.Buffer // Messages to be forwarded to the next peer
	Delay    time.Duration   // Run after x number of seconds
	Retries  int             // Number of retries for each task on run time
}

func CreateTask(name string, payload json.RawMessage, delay time.Duration) (*Task, error) {

	task := &Task{
		Name:    name,
		Delay:   delay,
		Retries: 0,
	}

	err := task.convertMessages(payload)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (task *Task) convertMessages(payload json.RawMessage) error {
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
		task.Messages = append(task.Messages, bytes.NewBufferString(string(jsonString)))
	}

	return nil
}
