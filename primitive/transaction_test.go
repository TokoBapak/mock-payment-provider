package primitive_test

import (
	"testing"
	"time"

	"mock-payment-provider/primitive"
)

func TestEntry_Expired(t *testing.T) {
	transaction1 := primitive.Transaction{ExpiresAt: time.Now().Add(time.Hour)}
	if transaction1.Expired() {
		t.Error("expecting transaction1 to not be expired, got expired")
	}

	transaction2 := primitive.Transaction{ExpiresAt: time.Now().Add(time.Hour * -1)}
	if !transaction2.Expired() {
		t.Error("expecting transaction2 to be expired, got not expired")
	}
}
