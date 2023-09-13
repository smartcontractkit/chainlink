package ccip

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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipevents"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
	plugintesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/plugins"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var defaultGasPrice = big.NewInt(3e9)

type commitTestHarness = struct {
	plugintesthelpers.CCIPPluginTestHarness
	plugin       *CommitReportingPlugin
	mockedGetFee *mock.Call
}

func setupCommitTestHarness(t *testing.T) commitTestHarness {
	th := plugintesthelpers.SetupCCIPTestHarness(t)

	sourceFeeEstimator := mocks.NewEvmFeeEstimator(t)

	mockedGetFee := sourceFeeEstimator.On(
		"GetFee",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Maybe().Return(gas.EvmFee{Legacy: assets.NewWei(defaultGasPrice)}, uint32(200e3), nil)

	lggr := logger.TestLogger(t)
	priceGetter := newMockPriceGetter()

	backendClient := client.NewSimulatedBackendClient(t, th.Dest.Chain, new(big.Int).SetUint64(th.Dest.ChainID))
	plugin := CommitReportingPlugin{
		config: CommitPluginConfig{
			lggr:                th.Lggr,
			sourceLP:            th.SourceLP,
			destLP:              th.DestLP,
			sourceEvents:        ccipevents.NewLogPollerClient(th.SourceLP, lggr, backendClient),
			destEvents:          ccipevents.NewLogPollerClient(th.DestLP, lggr, backendClient),
			offRamp:             th.Dest.OffRamp,
			onRampAddress:       th.Source.OnRamp.Address(),
			commitStore:         th.Dest.CommitStore,
			priceGetter:         priceGetter,
			sourceNative:        utils.RandomAddress(),
			sourceFeeEstimator:  sourceFeeEstimator,
			sourceChainSelector: th.Source.ChainSelector,
			destClient:          backendClient,
			sourceClient:        backendClient,
			leafHasher:          hashlib.NewLeafHasher(th.Source.ChainSelector, th.Dest.ChainSelector, th.Source.OnRamp.Address(), hashlib.NewKeccakCtx()),
		},
		inflightReports: newInflightCommitReportsContainer(time.Hour),
		onchainConfig:   th.CommitOnchainConfig,
		offchainConfig: ccipconfig.CommitOffchainConfig{
			SourceFinalityDepth:   0,
			DestFinalityDepth:     0,
			FeeUpdateDeviationPPB: 5e7,
			FeeUpdateHeartBeat:    models.MustMakeDuration(12 * time.Hour),
			MaxGasPrice:           200e9,
		},
		lggr:               th.Lggr,
		destPriceRegistry:  th.Dest.PriceRegistry,
		tokenDecimalsCache: cache.NewTokenToDecimals(th.Lggr, th.DestLP, th.Dest.OffRamp, th.Dest.PriceRegistry, backendClient, 0),
	}

	priceGetter.On("TokenPricesUSD", mock.Anything, mock.Anything).Return(map[common.Address]*big.Int{
		plugin.config.sourceNative:    big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1e18)),
		th.Source.LinkToken.Address(): big.NewInt(0).Mul(big.NewInt(200), big.NewInt(1e18)),
		th.Dest.LinkToken.Address():   big.NewInt(0).Mul(big.NewInt(200), big.NewInt(1e18)),
	}, nil)

	return commitTestHarness{
		CCIPPluginTestHarness: th,
		plugin:                &plugin,
		mockedGetFee:          mockedGetFee,
	}
}

func TestCommitReportSize(t *testing.T) {
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

func TestCommitReportEncoding(t *testing.T) {
	th := plugintesthelpers.SetupCCIPTestHarness(t)
	newTokenPrice := big.NewInt(9e18) // $9
	newGasPrice := big.NewInt(2000e9) // $2000 per eth * 1gwei

	// Send a report.
	mctx := hashlib.NewKeccakCtx()
	tree, err := merklemulti.NewTree(mctx, [][32]byte{mctx.Hash([]byte{0xaa})})
	require.NoError(t, err)
	report := commit_store.CommitStoreCommitReport{
		PriceUpdates: commit_store.InternalPriceUpdates{
			TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{
				{
					SourceToken: th.Dest.LinkToken.Address(),
					UsdPerToken: newTokenPrice,
				},
			},
			DestChainSelector: th.Source.ChainSelector,
			UsdPerUnitGas:     newGasPrice,
		},
		MerkleRoot: tree.Root(),
		Interval:   commit_store.CommitStoreInterval{Min: 1, Max: 10},
	}
	out, err := abihelpers.EncodeCommitReport(report)
	require.NoError(t, err)
	decodedReport, err := abihelpers.DecodeCommitReport(out)
	require.NoError(t, err)
	require.Equal(t, report, decodedReport)

	latestEpocAndRound, err := th.Dest.CommitStoreHelper.GetLatestPriceEpochAndRound(nil)
	require.NoError(t, err)

	tx, err := th.Dest.CommitStoreHelper.Report(th.Dest.User, out, big.NewInt(int64(latestEpocAndRound+1)))
	require.NoError(t, err)
	th.CommitAndPollLogs(t)
	res, err := th.Dest.Chain.TransactionReceipt(testutils.Context(t), tx.Hash())
	require.NoError(t, err)
	assert.Equal(t, uint64(1), res.Status)

	// Ensure root exists.
	ts, err := th.Dest.CommitStore.GetMerkleRoot(nil, tree.Root())
	require.NoError(t, err)
	require.NotEqual(t, ts.String(), "0")

	// Ensure price update went through
	destChainGasPrice, err := th.Dest.PriceRegistry.GetDestinationChainGasPrice(nil, th.Source.ChainSelector)
	require.NoError(t, err)
	assert.Equal(t, newGasPrice, destChainGasPrice.Value)

	linkTokenPrice, err := th.Dest.PriceRegistry.GetTokenPrice(nil, th.Dest.LinkToken.Address())
	require.NoError(t, err)
	assert.Equal(t, newTokenPrice, linkTokenPrice.Value)
}

func TestCommitObservation(t *testing.T) {
	th := setupCommitTestHarness(t)
	th.plugin.F = 1

	mb := th.GenerateAndSendMessageBatch(t, 1, 0, 0)

	tests := []struct {
		name            string
		commitStoreDown bool
		expected        *CommitObservation
		expectedError   bool
	}{
		{
			"base",
			false,
			&CommitObservation{
				Interval:          mb.Interval,
				SourceGasPriceUSD: new(big.Int).Mul(defaultGasPrice, big.NewInt(100)),
				TokenPricesUSD: map[common.Address]*big.Int{
					th.Dest.LinkToken.Address(): new(big.Int).Mul(big.NewInt(200), big.NewInt(1e18)),
				},
			},
			false,
		},
		{
			"commitStore down",
			true,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.commitStoreDown && !isCommitStoreDownNow(testutils.Context(t), th.Lggr, th.Dest.CommitStore) {
				_, err := th.Dest.CommitStore.Pause(th.Dest.User)
				require.NoError(t, err)
				th.CommitAndPollLogs(t)
			} else if !tt.commitStoreDown && isCommitStoreDownNow(testutils.Context(t), th.Lggr, th.Dest.CommitStore) {
				_, err := th.Dest.CommitStore.Unpause(th.Dest.User)
				require.NoError(t, err)
				th.CommitAndPollLogs(t)
			}

			gotObs, err := th.plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, types.Query{})

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			var decodedObservation *CommitObservation
			if gotObs != nil {
				decodedObservation = new(CommitObservation)
				err = json.Unmarshal(gotObs, decodedObservation)
				require.NoError(t, err)

			}
			assert.Equal(t, tt.expected, decodedObservation)
		})
	}
}

func TestCommitReport(t *testing.T) {
	th := setupCommitTestHarness(t)
	th.plugin.F = 1

	mb := th.GenerateAndSendMessageBatch(t, 1, 0, 0)

	tests := []struct {
		name          string
		observations  []CommitObservation
		shouldReport  bool
		commitReport  *commit_store.CommitStoreCommitReport
		expectedError bool
	}{
		{
			"base",
			[]CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 1, Max: 1}},
				{Interval: commit_store.CommitStoreInterval{Min: 1, Max: 1}},
			},
			true,
			&commit_store.CommitStoreCommitReport{
				MerkleRoot: mb.Root,
				Interval:   commit_store.CommitStoreInterval{Min: 1, Max: 1},
				PriceUpdates: commit_store.InternalPriceUpdates{
					TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
					DestChainSelector: 0,
					UsdPerUnitGas:     new(big.Int),
				},
			},
			false,
		},
		{
			"not enough observations",
			[]CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 1, Max: 1}},
			},
			false,
			nil,
			true,
		},
		{
			"empty",
			[]CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 0, Max: 0}},
				{Interval: commit_store.CommitStoreInterval{Min: 0, Max: 0}},
			},
			false,
			nil,
			false,
		},
		{
			"no leaves",
			[]CommitObservation{
				{Interval: commit_store.CommitStoreInterval{Min: 2, Max: 2}},
				{Interval: commit_store.CommitStoreInterval{Min: 2, Max: 2}},
			},
			false,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aos := make([]types.AttributedObservation, 0, len(tt.observations))
			for _, o := range tt.observations {
				obs, err := o.Marshal()
				require.NoError(t, err)
				aos = append(aos, types.AttributedObservation{Observation: obs})
			}
			gotShouldReport, gotReport, err := th.plugin.Report(testutils.Context(t), types.ReportTimestamp{}, types.Query{}, aos)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.shouldReport, gotShouldReport)

			var expectedReport types.Report
			if tt.commitReport != nil {
				expectedReport, err = abihelpers.EncodeCommitReport(*tt.commitReport)
				require.NoError(t, err)
			}
			assert.Equal(t, expectedReport, gotReport)
		})
	}
}

func TestCalculatePriceUpdates(t *testing.T) {
	t.Parallel()

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

func TestCalculateIntervalConsensus(t *testing.T) {
	t.Parallel()

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

func TestGeneratePriceUpdates(t *testing.T) {
	t.Parallel()

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
			priceGetter := newMockPriceGetter()
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
				priceGetter.On("TokenPricesUSD", tokens).Return(tc.priceGetterRespData, tc.priceGetterRespErr)
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

func TestUpdateTokenToDecimalMapping(t *testing.T) {
	th := plugintesthelpers.SetupCCIPTestHarness(t)

	destToken, _, _, err := link_token_interface.DeployLinkToken(th.Dest.User, th.Dest.Chain)
	require.NoError(t, err)

	feeToken, _, _, err := link_token_interface.DeployLinkToken(th.Dest.User, th.Dest.Chain)
	require.NoError(t, err)
	th.CommitAndPollLogs(t)

	tokens := []common.Address{}
	tokens = append(tokens, destToken)
	tokens = append(tokens, feeToken)

	mockOffRamp := &mock_contracts.EVM2EVMOffRampInterface{}
	mockOffRamp.On("GetDestinationTokens", mock.Anything).Return([]common.Address{destToken}, nil)
	mockOffRamp.On("Address").Return(common.Address{})

	mockPriceRegistry := &mock_contracts.PriceRegistryInterface{}
	mockPriceRegistry.On("GetFeeTokens", mock.Anything).Return([]common.Address{feeToken}, nil)
	mockPriceRegistry.On("Address").Return(common.Address{})

	backendClient := client.NewSimulatedBackendClient(t, th.Dest.Chain, new(big.Int).SetUint64(th.Dest.ChainID))
	plugin := CommitReportingPlugin{
		config: CommitPluginConfig{
			offRamp:    mockOffRamp,
			destClient: backendClient,
		},
		destPriceRegistry:  mockPriceRegistry,
		tokenDecimalsCache: cache.NewTokenToDecimals(th.Lggr, th.DestLP, mockOffRamp, mockPriceRegistry, backendClient, 0),
	}

	tokenMapping, err := plugin.tokenDecimalsCache.Get(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, len(tokens), len(tokenMapping))
	assert.Equal(t, uint8(18), tokenMapping[destToken])
	assert.Equal(t, uint8(18), tokenMapping[feeToken])
}

func TestCalculateUsdPer1e18TokenAmount(t *testing.T) {
	t.Parallel()

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

func TestShouldTransmitAcceptedReport(t *testing.T) {
	th := setupCommitTestHarness(t)
	tokenPrice := big.NewInt(9e18) // $9
	gasPrice := big.NewInt(1500e9) // $1500 per eth * 1gwei

	nextMinSeqNr := uint64(10)
	_, err := th.Dest.CommitStore.SetMinSeqNr(th.Dest.User, nextMinSeqNr)
	require.NoError(t, err)
	_, err = th.Dest.PriceRegistry.UpdatePrices(th.Dest.User, price_registry.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry.InternalTokenPriceUpdate{
			{SourceToken: th.Dest.LinkToken.Address(), UsdPerToken: tokenPrice},
		},
		DestChainSelector: th.Source.ChainSelector,
		UsdPerUnitGas:     gasPrice,
	})
	require.NoError(t, err)
	th.CommitAndPollLogs(t)
	round := uint8(1)

	tests := []struct {
		name       string
		seq        uint64
		gasPrice   *big.Int
		tokenPrice *big.Int
		expected   bool
	}{
		{"base", nextMinSeqNr, nil, nil, true},
		{"future", nextMinSeqNr + 10, nil, nil, true},
		{"empty", 0, nil, nil, false},
		{"gasPrice update", 0, big.NewInt(10), nil, true},
		{"gasPrice stale", 0, gasPrice, nil, false},
		{"tokenPrice update", 0, nil, big.NewInt(20), true},
		{"tokenPrice stale", 0, nil, tokenPrice, false},
		{"token price and gas price stale", 0, gasPrice, tokenPrice, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var destChainSelector uint64
			gasPrice := new(big.Int)
			if tt.gasPrice != nil {
				destChainSelector = th.Source.ChainSelector
				gasPrice = tt.gasPrice
			}

			var tokenPrices []commit_store.InternalTokenPriceUpdate
			if tt.tokenPrice != nil {
				tokenPrices = []commit_store.InternalTokenPriceUpdate{
					{SourceToken: th.Dest.LinkToken.Address(), UsdPerToken: tt.tokenPrice},
				}
			} else {
				tokenPrices = []commit_store.InternalTokenPriceUpdate{}
			}

			var root [32]byte
			if tt.seq > 0 {
				root = testutils.Random32Byte()
			}

			report, err := abihelpers.EncodeCommitReport(commit_store.CommitStoreCommitReport{
				PriceUpdates: commit_store.InternalPriceUpdates{
					TokenPriceUpdates: tokenPrices,
					DestChainSelector: destChainSelector,
					UsdPerUnitGas:     gasPrice,
				},
				MerkleRoot: root,
				Interval:   commit_store.CommitStoreInterval{Min: tt.seq, Max: tt.seq},
			})
			require.NoError(t, err)

			got, err := th.plugin.ShouldTransmitAcceptedReport(testutils.Context(t), types.ReportTimestamp{Epoch: 1, Round: round}, report)
			round++
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestShouldAcceptFinalizedReport(t *testing.T) {
	nextMinSeqNr := uint64(10)

	tests := []struct {
		name                     string
		seq                      uint64
		latestPriceEpochAndRound int64
		epoch                    uint32
		round                    uint8
		destChainSelector        int
		skipRoot                 bool
		expected                 bool
		err                      bool
	}{
		{
			name:  "future",
			seq:   nextMinSeqNr * 2,
			epoch: 1,
			round: 1,
		},
		{
			name:  "empty",
			epoch: 1,
			round: 2,
		},
		{
			name:  "stale",
			seq:   nextMinSeqNr - 1,
			epoch: 1,
			round: 3,
		},
		{
			name:     "base",
			seq:      nextMinSeqNr,
			epoch:    1,
			round:    4,
			expected: true,
		},
		{
			name:                     "price update - epoch and round is ok",
			seq:                      nextMinSeqNr,
			latestPriceEpochAndRound: int64(mergeEpochAndRound(2, 10)),
			epoch:                    2,
			round:                    11,
			destChainSelector:        rand.Int(),
			skipRoot:                 true,
			expected:                 true,
		},
		{
			name:                     "price update - epoch and round is behind",
			seq:                      nextMinSeqNr,
			latestPriceEpochAndRound: int64(mergeEpochAndRound(2, 10)),
			epoch:                    2,
			round:                    9,
			destChainSelector:        rand.Int(),
			skipRoot:                 true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := setupCommitTestHarness(t)
			_, err := th.Dest.CommitStore.SetMinSeqNr(th.Dest.User, nextMinSeqNr)
			require.NoError(t, err)

			_, err = th.Dest.CommitStoreHelper.SetLatestPriceEpochAndRound(th.Dest.User, big.NewInt(tt.latestPriceEpochAndRound))
			require.NoError(t, err)

			th.CommitAndPollLogs(t)

			var root [32]byte
			if tt.seq > 0 {
				root = testutils.Random32Byte()
			}

			r := commit_store.CommitStoreCommitReport{
				PriceUpdates: commit_store.InternalPriceUpdates{
					TokenPriceUpdates: []commit_store.InternalTokenPriceUpdate{},
					DestChainSelector: uint64(tt.destChainSelector),
					UsdPerUnitGas:     new(big.Int),
				},
				Interval: commit_store.CommitStoreInterval{Min: tt.seq, Max: tt.seq},
			}
			if !tt.skipRoot {
				r.MerkleRoot = root
			}
			report, err := abihelpers.EncodeCommitReport(r)
			require.NoError(t, err)

			got, err := th.plugin.ShouldAcceptFinalizedReport(
				testutils.Context(t),
				types.ReportTimestamp{Epoch: tt.epoch, Round: tt.round}, report)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expected, got)

			if got { // already added to inflight, should not be accepted again
				got, err = th.plugin.ShouldAcceptFinalizedReport(
					testutils.Context(t),
					types.ReportTimestamp{Epoch: tt.epoch, Round: tt.round}, report)
				require.NoError(t, err)
				assert.False(t, got)
			}
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

func TestNextMin(t *testing.T) {
	lggr := logger.TestLogger(t)
	commitStore := mock_contracts.CommitStoreInterface{}
	cp := CommitReportingPlugin{config: CommitPluginConfig{commitStore: &commitStore}, inflightReports: newInflightCommitReportsContainer(time.Hour)}
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
		commitStore.On("GetExpectedNextSequenceNumber", mock.Anything).Return(tc.onChainMin, nil)
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

func Test_isStaleReport(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)
	merkleRoot1 := utils.Keccak256Fixed([]byte("some merkle root 1"))
	merkleRoot2 := utils.Keccak256Fixed([]byte("some merkle root 2"))

	t.Run("empty report", func(t *testing.T) {
		commitStore := mock_contracts.NewCommitStoreInterface(t)
		r := &CommitReportingPlugin{config: CommitPluginConfig{commitStore: commitStore}}
		isStale := r.isStaleReport(ctx, lggr, commit_store.CommitStoreCommitReport{}, false, types.ReportTimestamp{})
		assert.True(t, isStale)
	})

	t.Run("merkle root", func(t *testing.T) {
		const expNextSeqNum = uint64(9)

		commitStore := mock_contracts.NewCommitStoreInterface(t)
		commitStore.On("GetExpectedNextSequenceNumber", mock.Anything).Return(expNextSeqNum, nil)

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
