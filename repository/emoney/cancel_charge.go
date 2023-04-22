package emoney

import "context"

func (r *Repository) CancelCharge(ctx context.Context, id string) error {
	// No need to do anything. Cancel
	return nil
}
