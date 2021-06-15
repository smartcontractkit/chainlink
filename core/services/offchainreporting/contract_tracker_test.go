package offchainreporting_test

import (
	"context"
	"math/big"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	ocrmocks "github.com/smartcontractkit/chainlink/core/services/offchainreporting/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

func Test_OCRContractTracker_LatestBlockHeight(t *testing.T) {
	t.Parallel()

	t.Run("on L2 chains, always returns 0", func(t *testing.T) {
		tracker := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			nil,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
			nil,
			nil,
			chains.OptimismMainnet,
			nil,
		)

		l, err := tracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(0), l)
	})

	t.Run("before first head incoming, looks up on-chain", func(t *testing.T) {
		ethClient := new(mocks.Client)

		ethClient.On("HeaderByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(&models.Head{Number: 42}, nil)

		tracker := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			nil,
			nil,
			ethClient,
			nil,
			42,
			*logger.Default,
			nil,
			nil,
			chains.EthMainnet,
			nil,
		)

		l, err := tracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)
	})

	t.Run("Before first head incoming, on client error returns error", func(t *testing.T) {
		ethClient := new(mocks.Client)

		tracker := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			nil,
			nil,
			ethClient,
			nil,
			42,
			*logger.Default,
			nil,
			nil,
			chains.EthMainnet,
			nil,
		)

		ethClient.On("HeaderByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(nil, nil).Once()

		_, err := tracker.LatestBlockHeight(context.Background())
		assert.EqualError(t, err, "got nil head")

		ethClient.On("HeaderByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(nil, errors.New("bar")).Once()

		_, err = tracker.LatestBlockHeight(context.Background())
		assert.EqualError(t, err, "bar")

		ethClient.AssertExpectations(t)
	})

	t.Run("after first head incoming, uses cached value", func(t *testing.T) {
		tracker := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			nil,
			nil,
			nil,
			nil,
			42,
			*logger.Default,
			nil,
			nil,
			chains.EthMainnet,
			nil,
		)

		tracker.OnNewLongestChain(context.Background(), models.Head{Number: 42})

		l, err := tracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)
	})
}

func Test_OCRContractTracker_HandleLog_OCRContractLatestRoundRequested(t *testing.T) {
	t.Parallel()

	fixtureLogAddress := gethCommon.HexToAddress("0x03bd0d5d39629423979f8a0e53dbce78c1791ebf")
	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(fixtureLogAddress, nil)
	require.NoError(t, err)
	s, c := cltest.NewStore(t)
	defer c()

	t.Run("does not update if contract address doesn't match", func(t *testing.T) {
		db := new(ocrmocks.OCRContractTrackerDB)
		lb := new(logmocks.Broadcaster)
		tracker := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			contractFilterer,
			nil,
			nil,
			lb,
			42,
			*logger.Default,
			s.DB,
			db,
			nil,
			nil,
		)
		require.NoError(t, err)
		logBroadcast := new(logmocks.Broadcast)

		rawLog := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_1.json")
		logBroadcast.On("RawLog").Return(rawLog)
		lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
		lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

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
		lb := new(logmocks.Broadcaster)
		tracker := offchainreporting.NewOCRContractTracker(
			mustNewContract(t, cltest.NewAddress()),
			contractFilterer,
			nil,
			nil,
			lb,
			42,
			*logger.Default,
			s.DB,
			db,
			nil,
			nil,
		)
		require.NoError(t, err)
		logBroadcast := new(logmocks.Broadcast)

		lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(true, nil)

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
		lb := new(logmocks.Broadcaster)
		tracker := offchainreporting.NewOCRContractTracker(
			contract,
			contractFilterer,
			nil,
			nil,
			lb,
			42,
			*logger.Default,
			s.DB,
			db,
			nil,
			nil,
		)
		require.NoError(t, err)

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		// Any round supercedes the 0 round

		rawLog := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_1.json")
		logBroadcast := new(logmocks.Broadcast)
		logBroadcast.On("RawLog").Return(rawLog)
		lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr offchainaggregator.OffchainAggregatorRoundRequested) bool {
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
		rawLog2 := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_9.json")
		logBroadcast2 := new(logmocks.Broadcast)
		logBroadcast2.On("RawLog").Return(rawLog2)
		lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr offchainaggregator.OffchainAggregatorRoundRequested) bool {
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
		rawLog3 := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_2_1.json")
		logBroadcast3 := new(logmocks.Broadcast)
		logBroadcast3.On("RawLog").Return(rawLog3)
		lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr offchainaggregator.OffchainAggregatorRoundRequested) bool {
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
		lb := new(logmocks.Broadcaster)
		tracker := offchainreporting.NewOCRContractTracker(
			contract,
			contractFilterer,
			nil,
			nil,
			lb,
			42,
			*logger.Default,
			s.DB,
			db,
			nil,
			nil,
		)
		require.NoError(t, err)

		rawLog := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_1.json")
		logBroadcast := new(logmocks.Broadcast)
		logBroadcast.On("RawLog").Return(rawLog)
		lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		db.On("SaveLatestRoundRequested", mock.Anything, mock.Anything).Return(errors.New("something exploded"))

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
		broadcaster := new(logmocks.Broadcaster)
		contract := mustNewContract(t, fixtureLogAddress)
		hb := new(htmocks.HeadBroadcaster)
		tracker := offchainreporting.NewOCRContractTracker(
			contract,
			contractFilterer,
			nil,
			nil,
			broadcaster,
			42,
			*logger.Default,
			s.DB,
			db,
			nil,
			hb,
		)
		require.NoError(t, err)

		rawLog := cltest.LogFromFixture(t, "../../testdata/jsonrpc/round_requested_log_1_1.json")
		rr := offchainaggregator.OffchainAggregatorRoundRequested{
			Requester:    cltest.NewAddress(),
			ConfigDigest: cltest.MakeConfigDigest(t),
			Epoch:        42,
			Round:        9,
			Raw:          rawLog,
		}

		eventuallyCloseLogBroadcaster := cltest.NewAwaiter()
		broadcaster.On("Register", tracker, mock.Anything).Return(func() { eventuallyCloseLogBroadcaster.ItHappened() })
		broadcaster.On("IsConnected").Return(true).Maybe()

		eventuallyCloseHeadBroadcaster := cltest.NewAwaiter()
		hb.On("Subscribe", tracker).Return(func() { eventuallyCloseHeadBroadcaster.ItHappened() })

		db.On("LoadLatestRoundRequested").Return(rr, nil)

		require.NoError(t, tracker.Start())

		configDigest, epoch, round, err := tracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, (ocrtypes.ConfigDigest)(rr.ConfigDigest).Hex(), configDigest.Hex())
		assert.Equal(t, rr.Epoch, epoch)
		assert.Equal(t, rr.Round, round)

		db.AssertExpectations(t)
		broadcaster.AssertExpectations(t)
		hb.AssertExpectations(t)

		require.NoError(t, tracker.Close())

		eventuallyCloseHeadBroadcaster.AssertHappened(t)
		eventuallyCloseLogBroadcaster.AssertHappened(t)
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
