package payment_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d *Dependency) MarkAsPaid(ctx context.Context, orderId string, paymentMethod primitive.PaymentType, paymentId string) error {
	// Get transaction from order id
	transaction, err := d.transactionRepository.GetByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return business.ErrTransactionNotFound
		}

		return fmt.Errorf("acquiring transaction: %w", err)
	}

	// Check whether transaction is already expired
	if transaction.ExpiresAt.Before(time.Now()) {
		return business.ErrCannotModifyStatus
	}

	// Mark as settled
	err = d.transactionRepository.UpdateStatus(ctx, orderId, primitive.TransactionStatusSettled)
	if err != nil {
		return fmt.Errorf("updating transaction status: %w", err)
	}

	switch paymentMethod {
	case primitive.PaymentTypeVirtualAccountBCA:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBNI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBRI:
		// TODO: add VirtualAccountPermata
		err := d.virtualAccountRepository.DeductCharge(ctx, paymentId)
		if err != nil {
			return fmt.Errorf("deducting virtual account charge: %w", err)
		}
	case primitive.PaymentTypeEMoneyQRIS:
		fallthrough
	case primitive.PaymentTypeEMoneyGopay:
		fallthrough
	case primitive.PaymentTypeEMoneyShopeePay:
		err := d.eMoneyRepository.DeductCharge(ctx, paymentId)
		if err != nil {
			return fmt.Errorf("deducting emoney charge: %w", err)
		}
	default:
		return fmt.Errorf("invalid payment type")
	}

	// TODO: send a webhook, in background
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Printf("Recovered from panic: %v", e)
			}
		}()

		ctx := context.Background()

		err := d.webhookClient.Send(ctx, []byte{})
		if err != nil {
			// TODO: proper error logging
			log.Printf("Encountered an error during sending webhook: %s", err.Error())
		}
	}()

	return nil
}
