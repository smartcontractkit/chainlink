package logprovider

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

var (
	defaultSampleSize = int64(10000)
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

	latest, err := r.poller.LatestBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block from poller: %w", err)
	}
	latestBlockNumber := latest.BlockNumber
	if latestBlockNumber <= blockSampleSize {
		return defaultBlockTime, nil
	}
	start, end := latestBlockNumber-blockSampleSize, latestBlockNumber
	startTime, endTime, err := r.getSampleTimestamps(ctx, uint64(start), uint64(end))
	if err != nil {
		return 0, err
	}

	return endTime.Sub(startTime) / time.Duration(blockSampleSize), nil
}

func (r *blockTimeResolver) getSampleTimestamps(ctx context.Context, start, end uint64) (time.Time, time.Time, error) {
	blocks, err := r.poller.GetBlocksRange(ctx, []uint64{start, end})
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to get block range from poller: %w", err)
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].BlockNumber < blocks[j].BlockNumber
	})
	if len(blocks) < 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to fetch blocks %d, %d from log poller", start, end)
	}
	return blocks[0].BlockTimestamp, blocks[1].BlockTimestamp, nil
}
