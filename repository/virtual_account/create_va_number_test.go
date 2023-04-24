package virtual_account_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/repository/virtual_account"
)

func TestRepository_CreateOrGetVirtualAccountNumber(t *testing.T) {
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Fatalf("creating virtual account repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty Unique ID", func(t *testing.T) {
		_, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, "")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if err.Error() != "customerUniqueField is empty" {
			t.Errorf("expecting error to be 'customerUniqueField is empty', instead got '%s'", err.Error())
		}
	})

	t.Run("Generate new", func(t *testing.T) {
		number, err := virtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, "johndoe@example.com")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if number == "" {
			t.Errorf("virtual account number is empty")
		}
	})
}
