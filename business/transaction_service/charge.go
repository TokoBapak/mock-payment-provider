package transaction_service

import (
	"context"

	"mock-payment-provider/business"
)

func (d Dependency) Charge(ctx context.Context, request business.ChargeRequest) (business.ChargeResponse, error) {
	// TODO implement me
	panic("implement me")
}
