package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/doyensec/safeurl"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestHTTPClient_Send(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)
	config := HTTPClientConfig{
		MaxResponseBytes: 1024,
		DefaultTimeout:   5 * time.Second,
	}

	// Define test cases
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		request       HTTPRequest
		expectedError error
		expectedResp  *HTTPResponse
	}{
		{
			name: "successful request",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					require.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: nil,
			expectedResp: &HTTPResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Length": "7"},
				Body:       []byte("success"),
			},
		},
		{
			name: "request timeout",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(10 * time.Second)
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					require.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 1 * time.Second,
			},
			expectedError: context.DeadlineExceeded,
			expectedResp:  nil,
		},
		{
			name: "server error",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					_, err2 := w.Write([]byte("error"))
					require.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: nil,
			expectedResp: &HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Length": "5"},
				Body:       []byte("error"),
			},
		},
		{
			name: "response too long",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write(make([]byte, 2048))
					require.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: &http.MaxBytesError{},
			expectedResp:  nil,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			u, err := url.Parse(server.URL)
			require.NoError(t, err)

			hostname, port := u.Hostname(), u.Port()
			portInt, err := strconv.ParseInt(port, 10, 32)
			require.NoError(t, err)

			safeConfig := safeurl.
				GetConfigBuilder().
				SetTimeout(config.DefaultTimeout).
				SetAllowedIPs(hostname).
				SetAllowedPorts(int(portInt)).
				Build()

			client := &httpClient{
				config: config,
				client: safeurl.Client(safeConfig),
				lggr:   lggr,
			}

			tt.request.URL = server.URL + tt.request.URL

			resp, err := client.Send(context.Background(), tt.request)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResp.StatusCode, resp.StatusCode)
				for k, v := range tt.expectedResp.Headers {
					value, ok := resp.Headers[k]
					require.True(t, ok)
					require.Equal(t, v, value)
				}
				require.Equal(t, tt.expectedResp.Body, resp.Body)
			}
		})
	}
}

func TestHTTPClient_BlocksUnallowed(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)
	config := HTTPClientConfig{
		MaxResponseBytes: 1024,
		DefaultTimeout:   5 * time.Second,
	}

	client, err := NewHTTPClient(config, lggr)
	require.NoError(t, err)

	// Define test cases
	tests := []struct {
		name          string
		request       HTTPRequest
		expectedError string
	}{
		{
			name: "blocked port",
			request: HTTPRequest{
				Method:  "GET",
				URL:     "http://127.0.0.1:8080",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: "port: 8080 not found in allowlist",
		},
		{
			name: "blocked scheme",
			request: HTTPRequest{
				Method:  "GET",
				URL:     "file://127.0.0.1",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: "scheme: file not found in allowlist",
		},
		{
			name: "explicitly blocked IP",
			request: HTTPRequest{
				Method:  "GET",
				URL:     "http://169.254.0.1",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: "ip: 169.254.0.1 not found in allowlist",
		},
		{
			name: "explicitly blocked IP - internal network",
			request: HTTPRequest{
				Method:  "GET",
				URL:     "http://169.254.0.1/endpoint",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: "ip: 169.254.0.1 not found in allowlist",
		},
		{
			name: "explicitly blocked IP - localhost",
			request: HTTPRequest{
				Method:  "GET",
				URL:     "http://127.0.0.1/endpoint",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name: "explicitly blocked IP - current network",
			request: HTTPRequest{
				Method:  "GET",
				URL:     "http://0.0.0.0/endpoint",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: "ip: 0.0.0.0 not found in allowlist",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Send(context.Background(), tt.request)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectedError)
		})
	}
}
