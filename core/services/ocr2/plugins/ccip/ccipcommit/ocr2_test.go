package ccipcommit

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	ccipcachemocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/pkg/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/pkg/merklemulti"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

func TestCommitReportingPlugin_Observation(t *testing.T) {
	sourceNativeTokenAddr := ccipcalc.HexToAddress("1000")

	bridgedTokens := []cciptypes.Address{
		ccipcalc.HexToAddress("2000"),
		ccipcalc.HexToAddress("3000"),
	}

	// Token price in 1e18 USD precision
	bridgedTokenPrices := map[cciptypes.Address]*big.Int{
		bridgedTokens[0]: big.NewInt(1),
		bridgedTokens[1]: big.NewInt(2e18),
	}

	bridgedTokenDecimals := map[cciptypes.Address]uint8{
		bridgedTokens[0]: 8,
		bridgedTokens[1]: 18,
	}

	// Token price of 1e18 token amount in 1e18 USD precision
	expectedEncodedTokenPrice := map[cciptypes.Address]*big.Int{
		bridgedTokens[0]: big.NewInt(1e10),
		bridgedTokens[1]: big.NewInt(2e18),
	}

	testCases := []struct {
		name              string
		epochAndRound     types.ReportTimestamp
		commitStorePaused bool
		sourceChainCursed bool
		commitStoreSeqNum uint64
		tokenPrices       map[cciptypes.Address]*big.Int
		sendReqs          []cciptypes.EVM2EVMMessageWithTxMeta
		tokenDecimals     map[cciptypes.Address]uint8
		fee               *big.Int

		expErr bool
		expObs ccip.CommitObservation
	}{
		{
			name:              "base report",
			commitStoreSeqNum: 54,
			tokenPrices: map[cciptypes.Address]*big.Int{
				bridgedTokens[0]:      bridgedTokenPrices[bridgedTokens[0]],
				bridgedTokens[1]:      bridgedTokenPrices[bridgedTokens[1]],
				sourceNativeTokenAddr: big.NewInt(2e18),
			},
			sendReqs: []cciptypes.EVM2EVMMessageWithTxMeta{
				{EVM2EVMMessage: cciptypes.EVM2EVMMessage{SequenceNumber: 54}},
				{EVM2EVMMessage: cciptypes.EVM2EVMMessage{SequenceNumber: 55}},
			},
			fee:           big.NewInt(2e18),
			tokenDecimals: bridgedTokenDecimals,
			expObs: ccip.CommitObservation{
				TokenPricesUSD:    expectedEncodedTokenPrice,
				SourceGasPriceUSD: big.NewInt(4e18),
				Interval: cciptypes.CommitStoreInterval{
					Min: 54,
					Max: 55,
				},
			},
		},
		{
			name:              "commit store is down",
			commitStorePaused: true,
			sourceChainCursed: false,
			expErr:            true,
		},
		{
			name:              "source chain is cursed",
			commitStorePaused: false,
			sourceChainCursed: true,
			expErr:            true,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
			commitStoreReader.On("IsDown", ctx).Return(tc.commitStorePaused, nil)
			commitStoreReader.On("IsDestChainHealthy", ctx).Return(true, nil)
			if !tc.commitStorePaused && !tc.sourceChainCursed {
				commitStoreReader.On("GetExpectedNextSequenceNumber", ctx).Return(tc.commitStoreSeqNum, nil)
			}

			onRampReader := ccipdatamocks.NewOnRampReader(t)
			onRampReader.On("IsSourceChainHealthy", ctx).Return(true, nil)
			onRampReader.On("IsSourceCursed", ctx).Return(tc.sourceChainCursed, nil)
			if len(tc.sendReqs) > 0 {
				onRampReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.commitStoreSeqNum, tc.commitStoreSeqNum+OnRampMessagesScanLimit, true).
					Return(tc.sendReqs, nil)
			}

			var destTokens []cciptypes.Address
			for tk := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
			}
			// ensure destTokens and destDecimals are in the same order, avoid flaky test from unordered map iteration
			sort.Slice(destTokens, func(i, j int) bool {
				return destTokens[i] < destTokens[j]
			})
			var destDecimals []uint8
			for _, token := range destTokens {
				destDecimals = append(destDecimals, tc.tokenDecimals[token])
			}

			priceGet := pricegetter.NewMockPriceGetter(t)
			if len(tc.tokenPrices) > 0 {
				queryTokens := ccipcommon.FlattenUniqueSlice([]cciptypes.Address{sourceNativeTokenAddr}, destTokens)
				priceGet.On("TokenPricesUSD", mock.Anything, queryTokens).Return(tc.tokenPrices, nil)
				priceGet.On("FilterConfiguredTokens", mock.Anything, destTokens).Return([]cciptypes.Address{
					bridgedTokens[0],
					bridgedTokens[1],
				}, []cciptypes.Address{}, nil)
			}

			gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
			if tc.fee != nil {
				var p = tc.fee
				var pUSD = ccipcalc.CalculateUsdPerUnitGas(p, tc.tokenPrices[sourceNativeTokenAddr])
				gasPriceEstimator.On("GetGasPrice", ctx).Return(p, nil)
				gasPriceEstimator.On("DenoteInUSD", p, tc.tokenPrices[sourceNativeTokenAddr]).Return(pUSD, nil)
			}

			offRampReader := ccipdatamocks.NewOffRampReader(t)
			offRampReader.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{
				DestinationTokens: destTokens,
			}, nil).Maybe()

			destPriceRegReader := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceRegReader.On("GetFeeTokens", ctx).Return(nil, nil).Maybe()
			destPriceRegReader.On("GetTokensDecimals", ctx, destTokens).Return(destDecimals, nil).Maybe()

			p := &CommitReportingPlugin{}
			p.lggr = logger.TestLogger(t)
			p.commitStoreReader = commitStoreReader
			p.onRampReader = onRampReader
			p.offRampReader = offRampReader
			p.destPriceRegistryReader = destPriceRegReader
			p.priceGetter = priceGet
			p.sourceNative = sourceNativeTokenAddr
			p.gasPriceEstimator = gasPriceEstimator
			p.metricsCollector = ccip.NoopMetricsCollector
			p.chainHealthcheck = cache.NewChainHealthcheck(p.lggr, onRampReader, commitStoreReader)

			obs, err := p.Observation(ctx, tc.epochAndRound, types.Query{})

			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tc.expObs.TokenPricesUSD != nil {
				// field ordering in mapping is not guaranteed, if TokenPricesUSD exists, unmarshal to compare mapping
				var obsStuct ccip.CommitObservation
				err = json.Unmarshal(obs, &obsStuct)
				assert.NoError(t, err)

				assert.Equal(t, tc.expObs, obsStuct)
			} else {
				// if TokenPricesUSD is nil, compare the bytes directly, marshal then unmarshal turns nil map to empty
				expObsBytes, err := tc.expObs.Marshal()
				assert.NoError(t, err)
				assert.Equal(t, expObsBytes, []byte(obs))
			}
		})
	}
}

func TestCommitReportingPlugin_Report(t *testing.T) {
	ctx := testutils.Context(t)
	sourceChainSelector := uint64(rand.Int())
	var gasPrice = big.NewInt(1)
	gasPriceHeartBeat := *config.MustNewDuration(time.Hour)

	t.Run("not enough observations", func(t *testing.T) {
		p := &CommitReportingPlugin{}
		p.lggr = logger.TestLogger(t)
		p.F = 1

		chainHealthcheck := ccipcachemocks.NewChainHealthcheck(t)
		chainHealthcheck.On("IsHealthy", ctx).Return(true, nil).Maybe()
		p.chainHealthcheck = chainHealthcheck

		o := ccip.CommitObservation{Interval: cciptypes.CommitStoreInterval{Min: 1, Max: 1}, SourceGasPriceUSD: big.NewInt(0)}
		obs, err := o.Marshal()
		assert.NoError(t, err)

		aos := []types.AttributedObservation{{Observation: obs}}

		gotSomeReport, gotReport, err := p.Report(ctx, types.ReportTimestamp{}, types.Query{}, aos)
		assert.False(t, gotSomeReport)
		assert.Nil(t, gotReport)
		assert.Error(t, err)
	})

	testCases := []struct {
		name              string
		observations      []ccip.CommitObservation
		f                 int
		gasPriceUpdates   []cciptypes.GasPriceUpdateWithTxMeta
		tokenDecimals     map[cciptypes.Address]uint8
		tokenPriceUpdates []cciptypes.TokenPriceUpdateWithTxMeta
		sendRequests      []cciptypes.EVM2EVMMessageWithTxMeta
		expCommitReport   *cciptypes.CommitStoreReport
		expSeqNumRange    cciptypes.CommitStoreInterval
		expErr            bool
	}{
		{
			name: "base",
			observations: []ccip.CommitObservation{
				{Interval: cciptypes.CommitStoreInterval{Min: 1, Max: 1}, SourceGasPriceUSD: gasPrice},
				{Interval: cciptypes.CommitStoreInterval{Min: 1, Max: 1}, SourceGasPriceUSD: gasPrice},
			},
			f: 1,
			sendRequests: []cciptypes.EVM2EVMMessageWithTxMeta{
				{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: 1,
					},
				},
			},
			gasPriceUpdates: []cciptypes.GasPriceUpdateWithTxMeta{
				{
					GasPriceUpdate: cciptypes.GasPriceUpdate{
						GasPrice: cciptypes.GasPrice{
							DestChainSelector: sourceChainSelector,
							Value:             big.NewInt(1),
						},
						TimestampUnixSec: big.NewInt(time.Now().Add(-2 * gasPriceHeartBeat.Duration()).Unix()),
					},
				},
			},
			expSeqNumRange: cciptypes.CommitStoreInterval{Min: 1, Max: 1},
			expCommitReport: &cciptypes.CommitStoreReport{
				MerkleRoot:  [32]byte{},
				Interval:    cciptypes.CommitStoreInterval{Min: 1, Max: 1},
				TokenPrices: nil,
				GasPrices:   []cciptypes.GasPrice{{DestChainSelector: sourceChainSelector, Value: gasPrice}},
			},
			expErr: false,
		},
		{
			name: "empty",
			observations: []ccip.CommitObservation{
				{Interval: cciptypes.CommitStoreInterval{Min: 0, Max: 0}, SourceGasPriceUSD: big.NewInt(0)},
				{Interval: cciptypes.CommitStoreInterval{Min: 0, Max: 0}, SourceGasPriceUSD: big.NewInt(0)},
			},
			gasPriceUpdates: []cciptypes.GasPriceUpdateWithTxMeta{
				{
					GasPriceUpdate: cciptypes.GasPriceUpdate{
						GasPrice: cciptypes.GasPrice{
							DestChainSelector: sourceChainSelector,
							Value:             big.NewInt(1),
						},
						TimestampUnixSec: big.NewInt(time.Now().Add(-gasPriceHeartBeat.Duration() / 2).Unix()),
					},
				},
			},
			f:      1,
			expErr: false,
		},
		{
			name: "no leaves",
			observations: []ccip.CommitObservation{
				{Interval: cciptypes.CommitStoreInterval{Min: 2, Max: 2}, SourceGasPriceUSD: big.NewInt(0)},
				{Interval: cciptypes.CommitStoreInterval{Min: 2, Max: 2}, SourceGasPriceUSD: big.NewInt(0)},
			},
			f:              1,
			sendRequests:   []cciptypes.EVM2EVMMessageWithTxMeta{{}},
			expSeqNumRange: cciptypes.CommitStoreInterval{Min: 2, Max: 2},
			expErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			destPriceRegistryReader := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceRegistryReader.On("GetGasPriceUpdatesCreatedAfter", ctx, sourceChainSelector, mock.Anything, 0).Return(tc.gasPriceUpdates, nil)
			destPriceRegistryReader.On("GetTokenPriceUpdatesCreatedAfter", ctx, mock.Anything, 0).Return(tc.tokenPriceUpdates, nil)

			onRampReader := ccipdatamocks.NewOnRampReader(t)
			if len(tc.sendRequests) > 0 {
				onRampReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.expSeqNumRange.Min, tc.expSeqNumRange.Max, true).Return(tc.sendRequests, nil)
			}

			gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
			gasPriceEstimator.On("Median", mock.Anything).Return(gasPrice, nil)
			if tc.gasPriceUpdates != nil {
				gasPriceEstimator.On("Deviates", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
			}

			destTokens := []cciptypes.Address{}
			for tk := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
			}
			sort.Slice(destTokens, func(i, j int) bool {
				return destTokens[i] < destTokens[j]
			})
			var destDecimals []uint8
			for _, token := range destTokens {
				destDecimals = append(destDecimals, tc.tokenDecimals[token])
			}

			destPriceRegistryReader.On("GetTokensDecimals", ctx, mock.MatchedBy(func(tokens []cciptypes.Address) bool {
				for _, token := range tokens {
					if !slices.Contains(destTokens, token) {
						return false
					}
				}
				return true
			})).Return(destDecimals, nil).Maybe()

			lp := mocks2.NewLogPoller(t)
			commitStoreReader, err := v1_2_0.NewCommitStore(logger.TestLogger(t), utils.RandomAddress(), nil, lp, nil, nil)
			assert.NoError(t, err)

			healthCheck := ccipcachemocks.NewChainHealthcheck(t)
			healthCheck.On("IsHealthy", ctx).Return(true, nil)

			p := &CommitReportingPlugin{}
			p.lggr = logger.TestLogger(t)
			p.destPriceRegistryReader = destPriceRegistryReader
			p.onRampReader = onRampReader
			p.sourceChainSelector = sourceChainSelector
			p.gasPriceEstimator = gasPriceEstimator
			p.offchainConfig.GasPriceHeartBeat = gasPriceHeartBeat.Duration()
			p.commitStoreReader = commitStoreReader
			p.F = tc.f
			p.metricsCollector = ccip.NoopMetricsCollector
			p.chainHealthcheck = healthCheck

			aos := make([]types.AttributedObservation, 0, len(tc.observations))
			for _, o := range tc.observations {
				obs, err2 := o.Marshal()
				assert.NoError(t, err2)
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
				encodedExpectedReport, err := encodeCommitReport(*tc.expCommitReport)
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
		p.metricsCollector = ccip.NoopMetricsCollector
		return p
	}

	t.Run("report cannot be decoded leads to error", func(t *testing.T) {
		p := newPlugin()

		encodedReport := []byte("whatever")

		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p.commitStoreReader = commitStoreReader
		commitStoreReader.On("DecodeCommitReport", mock.Anything, encodedReport).
			Return(cciptypes.CommitStoreReport{}, errors.New("unable to decode report"))

		_, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.Error(t, err)
	})

	t.Run("empty report should not be accepted", func(t *testing.T) {
		p := newPlugin()

		report := cciptypes.CommitStoreReport{}

		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p.commitStoreReader = commitStoreReader
		commitStoreReader.On("DecodeCommitReport", mock.Anything, mock.Anything).Return(report, nil)

		chainHealthCheck := ccipcachemocks.NewChainHealthcheck(t)
		chainHealthCheck.On("IsHealthy", ctx).Return(true, nil).Maybe()
		p.chainHealthcheck = chainHealthCheck

		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)
		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldAccept)
	})

	t.Run("stale report should not be accepted", func(t *testing.T) {
		onChainSeqNum := uint64(100)

		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p := newPlugin()

		p.commitStoreReader = commitStoreReader

		report := cciptypes.CommitStoreReport{
			GasPrices:  []cciptypes.GasPrice{{Value: big.NewInt(int64(rand.Int()))}},
			MerkleRoot: [32]byte{123}, // this report is considered non-empty since it has a merkle root
		}

		commitStoreReader.On("DecodeCommitReport", mock.Anything, mock.Anything).Return(report, nil)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(onChainSeqNum, nil)

		chainHealthCheck := ccipcachemocks.NewChainHealthcheck(t)
		chainHealthCheck.On("IsHealthy", ctx).Return(true, nil)
		p.chainHealthcheck = chainHealthCheck

		// stale since report interval is behind on chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum - 2, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)

		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldAccept)
	})

	t.Run("non-stale report should be accepted", func(t *testing.T) {
		onChainSeqNum := uint64(100)

		p := newPlugin()

		priceRegistryReader := ccipdatamocks.NewPriceRegistryReader(t)
		p.destPriceRegistryReader = priceRegistryReader

		p.lggr = logger.TestLogger(t)
		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p.commitStoreReader = commitStoreReader

		report := cciptypes.CommitStoreReport{
			Interval: cciptypes.CommitStoreInterval{
				Min: onChainSeqNum,
				Max: onChainSeqNum + 10,
			},
			TokenPrices: []cciptypes.TokenPrice{
				{
					Token: cciptypes.Address(utils.RandomAddress().String()),
					Value: big.NewInt(int64(rand.Int())),
				},
			},
			GasPrices: []cciptypes.GasPrice{
				{
					DestChainSelector: rand.Uint64(),
					Value:             big.NewInt(int64(rand.Int())),
				},
			},
			MerkleRoot: [32]byte{123},
		}
		commitStoreReader.On("DecodeCommitReport", mock.Anything, mock.Anything).Return(report, nil)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(onChainSeqNum, nil)

		// non-stale since report interval is not behind on-chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)

		chainHealthCheck := ccipcachemocks.NewChainHealthcheck(t)
		chainHealthCheck.On("IsHealthy", ctx).Return(true, nil)
		p.chainHealthcheck = chainHealthCheck

		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.True(t, shouldAccept)
	})
}

func TestCommitReportingPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	report := cciptypes.CommitStoreReport{
		TokenPrices: []cciptypes.TokenPrice{
			{Token: cciptypes.Address(utils.RandomAddress().String()), Value: big.NewInt(9e18)},
		},
		GasPrices: []cciptypes.GasPrice{
			{

				DestChainSelector: rand.Uint64(),
				Value:             big.NewInt(2000e9),
			},
		},
		MerkleRoot: [32]byte{123},
	}

	ctx := testutils.Context(t)
	p := &CommitReportingPlugin{}
	commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
	onChainSeqNum := uint64(100)
	commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(onChainSeqNum, nil)
	p.commitStoreReader = commitStoreReader
	p.lggr = logger.TestLogger(t)

	chainHealthCheck := ccipcachemocks.NewChainHealthcheck(t)
	chainHealthCheck.On("IsHealthy", ctx).Return(true, nil).Maybe()
	p.chainHealthcheck = chainHealthCheck

	t.Run("should transmit when report is not stale", func(t *testing.T) {
		// not-stale since report interval is not behind on chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)
		commitStoreReader.On("DecodeCommitReport", mock.Anything, encodedReport).Return(report, nil).Once()
		shouldTransmit, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.True(t, shouldTransmit)
	})

	t.Run("should not transmit when report is stale", func(t *testing.T) {
		// stale since report interval is behind on chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum - 2, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)
		commitStoreReader.On("DecodeCommitReport", mock.Anything, encodedReport).Return(report, nil).Once()
		shouldTransmit, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldTransmit)
	})

	t.Run("error when report cannot be decoded", func(t *testing.T) {
		reportBytes := []byte("whatever")
		commitStoreReader.On("DecodeCommitReport", mock.Anything, reportBytes).
			Return(cciptypes.CommitStoreReport{}, errors.New("decode error")).Once()
		_, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, reportBytes)
		assert.Error(t, err)
	})
}

func TestCommitReportingPlugin_extractObservationData(t *testing.T) {
	token1 := ccipcalc.HexToAddress("0xa")
	token2 := ccipcalc.HexToAddress("0xb")
	token1Price := big.NewInt(1)
	token2Price := big.NewInt(2)
	unsupportedToken := ccipcalc.HexToAddress("0xc")
	gasPrice := big.NewInt(100)

	tokenDecimals := make(map[cciptypes.Address]uint8)
	tokenDecimals[token1] = 18
	tokenDecimals[token2] = 18

	validInterval := cciptypes.CommitStoreInterval{Min: 1, Max: 2}
	zeroInterval := cciptypes.CommitStoreInterval{Min: 0, Max: 0}

	ob1 := ccip.CommitObservation{
		Interval: validInterval,
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1: token1Price,
			token2: token2Price,
		},
		SourceGasPriceUSD: gasPrice,
	}
	ob1Bytes, err := ob1.Marshal()
	assert.NoError(t, err)
	lggr := logger.TestLogger(t)
	observations := ccip.GetParsableObservations[ccip.CommitObservation](lggr, []types.AttributedObservation{
		{Observation: ob1Bytes},
		{Observation: ob1Bytes},
	})
	assert.Len(t, observations, 2)
	ob2 := observations[0]
	ob3 := observations[1]

	obWithNilGasPrice := ccip.CommitObservation{
		Interval: zeroInterval,
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1: token1Price,
			token2: token2Price,
		},
		SourceGasPriceUSD: nil,
	}
	obWithNilTokenPrice := ccip.CommitObservation{
		Interval: zeroInterval,
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1: token1Price,
			token2: nil,
		},
		SourceGasPriceUSD: gasPrice,
	}
	obMissingTokenPrices := ccip.CommitObservation{
		Interval:          zeroInterval,
		TokenPricesUSD:    map[cciptypes.Address]*big.Int{},
		SourceGasPriceUSD: gasPrice,
	}
	obWithUnsupportedToken := ccip.CommitObservation{
		Interval: zeroInterval,
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1:           token1Price,
			token2:           token2Price,
			unsupportedToken: token2Price,
		},
		SourceGasPriceUSD: gasPrice,
	}
	obEmpty := ccip.CommitObservation{
		Interval:          zeroInterval,
		TokenPricesUSD:    nil,
		SourceGasPriceUSD: nil,
	}

	testCases := []struct {
		name               string
		commitObservations []ccip.CommitObservation
		f                  int
		expIntervals       []cciptypes.CommitStoreInterval
		expGasPriceObs     []*big.Int
		expTokenPriceObs   map[cciptypes.Address][]*big.Int
		expValidObs        []ccip.CommitObservation
		expError           bool
	}{
		{
			name:               "base",
			commitObservations: []ccip.CommitObservation{ob1, ob2},
			f:                  1,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price},
				token2: {token2Price, token2Price},
			},
			expError: false,
		},
		{
			name:               "pass with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, ob3},
			f:                  2,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval, validInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD, ob3.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price, token1Price},
				token2: {token2Price, token2Price, token2Price},
			},
			expError: false,
		},
		{
			name:               "tolerate 1 faulty obs with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, ob3, obWithNilGasPrice},
			f:                  2,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval, validInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD, ob3.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price, token1Price, token1Price},
				token2: {token2Price, token2Price, token2Price, token2Price},
			},
			expError: false,
		},
		{
			name:               "tolerate 1 nil token price with f=1",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obWithNilTokenPrice},
			f:                  1,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD, obWithNilTokenPrice.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price, token1Price},
				token2: {token2Price, token2Price},
			},
			expError: false,
		},
		{
			name:               "tolerate 1 missing token prices with f=1",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obMissingTokenPrices},
			f:                  1,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD, obMissingTokenPrices.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price},
				token2: {token2Price, token2Price},
			},
			expError: false,
		},
		{
			name:               "tolerate 1 unsupported token with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obWithUnsupportedToken},
			f:                  2,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD, obWithUnsupportedToken.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price, token1Price},
				token2: {token2Price, token2Price, token2Price},
			},
			expError: false,
		},
		{
			name:               "tolerate mis-matched token observations with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obWithNilTokenPrice, obMissingTokenPrices},
			f:                  2,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, validInterval, zeroInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, ob2.SourceGasPriceUSD, obWithNilTokenPrice.SourceGasPriceUSD, obMissingTokenPrices.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price, token1Price},
			},
			expError: false,
		},
		{
			name:               "tolerate mis-matched token observations with f=2",
			commitObservations: []ccip.CommitObservation{ob1, obWithNilTokenPrice, obWithNilTokenPrice},
			f:                  2,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, zeroInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, obWithNilTokenPrice.SourceGasPriceUSD, obWithNilTokenPrice.SourceGasPriceUSD},
			expTokenPriceObs: map[cciptypes.Address][]*big.Int{
				token1: {token1Price, token1Price, token1Price},
			},
			expError: false,
		},
		{
			name:               "tolerate all tokens filtered out with f=2",
			commitObservations: []ccip.CommitObservation{ob1, obMissingTokenPrices, obMissingTokenPrices},
			f:                  2,
			expIntervals:       []cciptypes.CommitStoreInterval{validInterval, zeroInterval, zeroInterval},
			expGasPriceObs:     []*big.Int{ob1.SourceGasPriceUSD, obMissingTokenPrices.SourceGasPriceUSD, obMissingTokenPrices.SourceGasPriceUSD},
			expTokenPriceObs:   map[cciptypes.Address][]*big.Int{},
			expError:           false,
		},
		{
			name:               "not enough observations",
			commitObservations: []ccip.CommitObservation{ob1, ob2},
			f:                  2,
			expValidObs:        nil,
			expError:           true,
		},
		{
			name:               "too many faulty observations",
			commitObservations: []ccip.CommitObservation{obWithNilGasPrice, obWithNilTokenPrice, obEmpty, obEmpty, obEmpty},
			f:                  1,
			expValidObs:        nil,
			expError:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			intervals, gasPriceOps, tokenPriceOps, err := extractObservationData(logger.TestLogger(t), tc.f, tc.commitObservations)

			if tc.expError {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.expIntervals, intervals)
			assert.Equal(t, tc.expGasPriceObs, gasPriceOps)
			assert.Equal(t, tc.expTokenPriceObs, tokenPriceOps)
			assert.NoError(t, err)
		})
	}
}

func TestCommitReportingPlugin_calculatePriceUpdates(t *testing.T) {
	const defaultSourceChainSelector = 10 // we reuse this value across all test cases
	feeToken1 := ccipcalc.HexToAddress("0xa")
	feeToken2 := ccipcalc.HexToAddress("0xb")

	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	testCases := []struct {
		name                     string
		commitObservations       []ccip.CommitObservation
		f                        int
		latestGasPrice           update
		latestTokenPrices        map[cciptypes.Address]update
		gasPriceHeartBeat        config.Duration
		daGasPriceDeviationPPB   int64
		execGasPriceDeviationPPB int64
		tokenPriceHeartBeat      config.Duration
		tokenPriceDeviationPPB   uint32
		expTokenUpdates          []cciptypes.TokenPrice
		expGasUpdates            []cciptypes.GasPrice
	}{
		{
			name: "median",
			commitObservations: []ccip.CommitObservation{
				{SourceGasPriceUSD: big.NewInt(1)},
				{SourceGasPriceUSD: big.NewInt(2)},
				{SourceGasPriceUSD: big.NewInt(3)},
				{SourceGasPriceUSD: big.NewInt(4)},
			},
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute), // recent
				value:     val1e18(9),                        // median deviates
			},
			f:             2,
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: big.NewInt(3)}},
		},
		{
			name: "gas price update skipped because the latest is similar and was updated recently",
			commitObservations: []ccip.CommitObservation{
				{SourceGasPriceUSD: val1e18(11)},
				{SourceGasPriceUSD: val1e18(12)},
			},
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   20e7,
			execGasPriceDeviationPPB: 20e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute), // recent
				value:     val1e18(10),                       // latest value close to the update
			},
			f:             1,
			expGasUpdates: nil,
		},
		{
			name: "gas price update included, the latest is similar but was not updated recently",
			commitObservations: []ccip.CommitObservation{
				{SourceGasPriceUSD: val1e18(10)},
				{SourceGasPriceUSD: val1e18(11)},
			},
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   20e7,
			execGasPriceDeviationPPB: 20e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-90 * time.Minute), // recent
				value:     val1e18(9),                        // latest value close to the update
			},
			f:             1,
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: val1e18(11)}},
		},
		{
			name: "gas price update deviates from latest",
			commitObservations: []ccip.CommitObservation{
				{SourceGasPriceUSD: val1e18(10)},
				{SourceGasPriceUSD: val1e18(20)},
				{SourceGasPriceUSD: val1e18(20)},
			},
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   20e7,
			execGasPriceDeviationPPB: 20e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute), // recent
				value:     val1e18(11),                       // latest value close to the update
			},
			f:             2,
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: val1e18(20)}},
		},
		{
			name: "median one token",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: big.NewInt(10)}, SourceGasPriceUSD: val1e18(0)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: big.NewInt(12)}, SourceGasPriceUSD: val1e18(0)},
			},
			f: 1,
			expTokenUpdates: []cciptypes.TokenPrice{
				{Token: feeToken1, Value: big.NewInt(12)},
			},
			// We expect a gas update because no latest
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: big.NewInt(0)}},
		},
		{
			name: "median two tokens",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: big.NewInt(10), feeToken2: big.NewInt(13)}, SourceGasPriceUSD: val1e18(0)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: big.NewInt(12), feeToken2: big.NewInt(7)}, SourceGasPriceUSD: val1e18(0)},
			},
			f: 1,
			expTokenUpdates: []cciptypes.TokenPrice{
				{Token: feeToken1, Value: big.NewInt(12)},
				{Token: feeToken2, Value: big.NewInt(13)},
			},
			// We expect a gas update because no latest
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: big.NewInt(0)}},
		},
		{
			name: "token price update skipped because it is close to the latest",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(11)}, SourceGasPriceUSD: val1e18(0)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(12)}, SourceGasPriceUSD: val1e18(0)},
			},
			f:                        1,
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   20e7,
			execGasPriceDeviationPPB: 20e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestTokenPrices: map[cciptypes.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-30 * time.Minute),
					value:     val1e18(10),
				},
			},
			// We expect a gas update because no latest
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: big.NewInt(0)}},
		},
		{
			name: "gas price and token price both included because they are not close to the latest",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(20)}, SourceGasPriceUSD: val1e18(10)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(21)}, SourceGasPriceUSD: val1e18(11)},
			},
			f:                        1,
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   10e7,
			execGasPriceDeviationPPB: 10e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute),
				value:     val1e18(9),
			},
			latestTokenPrices: map[cciptypes.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-30 * time.Minute),
					value:     val1e18(9),
				},
			},
			expTokenUpdates: []cciptypes.TokenPrice{
				{Token: feeToken1, Value: val1e18(21)},
			},
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: val1e18(11)}},
		},
		{
			name: "gas price and token price both included because they not been updated recently",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(20)}, SourceGasPriceUSD: val1e18(10)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(21)}, SourceGasPriceUSD: val1e18(11)},
			},
			f:                        1,
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   10e7,
			execGasPriceDeviationPPB: 10e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(2 * time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-90 * time.Minute),
				value:     val1e18(11),
			},
			latestTokenPrices: map[cciptypes.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-4 * time.Hour),
					value:     val1e18(21),
				},
			},
			expTokenUpdates: []cciptypes.TokenPrice{
				{Token: feeToken1, Value: val1e18(21)},
			},
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: val1e18(11)}},
		},
		{
			name: "gas price included because it deviates from latest and token price skipped because it does not deviate",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(20)}, SourceGasPriceUSD: val1e18(10)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(21)}, SourceGasPriceUSD: val1e18(11)},
			},
			f:                        1,
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   10e7,
			execGasPriceDeviationPPB: 10e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(2 * time.Hour),
			tokenPriceDeviationPPB:   200e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute),
				value:     val1e18(9),
			},
			latestTokenPrices: map[cciptypes.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-30 * time.Minute),
					value:     val1e18(9),
				},
			},
			expGasUpdates: []cciptypes.GasPrice{{DestChainSelector: defaultSourceChainSelector, Value: val1e18(11)}},
		},
		{
			name: "gas price skipped because it does not deviate and token price included because it has not been updated recently",
			commitObservations: []ccip.CommitObservation{
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(20)}, SourceGasPriceUSD: val1e18(10)},
				{TokenPricesUSD: map[cciptypes.Address]*big.Int{feeToken1: val1e18(21)}, SourceGasPriceUSD: val1e18(11)},
			},
			f:                        1,
			gasPriceHeartBeat:        *config.MustNewDuration(time.Hour),
			daGasPriceDeviationPPB:   10e7,
			execGasPriceDeviationPPB: 10e7,
			tokenPriceHeartBeat:      *config.MustNewDuration(2 * time.Hour),
			tokenPriceDeviationPPB:   20e7,
			latestGasPrice: update{
				timestamp: time.Now().Add(-30 * time.Minute),
				value:     val1e18(11),
			},
			latestTokenPrices: map[cciptypes.Address]update{
				feeToken1: {
					timestamp: time.Now().Add(-4 * time.Hour),
					value:     val1e18(21),
				},
			},
			expTokenUpdates: []cciptypes.TokenPrice{
				{Token: feeToken1, Value: val1e18(21)},
			},
			expGasUpdates: nil,
		},
	}

	evmEstimator := mocks.NewEvmFeeEstimator(t)
	evmEstimator.On("L1Oracle").Return(nil)
	estimatorCSVer, _ := semver.NewVersion("1.2.0")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			estimator, _ := prices.NewGasPriceEstimatorForCommitPlugin(
				*estimatorCSVer,
				evmEstimator,
				nil,
				tc.daGasPriceDeviationPPB,
				tc.execGasPriceDeviationPPB,
			)

			r := &CommitReportingPlugin{
				lggr:                logger.TestLogger(t),
				sourceChainSelector: defaultSourceChainSelector,
				offchainConfig: cciptypes.CommitOffchainConfig{
					GasPriceHeartBeat:      tc.gasPriceHeartBeat.Duration(),
					TokenPriceHeartBeat:    tc.tokenPriceHeartBeat.Duration(),
					TokenPriceDeviationPPB: tc.tokenPriceDeviationPPB,
				},
				gasPriceEstimator: estimator,
				F:                 tc.f,
			}

			var gasPriceObs []*big.Int
			tokenPriceObs := make(map[cciptypes.Address][]*big.Int)
			for _, obs := range tc.commitObservations {
				gasPriceObs = append(gasPriceObs, obs.SourceGasPriceUSD)
				for token, price := range obs.TokenPricesUSD {
					tokenPriceObs[token] = append(tokenPriceObs[token], price)
				}
			}

			gotGas, gotTokens, err := r.calculatePriceUpdates(gasPriceObs, tokenPriceObs, tc.latestGasPrice, tc.latestTokenPrices)

			assert.Equal(t, tc.expGasUpdates, gotGas)
			assert.Equal(t, tc.expTokenUpdates, gotTokens)
			assert.NoError(t, err)
		})
	}
}

func TestCommitReportingPlugin_generatePriceUpdates(t *testing.T) {
	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	const nTokens = 10
	tokens := make([]cciptypes.Address, nTokens)
	for i := range tokens {
		tokens[i] = cciptypes.Address(utils.RandomAddress().String())
	}
	sort.Slice(tokens, func(i, j int) bool { return tokens[i] < tokens[j] })

	testCases := []struct {
		name                 string
		tokenDecimals        map[cciptypes.Address]uint8
		sourceNativeToken    cciptypes.Address
		priceGetterRespData  map[cciptypes.Address]*big.Int
		priceGetterRespErr   error
		feeEstimatorRespFee  *big.Int
		feeEstimatorRespErr  error
		maxGasPrice          uint64
		expSourceGasPriceUSD *big.Int
		expTokenPricesUSD    map[cciptypes.Address]*big.Int
		expErr               bool
	}{
		{
			name: "base",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(10),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(1000),
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			expErr: false,
		},
		{
			name: "price getter returned an error",
			tokenDecimals: map[cciptypes.Address]uint8{
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
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name: "price getter skipped source native price",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[2],
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name: "base",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(10),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(1000),
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			expErr: false,
		},
		{
			name: "dynamic fee cap overrides legacy",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(20),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(2000),
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
			},
			expErr: false,
		},
		{
			name: "nil gas price",
			tokenDecimals: map[cciptypes.Address]uint8{
				tokens[0]: 18,
				tokens[1]: 18,
			},
			sourceNativeToken: tokens[0],
			priceGetterRespData: map[cciptypes.Address]*big.Int{
				tokens[0]: val1e18(100),
				tokens[1]: val1e18(200),
				tokens[2]: val1e18(300), // price getter returned a price for this token even though we didn't request it (should be skipped)
			},
			feeEstimatorRespFee: nil,
			maxGasPrice:         1e18,
			expErr:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceGetter := pricegetter.NewMockPriceGetter(t)
			defer priceGetter.AssertExpectations(t)

			gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
			defer gasPriceEstimator.AssertExpectations(t)

			var destTokens []cciptypes.Address
			for tk := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
			}
			sort.Slice(destTokens, func(i, j int) bool {
				return destTokens[i] < destTokens[j]
			})
			var destDecimals []uint8
			for _, token := range destTokens {
				destDecimals = append(destDecimals, tc.tokenDecimals[token])
			}

			queryTokens := ccipcommon.FlattenUniqueSlice([]cciptypes.Address{tc.sourceNativeToken}, destTokens)

			if len(queryTokens) > 0 {
				priceGetter.On("TokenPricesUSD", mock.Anything, queryTokens).Return(tc.priceGetterRespData, tc.priceGetterRespErr)
			}

			if tc.maxGasPrice > 0 {
				gasPriceEstimator.On("GetGasPrice", mock.Anything).Return(tc.feeEstimatorRespFee, tc.feeEstimatorRespErr)
				if tc.feeEstimatorRespFee != nil {
					pUSD := ccipcalc.CalculateUsdPerUnitGas(tc.feeEstimatorRespFee, tc.expTokenPricesUSD[tc.sourceNativeToken])
					gasPriceEstimator.On("DenoteInUSD", mock.Anything, mock.Anything).Return(pUSD, nil)
				}
			}

			p := &CommitReportingPlugin{
				sourceNative:      tc.sourceNativeToken,
				priceGetter:       priceGetter,
				gasPriceEstimator: gasPriceEstimator,
			}

			destPriceReg := ccipdatamocks.NewPriceRegistryReader(t)
			destPriceReg.On("GetTokensDecimals", mock.Anything, destTokens).Return(destDecimals, nil).Maybe()
			p.destPriceRegistryReader = destPriceReg

			sourceGasPriceUSD, tokenPricesUSD, err := p.generatePriceUpdates(context.Background(), logger.TestLogger(t), destTokens)
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

func TestCommitReportingPlugin_isStaleReport(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)
	merkleRoot1 := utils.Keccak256Fixed([]byte("some merkle root 1"))

	t.Run("empty report", func(t *testing.T) {
		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		r := &CommitReportingPlugin{commitStoreReader: commitStoreReader}
		isStale := r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{}, types.ReportTimestamp{})
		assert.True(t, isStale)
	})

	t.Run("merkle root", func(t *testing.T) {
		const expNextSeqNum = uint64(9)
		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(expNextSeqNum, nil)

		r := &CommitReportingPlugin{
			commitStoreReader: commitStoreReader,
		}

		assert.False(t, r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{
			MerkleRoot: merkleRoot1,
			Interval:   cciptypes.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
		}, types.ReportTimestamp{}))

		assert.True(t, r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{
			MerkleRoot: merkleRoot1}, types.ReportTimestamp{}))
	})
}

func TestCommitReportingPlugin_calculateMinMaxSequenceNumbers(t *testing.T) {
	testCases := []struct {
		name              string
		commitStoreSeqNum uint64
		msgSeqNums        []uint64

		expQueryMin uint64 // starting seq num that is used in the query to get messages
		expMin      uint64
		expMax      uint64
		expErr      bool
	}{
		{
			name:              "happy flow",
			commitStoreSeqNum: 9,
			msgSeqNums:        []uint64{11, 12, 13, 14},
			expQueryMin:       9,
			expMin:            11,
			expMax:            14,
			expErr:            false,
		},
		{
			name:              "happy flow 2",
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
			commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
			commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(tc.commitStoreSeqNum, nil)
			p.commitStoreReader = commitStoreReader

			onRampReader := ccipdatamocks.NewOnRampReader(t)
			var sendReqs []cciptypes.EVM2EVMMessageWithTxMeta
			for _, seqNum := range tc.msgSeqNums {
				sendReqs = append(sendReqs, cciptypes.EVM2EVMMessageWithTxMeta{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{
						SequenceNumber: seqNum,
					},
				})
			}
			onRampReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.expQueryMin, tc.expQueryMin+OnRampMessagesScanLimit, true).Return(sendReqs, nil)
			p.onRampReader = onRampReader

			minSeqNum, maxSeqNum, _, err := p.calculateMinMaxSequenceNumbers(ctx, lggr)
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
	chainSelector := uint64(1234)

	testCases := []struct {
		name                string
		destGasPriceUpdates []update
		expUpdate           update
		expErr              bool
	}{
		{
			name: "happy path",
			destGasPriceUpdates: []update{
				{timestamp: now, value: big.NewInt(1000)},
			},
			expUpdate: update{timestamp: now, value: big.NewInt(1000)},
			expErr:    false,
		},
		{
			name: "happy path two updates",
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
			p.sourceChainSelector = chainSelector
			p.lggr = lggr
			destPriceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			p.destPriceRegistryReader = destPriceRegistry

			if len(tc.destGasPriceUpdates) > 0 {
				var events []cciptypes.GasPriceUpdateWithTxMeta
				for _, u := range tc.destGasPriceUpdates {
					events = append(events, cciptypes.GasPriceUpdateWithTxMeta{
						GasPriceUpdate: cciptypes.GasPriceUpdate{
							GasPrice:         cciptypes.GasPrice{Value: u.value},
							TimestampUnixSec: big.NewInt(u.timestamp.Unix()),
						},
					})
				}
				destReader := ccipdatamocks.NewPriceRegistryReader(t)
				destReader.On("GetGasPriceUpdatesCreatedAfter", ctx, chainSelector, mock.Anything, 0).Return(events, nil)
				p.destPriceRegistryReader = destReader
			}

			priceUpdate, err := p.getLatestGasPriceUpdate(ctx, time.Now())
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
	tk1 := cciptypes.Address(utils.RandomAddress().String())
	tk2 := cciptypes.Address(utils.RandomAddress().String())

	testCases := []struct {
		name                 string
		priceRegistryUpdates []cciptypes.TokenPriceUpdate
		expUpdates           map[cciptypes.Address]update
		expErr               bool
	}{
		{
			name: "happy path",
			priceRegistryUpdates: []cciptypes.TokenPriceUpdate{
				{
					TokenPrice: cciptypes.TokenPrice{
						Token: tk1,
						Value: big.NewInt(1000),
					},
					TimestampUnixSec: big.NewInt(now.Add(1 * time.Minute).Unix()),
				},
				{
					TokenPrice: cciptypes.TokenPrice{
						Token: tk2,
						Value: big.NewInt(2000),
					},
					TimestampUnixSec: big.NewInt(now.Add(2 * time.Minute).Unix()),
				},
			},
			expUpdates: map[cciptypes.Address]update{
				tk1: {timestamp: now.Add(1 * time.Minute), value: big.NewInt(1000)},
				tk2: {timestamp: now.Add(2 * time.Minute), value: big.NewInt(2000)},
			},
			expErr: false,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &CommitReportingPlugin{}

			priceReg := ccipdatamocks.NewPriceRegistryReader(t)
			p.destPriceRegistryReader = priceReg

			var events []cciptypes.TokenPriceUpdateWithTxMeta
			for _, up := range tc.priceRegistryUpdates {
				events = append(events, cciptypes.TokenPriceUpdateWithTxMeta{
					TokenPriceUpdate: up,
				})
			}

			priceReg.On("GetTokenPriceUpdatesCreatedAfter", ctx, mock.Anything, 0).Return(events, nil)

			updates, err := p.getLatestTokenPriceUpdates(ctx, now)
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
		rep, err := encodeCommitReport(cciptypes.CommitStoreReport{
			MerkleRoot:  root32,
			Interval:    cciptypes.CommitStoreInterval{Min: min, Max: max},
			TokenPrices: []cciptypes.TokenPrice{},
			GasPrices: []cciptypes.GasPrice{
				{
					DestChainSelector: 1337,
					Value:             big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
				},
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
		intervals  []cciptypes.CommitStoreInterval
		rangeLimit uint64
		f          int
		wantMin    uint64
		wantMax    uint64
		wantErr    bool
	}{
		{"no obs", []cciptypes.CommitStoreInterval{{Min: 0, Max: 0}}, 0, 0, 0, 0, false},
		{"basic", []cciptypes.CommitStoreInterval{
			{Min: 9, Max: 14},
			{Min: 10, Max: 12},
			{Min: 10, Max: 14},
		}, 0, 1, 10, 14, false},
		{"min > max", []cciptypes.CommitStoreInterval{
			{Min: 9, Max: 4},
			{Min: 10, Max: 4},
			{Min: 10, Max: 6},
		}, 0, 1, 0, 0, true},
		{
			"range limit", []cciptypes.CommitStoreInterval{
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
			report := cciptypes.CommitStoreReport{
				TokenPrices: []cciptypes.TokenPrice{},
				GasPrices: []cciptypes.GasPrice{
					{
						DestChainSelector: uint64(1337),
						Value:             big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
					},
				},
				MerkleRoot: tree.Root(),
				Interval:   cciptypes.CommitStoreInterval{Min: tc.min, Max: tc.max},
			}
			out, err := encodeCommitReport(report)
			require.NoError(t, err)

			fn, err := factory.CommitReportToEthTxMeta(ccipconfig.CommitStore, *semver.MustParse("1.0.0"))
			require.NoError(t, err)
			txMeta, err := fn(out)
			require.NoError(t, err)
			require.NotNil(t, txMeta)
			require.EqualValues(t, tc.expectedRange, txMeta.SeqNumbers)
		})
	}
}

// TODO should be removed, tests need to be updated to use the Reader interface.
// encodeCommitReport is only used in tests
func encodeCommitReport(report cciptypes.CommitStoreReport) ([]byte, error) {
	commitStoreABI := abihelpers.MustParseABI(commit_store.CommitStoreABI)
	return v1_2_0.EncodeCommitReport(abihelpers.MustGetEventInputs(v1_0_0.ReportAccepted, commitStoreABI), report)
}
