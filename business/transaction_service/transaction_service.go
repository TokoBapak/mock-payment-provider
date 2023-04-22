package transaction_service

import (
	"fmt"

	"mock-payment-provider/repository"
)

type Dependency struct {
	TransactionRepository    repository.TransactionRepository
	WebhookClient            repository.WebhookClient
	VirtualAccountRepository repository.VirtualAccountRepository
	EMoneyRepository         repository.EMoneyRepository
}

// NewTransactionService validates input from Dependency and return an error if
// any of it is nil. It implements business.Transaction interface.
func NewTransactionService(dependency Dependency) (*Dependency, error) {
	if dependency.TransactionRepository == nil {
		return &Dependency{}, fmt.Errorf("nil transaction repository")
	}

	if dependency.WebhookClient == nil {
		return &Dependency{}, fmt.Errorf("nil webhook client")
	}

	if dependency.VirtualAccountRepository == nil {
		return &Dependency{}, fmt.Errorf("nil virtual account repository")
	}

	if dependency.EMoneyRepository == nil {
		return &Dependency{}, fmt.Errorf("nil emoney repository")
	}

	return &Dependency{
		TransactionRepository:    dependency.TransactionRepository,
		WebhookClient:            dependency.WebhookClient,
		VirtualAccountRepository: dependency.VirtualAccountRepository,
		EMoneyRepository:         dependency.EMoneyRepository,
	}, nil
}
