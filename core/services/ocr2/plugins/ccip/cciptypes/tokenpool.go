package cciptypes

import "context"

type TokenPoolBatchedReader interface {
	GetInboundTokenPoolRateLimits(ctx context.Context, tokenPoolReaders []Address) ([]TokenBucketRateLimit, error)
}
