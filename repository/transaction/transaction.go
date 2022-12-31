package transaction

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mock-payment-provider/primitive"
)

var ErrDuplicate = errors.New("duplicate")
var ErrNotFound = errors.New("not found")
var ErrExpired = errors.New("expired")

type Transaction struct {
	OrderId           string
	TransactionAmount int64
	PaymentType       primitive.PaymentType
	TransactionStatus primitive.TransactionStatus
	TransactionTime   time.Time
}

type ITransactionRepository interface {
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
	GetByOrderId(ctx context.Context, orderId string) (Transaction, error)
}

type Repository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return &Repository{}, errors.New("db is nil")
	}

	return &Repository{db: db}, nil
}
