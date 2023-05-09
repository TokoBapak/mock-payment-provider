package virtual_account_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/virtual_account"
)

func TestRepository_GetByOrderId(t *testing.T) {
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Fatalf("creating virtual account repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty OrderID", func(t *testing.T) {
		_, err := virtualAccountRepository.GetByOrderId(ctx, "")
		if err.Error() != "orderId is empty" {
			t.Errorf("expecting an error of 'orderId is empty', instead got %s", err.Error())
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := virtualAccountRepository.GetByOrderId(ctx, "not-exists")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expecting an error of repository.ErrNotFound, instead got %v", err)
		}
	})

	t.Run("Happy Integration", func(t *testing.T) {
		vaNumber, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, "annedoe@example.com")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		orderId := uuid.NewString()
		_, err = virtualAccountRepository.CreateCharge(ctx, vaNumber, orderId, 50000, time.Now().Add(time.Hour))
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		transaction, err := virtualAccountRepository.GetByOrderId(ctx, orderId)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if transaction.OrderId != orderId {
			t.Errorf("expecting orderId to be %s, instead got %s", orderId, transaction.OrderId)
		}

		if transaction.VirtualAccountNumber != vaNumber {
			t.Errorf("expecting virtualAccountNumber to be %s, instead got %s", vaNumber, transaction.VirtualAccountNumber)
		}

		if transaction.ChargedAmount != 50000 {
			t.Errorf("expecting ChargedAmount to be %d, instead got %d", 50000, transaction.ChargedAmount)
		}
	})
}
