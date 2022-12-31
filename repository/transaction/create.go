package transaction

import (
	"context"
	"time"

	"mock-payment-provider/primitive"
)

func (r *Repository) Create(ctx context.Context, orderId string, amount int64, paymentType primitive.PaymentType, status primitive.TransactionStatus, expiredAt time.Time) error {
	// TODO implement me
	panic("implement me")
}
