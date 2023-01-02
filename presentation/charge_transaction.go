package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
)

func (p *Presenter) ChargeTransaction(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody schema.ChargeTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		responseBody, e := json.Marshal(schema.Error{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: "Malformed JSON",
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

	// Convert to common business schema
	chargeRequest := business.ChargeRequest{
		PaymentType:         PaymentTypeMap[requestBody.PaymentType],
		OrderId:             requestBody.Transaction.OrderId,
		TransactionAmount:   requestBody.Transaction.Amount,
		TransactionCurrency: CurrencyMap[requestBody.Transaction.Currency],

		Customer: business.CustomerInformation{
			FirstName:   requestBody.Customer.FirstName,
			LastName:    requestBody.Customer.LastName,
			Email:       requestBody.Customer.Email,
			PhoneNumber: requestBody.Customer.PhoneNumber,
			BillingAddress: business.Address{
				FirstName:   requestBody.Customer.BillingAddress.FirstName,
				LastName:    requestBody.Customer.BillingAddress.LastName,
				Email:       requestBody.Customer.BillingAddress.Email,
				Phone:       requestBody.Customer.BillingAddress.Phone,
				Address:     requestBody.Customer.BillingAddress.Address,
				PostalCode:  requestBody.Customer.BillingAddress.PostalCode,
				CountryCode: requestBody.Customer.BillingAddress.CountryCode,
			},
		},

		Seller: business.SellerInformation{
			FirstName:   requestBody.Seller.FirstName,
			LastName:    requestBody.Seller.LastName,
			Email:       requestBody.Seller.Email,
			PhoneNumber: requestBody.Seller.PhoneNumber,
			Address:     requestBody.Seller.Address,
		},
	}
	for _, item := range requestBody.Items {
		chargeRequest.ProductItems = append(chargeRequest.ProductItems, business.ProductItem{
			ID:       item.Id,
			Price:    int64(item.Price),
			Quantity: int64(item.Quantity),
			Name:     item.Name,
			Category: item.Category,
		})
	}

	// Call business function
	chargeResponse, err := p.transactionService.Charge(r.Context(), chargeRequest)
	if err != nil {
		// Handle known error
		if errors.Is(err, business.ErrDuplicateOrderId) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    406,
				StatusMessage: "Duplicate order ID. order_id has already been utilized previously.",
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(406)
			w.Write(responseBody)
			return
		}

		if errors.Is(err, business.ErrMismatchedTransactionAmount) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    400,
				StatusMessage: "Mismatched total transaction amount with the accumulated amount from product list",
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(responseBody)
			return
		}

		// if kind of error is RequestValidationError
		if errors.As(err, &business.RequestValidationError{}) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    400,
				StatusMessage: err.Error(),
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(responseBody)
			return
		}

		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "Internal Server Error.",
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

	// Send return output to the client
	responseBody, err := json.Marshal(schema.ChargeTransactionResponse{
		StatusCode:        http.StatusCreated,
		StatusMessage:     "The transaction is created successfully",
		OrderId:           chargeResponse.OrderId,
		GrossAmount:       chargeRequest.TransactionAmount,
		PaymentType:       chargeResponse.PaymentType.String(),
		TransactionTime:   chargeResponse.TransactionTime.Format(time.RFC3339),
		TransactionStatus: chargeResponse.TransactionStatus.String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseBody)
}
