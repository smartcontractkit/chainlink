package blockhashstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	_ Coordinator = &TestCoordinator{}
	_ BHS         = &TestBHS{}
)

func TestFeeder(t *testing.T) {
	tests := []struct {
		name           string
		requests       []Event
		fulfillments   []Event
		wait           int
		lookback       int
		latest         uint64
		bhs            TestBHS
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
			bhs:            TestBHS{Stored: []uint64{150}},
			expectedStored: []uint64{150},
		},
		{
			name:           "error checking if stored, store anyway",
			requests:       []Event{{Block: 150, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            TestBHS{ErrorsIsStored: []uint64{150}},
			expectedStored: []uint64{150},
			expectedErrMsg: "checking if stored: error checking if stored",
		},
		{
			name:           "error storing, continue to next block anyway",
			requests:       []Event{{Block: 150, ID: "request"}, {Block: 151, ID: "request"}},
			wait:           25,
			lookback:       100,
			latest:         200,
			bhs:            TestBHS{ErrorsStore: []uint64{150}},
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
			coordinator := &TestCoordinator{
				RequestEvents:     test.requests,
				FulfillmentEvents: test.fulfillments,
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

			require.ElementsMatch(t, test.expectedStored, test.bhs.Stored)
		})
	}
}

func TestFeeder_CachesStoredBlocks(t *testing.T) {
	coordinator := &TestCoordinator{
		RequestEvents: []Event{{Block: 100, ID: "request"}},
	}

	bhs := &TestBHS{}

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
	require.ElementsMatch(t, []uint64{100}, bhs.Stored)

	// Remove 100 from the BHS and try again, it should not be stored since it's cached in the
	// feeder
	bhs.Stored = nil
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, bhs.Stored)

	// Run the feeder on a later block and make sure the cache is pruned
	feeder.latestBlock = func(ctx context.Context) (uint64, error) {
		return 500, nil
	}
	require.NoError(t, feeder.Run(testutils.Context(t)))
	require.Empty(t, feeder.stored)
}
