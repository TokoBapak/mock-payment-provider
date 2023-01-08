package transaction

import (
	"database/sql"
	"errors"
)

var ErrDuplicate = errors.New("duplicate")
var ErrNotFound = errors.New("not found")
var ErrExpired = errors.New("expired")

type Repository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return &Repository{}, errors.New("db is nil")
	}

	return &Repository{db: db}, nil
}
