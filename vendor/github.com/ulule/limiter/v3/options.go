package limiter

import (
	"net"
)

// Option is a functional option.
type Option func(*Options)

// Options are limiter options.
type Options struct {
	// IPv4Mask defines the mask used to obtain a IPv4 address.
	IPv4Mask net.IPMask
	// IPv6Mask defines the mask used to obtain a IPv6 address.
	IPv6Mask net.IPMask
	// TrustForwardHeader enable parsing of X-Real-IP and X-Forwarded-For headers to obtain user IP.
	// Please be advised that using this option could be insecure (ie: spoofed) if your reverse
	// proxy is not configured properly to forward a trustworthy client IP.
	// Please read the section "Limiter behind a reverse proxy" in the README for further information.
	TrustForwardHeader bool
	// ClientIPHeader defines a custom header (likely defined by your CDN or Cloud provider) to obtain user IP.
	// If configured, this option will override "TrustForwardHeader" option.
	// Please be advised that using this option could be insecure (ie: spoofed) if your reverse
	// proxy is not configured properly to forward a trustworthy client IP.
	// Please read the section "Limiter behind a reverse proxy" in the README for further information.
	ClientIPHeader string
}

// WithIPv4Mask will configure the limiter to use given mask for IPv4 address.
func WithIPv4Mask(mask net.IPMask) Option {
	return func(o *Options) {
		o.IPv4Mask = mask
	}
}

// WithIPv6Mask will configure the limiter to use given mask for IPv6 address.
func WithIPv6Mask(mask net.IPMask) Option {
	return func(o *Options) {
		o.IPv6Mask = mask
	}
}

// WithTrustForwardHeader will configure the limiter to trust X-Real-IP and X-Forwarded-For headers.
// Please be advised that using this option could be insecure (ie: spoofed) if your reverse
// proxy is not configured properly to forward a trustworthy client IP.
// Please read the section "Limiter behind a reverse proxy" in the README for further information.
func WithTrustForwardHeader(enable bool) Option {
	return func(o *Options) {
		o.TrustForwardHeader = enable
	}
}

// WithClientIPHeader will configure the limiter to use a custom header to obtain user IP.
// Please be advised that using this option could be insecure (ie: spoofed) if your reverse
// proxy is not configured properly to forward a trustworthy client IP.
// Please read the section "Limiter behind a reverse proxy" in the README for further information.
func WithClientIPHeader(header string) Option {
	return func(o *Options) {
		o.ClientIPHeader = header
	}
}
