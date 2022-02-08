package monitoring

import (
	"context"
	"errors"
)

var (
	// ErrNoUpdate is an error value to suggest to the Pollers that the source
	// has not found an update and that's acceptable and should be treated as an error.
	ErrNoUpdate = errors.New("no updates found")
)

type Source interface {
	Fetch(context.Context) (interface{}, error)
}

type SourceFactory interface {
	NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error)
}
