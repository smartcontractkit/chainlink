package gas_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
)

func TestUniversalEstimatorGetLegacyGas(t *testing.T) {
	t.Parallel()

	var gasLimit uint64 = 21000
	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("fetches a new gas price when first called", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()

		cfg := gas.UniversalEstimatorConfig{}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gasPrice, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gasPrice)
	})

	t.Run("without forceRefetch enabled it fetches the cached gas price if not stale", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()

		cfg := gas.UniversalEstimatorConfig{CacheTimeout: 4 * time.Hour}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas1, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gas1)

		gas2, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gas2)
	})

	t.Run("without forceRefetch enabled it fetches the a new gas price if the cached one is stale", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(15), nil).Once()

		cfg := gas.UniversalEstimatorConfig{CacheTimeout: 0 * time.Second}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas1, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gas1)

		gas2, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(15), gas2)
	})

	t.Run("with forceRefetch enabled it updates the price even if not stale", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(15), nil).Once()

		cfg := gas.UniversalEstimatorConfig{CacheTimeout: 4 * time.Hour}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas1, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gas1)

		gas2, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice, feetypes.OptForceRefetch)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(15), gas2)
	})

	t.Run("will return max price if estimation exceeds it", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()

		cfg := gas.UniversalEstimatorConfig{}

		maxPrice := assets.NewWeiI(1)
		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas1, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, gas1)
	})
}

func TestUniversalEstimatorBumpLegacyGas(t *testing.T) {
	t.Parallel()

	var gasLimit uint64 = 21000
	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("bumps a previous attempt by BumpPercent", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		originalGasPrice := assets.NewWeiI(10)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()

		cfg := gas.UniversalEstimatorConfig{BumpPercent: 50}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gasPrice, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(15), gasPrice)
	})

	t.Run("fails if the original attempt is nil, or equal or higher than the max price", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)

		cfg := gas.UniversalEstimatorConfig{}
		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)

		var originalPrice *assets.Wei
		_, _, err := u.BumpLegacyGas(tests.Context(t), originalPrice, gasLimit, maxPrice, nil)
		assert.Error(t, err)

		originalPrice = assets.NewWeiI(100)
		_, _, err = u.BumpLegacyGas(tests.Context(t), originalPrice, gasLimit, maxPrice, nil)
		assert.Error(t, err)

	})

	t.Run("returns market gas price if bumped original fee is lower", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(80), nil).Once()
		originalGasPrice := assets.NewWeiI(10)

		cfg := gas.UniversalEstimatorConfig{}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(80), gas)
	})

	t.Run("returns max gas price if bumped original fee is higher", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1), nil).Once()
		originalGasPrice := assets.NewWeiI(10)

		cfg := gas.UniversalEstimatorConfig{BumpPercent: 50}

		maxPrice := assets.NewWeiI(14)
		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, gas)
	})

	t.Run("returns max gas price if the aggregation of max and original bumped fee is higher", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1), nil).Once()
		originalGasPrice := assets.NewWeiI(10)

		cfg := gas.UniversalEstimatorConfig{BumpPercent: 50}

		maxPrice := assets.NewWeiI(14)
		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, gas)
	})

	t.Run("fails if the bumped gas price is lower than the minimum bump percentage", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(100), nil).Once()
		originalGasPrice := assets.NewWeiI(100)

		cfg := gas.UniversalEstimatorConfig{BumpPercent: 20}

		// Price will be capped by the max price
		maxPrice := assets.NewWeiI(101)
		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.Error(t, err)
	})
}

func TestUniversalEstimatorGetDynamicFee(t *testing.T) {
	t.Parallel()

	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("fetches a new dynamic fee when first called", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		baseFee := big.NewInt(5)
		maxPriorityFeePerGas1 := big.NewInt(33)
		maxPriorityFeePerGas2 := big.NewInt(20)

		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas1, big.NewInt(5)}, {maxPriorityFeePerGas2, big.NewInt(5)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee, baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		blockHistoryLength := 2
		cfg := gas.UniversalEstimatorConfig{BlockHistoryRange: uint64(blockHistoryLength)}
		avrgPriorityFee := big.NewInt(0)
		avrgPriorityFee.Add(maxPriorityFeePerGas1, maxPriorityFeePerGas2).Div(avrgPriorityFee, big.NewInt(int64(blockHistoryLength)))
		maxFee := (*assets.Wei)(baseFee).AddPercentage(gas.BaseFeeBufferPercentage).Add((*assets.Wei)(avrgPriorityFee))

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		dynamicFee, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxFee, dynamicFee.FeeCap)
		assert.Equal(t, (*assets.Wei)(avrgPriorityFee), dynamicFee.TipCap)
	})

	t.Run("fails if BlockHistoryRange is zero and tries to fetch new prices", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)

		cfg := gas.UniversalEstimatorConfig{BlockHistoryRange: 0}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.Error(t, err)
	})

	t.Run("without forceRefetch enabled it fetches the cached dynamic fees if not stale", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		baseFee := big.NewInt(1)
		maxPriorityFeePerGas := big.NewInt(1)

		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(2)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{
			CacheTimeout:      4 * time.Hour,
			BlockHistoryRange: 1,
		}
		maxFee := (*assets.Wei)(baseFee).AddPercentage(gas.BaseFeeBufferPercentage).Add((*assets.Wei)(maxPriorityFeePerGas))

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		dynamicFee, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxFee, dynamicFee.FeeCap)
		assert.Equal(t, (*assets.Wei)(maxPriorityFeePerGas), dynamicFee.TipCap)

		dynamicFee2, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxFee, dynamicFee2.FeeCap)
		assert.Equal(t, (*assets.Wei)(maxPriorityFeePerGas), dynamicFee2.TipCap)

	})

	t.Run("fetches a new dynamic fee when first called", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		baseFee := big.NewInt(1)
		maxPriorityFeePerGas := big.NewInt(1)

		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(2)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{BlockHistoryRange: 1}
		maxFee := (*assets.Wei)(baseFee).AddPercentage(gas.BaseFeeBufferPercentage).Add((*assets.Wei)(maxPriorityFeePerGas))

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		dynamicFee, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxFee, dynamicFee.FeeCap)
		assert.Equal(t, (*assets.Wei)(maxPriorityFeePerGas), dynamicFee.TipCap)
	})

	t.Run("will return max price if tip cap or fee cap exceed it", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		baseFee := big.NewInt(1)
		maxPriorityFeePerGas := big.NewInt(3)
		maxPrice := assets.NewWeiI(2)

		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(5)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{BlockHistoryRange: 1}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		dynamicFee, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, dynamicFee.FeeCap)
		assert.Equal(t, maxPrice, dynamicFee.TipCap)
	})
}

func TestUniversalEstimatorBumpDynamicFee(t *testing.T) {
	t.Parallel()

	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("bumps a previous attempt by BumpPercent", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		originalFee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(20),
			TipCap: assets.NewWeiI(10),
		}

		// These values will be ignored because they are lower prices than the originalFee
		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{big.NewInt(5), big.NewInt(50)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{big.NewInt(5)},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{
			BlockHistoryRange: 2,
			BumpPercent:       50,
		}

		expectedFeeCap := originalFee.FeeCap.AddPercentage(cfg.BumpPercent)
		expectedTipCap := originalFee.TipCap.AddPercentage(cfg.BumpPercent)

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		dynamicFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedFeeCap, dynamicFee.FeeCap)
		assert.Equal(t, expectedTipCap, dynamicFee.TipCap)
	})

	t.Run("fails if the original attempt is invalid", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		maxPrice := assets.NewWeiI(20)
		cfg := gas.UniversalEstimatorConfig{BlockHistoryRange: 1}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		// nil original fee
		var originalFee gas.DynamicFee
		_, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.Error(t, err)

		// tip cap is higher than fee cap
		originalFee = gas.DynamicFee{
			FeeCap: assets.NewWeiI(10),
			TipCap: assets.NewWeiI(11),
		}
		_, err = u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.Error(t, err)

		// fee cap is equal or higher to max price
		originalFee = gas.DynamicFee{
			FeeCap: assets.NewWeiI(20),
			TipCap: assets.NewWeiI(10),
		}
		_, err = u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.Error(t, err)
	})

	t.Run("returns market prices bumped by BumpPercent if bumped original fee is lower", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		originalFee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(20),
			TipCap: assets.NewWeiI(10),
		}

		// Market fees
		baseFee := big.NewInt(5)
		maxPriorityFeePerGas := big.NewInt(33)
		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(100)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		maxFee := (*assets.Wei)(baseFee).AddPercentage(gas.BaseFeeBufferPercentage).Add((*assets.Wei)(maxPriorityFeePerGas))

		cfg := gas.UniversalEstimatorConfig{
			BlockHistoryRange: 1,
			BumpPercent:       50,
		}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		bumpedFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, (*assets.Wei)(maxPriorityFeePerGas), bumpedFee.TipCap)
		assert.Equal(t, maxFee, bumpedFee.FeeCap)
	})

	t.Run("fails if connectivity percentile value is reached", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		originalFee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(20),
			TipCap: assets.NewWeiI(10),
		}

		// Market fees
		baseFee := big.NewInt(5)
		maxPriorityFeePerGas := big.NewInt(33)
		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(30)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{
			BlockHistoryRange: 1,
			BumpPercent:       50,
		}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.Error(t, err)
	})

	t.Run("returns max price if the aggregation of max and original bumped fee is higher", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		originalFee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(20),
			TipCap: assets.NewWeiI(18),
		}

		maxPrice := assets.NewWeiI(25)
		// Market fees
		baseFee := big.NewInt(1)
		maxPriorityFeePerGas := big.NewInt(1)
		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(30)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{
			BlockHistoryRange: 1,
			BumpPercent:       50,
		}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		bumpedFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, bumpedFee.TipCap)
		assert.Equal(t, maxPrice, bumpedFee.FeeCap)
	})

	t.Run("fails if the bumped gas price is lower than the minimum bump percentage", func(t *testing.T) {
		client := mocks.NewUniversalEstimatorClient(t)
		originalFee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(20),
			TipCap: assets.NewWeiI(18),
		}

		maxPrice := assets.NewWeiI(21)
		// Market fees
		baseFee := big.NewInt(1)
		maxPriorityFeePerGas := big.NewInt(1)
		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(30)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil).Once()

		cfg := gas.UniversalEstimatorConfig{
			BlockHistoryRange: 1,
			BumpPercent:       50,
		}

		u := gas.NewUniversalEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.Error(t, err)
	})
}
