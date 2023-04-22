package emoney

import (
	"context"
	"time"
)

func (r *Repository) CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (id string, err error) {
	//TODO implement me
	panic("implement me")
}
