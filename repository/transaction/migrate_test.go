package transaction_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/repository/transaction"
)

func TestRepository_Migrate(t *testing.T) {
	repository, err := transaction.NewTransactionRepository(db)
	if err != nil {
		t.Fatalf("Creating transaction repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = repository.Migrate(ctx)
	if err != nil {
		t.Errorf("Migrating database: %s", err.Error())
	}
}
