package network_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

func TestHTTPClient_Send(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)
	config := network.HTTPClientConfig{
		MaxResponseBytes: 1024,
		DefaultTimeout:   5 * time.Second,
	}
	client, err := network.NewHTTPClient(config, lggr)
	require.NoError(t, err)

	// Define test cases
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		request       network.HTTPRequest
		expectedError error
		expectedResp  *network.HTTPResponse
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
			request: network.HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: nil,
			expectedResp: &network.HTTPResponse{
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
			request: network.HTTPRequest{
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
			request: network.HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: nil,
			expectedResp: &network.HTTPResponse{
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
			request: network.HTTPRequest{
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
