package oidcauth

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestCapDeviceExpiry(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)

	assert.Equal(t, now.Add(maxDeviceFlowLifetime), capDeviceExpiry(time.Time{}, now))
	assert.Equal(t, now.Add(maxDeviceFlowLifetime), capDeviceExpiry(now.Add(time.Hour), now))
	assert.Equal(t, now.Add(2*time.Minute), capDeviceExpiry(now.Add(2*time.Minute), now))
	past := now.Add(-time.Second)
	assert.Equal(t, past, capDeviceExpiry(past, now))
}

func TestVerificationURIAllowed(t *testing.T) {
	t.Parallel()
	const provider = "https://id.example.com/oauth2/default"

	assert.True(t, verificationURIAllowed(provider, "https://id.example.com/activate"))
	assert.True(t, verificationURIAllowed(provider, "https://ID.EXAMPLE.COM/activate?user_code=WDJB"))
	assert.False(t, verificationURIAllowed(provider, "https://evil.example/activate"))
	assert.False(t, verificationURIAllowed(provider, "http://id.example.com/activate"))
	assert.False(t, verificationURIAllowed(provider, ""))
	assert.False(t, verificationURIAllowed(provider, "not a url"))

	assert.True(t, verificationURIAllowed("http://127.0.0.1:8080", "http://127.0.0.1:8080/activate"))
	assert.True(t, verificationURIAllowed("http://localhost:8080", "http://localhost:8080/activate"))
	assert.False(t, verificationURIAllowed("http://10.0.0.1:8080", "http://10.0.0.1:8080/activate"))
}

func TestClientResolver_NoTrustedProxies(t *testing.T) {
	t.Parallel()
	r := clientResolver{}

	// With no trusted proxies the header is ignored, even from a private or
	// loopback peer. This is the forgery Halborn showed behind an L4 hop.
	for _, peer := range []string{"203.0.113.8:1234", "127.0.0.1:1234", "10.0.0.7:5555", "[fe80::1]:443"} {
		c := deviceRequest(t, peer, "198.51.100.1, 203.0.113.9")
		host, _, _ := net.SplitHostPort(peer)
		assert.Equal(t, host, r.clientAddr(c).String(), peer)
	}
}

func TestClientResolver_TrustedProxies(t *testing.T) {
	t.Parallel()
	r := clientResolver{trustedProxies: mustCIDRs(t, "10.0.0.0/8", "fd00::/8")}

	// One trusted hop: the rightmost entry is the client.
	assert.Equal(t, "203.0.113.9", r.clientAddr(deviceRequest(t, "10.0.0.5:443", "198.51.100.1, 203.0.113.9")).String())

	// Two trusted hops: skip the inner proxy and take the first untrusted entry.
	assert.Equal(t, "203.0.113.9", r.clientAddr(deviceRequest(t, "10.0.0.5:443", "198.51.100.1, 203.0.113.9, 10.1.2.3")).String())

	// Forging extra entries on the left changes nothing.
	a := r.clientAddr(deviceRequest(t, "10.0.0.5:443", "203.0.113.1, 203.0.113.9"))
	b := r.clientAddr(deviceRequest(t, "10.0.0.5:443", "203.0.113.2, 203.0.113.9"))
	assert.Equal(t, a, b)

	// A peer outside the trusted set is the client and its header is ignored.
	assert.Equal(t, "198.51.100.9", r.clientAddr(deviceRequest(t, "198.51.100.9:1", "203.0.113.1")).String())

	// A private peer that is not a configured proxy is also not believed.
	assert.Equal(t, "192.168.1.5", r.clientAddr(deviceRequest(t, "192.168.1.5:1", "203.0.113.1")).String())

	// Trusted peer without a header: the peer is the client.
	assert.Equal(t, "10.0.0.5", r.clientAddr(deviceRequest(t, "10.0.0.5:443", "")).String())

	// Every hop trusted: the leftmost is the client.
	assert.Equal(t, "10.9.9.9", r.clientAddr(deviceRequest(t, "10.0.0.5:443", "10.9.9.9, 10.1.1.1")).String())

	// A malformed entry from a trusted hop resolves to no address, never to
	// the proxy's own address.
	assert.False(t, r.clientAddr(deviceRequest(t, "10.0.0.5:443", "203.0.113.4, not-an-ip")).IsValid())

	// Several header lines are one list.
	c := deviceRequest(t, "10.0.0.5:443", "198.51.100.1")
	c.Request.Header.Add("X-Forwarded-For", "203.0.113.50")
	assert.Equal(t, "203.0.113.50", r.clientAddr(c).String())

	// ip:port entries and IPv4-mapped IPv6 normalize.
	assert.Equal(t, "203.0.113.7", r.clientAddr(deviceRequest(t, "10.0.0.5:443", "203.0.113.7:5555")).String())
	assert.Equal(t, "203.0.113.7", r.clientAddr(deviceRequest(t, "[::ffff:10.0.0.5]:443", "::ffff:203.0.113.7")).String())
}

func TestClientBucket(t *testing.T) {
	t.Parallel()
	// Addresses in one IPv6 /64 share a bucket.
	assert.Equal(t, clientBucket(netip.MustParseAddr("2001:db8::1")), clientBucket(netip.MustParseAddr("2001:db8::2")))
	assert.Equal(t, "2001:db8::/64", clientBucket(netip.MustParseAddr("2001:db8::ffff:1")))
	assert.NotEqual(t, clientBucket(netip.MustParseAddr("2001:db8:0:1::1")), clientBucket(netip.MustParseAddr("2001:db8:0:2::1")))
	// IPv4 stays per address.
	assert.Equal(t, "203.0.113.1", clientBucket(netip.MustParseAddr("203.0.113.1")))
	assert.Equal(t, "unknown", clientBucket(netip.Addr{}))
}

func TestParseCIDRs(t *testing.T) {
	t.Parallel()
	got, err := parseCIDRs("x", []string{"10.1.2.3/8", " 2001:db8::/32 "})
	require.NoError(t, err)
	assert.Equal(t, []netip.Prefix{netip.MustParsePrefix("10.0.0.0/8"), netip.MustParsePrefix("2001:db8::/32")}, got)
	_, err = parseCIDRs("x", []string{"10.0.0.1"})
	require.Error(t, err)
}

func TestHandleFingerprint(t *testing.T) {
	t.Parallel()
	fp := handleFingerprint("some-handle")
	assert.Len(t, fp, 16)
	assert.NotContains(t, fp, "some-handle")
	assert.Equal(t, fp, handleFingerprint("some-handle"))
	assert.NotEqual(t, fp, handleFingerprint("other-handle"))
}

func TestDeviceNetworkGuard(t *testing.T) {
	t.Parallel()
	run := func(oi *oidcAuthenticator, remote, xff string) int {
		w := httptest.NewRecorder()
		_, eng := gin.CreateTestContext(w)
		eng.POST("/oidc-device/start", oi.deviceNetworkGuard(), func(c *gin.Context) { c.Status(http.StatusOK) })
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/oidc-device/start", nil)
		req.RemoteAddr = remote
		if xff != "" {
			req.Header.Set("X-Forwarded-For", xff)
		}
		eng.ServeHTTP(w, req)
		return w.Code
	}
	lggr := logger.TestSugared(t)

	off := &oidcAuthenticator{lggr: lggr}
	assert.Equal(t, http.StatusForbidden, run(off, "10.8.0.4:1", ""), "empty allow-list turns device login off")

	on := &oidcAuthenticator{
		lggr:          lggr,
		deviceAllowed: mustCIDRs(t, "10.8.0.0/16"),
		clients:       clientResolver{trustedProxies: mustCIDRs(t, "10.0.0.0/24")},
	}
	assert.Equal(t, http.StatusOK, run(on, "10.8.0.4:1", ""))
	assert.Equal(t, http.StatusForbidden, run(on, "203.0.113.5:1", ""))
	// Through the trusted load balancer: the client decides.
	assert.Equal(t, http.StatusOK, run(on, "10.0.0.2:1", "10.8.3.3"))
	assert.Equal(t, http.StatusForbidden, run(on, "10.0.0.2:1", "203.0.113.5"))
	// A public caller cannot claim to be inside by forging the header.
	assert.Equal(t, http.StatusForbidden, run(on, "203.0.113.5:1", "10.8.3.3"))
	// A malformed header from the trusted proxy is refused.
	assert.Equal(t, http.StatusForbidden, run(on, "10.0.0.2:1", "garbage"))
}

func mustCIDRs(t *testing.T, list ...string) []netip.Prefix {
	t.Helper()
	p, err := parseCIDRs("test", list)
	require.NoError(t, err)
	return p
}

func TestDeviceFailureReason(t *testing.T) {
	t.Parallel()
	require.Empty(t, deviceFailureReason(nil))
	assert.Equal(t, "failed", deviceFailureReason(assert.AnError))
	assert.Equal(t, "denied", deviceFailureReason(stringError("access_denied")))
	assert.Equal(t, "abandoned", deviceFailureReason(errDeviceFlowAbandoned))
	assert.Equal(t, "expired", deviceFailureReason(stringError("context deadline exceeded")))
	assert.Equal(t, deviceReasonNoRole, deviceFailureReason(fmt.Errorf("%w: %w", errNoMatchingRole, ErrUserNoOIDCGroups)))
	for _, reason := range []string{deviceReasonNoRole, "denied", "expired", "abandoned", "failed"} {
		assert.NotEmpty(t, deviceReasonMessages[reason], reason)
	}
}

type stringError string

func (e stringError) Error() string { return string(e) }

func deviceRequest(t *testing.T, remoteAddr, xff string) *gin.Context {
	t.Helper()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/oidc-device/start", nil)
	c.Request.RemoteAddr = remoteAddr
	if xff != "" {
		c.Request.Header.Set("X-Forwarded-For", xff)
	}
	return c
}
