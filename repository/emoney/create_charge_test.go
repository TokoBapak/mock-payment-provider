package emoney_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"mock-payment-provider/repository/emoney"
)

func TestRepository_CreateCharge(t *testing.T) {
	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		t.Fatalf("creating emoney repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Happy", func(t *testing.T) {
		orderId := uuid.NewString()

		id, err := emoneyRepository.CreateCharge(ctx, orderId, 12345, time.Now().Add(time.Minute))
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if id == "" {
			t.Errorf("expecting id to be not empty")
		}
	})
}
