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
	// TODO: lanjutin switch case nya
	case PaymentTypeEMoneyQRIS:
		return "EMONEY_QRIS"
		// TODO: lanjutin switch case nya
	case PaymentTypeUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}
