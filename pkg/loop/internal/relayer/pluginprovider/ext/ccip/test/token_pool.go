package test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// TokenPoolBatchedReader is a static implementation of the TokenPoolBatchedReaderEvaluator interface.
var TokenPoolBatchedReader = staticTokenPoolBatchedReader{
	staticTokenPoolBatchedReaderConfig{
		getInboundTokenPoolRateLimitsRequest: []cciptypes.Address{
			cciptypes.Address("token pool batched reader request1"),
			cciptypes.Address("token pool batched reader request2"),
		},
		getInboundTokenPoolRateLimitsResponse: []cciptypes.TokenBucketRateLimit{
			{
				Tokens:      big.NewInt(7),
				LastUpdated: 3,
				IsEnabled:   true,
				Capacity:    big.NewInt(5),
				Rate:        big.NewInt(3),
			},
			{
				Tokens:      big.NewInt(11),
				LastUpdated: 3,
				IsEnabled:   true,
				Capacity:    big.NewInt(13),
				Rate:        big.NewInt(17),
			},
		},
	},
}

type TokenPoolBatchedReaderEvaluator interface {
	cciptypes.TokenPoolBatchedReader
	testtypes.Evaluator[cciptypes.TokenPoolBatchedReader]
}
type staticTokenPoolBatchedReader struct {
	staticTokenPoolBatchedReaderConfig
}

var _ TokenPoolBatchedReaderEvaluator = staticTokenPoolBatchedReader{}

// Close implements ccip.TokenPoolBatchedReader.
func (s staticTokenPoolBatchedReader) Close() error {
	return nil
}

// Evaluate implements types_test.Evaluator.
func (s staticTokenPoolBatchedReader) Evaluate(ctx context.Context, other cciptypes.TokenPoolBatchedReader) error {
	got, err := other.GetInboundTokenPoolRateLimits(ctx, s.getInboundTokenPoolRateLimitsRequest)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(got, s.getInboundTokenPoolRateLimitsResponse) {
		return fmt.Errorf("got %v, want %v", got, s.getInboundTokenPoolRateLimitsResponse)
	}
	return nil
}

// GetInboundTokenPoolRateLimits implements TokenPoolBatchedReaderEvaluator.
func (s staticTokenPoolBatchedReader) GetInboundTokenPoolRateLimits(ctx context.Context, tokenPoolReaders []cciptypes.Address) ([]cciptypes.TokenBucketRateLimit, error) {
	if !assert.ObjectsAreEqualValues(tokenPoolReaders, s.getInboundTokenPoolRateLimitsRequest) {
		return nil, fmt.Errorf("got %v, want %v", tokenPoolReaders, s.getInboundTokenPoolRateLimitsRequest)
	}
	return s.getInboundTokenPoolRateLimitsResponse, nil
}

type staticTokenPoolBatchedReaderConfig struct {
	getInboundTokenPoolRateLimitsRequest  cciptypes.Addresses
	getInboundTokenPoolRateLimitsResponse []cciptypes.TokenBucketRateLimit
}
