package presentation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
)

func (p *Presenter) InternalTransactionDetail(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())

	transactionId := r.URL.Query().Get("id")
	if transactionId == "" {
		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    400,
			StatusMessage: "Empty transaction id",
			Id:            "",
		})
		if err != nil {
			log.Err(err).Msg("marshaling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	transactionDetail, err := p.paymentService.GetDetail(r.Context(), transactionId)
	if err != nil {
		if errors.Is(err, business.ErrTransactionNotFound) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    400,
				StatusMessage: "Transaction was not found",
				Id:            "",
			})
			if err != nil {
				log.Err(err).Msg("marshaling json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBody)
			return
		}

		log.Err(err).Str("transaction_id", transactionId).Msg("executing business function")

		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    500,
			StatusMessage: err.Error(),
			Id:            "",
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

	responseBody, err := json.Marshal(schema.InternalTransactionDetailResponse{
		OrderId:              transactionDetail.OrderId,
		ChargedAmount:        transactionDetail.ChargedAmount,
		TransactionStatus:    transactionDetail.Status.String(),
		PaymentMethod:        transactionDetail.PaymentMethod.ToPaymentMethod(),
		Bank:                 transactionDetail.PaymentMethod.ToBank(),
		VirtualAccountNumber: transactionDetail.VirtualAccountNumber,
		EMoneyId:             transactionDetail.EMoneyID,
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
