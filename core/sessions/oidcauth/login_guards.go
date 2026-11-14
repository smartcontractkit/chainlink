package oidcauth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxDeviceFlowLifetime = 5 * time.Minute

// deviceReasonNoRole is the failure reason for a verified identity whose
// groups map to no node role.
const deviceReasonNoRole = "no_role"

// deviceReasonMessages are the operator-facing texts for each device-flow
// failure reason. The CLI keeps its own copy and never prints node text.
var deviceReasonMessages = map[string]string{
	deviceReasonNoRole: "Signed in, but none of your identity provider groups map to a node role",
	"denied":           "Authorization was denied at the identity provider",
	"expired":          "Authorization expired before it was approved",
	"abandoned":        "Authorization was abandoned",
	"failed":           "Authorization was not completed",
}

// capDeviceExpiry bounds a provider expiry. A missing expiry, or one longer
// than maxDeviceFlowLifetime, is shortened. An already-past expiry is left
// past so the caller can fail the flow instead of extending it.
func capDeviceExpiry(expiry, now time.Time) time.Time {
	if !expiry.IsZero() && !expiry.After(now) {
		return expiry
	}
	deadline := now.Add(maxDeviceFlowLifetime)
	if expiry.IsZero() || expiry.After(deadline) {
		return deadline
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

// clientResolver works out the client address for the pre-authentication
// OIDC routes. X-Forwarded-For is believed only when the hop that sent it is
// inside trustedProxies. The header is walked from the right, each trusted
// hop is skipped, and the first address outside the trusted set is the
// client. With no trusted proxies the socket peer is the client.
type clientResolver struct {
	trustedProxies []netip.Prefix
}

// parseCIDRs parses config prefixes. Config validation already rejects bad
// entries, so an error here means the config bypassed validation.
func parseCIDRs(name string, list []string) ([]netip.Prefix, error) {
	out := make([]netip.Prefix, 0, len(list))
	for _, v := range list {
		p, err := netip.ParsePrefix(strings.TrimSpace(v))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid CIDR %q: %w", name, v, err)
		}
		out = append(out, p.Masked())
	}
	return out, nil
}

func prefixesContain(list []netip.Prefix, addr netip.Addr) bool {
	if !addr.IsValid() {
		return false
	}
	return slices.ContainsFunc(list, func(p netip.Prefix) bool { return p.Contains(addr) })
}

// clientAddr returns the client address, or the zero Addr when it cannot be
// determined (unparseable peer, or a malformed entry in a trusted hop's
// X-Forwarded-For). Callers treat the zero Addr as untrusted.
func (r clientResolver) clientAddr(c *gin.Context) netip.Addr {
	if c.Request == nil {
		return netip.Addr{}
	}
	current := parseHostAddr(c.Request.RemoteAddr)
	if !current.IsValid() || !prefixesContain(r.trustedProxies, current) {
		return current
	}
	var entries []string
	for _, line := range c.Request.Header.Values("X-Forwarded-For") {
		entries = append(entries, strings.Split(line, ",")...)
	}
	for _, entry := range slices.Backward(entries) {
		next := parseHostAddr(strings.TrimSpace(entry))
		if !next.IsValid() {
			// A trusted hop passed on something that is not an address. Do not
			// fall back to the proxy's own address: that would let any caller
			// inherit the proxy's allow-list standing.
			return netip.Addr{}
		}
		current = next
		if !prefixesContain(r.trustedProxies, current) {
			return current
		}
	}
	// Every hop was a trusted proxy, so the leftmost one is the client.
	return current
}

// parseHostAddr accepts "ip", "ip:port" and "[ipv6]:port".
func parseHostAddr(s string) netip.Addr {
	if s == "" {
		return netip.Addr{}
	}
	if a, err := netip.ParseAddr(s); err == nil {
		return a.Unmap().WithZone("")
	}
	if ap, err := netip.ParseAddrPort(s); err == nil {
		return ap.Addr().Unmap().WithZone("")
	}
	return netip.Addr{}
}

// clientBucket is the key the rate limiters and concurrency caps count
// against. IPv6 clients are grouped by /64, the smallest block an ISP or
// cloud normally hands out, so one allocation cannot mint fresh buckets.
// Unresolvable clients share one bucket.
func clientBucket(addr netip.Addr) string {
	switch {
	case !addr.IsValid():
		return "unknown"
	case addr.Is4():
		return addr.String()
	default:
		return netip.PrefixFrom(addr, 64).Masked().String()
	}
}

func addrString(addr netip.Addr) string {
	if !addr.IsValid() {
		return ""
	}
	return addr.String()
}

// handleFingerprint is a short, non-reversible id for a device handle so
// audit events for one flow can be joined without logging the handle.
func handleFingerprint(handle string) string {
	sum := sha256.Sum256([]byte(handle))
	return hex.EncodeToString(sum[:8])
}

// truncateUserAgent bounds an attacker-supplied header before it reaches the
// audit log.
func truncateUserAgent(ua string) string {
	const maxUserAgent = 256
	if len(ua) > maxUserAgent {
		return ua[:maxUserAgent]
	}
	return ua
}

func deviceFailureReason(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, errNoMatchingRole) {
		return deviceReasonNoRole
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

func errorsIsDeadline(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "context canceled"))
}
