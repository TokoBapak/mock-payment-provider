package virtual_account_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"mock-payment-provider/repository/virtual_account"
)

func TestRepository_CreateCharge(t *testing.T) {
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Fatalf("creating virtual account repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty order ID", func(t *testing.T) {
		_, err := virtualAccountRepository.CreateCharge(ctx, "", "", 0, time.Now())
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if err.Error() != "orderId is empty" {
			t.Errorf("expecting an error of 'orderId is empty', instead got '%s'", err.Error())
		}
	})

	t.Run("Without valid VA number", func(t *testing.T) {
		_, err := virtualAccountRepository.CreateCharge(ctx, "12304560789", uuid.NewString(), 10000, time.Now().Add(time.Hour))
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})

	t.Run("With valid VA number", func(t *testing.T) {
		virtualAccountNumber, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, "annedoe@example.com")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		_, err = virtualAccountRepository.CreateCharge(ctx, virtualAccountNumber, uuid.NewString(), 10000, time.Now().Add(time.Hour))
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
