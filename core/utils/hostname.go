package utils

import (
	"net"
	"regexp"
	"strconv"
)

// Adapted from: https://github.com/go-playground/validator/blob/master/baked_in.go

var (
	// accepts hostname starting with a digit https://tools.ietf.org/html/rfc1123
	hostnameRegexRFC1123 = regexp.MustCompile(`^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*?$`)
)

// IsHostnamePort validates a <dns>:<port> combination for fields typically used for socket address.
func IsHostnamePort(val string) bool {
	host, port, err := net.SplitHostPort(val)
	if err != nil {
		return false
	}
	// Port must be a int <= 65535.
	if portNum, err := strconv.ParseInt(port, 10, 32); err != nil || portNum > 65535 || portNum < 1 {
		return false
	}

	// If host is specified, it should match a DNS name
	if host != "" {
		return hostnameRegexRFC1123.MatchString(host)
	}
	return true
}
