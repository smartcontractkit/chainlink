package logprovider

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestBlockTimeResolver_BlockTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		blockSampleSize int64
		latestBlock     int64
		latestBlockErr  error
		blocksRange     []logpoller.Block
		blocksRangeErr  error
		blockTime       time.Duration
		blockTimeErr    error
	}{
		{
			"latest block err",
			10,
			0,
			errors.New("test err"),
			nil,
			nil,
			0,
			errors.New("test err"),
		},
		{
			"block range err",
			10,
			20,
			nil,
			nil,
			errors.New("test err"),
			0,
			errors.New("test err"),
		},
		{
			"2 sec block time",
			4,
			20,
			nil,
			[]logpoller.Block{
				{BlockTimestamp: now.Add(-time.Second * (2 * 4)), BlockNumber: 16},
				{BlockTimestamp: now, BlockNumber: 20},
			},
			nil,
			2 * time.Second,
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			lp := new(lpmocks.LogPoller)
			resolver := newBlockTimeResolver(lp)

			lp.On("LatestBlock", mock.Anything).Return(logpoller.Block{BlockNumber: tc.latestBlock}, tc.latestBlockErr)
			lp.On("GetBlocksRange", mock.Anything, mock.Anything).Return(tc.blocksRange, tc.blocksRangeErr)

			blockTime, err := resolver.BlockTime(ctx, tc.blockSampleSize)
			if tc.blockTimeErr != nil {
				require.Error(t, err)
				return
			}
			require.Equal(t, tc.blockTime, blockTime)
		})
	}
}
