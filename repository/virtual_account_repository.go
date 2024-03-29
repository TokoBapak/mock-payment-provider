package repository

import (
	"context"
	"time"
)

type VirtualAccountRepository interface {
	// Migrate the database
	Migrate(ctx context.Context) error

	// CreateOrGetVirtualAccountNumber will accept the incoming customerUniqueField (can be anything ranging from
	// customer's unique ID to customer's email that's supposed to be unique). If the customerUniqueField haven't
	// been registered or submitted before, we will create a new virtual account number. Otherwise, we will
	// retrieve the virtual account number directly that's supposed to be exists.
	CreateOrGetVirtualAccountNumber(ctx context.Context, customerUniqueField string) (string, error)

	// CreateCharge create (or replace) the charged amount of the virtual account number.
	// If such virtual account number does not exist, it will create a new one.
	CreateCharge(ctx context.Context, virtualAccountNumber string, orderId string, amount int64, expiresAt time.Time) (account string, err error)

	// GetByVirtualAccountNumber acquires the current entry of the virtual account number.
	// It returns ErrNotFound if the entry was not found.
	GetByVirtualAccountNumber(ctx context.Context, virtualAccountNumber string) (Entry, error)

	// GetByOrderId acquires the current entry of the order id.
	// It returns ErrNotFound if the entry was not found
	GetByOrderId(ctx context.Context, orderId string) (Entry, error)

	// GetChargedAmount sees the amount that is charged to that specific virtual account number.
	// If the virtualAccountNumber does not exist, it will return an error of ErrNotFound
	GetChargedAmount(ctx context.Context, virtualAccountNumber string) (int64, error)

	// DeductCharge will free the virtual account number out of all charges. In other word,
	// it reset the charged amount to zero.
	DeductCharge(ctx context.Context, virtualAccountNumber string) error
}
