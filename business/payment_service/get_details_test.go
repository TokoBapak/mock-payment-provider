package payment_service_test

import (
	"context"
	"fmt"
	"mock-payment-provider/business/payment_service"
	"mock-payment-provider/primitive"
	"mock-payment-provider/repository"
	"testing"
	"time"
)

type MockTransactionRepository struct{}

func (m *MockTransactionRepository) Migrate(ctx context.Context) error {
	return nil
}
func (m *MockTransactionRepository) Create(ctx context.Context, params repository.CreateTransactionParam) error {
	return nil
}
func (m *MockTransactionRepository) UpdateStatus(ctx context.Context, orderId string, status primitive.TransactionStatus) error {
	return nil
}
func (m *MockTransactionRepository) GetByOrderId(ctx context.Context, orderId string) (primitive.Transaction, error) {
	if orderId != "order-id" {
		return primitive.Transaction{}, repository.ErrNotFound
	}
	return primitive.Transaction{
		OrderId:           "order-id",
		TransactionAmount: 50000,
		PaymentType:       primitive.PaymentTypeEMoneyGopay,
		TransactionStatus: primitive.TransactionStatusPending,
		TransactionTime:   time.Now(),
		ExpiresAt:         time.Now().Add(5 * time.Minute),
	}, nil
}

type MockWebhookClient struct{}

func (m *MockWebhookClient) Send(ctx context.Context, payload []byte) error {
	return nil
}

type MockEMoneyRepository struct{}

func (m *MockEMoneyRepository) Migrate(ctx context.Context) error {
	return nil
}
func (m *MockEMoneyRepository) CreateCharge(ctx context.Context, orderId string, amount int64, expiresAt time.Time) (string, error) {
	return "emoney-id", nil
}
func (m *MockEMoneyRepository) GetByID(ctx context.Context, id string) (repository.Entry, error) {
	if id != "emoney-id" {
		return repository.Entry{}, repository.ErrNotFound
	}
	return repository.Entry{
		VirtualAccountNumber: "1234567890",
		EMoneyID:             "emoney-id",
		OrderId:              "order-id",
		ChargedAmount:        50000,
		ExpiresAt:            time.Now().Add(5 * time.Minute),
	}, nil
}
func (m *MockEMoneyRepository) GetByOrderId(ctx context.Context, orderId string) (repository.Entry, error) {
	if orderId != "order-id" {
		return repository.Entry{}, repository.ErrNotFound
	}
	return repository.Entry{
		VirtualAccountNumber: "1234567890",
		EMoneyID:             "emoney-id",
		OrderId:              "order-id",
		ChargedAmount:        50000,
		ExpiresAt:            time.Now().Add(5 * time.Minute),
	}, nil
}
func (m *MockEMoneyRepository) CancelCharge(ctx context.Context, orderId string) error {
	return nil
}
func (m *MockEMoneyRepository) DeductCharge(ctx context.Context, orderId string) error {
	return nil
}

type MockVirtualAccountRepository struct{}

func (m *MockVirtualAccountRepository) Migrate(ctx context.Context) error {
	return nil
}
func (m *MockVirtualAccountRepository) CreateOrGetVirtualAccountNumber(ctx context.Context, customerUniqueField string) (string, error) {
	return "1234567890", nil
}
func (m *MockVirtualAccountRepository) CreateCharge(ctx context.Context, virtualAccountNumber string, orderId string, amount int64, expiresAt time.Time) (string, error) {
	return "1234567890", nil
}
func (m *MockVirtualAccountRepository) GetByVirtualAccountNumber(ctx context.Context, virtualAccountNumber string) (repository.Entry, error) {
	if virtualAccountNumber == "" {
		return repository.Entry{}, fmt.Errorf("virtualAccountNumber is empty")
	}
	if virtualAccountNumber == "empty-order-id" {
		return repository.Entry{
			VirtualAccountNumber: "1234567890",
			EMoneyID:             "emoney-id",
			OrderId:              "",
			ChargedAmount:        50000,
			ExpiresAt:            time.Now().Add(5 * time.Minute),
		}, nil
	}
	if virtualAccountNumber != "order-id" {
		return repository.Entry{}, repository.ErrNotFound
	}
	return repository.Entry{
		VirtualAccountNumber: "1234567890",
		EMoneyID:             "emoney-id",
		OrderId:              "order-id",
		ChargedAmount:        50000,
		ExpiresAt:            time.Now().Add(5 * time.Minute),
	}, nil
}
func (m *MockVirtualAccountRepository) GetByOrderId(ctx context.Context, orderId string) (repository.Entry, error) {
	if orderId == "" {
		return repository.Entry{}, repository.ErrNotFound
	}
	if orderId != "order-id" {
		return repository.Entry{}, repository.ErrNotFound
	}
	return repository.Entry{
		VirtualAccountNumber: "1234567890",
		EMoneyID:             "emoney-id",
		OrderId:              "order-id",
		ChargedAmount:        50000,
		ExpiresAt:            time.Now().Add(5 * time.Minute),
	}, nil
}
func (m *MockVirtualAccountRepository) GetChargedAmount(ctx context.Context, virtualAccountNumber string) (int64, error) {
	return 50000, nil
}
func (m *MockVirtualAccountRepository) DeductCharge(ctx context.Context, virtualAccountNumber string) error {
	return nil
}

func TestBusinessGetDetails(t *testing.T) {
	// Create a mock of the repository
	mockTransactionRepository := &MockTransactionRepository{}
	mockWebhookClient := &MockWebhookClient{}
	mockEMoneyRepository := &MockEMoneyRepository{}
	mockVirtualAccountRepository := &MockVirtualAccountRepository{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conf := payment_service.Config{
		ServerKey:                "server-key",
		TransactionRepository:    mockTransactionRepository,
		WebhookClient:            mockWebhookClient,
		EMoneyRepository:         mockEMoneyRepository,
		VirtualAccountRepository: mockVirtualAccountRepository,
	}

	deps, err := payment_service.NewPaymentService(conf)
	if err != nil {
		t.Fatalf("creating payment service: %s", err.Error())
	}

	t.Run("GetDetails should return the correct details", func(t *testing.T) {
		userId, err := deps.GetDetail(ctx, "order-id")
		if err != nil {
			t.Fatalf("error: %s", err.Error())
		}
		t.Logf("result: %v", userId)
	})

	t.Run("GetDetails should return error if the order id is not found", func(t *testing.T) {
		res, err := deps.GetDetail(ctx, "not-exist")
		if err == nil {
			t.Fatalf("error: %s", err.Error())
		}
		t.Logf("result: %v", res)
	})

	t.Run("GetDetails should return error if the order id is empty", func(t *testing.T) {
		res, err := deps.GetDetail(ctx, "")
		if err == nil {
			t.Fatalf("error: %s", err.Error())
		}
		t.Logf("result: %v", res)
	})
	t.Run("GetDetails should return error if return orderId of the virtual account is empty", func(t *testing.T) {
		res, err := deps.GetDetail(ctx, "empty-order-id")
		if err == nil {
			t.Fatalf("error: %s", err.Error())
		}
		t.Logf("result: %v", res)
	})

}
