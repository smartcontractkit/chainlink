package fakes

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	customhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestDirectHTTPAction_RequestHeaders(t *testing.T) {
	t.Run("MultiHeaders are sent in request", func(t *testing.T) {
		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		t.Cleanup(srv.Close)

		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:    srv.URL,
			Method: "GET",
			MultiHeaders: map[string]*customhttp.HeaderValues{
				"Authorization": {Values: []string{"Bearer test-token"}},
			},
		}
		metadata := commonCap.RequestMetadata{}

		result, err := action.SendRequest(context.Background(), metadata, input)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Bearer test-token", receivedAuth, "Authorization header should be sent")
	})

	t.Run("Headers (deprecated) are sent in request when MultiHeaders empty", func(t *testing.T) {
		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		t.Cleanup(srv.Close)

		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:     srv.URL,
			Method:  "GET",
			Headers: map[string]string{"Authorization": "Basic legacy-auth"},
		}
		metadata := commonCap.RequestMetadata{}

		result, err := action.SendRequest(context.Background(), metadata, input)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Basic legacy-auth", receivedAuth, "Authorization header should be sent via deprecated Headers")
	})
}

// Fixed self-signed ECDSA P-256 material (valid until 2100) used by the mTLS
// tests. The client cert (CN=test-client) doubles as its own CA so the server
// can verify it; the server cert (CN=127.0.0.1, IP SAN 127.0.0.1) doubles as
// its own CA so the client can verify the test server.
const (
	testClientCertPEM = `-----BEGIN CERTIFICATE-----
MIIBpDCCAUugAwIBAgIQNUFEoGtJzPWVdReYl4bKNzAKBggqhkjOPQQDAjAWMRQw
EgYDVQQDEwt0ZXN0LWNsaWVudDAgFw0yMDAxMDEwMDAwMDBaGA8yMTAwMDEwMTAw
MDAwMFowFjEUMBIGA1UEAxMLdGVzdC1jbGllbnQwWTATBgcqhkjOPQIBBggqhkjO
PQMBBwNCAAQB6kvfLo+M0LrgPfSvopjSCedHqLsXipTzXLLC1xwd1MAqky2TXWzh
yqvdKr21hUGVDd8FCOmrSkrTHUjhR0v2o3kwdzAOBgNVHQ8BAf8EBAMCAqQwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8wHQYD
VR0OBBYEFAeZ4IPcjmPgQSrOgLnLjc3F0u6BMBYGA1UdEQQPMA2CC3Rlc3QtY2xp
ZW50MAoGCCqGSM49BAMCA0cAMEQCIGIyAh1wNWKKa0MvGOYLcpkUndR4TgLfeIRc
HJ8fFIy5AiByWsuX6CCKvHE5jyafDxUiaRlPKO9OzrENpDqppX73pQ==
-----END CERTIFICATE-----
`
	testClientKeyPEM = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPZ0mut8ctRN+QGoZd9j+8iuGsnuyfKG12Ap9vjPq67hoAoGCCqGSM49
AwEHoUQDQgAEAepL3y6PjNC64D30r6KY0gnnR6i7F4qU81yywtccHdTAKpMtk11s
4cqr3Sq9tYVBlQ3fBQjpq0pK0x1I4UdL9g==
-----END EC PRIVATE KEY-----
`
	testServerCertPEM = `-----BEGIN CERTIFICATE-----
MIIBpDCCAUugAwIBAgIQBQ40ltFw9DgQ1NdykQUJXDAKBggqhkjOPQQDAjAUMRIw
EAYDVQQDEwkxMjcuMC4wLjEwIBcNMjAwMTAxMDAwMDAwWhgPMjEwMDAxMDEwMDAw
MDBaMBQxEjAQBgNVBAMTCTEyNy4wLjAuMTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABLu51hFtDNleEjaRmPvKKZQQ2narQ0q/GG2SJ1pIjN5O8LuyERGoq5C2d3Tj
xdDSas0C09epfv8Qp/qrQAlzkamjfTB7MA4GA1UdDwEB/wQEAwICpDAdBgNVHSUE
FjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4E
FgQUWXiLGyaH8brXSPqiLQLJwxwPPEQwGgYDVR0RBBMwEYIJMTI3LjAuMC4xhwR/
AAABMAoGCCqGSM49BAMCA0cAMEQCICea83WAn85mar5k5qS9XtIcPCIF2zSICYoh
IorjG1LeAiAecMcvVmaRf2bjEO6CYFOtP5yLgaXP/5bZCE7j6fOprw==
-----END CERTIFICATE-----
`
	testServerKeyPEM = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDgbSkGePmXcnAhSXPxly5AyrCyb4uKHV3FaVD1Ps3kxoAoGCCqGSM49
AwEHoUQDQgAEu7nWEW0M2V4SNpGY+8oplBDadqtDSr8YbZInWkiM3k7wu7IREair
kLZ3dOPF0NJqzQLT16l+/xCn+qtACXORqQ==
-----END EC PRIVATE KEY-----
`
)

func TestDirectHTTPAction_Mtls(t *testing.T) {
	t.Parallel()

	t.Run("presents the supplied client certificate to the server", func(t *testing.T) {
		t.Parallel()
		clientCAPool := x509.NewCertPool()
		require.True(t, clientCAPool.AppendCertsFromPEM([]byte(testClientCertPEM)))

		serverPair, err := tls.X509KeyPair([]byte(testServerCertPEM), []byte(testServerKeyPEM))
		require.NoError(t, err)

		peerCNCh := make(chan string, 1)
		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NotNil(t, r.TLS, "request must arrive over TLS")                          //nolint:testifylint // require in handler goroutine
			require.Len(t, r.TLS.PeerCertificates, 1, "client must present exactly one cert") //nolint:testifylint // require in handler goroutine
			peerCNCh <- r.TLS.PeerCertificates[0].Subject.CommonName
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("mtls-ok"))
		}))
		srv.TLS = &tls.Config{
			Certificates: []tls.Certificate{serverPair},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    clientCAPool,
			MinVersion:   tls.VersionTLS12,
		}
		srv.StartTLS()
		t.Cleanup(srv.Close)

		input := &customhttp.Request{
			Url:    srv.URL,
			Method: "GET",
			Mtls: &customhttp.MtlsAuth{
				Certificate: []byte(testClientCertPEM),
				PrivateKey:  []byte(testClientKeyPEM),
			},
		}

		client, err := newHTTPClient(input, 5*time.Second)
		require.NoError(t, err)

		// Trust the server's self-signed cert so the handshake completes; the
		// production path uses the system pool, which can't trust a test server.
		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)
		require.NotNil(t, transport.TLSClientConfig)
		require.Len(t, transport.TLSClientConfig.Certificates, 1, "client certificate must be installed on the transport")
		serverPool := x509.NewCertPool()
		require.True(t, serverPool.AppendCertsFromPEM([]byte(testServerCertPEM)))
		transport.TLSClientConfig.RootCAs = serverPool

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "test-client", <-peerCNCh, "server must observe the supplied client certificate's CN")
	})

	t.Run("no mtls configured leaves the default transport", func(t *testing.T) {
		t.Parallel()
		input := &customhttp.Request{Url: "https://example.com", Method: "GET"}

		client, err := newHTTPClient(input, 5*time.Second)
		require.NoError(t, err)
		assert.Nil(t, client.Transport, "no mtls auth should not configure a custom transport")
	})

	t.Run("invalid mtls material fails the request like the gateway client", func(t *testing.T) {
		t.Parallel()
		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:    "https://example.com",
			Method: "GET",
			Mtls: &customhttp.MtlsAuth{
				Certificate: []byte("not a pem certificate"),
				PrivateKey:  []byte("not a pem key"),
			},
		}

		result, err := action.SendRequest(context.Background(), commonCap.RequestMetadata{}, input)
		require.Error(t, err)
		assert.NotEqual(t, caperrors.InvalidArgument, err.Code())
		require.NotNil(t, result)
		require.NotNil(t, result.Response)
		assert.Equal(t, uint32(0), result.Response.StatusCode)
		assert.Contains(t, err.Error(), "failed to parse mtls auth into key pair")
	})
}

func TestDirectHTTPAction_ResponseHeadersAndMultiHeaders(t *testing.T) {
	t.Run("response has both Headers and MultiHeaders populated", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Add("Set-Cookie", "sessionid=abc123; Path=/")
			w.Header().Add("Set-Cookie", "csrf=xyz789; Path=/")
			w.Header().Add("X-Custom", "single-value")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
		t.Cleanup(srv.Close)

		lggr := logger.Test(t)
		action := NewDirectHTTPAction(lggr)
		require.NoError(t, action.Start(context.Background()))
		t.Cleanup(func() { _ = action.Close() })

		input := &customhttp.Request{
			Url:    srv.URL,
			Method: "GET",
		}
		metadata := commonCap.RequestMetadata{}

		result, err := action.SendRequest(context.Background(), metadata, input)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Response)

		resp := result.Response
		assert.Equal(t, uint32(200), resp.StatusCode)

		// Headers (comma-joined, backwards compat)
		require.NotNil(t, resp.Headers)                                    //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers, "Content-Type")                   //nolint:staticcheck // testing deprecated field
		assert.Equal(t, "application/json", resp.Headers["Content-Type"])  //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers, "Set-Cookie")                     //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers["Set-Cookie"], "sessionid=abc123") //nolint:staticcheck // testing deprecated field
		assert.Contains(t, resp.Headers["Set-Cookie"], "csrf=xyz789")      //nolint:staticcheck // testing deprecated field
		assert.Equal(t, "single-value", resp.Headers["X-Custom"])          //nolint:staticcheck // testing deprecated field

		// MultiHeaders (per-value slices)
		require.NotNil(t, resp.MultiHeaders)
		assert.Contains(t, resp.MultiHeaders, "Content-Type")
		assert.Equal(t, []string{"application/json"}, resp.MultiHeaders["Content-Type"].GetValues())

		setCookie := resp.MultiHeaders["Set-Cookie"]
		require.NotNil(t, setCookie)
		vals := setCookie.GetValues()
		require.Len(t, vals, 2)
		assert.Contains(t, vals, "sessionid=abc123; Path=/")
		assert.Contains(t, vals, "csrf=xyz789; Path=/")

		assert.Contains(t, resp.MultiHeaders, "X-Custom")
		assert.Equal(t, []string{"single-value"}, resp.MultiHeaders["X-Custom"].GetValues())
	})
}
