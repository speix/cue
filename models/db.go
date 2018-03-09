package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetQueues() ([]*Queue, error)
}

type DB struct {
	*sql.DB
}

func NewDB(dataSource string) (*DB, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
