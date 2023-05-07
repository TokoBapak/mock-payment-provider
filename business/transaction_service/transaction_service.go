package transaction_service

import (
	"fmt"

	"mock-payment-provider/repository"
)

type Config struct {
	ServerKey                string
	TransactionRepository    repository.TransactionRepository
	WebhookClient            repository.WebhookClient
	VirtualAccountRepository repository.VirtualAccountRepository
	EMoneyRepository         repository.EMoneyRepository
}

type Dependency struct {
	serverKey                string
	transactionRepository    repository.TransactionRepository
	webhookClient            repository.WebhookClient
	virtualAccountRepository repository.VirtualAccountRepository
	emoneyRepository         repository.EMoneyRepository
}

// NewTransactionService validates input from Dependency and return an error if
// any of it is nil. It implements business.Transaction interface.
func NewTransactionService(config Config) (*Dependency, error) {
	if config.TransactionRepository == nil {
		return &Dependency{}, fmt.Errorf("nil transaction repository")
	}

	if config.WebhookClient == nil {
		return &Dependency{}, fmt.Errorf("nil webhook client")
	}

	if config.VirtualAccountRepository == nil {
		return &Dependency{}, fmt.Errorf("nil virtual account repository")
	}

	if config.EMoneyRepository == nil {
		return &Dependency{}, fmt.Errorf("nil emoney repository")
	}

	return &Dependency{
		serverKey:                config.ServerKey,
		transactionRepository:    config.TransactionRepository,
		webhookClient:            config.WebhookClient,
		virtualAccountRepository: config.VirtualAccountRepository,
		emoneyRepository:         config.EMoneyRepository,
	}, nil
}
