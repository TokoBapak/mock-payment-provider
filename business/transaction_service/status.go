package transaction_service

import (
	"context"

	"mock-payment-provider/business"
)

func (d Dependency) GetStatus(ctx context.Context, orderId string) (business.GetStatusResponse, error) {
	// TODO implement me
	panic("implement me")
}
