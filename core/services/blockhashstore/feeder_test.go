package blockhashstore

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var (
	_ Coordinator = &testCoordinator{}
	_ BHS         = &testBHS{}
)

func TestFeeder(t *testing.T) {
	tests := []struct {
		name           string
		requests       []Event
		fulfillments   []Event
		wait           int
		lookback       int
		latest         uint64
		bhs            testBHS
		expectedStored []uint64
		expectedErrMsg string
	}{
		{
			name:           "single unfulfilled request",
			requests:       []Event{{Block: 150, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150},
		},
		{
			name:           "single fulfilled request",
			requests:       []Event{{Block: 150, ID: "request"}},
			fulfillments:   []Event{{Block: 155, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "single already fulfilled",
			requests:       []Event{{Block: 150, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            testBHS{stored: []uint64{150}},
			expectedStored: []uint64{150},
		},
		{
			name:           "error checking if stored, store anyway",
			requests:       []Event{{Block: 150, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            testBHS{errorsIsStored: []uint64{150}},
			expectedStored: []uint64{150},
			expectedErrMsg: "checking if stored: error checking if stored",
		},
		{
			name:           "error storing, continue to next block anyway",
			requests:       []Event{{Block: 150, ID: "request"}, {Block: 151, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            testBHS{errorsStore: []uint64{150}},
			expectedStored: []uint64{151},
			expectedErrMsg: "storing block: error storing",
		},
		{
			name: "multiple requests same block, some fulfilled",
			requests: []Event{
				{Block: 150, ID: "request1"},
				{Block: 150, ID: "request2"},
				{Block: 150, ID: "request3"}},
			fulfillments: []Event{
				{Block: 150, ID: "request1"},
				{Block: 150, ID: "request3"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150},
		},
		{
			name: "multiple requests same block, all fulfilled",
			requests: []Event{
				{Block: 150, ID: "request1"},
				{Block: 150, ID: "request2"},
				{Block: 150, ID: "request3"}},
			fulfillments: []Event{
				{Block: 150, ID: "request1"},
				{Block: 150, ID: "request2"},
				{Block: 150, ID: "request3"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "fulfillment no matching request no error",
			requests:       []Event{{Block: 150, ID: "request"}},
			fulfillments:   []Event{{Block: 199, ID: "request2"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150},
		},
		{
			name:           "multiple unfulfilled requests",
			requests:       []Event{{Block: 150, ID: "request1"}, {Block: 151, ID: "request2"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150, 151},
		},
		{
			name:           "multiple fulfilled requests",
			requests:       []Event{{Block: 150, ID: "request1"}, {Block: 151, ID: "request2"}},
			fulfillments:   []Event{{Block: 150, ID: "request1"}, {Block: 151, ID: "request2"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "recent unfulfilled request do not store",
			requests:       []Event{{Block: 185, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name:           "old unfulfilled request do not store",
			requests:       []Event{{Block: 99, ID: "request"}, {Block: 57, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{},
		},
		{
			name: "mixed",
			requests: []Event{
				// Block 150
				{Block: 150, ID: "request1"},
				{Block: 150, ID: "request2"},
				{Block: 150, ID: "request3"},

				// Block 151
				{Block: 151, ID: "request4"},
				{Block: 151, ID: "request5"},

				// Block 153
				{Block: 153, ID: "request6"},

				// Block 154
				{Block: 154, ID: "request7"}},
			fulfillments: []Event{
				// Block 150
				{Block: 150, ID: "request1"},
				{Block: 150, ID: "request2"},
				// request3 no fulfillment

				// Block 151
				{Block: 151, ID: "request4"},
				{Block: 151, ID: "request5"},

				// Block 153 - no fulfillment

				// Block 154
				{Block: 154, ID: "request7"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			expectedStored: []uint64{150, 153},
		},
		{
			name:           "lookback before 0th block",
			requests:       []Event{{Block: 20, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         50,
			expectedStored: []uint64{20},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			coordinator := &testCoordinator{
				requests:     test.requests,
				fulfillments: test.fulfillments,
			}

			feeder := NewFeeder(
				logger.TestLogger(t),
				coordinator,
				&test.bhs,
				test.wait,
				test.lookback,
				func(ctx context.Context) (uint64, error) {
					return test.latest, nil
				})

			err := feeder.Run(testutils.Context(t))
			if test.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expectedErrMsg)
			}

			require.ElementsMatch(t, test.expectedStored, test.bhs.stored)
		})
	}
}

func TestFeeder_CachesStoredBlocks(t *testing.T) {
	coordinator := &testCoordinator{
		requests: []Event{{Block: 100, ID: "request"}},
	}

	bhs := &testBHS{}

	feeder := NewFeeder(
		logger.TestLogger(t),
		coordinator,
		bhs,
		100,
		200,
		func(ctx context.Context) (uint64, error) {
			return 250, nil
		})

	// Should store block 100
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.ElementsMatch(t, []uint64{100}, bhs.stored)

	// Remove 100 from the BHS and try again, it should not be stored since it's cached in the
	// feeder
	bhs.stored = nil
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, bhs.stored)

	// Run the feeder on a later block and make sure the cache is pruned
	feeder.latestBlock = func(ctx context.Context) (uint64, error) {
		return 500, nil
	}
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, feeder.stored)
}

type testCoordinator struct {
	requests     []Event
	fulfillments []Event
}

func (t *testCoordinator) Requests(_ context.Context, fromBlock uint64, toBlock uint64) ([]Event, error) {
	var result []Event
	for _, req := range t.requests {
		if req.Block >= fromBlock && req.Block <= toBlock {
			result = append(result, req)
		}
	}
	return result, nil
}

func (t *testCoordinator) Fulfillments(_ context.Context, fromBlock uint64) ([]Event, error) {
	var result []Event
	for _, ful := range t.fulfillments {
		if ful.Block >= fromBlock {
			result = append(result, ful)
		}
	}
	return result, nil
}

type testBHS struct {
	stored []uint64

	// errorsStore defines which block numbers should return errors on Store.
	errorsStore []uint64

	// errorsIsStored defines which block numbers should return errors on IsStored.
	errorsIsStored []uint64
}

func (t *testBHS) Store(_ context.Context, blockNum uint64) error {
	for _, e := range t.errorsStore {
		if e == blockNum {
			return errors.New("error storing")
		}
	}

	t.stored = append(t.stored, blockNum)
	return nil
}

func (t *testBHS) IsStored(_ context.Context, blockNum uint64) (bool, error) {
	for _, e := range t.errorsIsStored {
		if e == blockNum {
			return false, errors.New("error checking if stored")
		}
	}

	for _, s := range t.stored {
		if s == blockNum {
			return true, nil
		}
	}
	return false, nil
}
