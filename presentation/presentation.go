package presentation

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
	"mock-payment-provider/primitive"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Presenter struct {
	transactionService business.Transaction
	paymentService     business.Payment
}

type Dependency struct {
	TransactionService business.Transaction
	PaymentService     business.Payment
	Logger             zerolog.Logger
}
type PresenterConfig struct {
	Hostname   string
	Port       string
	ServerKey  string
	Dependency *Dependency
}

func NewPresenter(config PresenterConfig) (*http.Server, error) {
	presenter := &Presenter{
		transactionService: config.Dependency.TransactionService,
		paymentService:     config.Dependency.PaymentService,
	}

	router := chi.NewRouter()

	router.Use(hlog.NewHandler(config.Dependency.Logger))
	router.Use(hlog.URLHandler("request_url"))

	router.Get("/", presenter.Index)

	// Internal routes
	router.Post("/internal/mark-as-paid", presenter.InternalMarkAsPaid)
	router.Get("/internal/transaction-detail", presenter.InternalTransactionDetail)

	// Apply authorization middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, _, ok := r.BasicAuth()
			if !ok || user != config.ServerKey {
				responseBody, err := json.Marshal(schema.Error{
					StatusCode:    401,
					StatusMessage: "Transaction cannot be authorized with the current client/server key.",
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(responseBody)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// External routes
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
