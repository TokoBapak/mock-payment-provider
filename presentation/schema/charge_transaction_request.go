package schema

type ChargeTransactionRequest struct {
	PaymentType        string `json:"payment_type"`
	TransactionDetails struct {
		OrderId     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
		Currency    string `json:"currency"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Email          string `json:"email"`
		PhoneNumber    string `json:"phone_number"`
		BillingAddress struct {
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			Phone       string `json:"phone"`
			Address     string `json:"address"`
			PostalCode  string `json:"postal_code"`
			CountryCode string `json:"country_code"`
		} `json:"billing_address"`
	} `json:"customer_details"`
	Seller struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
	} `json:"seller"`
	ItemDetails []struct {
		Id       string `json:"id"`
		Price    int64  `json:"price"`
		Quantity int    `json:"quantity"`
		Name     string `json:"name"`
		Category string `json:"category"`
	} `json:"item_details"`
	QRIS struct {
		Acquirer string `json:"acquirer"`
	}
	Gopay struct {
		EnableCallback bool   `json:"enable_callback"`
		CallbackURL    string `json:"callback_url"`
	}
	ShopeePay struct {
		CallbackURL string `json:"callback_url"`
	}
	BankTransfer struct {
		Bank    string `json:"bank"`
		Permata struct {
			RecipientName string `json:"recipient_name"`
		} `json:"permata"`
		VirtualAccountNumber string `json:"va_number"`
		FreeText             struct {
			Inquiry []struct {
				Indonesian string `json:"id"`
				English    string `json:"english"`
			} `json:"inquiry"`
			Payment []struct {
				Indonesian string `json:"id"`
				English    string `json:"english"`
			} `json:"payment"`
		} `json:"free_text"`
	} `json:"bank_transfer"`
	BCA struct {
		SubCompanyCode string `json:"sub_company_code"`
	} `json:"bca"`
}
