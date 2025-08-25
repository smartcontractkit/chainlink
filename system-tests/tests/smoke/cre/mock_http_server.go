package cre

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/stretchr/testify/require"
)

type MockServerRecorder struct {
	mu       sync.Mutex
	requests []RecordedRequest
}

type RecordedRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (r *MockServerRecorder) RecordRequest(method, url string, headers http.Header, body []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	headerMap := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			headerMap[key] = values[0]
		}
	}

	r.requests = append(r.requests, RecordedRequest{
		Method:  method,
		URL:     url,
		Headers: headerMap,
		Body:    string(body),
	})
}

func (r *MockServerRecorder) GetRequests() []RecordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	requests := make([]RecordedRequest, len(r.requests))
	copy(requests, r.requests)
	return requests
}

func (r *MockServerRecorder) GetRequestCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.requests)
}

func startMockHTTPServerOnPort(t *testing.T, port int) (*httptest.Server, *MockServerRecorder) {
	recorder := &MockServerRecorder{}
	mux := http.NewServeMux()

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		recorder.RecordRequest(r.Method, r.URL.String(), r.Header, body)

		framework.L.Info().Msgf("Mock server received order request: %s", string(body))

		response := map[string]interface{}{
			"orderId": "test-order-" + uuid.New().String()[0:8],
			"status":  "success",
			"message": "Order processed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	testServer := httptest.NewUnstartedServer(mux)
	testServer.Listener.Close()

	var err error
	testServer.Listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	require.NoError(t, err, "failed to listen on port %d", port)

	testServer.Start()

	framework.L.Info().Msgf("Mock HTTP server started on port %d at: %s", port, testServer.URL)
	return testServer, recorder
}
