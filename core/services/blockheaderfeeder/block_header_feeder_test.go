package blockheaderfeeder

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

type testCase struct {
	name                string
	requests            []blockhashstore.Event
	fulfillments        []blockhashstore.Event
	wait                int
	lookback            int
	latest              uint64
	alreadyStored       []uint64
	expectedStored      []uint64
	expectedErrMsg      string
	getBatchSize        uint16
	storeBatchSize      uint16
	getBatchCallCount   uint16
	storeBatchCallCount uint16
	storedEarliest      bool
	bhs                 blockhashstore.TestBHS
	batchBHS            blockhashstore.TestBatchBHS
}

func TestFeeder(t *testing.T) {
	tests := []testCase{
		{
			name:                "single missing block",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{150, 151, 152, 153, 154, 155},
			getBatchSize:        1,
			storeBatchSize:      1,
			getBatchCallCount:   5,
			storeBatchCallCount: 5,
			storedEarliest:      false,
		},
		{
			name:                "multiple missing blocks",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request"}, {Block: 149, ID: "request"}, {Block: 148, ID: "request"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{148, 149, 150, 151, 152, 153, 154, 155},
			getBatchSize:        1,
			storeBatchSize:      1,
			getBatchCallCount:   7,
			storeBatchCallCount: 7,
			storedEarliest:      false,
		},
		{
			name:                "single missing get batch size = 2",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{150, 151, 152, 153, 154, 155},
			getBatchSize:        2,
			storeBatchSize:      1,
			getBatchCallCount:   3,
			storeBatchCallCount: 5,
			storedEarliest:      false,
		},
		{
			name:                "single missing get and store batch size = 3",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{150, 151, 152, 153, 154, 155},
			getBatchSize:        3,
			storeBatchSize:      3,
			getBatchCallCount:   2,
			storeBatchCallCount: 2,
			storedEarliest:      false,
		},
		{
			name:              "single missing block store earliest",
			requests:          []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:              256,
			lookback:          500,
			latest:            450,
			getBatchSize:      10,
			getBatchCallCount: 5,
			storedEarliest:    true,
		},
		{
			name:         "request already fulfilled",
			requests:     []blockhashstore.Event{{Block: 150, ID: "request"}},
			fulfillments: []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:         256,
			lookback:     500,
			latest:       450,
		},
		{
			name:                "fulfillment no matching request no error",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request1"}},
			fulfillments:        []blockhashstore.Event{{Block: 153, ID: "request2"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{150, 151, 152, 153, 154, 155},
			getBatchSize:        1,
			storeBatchSize:      1,
			getBatchCallCount:   5,
			storeBatchCallCount: 5,
		},
		{
			name:                "error checking if stored, store subsequent blocks",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request1"}, {Block: 151, ID: "request2"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			bhs:                 blockhashstore.TestBHS{ErrorsIsStored: []uint64{150}},
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{151, 152, 153, 154, 155},
			getBatchSize:        1,
			storeBatchSize:      1,
			getBatchCallCount:   4,
			storeBatchCallCount: 4,
		},
		{
			name:                "another error checking if stored, store subsequent blocks",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request1"}, {Block: 151, ID: "request2"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			bhs:                 blockhashstore.TestBHS{ErrorsIsStored: []uint64{151}},
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{150, 151, 152, 153, 154, 155},
			getBatchSize:        1,
			storeBatchSize:      1,
			getBatchCallCount:   5,
			storeBatchCallCount: 5,
		},
		{
			name:              "error checking getBlockhashes, return with error",
			requests:          []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:              256,
			lookback:          500,
			latest:            450,
			alreadyStored:     []uint64{155},
			expectedStored:    []uint64{155},
			getBatchSize:      1,
			getBatchCallCount: 1,
			batchBHS:          blockhashstore.TestBatchBHS{GetBlockhashesError: errors.New("internal failure")},
			expectedErrMsg:    "finding earliest blocknumber with blockhash: fetching blockhashes: internal failure",
		},
		{
			name:                "error while storing block headers, return with error",
			requests:            []blockhashstore.Event{{Block: 150, ID: "request"}},
			wait:                256,
			lookback:            500,
			latest:              450,
			alreadyStored:       []uint64{155},
			expectedStored:      []uint64{155},
			getBatchSize:        1,
			storeBatchSize:      1,
			getBatchCallCount:   5,
			storeBatchCallCount: 1,
			batchBHS:            blockhashstore.TestBatchBHS{StoreVerifyHeadersError: errors.New("invalid header")},
			expectedErrMsg:      "store block headers: invalid header",
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.testFeeder)
	}
}

func (test testCase) testFeeder(t *testing.T) {
	lggr := logger.TestLogger(t)
	lggr.Debugf("running test case: %s", test.name)
	coordinator := &blockhashstore.TestCoordinator{
		RequestEvents:     test.requests,
		FulfillmentEvents: test.fulfillments,
	}

	test.batchBHS.Stored = append(test.batchBHS.Stored, test.alreadyStored...)

	blockHeaderProvider := &blockhashstore.TestBlockHeaderProvider{}
	fromAddress := "0x469aA2CD13e037DC5236320783dCfd0e641c0559"
	fromAddresses := []types.EIP55Address{types.EIP55Address(fromAddress)}
	ks := keystoremocks.NewEth(t)
	ks.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, mock.Anything).Maybe().Return(common.HexToAddress(fromAddress), nil)

	feeder := NewBlockHeaderFeeder(
		lggr,
		coordinator,
		&test.bhs,
		&test.batchBHS,
		blockHeaderProvider,
		test.wait,
		test.lookback,
		func(ctx context.Context) (uint64, error) {
			return test.latest, nil
		},
		ks,
		test.getBatchSize,
		test.storeBatchSize,
		fromAddresses,
		testutils.FixtureChainID,
	)

	err := feeder.Run(testutils.Context(t))
	if test.expectedErrMsg == "" {
		require.NoError(t, err)
	} else {
		require.EqualError(t, err, test.expectedErrMsg)
	}

	require.ElementsMatch(t, test.expectedStored, test.batchBHS.Stored)
	require.Equal(t, test.storedEarliest, test.bhs.StoredEarliest)
	require.Equal(t, test.getBatchCallCount, test.batchBHS.GetBlockhashesCallCounter)
	require.Equal(t, test.storeBatchCallCount, test.batchBHS.StoreVerifyHeaderCallCounter)
}

func TestFeeder_CachesStoredBlocks(t *testing.T) {
	coordinator := &blockhashstore.TestCoordinator{
		RequestEvents: []blockhashstore.Event{{Block: 74, ID: "request"}},
	}

	bhs := &blockhashstore.TestBHS{}
	batchBHS := &blockhashstore.TestBatchBHS{Stored: []uint64{75}}
	blockHeaderProvider := &blockhashstore.TestBlockHeaderProvider{}
	fromAddress := "0x469aA2CD13e037DC5236320783dCfd0e641c0559"
	fromAddresses := []types.EIP55Address{types.EIP55Address(fromAddress)}
	ks := keystoremocks.NewEth(t)
	ks.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, mock.Anything).Maybe().Return(common.HexToAddress(fromAddress), nil)

	feeder := NewBlockHeaderFeeder(
		logger.TestLogger(t),
		coordinator,
		bhs,
		batchBHS,
		blockHeaderProvider,
		20,
		30,
		func(ctx context.Context) (uint64, error) {
			return 100, nil
		},
		ks,
		1,
		1,
		fromAddresses,
		testutils.FixtureChainID,
	)

	// Should store block 74. block 75 was already stored from above
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.ElementsMatch(t, []uint64{74, 75}, batchBHS.Stored)

	// Run the feeder at a later block
	// cache should not be pruned yet because from block is lower than the stored blocks
	// latest block = 101
	// lookback block = 30
	// stored blocks = [74, 75]
	// from block = 71
	feeder.latestBlock = func(ctx context.Context) (uint64, error) {
		return 101, nil
	}
	// remove stored blocks
	batchBHS.Stored = nil
	require.NoError(t, feeder.Run(testutils.Context(t)))
	// nothing should be stored because of the feeder cache
	require.Empty(t, batchBHS.Stored)

	// Remove stored blocks from batchBHS and try again
	// for blocks 74, 75, nothing should be stored
	// because nothing was pruned above
	feeder.coordinator = &blockhashstore.TestCoordinator{
		RequestEvents: []blockhashstore.Event{
			{Block: 74, ID: "request1"},
			{Block: 75, ID: "request2"},
		},
	}
	batchBHS.Stored = nil
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, batchBHS.Stored)

	// Run the feeder at a later block. this time, the feeder cache will be pruned
	feeder.latestBlock = func(ctx context.Context) (uint64, error) {
		return 200, nil
	}
	batchBHS.Stored = []uint64{175}
	feeder.coordinator = &blockhashstore.TestCoordinator{RequestEvents: []blockhashstore.Event{{Block: 174, ID: "request"}}}
	require.NoError(t, feeder.Run(testutils.Context(t)))
	// nothing should be stored in this run because the cache will be pruned at the end of the current iteration.
	// in the next run, cache should be empty
	require.ElementsMatch(t, []uint64{174, 175}, batchBHS.Stored)

	// Rewind latest block
	feeder.coordinator = &blockhashstore.TestCoordinator{RequestEvents: []blockhashstore.Event{{Block: 74, ID: "request"}}}
	feeder.latestBlock = func(ctx context.Context) (uint64, error) {
		return 100, nil
	}
	batchBHS.Stored = []uint64{75}
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.ElementsMatch(t, []uint64{74, 75}, batchBHS.Stored)
}
