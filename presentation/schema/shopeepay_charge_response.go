package schema

type ShopeePayChargeSuccessResponse struct {
	StatusCode             string `json:"status_code"`
	StatusMessage          string `json:"status_message"`
	ChannelResponseCode    string `json:"channel_response_code"`
	ChannelResponseMessage string `json:"channel_response_message"`
	TransactionId          string `json:"transaction_id"`
	OrderId                string `json:"order_id"`
	MerchantId             string `json:"merchant_id"`
	GrossAmount            string `json:"gross_amount"`
	Currency               string `json:"currency"`
	PaymentType            string `json:"payment_type"`
	TransactionTime        string `json:"transaction_time"`
	TransactionStatus      string `json:"transaction_status"`
	FraudStatus            string `json:"fraud_status"`
	Actions                []struct {
		Name   string `json:"name"`
		Method string `json:"method"`
		Url    string `json:"url"`
	} `json:"actions"`
}

type ShopeePayChargePendingResponse struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionId     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderId           string `json:"order_id"`
	MerchantId        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type ShopeePayChargeSettlementResponse struct {
	TransactionTime          string `json:"transaction_time"`
	TransactionStatus        string `json:"transaction_status"`
	TransactionId            string `json:"transaction_id"`
	StatusMessage            string `json:"status_message"`
	StatusCode               string `json:"status_code"`
	SignatureKey             string `json:"signature_key"`
	SettlementTime           string `json:"settlement_time"`
	PaymentType              string `json:"payment_type"`
	OrderId                  string `json:"order_id"`
	MerchantId               string `json:"merchant_id"`
	GrossAmount              string `json:"gross_amount"`
	FraudStatus              string `json:"fraud_status"`
	Currency                 string `json:"currency"`
	ShopeepayReferenceNumber string `json:"shopeepay_reference_number"`
	ReferenceId              string `json:"reference_id"`
}

type ShopeePayChargeExpiredResponse struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionId     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderId           string `json:"order_id"`
	MerchantId        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}
