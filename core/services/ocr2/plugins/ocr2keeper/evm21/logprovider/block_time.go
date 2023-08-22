package logprovider

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	defaultSampleSize = int64(10)
	defaultBlockTime  = time.Second * 1
)

type blockTimeResolver struct {
	poller logpoller.LogPoller
}

func newBlockTimeResolver(poller logpoller.LogPoller) *blockTimeResolver {
	return &blockTimeResolver{
		poller: poller,
	}
}

func (r *blockTimeResolver) BlockTime(ctx context.Context, blockSampleSize int64) (time.Duration, error) {
	if blockSampleSize < 2 { // min 2 blocks range
		blockSampleSize = defaultSampleSize
	}

	latest, err := r.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block from poller: %w", err)
	}
	if latest < blockSampleSize {
		return defaultBlockTime, nil
	}
	blockTimes, err := r.getSampleTimestamps(ctx, blockSampleSize, latest)
	if err != nil {
		return 0, err
	}

	var sumDiff time.Duration
	for i := range blockTimes {
		if i != int(blockSampleSize-1) {
			sumDiff += blockTimes[i].Sub(blockTimes[i+1])
		}
	}

	return sumDiff / time.Duration(blockSampleSize-1), nil
}

func (r *blockTimeResolver) getSampleTimestamps(ctx context.Context, blockSampleSize, latest int64) ([]time.Time, error) {
	blockSample := make([]uint64, blockSampleSize)
	for i := range blockSample {
		blockSample[i] = uint64(latest - blockSampleSize + int64(i))
	}
	blocks, err := r.poller.GetBlocksRange(ctx, blockSample)
	if err != nil {
		return nil, fmt.Errorf("failed to get block range from poller: %w", err)
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].BlockNumber > blocks[j].BlockNumber
	})
	blockTimes := make([]time.Time, blockSampleSize)
	for i, b := range blocks {
		blockTimes[i] = b.BlockTimestamp
	}
	return blockTimes, nil
}
