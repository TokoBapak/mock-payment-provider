package repository

import (
	"context"
	"time"
)

type EMoneyRepository interface {
	// Migrate the database
	Migrate(ctx context.Context) error

	// CreateCharge saves the charge request and create a new unique ID. This unique ID
	// will be used as the ID to do things e-money related.
	CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (id string, err error)

	// GetByID acquires the current entry of the specified ID.
	// It returns ErrNotFound if the entry was not found.
	// It returns ErrExpired if the ID is expired
	GetByID(ctx context.Context, id string) (Entry, error)

	// GetByOrderId acquires the current entry of the order ID.
	// It returns ErrNotFound if the entry was not found.
	// It returns ErrExpired if the ID is expired
	GetByOrderId(ctx context.Context, orderId string) (Entry, error)

	// CancelCharge cancels a charge for the specified ID.
	CancelCharge(ctx context.Context, orderId string) error

	// DeductCharge will free the id of any charge and mark is as paid
	DeductCharge(ctx context.Context, orderId string) error
}
