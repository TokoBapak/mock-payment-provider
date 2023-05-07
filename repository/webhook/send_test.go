package webhook_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"mock-payment-provider/repository/webhook"
)

func TestClient_Send(t *testing.T) {
	webhookClient, err := webhook.NewWebhookClient(mockServerAddress)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	t.Run("Successful", func(t *testing.T) {
		payload := []byte("Hello world")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := webhookClient.Send(ctx, payload)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if !bytes.Equal(incomingRequests.lastRequest.Body, payload) {
			t.Errorf("expecting lastRequest.Body to equal payload, instead got %s vs %s", string(incomingRequests.lastRequest.Body), string(payload))
		}

		if incomingRequests.lastRequest.Method != http.MethodPost {
			t.Errorf("expecting lastRequest.Method to equal POST, instead got %s", incomingRequests.lastRequest.Method)
		}
	})
}
