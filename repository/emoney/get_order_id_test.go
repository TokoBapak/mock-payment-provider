package emoney_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/emoney"
)

func TestRepository_GetByOrderId(t *testing.T) {
	emoneyRepository, err := emoney.NewEmoneyRepository(db)
	if err != nil {
		t.Fatalf("creating emoney repository: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("Empty OrderId", func(t *testing.T) {
		_, err := emoneyRepository.GetByOrderId(ctx, "")
		if err.Error() != "orderId is empty" {
			t.Errorf("expecting an error of 'orderId is empty', instead got %s", err.Error())
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := emoneyRepository.GetByOrderId(ctx, "not-exists")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expecting an error of repository.ErrNotFound, instead got %v", err)
		}
	})

	t.Run("Normal Integration", func(t *testing.T) {
		orderId := uuid.NewString()
		id, err := emoneyRepository.CreateCharge(ctx, orderId, 50000, time.Now().Add(time.Hour))
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		entry, err := emoneyRepository.GetByOrderId(ctx, orderId)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if entry.OrderId != orderId {
			t.Errorf("expecting orderId to be %s. instead got %s", orderId, entry.OrderId)
		}

		if entry.EMoneyID != id {
			t.Errorf("expecting emoney id to be %s, instead got %s", id, entry.EMoneyID)
		}

		if entry.ChargedAmount != 50000 {
			t.Errorf("expecting charged amount to be 50000 instead got %d", entry.ChargedAmount)
		}
	})
}
