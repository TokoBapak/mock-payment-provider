package business

import (
	"context"

	"mock-payment-provider/primitive"
)

// Payment interface handles anything to do with the completion of a payment,
// from the customer's standpoint.
type Payment interface {
	// GetDetail GetDetails will acquire a payment detail from an ID coming from e-money or
	// virtual account payment.
	GetDetail(ctx context.Context, id string) (PaymentDetailsResponse, error)
	// MarkAsPaid will mark an order ID as paid. This one function must only be called
	// from the presentation that handles payment confirmation from the customer's
	// standpoint.
	MarkAsPaid(ctx context.Context, orderId string, paymentMethod primitive.PaymentType) error
}

type PaymentDetailsResponse struct {
	OrderId              string
	ChargedAmount        int64
	Status               primitive.TransactionStatus
	PaymentMethod        primitive.PaymentType
	VirtualAccountNumber string
	EMoneyID             string
}
