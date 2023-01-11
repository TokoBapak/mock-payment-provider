package transaction_service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
)

func (d Dependency) Charge(ctx context.Context, request business.ChargeRequest) (business.ChargeResponse, error) {
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
	case primitive.PaymentTypeVirtualAccountMandiri:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBRI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBNI:
		// Create new transaction
		err := d.TransactionRepository.Create(
			ctx,
			request.OrderId,
			request.TransactionAmount,
			request.PaymentType,
			primitive.TransactionStatusPending,
			time.Now().Add(time.Hour*24),
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
			request.OrderId,
			request.TransactionAmount,
			request.PaymentType,
			primitive.TransactionStatusPending,
			time.Now().Add(time.Hour*3),
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
					EMoneyActionType: business.EMoneyActionTypePay,
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
	// validate payment_type
	if request.PaymentType == primitive.PaymentTypeUnspecified {
		return &business.RequestValidationError{
			Reason: "payment_type is not valid",
		}
	}

	// validate transaction.order_id
	if request.OrderId == "" {
		return &business.RequestValidationError{
			Reason: "order_id is required",
		}
	}

	// valdiate transaction.amount
	if request.TransactionAmount <= 0 {
		return &business.RequestValidationError{
			Reason: "amount should be greater than 0",
		}
	}

	// validate transaction.currency
	if request.TransactionCurrency == primitive.CurrencyUnspecified {
		return &business.RequestValidationError{
			Reason: "transaction.currency is not valid",
		}
	}

	// validate customer.first_name
	if request.Customer.FirstName == "" {
		return &business.RequestValidationError{
			Reason: "customer.first_name is required",
		}
	}
	if len(request.Customer.FirstName) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.first_name must be less than 255 characters",
		}
	}

	// validate customer.last_name
	if len(request.Customer.LastName) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.last_name must be less than 255 characters",
		}
	}

	// validate customer.email
	if request.Customer.Email == "" {
		return &business.RequestValidationError{
			Reason: "customer.email is required",
		}
	}

	if len(request.Customer.Email) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.email must be less than 255 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).MatchString(request.Customer.Email); !ok {
		return &business.RequestValidationError{
			Reason: "customer.email is not valid",
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).MatchString(request.Customer.Email); !ok {
		return &business.RequestValidationError{
			Reason: "customer.email is not valid",
		}
	}

	// validate customer.phone_number
	if request.Customer.PhoneNumber == "" {
		return &business.RequestValidationError{
			Reason: "customer.phone_number is required",
		}
	}

	if len(request.Customer.PhoneNumber) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.phone_number must less than 255 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).MatchString(request.Customer.PhoneNumber); !ok {
		return &business.RequestValidationError{
			Reason: "customer.phone_number is not valid",
		}
	}

	// validate customer.billing_address.first_name
	if request.Customer.BillingAddress.FirstName == "" {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.first_name is required",
		}
	}

	if len(request.Customer.BillingAddress.FirstName) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.first_name must be less than 255 characters",
		}
	}

	// validate customer.billing_address.last_name
	if len(request.Customer.BillingAddress.LastName) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.last_name must be less than 255 characters",
		}
	}

	// validate customer.billing_address.email
	if request.Customer.BillingAddress.Email == "" {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.email is required",
		}
	}

	if len(request.Customer.BillingAddress.Email) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.email must be less than 255 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).MatchString(
		request.Customer.BillingAddress.Email,
	); !ok {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.email is not valid",
		}
	}

	// validate customer.billing_address.phone
	if request.Customer.BillingAddress.Phone == "" {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.phone is required",
		}
	}

	if len(request.Customer.BillingAddress.Phone) > 255 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.phone must be less than 255 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).MatchString(
		request.Customer.BillingAddress.Phone,
	); !ok {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.phone is not valid",
		}
	}

	// validate customer.bliing_address.address
	if request.Customer.BillingAddress.Address == "" {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.address is required",
		}
	}

	if len(request.Customer.BillingAddress.Address) > 500 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.address must be less than 500 characters",
		}
	}

	// validate customer.billing_address.postal_code
	if request.Customer.BillingAddress.PostalCode == "" {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.postal_code is required",
		}
	}

	if len(request.Customer.BillingAddress.PostalCode) > 10 { // less than equal uint64 characters length
		return &business.RequestValidationError{
			Reason: "customer.billing_address.postal_code must be less than 10 characters",
		}
	}

	if _, err := strconv.ParseUint(request.Customer.BillingAddress.PostalCode, 10, 64); err != nil {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.postal_code is not valid",
		}
	}

	// validate customer.billing_address.country_code
	if request.Customer.BillingAddress.CountryCode == "" {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.country_code is required",
		}
	}

	if len(request.Customer.BillingAddress.CountryCode) > 5 { // less than equal uint32 characters length
		return &business.RequestValidationError{
			Reason: "customer.billing_address.country_code must be less than 5 characters",
		}
	}

	if _, err := strconv.ParseUint(request.Customer.BillingAddress.CountryCode, 10, 32); err != nil {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.country_code is not valid",
		}
	}

	// validate seller.first_name
	if request.Seller.FirstName == "" {
		return &business.RequestValidationError{
			Reason: "seller.first_name is required",
		}
	}

	if len(request.Seller.FirstName) > 255 {
		return &business.RequestValidationError{
			Reason: "seller.first_name must be less than 255 characters",
		}
	}

	// validate seller.last_name
	if len(request.Seller.LastName) > 255 {
		return &business.RequestValidationError{
			Reason: "seller.last_name must be less than 255 characters",
		}
	}

	// validate seller.email
	if request.Seller.Email == "" {
		return &business.RequestValidationError{
			Reason: "seller.email is required",
		}
	}

	if len(request.Seller.Email) > 255 {
		return &business.RequestValidationError{
			Reason: "seller.email must be less than 255 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.EmailPattern).MatchString(
		request.Seller.Email,
	); !ok {
		return &business.RequestValidationError{
			Reason: "seller.email is not valid",
		}
	}

	// validate seller.phone_number
	if request.Seller.PhoneNumber == "" {
		return &business.RequestValidationError{
			Reason: "seller.phone_number is required",
		}
	}

	if len(request.Seller.PhoneNumber) > 255 {
		return &business.RequestValidationError{
			Reason: "seller.phone_number must less than 255 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).MatchString(
		request.Seller.PhoneNumber,
	); !ok {
		return &business.RequestValidationError{
			Reason: "seller.phone_number is not valid",
		}
	}

	// validate seller.address
	if request.Seller.Address == "" {
		return &business.RequestValidationError{
			Reason: "seller.address is required",
		}
	}

	if len(request.Seller.Address) > 500 {
		return &business.RequestValidationError{
			Reason: "seller.address must be less than 500 characters",
		}
	}

	// validate items.id
	if len(request.ProductItems) == 0 {
		return &business.RequestValidationError{
			Reason: "items.request_body must be greater than 0 length",
		}
	}

	// validate items
	for i, item := range request.ProductItems {
		// validate items.id
		if item.ID == "" {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.id is required", i),
			}
		}

		if len(item.ID) > 255 {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.id must be less than 255 characters", i),
			}
		}

		// validate items.price
		if item.Price <= 0 {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.price must be greather than 0", i),
			}
		}

		// validate itmes.quantity
		if item.Quantity <= 0 {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.quantity must be greater than 0", i),
			}
		}

		// validate items.name
		if item.Name == "" {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.quantity is required", i),
			}
		}

		if len(item.Name) > 255 {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.name must be less than 255 characters", i),
			}
		}

		// validate items.category
		if item.Category == "" {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.category is required", i),
			}
		}

		if len(item.Category) > 255 {
			return &business.RequestValidationError{
				Reason: fmt.Sprintf("items.%d.category must be less than 255 characters", i),
			}
		}
	}
	return nil
}
