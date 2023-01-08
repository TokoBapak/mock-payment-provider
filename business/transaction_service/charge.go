package transaction_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d Dependency) Charge(ctx context.Context, request business.ChargeRequest) (business.ChargeResponse, error) {
	// Validate the transaction amount and the amount of each product items
	var totalAmount int64
	for _, item := range request.ProductItems {
		totalAmount += item.Price * item.Quantity
	}

	if totalAmount != request.TransactionAmount {
		return business.ChargeResponse{}, business.ErrMismatchedTransactionAmount
	}

	switch request.PaymentType {
	case primitive.PaymentTypeVirtualAccountBCA:
		fallthrough
	case primitive.PaymentTypeVirtualAccountMandiri:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBRI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBNI:
		// Create new transaction
		err := d.TransactionRepository.Create(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			request.PaymentType,
			primitive.TransactionStatusPending,
			time.Now().Add(time.Hour*24),
		)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return business.ChargeResponse{}, business.ErrDuplicateOrderId
			}

			return business.ChargeResponse{}, fmt.Errorf("creating new transaction: %w", err)
		}

		// Create a virtual account entry
		virtualAccountNumber, err := d.VirtualAccountRepository.CreateCharge(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			time.Now().Add(time.Hour*24),
		)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("creating virtual account entry: %w", err)
		}

		return business.ChargeResponse{
			OrderId:           request.OrderId,
			TransactionAmount: request.TransactionAmount,
			PaymentType:       request.PaymentType,
			TransactionStatus: primitive.TransactionStatusPending,
			TransactionTime:   time.Now(),
			EMoneyAction:      []business.EMoneyAction{},
			VirtualAccountAction: business.VirtualAccountAction{
				Bank:                 request.PaymentType.String(),
				VirtualAccountNumber: virtualAccountNumber,
			},
		}, nil
	case primitive.PaymentTypeEMoneyQRIS:
		fallthrough
	case primitive.PaymentTypeEMoneyGopay:
		fallthrough
	case primitive.PaymentTypeEMoneyShopeePay:
		// Create new transaction
		err := d.TransactionRepository.Create(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			request.PaymentType,
			primitive.TransactionStatusPending,
			time.Now().Add(time.Hour*3),
		)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return business.ChargeResponse{}, business.ErrDuplicateOrderId
			}

			return business.ChargeResponse{}, fmt.Errorf("creating new transaction: %w", err)
		}

		// Create e-money entry
		id, err := d.EMoneyRepository.CreateCharge(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			time.Now().Add(time.Hour*3),
		)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("creating e-money entry: %w", err)
		}

		return business.ChargeResponse{
			OrderId:           request.OrderId,
			TransactionAmount: request.TransactionAmount,
			PaymentType:       request.PaymentType,
			TransactionStatus: primitive.TransactionStatusPending,
			TransactionTime:   time.Now(),
			EMoneyAction: []business.EMoneyAction{
				{
					EMoneyActionType: business.EMoneyActionTypePay,
					Method:           "GET",
					URL:              "/e-money/" + id + "/pay",
				},
				{
					EMoneyActionType: business.EMoneyActionTypeStatus,
					Method:           "GET",
					URL:              "/e-money/" + id + "/status",
				},
				{
					EMoneyActionType: business.EMoneyActionTypeCancel,
					Method:           "POST",
					URL:              "/e-money/" + id + "/cancel",
				},
			},
			VirtualAccountAction: business.VirtualAccountAction{},
		}, nil
	case primitive.PaymentTypeUnspecified:
		fallthrough
	default:
		return business.ChargeResponse{}, fmt.Errorf("invalid payment type")
	}
}
