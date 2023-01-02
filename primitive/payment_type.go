package primitive

type PaymentType uint8

const (
	PaymentTypeUnspecified PaymentType = iota
	PaymentTypeVirtualAccountBCA
	PaymentTypeVirtualAccountMandiri
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
	case PaymentTypeVirtualAccountMandiri:
		return "VIRTUAL_ACCOUNT_MANDIRI"
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

var PaymentTypeMap = map[string]PaymentType{
	"VIRTUAL_ACCOUNT_BCA":     PaymentTypeVirtualAccountBCA,
	"VIRTUAL_ACCOUNT_MANDIRI": PaymentTypeVirtualAccountMandiri,
	"VIRTUAL_ACCOUNT_BRI":     PaymentTypeVirtualAccountBRI,
	"VIRTUAL_ACCOUNT_BNI":     PaymentTypeVirtualAccountBNI,
	"E_MONEY_QRIS":            PaymentTypeEMoneyQRIS,
	"E_MONEY_GOPAY":           PaymentTypeEMoneyGopay,
	"E_MONEY_SHOPEE_PAY":      PaymentTypeEMoneyShopeePay,
}
