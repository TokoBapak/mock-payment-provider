package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

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
}
