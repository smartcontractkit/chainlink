package pricegetter

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib/rpclibmocks"
)

type testParameters struct {
	cfg                          config.DynamicPriceGetterConfig
	evmClients                   map[uint64]DynamicPriceGetterClient
	tokens                       []common.Address
	expectedTokenPrices          map[common.Address]big.Int
	invalidConfigErrorExpected   bool
	priceResolutionErrorExpected bool
}

func TestDynamicPriceGetter(t *testing.T) {
	tests := []struct {
		name  string
		param testParameters
	}{
		{
			name:  "aggregator_only_valid",
			param: testParamAggregatorOnly(t),
		},
		{
			name:  "aggregator_only_valid_multi",
			param: testParamAggregatorOnlyMulti(t),
		},
		{
			name:  "static_only_valid",
			param: testParamStaticOnly(),
		},
		{
			name:  "aggregator_and_static_valid",
			param: testParamAggregatorAndStaticValid(t),
		},
		{
			name:  "aggregator_and_static_token_collision",
			param: testParamAggregatorAndStaticTokenCollision(t),
		},
		{
			name:  "no_aggregator_for_token",
			param: testParamNoAggregatorForToken(t),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pg, err := NewDynamicPriceGetter(test.param.cfg, test.param.evmClients)
			if test.param.invalidConfigErrorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			ctx := testutils.Context(t)
			// Check configured token
			unconfiguredTk := cciptypes.Address(utils.RandomAddress().String())
			cfgTokens, uncfgTokens, err := pg.FilterConfiguredTokens(ctx, []cciptypes.Address{unconfiguredTk})
			require.NoError(t, err)
			assert.Equal(t, []cciptypes.Address{}, cfgTokens)
			assert.Equal(t, []cciptypes.Address{unconfiguredTk}, uncfgTokens)
			// Build list of tokens to query.
			tokens := make([]cciptypes.Address, 0, len(test.param.tokens))
			for _, tk := range test.param.tokens {
				tokenAddr := ccipcalc.EvmAddrToGeneric(tk)
				tokens = append(tokens, tokenAddr)
			}
			prices, err := pg.TokenPricesUSD(ctx, tokens)
			if test.param.priceResolutionErrorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			// we expect prices for at least all queried tokens (it is possible that additional tokens are returned).
			assert.True(t, len(prices) >= len(test.param.expectedTokenPrices))
			// Check prices are matching expected result.
			for tk, expectedPrice := range test.param.expectedTokenPrices {
				if prices[cciptypes.Address(tk.String())] == nil {
					assert.Fail(t, "Token price not found")
				}
				assert.Equal(t, 0, expectedPrice.Cmp(prices[cciptypes.Address(tk.String())]),
					"Token price mismatch: expected price %v, got %v", expectedPrice, *prices[cciptypes.Address(tk.String())])
			}
		})
	}
}

func testParamAggregatorOnly(t *testing.T) testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
	tk3 := utils.RandomAddress()
	tk4 := utils.RandomAddress()
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			tk1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk4: {
				ChainID:                   104,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(1396818990),
		StartedAt:       big.NewInt(1704896575),
		UpdatedAt:       big.NewInt(1704896575),
		AnsweredInRound: big.NewInt(1000),
	}
	// Real ETH/USD example from OP.
	round2 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(2000),
		Answer:          big.NewInt(238879815123),
		StartedAt:       big.NewInt(1704897197),
		UpdatedAt:       big.NewInt(1704897197),
		AnsweredInRound: big.NewInt(2000),
	}
	// Real LINK/ETH example from OP.
	round3 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(3000),
		Answer:          big.NewInt(4468862777874802),
		StartedAt:       big.NewInt(1715743907),
		UpdatedAt:       big.NewInt(1715743907),
		AnsweredInRound: big.NewInt(3000),
	}
	// Fake data for a token with more than 18 decimals.
	round4 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(4000),
		Answer:          multExp(big.NewInt(1234567890), 10), // 20 digits.
		StartedAt:       big.NewInt(1715753907),
		UpdatedAt:       big.NewInt(1715753907),
		AnsweredInRound: big.NewInt(4000),
	}
	evmClients := map[uint64]DynamicPriceGetterClient{
		uint64(101): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockClient(t, []uint8{18}, []aggregator_v3_interface.LatestRoundData{round3}),
		uint64(104): mockClient(t, []uint8{20}, []aggregator_v3_interface.LatestRoundData{round4}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *multExp(round1.Answer, 10),         // expected in 1e18 format.
		tk2: *multExp(round2.Answer, 10),         // expected in 1e18 format.
		tk3: *round3.Answer,                      // already in 1e18 format (contract decimals==18).
		tk4: *multExp(big.NewInt(1234567890), 8), // expected in 1e18 format.
	}
	return testParameters{
		cfg:                        cfg,
		evmClients:                 evmClients,
		tokens:                     []common.Address{tk1, tk2, tk3, tk4},
		expectedTokenPrices:        expectedTokenPrices,
		invalidConfigErrorExpected: false,
	}
}

// testParamAggregatorOnlyMulti test with several tokens on chain 102.
func testParamAggregatorOnlyMulti(t *testing.T) testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
	tk3 := utils.RandomAddress()
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			tk1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk3: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(1396818990),
		StartedAt:       big.NewInt(1704896575),
		UpdatedAt:       big.NewInt(1704896575),
		AnsweredInRound: big.NewInt(1000),
	}
	// Real ETH/USD example from OP.
	round2 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(2000),
		Answer:          big.NewInt(238879815123),
		StartedAt:       big.NewInt(1704897197),
		UpdatedAt:       big.NewInt(1704897197),
		AnsweredInRound: big.NewInt(2000),
	}
	round3 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(3000),
		Answer:          big.NewInt(238879815125),
		StartedAt:       big.NewInt(1704897198),
		UpdatedAt:       big.NewInt(1704897198),
		AnsweredInRound: big.NewInt(3000),
	}
	evmClients := map[uint64]DynamicPriceGetterClient{
		uint64(101): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockClient(t, []uint8{8, 8}, []aggregator_v3_interface.LatestRoundData{round2, round3}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *multExp(round1.Answer, 10),
		tk2: *multExp(round2.Answer, 10),
		tk3: *multExp(round3.Answer, 10),
	}
	return testParameters{
		cfg:                        cfg,
		evmClients:                 evmClients,
		invalidConfigErrorExpected: false,
		tokens:                     []common.Address{tk1, tk2, tk3},
		expectedTokenPrices:        expectedTokenPrices,
	}
}

func testParamStaticOnly() testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
	tk3 := utils.RandomAddress()
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			tk1: {
				ChainID: 101,
				Price:   big.NewInt(1_234_000),
			},
			tk2: {
				ChainID: 102,
				Price:   big.NewInt(2_234_000),
			},
			tk3: {
				ChainID: 103,
				Price:   big.NewInt(3_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	evmClients := map[uint64]DynamicPriceGetterClient{}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *cfg.StaticPrices[tk1].Price,
		tk2: *cfg.StaticPrices[tk2].Price,
		tk3: *cfg.StaticPrices[tk3].Price,
	}
	return testParameters{
		cfg:                 cfg,
		evmClients:          evmClients,
		tokens:              []common.Address{tk1, tk2, tk3},
		expectedTokenPrices: expectedTokenPrices,
	}
}

func testParamAggregatorAndStaticValid(t *testing.T) testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
	tk3 := utils.RandomAddress()
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			tk1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			tk3: {
				ChainID: 103,
				Price:   big.NewInt(1_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(1396818990),
		StartedAt:       big.NewInt(1704896575),
		UpdatedAt:       big.NewInt(1704896575),
		AnsweredInRound: big.NewInt(1000),
	}
	// Real ETH/USD example from OP.
	round2 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(2000),
		Answer:          big.NewInt(238879815123),
		StartedAt:       big.NewInt(1704897197),
		UpdatedAt:       big.NewInt(1704897197),
		AnsweredInRound: big.NewInt(2000),
	}
	evmClients := map[uint64]DynamicPriceGetterClient{
		uint64(101): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round2}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *multExp(round1.Answer, 10),
		tk2: *multExp(round2.Answer, 10),
		tk3: *cfg.StaticPrices[tk3].Price,
	}
	return testParameters{
		cfg:                 cfg,
		evmClients:          evmClients,
		tokens:              []common.Address{tk1, tk2, tk3},
		expectedTokenPrices: expectedTokenPrices,
	}
}

func testParamAggregatorAndStaticTokenCollision(t *testing.T) testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
	tk3 := utils.RandomAddress()
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			tk1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			tk3: {
				ChainID: 103,
				Price:   big.NewInt(1_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(1396818990),
		StartedAt:       big.NewInt(1704896575),
		UpdatedAt:       big.NewInt(1704896575),
		AnsweredInRound: big.NewInt(1000),
	}
	// Real ETH/USD example from OP.
	round2 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(2000),
		Answer:          big.NewInt(238879815123),
		StartedAt:       big.NewInt(1704897197),
		UpdatedAt:       big.NewInt(1704897197),
		AnsweredInRound: big.NewInt(2000),
	}
	round3 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(3000),
		Answer:          big.NewInt(238879815124),
		StartedAt:       big.NewInt(1704897198),
		UpdatedAt:       big.NewInt(1704897198),
		AnsweredInRound: big.NewInt(3000),
	}
	evmClients := map[uint64]DynamicPriceGetterClient{
		uint64(101): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round3}),
	}
	return testParameters{
		cfg:                        cfg,
		evmClients:                 evmClients,
		tokens:                     []common.Address{tk1, tk2, tk3},
		invalidConfigErrorExpected: true,
	}
}

func testParamNoAggregatorForToken(t *testing.T) testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
	tk3 := utils.RandomAddress()
	tk4 := utils.RandomAddress()
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			tk1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			tk2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			tk3: {
				ChainID: 103,
				Price:   big.NewInt(1_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(1396818990),
		StartedAt:       big.NewInt(1704896575),
		UpdatedAt:       big.NewInt(1704896575),
		AnsweredInRound: big.NewInt(1000),
	}
	// Real ETH/USD example from OP.
	round2 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(2000),
		Answer:          big.NewInt(238879815123),
		StartedAt:       big.NewInt(1704897197),
		UpdatedAt:       big.NewInt(1704897197),
		AnsweredInRound: big.NewInt(2000),
	}
	evmClients := map[uint64]DynamicPriceGetterClient{
		uint64(101): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockClient(t, []uint8{8}, []aggregator_v3_interface.LatestRoundData{round2}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *round1.Answer,
		tk2: *round2.Answer,
		tk3: *cfg.StaticPrices[tk3].Price,
		tk4: *big.NewInt(0),
	}
	return testParameters{
		cfg:                          cfg,
		evmClients:                   evmClients,
		tokens:                       []common.Address{tk1, tk2, tk3, tk4},
		expectedTokenPrices:          expectedTokenPrices,
		priceResolutionErrorExpected: true,
	}
}

func mockClient(t *testing.T, decimals []uint8, rounds []aggregator_v3_interface.LatestRoundData) DynamicPriceGetterClient {
	return DynamicPriceGetterClient{
		BatchCaller: mockCaller(t, decimals, rounds),
	}
}

func mockCaller(t *testing.T, decimals []uint8, rounds []aggregator_v3_interface.LatestRoundData) *rpclibmocks.EvmBatchCaller {
	caller := rpclibmocks.NewEvmBatchCaller(t)

	// Mock batch calls per chain: all decimals calls then all latestRoundData calls.
	dataAndErrs := make([]rpclib.DataAndErr, 0, len(decimals)+len(rounds))
	for _, d := range decimals {
		dataAndErrs = append(dataAndErrs, rpclib.DataAndErr{
			Outputs: []any{d},
		})
	}
	for _, round := range rounds {
		dataAndErrs = append(dataAndErrs, rpclib.DataAndErr{
			Outputs: []any{round.RoundId, round.Answer, round.StartedAt, round.UpdatedAt, round.AnsweredInRound},
		})
	}
	caller.On("BatchCall", mock.Anything, uint64(0), mock.Anything).Return(dataAndErrs, nil).Maybe()
	return caller
}

// multExp returns the result of multiplying x by 10^e.
func multExp(x *big.Int, e int64) *big.Int {
	return big.NewInt(0).Mul(x, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(e), nil))
}
