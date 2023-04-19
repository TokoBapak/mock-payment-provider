package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mock-payment-provider/presentation/schema"
)

func (p *Presenter) ExpireTransaction(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "order_id")
	if orderId == "" {
		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    404,
			StatusMessage: "Transaction doesn't exist.",
			Id:            uuid.NewString(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Do not complain. Do complain to Midtrans about the status code usage instead.
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
		return
	}

	// TODO
}
