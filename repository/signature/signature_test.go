package signature_test

import (
	"testing"

	"mock-payment-provider/repository/signature"
)

func TestGenerate(t *testing.T) {
	sig := signature.Generate("AABBCC", 200, 50_000, "SERVER KEY")
	if sig == "" {
		t.Error("expecting a signature, instead got empty string")
	}

	expect := "019dabd47e714a526adc8055c88534b9222fc0dff6d718edab908db5568377893573cf0f972e695b53ecdb5ffd919e1d47d6b92a8ceafb5c108e45062693fdbf"
	if sig != expect {
		t.Errorf("expecting signature to be '%s', instead got '%s'", expect, sig)
	}
}
