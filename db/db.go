package db

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

func New(filePath string) (*DB, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

type DB struct {
	*sql.DB
	sync.Mutex
}
