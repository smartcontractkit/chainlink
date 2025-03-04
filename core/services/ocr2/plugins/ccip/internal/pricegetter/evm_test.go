package pricegetter

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

type testParameters struct {
	cfg                          config.DynamicPriceGetterConfig
	contractReaders              map[uint64]types.ContractReader
	tokens                       []common.Address
	expectedTokenPrices          map[common.Address]big.Int
	expectedTokenPricesForAll    map[common.Address]big.Int
	evmCallErr                   bool
	invalidConfigErrorExpected   bool
	priceResolutionErrorExpected bool
}

var (
	TK1 common.Address
	TK2 common.Address
	TK3 common.Address
	TK4 common.Address
)

func init() {
	TK1 = utils.RandomAddress()
	TK2 = utils.RandomAddress()
	TK3 = utils.RandomAddress()
	TK4 = utils.RandomAddress()
}

func TestDynamicPriceGetterWithEmptyInput(t *testing.T) {
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
		{
			name:  "batchCall_returns_err",
			param: testParamBatchCallReturnsErr(t),
		},
		{
			name:  "less_inputs_than_defined_prices",
			param: testLessInputsThanDefinedPrices(t),
		},
		{
			name:  "get_all_tokens_aggregator_and_static",
			param: testGetAllTokensAggregatorAndStatic(t),
		},
		{
			name:  "get_all_tokens_aggregator_only",
			param: testGetAllTokensAggregatorOnly(t),
		},
		{
			name:  "get_all_tokens_static_only",
			param: testGetAllTokensStaticOnly(t),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pg, err := NewDynamicPriceGetter(test.param.cfg, test.param.contractReaders)
			if test.param.invalidConfigErrorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			ctx := testutils.Context(t)

			var prices map[cciptypes.Address]*big.Int
			var expectedTokens map[common.Address]big.Int
			if len(test.param.expectedTokenPricesForAll) == 0 {
				prices, err = pg.TokenPricesUSD(ctx, ccipcalc.EvmAddrsToGeneric(test.param.tokens...))
				if test.param.evmCallErr {
					require.Error(t, err)
					return
				}

				if test.param.priceResolutionErrorExpected {
					require.Error(t, err)
					return
				}
				expectedTokens = test.param.expectedTokenPrices
			} else {
				prices, err = pg.GetJobSpecTokenPricesUSD(ctx)
				expectedTokens = test.param.expectedTokenPricesForAll
			}

			require.NoError(t, err)
			// Ensure all expected prices are present.
			assert.True(t, len(prices) == len(expectedTokens))
			// Check prices are matching expected result.
			for tk, expectedPrice := range expectedTokens {
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
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK4: {
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockCR([]uint8{18}, cfg, []common.Address{TK3}, []aggregator_v3_interface.LatestRoundData{round3}),
		uint64(104): mockCR([]uint8{20}, cfg, []common.Address{TK4}, []aggregator_v3_interface.LatestRoundData{round4}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		TK1: *multExp(round1.Answer, 10),         // expected in 1e18 format.
		TK2: *multExp(round2.Answer, 10),         // expected in 1e18 format.
		TK3: *round3.Answer,                      // already in 1e18 format (contract decimals==18).
		TK4: *multExp(big.NewInt(1234567890), 8), // expected in 1e18 format.
	}
	return testParameters{
		cfg:                        cfg,
		contractReaders:            contractReaders,
		tokens:                     []common.Address{TK1, TK2, TK3, TK4},
		expectedTokenPrices:        expectedTokenPrices,
		invalidConfigErrorExpected: false,
	}
}

// testParamAggregatorOnlyMulti test with several tokens on chain 102.
func testParamAggregatorOnlyMulti(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK3: {
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8, 8}, cfg, []common.Address{TK2, TK3}, []aggregator_v3_interface.LatestRoundData{round2, round3}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		TK1: *multExp(round1.Answer, 10),
		TK2: *multExp(round2.Answer, 10),
		TK3: *multExp(round3.Answer, 10),
	}
	return testParameters{
		cfg:                        cfg,
		contractReaders:            contractReaders,
		invalidConfigErrorExpected: false,
		tokens:                     []common.Address{TK1, TK2, TK3},
		expectedTokenPrices:        expectedTokenPrices,
	}
}

func testParamStaticOnly() testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK1: {
				ChainID: 101,
				Price:   big.NewInt(1_234_000),
			},
			TK2: {
				ChainID: 102,
				Price:   big.NewInt(2_234_000),
			},
			TK3: {
				ChainID: 103,
				Price:   big.NewInt(3_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	contractReaders := map[uint64]types.ContractReader{}
	expectedTokenPrices := map[common.Address]big.Int{
		TK1: *cfg.StaticPrices[TK1].Price,
		TK2: *cfg.StaticPrices[TK2].Price,
		TK3: *cfg.StaticPrices[TK3].Price,
	}
	return testParameters{
		cfg:                 cfg,
		contractReaders:     contractReaders,
		tokens:              []common.Address{TK1, TK2, TK3},
		expectedTokenPrices: expectedTokenPrices,
	}
}

func testParamNoAggregatorForToken(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK3: {
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		TK1: *round1.Answer,
		TK2: *round2.Answer,
		TK3: *cfg.StaticPrices[TK3].Price,
		TK4: *big.NewInt(0),
	}
	return testParameters{
		cfg:                          cfg,
		contractReaders:              contractReaders,
		tokens:                       []common.Address{TK1, TK2, TK3, TK4},
		expectedTokenPrices:          expectedTokenPrices,
		priceResolutionErrorExpected: true,
	}
}

func testParamAggregatorAndStaticValid(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK3: {
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		TK1: *multExp(round1.Answer, 10),
		TK2: *multExp(round2.Answer, 10),
		TK3: *cfg.StaticPrices[TK3].Price,
	}
	return testParameters{
		cfg:                 cfg,
		contractReaders:     contractReaders,
		tokens:              []common.Address{TK1, TK2, TK3},
		expectedTokenPrices: expectedTokenPrices,
	}
}

func testParamAggregatorAndStaticTokenCollision(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK3: {
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockCR([]uint8{8}, cfg, []common.Address{TK3}, []aggregator_v3_interface.LatestRoundData{round3}),
	}
	return testParameters{
		cfg:                        cfg,
		contractReaders:            contractReaders,
		tokens:                     []common.Address{TK1, TK2, TK3},
		invalidConfigErrorExpected: true,
	}
}

func testParamBatchCallReturnsErr(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK3: {
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockErrCR(),
	}
	return testParameters{
		cfg:             cfg,
		contractReaders: contractReaders,
		tokens:          []common.Address{TK1, TK2, TK3},
		evmCallErr:      true,
	}
}

func testLessInputsThanDefinedPrices(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK4: {
				ChainID: 104,
				Price:   big.NewInt(1_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(3749350456),
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockCR([]uint8{8}, cfg, []common.Address{TK3}, []aggregator_v3_interface.LatestRoundData{round3}),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		TK1: *multExp(round1.Answer, 10),
		TK2: *multExp(round2.Answer, 10),
		TK3: *multExp(round3.Answer, 10),
	}
	return testParameters{
		cfg:                 cfg,
		contractReaders:     contractReaders,
		tokens:              []common.Address{TK1, TK2, TK3},
		expectedTokenPrices: expectedTokenPrices,
	}
}

func testGetAllTokensAggregatorAndStatic(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK4: {
				ChainID: 104,
				Price:   big.NewInt(1_234_000),
			},
		},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(3749350456),
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

	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockCR([]uint8{8}, cfg, []common.Address{TK3}, []aggregator_v3_interface.LatestRoundData{round3}),
	}
	expectedTokenPricesForAll := map[common.Address]big.Int{
		TK1: *multExp(round1.Answer, 10),
		TK2: *multExp(round2.Answer, 10),
		TK3: *multExp(round3.Answer, 10),
		TK4: *cfg.StaticPrices[TK4].Price,
	}
	return testParameters{
		cfg:                       cfg,
		expectedTokenPricesForAll: expectedTokenPricesForAll,
		contractReaders:           contractReaders,
	}
}

func testGetAllTokensAggregatorOnly(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			TK1: {
				ChainID:                   101,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK2: {
				ChainID:                   102,
				AggregatorContractAddress: utils.RandomAddress(),
			},
			TK3: {
				ChainID:                   103,
				AggregatorContractAddress: utils.RandomAddress(),
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{},
	}
	// Real LINK/USD example from OP.
	round1 := aggregator_v3_interface.LatestRoundData{
		RoundId:         big.NewInt(1000),
		Answer:          big.NewInt(3749350456),
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
	contractReaders := map[uint64]types.ContractReader{
		uint64(101): mockCR([]uint8{8}, cfg, []common.Address{TK1}, []aggregator_v3_interface.LatestRoundData{round1}),
		uint64(102): mockCR([]uint8{8}, cfg, []common.Address{TK2}, []aggregator_v3_interface.LatestRoundData{round2}),
		uint64(103): mockCR([]uint8{8}, cfg, []common.Address{TK3}, []aggregator_v3_interface.LatestRoundData{round3}),
	}

	expectedTokenPricesForAll := map[common.Address]big.Int{
		TK1: *multExp(round1.Answer, 10),
		TK2: *multExp(round2.Answer, 10),
		TK3: *multExp(round3.Answer, 10),
	}
	return testParameters{
		cfg:                       cfg,
		expectedTokenPricesForAll: expectedTokenPricesForAll,
		contractReaders:           contractReaders,
	}
}

func testGetAllTokensStaticOnly(t *testing.T) testParameters {
	cfg := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{},
		StaticPrices: map[common.Address]config.StaticPriceConfig{
			TK1: {
				ChainID: 101,
				Price:   big.NewInt(1_234_000),
			},
			TK2: {
				ChainID: 102,
				Price:   big.NewInt(2_234_000),
			},
			TK3: {
				ChainID: 103,
				Price:   big.NewInt(3_234_000),
			},
		},
	}

	contractReaders := map[uint64]types.ContractReader{}
	expectedTokenPricesForAll := map[common.Address]big.Int{
		TK1: *cfg.StaticPrices[TK1].Price,
		TK2: *cfg.StaticPrices[TK2].Price,
		TK3: *cfg.StaticPrices[TK3].Price,
	}
	return testParameters{
		cfg:                       cfg,
		contractReaders:           contractReaders,
		expectedTokenPricesForAll: expectedTokenPricesForAll,
	}
}

func mockCR(decimals []uint8, cfg config.DynamicPriceGetterConfig, addr []common.Address, rounds []aggregator_v3_interface.LatestRoundData) *mockContractReader {
	// Mock batch calls per chain: all decimals calls then all latestRoundData calls.
	bGLVR := make(types.BatchGetLatestValuesResult)

	for i := range len(decimals) {
		boundContract := types.BoundContract{
			Address: cfg.AggregatorPrices[addr[i]].AggregatorContractAddress.Hex(),
			Name:    fmt.Sprintf("%v_%v", OffchainAggregator, i),
		}
		bGLVR[boundContract] = types.ContractBatchResults{}
	}
	for i, d := range decimals {
		contractName := fmt.Sprintf("%v_%v", OffchainAggregator, i)
		readRes := types.BatchReadResult{
			ReadName: DecimalsMethodName,
		}
		readRes.SetResult(&d, nil)
		boundContract := types.BoundContract{
			Address: cfg.AggregatorPrices[addr[i]].AggregatorContractAddress.Hex(),
			Name:    contractName,
		}
		bGLVR[boundContract] = append(bGLVR[boundContract], readRes)
	}

	for i, r := range rounds {
		contractName := fmt.Sprintf("%v_%v", OffchainAggregator, i)
		readRes := types.BatchReadResult{
			ReadName: LatestRoundDataMethodName,
		}
		readRes.SetResult(&r, nil)
		boundContract := types.BoundContract{
			Address: cfg.AggregatorPrices[addr[i]].AggregatorContractAddress.Hex(),
			Name:    contractName,
		}
		bGLVR[boundContract] = append(bGLVR[boundContract], readRes)
	}

	return &mockContractReader{result: bGLVR}
}

func mockErrCR() *mockContractReader {
	return &mockContractReader{err: assert.AnError}
}

// multExp returns the result of multiplying x by 10^e.
func multExp(x *big.Int, e int64) *big.Int {
	return big.NewInt(0).Mul(x, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(e), nil))
}

type mockContractReader struct {
	types.UnimplementedContractReader
	result types.BatchGetLatestValuesResult
	err    error
}

func (m *mockContractReader) Bind(context.Context, []types.BoundContract) error {
	return nil
}

func (m *mockContractReader) BatchGetLatestValues(context.Context, types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}
