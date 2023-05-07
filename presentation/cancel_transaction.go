package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"

	"github.com/go-chi/chi/v5"
)

func (p *Presenter) CancelTransaction(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())

	orderId := chi.URLParam(r, "order_id")
	if orderId == "" {
		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    404,
			StatusMessage: "Transaction doesn't exist.",
			Id:            uuid.NewString(),
		})
		if err != nil {
			log.Err(err).Msg("marshaling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Do not complain. Do complain to Midtrans about the status code usage instead.
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
		return
	}

	// Call business logic
	cancelResponse, err := p.transactionService.Cancel(r.Context(), orderId)
	if err != nil {
		if errors.Is(err, business.ErrTransactionNotFound) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    404,
				StatusMessage: "Transaction doesn't exist.",
				Id:            uuid.NewString(),
			})
			if err != nil {
				log.Err(err).Msg("marshaling json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
			return
		}

		if errors.Is(err, business.ErrCannotModifyStatus) {
			responseBody, e := json.Marshal(schema.Error{
				StatusCode:    412,
				StatusMessage: "Merchant cannot modify the status of the transaction",
			})
			if e != nil {
				log.Err(err).Msg("marshaling json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
			w.Header().Set("Content-Type", "application/json")
			return
		}

		// TODO: send to logger

		responseBody, e := json.Marshal(schema.Error{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "internal server error",
		})
		if e != nil {
			log.Err(err).Msg("marshaling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)

		return
	}

	// Return output
	responseBody, e := json.Marshal(schema.CancelTransactionResponse{
		StatusCode:        http.StatusOK,
		StatusMessage:     "Success, transaction is canceled",
		OrderId:           orderId,
		PaymentType:       cancelResponse.PaymentType.String(),
		TransactionTime:   cancelResponse.TransactionTime.Format(time.DateTime),
		TransactionStatus: cancelResponse.TransactionStatus.String(),
	})
	if e != nil {
		log.Err(err).Msg("marshaling json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
