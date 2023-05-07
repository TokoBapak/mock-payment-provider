package schema

type InternalTransactionDetailResponse struct {
	OrderId              string `json:"order_id"`
	ChargedAmount        int64  `json:"charged_amount"`
	TransactionStatus    string `json:"transaction_status"`
	PaymentMethod        string `json:"payment_method"`
	Bank                 string `json:"bank"`
	VirtualAccountNumber string `json:"virtual_account_number,omitempty"`
	EMoneyId             string `json:"e_money_id,omitempty"`
}
