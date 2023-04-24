package virtual_account_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/repository/virtual_account"
)

func TestRepository_Migrate(t *testing.T) {
	virtualAccountRepository, err := virtual_account.NewVirtualAccountRepository(db)
	if err != nil {
		t.Fatalf("creating virtual account repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = virtualAccountRepository.Migrate(ctx)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
