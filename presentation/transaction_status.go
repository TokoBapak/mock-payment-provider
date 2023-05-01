package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
	"mock-payment-provider/repository/signature"

	"github.com/google/uuid"
)

func (p *Presenter) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
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

	status, err := p.transactionService.GetStatus(r.Context(), orderId)
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

		log.Err(err).Str("order_id", orderId).Msg("executing business function")

		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    500,
			StatusMessage: "Internal server error.",
			Id:            uuid.NewString(),
		})
		if err != nil {
			log.Err(err).Msg("marshaling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}

	signatureKey := signature.Generate(status.OrderId, 200, status.TransactionAmount, "")

	responseBody, err := json.Marshal(schema.TransactionStatusResponse{
		StatusCode:               "200",
		StatusMessage:            "Success, transaction found",
		TransactionId:            status.OrderId,
		MaskedCard:               "",
		OrderId:                  status.OrderId,
		PaymentType:              status.PaymentType.ToPaymentMethod(),
		TransactionTime:          status.TransactionTime.Format(time.DateTime),
		TransactionStatus:        status.TransactionStatus.String(),
		FraudStatus:              "accept",
		ApprovalCode:             "",
		SignatureKey:             signatureKey,
		Bank:                     status.PaymentType.ToBank(),
		GrossAmount:              status.TransactionAmount,
		ChannelResponseCode:      "",
		ChannelResponseMessage:   "",
		CardType:                 "",
		PaymentOptionType:        "",
		ShopeepayReferenceNumber: "",
		ReferenceId:              "",
	})
	if err != nil {
		log.Err(err).Msg("marshaling json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
