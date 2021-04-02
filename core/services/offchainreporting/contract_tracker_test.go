package offchainreporting_test

import (
	"context"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	ocrmocks "github.com/smartcontractkit/chainlink/core/services/offchainreporting/mocks"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mustNewContract(t *testing.T, address gethCommon.Address) *offchain_aggregator_wrapper.OffchainAggregator {
	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(address, nil)
	require.NoError(t, err)
	return contract
}

func Test_OCRContractTracker_HandleLog_OCRContractLatestRoundRequested(t *testing.T) {
	t.Parallel()

	fixtureLogAddress := gethCommon.HexToAddress("0x03bd0d5d39629423979f8a0e53dbce78c1791ebf")
	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(fixtureLogAddress, nil)
	require.NoError(t, err)

	t.Run("does not update if contract address doesn't match", func(t *testing.T) {
		db := new(ocrmocks.OCRContractTrackerDB)
		tracker, err := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
			db,
		)
		require.NoError(t, err)
		logBroadcast := new(mocks.Broadcast)

		rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log_1_1.json")
		logBroadcast.On("RawLog").Return(rawLog)
		logBroadcast.On("MarkConsumed").Return(nil)
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		tracker.HandleLog(logBroadcast)

		configDigest, epoch, round, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		logBroadcast.AssertExpectations(t)
		db.AssertExpectations(t)
	})

	t.Run("does nothing if log has already been consumed", func(t *testing.T) {
		db := new(ocrmocks.OCRContractTrackerDB)
		tracker, err := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
			db,
		)
		require.NoError(t, err)
		logBroadcast := new(mocks.Broadcast)

		logBroadcast.On("WasAlreadyConsumed").Return(true, nil)

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		tracker.HandleLog(logBroadcast)

		configDigest, epoch, round, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		logBroadcast.AssertExpectations(t)
		db.AssertExpectations(t)
	})

	t.Run("for new round requested log", func(t *testing.T) {
		db := new(ocrmocks.OCRContractTrackerDB)
		contract := mustNewContract(t, fixtureLogAddress)

		tracker, err := offchainreporting.NewOCRContractTracker(
			contract,
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
			db,
		)
		require.NoError(t, err)

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		// Any round supercedes the 0 round

		rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log_1_1.json")
		logBroadcast := new(mocks.Broadcast)
		logBroadcast.On("RawLog").Return(rawLog)
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)
		logBroadcast.On("MarkConsumed").Return(nil)

		db.On("SaveLatestRoundRequested", mock.MatchedBy(func(rr offchainaggregator.OffchainAggregatorRoundRequested) bool {
			return rr.Epoch == 1 && rr.Round == 1
		})).Return(nil)

		tracker.HandleLog(logBroadcast)

		db.AssertExpectations(t)

		configDigest, epoch, round, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 1, int(round))

		// Same round with higher epoch supercedes
		rawLog2 := cltest.LogFromFixture(t, "./testdata/round_requested_log_1_9.json")
		logBroadcast2 := new(mocks.Broadcast)
		logBroadcast2.On("RawLog").Return(rawLog2)
		logBroadcast2.On("WasAlreadyConsumed").Return(false, nil)
		logBroadcast2.On("MarkConsumed").Return(nil)

		db.On("SaveLatestRoundRequested", mock.MatchedBy(func(rr offchainaggregator.OffchainAggregatorRoundRequested) bool {
			return rr.Epoch == 1 && rr.Round == 9
		})).Return(nil)

		tracker.HandleLog(logBroadcast2)

		db.AssertExpectations(t)

		configDigest, epoch, round, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 9, int(round))

		logBroadcast.AssertExpectations(t)

		// Same round with lower epoch is ignored
		tracker.HandleLog(logBroadcast)

		db.AssertExpectations(t)

		configDigest, epoch, round, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 9, int(round))

		logBroadcast.AssertExpectations(t)

		// Higher epoch with lower round supercedes
		rawLog3 := cltest.LogFromFixture(t, "./testdata/round_requested_log_2_1.json")
		logBroadcast3 := new(mocks.Broadcast)
		logBroadcast3.On("RawLog").Return(rawLog3)
		logBroadcast3.On("WasAlreadyConsumed").Return(false, nil)
		logBroadcast3.On("MarkConsumed").Return(nil)

		db.On("SaveLatestRoundRequested", mock.MatchedBy(func(rr offchainaggregator.OffchainAggregatorRoundRequested) bool {
			return rr.Epoch == 2 && rr.Round == 1
		})).Return(nil)

		tracker.HandleLog(logBroadcast3)

		db.AssertExpectations(t)

		configDigest, epoch, round, err = tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "cccccccccccccccccccccccccccccccc", configDigest.Hex())
		assert.Equal(t, 2, int(epoch))
		assert.Equal(t, 1, int(round))

		logBroadcast.AssertExpectations(t)
		db.AssertExpectations(t)
	})

	t.Run("does mark consumed or update state if latest round fails to save", func(t *testing.T) {
		db := new(ocrmocks.OCRContractTrackerDB)
		contract := mustNewContract(t, fixtureLogAddress)

		tracker, err := offchainreporting.NewOCRContractTracker(
			contract,
			contractFilterer,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
			db,
		)
		require.NoError(t, err)

		rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log_1_1.json")
		logBroadcast := new(mocks.Broadcast)
		logBroadcast.On("RawLog").Return(rawLog)
		logBroadcast.On("WasAlreadyConsumed").Return(false, nil)

		db.On("SaveLatestRoundRequested", mock.Anything).Return(errors.New("something exploded"))

		tracker.HandleLog(logBroadcast)

		db.AssertExpectations(t)

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))
	})

	t.Run("restores latest round requested from database on start", func(t *testing.T) {
		db := new(ocrmocks.OCRContractTrackerDB)
		broadcaster := new(mocks.Broadcaster)
		contract := mustNewContract(t, fixtureLogAddress)
		tracker, err := offchainreporting.NewOCRContractTracker(
			contract,
			contractFilterer,
			nil,
			nil,
			broadcaster,
			42,
			*logger.Default,
			db,
		)
		require.NoError(t, err)

		rawLog := cltest.LogFromFixture(t, "./testdata/round_requested_log_1_1.json")
		rr := offchainaggregator.OffchainAggregatorRoundRequested{
			Requester:    cltest.NewAddress(),
			ConfigDigest: cltest.MakeConfigDigest(t),
			Epoch:        42,
			Round:        9,
			Raw:          rawLog,
		}

		broadcaster.On("Register", tracker, mock.Anything).Return(func() {})
		broadcaster.On("IsConnected").Return(true).Maybe()

		db.On("LoadLatestRoundRequested").Return(rr, nil)

		require.NoError(t, tracker.Start())

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, (ocrtypes.ConfigDigest)(rr.ConfigDigest).Hex(), configDigest.Hex())
		assert.Equal(t, rr.Epoch, epoch)
		assert.Equal(t, rr.Round, round)

		db.AssertExpectations(t)
		broadcaster.AssertExpectations(t)
	})
}

func Test_OCRContractTracker_IsLaterThan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		incoming types.Log
		existing types.Log
		expected bool
	}{
		{
			"incoming higher index than existing",
			types.Log{BlockNumber: 1, TxIndex: 1, Index: 2},
			types.Log{BlockNumber: 1, TxIndex: 1, Index: 1},
			true,
		},
		{
			"incoming lower index than existing",
			types.Log{BlockNumber: 1, TxIndex: 1, Index: 1},
			types.Log{BlockNumber: 1, TxIndex: 1, Index: 2},
			false,
		},
		{
			"incoming identical to existing",
			types.Log{BlockNumber: 1, TxIndex: 2, Index: 2},
			types.Log{BlockNumber: 1, TxIndex: 2, Index: 2},
			false,
		},
		{
			"incoming higher tx index than existing",
			types.Log{BlockNumber: 1, TxIndex: 2, Index: 2},
			types.Log{BlockNumber: 1, TxIndex: 1, Index: 2},
			true,
		},
		{
			"incoming lower tx index than existing",
			types.Log{BlockNumber: 1, TxIndex: 1, Index: 2},
			types.Log{BlockNumber: 1, TxIndex: 2, Index: 2},
			false,
		},
		{
			"incoming higher block number than existing",
			types.Log{BlockNumber: 3, TxIndex: 2, Index: 2},
			types.Log{BlockNumber: 2, TxIndex: 2, Index: 2},
			true,
		},
		{
			"incoming lower block number than existing",
			types.Log{BlockNumber: 2, TxIndex: 2, Index: 2},
			types.Log{BlockNumber: 3, TxIndex: 2, Index: 2},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := offchainreporting.IsLaterThan(test.incoming, test.existing)
			assert.Equal(t, test.expected, res)
		})
	}
}
