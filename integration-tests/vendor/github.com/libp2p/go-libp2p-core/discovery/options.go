package discovery

import "time"

// DiscoveryOpt is a single discovery option.
type Option func(opts *Options) error

// DiscoveryOpts is a set of discovery options.
type Options struct {
	Ttl   time.Duration
	Limit int

	// Other (implementation-specific) options
	Other map[interface{}]interface{}
}

// Apply applies the given options to this DiscoveryOpts
func (opts *Options) Apply(options ...Option) error {
	for _, o := range options {
		if err := o(opts); err != nil {
			return err
		}
	}
	return nil
}

// TTL is an option that provides a hint for the duration of an advertisement
func TTL(ttl time.Duration) Option {
	return func(opts *Options) error {
		opts.Ttl = ttl
		return nil
	}
}

// Limit is an option that provides an upper bound on the peer count for discovery
func Limit(limit int) Option {
	return func(opts *Options) error {
		opts.Limit = limit
		return nil
	}
}
