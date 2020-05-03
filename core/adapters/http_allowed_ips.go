package adapters

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isRestrictedIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() ||
		ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.Equal(net.IPv4bcast) ||
		ip.Equal(net.IPv4allsys) ||
		ip.Equal(net.IPv4allrouter) ||
		ip.Equal(net.IPv4zero) ||
		ip.IsMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// restrictedDialContext wraps the Dialer such that after successful connection,
// we check the IP.
// If the resolved IP is restricted, close the connection and return an error.
func restrictedDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	con, err := (&net.Dialer{
		// Defaults from GoLang standard http package
		// https://golang.org/pkg/net/http/#RoundTripper
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext(ctx, network, address)
	if err == nil {
		// If a connection could be established, ensure its not local or private
		a, _ := con.RemoteAddr().(*net.TCPAddr)

		if isRestrictedIP(a.IP) {
			defer logger.ErrorIfCalling(con.Close)
			return nil, fmt.Errorf("disallowed IP %s. Connections to local/private and multicast networks are disabled by default for security reasons. If you really want to allow this, consider using the httpgetwithunrestrictednetworkaccess or httppostwithunrestrictednetworkaccess adapter instead", a.IP.String())
		}
	}
	return con, err
}
