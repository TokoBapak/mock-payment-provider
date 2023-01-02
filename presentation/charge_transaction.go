package presentation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
	"mock-payment-provider/primitive"
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

	// Validate the request and Convert to common business schema
	chargeRequest := business.ChargeRequest{}
	if err := validateChargeTransaction(requestBody, &chargeRequest); err != nil {
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
	return
}

// validateChargeTransaction : handle validating the cahrge request data
func validateChargeTransaction(requestBody schema.ChargeTransactionRequest, validated *business.ChargeRequest) *RequestValidationError {
	// validate payment_type
	if primitive.PaymentTypeMap[requestBody.PaymentType] == 0 {
		return &RequestValidationError{
			Err: errors.New("payment_type is not valid"),
		}
	}

	// validate transaction.order_id
	if requestBody.Transaction.OrderId == "" {
		return &RequestValidationError{
			Err: errors.New("order_id is required"),
		}
	}

	// valdiate transaction.amount
	if requestBody.Transaction.Amount <= 0 {
		return &RequestValidationError{
			Err: errors.New("amount should be greater than 0"),
		}
	}

	// validate transaction.currency
	if primitive.CurrencyMap[requestBody.Transaction.Currency] == 0 {
		return &RequestValidationError{
			Err: errors.New("transaction.currency is not valid"),
		}
	}

	// validate customer.first_name
	if requestBody.Customer.FirstName == "" {
		return &RequestValidationError{
			Err: errors.New("customer.first_name is required"),
		}
	}
	if len(requestBody.Customer.FirstName) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.first_name must be less than 255 characters"),
		}
	}

	// validate customer.last_name
	if len(requestBody.Customer.LastName) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.last_name must be less than 255 characters"),
		}
	}

	// validate customer.email
	if requestBody.Customer.Email == "" {
		return &RequestValidationError{
			Err: errors.New("customer.email is required"),
		}
	}

	if len(requestBody.Customer.Email) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.email must be less than 255 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).Match([]byte(requestBody.Customer.Email)); !ok {
		return &RequestValidationError{
			Err: errors.New("customer.email is not valid"),
		}
	}

	// validate customer.phone_number
	if requestBody.Customer.PhoneNumber == "" {
		return &RequestValidationError{
			Err: errors.New("customer.phone_number is required"),
		}
	}

	if len(requestBody.Customer.PhoneNumber) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.phone_number must less than 255 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).Match([]byte(requestBody.Customer.PhoneNumber)); !ok {
		return &RequestValidationError{
			Err: errors.New("customer.phone_number is not valid"),
		}
	}

	// validate customer.billing_address.first_name
	if requestBody.Customer.BillingAddress.FirstName == "" {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.first_name is required"),
		}
	}

	if len(requestBody.Customer.BillingAddress.FirstName) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.first_name must be less than 255 characters"),
		}
	}

	// validate customer.billing_address.last_name
	if len(requestBody.Customer.BillingAddress.LastName) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.last_name must be less than 255 characters"),
		}
	}

	// validate customer.billing_address.email
	if requestBody.Customer.BillingAddress.Email == "" {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.email is required"),
		}
	}

	if len(requestBody.Customer.BillingAddress.Email) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.email must be less than 255 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).Match(
		[]byte(requestBody.Customer.BillingAddress.Email),
	); !ok {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.email is not valid"),
		}
	}

	// validate customer.billing_address.phone
	if requestBody.Customer.BillingAddress.Phone == "" {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.phone is required"),
		}
	}

	if len(requestBody.Customer.BillingAddress.Phone) > 255 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.phone must be less than 255 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).Match(
		[]byte(requestBody.Customer.BillingAddress.Phone),
	); !ok {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.phone is not valid"),
		}
	}

	// validate customer.bliing_address.address
	if requestBody.Customer.BillingAddress.Address == "" {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.address is required"),
		}
	}

	if len(requestBody.Customer.BillingAddress.Address) > 500 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.address must be less than 500 characters"),
		}
	}

	// validate customer.billing_address.postal_code
	if requestBody.Customer.BillingAddress.PostalCode == "" {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.postal_code is required"),
		}
	}

	if len(requestBody.Customer.BillingAddress.PostalCode) > 10 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.postal_code must be less than 10 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.PostalCodePattern).Match(
		[]byte(requestBody.Customer.BillingAddress.PostalCode),
	); !ok {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.postal_code is not valid"),
		}
	}

	// validate customer.billing_address.country_code
	if requestBody.Customer.BillingAddress.CountryCode == "" {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.country_code is required"),
		}
	}

	if len(requestBody.Customer.BillingAddress.CountryCode) > 5 {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.country_code must be less than 5 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.CountryCodePattern).Match(
		[]byte(requestBody.Customer.BillingAddress.CountryCode),
	); !ok {
		return &RequestValidationError{
			Err: errors.New("customer.billing_address.country_code is not valid"),
		}
	}

	// validate seller.first_name
	if requestBody.Seller.FirstName == "" {
		return &RequestValidationError{
			Err: errors.New("seller.first_name is required"),
		}
	}

	if len(requestBody.Seller.FirstName) > 255 {
		return &RequestValidationError{
			Err: errors.New("seller.first_name must be less than 255 characters"),
		}
	}

	// validate seller.last_name
	if len(requestBody.Seller.LastName) > 255 {
		return &RequestValidationError{
			Err: errors.New("seller.last_name must be less than 255 characters"),
		}
	}

	// validate seller.email
	if requestBody.Seller.Email == "" {
		return &RequestValidationError{
			Err: errors.New("seller.email is required"),
		}
	}

	if len(requestBody.Seller.Email) > 255 {
		return &RequestValidationError{
			Err: errors.New("seller.email must be less than 255 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).Match(
		[]byte(requestBody.Seller.Email),
	); !ok {
		return &RequestValidationError{
			Err: errors.New("seller.email is not valid"),
		}
	}

	// validate seller.phone_number
	if requestBody.Seller.PhoneNumber == "" {
		return &RequestValidationError{
			Err: errors.New("seller.phone_number is required"),
		}
	}

	if len(requestBody.Seller.PhoneNumber) > 255 {
		return &RequestValidationError{
			Err: errors.New("seller.phone_number must less than 255 characters"),
		}
	}

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).Match(
		[]byte(requestBody.Seller.PhoneNumber),
	); !ok {
		return &RequestValidationError{
			Err: errors.New("seller.phone_number is not valid"),
		}
	}

	// validate seller.address
	if requestBody.Seller.Address == "" {
		return &RequestValidationError{
			Err: errors.New("seller.address is required"),
		}
	}

	if len(requestBody.Seller.Address) > 500 {
		return &RequestValidationError{
			Err: errors.New("seller.address must be less than 500 characters"),
		}
	}

	// validate items.id
	if len(requestBody.Items) == 0 {
		return &RequestValidationError{
			Err: errors.New("items.request_body must be greater than 0 length"),
		}
	}

	// validate items
	for i, item := range requestBody.Items {
		// validate items.id
		if item.Id == "" {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.id is required", i),
			}
		}

		if len(item.Id) > 255 {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.id must be less than 255 characters", i),
			}
		}

		// validate items.price
		if item.Price <= 0 {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.price must be greather than 0", i),
			}
		}

		// validate itmes.quantity
		if item.Quantity <= 0 {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.quantity must be greater than 0", i),
			}
		}

		// validate items.name
		if item.Name == "" {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.quantity is required", i),
			}
		}

		if len(item.Name) > 255 {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.name must be less than 255 characters", i),
			}
		}

		// validate items.category
		if item.Category == "" {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.category is required", i),
			}
		}

		if len(item.Category) > 255 {
			return &RequestValidationError{
				Err: fmt.Errorf("items.%d.category must be less than 255 characters", i),
			}
		}
	}

	if validated != nil {
		validated.PaymentType = primitive.PaymentTypeMap[requestBody.PaymentType]
		validated.OrderId = requestBody.Transaction.OrderId
		validated.TransactionAmount = requestBody.Transaction.Amount
		validated.TransactionCurrency = primitive.CurrencyMap[requestBody.Transaction.Currency]

		validated.Customer.FirstName = requestBody.Customer.FirstName
		validated.Customer.LastName = requestBody.Customer.LastName
		validated.Customer.Email = requestBody.Customer.Email
		validated.Customer.PhoneNumber = requestBody.Customer.PhoneNumber

		validated.Customer.BillingAddress.FirstName = requestBody.Customer.BillingAddress.FirstName
		validated.Customer.BillingAddress.LastName = requestBody.Customer.BillingAddress.LastName
		validated.Customer.BillingAddress.Email = requestBody.Customer.BillingAddress.Email
		validated.Customer.BillingAddress.Phone = requestBody.Customer.BillingAddress.Phone
		validated.Customer.BillingAddress.Address = requestBody.Customer.BillingAddress.Address
		validated.Customer.BillingAddress.PostalCode = requestBody.Customer.BillingAddress.PostalCode
		validated.Customer.BillingAddress.CountryCode = requestBody.Customer.BillingAddress.CountryCode

		validated.Seller.FirstName = requestBody.Seller.FirstName
		validated.Seller.LastName = requestBody.Seller.LastName
		validated.Seller.Email = requestBody.Seller.Email
		validated.Seller.PhoneNumber = requestBody.Seller.PhoneNumber
		validated.Seller.Address = requestBody.Seller.Address

		for _, item := range requestBody.Items {
			validated.ProductItems = append(validated.ProductItems, business.ProductItem{
				ID:       item.Id,
				Price:    int64(item.Price),
				Quantity: int64(item.Quantity),
				Name:     item.Name,
				Category: item.Category,
			})
		}

	}

	return nil
}
