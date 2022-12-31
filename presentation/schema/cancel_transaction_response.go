package schema

type CancelTransactionResponse struct {
	StatusCode        int    `json:"status_code,string"`
	StatusMessage     string `json:"status_message"`
	OrderId           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
}
