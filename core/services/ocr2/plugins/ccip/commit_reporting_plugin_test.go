package ccip

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestCommitReportingPlugin_Observation(t *testing.T) {
	sourceNativeTokenAddr := common.HexToAddress("1000")
	someTokenAddr := common.HexToAddress("2000")

	testCases := []struct {
		name                string
		epochAndRound       types.ReportTimestamp
		commitStoreIsPaused bool
		commitStoreSeqNum   uint64
		tokenPrices         map[common.Address]*big.Int
		sendReqs            []ccipdata.Event[ccipdata.EVM2EVMMessage]
		tokenDecimals       map[common.Address]uint8
		fee                 *big.Int

		expErr bool
		expObs CommitObservation
	}{
		{
			name:              "base report",
			commitStoreSeqNum: 54,
			tokenPrices: map[common.Address]*big.Int{
				someTokenAddr:         big.NewInt(2),
				sourceNativeTokenAddr: big.NewInt(2),
			},
			sendReqs: []ccipdata.Event[ccipdata.EVM2EVMMessage]{
				{Data: ccipdata.EVM2EVMMessage{SequenceNumber: 54}},
				{Data: ccipdata.EVM2EVMMessage{SequenceNumber: 55}},
			},
			fee: big.NewInt(100),
			tokenDecimals: map[common.Address]uint8{
				someTokenAddr: 8,
			},
			expObs: CommitObservation{
				TokenPricesUSD: map[common.Address]*big.Int{
					someTokenAddr: big.NewInt(20000000000),
				},
				SourceGasPriceUSD: big.NewInt(0),
				Interval: commit_store.CommitStoreInterval{
					Min: 54,
					Max: 55,
				},
			},
		},
		{
			name:                "commit store is down",
			commitStoreIsPaused: true,
			expErr:              true,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceFinalityDepth := 10

			commitStore, _ := testhelpers.NewFakeCommitStore(t, tc.commitStoreSeqNum)
			commitStore.SetPaused(tc.commitStoreIsPaused)

			onRampReader := ccipdata.NewMockOnRampReader(t)
			if len(tc.sendReqs) > 0 {
				onRampReader.On("GetSendRequestsGteSeqNum", ctx, tc.commitStoreSeqNum, sourceFinalityDepth).
					Return(tc.sendReqs, nil)
			}

			tokenDecimalsCache := cache.NewMockAutoSync[map[common.Address]uint8](t)
			if len(tc.tokenDecimals) > 0 {
				tokenDecimalsCache.On("Get", ctx).Return(tc.tokenDecimals, nil)
			}

			priceGet := pricegetter.NewMockPriceGetter(t)
			if len(tc.tokenPrices) > 0 {
				addrs := []common.Address{sourceNativeTokenAddr}
				for addr := range tc.tokenDecimals {
					addrs = append(addrs, addr)
				}
				priceGet.On("TokenPricesUSD", mock.Anything, addrs).Return(tc.tokenPrices, nil)
			}

			sourceFeeEst := mocks.NewEvmFeeEstimator(t)
			if tc.fee != nil {
				sourceFeeEst.On("GetFee", ctx, []byte(nil), uint32(0), assets.NewWei(big.NewInt(0))).
					Return(gas.EvmFee{Legacy: assets.NewWei(tc.fee)}, uint32(0), nil)
			}

			p := &CommitReportingPlugin{}
			p.lggr = logger.TestLogger(t)
			p.inflightReports = newInflightCommitReportsContainer(time.Hour)
			p.config.commitStore = commitStore
			p.offchainConfig.SourceFinalityDepth = uint32(sourceFinalityDepth)
			p.config.onRampReader = onRampReader
			p.tokenDecimalsCache = tokenDecimalsCache
			p.config.priceGetter = priceGet
			p.config.sourceFeeEstimator = sourceFeeEst
			p.config.sourceNative = sourceNativeTokenAddr

			obs, err := p.Observation(ctx, tc.epochAndRound, types.Query{})

			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			expObsBytes, err := tc.expObs.Marshal()
			assert.NoError(t, err)
			assert.Equal(t, expObsBytes, []byte(obs))
		})
	}
}

func TestCommitReportingPlugin_Report(t *testing.T) {

	testCases := []struct {
		name              string
		observations      []CommitObservation
		f                 int
		gasPriceUpdates   []ccipdata.Event[price_registry.PriceRegistryUsdPerUnitGasUpdated]
		tokenPriceUpdates []ccipdata.Event[price_registry.PriceRegistryUsdPerTokenUpdated]
		sendRequests      []ccipdata.Event[ccipdata.EVM2EVMMessage]

		expCommitReport *commit_store.CommitStoreCommitReport
		expSeqNumRange  commit_store.CommitStoreInterval
		expErr          bool
	}{
		{
			name: "base",
			observations: []CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 1, Max: 1}},
				{Interval: commit_store.CommitStoreInterval{Min: 1, Max: 1}},
			},
			f: 1,
			sendRequests: []ccipdata.Event[ccipdata.EVM2EVMMessage]{
				{
					Data: ccipdata.EVM2EVMMessage{
						SequenceNumber: 1,
					},
				},
			},
			expSeqNumRange: commit_store.CommitStoreInterval{Min: 1, Max: 1},
			expCommitReport: &commit_store.CommitStoreCommitReport{
				MerkleRoot: [32]byte{},
				Interval:   commit_store.CommitStoreInterval{Min: 1, Max: 1},
				PriceUpdates: commit_store.InternalPriceUpdates{
					TokenPriceUpdates: nil,
					DestChainSelector: 0,
					UsdPerUnitGas:     big.NewInt(0),
				},
			},
			expErr: false,
		},
		{
			name: "not enough observations",
			observations: []CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 1, Max: 1}},
			},
			f:              1,
			sendRequests:   []ccipdata.Event[ccipdata.EVM2EVMMessage]{{}},
			expSeqNumRange: commit_store.CommitStoreInterval{Min: 1, Max: 1},
			expErr:         true,
		},
		{
			name: "empty",
			observations: []CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 0, Max: 0}},
				{Interval: commit_store.CommitStoreInterval{Min: 0, Max: 0}},
			},
			f:      1,
			expErr: false,
		},
		{
			name: "no leaves",
			observations: []CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 2, Max: 2}},
				{Interval: commit_store.CommitStoreInterval{Min: 2, Max: 2}},
			},
			f:              1,
			sendRequests:   []ccipdata.Event[ccipdata.EVM2EVMMessage]{{}},
			expSeqNumRange: commit_store.CommitStoreInterval{Min: 2, Max: 2},
			expErr:         true,
		},
	}

	ctx := testutils.Context(t)
	sourceChainSelector := rand.Int()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			destPriceRegistry, destPriceRegistryAddress := testhelpers.NewFakePriceRegistry(t)

			destReader := ccipdata.NewMockReader(t)
			destReader.On("GetGasPriceUpdatesCreatedAfter", ctx, destPriceRegistryAddress, uint64(sourceChainSelector), mock.Anything, 0).Return(tc.gasPriceUpdates, nil)
			destReader.On("GetTokenPriceUpdatesCreatedAfter", ctx, destPriceRegistryAddress, mock.Anything, 0).Return(tc.tokenPriceUpdates, nil)

			onRampReader := ccipdata.NewMockOnRampReader(t)
			if len(tc.sendRequests) > 0 {
				onRampReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.expSeqNumRange.Min, tc.expSeqNumRange.Max, 0).Return(tc.sendRequests, nil)
			}

			p := &CommitReportingPlugin{}
			p.lggr = logger.TestLogger(t)
			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			p.destPriceRegistry = destPriceRegistry
			p.config.destReader = destReader
			p.config.onRampReader = onRampReader
			p.config.sourceChainSelector = uint64(sourceChainSelector)

			aos := make([]types.AttributedObservation, 0, len(tc.observations))
			for _, o := range tc.observations {
				obs, err := o.Marshal()
				assert.NoError(t, err)
				aos = append(aos, types.AttributedObservation{Observation: obs})
			}

			gotSomeReport, gotReport, err := p.Report(ctx, types.ReportTimestamp{}, types.Query{}, aos)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tc.expCommitReport != nil {
				assert.True(t, gotSomeReport)
				encodedExpectedReport, err := abihelpers.EncodeCommitReport(*tc.expCommitReport)
				assert.NoError(t, err)
				assert.Equal(t, types.Report(encodedExpectedReport), gotReport)
			}
		})
	}
}

func TestCommitReportingPlugin_ShouldAcceptFinalizedReport(t *testing.T) {
	ctx := testutils.Context(t)

	newPlugin := func() *CommitReportingPlugin {
		p := &CommitReportingPlugin{}
		p.lggr = logger.TestLogger(t)
		p.inflightReports = newInflightCommitReportsContainer(time.Minute)
		return p
	}

	t.Run("report cannot be decoded leads to error", func(t *testing.T) {
		p := newPlugin()
		encodedReport := []byte("whatever")
		_, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.Error(t, err)
	})

	t.Run("empty report should not be accepted", func(t *testing.T) {
		p := newPlugin()
		report := commit_store.CommitStoreCommitReport{
			// UsdPerUnitGas is mandatory otherwise report cannot be encoded/decoded
			PriceUpdates: commit_store.InternalPriceUpdates{UsdPerUnitGas: big.NewInt(int64(rand.Int()))},
		}
		encodedReport, err := abihelpers.EncodeCommitReport(report)
		assert.NoError(t, err)
		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldAccept)
	})

	t.Run("stale report should not be accepted", func(t *testing.T) {
		onChainSeqNum := uint64(100)

		commitStore, _ := testhelpers.NewFakeCommitStore(t, onChainSeqNum)

		p := newPlugin()
		p.config.commitStore = commitStore

		report := commit_store.CommitStoreCommitReport{
			PriceUpdates: commit_store.InternalPriceUpdates{UsdPerUnitGas: big.NewInt(int64(rand.Int()))},
			MerkleRoot:   [32]byte{123}, // this report is considered non-empty since it has a merkle root
		}

		// stale since report interval is behind on chain seq num
		report.Interval = commit_store.CommitStoreInterval{Min: onChainSeqNum - 2, Max: onChainSeqNum + 10}
		encodedReport, err := abihelpers.EncodeCommitReport(report)
		assert.NoError(t, err)

		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldAccept)
	})

	t.Run("non-stale report should be accepted and added inflight", func(t *testing.T) {
		onChainSeqNum := uint64(100)

		commitStore, _ := testhelpers.NewFakeCommitStore(t, onChainSeqNum)

		p := newPlugin()
		p.config.commitStore = commitStore

		report := commit_store.CommitStoreCommitReport{
			PriceUpdates: commit_store.InternalPriceUpdates{
				TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{
					{
						SourceToken: utils.RandomAddress(),
						UsdPerToken: big.NewInt(int64(rand.Int())),
					},
				},
				DestChainSelector: rand.Uint64(),
				UsdPerUnitGas:     big.NewInt(int64(rand.Int())),
			},
			MerkleRoot: [32]byte{123},
		}

		// non-stale since report interval is not behind on-chain seq num
		report.Interval = commit_store.CommitStoreInterval{Min: onChainSeqNum, Max: onChainSeqNum + 10}
		encodedReport, err := abihelpers.EncodeCommitReport(report)
		assert.NoError(t, err)

		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.True(t, shouldAccept)

		// make sure that the report was added inflight
		tokenPriceUpdates := p.inflightReports.latestInflightTokenPriceUpdates()
		priceUpdate := tokenPriceUpdates[report.PriceUpdates.TokenPriceUpdates[0].SourceToken]
		assert.Equal(t, report.PriceUpdates.TokenPriceUpdates[0].UsdPerToken.Uint64(), priceUpdate.value.Uint64())
	})
}

func TestCommitReportingPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	report := commit_store.CommitStoreCommitReport{
		PriceUpdates: commit_store.InternalPriceUpdates{
			TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{
				{SourceToken: utils.RandomAddress(), UsdPerToken: big.NewInt(9e18)},
			},
			DestChainSelector: rand.Uint64(),
			UsdPerUnitGas:     big.NewInt(2000e9),
		},
		MerkleRoot: [32]byte{123},
	}

	ctx := testutils.Context(t)
	p := &CommitReportingPlugin{}
	commitStore, _ := testhelpers.NewFakeCommitStore(t, 0)
	p.config.commitStore = commitStore
	p.inflightReports = newInflightCommitReportsContainer(time.Minute)
	p.lggr = logger.TestLogger(t)

	t.Run("should transmit when report is not stale", func(t *testing.T) {
		onChainSeqNum := uint64(100)
		commitStore.SetNextSequenceNumber(onChainSeqNum)
		// not-stale since report interval is not behind on chain seq num
		report.Interval = commit_store.CommitStoreInterval{Min: onChainSeqNum, Max: onChainSeqNum + 10}
		encodedReport, err := abihelpers.EncodeCommitReport(report)
		assert.NoError(t, err)
		shouldTransmit, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.True(t, shouldTransmit)
	})

	t.Run("should not transmit when report is stale", func(t *testing.T) {
		onChainSeqNum := uint64(100)
		commitStore.SetNextSequenceNumber(onChainSeqNum)
		// stale since report interval is behind on chain seq num
		report.Interval = commit_store.CommitStoreInterval{Min: onChainSeqNum - 2, Max: onChainSeqNum + 10}
		encodedReport, err := abihelpers.EncodeCommitReport(report)
		assert.NoError(t, err)
		shouldTransmit, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldTransmit)
	})

	t.Run("error when report cannot be decoded", func(t *testing.T) {
		_, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, []byte("whatever"))
		assert.Error(t, err)
	})
}

func TestCommitReportingPlugin_calculatePriceUpdates(t *testing.T) {
	const defaultSourceChainSelector = 10 // we reuse this value across all test cases
	feeToken1 := common.HexToAddress("0xa")
	feeToken2 := common.HexToAddress("0xb")
	zero := big.NewInt(0)

	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	testCases := []struct {
		name                  string
		commitObservations    []CommitObservation
		f                     int
		latestGasPrice        update
		latestTokenPrices     map[common.Address]update
		feeUpdateHeartBeat    models.Duration
		feeUpdateDeviationPPB uint32
		expGas                *big.Int
		expTokenUpdates       []commit_store.InternalTokenPriceUpdate
		expDestChainSel       uint64
	}{
		{
			name: "median",
			commitObservations: []CommitObservation{
				{SourceGasPriceUSD: big.NewInt(1)},
				{SourceGasPriceUSD: big.NewInt(2)},
				{SourceGasPriceUSD: big.NewInt(3)},
				{SourceGasPriceUSD: big.NewInt(4)},
			},
			f:               2,
			expGas:          big.NewInt(3),
			expDestChainSel: defaultSourceChainSelector,
		},
		{
			name: "insufficient",
			commitObservations: []CommitObservation{
				{SourceGasPriceUSD: nil},
				{SourceGasPriceUSD: nil},
				{SourceGasPriceUSD: big.NewInt(3)},
			},
			f:      1,
			expGas: big.NewInt(0),
		},
		{
			name: "median including empties",
			commitObservations: []CommitObservation{
				{SourceGasPriceUSD: nil},
				{SourceGasPriceUSD: big.NewInt(1)},
				{SourceGasPriceUSD: big.NewInt(2)},
			},
			f:               1,
			expGas:          big.NewInt(2),
			expDestChainSel: defaultSourceChainSelector,
		},
		{
			name: "gas price update skipped because the latest is similar and was updated recently",
			commitObservations: []CommitObservation{
				{SourceGasPriceUSD: val1e18(10)},
				{SourceGasPriceUSD: val1e18(11)},
			},
			feeUpdateHeartBeat:    models.MustMakeDuration(time.Hour),
			feeUpdateDeviationPPB: 20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute), // recent
				value:     val1e18(9),                        // latest value close to the update
			},
			f:               1,
			expGas:          zero,
			expDestChainSel: 0,
		},
		{
			name: "gas price update included, the latest is similar but was not updated recently",
			commitObservations: []CommitObservation{
				{SourceGasPriceUSD: val1e18(10)},
				{SourceGasPriceUSD: val1e18(11)},
			},
			feeUpdateHeartBeat:    models.MustMakeDuration(time.Hour),
			feeUpdateDeviationPPB: 20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-90 * time.Minute), // recent
				value:     val1e18(9),                        // latest value close to the update
			},
			f:               1,
			expGas:          val1e18(11),
			expDestChainSel: defaultSourceChainSelector,
		},
		{
			name: "gas price update deviates from latest",
			commitObservations: []CommitObservation{
				{SourceGasPriceUSD: val1e18(10)},
				{SourceGasPriceUSD: val1e18(20)},
				{SourceGasPriceUSD: val1e18(20)},
			},
			feeUpdateHeartBeat:    models.MustMakeDuration(time.Hour),
			feeUpdateDeviationPPB: 20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute), // recent
				value:     val1e18(11),                       // latest value close to the update
			},
			f:               2,
			expGas:          val1e18(20),
			expDestChainSel: defaultSourceChainSelector,
		},
		{
			name: "median one token",
			commitObservations: []CommitObservation{
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: big.NewInt(10)}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: big.NewInt(12)}},
			},
			f: 1,
			expTokenUpdates: []commit_store.InternalTokenPriceUpdate{
				{SourceToken: feeToken1, UsdPerToken: big.NewInt(12)},
			},
			expGas:          zero,
			expDestChainSel: 0,
		},
		{
			name: "median two tokens, including nil",
			commitObservations: []CommitObservation{
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: nil, feeToken2: nil}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: big.NewInt(10), feeToken2: big.NewInt(13)}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: big.NewInt(12), feeToken2: big.NewInt(7)}},
			},
			f: 1,
			expTokenUpdates: []commit_store.InternalTokenPriceUpdate{
				{SourceToken: feeToken1, UsdPerToken: big.NewInt(12)},
				{SourceToken: feeToken2, UsdPerToken: big.NewInt(13)},
			},
			expGas:          zero,
			expDestChainSel: 0,
		},
		{
			name: "only one token with enough votes",
			commitObservations: []CommitObservation{
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: big.NewInt(10)}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: big.NewInt(12), feeToken2: big.NewInt(7)}},
			},
			f: 1,
			expTokenUpdates: []commit_store.InternalTokenPriceUpdate{
				{SourceToken: feeToken1, UsdPerToken: big.NewInt(12)},
			},
			expGas:          zero,
			expDestChainSel: 0,
		},
		{
			name: "token price update skipped because it is close to the latest",
			commitObservations: []CommitObservation{
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: val1e18(10)}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: val1e18(11), feeToken2: val1e18(7)}},
			},
			f:                     1,
			feeUpdateHeartBeat:    models.MustMakeDuration(time.Hour),
			feeUpdateDeviationPPB: 20e7,
			latestTokenPrices: map[common.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-30 * time.Minute),
					value:     val1e18(9),
				},
			},
			expGas:          zero,
			expDestChainSel: 0,
		},
		{
			name: "token price update is close to the latest but included because it has not been updated recently",
			commitObservations: []CommitObservation{
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: val1e18(10)}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: val1e18(11), feeToken2: val1e18(7)}},
			},
			f:                     1,
			feeUpdateHeartBeat:    models.MustMakeDuration(50 * time.Minute),
			feeUpdateDeviationPPB: 20e7,
			latestTokenPrices: map[common.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-1 * time.Hour),
					value:     val1e18(9),
				},
			},
			expTokenUpdates: []commit_store.InternalTokenPriceUpdate{
				{SourceToken: feeToken1, UsdPerToken: val1e18(11)},
			},
			expGas:          zero,
			expDestChainSel: 0,
		},
		{
			name: "token price update included because it is not close to the latest",
			commitObservations: []CommitObservation{
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: val1e18(20)}},
				{TokenPricesUSD: map[common.Address]*big.Int{feeToken1: val1e18(21), feeToken2: val1e18(7)}},
			},
			f:                     1,
			feeUpdateHeartBeat:    models.MustMakeDuration(time.Hour),
			feeUpdateDeviationPPB: 20e7,
			latestTokenPrices: map[common.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-30 * time.Minute),
					value:     val1e18(9),
				},
			},
			expTokenUpdates: []commit_store.InternalTokenPriceUpdate{
				{SourceToken: feeToken1, UsdPerToken: val1e18(21)},
			},
			expGas:          zero,
			expDestChainSel: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &CommitReportingPlugin{
				lggr:   logger.TestLogger(t),
				config: CommitPluginConfig{sourceChainSelector: defaultSourceChainSelector},
				offchainConfig: ccipconfig.CommitOffchainConfig{
					FeeUpdateHeartBeat:    tc.feeUpdateHeartBeat,
					FeeUpdateDeviationPPB: tc.feeUpdateDeviationPPB,
				},
				F: tc.f,
			}
			got := r.calculatePriceUpdates(tc.commitObservations, tc.latestGasPrice, tc.latestTokenPrices)

			assert.Equal(t, tc.expGas, got.UsdPerUnitGas)
			assert.Equal(t, tc.expTokenUpdates, got.TokenPriceUpdates)
			assert.Equal(t, tc.expDestChainSel, got.DestChainSelector)
		})
	}
}

func TestCommitReportingPlugin_generatePriceUpdates(t *testing.T) {
	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	const nTokens = 10
	tokens := make([]common.Address, nTokens)
	for i := range tokens {
		tokens[i] = utils.RandomAddress()
	}
	sort.Slice(tokens, func(i, j int) bool { return tokens[i].String() < tokens[j].String() })

	testCases := []struct {
		name                      string
		tokenDecimals             map[common.Address]uint8
		sourceNativeToken         common.Address
		priceGetterRespData       map[common.Address]*big.Int
		priceGetterRespErr        error
		sourceFeeEstimatorRespFee gas.EvmFee
		sourceFeeEstimatorRespErr error
		maxGasPrice               uint64
		expSourceGasPriceUSD      *big.Int
		expTokenPricesUSD         map[common.Address]*big.Int
		expErr                    bool
	}{
		{
			name: "base",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			priceGetterRespErr: nil,
			sourceFeeEstimatorRespFee: gas.EvmFee{
				Legacy:        assets.NewWei(big.NewInt(10)),
				DynamicFeeCap: nil,
				DynamicTipCap: nil,
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               1e18,
			expSourceGasPriceUSD:      big.NewInt(1000),
			expTokenPricesUSD: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			expErr: false,
		},
		{
			name: "price getter returned an error",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken:   tokens[0],
			priceGetterRespData: nil,
			priceGetterRespErr:  fmt.Errorf("some random network error"),
			expErr:              true,
		},
		{
			name: "price getter skipped a requested price",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name: "price getter skipped source native price",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[2],
			priceGetterRespData: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name: "base",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it
			},
			priceGetterRespErr: nil,
			sourceFeeEstimatorRespFee: gas.EvmFee{
				Legacy:        assets.NewWei(big.NewInt(10)),
				DynamicFeeCap: nil,
				DynamicTipCap: nil,
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               1e18,
			expSourceGasPriceUSD:      big.NewInt(1000),
			expTokenPricesUSD: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			expErr: false,
		},
		{
			name: "dynamic fee cap overrides legacy",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			priceGetterRespErr: nil,
			sourceFeeEstimatorRespFee: gas.EvmFee{
				Legacy:        assets.NewWei(big.NewInt(10)),
				DynamicFeeCap: assets.NewWei(big.NewInt(20)),
				DynamicTipCap: nil,
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               1e18,
			expSourceGasPriceUSD:      big.NewInt(2000),
			expTokenPricesUSD: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			expErr: false,
		},
		{
			name: "nil gas price",
			tokenDecimals: map[common.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[common.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			sourceFeeEstimatorRespFee: gas.EvmFee{
				Legacy:        nil,
				DynamicFeeCap: nil,
				DynamicTipCap: nil,
			},
			maxGasPrice: 1e18,
			expErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceGetter := pricegetter.NewMockPriceGetter(t)
			defer priceGetter.AssertExpectations(t)

			sourceFeeEstimator := mocks.NewEvmFeeEstimator(t)
			defer sourceFeeEstimator.AssertExpectations(t)

			tokens := make([]common.Address, 0, len(tc.tokenDecimals))
			for tk := range tc.tokenDecimals {
				tokens = append(tokens, tk)
			}
			tokens = append(tokens, tc.sourceNativeToken)
			sort.Slice(tokens, func(i, j int) bool { return tokens[i].String() < tokens[j].String() })

			if len(tokens) > 0 {
				priceGetter.On("TokenPricesUSD", mock.Anything, tokens).Return(tc.priceGetterRespData, tc.priceGetterRespErr)
			}

			if tc.maxGasPrice > 0 {
				sourceFeeEstimator.On("GetFee", mock.Anything, []byte(nil), uint32(0), assets.NewWei(big.NewInt(int64(tc.maxGasPrice)))).Return(
					tc.sourceFeeEstimatorRespFee, uint32(0), tc.sourceFeeEstimatorRespErr)
			}

			p := &CommitReportingPlugin{
				config: CommitPluginConfig{
					sourceNative:       tc.sourceNativeToken,
					priceGetter:        priceGetter,
					sourceFeeEstimator: sourceFeeEstimator,
				},
				offchainConfig: ccipconfig.CommitOffchainConfig{MaxGasPrice: tc.maxGasPrice},
			}

			sourceGasPriceUSD, tokenPricesUSD, err := p.generatePriceUpdates(context.Background(), logger.TestLogger(t), tc.tokenDecimals)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, tc.expSourceGasPriceUSD.Cmp(sourceGasPriceUSD) == 0)
			assert.True(t, reflect.DeepEqual(tc.expTokenPricesUSD, tokenPricesUSD))
		})
	}
}

func TestCommitReportingPlugin_nextMinSeqNum(t *testing.T) {
	lggr := logger.TestLogger(t)
	root1 := utils.Keccak256Fixed(hexutil.MustDecode("0xaa"))

	var tt = []struct {
		onChainMin          uint64
		inflight            []commit_store.CommitStoreCommitReport
		expectedOnChainMin  uint64
		expectedInflightMin uint64
	}{
		{
			onChainMin:          uint64(1),
			inflight:            nil,
			expectedInflightMin: uint64(1),
			expectedOnChainMin:  uint64(1),
		},
		{
			onChainMin: uint64(1),
			inflight: []commit_store.CommitStoreCommitReport{
				{Interval: commit_store.CommitStoreInterval{Min: uint64(1), Max: uint64(2)}, MerkleRoot: root1}},
			expectedInflightMin: uint64(3),
			expectedOnChainMin:  uint64(1),
		},
		{
			onChainMin: uint64(1),
			inflight: []commit_store.CommitStoreCommitReport{
				{Interval: commit_store.CommitStoreInterval{Min: uint64(3), Max: uint64(4)}, MerkleRoot: root1}},
			expectedInflightMin: uint64(5),
			expectedOnChainMin:  uint64(1),
		},
		{
			onChainMin: uint64(1),
			inflight: []commit_store.CommitStoreCommitReport{
				{Interval: commit_store.CommitStoreInterval{Min: uint64(1), Max: uint64(MaxInflightSeqNumGap + 2)}, MerkleRoot: root1}},
			expectedInflightMin: uint64(1),
			expectedOnChainMin:  uint64(1),
		},
	}
	for _, tc := range tt {
		commitStore, _ := testhelpers.NewFakeCommitStore(t, tc.onChainMin)
		cp := CommitReportingPlugin{config: CommitPluginConfig{commitStore: commitStore}, inflightReports: newInflightCommitReportsContainer(time.Hour)}
		epochAndRound := uint64(1)
		for _, rep := range tc.inflight {
			rc := rep
			require.NoError(t, cp.inflightReports.add(lggr, rc, epochAndRound))
			epochAndRound++
		}
		t.Log("inflight", cp.inflightReports.maxInflightSeqNr())
		inflightMin, onchainMin, err := cp.nextMinSeqNum(context.Background(), lggr)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedInflightMin, inflightMin)
		assert.Equal(t, tc.expectedOnChainMin, onchainMin)
		cp.inflightReports.reset(lggr)
	}
}

func TestCommitReportingPlugin_isStaleReport(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)
	merkleRoot1 := utils.Keccak256Fixed([]byte("some merkle root 1"))
	merkleRoot2 := utils.Keccak256Fixed([]byte("some merkle root 2"))

	t.Run("empty report", func(t *testing.T) {
		commitStore, _ := testhelpers.NewFakeCommitStore(t, 1)
		r := &CommitReportingPlugin{config: CommitPluginConfig{commitStore: commitStore}}
		isStale := r.isStaleReport(ctx, lggr, commit_store.CommitStoreCommitReport{}, false, types.ReportTimestamp{})
		assert.True(t, isStale)
	})

	t.Run("merkle root", func(t *testing.T) {
		const expNextSeqNum = uint64(9)
		commitStore, _ := testhelpers.NewFakeCommitStore(t, expNextSeqNum)

		r := &CommitReportingPlugin{
			config: CommitPluginConfig{commitStore: commitStore},
			inflightReports: &inflightCommitReportsContainer{
				inFlight: map[[32]byte]InflightCommitReport{
					merkleRoot2: {
						report: commit_store.CommitStoreCommitReport{
							Interval: commit_store.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
						},
					},
				},
			},
		}

		assert.False(t, r.isStaleReport(ctx, lggr, commit_store.CommitStoreCommitReport{
			MerkleRoot: merkleRoot1,
			Interval:   commit_store.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
		}, false, types.ReportTimestamp{}))

		assert.True(t, r.isStaleReport(ctx, lggr, commit_store.CommitStoreCommitReport{
			MerkleRoot: merkleRoot1,
			Interval:   commit_store.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
		}, true, types.ReportTimestamp{}))

		assert.True(t, r.isStaleReport(ctx, lggr, commit_store.CommitStoreCommitReport{
			MerkleRoot: merkleRoot1}, false, types.ReportTimestamp{}))
	})
}

func TestCommitReportingPlugin_calculateMinMaxSequenceNumbers(t *testing.T) {
	testCases := []struct {
		name              string
		commitStoreSeqNum uint64
		inflightSeqNum    uint64
		msgSeqNums        []uint64

		expQueryMin uint64 // starting seq num that is used in the query to get messages
		expMin      uint64
		expMax      uint64
		expErr      bool
	}{
		{
			name:              "happy flow inflight",
			commitStoreSeqNum: 9,
			inflightSeqNum:    10,
			msgSeqNums:        []uint64{11, 12, 13, 14},
			expQueryMin:       11, // inflight+1
			expMin:            11,
			expMax:            14,
			expErr:            false,
		},
		{
			name:              "happy flow no inflight",
			commitStoreSeqNum: 9,
			msgSeqNums:        []uint64{11, 12, 13, 14},
			expQueryMin:       9, // from commit store
			expMin:            11,
			expMax:            14,
			expErr:            false,
		},
		{
			name:              "gap in msg seq nums",
			commitStoreSeqNum: 10,
			inflightSeqNum:    9,
			expQueryMin:       10,
			msgSeqNums:        []uint64{11, 12, 14},
			expErr:            true,
		},
		{
			name:              "no new messages",
			commitStoreSeqNum: 9,
			msgSeqNums:        []uint64{},
			expQueryMin:       9,
			expMin:            0,
			expMax:            0,
			expErr:            false,
		},
		{
			name:              "unordered seq nums",
			commitStoreSeqNum: 9,
			msgSeqNums:        []uint64{11, 13, 14, 10},
			expQueryMin:       9,
			expErr:            true,
		},
	}

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &CommitReportingPlugin{}
			commitStore, _ := testhelpers.NewFakeCommitStore(t, tc.commitStoreSeqNum)
			p.config.commitStore = commitStore

			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			if tc.inflightSeqNum > 0 {
				p.inflightReports.inFlight[[32]byte{}] = InflightCommitReport{
					report: commit_store.CommitStoreCommitReport{
						Interval: commit_store.CommitStoreInterval{
							Min: tc.inflightSeqNum,
							Max: tc.inflightSeqNum,
						},
					},
				}
			}

			onRampReader := ccipdata.NewMockOnRampReader(t)
			var sendReqs []ccipdata.Event[ccipdata.EVM2EVMMessage]
			for _, seqNum := range tc.msgSeqNums {
				sendReqs = append(sendReqs, ccipdata.Event[ccipdata.EVM2EVMMessage]{
					Data: ccipdata.EVM2EVMMessage{
						SequenceNumber: seqNum,
					},
				})
			}
			onRampReader.On("GetSendRequestsGteSeqNum", ctx, tc.expQueryMin, 0).Return(sendReqs, nil)
			p.config.onRampReader = onRampReader

			minSeqNum, maxSeqNum, err := p.calculateMinMaxSequenceNumbers(ctx, lggr)
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.expMin, minSeqNum)
			assert.Equal(t, tc.expMax, maxSeqNum)
		})
	}
}

func TestCommitReportingPlugin_getLatestGasPriceUpdate(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name                   string
		checkInflight          bool
		inflightGasPriceUpdate *update
		destGasPriceUpdates    []update
		expUpdate              update
		expErr                 bool
	}{
		{
			name:                   "only inflight gas price",
			checkInflight:          true,
			inflightGasPriceUpdate: &update{timestamp: now, value: big.NewInt(1000)},
			expUpdate:              update{timestamp: now, value: big.NewInt(1000)},
			expErr:                 false,
		},
		{
			name:                   "inflight price is nil",
			checkInflight:          true,
			inflightGasPriceUpdate: &update{timestamp: now, value: nil},
			destGasPriceUpdates: []update{
				{timestamp: now.Add(time.Minute), value: big.NewInt(2000)},
				{timestamp: now.Add(2 * time.Minute), value: big.NewInt(3000)},
			},
			expUpdate: update{timestamp: now.Add(2 * time.Minute), value: big.NewInt(3000)},
			expErr:    false,
		},
		{
			name:                   "inflight updates are skipped",
			checkInflight:          false,
			inflightGasPriceUpdate: &update{timestamp: now, value: big.NewInt(1000)},
			destGasPriceUpdates: []update{
				{timestamp: now.Add(time.Minute), value: big.NewInt(2000)},
				{timestamp: now.Add(2 * time.Minute), value: big.NewInt(3000)},
			},
			expUpdate: update{timestamp: now.Add(2 * time.Minute), value: big.NewInt(3000)},
			expErr:    false,
		},
	}

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &CommitReportingPlugin{}
			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			p.lggr = lggr
			destPriceRegistry, _ := testhelpers.NewFakePriceRegistry(t)
			p.destPriceRegistry = destPriceRegistry

			if tc.inflightGasPriceUpdate != nil {
				p.inflightReports.inFlightPriceUpdates = append(
					p.inflightReports.inFlightPriceUpdates,
					InflightPriceUpdate{
						createdAt: tc.inflightGasPriceUpdate.timestamp,
						priceUpdates: commit_store.InternalPriceUpdates{
							DestChainSelector: 1234,
							UsdPerUnitGas:     tc.inflightGasPriceUpdate.value,
						},
					},
				)
			}

			if len(tc.destGasPriceUpdates) > 0 {
				var events []ccipdata.Event[price_registry.PriceRegistryUsdPerUnitGasUpdated]
				for _, u := range tc.destGasPriceUpdates {
					events = append(events, ccipdata.Event[price_registry.PriceRegistryUsdPerUnitGasUpdated]{
						Data: price_registry.PriceRegistryUsdPerUnitGasUpdated{
							Value:     u.value,
							Timestamp: big.NewInt(u.timestamp.Unix()),
						},
					})
				}
				destReader := ccipdata.NewMockReader(t)
				destReader.On("GetGasPriceUpdatesCreatedAfter", ctx, mock.Anything, uint64(0), mock.Anything, 0).Return(events, nil)
				p.config.destReader = destReader
			}

			priceUpdate, err := p.getLatestGasPriceUpdate(ctx, time.Now(), tc.checkInflight)
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expUpdate.timestamp.Truncate(time.Second), priceUpdate.timestamp.Truncate(time.Second))
			assert.Equal(t, tc.expUpdate.value.Uint64(), priceUpdate.value.Uint64())
		})
	}
}

func TestCommitReportingPlugin_getLatestTokenPriceUpdates(t *testing.T) {
	now := time.Now()
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()

	testCases := []struct {
		name                 string
		priceRegistryUpdates []price_registry.PriceRegistryUsdPerTokenUpdated
		checkInflight        bool
		inflightUpdates      map[common.Address]update
		expUpdates           map[common.Address]update
		expErr               bool
	}{
		{
			name: "ignore inflight updates",
			priceRegistryUpdates: []price_registry.PriceRegistryUsdPerTokenUpdated{
				{
					Token:     tk1,
					Value:     big.NewInt(1000),
					Timestamp: big.NewInt(now.Add(1 * time.Minute).Unix()),
				},
				{
					Token:     tk2,
					Value:     big.NewInt(2000),
					Timestamp: big.NewInt(now.Add(2 * time.Minute).Unix()),
				},
			},
			checkInflight: false,
			expUpdates: map[common.Address]update{
				tk1: {timestamp: now.Add(1 * time.Minute), value: big.NewInt(1000)},
				tk2: {timestamp: now.Add(2 * time.Minute), value: big.NewInt(2000)},
			},
			expErr: false,
		},
		{
			name: "consider inflight updates",
			priceRegistryUpdates: []price_registry.PriceRegistryUsdPerTokenUpdated{
				{
					Token:     tk1,
					Value:     big.NewInt(1000),
					Timestamp: big.NewInt(now.Add(1 * time.Minute).Unix()),
				},
				{
					Token:     tk2,
					Value:     big.NewInt(2000),
					Timestamp: big.NewInt(now.Add(2 * time.Minute).Unix()),
				},
			},
			checkInflight: true,
			inflightUpdates: map[common.Address]update{
				tk1: {timestamp: now, value: big.NewInt(500)}, // inflight but older
				tk2: {timestamp: now.Add(4 * time.Minute), value: big.NewInt(4000)},
			},
			expUpdates: map[common.Address]update{
				tk1: {timestamp: now.Add(1 * time.Minute), value: big.NewInt(1000)},
				tk2: {timestamp: now.Add(4 * time.Minute), value: big.NewInt(4000)},
			},
			expErr: false,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &CommitReportingPlugin{}

			priceReg, priceRegAddr := testhelpers.NewFakePriceRegistry(t)
			p.destPriceRegistry = priceReg

			destReader := ccipdata.NewMockReader(t)
			var events []ccipdata.Event[price_registry.PriceRegistryUsdPerTokenUpdated]
			for _, up := range tc.priceRegistryUpdates {
				events = append(events, ccipdata.Event[price_registry.PriceRegistryUsdPerTokenUpdated]{
					Data: price_registry.PriceRegistryUsdPerTokenUpdated{
						Token:     up.Token,
						Value:     up.Value,
						Timestamp: up.Timestamp,
					},
				})
			}
			destReader.On("GetTokenPriceUpdatesCreatedAfter", ctx, priceRegAddr, mock.Anything, 0).Return(events, nil)
			p.config.destReader = destReader

			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			if len(tc.inflightUpdates) > 0 {
				for tk, upd := range tc.inflightUpdates {
					p.inflightReports.inFlightPriceUpdates = append(p.inflightReports.inFlightPriceUpdates, InflightPriceUpdate{
						createdAt: upd.timestamp,
						priceUpdates: commit_store.InternalPriceUpdates{
							TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{
								{SourceToken: tk, UsdPerToken: upd.value},
							},
						},
					})
				}
			}

			updates, err := p.getLatestTokenPriceUpdates(ctx, now, tc.checkInflight)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expUpdates), len(updates))
			for k, v := range updates {
				assert.Equal(t, tc.expUpdates[k].timestamp.Truncate(time.Second), v.timestamp.Truncate(time.Second))
				assert.Equal(t, tc.expUpdates[k].value.Uint64(), v.value.Uint64())
			}
		})
	}

}

func Test_commitReportSize(t *testing.T) {
	testParams := gopter.DefaultTestParameters()
	testParams.MinSuccessfulTests = 100
	p := gopter.NewProperties(testParams)
	p.Property("bounded commit report size", prop.ForAll(func(root []byte, min, max uint64) bool {
		var root32 [32]byte
		copy(root32[:], root)
		rep, err := abihelpers.EncodeCommitReport(commit_store.CommitStoreCommitReport{
			MerkleRoot: root32,
			Interval:   commit_store.CommitStoreInterval{Min: min, Max: max},
			PriceUpdates: commit_store.InternalPriceUpdates{
				TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
				DestChainSelector: 1337,
				UsdPerUnitGas:     big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
			},
		})
		require.NoError(t, err)
		return len(rep) <= MaxCommitReportLength
	}, gen.SliceOfN(32, gen.UInt8()), gen.UInt64(), gen.UInt64()))
	p.TestingRun(t)
}

func Test_calculateIntervalConsensus(t *testing.T) {
	tests := []struct {
		name       string
		intervals  []commit_store.CommitStoreInterval
		rangeLimit uint64
		f          int
		wantMin    uint64
		wantMax    uint64
		wantErr    bool
	}{
		{"no obs", []commit_store.CommitStoreInterval{{Min: 0, Max: 0}}, 0, 0, 0, 0, false},
		{"basic", []commit_store.CommitStoreInterval{
			{Min: 9, Max: 14},
			{Min: 10, Max: 12},
			{Min: 10, Max: 14},
		}, 0, 1, 10, 14, false},
		{"not enough intervals", []commit_store.CommitStoreInterval{}, 0, 1, 0, 0, true},
		{"min > max", []commit_store.CommitStoreInterval{
			{Min: 9, Max: 4},
			{Min: 10, Max: 4},
			{Min: 10, Max: 6},
		}, 0, 1, 0, 0, true},
		{
			"range limit", []commit_store.CommitStoreInterval{
				{Min: 10, Max: 100},
				{Min: 1, Max: 1000},
			}, 256, 1, 10, 265, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateIntervalConsensus(tt.intervals, tt.f, tt.rangeLimit)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantMin, got.Min)
			assert.Equal(t, tt.wantMax, got.Max)
		})
	}
}

func Test_calculateUsdPer1e18TokenAmount(t *testing.T) {
	tests := []struct {
		name       string
		price      *big.Int
		decimal    uint8
		wantResult *big.Int
	}{
		{
			name:       "18-decimal token, $6.5 per token",
			price:      big.NewInt(65e17),
			decimal:    18,
			wantResult: big.NewInt(65e17),
		},
		{
			name:       "6-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    6,
			wantResult: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e12)), // 1e30
		},
		{
			name:       "0-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    0,
			wantResult: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e18)), // 1e36
		},
		{
			name:       "36-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    36,
			wantResult: big.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateUsdPer1e18TokenAmount(tt.price, tt.decimal)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestCommitReportToEthTxMeta(t *testing.T) {
	mctx := hashlib.NewKeccakCtx()
	tree, err := merklemulti.NewTree(mctx, [][32]byte{mctx.Hash([]byte{0xaa})})
	require.NoError(t, err)

	tests := []struct {
		name          string
		min, max      uint64
		expectedRange []uint64
	}{
		{
			"happy flow",
			1, 10,
			[]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			"same sequence",
			1, 1,
			[]uint64{1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			report := commit_store.CommitStoreCommitReport{
				PriceUpdates: commit_store.InternalPriceUpdates{
					TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
					DestChainSelector: uint64(1337),
					UsdPerUnitGas:     big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
				},
				MerkleRoot: tree.Root(),
				Interval:   commit_store.CommitStoreInterval{Min: tc.min, Max: tc.max},
			}
			out, err := abihelpers.EncodeCommitReport(report)
			require.NoError(t, err)

			txMeta, err := CommitReportToEthTxMeta(out)
			require.NoError(t, err)
			require.NotNil(t, txMeta)
			require.EqualValues(t, tc.expectedRange, txMeta.SeqNumbers)
		})
	}
}
