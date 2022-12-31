package schema

type ChargeTransactionResponse struct {
	StatusCode        int    `json:"status_code,string"`
	StatusMessage     string `json:"status_message"`
	OrderId           string `json:"order_id"`
	GrossAmount       int64  `json:"gross_amount,string"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
}
