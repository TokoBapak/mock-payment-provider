package presentation

import (
	"net"
	"net/http"
	"time"

	"mock-payment-provider/business"

	"github.com/go-chi/chi/v5"
)

type Presenter struct {
	transactionService business.Transaction
}

type Dependency struct {
	TransactionService business.Transaction
}
type PresenterConfig struct {
	Hostname   string
	Port       string
	Dependency *Dependency
}

func NewPresenter(config PresenterConfig) (*http.Server, error) {
	presenter := Presenter{
		transactionService: config.Dependency.TransactionService,
	}

	router := chi.NewRouter()

	router.Post("/charge", presenter.ChargeTransaction)
	router.Post("/{order_id}/cancel", presenter.CancelTransaction)
	router.Get("/{order_id}/status", presenter.GetTransactionStatus)

	server := &http.Server{
		Addr:              net.JoinHostPort(config.Hostname, config.Port),
		Handler:           router,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
	}

	return server, nil
}
