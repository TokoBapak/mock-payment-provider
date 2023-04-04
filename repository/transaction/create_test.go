package transaction_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/repository"
	"mock-payment-provider/repository/transaction"
)

func TestRepository_Create(t *testing.T) {
	transactionRepository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		t.Fatalf("Creating transaction repository: %s", err.Error())
	}

	t.Run("Happy Case", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     "A",
			Amount:      10_000,
			PaymentType: 2,
			Status:      1,
			ExpiredAt:   time.Now().Add(time.Hour * 3),
		})
		if err != nil {
			t.Errorf("Creating entry to transaction log: %s", err.Error())
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     "3851b601-686d-4f29-8ff0-ac13d34f516c",
			Amount:      10_000,
			PaymentType: 1,
			Status:      1,
			ExpiredAt:   time.Now().Add(time.Hour * 3),
		})
		if err != nil {
			t.Errorf("Creating entry to transaction log: %s", err.Error())
		}

		err = transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     "3851b601-686d-4f29-8ff0-ac13d34f516c",
			Amount:      10_000,
			PaymentType: 1,
			Status:      1,
			ExpiredAt:   time.Now().Add(time.Hour * 3),
		})
		if err == nil {
			t.Errorf("Expecting an error, got nil instead")
		}

		t.Logf(err.Error())
	})
}
