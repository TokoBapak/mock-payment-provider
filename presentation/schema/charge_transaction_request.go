package schema

type ChargeTransactionRequest struct {
	PaymentType string `json:"payment_type"`
	Transaction struct {
		OrderId  string `json:"order_id"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	} `json:"transaction"`
	Customer struct {
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
	} `json:"customer"`
	Seller struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
	} `json:"seller"`
	Items []struct {
		Id       string `json:"id"`
		Price    int    `json:"price"`
		Quantity int    `json:"quantity"`
		Name     string `json:"name"`
		Category string `json:"category"`
	} `json:"items"`
}
