package primitive_test

import (
	"testing"

	"mock-payment-provider/primitive"
)

func TestPaymentType_String(t *testing.T) {
	t.Run("PaymentTypeVirtualAccountBCA", func(t *testing.T) {
		if primitive.PaymentTypeVirtualAccountBCA.String() != "VIRTUAL_ACCOUNT_BCA" {
			t.Errorf("expecting PaymentTypeVirtualAccountBCA.String() to be 'VIRTUAL_ACCOUNT_BCA', instead got %s", primitive.PaymentTypeVirtualAccountBCA.String())
		}
	})

	t.Run("PaymentTypeVirtualAccountMandiri", func(t *testing.T) {
		if primitive.PaymentTypeVirtualAccountMandiri.String() != "VIRTUAL_ACCOUNT_MANDIRI" {
			t.Errorf("expecting PaymentTypeVirtualAccountMandiri.String() to be 'VIRTUAL_ACCOUNT_MANDIRI', instead got %s", primitive.PaymentTypeVirtualAccountMandiri.String())
		}
	})

	t.Run("PaymentTypeVirtualAccountBRI", func(t *testing.T) {
		if primitive.PaymentTypeVirtualAccountBRI.String() != "VIRTUAL_ACCOUNT_BRI" {
			t.Errorf("expecting PaymentTypeVirtualAccountBRI.String() to be 'VIRTUAL_ACCOUNT_BRI', instead got %s", primitive.PaymentTypeVirtualAccountBRI.String())
		}
	})

	t.Run("PaymentTypeVirtualAccountBNI", func(t *testing.T) {
		if primitive.PaymentTypeVirtualAccountBNI.String() != "VIRTUAL_ACCOUNT_BNI" {
			t.Errorf("expecting PaymentTypeVirtualAccountBNI.String() to be 'VIRTUAL_ACCOUNT_BNI', instead got %s", primitive.PaymentTypeVirtualAccountBNI.String())
		}
	})

	t.Run("PaymentTypeEMoneyQRIS", func(t *testing.T) {
		if primitive.PaymentTypeEMoneyQRIS.String() != "E_MONEY_QRIS" {
			t.Errorf("expecting PaymentTypeEMoneyQRIS.String() to be 'E_MONEY_QRIS', instead got %s", primitive.PaymentTypeEMoneyQRIS.String())
		}
	})

	t.Run("PaymentTypeEMoneyGopay", func(t *testing.T) {
		if primitive.PaymentTypeEMoneyGopay.String() != "E_MONEY_GOPAY" {
			t.Errorf("expecting PaymentTypeEMoneyGopay.String() to be 'E_MONEY_GOPAY', instead got %s", primitive.PaymentTypeEMoneyGopay.String())
		}
	})

	t.Run("PaymentTypeEMoneyShopeePay", func(t *testing.T) {
		if primitive.PaymentTypeEMoneyShopeePay.String() != "E_MONEY_SHOPEE_PAY" {
			t.Errorf("expecting PaymentTypeEMoneyShopeePay.String() to be 'E_MONEY_SHOPEE_PAY', instead got %s", primitive.PaymentTypeEMoneyShopeePay.String())
		}
	})

	t.Run("PaymentTypeUnspecified", func(t *testing.T) {
		if primitive.PaymentTypeUnspecified.String() != "UNSPECIFIED" {
			t.Errorf("expecting PaymentTypeUnspecified.String() to be 'UNSPECIFIED', instead got %s", primitive.PaymentTypeUnspecified.String())
		}
	})
}
