package ccipcommit

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	mocks2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

func TestCommitReportingPlugin_Observation(t *testing.T) {
	sourceNativeTokenAddr := ccipcalc.HexToAddress("1000")
	someTokenAddr := ccipcalc.HexToAddress("2000")

	testCases := []struct {
		name                string
		epochAndRound       types.ReportTimestamp
		commitStoreIsPaused bool
		commitStoreSeqNum   uint64
		tokenPrices         map[cciptypes.Address]*big.Int
		sendReqs            []cciptypes.EVM2EVMMessageWithTxMeta
		tokenDecimals       map[cciptypes.Address]uint8
		fee                 *big.Int

		expErr bool
		expObs ccip.CommitObservation
	}{
		{
			name:              "base report",
			commitStoreSeqNum: 54,
			tokenPrices: map[cciptypes.Address]*big.Int{
				someTokenAddr:         big.NewInt(2),
				sourceNativeTokenAddr: big.NewInt(2),
			},
			sendReqs: []cciptypes.EVM2EVMMessageWithTxMeta{
				{EVM2EVMMessage: cciptypes.EVM2EVMMessage{SequenceNumber: 54}},
				{EVM2EVMMessage: cciptypes.EVM2EVMMessage{SequenceNumber: 55}},
			},
			fee: big.NewInt(100),
			tokenDecimals: map[cciptypes.Address]uint8{
				someTokenAddr: 8,
			},
			expObs: ccip.CommitObservation{
				TokenPricesUSD: map[cciptypes.Address]*big.Int{
					someTokenAddr: big.NewInt(20000000000),
				},
				SourceGasPriceUSD: big.NewInt(0),
				Interval: cciptypes.CommitStoreInterval{
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
			commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
			commitStoreReader.On("IsDown", ctx).Return(tc.commitStoreIsPaused, nil)
			if !tc.commitStoreIsPaused {
				commitStoreReader.On("GetExpectedNextSequenceNumber", ctx).Return(tc.commitStoreSeqNum, nil)
			}

			onRampReader := ccipdatamocks.NewOnRampReader(t)
			if len(tc.sendReqs) > 0 {
				onRampReader.On("GetSendRequestsBetweenSeqNums", ctx, tc.commitStoreSeqNum, tc.commitStoreSeqNum+OnRampMessagesScanLimit, true).
					Return(tc.sendReqs, nil)
			}

			priceGet := pricegetter.NewMockPriceGetter(t)
			if len(tc.tokenPrices) > 0 {
				addrs := []cciptypes.Address{sourceNativeTokenAddr}
				for addr := range tc.tokenDecimals {
					addrs = append(addrs, addr)
				}
				priceGet.On("TokenPricesUSD", mock.Anything, addrs).Return(tc.tokenPrices, nil)
			}

			gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
			if tc.fee != nil {
				var p = tc.fee
				var pUSD = ccipcalc.CalculateUsdPerUnitGas(p, tc.tokenPrices[sourceNativeTokenAddr])
				gasPriceEstimator.On("GetGasPrice", ctx).Return(p, nil)
				gasPriceEstimator.On("DenoteInUSD", p, tc.tokenPrices[sourceNativeTokenAddr]).Return(pUSD, nil)
			}

			destTokens := make([]cciptypes.Address, 0)
			destDecimals := make([]uint8, 0)
			for tk, d := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
				destDecimals = append(destDecimals, d)
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
			p.inflightReports = newInflightCommitReportsContainer(time.Hour)
			p.commitStoreReader = commitStoreReader
			p.onRampReader = onRampReader
			p.offRampReader = offRampReader
			p.destPriceRegistryReader = destPriceRegReader
			p.priceGetter = priceGet
			p.sourceNative = sourceNativeTokenAddr
			p.gasPriceEstimator = gasPriceEstimator
			p.metricsCollector = ccip.NoopMetricsCollector

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
	ctx := testutils.Context(t)
	sourceChainSelector := uint64(rand.Int())
	var gasPrice = big.NewInt(1)
	gasPriceHeartBeat := *config.MustNewDuration(time.Hour)

	t.Run("not enough observations", func(t *testing.T) {
		p := &CommitReportingPlugin{}
		p.lggr = logger.TestLogger(t)
		p.F = 1

		offRampReader := ccipdatamocks.NewOffRampReader(t)
		destPriceRegReader := ccipdatamocks.NewPriceRegistryReader(t)
		p.offRampReader = offRampReader
		p.destPriceRegistryReader = destPriceRegReader
		offRampReader.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{}, nil).Maybe()
		destPriceRegReader.On("GetFeeTokens", ctx).Return(nil, nil).Maybe()

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

			destTokens := make([]cciptypes.Address, 0)
			destDecimals := make([]uint8, 0)
			for tk, d := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
				destDecimals = append(destDecimals, d)
			}

			offRampReader := ccipdatamocks.NewOffRampReader(t)
			offRampReader.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{
				DestinationTokens: destTokens,
			}, nil).Maybe()

			destPriceRegistryReader.On("GetFeeTokens", ctx).Return(nil, nil).Maybe()
			destPriceRegistryReader.On("GetTokensDecimals", ctx, destTokens).Return(destDecimals, nil).Maybe()

			lp := mocks2.NewLogPoller(t)
			commitStoreReader, err := v1_2_0.NewCommitStore(logger.TestLogger(t), utils.RandomAddress(), nil, lp, nil)
			assert.NoError(t, err)

			p := &CommitReportingPlugin{}
			p.lggr = logger.TestLogger(t)
			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			p.destPriceRegistryReader = destPriceRegistryReader
			p.onRampReader = onRampReader
			p.sourceChainSelector = sourceChainSelector
			p.offRampReader = offRampReader
			p.gasPriceEstimator = gasPriceEstimator
			p.offchainConfig.GasPriceHeartBeat = gasPriceHeartBeat.Duration()
			p.commitStoreReader = commitStoreReader
			p.F = tc.f
			p.metricsCollector = ccip.NoopMetricsCollector

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
		p.inflightReports = newInflightCommitReportsContainer(time.Minute)
		p.metricsCollector = ccip.NoopMetricsCollector
		return p
	}

	t.Run("report cannot be decoded leads to error", func(t *testing.T) {
		p := newPlugin()

		encodedReport := []byte("whatever")

		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p.commitStoreReader = commitStoreReader
		commitStoreReader.On("DecodeCommitReport", encodedReport).
			Return(cciptypes.CommitStoreReport{}, errors.New("unable to decode report"))

		_, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.Error(t, err)
	})

	t.Run("empty report should not be accepted", func(t *testing.T) {
		p := newPlugin()

		report := cciptypes.CommitStoreReport{}

		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p.commitStoreReader = commitStoreReader
		commitStoreReader.On("DecodeCommitReport", mock.Anything).Return(report, nil)

		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)
		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldAccept)
	})

	t.Run("stale report should not be accepted", func(t *testing.T) {
		onChainSeqNum := uint64(100)

		//_, _ := testhelpers.NewFakeCommitStore(t, onChainSeqNum)

		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		p := newPlugin()

		p.commitStoreReader = commitStoreReader

		report := cciptypes.CommitStoreReport{
			GasPrices:  []cciptypes.GasPrice{{Value: big.NewInt(int64(rand.Int()))}},
			MerkleRoot: [32]byte{123}, // this report is considered non-empty since it has a merkle root
		}

		commitStoreReader.On("DecodeCommitReport", mock.Anything).Return(report, nil)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(onChainSeqNum, nil)

		// stale since report interval is behind on chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum - 2, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)

		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldAccept)
	})

	t.Run("non-stale report should be accepted and added inflight", func(t *testing.T) {
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
		commitStoreReader.On("DecodeCommitReport", mock.Anything).Return(report, nil)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(onChainSeqNum, nil)

		// non-stale since report interval is not behind on-chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)

		shouldAccept, err := p.ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.True(t, shouldAccept)

		// make sure that the report was added inflight
		tokenPriceUpdates := p.inflightReports.latestInflightTokenPriceUpdates()
		priceUpdate := tokenPriceUpdates[report.TokenPrices[0].Token]
		assert.Equal(t, report.TokenPrices[0].Value.Uint64(), priceUpdate.value.Uint64())
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
	p.inflightReports = newInflightCommitReportsContainer(time.Minute)
	p.lggr = logger.TestLogger(t)

	t.Run("should transmit when report is not stale", func(t *testing.T) {
		// not-stale since report interval is not behind on chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)
		commitStoreReader.On("DecodeCommitReport", encodedReport).Return(report, nil).Once()
		shouldTransmit, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.True(t, shouldTransmit)
	})

	t.Run("should not transmit when report is stale", func(t *testing.T) {
		// stale since report interval is behind on chain seq num
		report.Interval = cciptypes.CommitStoreInterval{Min: onChainSeqNum - 2, Max: onChainSeqNum + 10}
		encodedReport, err := encodeCommitReport(report)
		assert.NoError(t, err)
		commitStoreReader.On("DecodeCommitReport", encodedReport).Return(report, nil).Once()
		shouldTransmit, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, encodedReport)
		assert.NoError(t, err)
		assert.False(t, shouldTransmit)
	})

	t.Run("error when report cannot be decoded", func(t *testing.T) {
		reportBytes := []byte("whatever")
		commitStoreReader.On("DecodeCommitReport", reportBytes).
			Return(cciptypes.CommitStoreReport{}, errors.New("decode error")).Once()
		_, err := p.ShouldTransmitAcceptedReport(ctx, types.ReportTimestamp{}, reportBytes)
		assert.Error(t, err)
	})
}

func TestCommitReportingPlugin_validateObservations(t *testing.T) {
	ctx := context.Background()

	token1 := ccipcalc.HexToAddress("0xa")
	token2 := ccipcalc.HexToAddress("0xb")
	token1Price := big.NewInt(1)
	token2Price := big.NewInt(2)
	unsupportedToken := ccipcalc.HexToAddress("0xc")
	gasPrice := big.NewInt(100)

	tokenDecimals := make(map[cciptypes.Address]uint8)
	tokenDecimals[token1] = 18
	tokenDecimals[token2] = 18
	destTokens := []cciptypes.Address{token1, token2}

	ob1 := ccip.CommitObservation{
		Interval: cciptypes.CommitStoreInterval{Min: 0, Max: 0},
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1: token1Price,
			token2: token2Price,
		},
		SourceGasPriceUSD: gasPrice,
	}
	ob1Bytes, err := ob1.Marshal()
	assert.NoError(t, err)
	var ob2, ob3 ccip.CommitObservation
	_ = json.Unmarshal(ob1Bytes, &ob2)
	_ = json.Unmarshal(ob1Bytes, &ob3)

	obWithNilGasPrice := ccip.CommitObservation{
		Interval: cciptypes.CommitStoreInterval{Min: 0, Max: 0},
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1: token1Price,
			token2: token2Price,
		},
		SourceGasPriceUSD: nil,
	}
	obWithNilTokenPrice := ccip.CommitObservation{
		Interval: cciptypes.CommitStoreInterval{Min: 0, Max: 0},
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1: token1Price,
			token2: nil,
		},
		SourceGasPriceUSD: gasPrice,
	}
	obMissingTokenPrices := ccip.CommitObservation{
		Interval:          cciptypes.CommitStoreInterval{Min: 0, Max: 0},
		TokenPricesUSD:    map[cciptypes.Address]*big.Int{},
		SourceGasPriceUSD: gasPrice,
	}
	obWithUnsupportedToken := ccip.CommitObservation{
		Interval: cciptypes.CommitStoreInterval{Min: 0, Max: 0},
		TokenPricesUSD: map[cciptypes.Address]*big.Int{
			token1:           token1Price,
			token2:           token2Price,
			unsupportedToken: token2Price,
		},
		SourceGasPriceUSD: gasPrice,
	}
	obEmpty := ccip.CommitObservation{
		Interval:          cciptypes.CommitStoreInterval{Min: 0, Max: 0},
		TokenPricesUSD:    nil,
		SourceGasPriceUSD: nil,
	}

	testCases := []struct {
		name               string
		commitObservations []ccip.CommitObservation
		f                  int
		expValidObs        []ccip.CommitObservation
		expError           bool
	}{
		{
			name:               "base",
			commitObservations: []ccip.CommitObservation{ob1, ob2},
			f:                  1,
			expValidObs:        []ccip.CommitObservation{ob1, ob2},
			expError:           false,
		},
		{
			name:               "pass with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, ob3},
			f:                  2,
			expValidObs:        []ccip.CommitObservation{ob1, ob2, ob3},
			expError:           false,
		},
		{
			name:               "tolerate 1 nil gas price with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, ob3, obWithNilGasPrice},
			f:                  2,
			expValidObs:        []ccip.CommitObservation{ob1, ob2, ob3},
			expError:           false,
		},
		{
			name:               "tolerate 1 nil token price with f=1",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obWithNilTokenPrice},
			f:                  1,
			expValidObs:        []ccip.CommitObservation{ob1, ob2},
			expError:           false,
		},
		{
			name:               "tolerate 1 missing token prices with f=1",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obMissingTokenPrices},
			f:                  1,
			expValidObs:        []ccip.CommitObservation{ob1, ob2},
			expError:           false,
		},
		{
			name:               "tolerate 1 unsupported token with f=1",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obWithUnsupportedToken},
			f:                  1,
			expValidObs:        []ccip.CommitObservation{ob1, ob2},
			expError:           false,
		},
		{
			name:               "not enough valid observations",
			commitObservations: []ccip.CommitObservation{ob1, ob2},
			f:                  2,
			expValidObs:        nil,
			expError:           true,
		},
		{
			name:               "too many faulty observations with f=2",
			commitObservations: []ccip.CommitObservation{ob1, ob2, obMissingTokenPrices, obWithUnsupportedToken},
			f:                  2,
			expValidObs:        nil,
			expError:           true,
		},
		{
			name:               "too many faulty observations with f=1",
			commitObservations: []ccip.CommitObservation{ob1, obEmpty},
			f:                  1,
			expValidObs:        nil,
			expError:           true,
		},
		{
			name:               "all faulty observations",
			commitObservations: []ccip.CommitObservation{obWithNilGasPrice, obWithNilTokenPrice, obMissingTokenPrices, obWithUnsupportedToken, obEmpty},
			f:                  1,
			expValidObs:        nil,
			expError:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obs, err := validateObservations(ctx, logger.TestLogger(t), destTokens, tc.f, tc.commitObservations)

			if tc.expError {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.expValidObs, obs)
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
			gotTokens, gotGas, err := r.calculatePriceUpdates(tc.commitObservations, tc.latestGasPrice, tc.latestTokenPrices)

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

			tokens := make([]cciptypes.Address, 0, len(tc.tokenDecimals))
			for tk := range tc.tokenDecimals {
				tokens = append(tokens, tk)
			}
			tokens = ccipcommon.FlattenUniqueSlice(tokens, []cciptypes.Address{tc.sourceNativeToken})
			sort.Slice(tokens, func(i, j int) bool { return tokens[i] < tokens[j] })

			if len(tokens) > 0 {
				priceGetter.On("TokenPricesUSD", mock.Anything, tokens).Return(tc.priceGetterRespData, tc.priceGetterRespErr)
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

			destTokens := make([]cciptypes.Address, 0, len(tc.tokenDecimals))
			destDecimals := make([]uint8, 0, len(tc.tokenDecimals))
			for tk, d := range tc.tokenDecimals {
				destTokens = append(destTokens, tk)
				destDecimals = append(destDecimals, d)
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

func TestCommitReportingPlugin_nextMinSeqNum(t *testing.T) {
	lggr := logger.TestLogger(t)
	root1 := utils.Keccak256Fixed(hexutil.MustDecode("0xaa"))

	var tt = []struct {
		onChainMin          uint64
		inflight            []cciptypes.CommitStoreReport
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
			inflight: []cciptypes.CommitStoreReport{
				{Interval: cciptypes.CommitStoreInterval{Min: uint64(1), Max: uint64(2)}, MerkleRoot: root1}},
			expectedInflightMin: uint64(3),
			expectedOnChainMin:  uint64(1),
		},
		{
			onChainMin: uint64(1),
			inflight: []cciptypes.CommitStoreReport{
				{Interval: cciptypes.CommitStoreInterval{Min: uint64(3), Max: uint64(4)}, MerkleRoot: root1}},
			expectedInflightMin: uint64(5),
			expectedOnChainMin:  uint64(1),
		},
		{
			onChainMin: uint64(1),
			inflight: []cciptypes.CommitStoreReport{
				{Interval: cciptypes.CommitStoreInterval{Min: uint64(1), Max: uint64(MaxInflightSeqNumGap + 2)}, MerkleRoot: root1}},
			expectedInflightMin: uint64(1),
			expectedOnChainMin:  uint64(1),
		},
	}
	for _, tc := range tt {
		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(tc.onChainMin, nil).Maybe()
		cp := CommitReportingPlugin{commitStoreReader: commitStoreReader, inflightReports: newInflightCommitReportsContainer(time.Hour)}
		epochAndRound := uint64(1)
		for _, rep := range tc.inflight {
			rc := rep
			rc.GasPrices = []cciptypes.GasPrice{{}}
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
		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		r := &CommitReportingPlugin{commitStoreReader: commitStoreReader}
		isStale := r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{}, false, types.ReportTimestamp{})
		assert.True(t, isStale)
	})

	t.Run("merkle root", func(t *testing.T) {
		const expNextSeqNum = uint64(9)
		commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
		commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(expNextSeqNum, nil)

		r := &CommitReportingPlugin{
			commitStoreReader: commitStoreReader,
			inflightReports: &inflightCommitReportsContainer{
				inFlight: map[[32]byte]InflightCommitReport{
					merkleRoot2: {
						report: cciptypes.CommitStoreReport{
							Interval: cciptypes.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
						},
					},
				},
			},
		}

		assert.False(t, r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{
			MerkleRoot: merkleRoot1,
			Interval:   cciptypes.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
		}, false, types.ReportTimestamp{}))

		assert.True(t, r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{
			MerkleRoot: merkleRoot1,
			Interval:   cciptypes.CommitStoreInterval{Min: expNextSeqNum + 1, Max: expNextSeqNum + 10},
		}, true, types.ReportTimestamp{}))

		assert.True(t, r.isStaleReport(ctx, lggr, cciptypes.CommitStoreReport{
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
			commitStoreReader := ccipdatamocks.NewCommitStoreReader(t)
			commitStoreReader.On("GetExpectedNextSequenceNumber", mock.Anything).Return(tc.commitStoreSeqNum, nil)
			p.commitStoreReader = commitStoreReader

			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			if tc.inflightSeqNum > 0 {
				p.inflightReports.inFlight[[32]byte{}] = InflightCommitReport{
					report: cciptypes.CommitStoreReport{
						Interval: cciptypes.CommitStoreInterval{
							Min: tc.inflightSeqNum,
							Max: tc.inflightSeqNum,
						},
					},
				}
			}

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
			inflightGasPriceUpdate: nil,
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
			p.sourceChainSelector = chainSelector
			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			p.lggr = lggr
			destPriceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			p.destPriceRegistryReader = destPriceRegistry

			if tc.inflightGasPriceUpdate != nil {
				p.inflightReports.inFlightPriceUpdates = append(
					p.inflightReports.inFlightPriceUpdates,
					InflightPriceUpdate{
						createdAt: tc.inflightGasPriceUpdate.timestamp,
						gasPrices: []cciptypes.GasPrice{{
							DestChainSelector: chainSelector,
							Value:             tc.inflightGasPriceUpdate.value,
						}},
					},
				)
			}

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
	tk1 := cciptypes.Address(utils.RandomAddress().String())
	tk2 := cciptypes.Address(utils.RandomAddress().String())

	testCases := []struct {
		name                 string
		priceRegistryUpdates []cciptypes.TokenPriceUpdate
		checkInflight        bool
		inflightUpdates      map[cciptypes.Address]update
		expUpdates           map[cciptypes.Address]update
		expErr               bool
	}{
		{
			name: "ignore inflight updates",
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
			checkInflight: false,
			expUpdates: map[cciptypes.Address]update{
				tk1: {timestamp: now.Add(1 * time.Minute), value: big.NewInt(1000)},
				tk2: {timestamp: now.Add(2 * time.Minute), value: big.NewInt(2000)},
			},
			expErr: false,
		},
		{
			name: "consider inflight updates",
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
			checkInflight: true,
			inflightUpdates: map[cciptypes.Address]update{
				tk1: {timestamp: now, value: big.NewInt(500)}, // inflight but older
				tk2: {timestamp: now.Add(4 * time.Minute), value: big.NewInt(4000)},
			},
			expUpdates: map[cciptypes.Address]update{
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

			//_, priceRegAddr := testhelpers.NewFakePriceRegistry(t)
			priceReg := ccipdatamocks.NewPriceRegistryReader(t)
			p.destPriceRegistryReader = priceReg

			//destReader := ccipdata.NewMockReader(t)
			var events []cciptypes.TokenPriceUpdateWithTxMeta
			for _, up := range tc.priceRegistryUpdates {
				events = append(events, cciptypes.TokenPriceUpdateWithTxMeta{
					TokenPriceUpdate: up,
				})
			}
			//destReader.On("GetTokenPriceUpdatesCreatedAfter", ctx, priceRegAddr, mock.Anything, 0).Return(events, nil)
			priceReg.On("GetTokenPriceUpdatesCreatedAfter", ctx, mock.Anything, 0).Return(events, nil)

			p.inflightReports = newInflightCommitReportsContainer(time.Minute)
			if len(tc.inflightUpdates) > 0 {
				for tk, upd := range tc.inflightUpdates {
					p.inflightReports.inFlightPriceUpdates = append(p.inflightReports.inFlightPriceUpdates, InflightPriceUpdate{
						createdAt: upd.timestamp,
						tokenPrices: []cciptypes.TokenPrice{
							{Token: tk, Value: upd.value},
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
