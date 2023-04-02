package schema

type PermataVirtualAccountChargeSuccessResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PermataVaNumber   string `json:"permata_va_number"`
}

type PermataVirtualAccountChargePendingResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PermataVaNumber   string `json:"permata_va_number"`
	SignatureKey      string `json:"signature_key"`
}

type PermataVirtualAccountChargeSettlementResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PermataVaNumber   string `json:"permata_va_number"`
	SignatureKey      string `json:"signature_key"`
}

type PermataVirtualAccountChargeExpiredResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PermataVaNumber   string `json:"permata_va_number"`
	SignatureKey      string `json:"signature_key"`
}
