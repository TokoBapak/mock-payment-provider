package webhook_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mockServerAddress string
var incomingRequests = &MockRequest{}

type MockRequest struct {
	requests    []requestResult
	lastRequest requestResult
}

type requestResult struct {
	Method string
	URL    string
	Body   []byte
	Header http.Header
}

func MockWebhookTargetServer(incomingRequests *MockRequest) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var requestBody []byte = nil
		if r.Body != nil {
			var err error = nil
			requestBody, err = io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
		}

		result := requestResult{
			URL:    r.URL.String(),
			Body:   requestBody,
			Header: r.Header,
			Method: r.Method,
		}
		incomingRequests.requests = append(incomingRequests.requests, result)
		incomingRequests.lastRequest = result
		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewServer(mux)
}

func TestMain(m *testing.M) {
	mockServer := MockWebhookTargetServer(incomingRequests)
	mockServerAddress = mockServer.URL
	exitCode := m.Run()

	mockServer.Close()
	os.Exit(exitCode)
}
