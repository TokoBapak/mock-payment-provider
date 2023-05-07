package emoney

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewEmoneyRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return &Repository{}, fmt.Errorf("db is nil")
	}
	return &Repository{db: db}, nil
}
