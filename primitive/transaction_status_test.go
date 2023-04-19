package primitive_test

import (
	"testing"

	"mock-payment-provider/primitive"
)

func TestTransactionStatus_String(t *testing.T) {
	t.Run("TransactionStatusPending", func(t *testing.T) {
		if primitive.TransactionStatusPending.String() != "pending" {
			t.Errorf("expecting TransactionStatusPending.String() to be 'pending', instead got %s", primitive.TransactionStatusPending.String())
		}
	})

	t.Run("TransactionStatusDenied", func(t *testing.T) {
		if primitive.TransactionStatusDenied.String() != "denied" {
			t.Errorf("expecting TransactionStatusDenied.String() to be 'denied', instead got %s", primitive.TransactionStatusDenied.String())
		}
	})

	t.Run("TransactionStatusSettled", func(t *testing.T) {
		if primitive.TransactionStatusSettled.String() != "settled" {
			t.Errorf("expecting TransactionStatusSettled.String() to be 'settled', instead got %s", primitive.TransactionStatusSettled.String())
		}
	})

	t.Run("TransactionStatusExpired", func(t *testing.T) {
		if primitive.TransactionStatusExpired.String() != "expired" {
			t.Errorf("expecting TransactionStatusExpired.String() to be 'expired', instead got %s", primitive.TransactionStatusExpired.String())
		}
	})

	t.Run("TransactionStatusUnspecified", func(t *testing.T) {
		if primitive.TransactionStatusUnspecified.String() != "UNSPECIFIED" {
			t.Errorf("expecting TransactionStatusUnspecified.String() to be 'UNSPECIFIED', instead got %s", primitive.TransactionStatusUnspecified.String())
		}
	})
}
