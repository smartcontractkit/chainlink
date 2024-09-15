package cache_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
)

func TestLogpollerEventsBased(t *testing.T) {
	ctx := testutils.Context(t)
	lp := lpmocks.NewLogPoller(t)
	observedEvents := []common.Hash{
		utils.Bytes32FromString("event a"),
		utils.Bytes32FromString("event b"),
	}
	contractAddress := utils.RandomAddress()
	c := cache.NewLogpollerEventsBased[[]int](lp, observedEvents, contractAddress)

	testRounds := []struct {
		logPollerLatestBlock int64 // latest block that logpoller parsed
		latestEventBlock     int64 // latest block that an event was seen
		stateLatestBlock     int64 // block of the current cached value (before run)
		shouldSync           bool  // whether we expect sync to happen in this round
		syncData             []int // data returned after sync
		expData              []int // expected data that cache will return
	}{
		{
			// this is the first 'Get' call to our cache, an event was seen at block 800
			// and now log poller has reached block 1000.
			logPollerLatestBlock: 1000,
			latestEventBlock:     800,
			stateLatestBlock:     0,
			shouldSync:           true,
			syncData:             []int{1, 2, 3},
			expData:              []int{1, 2, 3},
		},
		{
			// log poller moved a few blocks and there weren't any new events
			logPollerLatestBlock: 1010,
			latestEventBlock:     800,
			stateLatestBlock:     1000,
			shouldSync:           false,
			expData:              []int{1, 2, 3},
		},
		{
			// log poller moved a few blocks and there was a new event
			logPollerLatestBlock: 1020,
			latestEventBlock:     1020,
			stateLatestBlock:     1010,
			shouldSync:           true,
			syncData:             []int{111},
			expData:              []int{111},
		},
		{
			// log poller moved a few more blocks and there was another new event
			logPollerLatestBlock: 1050,
			latestEventBlock:     1040,
			stateLatestBlock:     1020,
			shouldSync:           true,
			syncData:             []int{222},
			expData:              []int{222},
		},
		{
			// log poller moved a few more blocks and there wasn't any new event
			logPollerLatestBlock: 1100,
			latestEventBlock:     1040,
			stateLatestBlock:     1050,
			shouldSync:           false,
			expData:              []int{222},
		},
		{
			// log poller moved a few more blocks and there wasn't any new event
			logPollerLatestBlock: 1300,
			latestEventBlock:     1040,
			stateLatestBlock:     1100,
			shouldSync:           false,
			expData:              []int{222},
		},
		{
			// log poller moved a few more blocks and there was a new event
			// more recent than latest block (for whatever internal reason)
			logPollerLatestBlock: 1300,
			latestEventBlock:     1305,
			stateLatestBlock:     1300,
			shouldSync:           true,
			syncData:             []int{666},
			expData:              []int{666},
		},
		{
			// log poller moved a few more blocks and there wasn't any new event
			logPollerLatestBlock: 1300,
			latestEventBlock:     1305,
			stateLatestBlock:     1305, // <-- that's what we are testing in this round
			shouldSync:           false,
			expData:              []int{666},
		},
	}

	for _, round := range testRounds {
		lp.On("LatestBlock", mock.Anything).
			Return(logpoller.LogPollerBlock{FinalizedBlockNumber: round.logPollerLatestBlock}, nil).Once()

		if round.stateLatestBlock > 0 {
			lp.On(
				"LatestBlockByEventSigsAddrsWithConfs",
				mock.Anything,
				round.stateLatestBlock,
				observedEvents,
				[]common.Address{contractAddress},
				evmtypes.Finalized,
			).Return(round.latestEventBlock, nil).Once()
		}

		data, err := c.Get(ctx, func(ctx context.Context) ([]int, error) { return round.syncData, nil })
		assert.NoError(t, err)
		assert.Equal(t, round.expData, data)
	}
}
