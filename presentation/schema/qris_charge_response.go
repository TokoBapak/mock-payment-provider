package schema

type QRISChargeSuccessResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	MerchantId        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	Currency          string `json:"currency"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	Acquirer          string `json:"acquirer"`
	Actions           []struct {
		Name   string `json:"name"`
		Method string `json:"method"`
		Url    string `json:"url"`
	} `json:"actions"`
}

type QRISChargePendingResponse struct {
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
	Acquirer          string `json:"acquirer"`
}

type QRISChargeSettlementResponse struct {
	TransactionType          string `json:"transaction_type"`
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
	Issuer                   string `json:"issuer"`
	GrossAmount              string `json:"gross_amount"`
	FraudStatus              string `json:"fraud_status"`
	Currency                 string `json:"currency"`
	Acquirer                 string `json:"acquirer"`
	ShopeepayReferenceNumber string `json:"shopeepay_reference_number"`
	ReferenceId              string `json:"reference_id"`
}

type QRISChargeExpiredResponse struct {
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
	Acquirer          string `json:"acquirer"`
}
