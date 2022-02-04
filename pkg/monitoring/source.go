package monitoring

import (
	"context"
)

type Source interface {
	Fetch(context.Context) (interface{}, error)
}

type SourceFactory interface {
	NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error)
}
