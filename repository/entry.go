package repository

import "time"

type Entry struct {
	VirtualAccountNumber string
	EMoneyID             string
	OrderId              string
	ChargedAmount        int64
	ExpiresAt            time.Time
}

func (e Entry) Expired() bool {
	return e.ExpiresAt.Before(time.Now())
}
