package schema

type TransactionStatusResponse struct {
	StatusCode               string `json:"status_code"`
	StatusMessage            string `json:"status_message"`
	TransactionId            string `json:"transaction_id"`
	MaskedCard               string `json:"masked_card"`
	OrderId                  string `json:"order_id"`
	PaymentType              string `json:"payment_type"`
	TransactionTime          string `json:"transaction_time"`
	TransactionStatus        string `json:"transaction_status"`
	FraudStatus              string `json:"fraud_status"`
	ApprovalCode             string `json:"approval_code"`
	SignatureKey             string `json:"signature_key"`
	Bank                     string `json:"bank"`
	GrossAmount              int64  `json:"gross_amount,string"`
	ChannelResponseCode      string `json:"channel_response_code"`
	ChannelResponseMessage   string `json:"channel_response_message"`
	CardType                 string `json:"card_type"`
	PaymentOptionType        string `json:"payment_option_type"`
	ShopeepayReferenceNumber string `json:"shopeepay_reference_number"`
	ReferenceId              string `json:"reference_id"`
}
