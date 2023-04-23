package schema

import "mock-payment-provider/primitive"

type InternalMarkAsPaidRequest struct {
	OrderId       string                `json:"order_id"`
	PaymentMethod primitive.PaymentType `json:"payment_method"`
}
