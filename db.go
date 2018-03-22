package main

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateQueues() ([]*Queue, error)
}

type DB struct {
	*sqlx.DB
}

func NewDB(dataSource string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
