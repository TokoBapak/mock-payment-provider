package business

import (
	"context"
	"time"

	"mock-payment-provider/primitive"
)

type Transaction interface {
	Charge(ctx context.Context, request ChargeRequest) (ChargeResponse, error)
	Cancel(ctx context.Context, orderId string) (CancelResponse, error)
	GetStatus(ctx context.Context, orderId string) (GetStatusResponse, error)
}

type ProductItem struct {
	ID    string
	Price int64
	// Cannot be lower than 0
	Quantity int64
	Name     string
	Category string
}

type Address struct {
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Address     string
	PostalCode  string
	CountryCode string
}

type CustomerInformation struct {
	FirstName      string
	LastName       string
	Email          string
	PhoneNumber    string
	BillingAddress Address
}

type SellerInformation struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Address     string
}

type VirtualAccountAction struct {
	Bank                 string
	VirtualAccountNumber string
}

type EMoneyActionType uint8

const (
	EMoneyActionTypeUnspecified EMoneyActionType = iota
	EMoneyActionTypePay
	EMoneyActionTypeStatus
	EMoneyActionTypeCancel
)

func (e EMoneyActionType) String() string {
	switch e {
	case EMoneyActionTypePay:
		return "CODE"
	case EMoneyActionTypeStatus:
		return "STATUS"
	case EMoneyActionTypeCancel:
		return "CANCEL"
	case EMoneyActionTypeUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}

type EMoneyAction struct {
	EMoneyActionType EMoneyActionType
	Method           string
	URL              string
}

type ChargeRequest struct {
	PaymentType         primitive.PaymentType
	OrderId             string
	TransactionAmount   int64
	TransactionCurrency primitive.Currency
	Customer            CustomerInformation
	Seller              SellerInformation
	ProductItems        []ProductItem
}

type ChargeResponse struct {
	OrderId              string
	TransactionAmount    int64
	PaymentType          primitive.PaymentType
	TransactionStatus    primitive.TransactionStatus
	TransactionTime      time.Time
	EMoneyAction         []EMoneyAction
	VirtualAccountAction VirtualAccountAction
}

type CancelResponse struct {
	OrderId           string
	TransactionAmount int64
	PaymentType       primitive.PaymentType
	TransactionStatus primitive.TransactionStatus
	TransactionTime   time.Time
}

type GetStatusResponse struct {
	OrderId           string
	TransactionStatus primitive.TransactionStatus
}
