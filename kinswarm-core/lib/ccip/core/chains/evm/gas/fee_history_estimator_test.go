package gas_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
)

func TestFeeHistoryEstimatorLifecycle(t *testing.T) {
	t.Parallel()
	var gasLimit uint64 = 21000
	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("fails if you fetch gas price before the estimator starts", func(t *testing.T) {
		cfg := gas.FeeHistoryEstimatorConfig{
			BumpPercent:      20,
			RewardPercentile: 60,
			EIP1559:          false,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), nil, cfg, chainID, nil)
		_, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.ErrorContains(t, err, "gas price not set")
	})

	t.Run("fails to start if BumpPercent is lower than the minimum cap", func(t *testing.T) {
		cfg := gas.FeeHistoryEstimatorConfig{BumpPercent: 9}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), nil, cfg, chainID, nil)
		assert.ErrorContains(t, u.Start(tests.Context(t)), "BumpPercent")
	})

	t.Run("fails to start if RewardPercentile is higher than ConnectivityPercentile in EIP-1559", func(t *testing.T) {
		cfg := gas.FeeHistoryEstimatorConfig{
			BumpPercent:      20,
			RewardPercentile: 99,
			EIP1559:          true,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), nil, cfg, chainID, nil)
		assert.ErrorContains(t, u.Start(tests.Context(t)), "RewardPercentile")
	})

	t.Run("starts if configs are correct", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Maybe()

		cfg := gas.FeeHistoryEstimatorConfig{
			BumpPercent:      20,
			RewardPercentile: 10,
			CacheTimeout:     10 * time.Second,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), nil, cfg, chainID, nil)
		err := u.Start(tests.Context(t))
		assert.NoError(t, err)
		err = u.Close()
		assert.NoError(t, err)
	})
}

func TestFeeHistoryEstimatorGetLegacyGas(t *testing.T) {
	t.Parallel()

	var gasLimit uint64 = 21000
	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("fetches a new gas price when first called", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()

		cfg := gas.FeeHistoryEstimatorConfig{}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, err := u.RefreshGasPrice()
		assert.NoError(t, err)
		gasPrice, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gasPrice)
	})

	t.Run("will return max price if estimation exceeds it", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil).Once()

		cfg := gas.FeeHistoryEstimatorConfig{}

		maxPrice := assets.NewWeiI(1)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, err := u.RefreshGasPrice()
		assert.NoError(t, err)
		gas1, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, gas1)
	})

	t.Run("fails if gas price has not been set yet", func(t *testing.T) {
		cfg := gas.FeeHistoryEstimatorConfig{}

		maxPrice := assets.NewWeiI(1)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), nil, cfg, chainID, nil)
		_, _, err := u.GetLegacyGas(tests.Context(t), nil, gasLimit, maxPrice)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "gas price not set")
	})
}

func TestFeeHistoryEstimatorBumpLegacyGas(t *testing.T) {
	t.Parallel()

	var gasLimit uint64 = 21000
	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("bumps a previous attempt by BumpPercent", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		originalGasPrice := assets.NewWeiI(10)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(10), nil)

		cfg := gas.FeeHistoryEstimatorConfig{BumpPercent: 50, CacheTimeout: 5 * time.Second}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		servicetest.RunHealthy(t, u)
		gasPrice, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(15), gasPrice)
	})

	t.Run("fails if the original attempt is nil, or equal or higher than the max price", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)

		cfg := gas.FeeHistoryEstimatorConfig{}
		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)

		var originalPrice *assets.Wei
		_, _, err := u.BumpLegacyGas(tests.Context(t), originalPrice, gasLimit, maxPrice, nil)
		assert.Error(t, err)

		originalPrice = assets.NewWeiI(100)
		_, _, err = u.BumpLegacyGas(tests.Context(t), originalPrice, gasLimit, maxPrice, nil)
		assert.Error(t, err)
	})

	t.Run("returns market gas price if bumped original fee is lower", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(80), nil).Once()
		originalGasPrice := assets.NewWeiI(10)

		cfg := gas.FeeHistoryEstimatorConfig{}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(80), gas)
	})

	t.Run("returns max gas price if bumped original fee is higher", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1), nil).Once()
		originalGasPrice := assets.NewWeiI(10)

		cfg := gas.FeeHistoryEstimatorConfig{BumpPercent: 50}

		maxPrice := assets.NewWeiI(14)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, gas)
	})

	t.Run("returns max gas price if the aggregation of max and original bumped fee is higher", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(1), nil).Once()
		originalGasPrice := assets.NewWeiI(10)

		cfg := gas.FeeHistoryEstimatorConfig{BumpPercent: 50}

		maxPrice := assets.NewWeiI(14)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		gas, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, gas)
	})

	t.Run("fails if the bumped gas price is lower than the minimum bump percentage", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		client.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(100), nil).Once()
		originalGasPrice := assets.NewWeiI(100)

		cfg := gas.FeeHistoryEstimatorConfig{BumpPercent: 20}

		// Price will be capped by the max price
		maxPrice := assets.NewWeiI(101)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		_, _, err := u.BumpLegacyGas(tests.Context(t), originalGasPrice, gasLimit, maxPrice, nil)
		assert.Error(t, err)
	})
}

func TestFeeHistoryEstimatorGetDynamicFee(t *testing.T) {
	t.Parallel()

	maxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("fetches a new dynamic fee when first called", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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
		cfg := gas.FeeHistoryEstimatorConfig{BlockHistorySize: uint64(blockHistoryLength)}
		avrgPriorityFee := big.NewInt(0)
		avrgPriorityFee.Add(maxPriorityFeePerGas1, maxPriorityFeePerGas2).Div(avrgPriorityFee, big.NewInt(int64(blockHistoryLength)))
		maxFee := (*assets.Wei)(baseFee).AddPercentage(gas.BaseFeeBufferPercentage).Add((*assets.Wei)(avrgPriorityFee))

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		dynamicFee, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxFee, dynamicFee.FeeCap)
		assert.Equal(t, (*assets.Wei)(avrgPriorityFee), dynamicFee.TipCap)
	})

	t.Run("fails if dynamic prices have not been set yet", func(t *testing.T) {
		cfg := gas.FeeHistoryEstimatorConfig{}

		maxPrice := assets.NewWeiI(1)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), nil, cfg, chainID, nil)
		_, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "dynamic price not set")
	})

	t.Run("will return max price if tip cap or fee cap exceed it", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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

		cfg := gas.FeeHistoryEstimatorConfig{BlockHistorySize: 1}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		dynamicFee, err := u.GetDynamicFee(tests.Context(t), maxPrice)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, dynamicFee.FeeCap)
		assert.Equal(t, maxPrice, dynamicFee.TipCap)
	})
}

func TestFeeHistoryEstimatorBumpDynamicFee(t *testing.T) {
	t.Parallel()

	globalMaxPrice := assets.NewWeiI(100)
	chainID := big.NewInt(0)

	t.Run("bumps a previous attempt by BumpPercent", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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

		cfg := gas.FeeHistoryEstimatorConfig{
			BlockHistorySize: 2,
			BumpPercent:      50,
		}

		expectedFeeCap := originalFee.FeeCap.AddPercentage(cfg.BumpPercent)
		expectedTipCap := originalFee.TipCap.AddPercentage(cfg.BumpPercent)

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		dynamicFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, globalMaxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedFeeCap, dynamicFee.FeeCap)
		assert.Equal(t, expectedTipCap, dynamicFee.TipCap)
	})

	t.Run("fails if the original attempt is invalid", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		maxPrice := assets.NewWeiI(20)
		cfg := gas.FeeHistoryEstimatorConfig{BlockHistorySize: 1}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
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

	t.Run("returns market prices if bumped original fee is lower", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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

		cfg := gas.FeeHistoryEstimatorConfig{
			BlockHistorySize: 1,
			BumpPercent:      50,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		bumpedFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, globalMaxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, (*assets.Wei)(maxPriorityFeePerGas), bumpedFee.TipCap)
		assert.Equal(t, maxFee, bumpedFee.FeeCap)
	})

	t.Run("fails if connectivity percentile value is reached", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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

		cfg := gas.FeeHistoryEstimatorConfig{
			BlockHistorySize: 1,
			BumpPercent:      50,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		_, err = u.BumpDynamicFee(tests.Context(t), originalFee, globalMaxPrice, nil)
		assert.Error(t, err)
	})

	t.Run("returns max price if the aggregation of max and original bumped fee is higher", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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

		cfg := gas.FeeHistoryEstimatorConfig{
			BlockHistorySize: 1,
			BumpPercent:      50,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		bumpedFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, maxPrice, bumpedFee.TipCap)
		assert.Equal(t, maxPrice, bumpedFee.FeeCap)
	})

	t.Run("fails if the bumped gas price is lower than the minimum bump percentage", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
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

		cfg := gas.FeeHistoryEstimatorConfig{
			BlockHistorySize: 1,
			BumpPercent:      50,
		}

		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		err := u.RefreshDynamicPrice()
		assert.NoError(t, err)
		_, err = u.BumpDynamicFee(tests.Context(t), originalFee, maxPrice, nil)
		assert.Error(t, err)
	})

	t.Run("ignores maxPriorityFeePerGas if there is no mempool and forces refetch", func(t *testing.T) {
		client := mocks.NewFeeHistoryEstimatorClient(t)
		originalFee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(40),
			TipCap: assets.NewWeiI(0),
		}

		// Market fees
		baseFee := big.NewInt(10)
		maxPriorityFeePerGas := big.NewInt(0)
		feeHistoryResult := &ethereum.FeeHistory{
			OldestBlock:  big.NewInt(1),
			Reward:       [][]*big.Int{{maxPriorityFeePerGas, big.NewInt(0)}}, // first one represents market price and second one connectivity price
			BaseFee:      []*big.Int{baseFee},
			GasUsedRatio: nil,
		}
		client.On("FeeHistory", mock.Anything, mock.Anything, mock.Anything).Return(feeHistoryResult, nil)

		cfg := gas.FeeHistoryEstimatorConfig{
			BlockHistorySize: 0,
			BumpPercent:      20,
			CacheTimeout:     10 * time.Second,
			EIP1559:          true,
		}

		maxFeePerGas := assets.NewWei(baseFee).AddPercentage(gas.BaseFeeBufferPercentage)
		u := gas.NewFeeHistoryEstimator(logger.Test(t), client, cfg, chainID, nil)
		servicetest.RunHealthy(t, u)
		bumpedFee, err := u.BumpDynamicFee(tests.Context(t), originalFee, globalMaxPrice, nil)
		assert.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(0), (*assets.Wei)(maxPriorityFeePerGas))
		assert.Equal(t, maxFeePerGas, bumpedFee.FeeCap)
	})
}
