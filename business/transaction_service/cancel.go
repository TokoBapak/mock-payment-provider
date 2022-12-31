package transaction_service

import (
	"context"

	"mock-payment-provider/business"
)

func (d Dependency) Cancel(ctx context.Context, orderId string) (business.CancelResponse, error) {
	// TODO implement me
	panic("implement me")
}
