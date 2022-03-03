package monitoring

import (
	"context"
	"errors"
)

var (
	// ErrNoUpdate is an error value interpreted by a Poller to mean that the
	// Fetch() was successful but a new value was not found.
	// The pollers will not report this as an error!
	ErrNoUpdate = errors.New("no updates found")
)

// Source is an abstraction for reading data from a remote API, usually a chain RPC endpoint.
type Source interface {
	// Fetch must be thread-safe!
	// There is no guarantee on the ordering of Fetch() calls for the same source instance.
	Fetch(context.Context) (interface{}, error)
}

type SourceFactory interface {
	NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error)
	// GetType should return a namespace for all the source instances produced by this factory.
	GetType() string
}
