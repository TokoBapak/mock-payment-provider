package repository

import (
	"context"
	"time"

	"mock-payment-provider/primitive"
)

type TransactionRepository interface {
	// Migrate the database
	Migrate(ctx context.Context) error
	// Create creates a new entry of transaction. If OrderId already exists,
	// it will return ErrDuplicate
	Create(ctx context.Context, orderId string, amount int64, paymentType primitive.PaymentType, status primitive.TransactionStatus, expiredAt time.Time) error
	// UpdateStatus will update the status. If the transaction has expired, it
	// will return ErrExpired. If the transaction was not found, it will return
	// ErrNotFound.
	UpdateStatus(ctx context.Context, orderId string, status primitive.TransactionStatus) error
	// GetByOrderId will get a transaction based on the order ID. It will return
	// ErrNotFound if the transaction can't be found.
	GetByOrderId(ctx context.Context, orderId string) (primitive.Transaction, error)
}
