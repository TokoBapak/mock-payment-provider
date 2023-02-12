package transaction_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/transaction"
)

func TestRepository_UpdateStatus(t *testing.T) {
	transactionRepository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		t.Fatalf("Creating transaction repository: %s", err.Error())
	}

	t.Run("Empty Order ID", func(t *testing.T) {
		err := transactionRepository.UpdateStatus(context.Background(), "", primitive.TransactionStatusDenied)
		if err == nil {
			t.Errorf("expecting an error, got nil instead")
		}

		if err.Error() != "empty order id" {
			t.Errorf("expecting an error of empty order id, instead got %s", err.Error())
		}
	})

	t.Run("Happy Case", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		err := transactionRepository.Create(ctx, repository.CreateTransactionParam{
			OrderID:     "d41d8cd98f00b204e9800998ecf8427e",
			Amount:      100_000,
			PaymentType: primitive.PaymentTypeEMoneyQRIS,
			Status:      primitive.TransactionStatusPending,
			ExpiredAt:   time.Now().Add(time.Hour),
		})
		if err != nil {
			t.Fatalf("creating an entry: %s", err.Error())
		}

		err = transactionRepository.UpdateStatus(ctx, "d41d8cd98f00b204e9800998ecf8427e", primitive.TransactionStatusSettled)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
