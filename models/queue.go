package models

type Queue struct {
	Name  string
	Mode  string
	Tasks chan Task
}

func NewQueue(name string, mode string) *Queue {
	return &Queue{
		Name:  name,
		Mode:  mode,
		Tasks: make(chan Task, 100),
	}
}
