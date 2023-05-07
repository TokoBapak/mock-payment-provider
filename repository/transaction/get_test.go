package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/transaction"
)

func TestRepository_GetByOrderId(t *testing.T) {
	transactionRepository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		t.Fatalf("Creating new transaction repository: %s", err.Error())
	}

	t.Run("Empty Order ID", func(t *testing.T) {
		_, err := transactionRepository.GetByOrderId(context.Background(), "")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if err.Error() != "empty order id" {
			t.Errorf("expecting .Error() to be 'empty order id', instead got %s", err.Error())
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err := transactionRepository.GetByOrderId(ctx, "NOT-FOUND")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expecting err to be repository.ErrNotFound, instead got %s", err.Error())
		}
	})

	t.Run("Happy Integration", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		orderId := uuid.NewString()
		expiredAt := time.Now().Add(time.Hour * 3)

		err := transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     orderId,
			Amount:      10_000,
			PaymentType: 2,
			Status:      1,
			ExpiredAt:   expiredAt,
		})
		if err != nil {
			t.Errorf("Creating entry to transaction log: %s", err.Error())
		}

		entry, err := transactionRepository.GetByOrderId(ctx, orderId)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if entry.OrderId != orderId {
			t.Errorf("expecting orderId to be %s, instead got %s", orderId, entry.OrderId)
		}

		if entry.PaymentType != 2 {
			t.Errorf("expecting payment type to be 2, instead got %d", entry.PaymentType)
		}

		if entry.TransactionStatus != 1 {
			t.Errorf("expecting transaction status to be 1, instead got %d", entry.TransactionStatus)
		}

		if entry.TransactionAmount != 10_000 {
			t.Errorf("expecting transaction amount to be 10_000, instead got %d", entry.TransactionAmount)
		}

		if !entry.ExpiresAt.Equal(expiredAt) {
			t.Errorf("expecting expires at to be %s, instead got %s", expiredAt.String(), entry.ExpiresAt.String())
		}
	})
}
