package transaction

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return &Repository{}, errors.New("db is nil")
	}

	return &Repository{db: db}, nil
}
