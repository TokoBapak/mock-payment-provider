package repository_test

import (
	"testing"
	"time"

	"mock-payment-provider/repository"
)

func TestEntry_Expired(t *testing.T) {
	entry1 := repository.Entry{ExpiresAt: time.Now().Add(time.Hour)}
	if entry1.Expired() {
		t.Error("expecting entry1 to not be expired, got expired")
	}

	entry2 := repository.Entry{ExpiresAt: time.Now().Add(time.Hour * -1)}
	if !entry2.Expired() {
		t.Error("expecting entry2 to be expired, got not expired")
	}
}
