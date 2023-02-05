package primitive_test

import (
	"testing"

	"mock-payment-provider/primitive"
)

func TestCurrency_String(t *testing.T) {
	t.Run("IDR", func(t *testing.T) {
		if primitive.CurrencyIDR.String() != "IDR" {
			t.Errorf("expecting CurrencyIDR.String() to be 'IDR', instead got %s", primitive.CurrencyIDR.String())
		}
	})

	t.Run("USD", func(t *testing.T) {
		if primitive.CurrencyUSD.String() != "USD" {
			t.Errorf("expecting CurrencyUSD.String() to be 'USD', instead got %s", primitive.CurrencyUSD.String())
		}
	})

	t.Run("UNSPECIFIED", func(t *testing.T) {
		if primitive.CurrencyUnspecified.String() != "UNSPECIFIED" {
			t.Errorf("expecting CurrencyUnspecified.String() to be 'UNSPECIFIED', instead got %s", primitive.CurrencyUnspecified.String())
		}
	})
}
