package primitive

type PaymentType uint8

const (
	PaymentTypeUnspecified PaymentType = iota
	PaymentTypeVirtualAccountBCA
	PaymentTypeVirtualAccountPermata
	PaymentTypeVirtualAccountBRI
	PaymentTypeVirtualAccountBNI
	PaymentTypeEMoneyQRIS
	PaymentTypeEMoneyGopay
	PaymentTypeEMoneyShopeePay
)

func (p PaymentType) String() string {
	switch p {
	case PaymentTypeVirtualAccountBCA:
		return "VIRTUAL_ACCOUNT_BCA"
	case PaymentTypeVirtualAccountPermata:
		return "VIRTUAL_ACCOUNT_PERMATA"
	case PaymentTypeVirtualAccountBRI:
		return "VIRTUAL_ACCOUNT_BRI"
	case PaymentTypeVirtualAccountBNI:
		return "VIRTUAL_ACCOUNT_BNI"
	case PaymentTypeEMoneyQRIS:
		return "E_MONEY_QRIS"
	case PaymentTypeEMoneyGopay:
		return "E_MONEY_GOPAY"
	case PaymentTypeEMoneyShopeePay:
		return "E_MONEY_SHOPEE_PAY"
	case PaymentTypeUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}

func (p PaymentType) ToPaymentMethod() string {
	switch p {
	case PaymentTypeVirtualAccountBCA:
		fallthrough
	case PaymentTypeVirtualAccountPermata:
		fallthrough
	case PaymentTypeVirtualAccountBRI:
		fallthrough
	case PaymentTypeVirtualAccountBNI:
		return "bank_transfer"
	case PaymentTypeEMoneyQRIS:
		return "qris"
	case PaymentTypeEMoneyGopay:
		return "gopay"
	case PaymentTypeEMoneyShopeePay:
		return "shopeepay"
	case PaymentTypeUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}

func (p PaymentType) ToBank() string {
	switch p {
	case PaymentTypeVirtualAccountBCA:
		return "bca"
	case PaymentTypeVirtualAccountPermata:
		return "permata"
	case PaymentTypeVirtualAccountBRI:
		return "bri"
	case PaymentTypeVirtualAccountBNI:
		return "bni"
	default:
		return ""
	}
}
