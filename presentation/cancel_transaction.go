package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"

	"github.com/go-chi/chi/v5"
)

func (p *Presenter) CancelTransaction(w http.ResponseWriter, r *http.Request) {
	// Parse input
	orderId := chi.URLParam(r, "order_id")
	if orderId == "" {
		responseBody, e := json.Marshal(schema.Error{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: "Empty Order ID",
		})
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	// Call business logic
	cancelResponse, err := p.transactionService.Cancel(r.Context(), orderId)
	if err != nil {
		if errors.Is(err, business.ErrTransactionNotFound) {
			// TODO: handle error, return 404
		}

		// TODO: handle error, return 500
	}

	// Return output
	responseBody, e := json.Marshal(schema.CancelTransactionResponse{
		StatusCode:        http.StatusOK,
		StatusMessage:     "Success, transaction is canceled",
		OrderId:           orderId,
		PaymentType:       cancelResponse.PaymentType.String(),
		TransactionTime:   cancelResponse.TransactionTime.Format(time.RFC3339),
		TransactionStatus: cancelResponse.TransactionStatus.String(),
	})
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
	return
}
