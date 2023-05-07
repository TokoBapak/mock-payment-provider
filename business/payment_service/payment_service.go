package payment_service

import (
	"fmt"

	"mock-payment-provider/repository"
)

type Config struct {
	ServerKey                string
	TransactionRepository    repository.TransactionRepository
	WebhookClient            repository.WebhookClient
	EMoneyRepository         repository.EMoneyRepository
	VirtualAccountRepository repository.VirtualAccountRepository
}

type Dependency struct {
	serverKey                string
	transactionRepository    repository.TransactionRepository
	webhookClient            repository.WebhookClient
	eMoneyRepository         repository.EMoneyRepository
	virtualAccountRepository repository.VirtualAccountRepository
}

func NewPaymentService(config Config) (*Dependency, error) {
	if config.TransactionRepository == nil {
		return nil, fmt.Errorf("nil transaction repository")
	}

	if config.WebhookClient == nil {
		return nil, fmt.Errorf("nil webhook client")
	}

	if config.EMoneyRepository == nil {
		return nil, fmt.Errorf("nil emoney repository")
	}

	if config.VirtualAccountRepository == nil {
		return nil, fmt.Errorf("nil virtual account repository")
	}

	return &Dependency{
		serverKey:                config.ServerKey,
		transactionRepository:    config.TransactionRepository,
		webhookClient:            config.WebhookClient,
		eMoneyRepository:         config.EMoneyRepository,
		virtualAccountRepository: config.VirtualAccountRepository,
	}, nil
}
