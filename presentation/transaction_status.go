package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"

	"github.com/google/uuid"
)

func (p *Presenter) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
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

	status, err := p.transactionService.GetStatus(r.Context(), orderId)
	if err != nil {
		if errors.Is(err, business.ErrTransactionNotFound) {
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
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
			return
		}

		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    500,
			StatusMessage: "Internal server error.",
			Id:            uuid.NewString(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}

	responseBody, err := json.Marshal(schema.TransactionStatusResponse{
		StatusCode:               "200",
		StatusMessage:            "Success, transaction found",
		TransactionId:            "",
		MaskedCard:               "",
		OrderId:                  status.OrderId,
		PaymentType:              status.PaymentType.String(),
		TransactionTime:          status.TransactionTime.Format(time.RFC3339),
		TransactionStatus:        status.TransactionStatus.String(),
		FraudStatus:              "",
		ApprovalCode:             "",
		SignatureKey:             "",
		Bank:                     "",
		GrossAmount:              status.TransactionAmount,
		ChannelResponseCode:      "",
		ChannelResponseMessage:   "",
		CardType:                 "",
		PaymentOptionType:        "",
		ShopeepayReferenceNumber: "",
		ReferenceId:              "",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseBody)
}
