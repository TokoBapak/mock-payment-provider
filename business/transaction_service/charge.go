package transaction_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"mock-payment-provider/presentation/schema"
	"mock-payment-provider/repository/signature"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d *Dependency) Charge(ctx context.Context, request business.ChargeRequest) (business.ChargeResponse, error) {
	// Validate the request payload
	if err := ValidateChargeRequest(request); err != nil {
		return business.ChargeResponse{}, err
	}

	// Validate the transaction amount and the amount of each product items
	var totalAmount int64
	for _, item := range request.ProductItems {
		totalAmount += item.Price * item.Quantity
	}

	if totalAmount != request.TransactionAmount {
		return business.ChargeResponse{}, business.ErrMismatchedTransactionAmount
	}

	transactionTime := time.Now()
	switch request.PaymentType {
	case primitive.PaymentTypeVirtualAccountBCA:
		fallthrough
	case primitive.PaymentTypeVirtualAccountPermata:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBRI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBNI:
		// Create new transaction
		expiredAt := time.Now().Add(time.Hour * 24)
		err := d.TransactionRepository.Create(
			ctx,
			repository.CreateTransactionParam{
				OrderID:     request.OrderId,
				Amount:      request.TransactionAmount,
				PaymentType: request.PaymentType,
				Status:      primitive.TransactionStatusPending,
				ExpiredAt:   expiredAt,
			},
		)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return business.ChargeResponse{}, business.ErrDuplicateOrderId
			}

			return business.ChargeResponse{}, fmt.Errorf("creating new transaction: %w", err)
		}

		// Acquire virtual account number from customer email
		virtualAccountNumber, err := d.VirtualAccountRepository.CreateOrGetVirtualAccountNumber(ctx, request.Customer.Email)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("acquiring virtual account number for %s: %w", request.Customer.Email, err)
		}

		// Create a virtual account entry
		_, err = d.VirtualAccountRepository.CreateCharge(
			ctx,
			virtualAccountNumber,
			request.OrderId,
			request.TransactionAmount,
			expiredAt,
		)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("creating virtual account entry: %w", err)
		}

		go func() {
			// Send a PENDING webhook
			payload, err := buildPendingWebhookMessage(pendingWebhookParameters{
				TransactionTime:      transactionTime,
				GrossAmount:          totalAmount,
				OrderId:              request.OrderId,
				PaymentType:          request.PaymentType,
				VirtualAccountNumber: virtualAccountNumber,
			})
			if err != nil {
				// TODO: properly log errors
				return
			}

			// Sleep for 10 seconds to make sure client has received the response
			time.Sleep(time.Second * 10)

			ctx := context.Background()

			err = d.WebhookClient.Send(ctx, payload)
			if err != nil {
				// TODO: properly log errors
			}
		}()

		go func(expiresAt time.Time, orderId string) {
			// Send a EXPIRED webhook

			time.Sleep(time.Until(expiredAt))
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			transaction, err := d.TransactionRepository.GetByOrderId(ctx, orderId)
			if err != nil {
				// TODO: properly log errors
				return
			}

			if transaction.TransactionStatus != primitive.TransactionStatusPending {
				// Do nothing
				return
			}

			// Update transaction status to expired
			err = d.TransactionRepository.UpdateStatus(ctx, orderId, primitive.TransactionStatusExpired)
			if err != nil {
				// TODO: properly log errors
				return
			}

			// Send webhook
			ctx = context.Background()

			payload, err := buildExpiredWebhookMessage(expiredWebhookParameters{
				TransactionTime: transactionTime,
				GrossAmount:     totalAmount,
				OrderId:         request.OrderId,
				PaymentType:     request.PaymentType,
			})
			if err != nil {
				// TODO: properly log errors
				return
			}

			err = d.WebhookClient.Send(ctx, payload)
			if err != nil {
				// TODO: properly log errors
			}
		}(expiredAt, request.OrderId)

		return business.ChargeResponse{
			OrderId:           request.OrderId,
			TransactionAmount: request.TransactionAmount,
			PaymentType:       request.PaymentType,
			TransactionStatus: primitive.TransactionStatusPending,
			TransactionTime:   time.Now(),
			EMoneyAction:      []business.EMoneyAction{},
			VirtualAccountAction: business.VirtualAccountAction{
				Bank:                 request.PaymentType.String(),
				VirtualAccountNumber: virtualAccountNumber,
			},
		}, nil
	case primitive.PaymentTypeEMoneyQRIS:
		fallthrough
	case primitive.PaymentTypeEMoneyGopay:
		fallthrough
	case primitive.PaymentTypeEMoneyShopeePay:
		// Create new transaction
		expiredAt := time.Now().Add(time.Hour * 3)
		err := d.TransactionRepository.Create(
			ctx,
			repository.CreateTransactionParam{
				OrderID:     request.OrderId,
				Amount:      request.TransactionAmount,
				PaymentType: request.PaymentType,
				Status:      primitive.TransactionStatusPending,
				ExpiredAt:   expiredAt,
			},
		)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return business.ChargeResponse{}, business.ErrDuplicateOrderId
			}

			return business.ChargeResponse{}, fmt.Errorf("creating new transaction: %w", err)
		}

		// Create e-money entry
		id, err := d.EMoneyRepository.CreateCharge(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			expiredAt,
		)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("creating e-money entry: %w", err)
		}

		go func() {
			payload, err := buildPendingWebhookMessage(pendingWebhookParameters{
				TransactionTime:      transactionTime,
				GrossAmount:          totalAmount,
				OrderId:              request.OrderId,
				PaymentType:          request.PaymentType,
				VirtualAccountNumber: "",
			})
			if err != nil {
				// TODO: properly log errors
				return
			}

			// Sleep for 10 seconds to make sure client has received the response
			time.Sleep(time.Second * 10)

			ctx := context.Background()

			err = d.WebhookClient.Send(ctx, payload)
			if err != nil {
				// TODO: properly log errors
			}
		}()

		go func(expiresAt time.Time, orderId string) {
			// Send a EXPIRED webhook

			time.Sleep(time.Until(expiredAt))
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			transaction, err := d.TransactionRepository.GetByOrderId(ctx, orderId)
			if err != nil {
				// TODO: properly log errors
				return
			}

			if transaction.TransactionStatus != primitive.TransactionStatusPending {
				// Do nothing
				return
			}

			// Update transaction status to expired
			err = d.TransactionRepository.UpdateStatus(ctx, orderId, primitive.TransactionStatusExpired)
			if err != nil {
				// TODO: properly log errors
				return
			}

			// Send webhook
			ctx = context.Background()

			payload, err := buildExpiredWebhookMessage(expiredWebhookParameters{
				TransactionTime: transactionTime,
				GrossAmount:     totalAmount,
				OrderId:         request.OrderId,
				PaymentType:     request.PaymentType,
			})
			if err != nil {
				// TODO: properly log errors
				return
			}

			err = d.WebhookClient.Send(ctx, payload)
			if err != nil {
				// TODO: properly log errors
			}
		}(expiredAt, request.OrderId)

		return business.ChargeResponse{
			OrderId:           request.OrderId,
			TransactionAmount: request.TransactionAmount,
			PaymentType:       request.PaymentType,
			TransactionStatus: primitive.TransactionStatusPending,
			TransactionTime:   time.Now(),
			EMoneyAction: []business.EMoneyAction{
				{
					EMoneyActionType: business.EMoneyActionTypeGenerateQRCode,
					Method:           "GET",
					URL:              "/e-money/" + id + "/pay",
				},
				{
					EMoneyActionType: business.EMoneyActionTypeStatus,
					Method:           "GET",
					URL:              "/e-money/" + id + "/status",
				},
				{
					EMoneyActionType: business.EMoneyActionTypeCancel,
					Method:           "POST",
					URL:              "/e-money/" + id + "/cancel",
				},
			},
			VirtualAccountAction: business.VirtualAccountAction{},
		}, nil
	case primitive.PaymentTypeUnspecified:
		fallthrough
	default:
		return business.ChargeResponse{}, fmt.Errorf("invalid payment type")
	}
}

func ValidateChargeRequest(request business.ChargeRequest) *business.RequestValidationError {
	var issues []business.RequestValidationIssue

	// validate payment_type
	if request.PaymentType == primitive.PaymentTypeUnspecified {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeInvalidValue,
			Field:   "payment_type",
			Message: "must be valid value",
		})
	}

	// validate transaction.order_id
	if request.OrderId == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "order_id",
			Message: "can not be empty",
		})
	}

	// valdiate transaction.amount
	if request.TransactionAmount <= 0 {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeInvalidValue,
			Field:   "amount",
			Message: "must be greater than 0",
		})
	}

	// validate transaction.currency
	if request.TransactionCurrency == primitive.CurrencyUnspecified {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeInvalidValue,
			Field:   "currency",
			Message: "must be valid value",
		})
	}

	// validate customer.first_name
	if request.Customer.FirstName == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.first_name",
			Message: "can not be empty",
		})
	} else {
		if len(request.Customer.FirstName) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.first_name",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate customer.last_name
	if len(request.Customer.LastName) > 255 {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeTooLong,
			Field:   "customer.last_name",
			Message: "maximum of 255 characters length",
		})
	}

	// validate customer.email
	if request.Customer.Email == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.email",
			Message: "can not be empty",
		})
	} else {
		if ok := primitive.EmailPattern.MatchString(request.Customer.Email); !ok {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "customer.email",
				Message: "must be valid email",
			})
		}

		if len(request.Customer.Email) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.email",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate customer.phone_number
	if request.Customer.PhoneNumber == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.phone_number",
			Message: "can not be empty",
		})
	} else {
		if ok := primitive.PhoneNumberPattern.MatchString(request.Customer.PhoneNumber); !ok {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "customer.phone_number",
				Message: "must be valid phone_number",
			})
		}

		if len(request.Customer.PhoneNumber) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.email",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate customer.billing_address.first_name
	if request.Customer.BillingAddress.FirstName == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.billing_address.first_name",
			Message: "can not be empty",
		})
	} else {
		if len(request.Customer.BillingAddress.FirstName) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.billing_address.first_name",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate customer.billing_address.last_name
	if len(request.Customer.BillingAddress.LastName) > 255 {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeTooLong,
			Field:   "customer.billing_address.last_name",
			Message: "maximum of 255 characters length",
		})
	}

	// validate customer.billing_address.email
	if request.Customer.BillingAddress.Email == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.phone_number",
			Message: "can not be empty",
		})
	} else {
		if ok := primitive.EmailPattern.MatchString(
			request.Customer.BillingAddress.Email,
		); !ok {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "customer.billing_address.email",
				Message: "must be a valid email",
			})
		}

		if len(request.Customer.BillingAddress.Email) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.billing_address.email",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate customer.billing_address.phone
	if request.Customer.BillingAddress.Phone == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.billing_address.phone",
			Message: "can not be empty",
		})
	} else {
		if ok := primitive.PhoneNumberPattern.MatchString(
			request.Customer.BillingAddress.Phone,
		); !ok {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "customer.billing_address.phone",
				Message: "must be a valid phone number",
			})
		}

		if len(request.Customer.BillingAddress.Phone) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.billing_address.phone",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate customer.bliing_address.address
	if request.Customer.BillingAddress.Address == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.billing_address.address",
			Message: "can not be empty",
		})
	} else {
		if len(request.Customer.BillingAddress.Address) > 500 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.billing_address.address",
				Message: "maximum of 500 characters length",
			})
		}
	}

	// validate customer.billing_address.postal_code
	if request.Customer.BillingAddress.PostalCode == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.billing_address.postal_code",
			Message: "can not be empty",
		})
	} else {
		if _, err := strconv.ParseUint(request.Customer.BillingAddress.PostalCode, 10, 64); err != nil {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "customer.billing_address.postal_code",
				Message: "must be a valid postal code",
			})
		}

		if len(request.Customer.BillingAddress.PostalCode) > 10 { // less than equal uint64 characters length
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.billing_address.postal_code",
				Message: "maximum of 10 characters length",
			})
		}
	}

	// validate customer.billing_address.country_code
	if request.Customer.BillingAddress.CountryCode == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "customer.billing_address.country_code",
			Message: "can not be empty",
		})
	} else {
		if _, err := strconv.ParseUint(request.Customer.BillingAddress.CountryCode, 10, 32); err != nil {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "customer.billing_address.country_code",
				Message: "must be a valid country code",
			})
		}

		if len(request.Customer.BillingAddress.CountryCode) > 5 { // less than equal uint32 characters length
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "customer.billing_address.country_code",
				Message: "maximum of 5 characters length",
			})
		}
	}

	// validate seller.first_name
	if request.Seller.FirstName == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "seller.first_name",
			Message: "can not be empty",
		})
	} else {
		if len(request.Seller.FirstName) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "seller.first_name",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate seller.last_name
	if len(request.Seller.LastName) > 255 {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeTooLong,
			Field:   "seller.last_name",
			Message: "maximum of 255 characters length",
		})
	}

	// validate seller.email
	if request.Seller.Email == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "seller.email",
			Message: "can not be empty",
		})
	} else {
		if ok := primitive.EmailPattern.MatchString(
			request.Seller.Email,
		); !ok {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "seller.email",
				Message: "must be a valid value",
			})
		}

		if len(request.Seller.Email) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "seller.email",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate seller.phone_number
	if request.Seller.PhoneNumber == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "seller.email",
			Message: "can not be empty",
		})
	} else {
		if ok := primitive.PhoneNumberPattern.MatchString(
			request.Seller.PhoneNumber,
		); !ok {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "seller.phone_number",
				Message: "must be a valid value",
			})
		}

		if len(request.Seller.PhoneNumber) > 255 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "seller.phone_number",
				Message: "maximum of 255 characters length",
			})
		}
	}

	// validate seller.address
	if request.Seller.Address == "" {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "seller.address",
			Message: "can not be empty",
		})
	} else {
		if len(request.Seller.Address) > 500 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeTooLong,
				Field:   "seller.address",
				Message: "maximum of 500 characters length",
			})
		}
	}

	// validate items
	if len(request.ProductItems) == 0 {
		issues = append(issues, business.RequestValidationIssue{
			Code:    business.RequestValidationCodeRequired,
			Field:   "items",
			Message: "can not be empty",
		})
	}

	// validate items
	for _, item := range request.ProductItems {
		// validate items.id
		if item.ID == "" {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeRequired,
				Field:   "items.id",
				Message: "can not be empty",
			})
		} else {
			if len(item.ID) > 255 {
				issues = append(issues, business.RequestValidationIssue{
					Code:    business.RequestValidationCodeTooLong,
					Field:   "items.id",
					Message: "maximum of 255 characters length",
				})
			}
		}

		// validate items.price
		if item.Price <= 0 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "items.price",
				Message: "must be greater than 0",
			})
		}

		// validate itmes.quantity
		if item.Quantity <= 0 {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeInvalidValue,
				Field:   "items.quantity",
				Message: "must be greater than 0",
			})
		}

		// validate items.name
		if item.Name == "" {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeRequired,
				Field:   "items.name",
				Message: "can not be empty",
			})
		} else {
			if len(item.Name) > 255 {
				issues = append(issues, business.RequestValidationIssue{
					Code:    business.RequestValidationCodeTooLong,
					Field:   "items.name",
					Message: "maximum of 255 characters length",
				})
			}

		}

		// validate items.category
		if item.Category == "" {
			issues = append(issues, business.RequestValidationIssue{
				Code:    business.RequestValidationCodeRequired,
				Field:   "items.category",
				Message: "can not be empty",
			})
		} else {
			if len(item.Category) > 255 {
				issues = append(issues, business.RequestValidationIssue{
					Code:    business.RequestValidationCodeTooLong,
					Field:   "items.category",
					Message: "maximum of 255 characters length",
				})
			}
		}
	}

	if len(issues) > 0 {
		return &business.RequestValidationError{Issues: issues}
	}

	return nil
}

type pendingWebhookParameters struct {
	TransactionTime      time.Time
	GrossAmount          int64
	OrderId              string
	PaymentType          primitive.PaymentType
	VirtualAccountNumber string
}

func buildPendingWebhookMessage(parameters pendingWebhookParameters) ([]byte, error) {
	// TODO: should we include server key?
	signatureKey := signature.Generate(parameters.OrderId, 200, parameters.GrossAmount, "")
	switch parameters.PaymentType {
	case primitive.PaymentTypeVirtualAccountBCA:
		return json.Marshal(schema.BCAVirtualAccountChargePendingResponse{
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     "bca",
					VaNumber: parameters.VirtualAccountNumber,
				},
			},
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.String(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusPending.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
		})
	case primitive.PaymentTypeVirtualAccountBRI:
		return json.Marshal(schema.BRIVirtualAccountChargePendingResponse{
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     "bri",
					VaNumber: parameters.VirtualAccountNumber,
				},
			},
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.String(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusPending.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
		})
	case primitive.PaymentTypeVirtualAccountBNI:
		return json.Marshal(schema.BNIVirtualAccountChargePendingResponse{
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     "bni",
					VaNumber: parameters.VirtualAccountNumber,
				},
			},
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.String(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusPending.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
		})
	case primitive.PaymentTypeVirtualAccountPermata:
		return json.Marshal(schema.PermataVirtualAccountChargePendingResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusPending.String(),
			FraudStatus:       "accept",
			PermataVaNumber:   parameters.VirtualAccountNumber,
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeEMoneyQRIS:
		return json.Marshal(schema.QRISChargePendingResponse{
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusPending.String(),
			TransactionId:     parameters.OrderId,
			StatusMessage:     "midtrans payment notification",
			StatusCode:        "200",
			SignatureKey:      signatureKey,
			PaymentType:       parameters.PaymentType.String(),
			OrderId:           parameters.OrderId,
			MerchantId:        "",
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			FraudStatus:       "accept",
			Currency:          "IDR",
			Acquirer:          "nobu",
		})
	case primitive.PaymentTypeEMoneyGopay:
		return json.Marshal(schema.GopayChargePendingResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusPending.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeEMoneyShopeePay:
		return json.Marshal(schema.ShopeePayChargePendingResponse{
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusPending.String(),
			TransactionId:     parameters.OrderId,
			StatusMessage:     "midtrans payment notification",
			StatusCode:        "200",
			SignatureKey:      signatureKey,
			PaymentType:       parameters.PaymentType.String(),
			OrderId:           parameters.OrderId,
			MerchantId:        "",
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			FraudStatus:       "accept",
			Currency:          "IDR",
		})
	default:
		return nil, fmt.Errorf("unknown payment type")
	}
}

type expiredWebhookParameters struct {
	TransactionTime time.Time
	GrossAmount     int64
	OrderId         string
	PaymentType     primitive.PaymentType
}

func buildExpiredWebhookMessage(parameters expiredWebhookParameters) ([]byte, error) {
	// TODO: should we include server key?
	signatureKey := signature.Generate(parameters.OrderId, 200, parameters.GrossAmount, "")
	switch parameters.PaymentType {
	case primitive.PaymentTypeVirtualAccountBCA:
		return json.Marshal(schema.BCAVirtualAccountChargeExpiredResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeVirtualAccountBRI:
		return json.Marshal(schema.BRIVirtualAccountChargeExpiredResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeVirtualAccountBNI:
		return json.Marshal(schema.BRIVirtualAccountChargeExpiredResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeVirtualAccountPermata:
		return json.Marshal(schema.PermataVirtualAccountChargeExpiredResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeEMoneyQRIS:
		return json.Marshal(schema.QRISChargeExpiredResponse{
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			TransactionId:     parameters.OrderId,
			StatusMessage:     "midtrans payment notification",
			StatusCode:        "200",
			SignatureKey:      signatureKey,
			PaymentType:       parameters.PaymentType.String(),
			OrderId:           parameters.OrderId,
			MerchantId:        "",
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			FraudStatus:       "accept",
			Currency:          "IDR",
			Acquirer:          "nobu",
		})
	case primitive.PaymentTypeEMoneyGopay:
		return json.Marshal(schema.GopayChargeExpiredResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.String(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeEMoneyShopeePay:
		return json.Marshal(schema.ShopeePayChargeExpiredResponse{
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusExpired.String(),
			TransactionId:     parameters.OrderId,
			StatusMessage:     "midtrans payment notification",
			StatusCode:        "200",
			SignatureKey:      signatureKey,
			PaymentType:       parameters.PaymentType.String(),
			OrderId:           parameters.OrderId,
			MerchantId:        "",
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			FraudStatus:       "accept",
			Currency:          "IDR",
		})
	default:
		return nil, fmt.Errorf("unknown payment type")
	}
}
