package transaction_service

import (
	"context"
	"errors"
	"fmt"

	"mock-payment-provider/business"
	transaction_repository "mock-payment-provider/repository/transaction"
)

func (d Dependency) GetStatus(ctx context.Context, orderId string) (business.GetStatusResponse, error) {
	if orderId == "" {
		return business.GetStatusResponse{}, fmt.Errorf("empty order id")
	}

	transaction, err := d.TransactionRepository.GetByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, transaction_repository.ErrNotFound) {
			return business.GetStatusResponse{}, business.ErrTransactionNotFound
		}

		return business.GetStatusResponse{}, fmt.Errorf("acquiring transaction by order id: %w", err)
	}

	return business.GetStatusResponse{
		OrderId:           orderId,
		TransactionStatus: transaction.TransactionStatus,
	}, nil
}
