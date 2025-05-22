//go:debug netdns=go
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestHTTPClient_Send(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)
	// Define test cases
	tests := []struct {
		name             string
		setupServer      func() *httptest.Server
		configOption     func(*HTTPClientConfig)
		request          HTTPRequest
		giveMaxRespBytes uint32
		expectedError    error
		expectedResp     *HTTPResponse
	}{
		{
			name: "successful request",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					assert.NoError(t, err2)
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
			name: "transmits headers",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "bar", r.Header.Get("foo"))
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					assert.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method: "GET",
				URL:    "/",
				Headers: map[string]string{
					"foo": "bar",
				},
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
			name: "context canceled due to timeout passed in request",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(10 * time.Second)
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					assert.NoError(t, err2)
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
			name: "context canceled due to default timeout",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(10 * time.Second)
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					assert.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
			},
			expectedError: context.DeadlineExceeded,
		},
		{
			name: "success with long timeout passed in request",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(1 * time.Second)
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					assert.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedResp: &HTTPResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Length": "7"},
				Body:       []byte("success"),
			},
		},
		{
			name: "fails with long timeout capped by default",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(1 * time.Second)
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write([]byte("success"))
					assert.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 5 * time.Second,
			},
			configOption: func(hc *HTTPClientConfig) {
				hc.maxRequestDuration = 100 * time.Millisecond
			},
			expectedError: context.DeadlineExceeded,
		},
		{
			name: "server error",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					_, err2 := w.Write([]byte("error"))
					assert.NoError(t, err2)
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
			name: "response too long with non-default config",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write(make([]byte, 2048))
					assert.NoError(t, err2)
				}))
			},
			giveMaxRespBytes: 1024,
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
			name: "success with long response and default config",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, err2 := w.Write(make([]byte, 2048))
					assert.NoError(t, err2)
				}))
			},
			request: HTTPRequest{
				Method:  "GET",
				URL:     "/",
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 2 * time.Second,
			},
			expectedResp: &HTTPResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Length": "2048"},
				Body:       make([]byte, 2048),
			},
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

			config := &HTTPClientConfig{
				MaxResponseBytes: tt.giveMaxRespBytes,
				AllowedIPs:       []string{hostname},
				AllowedPorts:     []int{int(portInt)},
			}

			client, err := NewHTTPClient(*config, lggr)
			require.NoError(t, err)

			if tt.configOption != nil {
				hc, ok := client.(*httpClient)
				require.True(t, ok)
				tt.configOption(&hc.config)
			}

			tt.request.URL = server.URL + tt.request.URL

			resp, err := client.Send(context.Background(), tt.request)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedResp.StatusCode, resp.StatusCode)
			for k, v := range tt.expectedResp.Headers {
				value, ok := resp.Headers[k]
				require.True(t, ok)
				require.Equal(t, v, value)
			}
			require.Equal(t, tt.expectedResp.Body, resp.Body)
		})
	}
}

// IMPORTANT: The behaviour of Go's network stack is heavily dependent on the platform;
// this means that the errors returned can change depending on whether the tests are
// run on osx or on linux.
func TestHTTPClient_BlocksUnallowed(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)
	// Define test cases
	tests := []struct {
		name          string
		url           string
		expectedError string
		blockPort     bool
	}{
		{
			name:          "blocked port",
			url:           "http://177.0.0.1:8080",
			expectedError: "port: 8080 not found in allowlist",
			blockPort:     true,
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
			url:           "http://169.254.0.1",
			expectedError: "ip: 169.254.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback",
			url:           "http://127.0.0.1",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback without scheme",
			url:           "127.0.0.1",
			expectedError: "host:  is not valid",
		},
		{
			name:          "explicitly blocked IP - loopback",
			url:           "https://⑫7.0.0.1",
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
		{
			name:          "explicitly blocked IP - loopback shortened",
			url:           "https://127.1",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - loopback shortened",
			url:           "https://127.0.1",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - loopback hex encoded with separators",
			url:           `https://0x7F.0x00.0x00.0x01`,
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - loopback octal encoded",
			url:           `https://0177.0000.0000.0001`,
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - loopback binary encoded",
			url:           `https://01111111.00000000.00000000.00000001`,
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - loopback - dword no escape",
			url:           "https://2130706433",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - loopback - dword with overflow no escape",
			url:           "https://45080379393",
			expectedError: "no such host",
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
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - current network - hex",
			url:           "http://0x00.0x00.0x00.0x01",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - current network - binary",
			url:           "http://00000000.00000000.00000000.00000001",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - current network - shortened",
			url:           "http://1",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - current network - shortened",
			url:           "http://0.1",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - current network - shortened",
			url:           "http://0.0.1",
			expectedError: "no such host",
		},
		{
			name:          "explicitly blocked IP - dword",
			url:           "http://42949672961",
			expectedError: "no such host",
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
			testURL, err := url.Parse(tt.url)
			require.NoError(t, err)

			if testURL.Port() == "" {
				// Setup a test server so the request succeeds if we don't block it, then modify the URL to add the port to it.
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				defer server.Close()

				u, ierr := url.Parse(server.URL)
				require.NoError(t, ierr)

				testURL.Host = testURL.Hostname() + ":" + u.Port()
			}

			portInt, err := strconv.ParseInt(testURL.Port(), 10, 64)
			require.NoError(t, err)

			allowedPorts := []int{}
			if !tt.blockPort {
				allowedPorts = []int{int(portInt)}
			}

			config := HTTPClientConfig{
				MaxResponseBytes: 1024,
				DefaultTimeout:   5 * time.Second,
				AllowedPorts:     allowedPorts,
			}

			client, err := NewHTTPClient(config, lggr)
			require.NoError(t, err)

			require.NoError(t, err)
			_, err = client.Send(context.Background(), HTTPRequest{
				Method:  "GET",
				URL:     testURL.String(),
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 1 * time.Second,
			})
			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectedError)
		})
	}
}

func TestHTTPClient_AllowedIPsCIDR(t *testing.T) {
	t.Parallel()

	// Setup the test environment
	lggr := logger.Test(t)

	// Start a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	hostname, port := u.Hostname(), u.Port()
	t.Logf("hostname: %s, port: %s", hostname, port)
	portInt, err := strconv.ParseInt(port, 10, 32)
	require.NoError(t, err)

	// Define test cases
	tests := []struct {
		name          string
		allowedCIDRs  []string
		expectedError string
	}{
		{
			name:          "allowed CIDR block",
			allowedCIDRs:  []string{"127.0.0.1/32"},
			expectedError: "",
		},
		{
			name:          "blocked CIDR block",
			allowedCIDRs:  []string{"192.168.1.0/24"},
			expectedError: "ip: 127.0.0.1 not found in allowlist",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := HTTPClientConfig{
				MaxResponseBytes: 1024,
				DefaultTimeout:   5 * time.Second,
				AllowedPorts:     []int{int(portInt)},
				AllowedIPsCIDR:   tt.allowedCIDRs,
			}

			client, err := NewHTTPClient(config, lggr)
			require.NoError(t, err)

			_, err = client.Send(context.Background(), HTTPRequest{
				Method:  "GET",
				URL:     server.URL,
				Headers: map[string]string{},
				Body:    nil,
				Timeout: 1 * time.Second,
			})

			if tt.expectedError != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_ConfigApplyDefaults(t *testing.T) {
	t.Parallel()
	t.Run("successfully overrides defaults", func(t *testing.T) {
		config := HTTPClientConfig{
			MaxResponseBytes: 1024,
			DefaultTimeout:   5 * time.Second,
		}
		config.ApplyDefaults()
		require.Equal(t, uint32(1024), config.MaxResponseBytes)
		require.Equal(t, 5*time.Second, config.DefaultTimeout)
	})

	t.Run("successfully sets default values", func(t *testing.T) {
		config := HTTPClientConfig{}
		config.ApplyDefaults()
		require.Equal(t, defaultMaxResponseBytes, config.MaxResponseBytes) // 30MB
		require.Equal(t, defaultTimeout, config.DefaultTimeout)
		require.Equal(t, defaultAllowedPorts, config.AllowedPorts)
		require.Equal(t, defaultAllowedSchemes, config.AllowedSchemes)
	})
}
