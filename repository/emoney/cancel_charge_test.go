package emoney_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/repository/emoney"
)

func TestRepository_CancelCharge(t *testing.T) {
	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		t.Fatalf("creating emoney repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = emoneyRepository.CancelCharge(ctx, "any")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
