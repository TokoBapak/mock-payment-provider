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
	// CheckPaidStatus checks whether an ID is paid
	CheckPaidStatus(ctx context.Context, id string) (paid bool, err error)
	// CancelCharge cancels a charge for the specified ID.
	CancelCharge(ctx context.Context, id string) error
	// DeductCharge will free the id of any charge and mark is as paid
	DeductCharge(ctx context.Context, id string) error
}
