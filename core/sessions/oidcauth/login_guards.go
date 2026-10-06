package oidcauth

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxDeviceFlowLifetime = 5 * time.Minute

// capDeviceExpiry bounds a provider expiry. A missing expiry, or one longer
// than maxDeviceFlowLifetime, is shortened. An already-past expiry is left
// past so the caller can fail the flow instead of extending it.
func capDeviceExpiry(expiry, now time.Time) time.Time {
	if !expiry.IsZero() && !expiry.After(now) {
		return expiry
	}
	cap := now.Add(maxDeviceFlowLifetime)
	if expiry.IsZero() || expiry.After(cap) {
		return cap
	}
	return expiry
}

// verificationURIAllowed reports whether uri is safe to hand to an operator.
// The host must be the configured provider host, and the scheme must be https
// unless the host is loopback (local providers).
func verificationURIAllowed(providerURL, uri string) bool {
	if uri == "" {
		return false
	}
	provider, err := url.Parse(providerURL)
	if err != nil || provider.Hostname() == "" {
		return false
	}
	parsed, err := url.Parse(uri)
	if err != nil || parsed.Hostname() == "" {
		return false
	}
	if !strings.EqualFold(parsed.Hostname(), provider.Hostname()) {
		return false
	}
	if parsed.Scheme == "https" {
		return true
	}
	if parsed.Scheme != "http" {
		return false
	}
	ip := net.ParseIP(parsed.Hostname())
	return (ip != nil && ip.IsLoopback()) || parsed.Hostname() == "localhost"
}

// deviceClientIP is the address the device-flow caps use. Behind a private
// load balancer the socket peer is the balancer, so the rightmost
// X-Forwarded-For entry is the address that balancer observed. A public peer
// is used as-is. The header is ignored in that case, because the caller can
// set it.
func deviceClientIP(c *gin.Context) string {
	peer := remotePeer(c)
	parsed := net.ParseIP(peer)
	if parsed == nil || !trustedProxyPeer(parsed) {
		if peer != "" {
			return peer
		}
		return c.ClientIP()
	}
	xff := c.GetHeader("X-Forwarded-For")
	if xff == "" {
		return peer
	}
	parts := strings.Split(xff, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		ip := net.ParseIP(strings.TrimSpace(parts[i]))
		if ip != nil {
			return ip.String()
		}
	}
	return peer
}

func remotePeer(c *gin.Context) string {
	if c.Request == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}

func trustedProxyPeer(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}

func deviceFailureReason(err error) string {
	if err == nil {
		return ""
	}
	if errorsIsDeadline(err) {
		return "expired"
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "abandon"):
		return "abandoned"
	case strings.Contains(msg, "expired") || strings.Contains(msg, "timeout"):
		return "expired"
	case strings.Contains(msg, "denied") || strings.Contains(msg, "access_denied"):
		return "denied"
	default:
		return "failed"
	}
}

func unexpectedGroup(claims []string, allowed ...string) string {
	allow := make(map[string]struct{}, len(allowed))
	for _, name := range allowed {
		if name != "" {
			allow[name] = struct{}{}
		}
	}
	for _, claim := range claims {
		if _, ok := allow[claim]; !ok {
			return claim
		}
	}
	return ""
}

func errorsIsDeadline(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "context canceled"))
}
