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
	TrustForwardHeader bool
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
func WithTrustForwardHeader(enable bool) Option {
	return func(o *Options) {
		o.TrustForwardHeader = enable
	}
}
