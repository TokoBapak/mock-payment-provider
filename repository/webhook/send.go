package webhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *Client) Send(ctx context.Context, payload []byte) error {
	var retryCounter int = 0
	for retryCounter < 5 {
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.targetUrl, bytes.NewReader(payload))
		if err != nil {
			return fmt.Errorf("creating new request: %w", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return fmt.Errorf("executing http request: %w", err)
		}

		if response.StatusCode < 300 {
			return nil
		}

		if response.StatusCode >= 400 && response.StatusCode < 499 {
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("reading response body: %w", err)
			}

			return fmt.Errorf("request error: %s", string(responseBody))
		}

		retryCounter += 1
		time.Sleep(time.Second * 10 * time.Duration(retryCounter))
		continue
	}

	return fmt.Errorf("too many retries")
}
