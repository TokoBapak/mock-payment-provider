package schema

type BCAVirtualAccountChargeSuccessResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	VaNumbers         []struct {
		Bank     string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
	FraudStatus string `json:"fraud_status"`
	Currency    string `json:"currency"`
}

type BCAVirtualAccountChargePendingResponse struct {
	VaNumbers []struct {
		Bank     string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
	TransactionTime   string `json:"transaction_time"`
	GrossAmount       string `json:"gross_amount"`
	OrderId           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	TransactionId     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	StatusMessage     string `json:"status_message"`
}

type BCAVirtualAccountChargeSettlementResponse struct {
	VaNumbers []struct {
		Bank     string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
	TransactionTime   string `json:"transaction_time"`
	GrossAmount       string `json:"gross_amount"`
	OrderId           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	TransactionId     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	StatusMessage     string `json:"status_message"`
}

type BCAVirtualAccountChargeExpiredResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	SignatureKey      string `json:"signature_key"`
}
