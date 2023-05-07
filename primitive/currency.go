package primitive

type Currency uint8

const (
	CurrencyUnspecified Currency = iota
	CurrencyIDR
	CurrencyUSD
)

func (c Currency) String() string {
	switch c {
	case CurrencyIDR:
		return "IDR"
	case CurrencyUSD:
		return "USD"
	case CurrencyUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}
