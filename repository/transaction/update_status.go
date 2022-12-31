package transaction

import (
	"context"

	"mock-payment-provider/primitive"
)

func (r *Repository) UpdateStatus(ctx context.Context, orderId string, status primitive.TransactionStatus) error {
	// TODO implement me
	panic("implement me")
}
