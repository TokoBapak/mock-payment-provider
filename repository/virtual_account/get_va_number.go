package virtual_account

import (
	"context"

	"mock-payment-provider/repository"
)

func (r *Repository) GetByVirtualAccountNumber(ctx context.Context, virtualAccountNumber string) (repository.Entry, error) {
	//TODO implement me
	panic("implement me")
}
