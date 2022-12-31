package transaction_service

import (
	"fmt"

	"mock-payment-provider/repository/transaction"
	"mock-payment-provider/repository/webhook"
)

type Dependency struct {
	TransactionRepository transaction.ITransactionRepository
	WebhookClient         *webhook.Client
}

func NewTransactionService(dependency Dependency) (*Dependency, error) {
	if dependency.TransactionRepository == nil {
		return &Dependency{}, fmt.Errorf("nil transaction repository")
	}

	if dependency.WebhookClient == nil {
		return &Dependency{}, fmt.Errorf("webhook client is nil")
	}

	return &Dependency{
		TransactionRepository: dependency.TransactionRepository,
		WebhookClient:         dependency.WebhookClient,
	}, nil
}
