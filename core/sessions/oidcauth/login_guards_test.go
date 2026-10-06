package oidcauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestDeviceClientIP(t *testing.T) {
	t.Parallel()

	public := deviceRequest(t, "203.0.113.8:1234", "198.51.100.1, 203.0.113.9")
	assert.Equal(t, "203.0.113.8", deviceClientIP(public))

	loopback := deviceRequest(t, "127.0.0.1:1234", "198.51.100.1, 203.0.113.9")
	assert.Equal(t, "203.0.113.9", deviceClientIP(loopback))

	private := deviceRequest(t, "10.0.0.5:443", "203.0.113.4, not-an-ip")
	assert.Equal(t, "203.0.113.4", deviceClientIP(private))

	noHeader := deviceRequest(t, "10.0.0.5:443", "")
	assert.Equal(t, "10.0.0.5", deviceClientIP(noHeader))

	linkLocal := deviceRequest(t, "[fe80::1]:443", "203.0.113.7")
	assert.Equal(t, "203.0.113.7", deviceClientIP(linkLocal))
}

func TestUnexpectedGroup(t *testing.T) {
	t.Parallel()
	allowed := []string{AdminClaim, EditorClaim, RunnerClaim, ReadClaim}
	assert.Empty(t, unexpectedGroup([]string{AdminClaim, EditorClaim}, allowed...))
	assert.Equal(t, "Contractors", unexpectedGroup([]string{AdminClaim, "Contractors"}, allowed...))
	assert.Empty(t, unexpectedGroup(nil, allowed...))
	assert.Equal(t, "Contractors", unexpectedGroup([]string{"Contractors"}, allowed...))
}

func TestDeviceFailureReason(t *testing.T) {
	t.Parallel()
	require.Empty(t, deviceFailureReason(nil))
	assert.Equal(t, "failed", deviceFailureReason(assert.AnError))
	assert.Equal(t, "denied", deviceFailureReason(errString("access_denied")))
	assert.Equal(t, "abandoned", deviceFailureReason(errDeviceFlowAbandoned))
	assert.Equal(t, "expired", deviceFailureReason(errString("context deadline exceeded")))
}

type errString string

func (e errString) Error() string { return string(e) }

func deviceRequest(t *testing.T, remoteAddr, xff string) *gin.Context {
	t.Helper()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/oidc-device/start", nil)
	c.Request.RemoteAddr = remoteAddr
	if xff != "" {
		c.Request.Header.Set("X-Forwarded-For", xff)
	}
	return c
}
