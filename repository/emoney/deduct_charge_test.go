package emoney_test

import (
	"context"
	"testing"
	"time"

	"mock-payment-provider/repository/emoney"
)

func TestRepository_DeductCharge(t *testing.T) {
	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		t.Fatalf("creating emoney repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty OrderId", func(t *testing.T) {
		err := emoneyRepository.DeductCharge(ctx, "")
		if err.Error() != "orderId is empty" {
			t.Errorf("expecting an error of 'orderId is empty', instead got %s", err.Error())
		}
	})

	t.Run("Normal", func(t *testing.T) {
		err := emoneyRepository.DeductCharge(ctx, "any")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
