package transaction_service

import (
	"context"
	"fmt"
	"regexp"

	"mock-payment-provider/business"
	"mock-payment-provider/primitive"
)

func (d Dependency) Charge(ctx context.Context, request business.ChargeRequest) (business.ChargeResponse, error) {
	// valiate the request
	if err := ValidateChageRequest(request); err != nil {
		return business.ChargeResponse{}, err
	}

	// TODO implement me
	panic("implement me")
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
	if request.TransactionCurrency == 0 {
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

	if ok := regexp.MustCompile(primitive.EmailPattern).Match([]byte(request.Customer.Email)); !ok {
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

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).Match([]byte(request.Customer.PhoneNumber)); !ok {
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

	if ok := regexp.MustCompile(primitive.EmailPattern).Match(
		[]byte(request.Customer.BillingAddress.Email),
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

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).Match(
		[]byte(request.Customer.BillingAddress.Phone),
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

	if len(request.Customer.BillingAddress.PostalCode) > 10 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.postal_code must be less than 10 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.PostalCodePattern).Match(
		[]byte(request.Customer.BillingAddress.PostalCode),
	); !ok {
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

	if len(request.Customer.BillingAddress.CountryCode) > 5 {
		return &business.RequestValidationError{
			Reason: "customer.billing_address.country_code must be less than 5 characters",
		}
	}

	if ok := regexp.MustCompile(primitive.CountryCodePattern).Match(
		[]byte(request.Customer.BillingAddress.CountryCode),
	); !ok {
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

	if ok := regexp.MustCompile(primitive.EmailPattern).Match(
		[]byte(request.Seller.Email),
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

	if ok := regexp.MustCompile(primitive.PhoneNumberPattern).Match(
		[]byte(request.Seller.PhoneNumber),
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
