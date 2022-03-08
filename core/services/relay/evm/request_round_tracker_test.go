package evm_test

import (
	"context"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/relay/evm"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	htmocks "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/mocks"
	logmocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	ocrmocks "github.com/smartcontractkit/chainlink/core/services/ocr2/mocks"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/testhelpers"
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
	db                  *ocrmocks.OCRContractTrackerDB
	lb                  *logmocks.Broadcaster
	hb                  *htmocks.HeadBroadcaster
	ec                  *evmmocks.Client
	requestRoundTracker *evm.RequestRoundTracker
	configTracker       *evm.ConfigTracker
}

func newContractTrackerUni(t *testing.T, opts ...interface{}) (uni contractTrackerUni) {
	var chain evmconfig.ChainScopedConfig
	var filterer *ocr2aggregator.OCR2AggregatorFilterer
	var contract *offchain_aggregator_wrapper.OffchainAggregator
	for _, opt := range opts {
		switch v := opt.(type) {
		case evmconfig.ChainScopedConfig:
			chain = v
		case *ocr2aggregator.OCR2AggregatorFilterer:
			filterer = v
		case *offchain_aggregator_wrapper.OffchainAggregator:
			contract = v
		default:
			t.Fatalf("unrecognised option type %T", v)
		}
	}
	if chain == nil {
		chain = evmtest.NewChainScopedConfig(t, configtest.NewTestGeneralConfig(t))
	}
	if filterer == nil {
		filterer = mustNewFilterer(t, testutils.NewAddress())
	}
	if contract == nil {
		contract = mustNewContract(t, testutils.NewAddress())
	}
	uni.db = new(ocrmocks.OCRContractTrackerDB)
	uni.lb = new(logmocks.Broadcaster)
	uni.hb = new(htmocks.HeadBroadcaster)
	uni.ec = new(evmmocks.Client)

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	uni.requestRoundTracker = evm.NewRequestRoundTracker(
		contract,
		filterer,
		uni.ec,
		uni.lb,
		42,
		lggr,
		db,
		uni.db,
		chain,
	)
	contractABI, err := abi.JSON(strings.NewReader(offchain_aggregator_wrapper.OffchainAggregatorABI))
	require.NoError(t, err)
	uni.configTracker = evm.NewConfigTracker(lggr, contractABI, uni.ec, contract.Address(), chain.ChainType(), uni.hb)

	t.Cleanup(func() {
		uni.db.AssertExpectations(t)
		uni.lb.AssertExpectations(t)
		uni.hb.AssertExpectations(t)
		uni.ec.AssertExpectations(t)
	})

	return uni
}

func Test_OCRContractTracker_HandleLog_OCRContractLatestRoundRequested(t *testing.T) {
	t.Parallel()

	fixtureLogAddress := gethCommon.HexToAddress("0x03bd0d5d39629423979f8a0e53dbce78c1791ebf")
	fixtureFilterer := mustNewFilterer(t, fixtureLogAddress)
	fixtureContract := mustNewContract(t, fixtureLogAddress)

	t.Run("does not update if contract address doesn't match", func(t *testing.T) {
		uni := newContractTrackerUni(t)
		logBroadcast := new(logmocks.Broadcast)

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		logBroadcast.On("RawLog").Return(rawLog)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		uni.requestRoundTracker.HandleLog(logBroadcast)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		logBroadcast.AssertExpectations(t)
		uni.db.AssertExpectations(t)
	})

	t.Run("does nothing if log has already been consumed", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)
		logBroadcast := new(logmocks.Broadcast)

		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(true, nil)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		uni.requestRoundTracker.HandleLog(logBroadcast)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		logBroadcast.AssertExpectations(t)
		uni.db.AssertExpectations(t)
	})

	t.Run("for new round requested log", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
		require.Equal(t, 0, int(round))
		require.Equal(t, 0, int(epoch))

		// Any round supercedes the 0 round

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		logBroadcast := new(logmocks.Broadcast)
		logBroadcast.On("RawLog").Return(rawLog)
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr ocr2aggregator.OCR2AggregatorRoundRequested) bool {
			return rr.Epoch == 1 && rr.Round == 1
		})).Return(nil)

		uni.requestRoundTracker.HandleLog(logBroadcast)

		uni.db.AssertExpectations(t)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 1, int(round))

		// Same round with higher epoch supercedes
		rawLog2 := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_9.json")
		logBroadcast2 := new(logmocks.Broadcast)
		logBroadcast2.On("RawLog").Return(rawLog2)
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr ocr2aggregator.OCR2AggregatorRoundRequested) bool {
			return rr.Epoch == 1 && rr.Round == 9
		})).Return(nil)

		uni.requestRoundTracker.HandleLog(logBroadcast2)

		uni.db.AssertExpectations(t)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 9, int(round))

		logBroadcast.AssertExpectations(t)

		// Same round with lower epoch is ignored
		uni.requestRoundTracker.HandleLog(logBroadcast)

		uni.db.AssertExpectations(t)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", configDigest.Hex())
		assert.Equal(t, 1, int(epoch))
		assert.Equal(t, 9, int(round))

		logBroadcast.AssertExpectations(t)

		// Higher epoch with lower round supercedes
		rawLog3 := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_2_1.json")
		rawLog3.Address = fixtureContract.Address()
		logBroadcast3 := new(logmocks.Broadcast)
		logBroadcast3.On("RawLog").Return(rawLog3)
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.MatchedBy(func(rr ocr2aggregator.OCR2AggregatorRoundRequested) bool {
			return rr.Epoch == 2 && rr.Round == 1
		})).Return(nil)

		uni.requestRoundTracker.HandleLog(logBroadcast3)

		uni.db.AssertExpectations(t)

		configDigest, epoch, round, err = uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", configDigest.Hex())
		assert.Equal(t, 2, int(epoch))
		assert.Equal(t, 1, int(round))

		logBroadcast.AssertExpectations(t)
		uni.db.AssertExpectations(t)
	})

	t.Run("does not mark consumed or update state if latest round fails to save", func(t *testing.T) {
		uni := newContractTrackerUni(t, fixtureFilterer, fixtureContract)

		rawLog := cltest.LogFromFixture(t, "../../../testdata/jsonrpc/ocr2_round_requested_log_1_1.json")
		rawLog.Address = fixtureContract.Address()
		logBroadcast := new(logmocks.Broadcast)
		logBroadcast.On("RawLog").Return(rawLog)
		uni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		uni.db.On("SaveLatestRoundRequested", mock.Anything, mock.Anything).Return(errors.New("something exploded"))

		uni.requestRoundTracker.HandleLog(logBroadcast)

		uni.db.AssertExpectations(t)

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
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

		uni.db.On("LoadLatestRoundRequested").Return(rr, nil)

		require.NoError(t, uni.requestRoundTracker.Start())

		configDigest, epoch, round, err := uni.requestRoundTracker.LatestRoundRequested(context.Background(), 0)
		require.NoError(t, err)
		assert.Equal(t, (ocrtypes.ConfigDigest)(rr.ConfigDigest).Hex(), configDigest.Hex())
		assert.Equal(t, rr.Epoch, epoch)
		assert.Equal(t, rr.Round, round)

		uni.db.AssertExpectations(t)
		uni.lb.AssertExpectations(t)
		uni.hb.AssertExpectations(t)

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
