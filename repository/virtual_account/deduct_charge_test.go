package virtual_account_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"mock-payment-provider/repository/virtual_account"
)

func TestRepository_DeductCharge(t *testing.T) {
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Fatalf("creating virtual account repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty virtual account number", func(t *testing.T) {
		err := virtualAccountRepository.DeductCharge(ctx, "")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if err.Error() != "empty virtual account number" {
			t.Errorf("expecting an error of 'empty virtual account number', got %s instead", err.Error())
		}
	})

	t.Run("Random", func(t *testing.T) {
		err := virtualAccountRepository.DeductCharge(ctx, uuid.NewString())
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
