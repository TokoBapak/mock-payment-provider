package payment_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"mock-payment-provider/business"
	"mock-payment-provider/presentation/schema"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
	"mock-payment-provider/repository/signature"
)

func (d *Dependency) MarkAsPaid(ctx context.Context, orderId string, paymentMethod primitive.PaymentType) error {
	// Get transaction from order id
	transaction, err := d.transactionRepository.GetByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return business.ErrTransactionNotFound
		}

		return fmt.Errorf("acquiring transaction: %w", err)
	}

	// Check whether transaction is already expired
	if transaction.Expired() {
		return business.ErrCannotModifyStatus
	}

	// Check previous transaction status. We can only change it to settled only
	// if the previous status is pending.
	if transaction.TransactionStatus != primitive.TransactionStatusPending {
		return business.ErrCannotModifyStatus
	}

	// Mark as settled
	err = d.transactionRepository.UpdateStatus(ctx, orderId, primitive.TransactionStatusSettled)
	if err != nil {
		return fmt.Errorf("updating transaction status: %w", err)
	}

	var virtualAccountNumber = ""

	if paymentMethod == primitive.PaymentTypeUnspecified {
		paymentMethod = transaction.PaymentType
	}

	switch paymentMethod {
	case primitive.PaymentTypeVirtualAccountBCA:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBNI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountBRI:
		fallthrough
	case primitive.PaymentTypeVirtualAccountPermata:
		// Acquire virtual account number
		virtualAccountEntry, err := d.virtualAccountRepository.GetByOrderId(ctx, orderId)
		if err != nil {
			return fmt.Errorf("acquiring virtual account entry from order id: %w", err)
		}

		virtualAccountNumber = virtualAccountEntry.VirtualAccountNumber

		err = d.virtualAccountRepository.DeductCharge(ctx, virtualAccountNumber)
		if err != nil {
			return fmt.Errorf("deducting virtual account charge: %w", err)
		}
	case primitive.PaymentTypeEMoneyQRIS:
		fallthrough
	case primitive.PaymentTypeEMoneyGopay:
		fallthrough
	case primitive.PaymentTypeEMoneyShopeePay:
		err := d.eMoneyRepository.DeductCharge(ctx, orderId)
		if err != nil {
			return fmt.Errorf("deducting emoney charge: %w", err)
		}
	default:
		return fmt.Errorf("invalid payment type")
	}

	go func() {
		payload, err := d.buildSettlementMessage(settlementMessageParameters{
			PaymentType:          paymentMethod,
			OrderId:              orderId,
			TransactionTime:      transaction.TransactionTime,
			GrossAmount:          transaction.TransactionAmount,
			VirtualAccountNumber: virtualAccountNumber,
		})
		if err != nil {
			// TODO: proper error logging
			log.Printf("Encountered an error during marshaling json: %s", err.Error())
			return
		}

		ctx := context.Background()

		err = d.webhookClient.Send(ctx, payload)
		if err != nil {
			// TODO: proper error logging
			log.Printf("Encountered an error during sending webhook: %s", err.Error())
		}
	}()

	return nil
}

type settlementMessageParameters struct {
	PaymentType          primitive.PaymentType
	OrderId              string
	TransactionTime      time.Time
	GrossAmount          int64
	VirtualAccountNumber string
}

func (d *Dependency) buildSettlementMessage(parameters settlementMessageParameters) ([]byte, error) {
	signatureKey := signature.Generate(parameters.OrderId, 200, parameters.GrossAmount, d.serverKey)

	switch parameters.PaymentType {
	case primitive.PaymentTypeVirtualAccountBCA:
		return json.Marshal(schema.BCAVirtualAccountChargeSettlementResponse{
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     "bca",
					VaNumber: parameters.VirtualAccountNumber,
				},
			},
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.ToPaymentMethod(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusSettled.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
		})
	case primitive.PaymentTypeVirtualAccountBRI:
		return json.Marshal(schema.BRIVirtualAccountChargeSettlementResponse{
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     "bri",
					VaNumber: parameters.VirtualAccountNumber,
				},
			},
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.ToPaymentMethod(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusSettled.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
		})
	case primitive.PaymentTypeVirtualAccountBNI:
		return json.Marshal(schema.BNIVirtualAccountChargeSettlementResponse{
			VaNumbers: []struct {
				Bank     string `json:"bank"`
				VaNumber string `json:"va_number"`
			}{
				{
					Bank:     "bni",
					VaNumber: parameters.VirtualAccountNumber,
				},
			},
			PaymentAmounts: []struct {
				PaidAt string `json:"paid_at"`
				Amount string `json:"amount"`
			}{
				{
					PaidAt: time.Now().Format(time.DateTime),
					Amount: strconv.FormatInt(parameters.GrossAmount, 10),
				},
			},
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.ToPaymentMethod(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusSettled.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
		})
	case primitive.PaymentTypeVirtualAccountPermata:
		return json.Marshal(schema.PermataVirtualAccountChargeSettlementResponse{
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			OrderId:           parameters.OrderId,
			PaymentType:       parameters.PaymentType.ToPaymentMethod(),
			SignatureKey:      signatureKey,
			StatusCode:        "200",
			TransactionId:     parameters.OrderId,
			TransactionStatus: primitive.TransactionStatusSettled.String(),
			FraudStatus:       "accept",
			StatusMessage:     "midtrans payment notification",
			PermataVaNumber:   parameters.VirtualAccountNumber,
		})
	case primitive.PaymentTypeEMoneyQRIS:
		return json.Marshal(schema.QRISChargeSettlementResponse{
			TransactionType:          "",
			TransactionTime:          parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus:        primitive.TransactionStatusSettled.String(),
			TransactionId:            parameters.OrderId,
			StatusMessage:            "midtrans payment notification",
			StatusCode:               "200",
			SignatureKey:             signatureKey,
			SettlementTime:           time.Now().Format(time.DateTime),
			PaymentType:              parameters.PaymentType.ToPaymentMethod(),
			OrderId:                  parameters.OrderId,
			MerchantId:               "",
			Issuer:                   "",
			GrossAmount:              strconv.FormatInt(parameters.GrossAmount, 10),
			FraudStatus:              "accept",
			Currency:                 "IDR",
			Acquirer:                 "nobu",
			ShopeepayReferenceNumber: "",
			ReferenceId:              "",
		})
	case primitive.PaymentTypeEMoneyGopay:
		return json.Marshal(schema.GopayChargeSettlementResponse{
			StatusCode:        "200",
			StatusMessage:     "midtrans payment notification",
			TransactionId:     parameters.OrderId,
			OrderId:           parameters.OrderId,
			GrossAmount:       strconv.FormatInt(parameters.GrossAmount, 10),
			PaymentType:       parameters.PaymentType.ToPaymentMethod(),
			TransactionTime:   parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus: primitive.TransactionStatusSettled.String(),
			SignatureKey:      signatureKey,
		})
	case primitive.PaymentTypeEMoneyShopeePay:
		return json.Marshal(schema.ShopeePayChargeSettlementResponse{
			TransactionTime:          parameters.TransactionTime.Format(time.DateTime),
			TransactionStatus:        primitive.TransactionStatusSettled.String(),
			TransactionId:            parameters.OrderId,
			StatusMessage:            "midtrans payment notification",
			StatusCode:               "200",
			SignatureKey:             signatureKey,
			SettlementTime:           time.Now().Format(time.DateTime),
			PaymentType:              parameters.PaymentType.ToPaymentMethod(),
			OrderId:                  parameters.OrderId,
			MerchantId:               "",
			GrossAmount:              strconv.FormatInt(parameters.GrossAmount, 10),
			FraudStatus:              "accept",
			Currency:                 "IDR",
			ShopeepayReferenceNumber: "",
			ReferenceId:              "",
		})
	default:
		return nil, fmt.Errorf("invalid payment type")
	}
}
