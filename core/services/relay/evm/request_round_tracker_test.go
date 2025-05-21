package evm_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/v2/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mocks"
)

func mustNewContract(t *testing.T, address gethCommon.Address) *offchain_aggregator_wrapper.OffchainAggregator {
	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(address, nil)
	require.NoError(t, err)
	return contract
}

func mustNewFilterer(t *testing.T, address gethCommon.Address) *ocr2aggregator.OCR2AggregatorFilterer {
	filterer, err := ocr2aggregator.NewOCR2AggregatorFilterer(testutils.NewAddress(), nil)
	require.NoError(t, err)
	return filterer
}

type contractTrackerUni struct {
	db                  *mocks.RequestRoundDB
	lb                  *logmocks.Broadcaster
	hb                  *headstest.Broadcaster[*evmtypes.Head, common.Hash]
	ec                  *clienttest.Client
	requestRoundTracker *evm.RequestRoundTracker
}

func newContractTrackerUni(t *testing.T, opts ...interface{}) (uni contractTrackerUni) {
	var filterer *ocr2aggregator.OCR2AggregatorFilterer
	var contract *offchain_aggregator_wrapper.OffchainAggregator
	for _, opt := range opts {
		switch v := opt.(type) {
		case *ocr2aggregator.OCR2AggregatorFilterer:
			filterer = v
		case *offchain_aggregator_wrapper.OffchainAggregator:
			contract = v
		default:
			t.Fatalf("unrecognised option type %T", v)
		}
	}
	config := configtest.NewTestGeneralConfig(t)
	chain := evmtest.NewChainScopedConfig(t, config)
	if filterer == nil {
		filterer = mustNewFilterer(t, testutils.NewAddress())
	}
	if contract == nil {
		contract = mustNewContract(t, testutils.NewAddress())
	}
	uni.db = mocks.NewRequestRoundDB(t)
	uni.lb = logmocks.NewBroadcaster(t)
	uni.hb = headstest.NewBroadcaster[*evmtypes.Head, common.Hash](t)
	uni.ec = clienttest.NewClient(t)

	db := pgtest.NewSqlxDB(t)
	lggr := logger.Test(t)
	uni.requestRoundTracker = evm.NewRequestRoundTracker(
		contract,
		filterer,
		uni.ec,
		uni.lb,
		42,
		lggr,
		db,
		uni.db,
		chain.EVM(),
	)

	return uni
}

func Test_OCRContractTracker_HandleLog_OCRContractLatestRoundRequested(t *testing.T) {
	t.Parallel()

	fixtureLogAddress := gethCommon.HexToAddress("0x03bd0d5d39629423979f8a0e53dbce78c1791ebf")
	fixtureFilterer := mustNewFilterer(t, fixtureLogAddress)
	fixtureContract := mustNewContract(t, fixtureLogAddress)

	t.Run("does not update if contract address doesn't match", func(t *testing.T) {
		uni := newContractTrackerUni(t)
		logBroadcast := logmocks.NewBroadcast(t)

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		logBroadcast.On("RawLog").Return(rawLog).Maybe()
		logBroadcast.On("String").Return("").Maybe()
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))
	})

	t.Run("does nothing if log has already been consumed", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)
		logBroadcast := logmocks.NewBroadcast(t)
		logBroadcast.On("String").Return("").Maybe()

		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(true, nil)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))
	})

	t.Run("for new round requested log", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		// Any round supercedes the 0 round

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		logBroadcast := logmocks.NewBroadcast(t)
		logBroadcast.On("RawLog").Return(rawLog).Maybe()
		logBroadcast.On("String").Return("").Maybe()
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr ocr2aggregator.OCR2AggregatorRoundRequested) bool {
			return rr.Epoch == 1 && rr.Round == 1
		})).Return(nil)
		uni.db.On("WithDataSource", mock.Anything).Return(uni.db)

		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		assert.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 1, int(round))

		// Same round with higher epoch supercedes
		rawLog2 := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_9.json")
		logBroadcast2 := logmocks.NewBroadcast(t)
		logBroadcast2.On("RawLog").Return(rawLog2).Maybe()
		logBroadcast2.On("String").Return("").Maybe()
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr ocr2aggregator.OCR2AggregatorRoundRequested) bool {
			return rr.Epoch == 1 && rr.Round == 9
		})).Return(nil)

		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast2)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 9, int(round))

		// Same round with lower epoch is ignored
		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 9, int(round))

		// Higher epoch with lower round supercedes
		rawLog3 := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_2_1.json")
		rawLog3.Address = fixtureContract.Address()
		logBroadcast3 := logmocks.NewBroadcast(t)
		logBroadcast3.On("RawLog").Return(rawLog3).Maybe()
		logBroadcast3.On("String").Return("").Maybe()
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr ocr2aggregator.OCR2AggregatorRoundRequested) bool {
			return rr.Epoch == 2 && rr.Round == 1
		})).Return(nil)

		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast3)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		assert.Equal(t, "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", configDigest.Hex())
		assert.Equal(t, 2, int(epoch))
		assert.Equal(t, 1, int(round))
	})

	t.Run("does not mark consumed or update state if latest round fails to save", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		rawLog.Address = fixtureContract.Address()
		logBroadcast := logmocks.NewBroadcast(t)
		logBroadcast.On("RawLog").Return(rawLog).Maybe()
		logBroadcast.On("String").Return("").Maybe()
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.Anything).Return(errors.New("something exploded"))
		uni.db.On("WithDataSource", mock.Anything).Return(uni.db)

		uni.requestRoundTracker.HandleLog(t.Context(), logBroadcast)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))
	})

	t.Run("restores latest round requested from database on start", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		rr := ocr2aggregator.OCR2AggregatorRoundRequested{
			Requester:    testutils.NewAddress(),
			ConfigDigest: testhelpers.MakeConfigDigest(t),
			Epoch:        42,
			Round:        9,
			Raw:          rawLog,
		}

		eventuallyCloseLogBroadcaster := cltest.NewAwaiter()
		uni.lb.On("Register", uni.requestRoundTracker, mock.Anything).Return(func() { eventuallyCloseLogBroadcaster.ItHappened() })
		uni.lb.On("IsConnected").Return(true).Maybe()

		uni.db.On("LoadLatestRoundRequested", mock.Anything).Return(rr, nil)

		ctx := testutils.Context(t)
		require.NoError(t, uni.requestRoundTracker.Start(ctx))

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(testutils.Context(t), 0)
		require.NoError(t, err)
		assert.Equal(t, (ocrtypes.ConfigDigest)(rr.ConfigDigest).Hex(), configDigest.Hex())
		assert.Equal(t, rr.Epoch, epoch)
		assert.Equal(t, rr.Round, round)

		require.NoError(t, uni.requestRoundTracker.Close())

		eventuallyCloseLogBroadcaster.AssertHappened(t, true)
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
			res := evm.IsLaterThan(test.incoming, test.existing)
			assert.Equal(t, test.expected, res)
		})
	}
}
