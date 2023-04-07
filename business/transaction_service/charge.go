package transaction_service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d Dependency) Charge(ctx context.Context, request business.ChargeRequest) (business.ChargeResponse, error) {
	// Validate the request payload
	if err := ValidateChageRequest(request); err != nil {
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

	switch request.PaymentType {
	case primitive.PaymentTypeVirtualAccountBCA:
		fallthrough
	case primitive.PaymentTypeVirtualAccountPermata:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBRI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBNI:
		// Create new transaction
		err := d.TransactionRepository.Create(
			ctx,
			repository.CreateTransactionParam{
				OrderID:     request.OrderId,
				Amount:      request.TransactionAmount,
				PaymentType: request.PaymentType,
				Status:      primitive.TransactionStatusPending,
				ExpiredAt:   time.Now().Add(time.Hour * 24),
			},
		)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return business.ChargeResponse{}, business.ErrDuplicateOrderId
			}

			return business.ChargeResponse{}, fmt.Errorf("creating new transaction: %w", err)
		}

		// Create a virtual account entry
		virtualAccountNumber, err := d.VirtualAccountRepository.CreateCharge(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			time.Now().Add(time.Hour*24),
		)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("creating virtual account entry: %w", err)
		}

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
		err := d.TransactionRepository.Create(
			ctx,
			repository.CreateTransactionParam{
				OrderID:     request.OrderId,
				Amount:      request.TransactionAmount,
				PaymentType: request.PaymentType,
				Status:      primitive.TransactionStatusPending,
				ExpiredAt:   time.Now().Add(time.Hour * 3),
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
			time.Now().Add(time.Hour*3),
		)
		if err != nil {
			return business.ChargeResponse{}, fmt.Errorf("creating e-money entry: %w", err)
		}

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

func ValidateChageRequest(request business.ChargeRequest) *business.RequestValidationError {
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
