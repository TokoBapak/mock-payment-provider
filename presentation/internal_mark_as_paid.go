package presentation

import (
	"encoding/json"
	"errors"
	"net/http"

	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
)

func (p *Presenter) InternalMarkAsPaid(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody schema.InternalMarkAsPaidRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    400,
			StatusMessage: "Invalid request body",
			Id:            "",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	err = p.paymentService.MarkAsPaid(r.Context(), requestBody.OrderId, requestBody.PaymentMethod)
	if err != nil {
		if errors.Is(err, business.ErrCannotModifyStatus) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    400,
				StatusMessage: "Transaction is not from PENDING status",
				Id:            "",
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBody)
			return
		}

		if errors.Is(err, business.ErrTransactionNotFound) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    400,
				StatusMessage: "Transaction was not found",
				Id:            "",
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBody)
			return
		}

		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    500,
			StatusMessage: err.Error(),
			Id:            "",
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

	w.WriteHeader(http.StatusOK)
}
