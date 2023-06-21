package payment_service_test

import (
	"context"
	"errors"
	"mock-payment-provider/business"
	"mock-payment-provider/business/payment_service"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/emoney"
	"mock-payment-provider/repository/transaction"
	"mock-payment-provider/repository/virtual_account"
	"mock-payment-provider/repository/webhook"
	"testing"
	"time"
)

func TestMarkAsPaid(t *testing.T) {
	ctx, setupCancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer setupCancel()

	transactionRepository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		t.Errorf("creating transaction repository: %s", err.Error())
	}
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Errorf("creating virtual account repository: %s", err.Error())
	}

	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		t.Errorf("creating emoney repository: %s", err.Error())
	}
	webhookClient, err := webhook.NewWebhookClient(cfg.webhookTargetURL)
	if err != nil {
		t.Errorf("creating webhook client: %s", err.Error())
	}

	paymentService, err := payment_service.NewPaymentService(payment_service.Config{
		ServerKey:                cfg.serverKey,
		TransactionRepository:    transactionRepository,
		EMoneyRepository:         emoneyRepository,
		WebhookClient:            webhookClient,
		VirtualAccountRepository: virtualAccountRepository,
	})
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	t.Run("MarkAsPaid should return not found if the order id is empty", func(t *testing.T) {
		err = paymentService.MarkAsPaid(ctx, "", primitive.PaymentTypeUnspecified)
		if err == nil {
			t.Errorf("expecting error to be not nil, but got nil")
		}
	})
	t.Run("MarkAsPaid should return err 'not found' if the order id is not found", func(t *testing.T) {
		err = paymentService.MarkAsPaid(ctx, "not-exist", primitive.PaymentTypeUnspecified)
		if err == nil {
			t.Errorf("expecting error to be not nil, but got nil")
		}
		if !errors.Is(err, business.ErrTransactionNotFound) {
			t.Errorf("expecting error %s, instead got %v", business.ErrTransactionNotFound, err)
		}
	})

	t.Run("MarkAsPaid should return err 'cannot modify status' if the transaction is already expired", func(t *testing.T) {
		orderId := "order-id-expired"
		err = transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     orderId,
			Amount:      50000,
			PaymentType: primitive.PaymentTypeEMoneyQRIS,
			Status:      primitive.TransactionStatusExpired,
			ExpiredAt:   time.Now().Add(-time.Minute),
		})
		err = paymentService.MarkAsPaid(ctx, orderId, primitive.PaymentTypeEMoneyQRIS)
		if err == nil {
			t.Errorf("expecting error to be not nil, but got nil")
		}
		if !errors.Is(err, business.ErrCannotModifyStatus) {
			t.Errorf("expecting error %s, instead got %v", business.ErrCannotModifyStatus, err)
		}

	})

	t.Run("MarkAsPaid should return err 'cannot modify status' if the previous status is not pending", func(t *testing.T) {
		orderId := "order-id-2"
		err = transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     orderId,
			Amount:      50000,
			PaymentType: primitive.PaymentTypeEMoneyQRIS,
			Status:      primitive.TransactionStatusSettled,
			ExpiredAt:   time.Now().Add(time.Hour),
		})
		err = paymentService.MarkAsPaid(ctx, orderId, primitive.PaymentTypeEMoneyQRIS)
		if err == nil {
			t.Errorf("expecting error to be not nil, but got nil")
		}
		if !errors.Is(err, business.ErrCannotModifyStatus) {
			t.Errorf("expecting error %s, instead got %v", business.ErrCannotModifyStatus, err)
		}
	})

	t.Run("MarkAsPaid should return err if the payment method is not supported", func(t *testing.T) {
		orderId := "order-id-not-supported"
		err = transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     orderId,
			Amount:      50000,
			PaymentType: primitive.PaymentTypeEMoneyQRIS,
			Status:      primitive.TransactionStatusPending,
			ExpiredAt:   time.Now().Add(time.Hour),
		})
		err = paymentService.MarkAsPaid(ctx, orderId, primitive.PaymentTypeEMoneyQRIS)
		if err != nil {
			t.Errorf("expecting error to be nil, but got %v", err)
		}
	})

	t.Run("MarkAsPaid should return err when use PaymentTypeVirtualAccountPermata", func(t *testing.T) {
		orderId := "order-id-3"
		err = transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     orderId,
			Amount:      50000,
			PaymentType: primitive.PaymentTypeVirtualAccountPermata,
			Status:      primitive.TransactionStatusPending,
			ExpiredAt:   time.Now().Add(time.Hour),
		})
		err = paymentService.MarkAsPaid(ctx, orderId, primitive.PaymentTypeVirtualAccountPermata)
		if err != nil {
			t.Errorf("expecting error to be nil, but got %v", err)
		}
	})
}
