package environment

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

func rewriteAddressHost(rawAddress, host string) (string, error) {
	return rewriteAddressHostWithPolicy(rawAddress, host, true)
}

func rewriteURLHost(rawURL, host string) (string, error) {
	return rewriteAddressHostWithPolicy(rawURL, host, false)
}

func rewriteAddressHostWithPolicy(rawAddress, host string, requireExplicitPort bool) (string, error) {
	trimmed := strings.TrimSpace(rawAddress)
	if trimmed == "" {
		return "", nil
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return "", fmt.Errorf("failed to parse address %q: %w", rawAddress, err)
		}
		port := parsed.Port()
		if port == "" {
			if requireExplicitPort {
				return "", fmt.Errorf("address %q must include a port", rawAddress)
			}
			parsed.Host = host
			return parsed.String(), nil
		}
		parsed.Host = net.JoinHostPort(host, port)
		return parsed.String(), nil
	}
	_, port, err := net.SplitHostPort(trimmed)
	if err != nil {
		return "", fmt.Errorf("failed to parse host:port %q: %w", rawAddress, err)
	}
	return net.JoinHostPort(host, port), nil
}
