//go:build !go1.26

package network

var (
	expErrIPV6Blocked = "ipv6 blocked"
	expErrNotAllowed  = "not found in allowlist"
)
