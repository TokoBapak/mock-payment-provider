package webhook

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) Send(ctx context.Context, payload []byte) error {
	// Midtrans retry rules:
	//
	// for 2xx: No retries, it is considered success.
	// for 500: Retry only once.
	// for 503: Retry four times.
	// for 400/404: Retry two times.
	// for 301/302/303: No retries. We suggest to update the Notification endpoint in Settings instead of replying to these status code.
	// for 307/308: Follow the new URL with POST method and same notification body. Max redirect is five times.
	// for all other failures: Retry five times.
	//
	// Different retry intervals from 1st time to 5th time (2m, 10m, 30m, 1.5hour, 3.5hour).
	// Put a time shift for each retry based on the above interval. For example, for the first time,
	// retry might be two minutes after the job failed. The second retry might be ten minutes after
	// the first retry is failed and so on.

	var retryDuration = map[int]time.Duration{
		0: time.Minute * 2,
		1: time.Minute * 10,
		3: time.Minute * 30,
		4: time.Hour + (time.Minute * 30),
		5: (time.Hour * 3) + (time.Minute * 30),
	}
	var retryCounter int = 0
	var maximumRetry int = 5
	var initialRetrySet = false
	for retryCounter < maximumRetry {
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

		if response.StatusCode < 400 {
			return nil
		}

		if response.StatusCode == 400 || response.StatusCode == 404 {
			if !initialRetrySet {
				maximumRetry = 2
				initialRetrySet = true
			}
		}

		if response.StatusCode == 500 {
			if !initialRetrySet {
				maximumRetry = 1
				initialRetrySet = true
			}
		}

		if response.StatusCode == 503 {
			if !initialRetrySet {
				maximumRetry = 4
				initialRetrySet = true
			}
		}

		time.Sleep(retryDuration[retryCounter])
		retryCounter += 1
		continue
	}

	return fmt.Errorf("too many retries")
}
