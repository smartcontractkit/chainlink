package limiter

import (
	"context"
)

// -----------------------------------------------------------------
// Context
// -----------------------------------------------------------------

// Context is the limit context.
type Context struct {
	Limit     int64
	Remaining int64
	Reset     int64
	Reached   bool
}

// -----------------------------------------------------------------
// Limiter
// -----------------------------------------------------------------

// Limiter is the limiter instance.
type Limiter struct {
	Store   Store
	Rate    Rate
	Options Options
}

// New returns an instance of Limiter.
func New(store Store, rate Rate, options ...Option) *Limiter {
	opt := Options{
		IPv4Mask:           DefaultIPv4Mask,
		IPv6Mask:           DefaultIPv6Mask,
		TrustForwardHeader: false,
	}
	for _, o := range options {
		o(&opt)
	}
	return &Limiter{
		Store:   store,
		Rate:    rate,
		Options: opt,
	}
}

// Get returns the limit for given identifier.
func (limiter *Limiter) Get(ctx context.Context, key string) (Context, error) {
	return limiter.Store.Get(ctx, key, limiter.Rate)
}

// Peek returns the limit for given identifier, without modification on current values.
func (limiter *Limiter) Peek(ctx context.Context, key string) (Context, error) {
	return limiter.Store.Peek(ctx, key, limiter.Rate)
}

// Reset sets the limit for given identifier to zero.
func (limiter *Limiter) Reset(ctx context.Context, key string) (Context, error) {
	return limiter.Store.Reset(ctx, key, limiter.Rate)
}

// Increment increments the limit by given count & gives back the new limit for given identifier
func (limiter *Limiter) Increment(ctx context.Context, key string, count int64) (Context, error) {
	return limiter.Store.Increment(ctx, key, count, limiter.Rate)
}
