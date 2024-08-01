package ccip

import (
	"context"
	"io"
)

type TokenPoolBatchedReader interface {
	GetInboundTokenPoolRateLimits(ctx context.Context, tokenPoolReaders []Address) ([]TokenBucketRateLimit, error)
	io.Closer
}
