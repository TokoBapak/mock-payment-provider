package presentation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
			log.Err(err).Msg("marshaling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	paymentType, err := parseValidPaymentMethod(requestBody)
	if err != nil {
		responseBody, e := json.Marshal(schema.Error{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: err.Error(),
		})
		if e != nil {
			log.Err(err).Msg("marshaling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	var callbackURL string
	if paymentType == primitive.PaymentTypeEMoneyGopay && requestBody.Gopay.EnableCallback {
		callbackURL = requestBody.Gopay.CallbackURL
	} else if paymentType == primitive.PaymentTypeEMoneyShopeePay {
		callbackURL = requestBody.ShopeePay.CallbackURL
	}

	var productItems []business.ProductItem
	for _, product := range requestBody.ItemDetails {
		productItems = append(productItems, business.ProductItem{
			ID:       product.Id,
			Price:    product.Price,
			Quantity: product.Quantity,
			Name:     product.Name,
			Category: product.Category,
		})
	}

	// Convert to common business schema
	chargeRequest := business.ChargeRequest{
		PaymentType:         paymentType,
		OrderId:             requestBody.TransactionDetails.OrderId,
		TransactionAmount:   requestBody.TransactionDetails.GrossAmount,
		TransactionCurrency: currencyMap[requestBody.TransactionDetails.Currency],
		Customer: business.CustomerInformation{
			FirstName:   requestBody.CustomerDetails.FirstName,
			LastName:    requestBody.CustomerDetails.LastName,
			Email:       requestBody.CustomerDetails.Email,
			PhoneNumber: requestBody.CustomerDetails.PhoneNumber,
			BillingAddress: business.Address{
				FirstName:   requestBody.CustomerDetails.BillingAddress.FirstName,
				LastName:    requestBody.CustomerDetails.BillingAddress.LastName,
				Email:       requestBody.CustomerDetails.BillingAddress.Email,
				Phone:       requestBody.CustomerDetails.BillingAddress.Phone,
				Address:     requestBody.CustomerDetails.BillingAddress.Address,
				PostalCode:  requestBody.CustomerDetails.BillingAddress.PostalCode,
				CountryCode: requestBody.CustomerDetails.BillingAddress.CountryCode,
			},
		},
		Seller: business.SellerInformation{
			FirstName:   requestBody.Seller.FirstName,
			LastName:    requestBody.Seller.LastName,
			Email:       requestBody.Seller.Email,
			PhoneNumber: requestBody.Seller.PhoneNumber,
			Address:     requestBody.Seller.Address,
		},
		ProductItems: productItems,
		BankTransferOptions: business.BankTransferOptions{
			VirtualAccountNumber: requestBody.BankTransfer.VirtualAccountNumber,
			RecipientName:        requestBody.BankTransfer.Permata.RecipientName,
		},
		EMoneyOptions: business.EMoneyOptions{
			CallbackURL: callbackURL,
		},
	}
	for _, item := range requestBody.ItemDetails {
		chargeRequest.ProductItems = append(chargeRequest.ProductItems, business.ProductItem{
			ID:       item.Id,
			Price:    item.Price,
			Quantity: item.Quantity,
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
				Id:            uuid.NewString(),
			})
			if err != nil {
				log.Err(err).Msg("marshaling json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			// Do not complain. See Midtrans API if you don't believe me.
			w.WriteHeader(200)
			w.Write(responseBody)
			return
		}

		if errors.Is(err, business.ErrMismatchedTransactionAmount) {
			responseBody, err := json.Marshal(schema.Error{
				StatusCode:    400,
				StatusMessage: "Mismatched total transaction amount with the accumulated amount from product list",
			})
			if err != nil {
				log.Err(err).Msg("marshaling json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(responseBody)
			return
		}

		// if kind of error is RequestValidationError
		var requestValidationError *business.RequestValidationError
		if errors.As(err, &requestValidationError) {
			// make the specific struct for request validation error
			validationError := schema.ValidationError{
				Error: schema.Error{
					StatusCode:    400,
					StatusMessage: "some request validation is failed",
				},
			}
			if err, ok := err.(*business.RequestValidationError); ok {
				for _, issue := range err.Issues {
					validationError.Issues = append(validationError.Issues, schema.ValidationIssue{
						Field:   issue.Field,
						Code:    issue.Code.String(),
						Message: fmt.Sprintf("%s %s", issue.Field, issue.Message),
					})
				}

			}

			responseBody, err := json.Marshal(validationError)
			if err != nil {
				log.Err(err).Msg("marshaling json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(responseBody)
			return
		}

		log.Err(err).Any("charge_request", chargeRequest).Msg("executing business function")

		responseBody, err := json.Marshal(schema.Error{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "Internal Server Error.",
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

	var emoneyActions []struct {
		Name   string `json:"name"`
		Method string `json:"method"`
		Url    string `json:"url"`
	}

	for _, emoneyAction := range chargeResponse.EMoneyAction {
		emoneyActions = append(emoneyActions, struct {
			Name   string `json:"name"`
			Method string `json:"method"`
			Url    string `json:"url"`
		}{
			Name:   emoneyAction.EMoneyActionType.String(),
			Method: emoneyAction.Method,
			Url:    emoneyAction.URL,
		})
	}

	// Send return output to the client
	switch chargeResponse.PaymentType {
	case primitive.PaymentTypeEMoneyGopay:
		responseBody, err := json.Marshal(schema.GopayChargeSuccessResponse{
			StatusCode:             "201",
			StatusMessage:          "GO-PAY billing created",
			TransactionId:          chargeResponse.OrderId,
			OrderId:                chargeResponse.OrderId,
			GrossAmount:            strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			PaymentType:            chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:        chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus:      chargeResponse.TransactionStatus.String(),
			Actions:                emoneyActions,
			ChannelResponseCode:    "200",
			ChannelResponseMessage: "Success",
			Currency:               "IDR",
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
	case primitive.PaymentTypeEMoneyShopeePay:
		responseBody, err := json.Marshal(schema.ShopeePayChargeSuccessResponse{
			StatusCode:             "201",
			StatusMessage:          "ShopeePay transaction is created",
			ChannelResponseCode:    "0",
			ChannelResponseMessage: "success",
			TransactionId:          chargeResponse.OrderId,
			OrderId:                chargeResponse.OrderId,
			MerchantId:             "MOCK",
			GrossAmount:            strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			Currency:               "IDR",
			PaymentType:            chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:        chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus:      chargeResponse.TransactionStatus.String(),
			FraudStatus:            "accept",
			Actions:                emoneyActions,
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
	case primitive.PaymentTypeEMoneyQRIS:
		responseBody, err := json.Marshal(schema.QRISChargeSuccessResponse{
			StatusCode:        "201",
			StatusMessage:     "QRIS transaction is created",
			TransactionId:     chargeResponse.OrderId,
			OrderId:           chargeResponse.OrderId,
			MerchantId:        "MOCK",
			GrossAmount:       strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			Currency:          "IDR",
			PaymentType:       chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:   chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus: chargeResponse.TransactionStatus.String(),
			FraudStatus:       "accept",
			Acquirer:          "nobu",
			Actions:           nil,
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
	case primitive.PaymentTypeVirtualAccountBCA:
		responseBody, err := json.Marshal(schema.BCAVirtualAccountChargeSuccessResponse{
			StatusCode:        "201",
			StatusMessage:     "Success, Bank Transfer transaction is created",
			TransactionId:     chargeResponse.OrderId,
			OrderId:           chargeResponse.OrderId,
			GrossAmount:       strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			PaymentType:       chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:   chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus: chargeResponse.TransactionStatus.String(),
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     chargeResponse.VirtualAccountAction.Bank,
					VaNumber: chargeResponse.VirtualAccountAction.VirtualAccountNumber,
				},
			},
			FraudStatus: "accept",
			Currency:    "IDR",
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
	case primitive.PaymentTypeVirtualAccountBRI:
		responseBody, err := json.Marshal(schema.BRIVirtualAccountChargeSuccessResponse{
			StatusCode:        "201",
			StatusMessage:     "Success, Bank Transfer transaction is created",
			TransactionId:     chargeResponse.OrderId,
			OrderId:           chargeResponse.OrderId,
			GrossAmount:       strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			PaymentType:       chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:   chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus: chargeResponse.TransactionStatus.String(),
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     chargeResponse.VirtualAccountAction.Bank,
					VaNumber: chargeResponse.VirtualAccountAction.VirtualAccountNumber,
				},
			},
			FraudStatus: "accept",
			Currency:    "IDR",
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
	case primitive.PaymentTypeVirtualAccountBNI:
		responseBody, err := json.Marshal(schema.BNIVirtualAccountChargeSuccessResponse{
			StatusCode:        "201",
			StatusMessage:     "Success, Bank Transfer transaction is created",
			TransactionId:     chargeResponse.OrderId,
			OrderId:           chargeResponse.OrderId,
			GrossAmount:       strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			PaymentType:       chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:   chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus: chargeResponse.TransactionStatus.String(),
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     chargeResponse.VirtualAccountAction.Bank,
					VaNumber: chargeResponse.VirtualAccountAction.VirtualAccountNumber,
				},
			},
			FraudStatus: "accept",
			Currency:    "IDR",
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
	case primitive.PaymentTypeVirtualAccountPermata:
		responseBody, err := json.Marshal(schema.PermataVirtualAccountChargeSuccessResponse{
			StatusCode:        "201",
			StatusMessage:     "Success, PERMATA VA transaction is successful",
			TransactionId:     chargeResponse.OrderId,
			OrderId:           chargeResponse.OrderId,
			GrossAmount:       strconv.FormatInt(chargeResponse.TransactionAmount, 10),
			PaymentType:       chargeResponse.PaymentType.ToPaymentMethod(),
			TransactionTime:   chargeResponse.TransactionTime.Format(time.DateTime),
			TransactionStatus: chargeResponse.TransactionStatus.String(),
			FraudStatus:       "accept",
			PermataVaNumber:   chargeResponse.VirtualAccountAction.VirtualAccountNumber,
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

	// This shouldn't happen.
	// Running to this line of code means something's wrong when
	// validating the payment type (or processing the payment type
	// returned by the business logic).
	log.Error().Str("payment_type", chargeResponse.PaymentType.String()).Msg("unexpected payment type")
	w.WriteHeader(http.StatusInternalServerError)
}

func parseValidPaymentMethod(r schema.ChargeTransactionRequest) (primitive.PaymentType, error) {
	switch r.PaymentType {
	case "gopay":
		return primitive.PaymentTypeEMoneyGopay, nil
	case "shopeepay":
		return primitive.PaymentTypeEMoneyShopeePay, nil
	case "qris":
		return primitive.PaymentTypeEMoneyQRIS, nil
	case "bank_transfer":
		switch r.BankTransfer.Bank {
		case "bca":
			return primitive.PaymentTypeVirtualAccountBCA, nil
		case "bri":
			return primitive.PaymentTypeVirtualAccountBNI, nil
		case "permata":
			return primitive.PaymentTypeVirtualAccountPermata, nil
		case "bni":
			return primitive.PaymentTypeVirtualAccountBNI, nil
		default:
			return primitive.PaymentTypeUnspecified, fmt.Errorf("invalid bank name")
		}
	default:
		return primitive.PaymentTypeUnspecified, fmt.Errorf("unknown payment type")
	}
}
