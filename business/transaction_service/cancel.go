package transaction_service

import (
	"context"
	"errors"
	"fmt"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d Dependency) Cancel(ctx context.Context, orderId string) (business.CancelResponse, error) {
	if orderId == "" {
		return business.CancelResponse{}, fmt.Errorf("empty order id")
	}

	// Check for transaction status. If it's cancelled before or expired, then we shouldn't cancel it
	transactionStatus, err := d.TransactionRepository.GetByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return business.CancelResponse{}, business.ErrTransactionNotFound
		}

		return business.CancelResponse{}, fmt.Errorf("acquiring transaction status: %w", err)
	}

	// We can't cancel any transaction that has been canceled, settled, or expired
	if transactionStatus.TransactionStatus == primitive.TransactionStatusExpired ||
		transactionStatus.TransactionStatus == primitive.TransactionStatusSettled ||
		transactionStatus.TransactionStatus == primitive.TransactionStatusCanceled {
		return business.CancelResponse{}, business.ErrCannotModifyStatus
	}

	// Cancel the transaction
	err = d.TransactionRepository.UpdateStatus(ctx, orderId, primitive.TransactionStatusCanceled)
	if err != nil {
		return business.CancelResponse{}, fmt.Errorf("modifying the transaction status to canceled: %w", err)
	}

	return business.CancelResponse{
		OrderId:           orderId,
		TransactionAmount: transactionStatus.TransactionAmount,
		PaymentType:       transactionStatus.PaymentType,
		TransactionStatus: primitive.TransactionStatusCanceled,
		TransactionTime:   transactionStatus.TransactionTime,
	}, nil
}
