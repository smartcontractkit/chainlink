package pricegetter

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib/rpclibmocks"
)

type testParameters struct {
	cfg                          config.DynamicPriceGetterConfig
	evmClients                   map[uint64]DynamicPriceGetterClient
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
			// Build list of tokens to query.
			tokens := make([]cciptypes.Address, 0, len(test.param.expectedTokenPrices))
			for tk := range test.param.expectedTokenPrices {
				tokens = append(tokens, cciptypes.Address(tk.String()))
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
				assert.Equal(t, expectedPrice, *prices[cciptypes.Address(tk.String())])
			}
		})
	}
}

func testParamAggregatorOnly(t *testing.T) testParameters {
	tk1 := utils.RandomAddress()
	tk2 := utils.RandomAddress()
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
	evmClients := map[uint64]DynamicPriceGetterClient{
		uint64(101): mockClientFromRound(t, round1),
		uint64(102): mockClientFromRound(t, round2),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *round1.Answer,
		tk2: *round2.Answer,
	}
	return testParameters{
		cfg:                        cfg,
		evmClients:                 evmClients,
		invalidConfigErrorExpected: false,
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
		uint64(101): mockClientFromRound(t, round1),
		uint64(102): mockClientFromRound(t, round2),
	}
	expectedTokenPrices := map[common.Address]big.Int{
		tk1: *round1.Answer,
		tk2: *round2.Answer,
		tk3: *cfg.StaticPrices[tk3].Price,
	}
	return testParameters{
		cfg:                 cfg,
		evmClients:          evmClients,
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
		uint64(101): mockClientFromRound(t, round1),
		uint64(102): mockClientFromRound(t, round2),
		uint64(103): mockClientFromRound(t, round3),
	}
	return testParameters{
		cfg:                        cfg,
		evmClients:                 evmClients,
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
		uint64(101): mockClientFromRound(t, round1),
		uint64(102): mockClientFromRound(t, round2),
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
		expectedTokenPrices:          expectedTokenPrices,
		priceResolutionErrorExpected: true,
	}
}

func mockClientFromRound(t *testing.T, round aggregator_v3_interface.LatestRoundData) DynamicPriceGetterClient {
	return DynamicPriceGetterClient{
		BatchCaller: mockCallerFromRound(t, round),
	}
}

func mockCallerFromRound(t *testing.T, round aggregator_v3_interface.LatestRoundData) *rpclibmocks.EvmBatchCaller {
	caller := rpclibmocks.NewEvmBatchCaller(t)
	caller.On("BatchCall", mock.Anything, uint64(0), mock.Anything).Return(
		[]rpclib.DataAndErr{
			{
				Outputs: []any{round.RoundId, round.Answer, round.StartedAt, round.UpdatedAt, round.AnsweredInRound},
			},
		},
		nil,
	).Maybe()
	return caller
}
