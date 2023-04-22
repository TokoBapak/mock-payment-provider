package transaction_service

import (
	"context"
	"errors"
	"fmt"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d *Dependency) Expire(ctx context.Context, orderId string) (business.ExpireResponse, error) {
	if orderId == "" {
		return business.ExpireResponse{}, fmt.Errorf("empty order id")
	}

	// Check for transaction status. If it's cancelled before or expired, then we shouldn't cancel it
	transactionStatus, err := d.TransactionRepository.GetByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return business.ExpireResponse{}, business.ErrTransactionNotFound
		}

		return business.ExpireResponse{}, fmt.Errorf("acquiring transaction status: %w", err)
	}

	// We can't expire any transaction that's not pending
	if transactionStatus.TransactionStatus != primitive.TransactionStatusPending {
		return business.ExpireResponse{}, business.ErrCannotModifyStatus
	}

	// Cancel the transaction
	err = d.TransactionRepository.UpdateStatus(ctx, orderId, primitive.TransactionStatusExpired)
	if err != nil {
		return business.ExpireResponse{}, fmt.Errorf("modifying the transaction status to canceled: %w", err)
	}

	return business.ExpireResponse{
		OrderId:           orderId,
		TransactionAmount: transactionStatus.TransactionAmount,
		PaymentType:       transactionStatus.PaymentType,
		TransactionStatus: primitive.TransactionStatusCanceled,
		TransactionTime:   transactionStatus.TransactionTime,
	}, nil
}
