package transaction_service_test

import (
	"errors"
	"mock-payment-provider/business"
	"mock-payment-provider/business/transaction_service"
	"testing"
)

func TestValidateChargeRequest(t *testing.T) {
	type test struct {
		param     business.ChargeRequest
		expectErr *business.RequestValidationError
	}

	// provide correct value
	request := business.ChargeRequest{
		PaymentType:         1,
		OrderId:             "20220103",
		TransactionAmount:   50000,
		TransactionCurrency: 1,
		Customer: business.CustomerInformation{
			FirstName:   "tony",
			LastName:    "stark",
			Email:       "tonystark01@email.com",
			PhoneNumber: "+62123456789",
			BillingAddress: business.Address{
				FirstName:   "tony",
				LastName:    "stark",
				Email:       "tonystark01@email.com",
				Phone:       "+628123456789",
				Address:     "Jl. Kenangan",
				PostalCode:  "55123",
				CountryCode: "62",
			},
		},
		Seller: business.SellerInformation{
			FirstName:   "tom",
			LastName:    "holland",
			Email:       "tomholland01@email.com",
			PhoneNumber: "+62123456780",
			Address:     "Jl. Nin Aja Dulu",
		},
		ProductItems: []business.ProductItem{
			{
				ID:       "A123",
				Price:    25000,
				Quantity: 1,
				Name:     "Keyboard",
				Category: "Electronic",
			},
			{
				ID:       "A124",
				Price:    25000,
				Quantity: 1,
				Name:     "Mouse",
				Category: "Electronic",
			},
		},
	}

	// test positive case
	t.Run("positive test case", func(t *testing.T) {
		// action
		err := transaction_service.ValidateChageRequest(request)

		// assert
		if err != nil {
			t.Errorf("expect error nil, but got %T instead", err)
		}
	})

	// test payment_type
	t.Run("payment_type", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		// mock
		mock.PaymentType = 0

		// action
		err := transaction_service.ValidateChageRequest(mock)

		// assert
		if err == nil {
			t.Errorf("expect errors as *business.RequestValidationError when the given payment_type is invalid, instead got %T", err)
		}
		if !errors.As(err, &requestValidationError) {
			t.Errorf("expect errors as *business.RequestValidationError when the given payment_type is invalid, instead got %T", err)
		}
	})

	// test order_id
	t.Run("order_id", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			// empty string
			mock.OrderId = ""
			err := transaction_service.ValidateChageRequest(mock)

			// assert
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError when the given OrderId is empty string, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError when the given OrderId is empty string, instead got %T", err)
			}
		})
	})

	// test transaction_amount
	t.Run("transaction_amount", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("greater than 0", func(t *testing.T) {
			// less than 0
			mock.TransactionAmount = -1
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given TransactionAmount less than 0, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given TransactionAmount less than 0, instead got %T", err)
			}

			// equal 0
			mock.TransactionAmount = 0
			err = transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given TransactionAmount less than 0, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given TransactionAmount less than 0, instead got %T", err)
			}
		})
	})

	// test transaction_currency
	t.Run("transaction_currency", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("invalid value", func(t *testing.T) {
			mock.TransactionCurrency = 0
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given TarnsactionCurrency is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given TarnsactionCurrency is invalid, instead got %T", err)
			}
		})
	})

	// test Customer.FirstName
	t.Run("Customer.FirstName", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.FirstName = ""
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.FirstName is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.FirstName is empty, instead got %T", err)
			}
		})

		t.Run("less than 255 characters", func(t *testing.T) {
			// generate more than 255 characters
			for i := 0; i < 260; i++ {
				mock.Customer.FirstName += "a"
			}
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.FirstName greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.FirstName greater than 255 characters, instead got %T", err)
			}
		})
	})

	// test Customer.LastName
	t.Run("Customer.LastName", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("less than 255 characters", func(t *testing.T) {
			// generate more than 255 characters
			for i := 0; i < 260; i++ {
				mock.Customer.LastName += "a"
			}
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.LastName greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.LastName greater than 255 characters, instead got %T", err)
			}
		})
	})

	// test Customer.Email
	t.Run("Customer.Email", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.Email = ""
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is empty, instead got %T", err)
			}
		})

		t.Run("less than 255", func(t *testing.T) {
			for i := 0; i < 260; i++ {
				mock.Customer.Email += "a"
			}
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is greater than 255 characters, instead got %T", err)
			}
		})

		t.Run("valid email", func(t *testing.T) {
			mock.Customer.Email = "aaa@.com"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}

			mock.Customer.Email = "aaa.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}

			mock.Customer.Email = "aaa@com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}

			mock.Customer.Email = "@fefe@email.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.Email is invalid, instead got %T", err)
			}

			mock.Customer.Email = "@fefe@aaa.bbb.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors *business.RequestValidationError"+
					"when the given Customer.Email is valid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors *business.RequestValidationError"+
					"when the given Customer.Email is valid, instead got %T", err)
			}

			mock.Customer.Email = "fefe@aaa.bbb.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err != nil {
				t.Errorf("expect errors nil"+
					"when the given Customer.Email is valid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors nil"+
					"when the given Customer.Email is valid, instead got %T", err)
			}
		})
	})

	// test Customer.PhoneNumber
	t.Run("Customer.PhoneNumber", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.PhoneNumber = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.PhoneNumber is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.PhoneNumber is empty, instead got %T", err)
			}
		})

		t.Run("less than 255", func(t *testing.T) {
			for i := 0; i < 260; i++ {
				mock.Customer.PhoneNumber += "0"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.PhoneNumber is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.PhoneNumber is greater than 255 characters, instead got %T", err)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			mock.Customer.PhoneNumber = "123456"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.PhoneNumber is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.PhoneNumber is invalid, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.FirstName
	t.Run("Customer.BillingAddress.FirstName", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.BillingAddress.FirstName = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.FirstName is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.FirstName is empty, instead got %T", err)
			}
		})

		t.Run("less than 255 characters", func(t *testing.T) {
			for i := 0; i <= 26; i++ {
				mock.Customer.BillingAddress.FirstName += "aaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.FirstName is greater than 255 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.FirstName is greater than 255 characters length, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.LastName
	t.Run("Customer.BillingAddress.LastName", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("less than 255 characters", func(t *testing.T) {
			for i := 0; i <= 26; i++ {
				mock.Customer.BillingAddress.LastName += "aaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.LastName is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.LastName is greater than 255 characters, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.Email
	t.Run("Customer.BillingAddress.Email", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.BillingAddress.Email = ""
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is empty, instead got %T", err)
			}
		})

		t.Run("less than 255", func(t *testing.T) {
			for i := 0; i < 26; i++ {
				mock.Customer.BillingAddress.Email += "aaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is greater than 255 characters, instead got %T", err)
			}
		})

		t.Run("valid email", func(t *testing.T) {
			mock.Customer.BillingAddress.Email = "aaa@.com"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}

			mock.Customer.BillingAddress.Email = "aaa.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}

			mock.Customer.BillingAddress.Email = "aaa@com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}

			mock.Customer.BillingAddress.Email = "@fefe@email.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Email is invalid, instead got %T", err)
			}

			mock.Customer.BillingAddress.Email = "@fefe@aaa.bbb.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors *business.RequestValidationErrorl"+
					"when the given Customer.BillingAddress.Email is valid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors *business.RequestValidationErrorl"+
					"when the given Customer.BillingAddress.Email is valid, instead got %T", err)
			}

			mock.Customer.BillingAddress.Email = "fefe@aaa.bbb.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err != nil {
				t.Errorf("expect errors nil"+
					"when the given Customer.BillingAddress.Email is valid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors nil"+
					"when the given Customer.BillingAddress.Email is valid, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.Phone
	t.Run("Customer.BillingAddress.Phone", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.BillingAddress.Phone = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Phone is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Phone is empty, instead got %T", err)
			}
		})

		t.Run("less than 255", func(t *testing.T) {
			for i := 0; i < 260; i++ {
				mock.Customer.BillingAddress.Phone += "0"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Phone is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Phone is greater than 255 characters, instead got %T", err)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			mock.Customer.BillingAddress.Phone = "123456"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Phone is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Phone is invalid, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.Address
	t.Run("Customer.BillingAddress.Address", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.BillingAddress.Address = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Address is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Address is empty, instead got %T", err)
			}
		})

		t.Run("less than 500 characters", func(t *testing.T) {
			for i := 0; i <= 50; i++ {
				mock.Customer.BillingAddress.Address += "aaaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Address is greater than 500 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.Address is greater than 500 characters length, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.PostalCode
	t.Run("Customer.BillingAddress.PostalCode", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.BillingAddress.PostalCode = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					" when the given Customer.BillingAddress.PostalCode is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					" when the given Customer.BillingAddress.PostalCode is empty, instead got %T", err)
			}
		})

		t.Run("less than 10 characters length", func(t *testing.T) {
			mock.Customer.BillingAddress.PostalCode = "123456789123456789"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.PostalCode is greater than 10 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.PostalCode is greater than 10 characters length, instead got %T", err)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			mock.Customer.BillingAddress.PostalCode = "abcdefg"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.PostalCode is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.PostalCode is invalid, instead got %T", err)
			}
		})
	})

	// test Customer.BillingAddress.CountryCode
	t.Run("Customer.BillingAddress.CountryCode", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Customer.BillingAddress.CountryCode = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.CountryCode is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.CountryCode is empty, instead got %T", err)
			}
		})

		t.Run("less than 5 characters length", func(t *testing.T) {
			mock.Customer.BillingAddress.CountryCode = "123456789"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.CountryCode is greater than 5 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.CountryCode is greater than 5 characters length, instead got %T", err)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			mock.Customer.BillingAddress.CountryCode = "abc"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.CountryCode is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Customer.BillingAddress.CountryCode is invalid, instead got %T", err)
			}
		})
	})

	// test Seller.FirstName
	t.Run("Seller.FirstName", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Seller.FirstName = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.FirstName is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.FirstName is empty, instead got %T", err)
			}
		})

		t.Run("less than 255 characters length", func(t *testing.T) {
			for i := 0; i <= 26; i++ {
				mock.Seller.FirstName += "aaaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.FirstName is greater than 255 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.FirstName is greater than 255 characters length, instead got %T", err)
			}
		})

	})

	// test Seller.LastName
	t.Run("Seller.LastName", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("less than 255 characters length", func(t *testing.T) {
			for i := 0; i <= 26; i++ {
				mock.Seller.LastName += "aaaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.LastName is greater than 255 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.LastName is greater than 255 characters length, instead got %T", err)
			}
		})
	})

	// test Seller.Email
	t.Run("Seller.Email", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Seller.Email = ""
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is empty, instead got %T", err)
			}
		})

		t.Run("less than 255", func(t *testing.T) {
			for i := 0; i < 26; i++ {
				mock.Seller.Email += "aaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)

			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is greater than 255 characters, instead got %T", err)
			}
		})

		t.Run("valid email", func(t *testing.T) {
			mock.Seller.Email = "aaa@.com"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}

			mock.Seller.Email = "aaa.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}

			mock.Seller.Email = "aaa@com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}

			mock.Seller.Email = "@fefe@email.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Email is invalid, instead got %T", err)
			}

			mock.Seller.Email = "@fefe@aaa.bbb.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors *business.RequestValidationError"+
					"when the given Seller.Email is valid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors *business.RequestValidationError"+
					"when the given Seller.Email is valid, instead got %T", err)
			}

			mock.Seller.Email = "fefe@aaa.bbb.com"
			err = transaction_service.ValidateChageRequest(mock)
			if err != nil {
				t.Errorf("expect errors nil"+
					"when the given Seller.Email is valid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors nil"+
					"when the given Seller.Email is valid, instead got %T", err)
			}
		})
	})

	// test Seller.PhoneNumber
	t.Run("Seller.PhoneNumber", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Seller.PhoneNumber = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.PhoneNumber is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.PhoneNumber is empty, instead got %T", err)
			}
		})

		t.Run("less than 255", func(t *testing.T) {
			for i := 0; i < 260; i++ {
				mock.Seller.PhoneNumber += "0"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.PhoneNumber is greater than 255 characters, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.PhoneNumber is greater than 255 characters, instead got %T", err)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			mock.Seller.PhoneNumber = "123456"
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.PhoneNumber is invalid, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.PhoneNumber is invalid, instead got %T", err)
			}
		})
	})

	// test Seller.Address
	t.Run("Seller.Address", func(t *testing.T) {
		// arrange
		mock := request
		var requestValidationError *business.RequestValidationError

		t.Run("required", func(t *testing.T) {
			mock.Seller.Address = ""
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Address is empty, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Address is empty, instead got %T", err)
			}
		})

		t.Run("less than 500 characters", func(t *testing.T) {
			for i := 0; i <= 50; i++ {
				mock.Seller.Address += "aaaaaaaaaaa"
			}
			err := transaction_service.ValidateChageRequest(mock)
			if err == nil {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Address is greater than 500 characters length, instead got %T", err)
			}
			if !errors.As(err, &requestValidationError) {
				t.Errorf("expect errors as *business.RequestValidationError"+
					"when the given Seller.Address is greater than 500 characters length, instead got %T", err)
			}
		})
	})

	// test ProductItems
	t.Run("ProductItems", func(t *testing.T) {

		t.Run("ProductItems[0].ID", func(t *testing.T) {
			// arrange
			mock := request
			var requestValidationError *business.RequestValidationError
			t.Run("required", func(t *testing.T) {
				mock.ProductItems[0].ID = ""
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].ID is empty, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].ID is empty, instead got %T", err)
				}
			})
		})

		t.Run("ProductItems[0].Price", func(t *testing.T) {
			// arrange
			mock := request
			var requestValidationError *business.RequestValidationError

			t.Run("greater than 0", func(t *testing.T) {
				mock.ProductItems[0].Price = 0
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Price is 0, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Price is 0, instead got %T", err)
				}

				mock.ProductItems[0].Price = -1
				err = transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].ID is less than 0, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].ID is less than 0, instead got %T", err)
				}
			})
		})

		t.Run("ProductItems[0].Quantity", func(t *testing.T) {
			// arrange
			mock := request
			var requestValidationError *business.RequestValidationError

			t.Run("greater than 0", func(t *testing.T) {
				mock.ProductItems[0].Quantity = 0
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Quantity is 0, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Quantity is 0, instead got %T", err)
				}

				mock.ProductItems[0].Quantity = -1
				err = transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Quantity less than 0, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Quantity less than 0, instead got %T", err)
				}
			})
		})

		t.Run("ProducItems[0].Name", func(t *testing.T) {
			// arrange
			mock := request
			var requestValidationError *business.RequestValidationError

			t.Run("required", func(t *testing.T) {
				mock.ProductItems[0].Name = ""
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Name is empty, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Name is empty, instead got %T", err)
				}
			})

			t.Run("less than 255 characters length", func(t *testing.T) {
				for i := 0; i <= 26; i++ {
					mock.ProductItems[0].Name += "aaaaaaaaaaa"
				}
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Name is greater than 255 characters length, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Name is greater than 255 characters length, instead got %T", err)
				}
			})
		})

		t.Run("ProductItems[0].Cateogry", func(t *testing.T) {
			// arrange
			mock := request
			var requestValidationError *business.RequestValidationError

			t.Run("required", func(t *testing.T) {
				mock.ProductItems[0].Category = ""
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Category is empty, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Category is empty, instead got %T", err)
				}
			})

			t.Run("less than 255 characters length", func(t *testing.T) {
				for i := 0; i <= 26; i++ {
					mock.ProductItems[0].Category += "aaaaaaaaaaa"
				}
				err := transaction_service.ValidateChageRequest(mock)
				if err == nil {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Category is greater than 255 characters length, instead got %T", err)
				}
				if !errors.As(err, &requestValidationError) {
					t.Errorf("expect errors as *business.RequestValidatorError"+
						"when the given ProductItems[0].Category is greater than 255 characters length, instead got %T", err)
				}
			})
		})
	})

}
