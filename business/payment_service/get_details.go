package payment_service

import (
	"context"
	"errors"
	"fmt"

	"mock-payment-provider/business"
	"mock-payment-provider/repository"
)

func (d *Dependency) GetDetail(ctx context.Context, id string) (business.PaymentDetailsResponse, error) {
	// Set up an entry
	var entry repository.Entry
	var err error = nil

	// Try virtual account
	entry, err = d.virtualAccountRepository.GetByVirtualAccountNumber(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return business.PaymentDetailsResponse{}, fmt.Errorf("acquiring from virtual account store: %w", err)
		}

		// Try e-money
		entry, err = d.eMoneyRepository.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) || errors.Is(err, repository.ErrExpired) {
				return business.PaymentDetailsResponse{}, business.ErrTransactionNotFound
			}

			return business.PaymentDetailsResponse{}, fmt.Errorf("acquiring from emoney store: %w", err)
		}
	}

	// Acquire more data from transaction repository
	transaction, err := d.transactionRepository.GetByOrderId(ctx, entry.OrderId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return business.PaymentDetailsResponse{}, business.ErrTransactionNotFound
		}

		return business.PaymentDetailsResponse{}, fmt.Errorf("acquiring additional transaction information: %w", err)
	}

	return business.PaymentDetailsResponse{
		OrderId:              entry.OrderId,
		ChargedAmount:        transaction.TransactionAmount,
		Status:               transaction.TransactionStatus,
		PaymentMethod:        transaction.PaymentType,
		VirtualAccountNumber: entry.VirtualAccountNumber,
		EMoneyID:             entry.EMoneyID,
	}, nil
}
