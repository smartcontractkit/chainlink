package ccipdata_test

import (
	"context"
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

type priceRegReaderTH struct {
	lp      logpoller.LogPollerTest
	ec      client.Client
	lggr    logger.Logger
	user    *bind.TransactOpts
	readers map[string]ccipdata.PriceRegistryReader

	// Expected state
	blockTs              []uint64
	expectedFeeTokens    []common.Address
	expectedGasUpdates   map[uint64][]cciptypes.GasPrice
	expectedTokenUpdates map[uint64][]cciptypes.TokenPrice
	destSelectors        []uint64
}

func commitAndGetBlockTs(ec *client.SimulatedBackendClient) uint64 {
	h := ec.Commit()
	b, _ := ec.BlockByHash(context.Background(), h)
	return b.Time()
}

func newSim(t *testing.T) (*bind.TransactOpts, *client.SimulatedBackendClient) {
	user := evmtestutils.MustNewSimTransactor(t)
	sim := simulated.NewBackend(map[common.Address]types.Account{
		user.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	ec := client.NewSimulatedBackendClient(t, sim, testutils.SimulatedChainID)
	return user, ec
}

// setupPriceRegistryReaderTH instantiates all versions of the price registry reader
// with a snapshot of data so reader tests can do multi-version assertions.
func setupPriceRegistryReaderTH(t *testing.T) priceRegReaderTH {
	user, ec := newSim(t)
	lggr := logger.Test(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	headTracker := headstest.NewSimulatedHeadTracker(ec, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	if lpOpts.PollPeriod == 0 {
		lpOpts.PollPeriod = 1 * time.Hour
	}
	// TODO: We should be able to use an in memory log poller ORM here to speed up the tests.
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, pgtest.NewSqlxDB(t), lggr), ec, lggr, headTracker, lpOpts)

	feeTokens := []common.Address{utils.RandomAddress(), utils.RandomAddress()}
	dest1 := uint64(10)
	dest2 := uint64(11)
	gasPriceUpdatesBlock1 := []cciptypes.GasPrice{
		{
			DestChainSelector: dest1,
			Value:             big.NewInt(11),
		},
	}
	gasPriceUpdatesBlock2 := []cciptypes.GasPrice{
		{
			DestChainSelector: dest1,          // Reset same gas price
			Value:             big.NewInt(12), // Intentionally different from block1
		},
		{
			DestChainSelector: dest2, // Set gas price for different chain
			Value:             big.NewInt(12),
		},
	}
	token1 := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
	token2 := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
	tokenPriceUpdatesBlock1 := []cciptypes.TokenPrice{
		{
			Token: token1,
			Value: big.NewInt(12),
		},
	}
	tokenPriceUpdatesBlock2 := []cciptypes.TokenPrice{
		{
			Token: token1,
			Value: big.NewInt(13), // Intentionally change token1 value
		},
		{
			Token: token2,
			Value: big.NewInt(12), // Intentionally set a same value different token
		},
	}
	ctx := testutils.Context(t)
	addr2, _, _, err := price_registry_1_2_0.DeployPriceRegistry(user, ec, nil, feeTokens, 1000)
	require.NoError(t, err)
	commitAndGetBlockTs(ec) // Deploy these
	pr12r, err := factory.NewPriceRegistryReader(ctx, lggr, factory.NewEvmVersionFinder(), ccipcalc.EvmAddrToGeneric(addr2), lp, ec)
	require.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(pr12r).String(), reflect.TypeOf(&v1_2_0.PriceRegistry{}).String())
	// Apply block1.
	v1_2_0.ApplyPriceRegistryUpdate(t, user, addr2, ec, gasPriceUpdatesBlock1, tokenPriceUpdatesBlock1)
	b1 := commitAndGetBlockTs(ec)
	// Apply block2
	v1_2_0.ApplyPriceRegistryUpdate(t, user, addr2, ec, gasPriceUpdatesBlock2, tokenPriceUpdatesBlock2)
	b2 := commitAndGetBlockTs(ec)

	// Capture all lp data.
	lp.PollAndSaveLogs(ctx, 1)

	return priceRegReaderTH{
		lp:   lp,
		ec:   ec,
		lggr: lggr,
		user: user,
		readers: map[string]ccipdata.PriceRegistryReader{
			ccipdata.V1_2_0: pr12r,
		},
		expectedFeeTokens: feeTokens,
		expectedGasUpdates: map[uint64][]cciptypes.GasPrice{
			b1: gasPriceUpdatesBlock1,
			b2: gasPriceUpdatesBlock2,
		},
		expectedTokenUpdates: map[uint64][]cciptypes.TokenPrice{
			b1: tokenPriceUpdatesBlock1,
			b2: tokenPriceUpdatesBlock2,
		},
		blockTs:       []uint64{b1, b2},
		destSelectors: []uint64{dest1, dest2},
	}
}

func testPriceRegistryReader(t *testing.T, th priceRegReaderTH, pr ccipdata.PriceRegistryReader) {
	ctx := testutils.Context(t)
	// Assert have expected fee tokens.
	gotFeeTokens, err := pr.GetFeeTokens(ctx)
	require.NoError(t, err)
	evmAddrs, err := ccipcalc.GenericAddrsToEvm(gotFeeTokens...)
	require.NoError(t, err)
	assert.Equal(t, th.expectedFeeTokens, evmAddrs)

	// Note unsupported chain selector simply returns an empty set not an error
	gasUpdates, err := pr.GetGasPriceUpdatesCreatedAfter(ctx, 1e6, time.Unix(0, 0), 0)
	require.NoError(t, err)
	assert.Empty(t, gasUpdates)

	for i, ts := range th.blockTs {
		// Should see all updates >= ts.
		var expectedGas []cciptypes.GasPrice
		var expectedDest0Gas []cciptypes.GasPrice
		var expectedToken []cciptypes.TokenPrice
		for j := i; j < len(th.blockTs); j++ {
			expectedGas = append(expectedGas, th.expectedGasUpdates[th.blockTs[j]]...)
			for _, g := range th.expectedGasUpdates[th.blockTs[j]] {
				if g.DestChainSelector == th.destSelectors[0] {
					expectedDest0Gas = append(expectedDest0Gas, g)
				}
			}
			expectedToken = append(expectedToken, th.expectedTokenUpdates[th.blockTs[j]]...)
		}
		if ts > math.MaxInt64 {
			t.Fatalf("timestamp overflows int64: %d", ts)
		}
		unixTS := time.Unix(int64(ts-1), 0) //nolint:gosec // G115 false positive
		gasUpdates, err = pr.GetAllGasPriceUpdatesCreatedAfter(ctx, unixTS, 0)
		require.NoError(t, err)
		assert.Len(t, gasUpdates, len(expectedGas))

		gasUpdates, err = pr.GetGasPriceUpdatesCreatedAfter(ctx, th.destSelectors[0], unixTS, 0)
		require.NoError(t, err)
		assert.Len(t, gasUpdates, len(expectedDest0Gas))

		tokenUpdates, err2 := pr.GetTokenPriceUpdatesCreatedAfter(ctx, unixTS, 0)
		require.NoError(t, err2)
		assert.Len(t, tokenUpdates, len(expectedToken))
	}

	// Empty token set should return empty set no error.
	gotEmpty, err := pr.GetTokenPrices(ctx, []cciptypes.Address{})
	require.NoError(t, err)
	assert.Empty(t, gotEmpty)

	// We expect latest token prices to apply
	allTokenUpdates, err := pr.GetTokenPriceUpdatesCreatedAfter(ctx, time.Unix(0, 0), 0)
	require.NoError(t, err)
	// Build latest map
	latest := make(map[cciptypes.Address]*big.Int)
	// Comes back in ascending order (oldest first)
	var allTokens []cciptypes.Address
	for i := len(allTokenUpdates) - 1; i >= 0; i-- {
		assert.NoError(t, err)
		_, have := latest[allTokenUpdates[i].Token]
		if have {
			continue
		}
		latest[allTokenUpdates[i].Token] = allTokenUpdates[i].Value
		allTokens = append(allTokens, allTokenUpdates[i].Token)
	}
	tokenPrices, err := pr.GetTokenPrices(ctx, allTokens)
	require.NoError(t, err)
	require.Len(t, tokenPrices, len(allTokens))
	for _, p := range tokenPrices {
		assert.Equal(t, p.Value, latest[p.Token])
	}
}

func TestPriceRegistryReader(t *testing.T) {
	th := setupPriceRegistryReaderTH(t)
	// Assert all readers produce the same expected results.
	for version, pr := range th.readers {
		pr := pr
		t.Run("PriceRegistryReader"+version, func(t *testing.T) {
			testPriceRegistryReader(t, th, pr)
		})
	}
}

func TestNewPriceRegistryReader(t *testing.T) {
	var tt = []struct {
		typeAndVersion string
		expectedErr    string
	}{
		{
			typeAndVersion: "EVM2EVMOffRamp 1.2.0",
			expectedErr:    "expected PriceRegistry got EVM2EVMOffRamp",
		},
		{
			typeAndVersion: "PriceRegistry 1.2.0",
			expectedErr:    "",
		},
		{
			typeAndVersion: "PriceRegistry 1.6.0",
			expectedErr:    "",
		},
		{
			typeAndVersion: "PriceRegistry 2.0.0",
			expectedErr:    "unsupported price registry version 2.0.0",
		},
	}
	ctx := testutils.Context(t)
	for _, tc := range tt {
		t.Run(tc.typeAndVersion, func(t *testing.T) {
			b, err := utils.ABIEncode(`[{"type":"string"}]`, tc.typeAndVersion)
			require.NoError(t, err)
			c := clienttest.NewClient(t)
			c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(b, nil)
			addr := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
			lp := lpmocks.NewLogPoller(t)
			lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Maybe()
			_, err = factory.NewPriceRegistryReader(ctx, logger.Test(t), factory.NewEvmVersionFinder(), addr, lp, c)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
