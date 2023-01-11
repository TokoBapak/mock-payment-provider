package repository

import "context"

type WebhookClient interface {
	Send(ctx context.Context, payload []byte) error
}
