package network

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/doyensec/safeurl"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
)

// testMtlsLimiter returns a generously-sized concurrency limiter for tests that
// only need mTLS construction to succeed (the limit itself is not exercised).
func testMtlsLimiter() limits.ResourcePoolLimiter[int] {
	return limits.GlobalResourcePoolLimiter[int](50)
}

// doerFunc adapts a function to the httpDoer interface.
type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

// underlyingSafeURLClient unwraps the *httpClient's doer down to the safeurl
// client, transparently stepping through the concurrencyLimitedClient decorator
// that the mTLS path installs.
func underlyingSafeURLClient(t *testing.T, hc *httpClient) *safeurl.WrappedClient {
	t.Helper()
	switch c := hc.client.(type) {
	case *safeurl.WrappedClient:
		return c
	case *concurrencyLimitedClient:
		wc, ok := c.client.(*safeurl.WrappedClient)
		require.True(t, ok, "expected *safeurl.WrappedClient under the limiter, got %T", c.client)
		return wc
	default:
		t.Fatalf("unexpected client type %T", hc.client)
		return nil
	}
}

// pemKeyPair generates a fresh ECDSA P-256 key and a self-signed certificate for
// the given common name. Returns PEM-encoded certificate and private key bytes
// in the same form `tls.X509KeyPair` accepts.
func pemKeyPair(t *testing.T, commonName string, ipSANs ...net.IP) (certPEM, keyPEM []byte, parsed *x509.Certificate) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           ipSANs,
		DNSNames:              []string{commonName},
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	require.NoError(t, err)

	parsed, err = x509.ParseCertificate(der)
	require.NoError(t, err)

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	keyDER, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return certPEM, keyPEM, parsed
}

func TestNewHTTPClientWithOptions_MtlsRejectsInvalidPEM(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	_, err := NewHTTPClient(HTTPClientConfig{
		Mtls: &gateway.MtlsAuth{
			Certificate: []byte("not a pem certificate"),
			PrivateKey:  []byte("not a pem key"),
		},
	}, lggr)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse MtlsAuth into KeyPair")
}

func TestNewHTTPClientWithOptions_MtlsRejectsMismatchedKeyPair(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	certPEM, _, _ := pemKeyPair(t, "client-a")
	_, keyPEMOther, _ := pemKeyPair(t, "client-b")

	_, err := NewHTTPClient(HTTPClientConfig{
		Mtls: &gateway.MtlsAuth{
			Certificate: certPEM,
			PrivateKey:  keyPEMOther,
		},
	}, lggr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse MtlsAuth into KeyPair")
}

func TestNewHTTPClientWithOptions_MtlsValidKeyPair(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	certPEM, keyPEM, _ := pemKeyPair(t, "client")

	client, err := NewHTTPClient(HTTPClientConfig{
		Mtls:               &gateway.MtlsAuth{Certificate: certPEM, PrivateKey: keyPEM},
		ConcurrencyLimiter: testMtlsLimiter(),
	}, lggr)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewHTTPClientFactory_MtlsFlow(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	factory := NewHTTPClientFactory(HTTPClientConfig{}, lggr)

	t.Run("nil config returns a working non-mtls client", func(t *testing.T) {
		client, err := factory(HTTPClientConfig{})
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("valid mtls config returns a working client", func(t *testing.T) {
		certPEM, keyPEM, _ := pemKeyPair(t, "client")
		client, err := factory(HTTPClientConfig{
			Mtls:               &gateway.MtlsAuth{Certificate: certPEM, PrivateKey: keyPEM},
			ConcurrencyLimiter: testMtlsLimiter(),
		})
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("invalid mtls config surface an error", func(t *testing.T) {
		client, err := factory(HTTPClientConfig{
			Mtls: &gateway.MtlsAuth{Certificate: []byte("garbage"), PrivateKey: []byte("garbage")},
		})
		require.Error(t, err)
		require.Nil(t, client)
	})
}

// TestHTTPClient_MtlsPresentsCertificateToServer wires up an httptest TLS server
// that requires a client certificate signed by a trusted CA, and verifies that
// the http client created via the mTLS path successfully presents the supplied
// client certificate end-to-end.
//
// We reach into the constructed client to install the server's self-signed CA
// into RootCAs after construction — this is the only way to test against a
// local httptest server, since the production mTLS path does not (and should
// not) accept arbitrary root CAs from callers.
func TestHTTPClient_MtlsPresentsCertificateToServer(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	// Generate a self-signed client cert that will double as its own CA.
	clientCertPEM, clientKeyPEM, clientCert := pemKeyPair(t, "test-client")
	clientCAPool := x509.NewCertPool()
	clientCAPool.AddCert(clientCert)

	// Server presents its own self-signed cert.
	serverCertPEM, serverKeyPEM, serverCert := pemKeyPair(t, "127.0.0.1", net.IPv4(127, 0, 0, 1))
	serverPair, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	require.NoError(t, err)

	peerCNCh := make(chan string, 1)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NotNil(t, r.TLS, "request must arrive over TLS")                          //nolint:testifylint // require should be allowed in a test file
		require.Len(t, r.TLS.PeerCertificates, 1, "client must present exactly one cert") //nolint:testifylint // require should be allowed in a test file
		peerCNCh <- r.TLS.PeerCertificates[0].Subject.CommonName
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("mtls-ok"))
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverPair},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	defer server.Close()

	serverPool := x509.NewCertPool()
	serverPool.AddCert(serverCert)

	u, err := url.Parse(server.URL)
	require.NoError(t, err)
	host, portStr := u.Hostname(), u.Port()
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)

	client, err := NewHTTPClient(
		HTTPClientConfig{
			AllowedIPs:         []string{host},
			AllowedPorts:       []int{port},
			Mtls:               &gateway.MtlsAuth{Certificate: clientCertPEM, PrivateKey: clientKeyPEM},
			ConcurrencyLimiter: testMtlsLimiter(),
		},
		lggr,
	)
	require.NoError(t, err)

	hc, ok := client.(*httpClient)
	require.True(t, ok)
	transport, ok := underlyingSafeURLClient(t, hc).Client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport.TLSClientConfig)
	transport.TLSClientConfig.RootCAs = serverPool

	resp, err := client.Send(t.Context(), HTTPRequest{
		Method:  "GET",
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("mtls-ok"), resp.Body)
	require.Equal(t, "test-client", <-peerCNCh, "server must observe the supplied client certificate's CN")
}

// TestHTTPClient_NoMtls_RejectedByMtlsServer is the negative twin of the above:
// the same mTLS-requiring server should reject a non-mTLS client.
func TestHTTPClient_NoMtls_RejectedByMtlsServer(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	_, _, clientCAOnly := pemKeyPair(t, "ca-only")
	clientCAPool := x509.NewCertPool()
	clientCAPool.AddCert(clientCAOnly)

	serverCertPEM, serverKeyPEM, serverCert := pemKeyPair(t, "127.0.0.1", net.IPv4(127, 0, 0, 1))
	serverPair, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	require.NoError(t, err)

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverPair},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	defer server.Close()

	u, err := url.Parse(server.URL)
	require.NoError(t, err)
	host, portStr := u.Hostname(), u.Port()
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)

	client, err := NewHTTPClient(
		HTTPClientConfig{
			AllowedIPs:   []string{host},
			AllowedPorts: []int{port},
		}, // no Mtls
		lggr,
	)
	require.NoError(t, err)

	// Trust the server cert so we get past server-cert verification — the TLS
	// failure should come from the missing *client* cert.
	hc := client.(*httpClient)
	transport := underlyingSafeURLClient(t, hc).Client.Transport.(*http.Transport)
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	pool := x509.NewCertPool()
	pool.AddCert(serverCert)
	transport.TLSClientConfig.RootCAs = pool

	_, err = client.Send(t.Context(), HTTPRequest{
		Method:  "GET",
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.Error(t, err, "non-mTLS client must be rejected by an mTLS-requiring server")
}

// TestHTTPClient_MtlsDisablesKeepAlives ensures the mTLS-configured transport
// disables connection reuse, which is the property that prevents auth'd
// connections from being shared across users.
func TestHTTPClient_MtlsDisablesKeepAlives(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	certPEM, keyPEM, _ := pemKeyPair(t, "client")

	client, err := NewHTTPClient(HTTPClientConfig{
		Mtls:               &gateway.MtlsAuth{Certificate: certPEM, PrivateKey: keyPEM},
		ConcurrencyLimiter: testMtlsLimiter(),
	}, lggr)
	require.NoError(t, err)

	hc, ok := client.(*httpClient)
	require.True(t, ok)

	sc := underlyingSafeURLClient(t, hc)
	require.NotNil(t, sc.Client)
	transport, ok := sc.Client.Transport.(*http.Transport)
	require.True(t, ok, "expected *http.Transport, got %T", sc.Client.Transport)

	require.True(t, transport.DisableKeepAlives, "keep-alives must be disabled to prevent auth'd connection reuse across users")
	require.Equal(t, 10*time.Second, transport.TLSHandshakeTimeout, "TLS handshake timeout should be set (safeurl defaults to 0 == no timeout)")
	require.NotNil(t, transport.TLSClientConfig)
	require.Len(t, transport.TLSClientConfig.Certificates, 1, "client certificate must be installed on the transport")
}

func TestNewHTTPClient_MtlsRequiresConcurrencyLimiter(t *testing.T) {
	t.Parallel()
	certPEM, keyPEM, _ := pemKeyPair(t, "client")

	_, err := NewHTTPClient(HTTPClientConfig{
		Mtls: &gateway.MtlsAuth{Certificate: certPEM, PrivateKey: keyPEM},
	}, logger.Test(t))
	require.Error(t, err)
	require.Contains(t, err.Error(), "mtls requires a ConcurrencyLimiter")
}

// TestConcurrencyLimitedClient_BlocksWhenSaturated verifies a held slot blocks a
// second request until its context expires, and that completing the first frees
// the slot.
func TestConcurrencyLimitedClient_BlocksWhenSaturated(t *testing.T) {
	t.Parallel()

	entered := make(chan struct{})
	release := make(chan struct{})
	wrapped := &concurrencyLimitedClient{
		client: doerFunc(func(*http.Request) (*http.Response, error) {
			entered <- struct{}{}
			<-release
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
		}),
		limiter: limits.GlobalResourcePoolLimiter[int](1),
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	// First request takes the only slot and blocks inside Do.
	done := make(chan struct{})
	go func() { defer close(done); _, _ = wrapped.Do(req) }()
	<-entered

	avail, err := wrapped.limiter.Available(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, avail, "the single slot should be held")

	// Second request is denied once its (short) context expires.
	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = wrapped.Do(req.Clone(ctx2))
	require.ErrorIs(t, err, ErrBlockedRequest)
	require.Contains(t, err.Error(), "mtls concurrency limit exceeded")

	// Releasing the first request frees the slot.
	close(release)
	<-done
	avail, err = wrapped.limiter.Available(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, avail, "slot must be freed after the request completes")
}

// TestHTTPClient_MtlsConcurrencyLimit_RespectsRequestTimeout proves the wait is
// bounded by the (capped) request timeout rather than blocking indefinitely:
// with a zero-capacity pool no slot can ever be acquired, so Send must return
// the concurrency error around the request timeout.
func TestHTTPClient_MtlsConcurrencyLimit_RespectsRequestTimeout(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	certPEM, keyPEM, _ := pemKeyPair(t, "client")

	client, err := NewHTTPClient(HTTPClientConfig{
		Mtls:               &gateway.MtlsAuth{Certificate: certPEM, PrivateKey: keyPEM},
		ConcurrencyLimiter: limits.GlobalResourcePoolLimiter[int](0),
	}, lggr)
	require.NoError(t, err)

	start := time.Now()
	_, err = client.Send(t.Context(), HTTPRequest{
		Method:  http.MethodGet,
		URL:     "https://example.com",
		Timeout: 100 * time.Millisecond,
	})
	require.ErrorIs(t, err, ErrBlockedRequest)
	require.Contains(t, err.Error(), "mtls concurrency limit exceeded")
	require.Less(t, time.Since(start), time.Second, "wait must be bounded by the request timeout")
}
