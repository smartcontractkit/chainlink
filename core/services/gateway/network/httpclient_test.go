package network

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestHTTPClient_Send(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)
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
		{
			name: "redirects are blocked",
			setupServer: func() *httptest.Server {
				count := 0
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					count++
					if count <= 1 {
						http.Redirect(w, r, "/", http.StatusMovedPermanently)
					} else {
						w.WriteHeader(http.StatusOK)
					}
					count++
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedError: errors.New("redirects are not allowed"),
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

			config := HTTPClientConfig{
				MaxResponseBytes: 1024,
				DefaultTimeout:   5 * time.Second,
				allowedIPs:       []string{hostname},
				AllowedPorts:     []int{int(portInt)},
			}

			client, err := NewHTTPClient(config, lggr)
			require.NoError(t, err)

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
		BlockedIPs:       []string{"177.0.0.1"},
	}

	client, err := NewHTTPClient(config, lggr)
	require.NoError(t, err)

	// Define test cases
	tests := []struct {
		name          string
		url           string
		expectedError string
	}{
		{
			name:          "blocked port",
			url:           "http://127.0.0.1:8080",
			expectedError: "port: 8080 not found in allowlist",
		},
		{
			name:          "blocked scheme",
			url:           "file://127.0.0.1",
			expectedError: "scheme: file not found in allowlist",
		},
		{
			name:          "explicitly blocked IP",
			url:           "http://169.254.0.1",
			expectedError: "ip: 169.254.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - internal network",
			url:           "http://169.254.0.1/endpoint",
			expectedError: "ip: 169.254.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback",
			url:           "http://127.0.0.1/endpoint",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback without scheme",
			url:           "127.0.0.1",
			expectedError: "host:  is not valid",
		},
		{
			name:          "explicitly blocked IP - loopback",
			url:           "https://⑫7.0.0.1/",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback shortened",
			url:           "https://127.1",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback shortened",
			url:           "https://127.0.1",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback hex no separators",
			url:           "https://0x17F000001/",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback hex encoded with separators",
			url:           `https://0x7F.0x00.0x00.0x01`,
			expectedError: "127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback octal encoded",
			url:           `https://\\0177.0000.0000.0001`,
			expectedError: "invalid character",
		},
		{
			name: "explicitly blocked IP - loopback octal encoded",
			url:  `https://0177.0000.0000.0001`,
			// 0177.0000.0000.0001 is 127.0.0.1 octal encoded
			// however Go interprets this as `177.0.0.1`
			// in the test setup we add that URL to the list of blocked IPs
			// if the test confirms that 177.0.0.1 is blocked, then that
			// means we cannot access 127.0.0.1 via its octal encoding.
			expectedError: "177.0.0.1 not found in allowlist",
		},
		{
			name: "explicitly blocked IP - loopback binary encoded",
			url:  `https://01111111.00000000.00000000.00000001`,
			// fails with dial error
			expectedError: "no such host",
		},
		{
			name: "explicitly blocked IP - loopback - dword",
			url:  "https://\\2130706433/",
			// this fails with a parsing error
			expectedError: "invalid character",
		},
		{
			name: "explicitly blocked IP - loopback - dword with overflow",
			url:  "https://\\45080379393/",
			// this fails with a parsing error
			expectedError: "invalid character",
		},
		{
			name:          "explicitly blocked IP - loopback - dword no escape",
			url:           "https://2130706433/",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback - dword with overflow no escape",
			url:           "https://45080379393/",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback - ipv6",
			url:           `https://[::1]`,
			expectedError: "ipv6 blocked",
		},
		{
			name:          "explicitly blocked IP - loopback ipv6 mapped ipv4",
			url:           `https://[::FFF:7F00:0001]`,
			expectedError: "ipv6 blocked",
		},
		{
			name:          "explicitly blocked IP - loopback ipv6 mapped ipv4",
			url:           `https://[::FFFF:127.0.0.1]`,
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback long-form",
			url:           `https://[0000:0000:0000:0000:0000:0000:0000:0001]`,
			expectedError: "ipv6 blocked",
		},
		{
			name:          "explicitly blocked IP - current network",
			url:           "http://0.0.0.0/endpoint",
			expectedError: "ip: 0.0.0.0 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - octal",
			url:           "http://0000.0000.0000.0001",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - hex",
			url:           "http://0x00.0x00.0x00.0x01",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - hex no separators",
			url:           "http://0x00000001",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - binary",
			url:           "http://00000000.00000000.00000000.00000001",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - shortened",
			url:           "http://1",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - shortened",
			url:           "http://0.1",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - current network - shortened",
			url:           "http://0.0.1",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - dword",
			url:           "http://42949672961",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - ipv6 mapped",
			url:           "http://[::FFFF:0000:0001]",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - ipv6 mapped",
			url:           "http://[::FFFF:0.0.0.1]",
			expectedError: "ip: 0.0.0.1 not found in allowlist",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Send(context.Background(), HTTPRequest{
				Method:  "GET",
				URL:     tt.url,
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 10 * time.Millisecond,
			})
			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectedError)
		})
	}
}
