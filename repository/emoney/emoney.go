package emoney

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mock-payment-provider/repository"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Migrate(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (id string, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetByID(ctx context.Context, id string) (repository.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) CheckPaidStatus(ctx context.Context, id string) (paid bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) CancelCharge(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) DeductCharge(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewEmoneyRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return &Repository{}, fmt.Errorf("db is nil")
	}
	return &Repository{db: db}, nil
}
