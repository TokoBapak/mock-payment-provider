package emoney

import (
	"context"
	"fmt"
)

func (r *Repository) DeductCharge(ctx context.Context, orderId string) error {
	if orderId == "" {
		return fmt.Errorf("orderId is empty")
	}

	// No need to do anything. Trust me.
	return nil
}
