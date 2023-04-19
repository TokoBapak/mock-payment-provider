package presentation

import (
	"net"
	"net/http"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"

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
	router.Post("/{order_id}/expire", presenter.ExpireTransaction)

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

var paymentTypeMap = map[string]primitive.PaymentType{
	"VIRTUAL_ACCOUNT_BCA":     primitive.PaymentTypeVirtualAccountBCA,
	"VIRTUAL_ACCOUNT_PERMATA": primitive.PaymentTypeVirtualAccountPermata,
	"VIRTUAL_ACCOUNT_BRI":     primitive.PaymentTypeVirtualAccountBRI,
	"VIRTUAL_ACCOUNT_BNI":     primitive.PaymentTypeVirtualAccountBNI,
	"E_MONEY_QRIS":            primitive.PaymentTypeEMoneyQRIS,
	"E_MONEY_GOPAY":           primitive.PaymentTypeEMoneyGopay,
	"E_MONEY_SHOPEE_PAY":      primitive.PaymentTypeEMoneyShopeePay,
}

var currencyMap = map[string]primitive.Currency{
	"IDR": primitive.CurrencyIDR,
	"USD": primitive.CurrencyUSD,
}
