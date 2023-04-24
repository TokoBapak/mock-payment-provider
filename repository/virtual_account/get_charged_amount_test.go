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

func TestRepository_GetChargedAmount(t *testing.T) {
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Fatalf("creating virtual account repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty virtual account number", func(t *testing.T) {
		_, err := virtualAccountRepository.GetChargedAmount(context.Background(), "")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if err.Error() != "virtualAccountNumber is empty" {
			t.Errorf("expecting an error of 'virtualAccountNumber is empty', instead got %s", err.Error())
		}
	})

	t.Run("Not found, VA does not exists", func(t *testing.T) {
		_, err := virtualAccountRepository.GetChargedAmount(ctx, uuid.NewString())
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expecting repository.ErrNotFound, instead got %v", err)
		}
	})

	t.Run("Not found, entry does not exists", func(t *testing.T) {
		number, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, uuid.NewString())
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		_, err = virtualAccountRepository.GetChargedAmount(ctx, number)
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expecting repository.ErrNotFound, instead got %v", err)
		}
	})

	t.Run("Normal", func(t *testing.T) {
		number, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, uuid.NewString())
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		_, err = virtualAccountRepository.CreateCharge(ctx, number, uuid.NewString(), 10000, time.Now().Add(time.Hour))
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		chargedAmount, err := virtualAccountRepository.GetChargedAmount(ctx, number)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if chargedAmount != 10000 {
			t.Errorf("expecting charged amount to be '10000', instead got %d", chargedAmount)
		}
	})
}
