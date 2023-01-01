package repository

import (
	"context"
	"time"
)

type VirtualAccountRepository interface {
	// Migrate the database
	Migrate(ctx context.Context) error
	// CreateCharge create (or replace) the charged amount of the virtual account number.
	// If such virtual account number does not exist, it will create a new one.
	CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (account string, err error)
	// GetChargedAmount sees the amount that is charged to that specific virtual account number.
	// If the virtualAccountNumber does not exists, it will return an error of ErrNotFound
	GetChargedAmount(ctx context.Context, virtualAccountNumber string) (int64, error)
	// DeductCharge will free the virtual account number out of all charges. In other word,
	// it reset the charged amount to zero.
	DeductCharge(ctx context.Context, virtualAccountNumber string) error
}
