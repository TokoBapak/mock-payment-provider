package primitive

import "time"

type Transaction struct {
	OrderId           string
	TransactionAmount int64
	PaymentType       PaymentType
	TransactionStatus TransactionStatus
	TransactionTime   time.Time
	ExpiresAt         time.Time
}

func (t Transaction) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}
