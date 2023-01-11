package transaction

import (
	"context"

	"mock-payment-provider/primitive"
)

func (r *Repository) GetByOrderId(ctx context.Context, orderId string) (primitive.Transaction, error) {
	// TODO implement me
	panic("implement me")
}
