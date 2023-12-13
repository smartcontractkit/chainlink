package gas_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

type blockHistoryConfig struct {
	v uint16
}

func (b *blockHistoryConfig) EIP1559FeeCapBufferBlocks() uint16 {
	return b.v
}

func Test_FixedPriceEstimator(t *testing.T) {
	t.Parallel()
	maxGasPrice := assets.NewWeiI(1000000)

	t.Run("GetLegacyGas returns EvmGasPriceDefault from config, with multiplier applied", func(t *testing.T) {
		config := &gas.MockGasEstimatorConfig{}
		f := gas.NewFixedPriceEstimator(config, &blockHistoryConfig{}, logger.Test(t))

		config.PriceDefaultF = assets.NewWeiI(42)
		config.LimitMultiplierF = float32(1.1)
		config.PriceMaxF = maxGasPrice

		gasPrice, gasLimit, err := f.GetLegacyGas(testutils.Context(t), nil, 100000, maxGasPrice)
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))
		assert.Equal(t, assets.NewWeiI(42), gasPrice)
	})

	t.Run("GetLegacyGas returns user specified maximum gas price", func(t *testing.T) {
		config := &gas.MockGasEstimatorConfig{}
		config.PriceDefaultF = assets.NewWeiI(42)
		config.LimitMultiplierF = float32(1.1)
		config.PriceMaxF = assets.NewWeiI(35)
		f := gas.NewFixedPriceEstimator(config, &blockHistoryConfig{}, logger.Test(t))

		gasPrice, gasLimit, err := f.GetLegacyGas(testutils.Context(t), nil, 100000, assets.NewWeiI(30))
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))
		assert.Equal(t, assets.NewWeiI(30), gasPrice)
	})

	t.Run("GetLegacyGas returns global maximum gas price", func(t *testing.T) {
		config := &gas.MockGasEstimatorConfig{}
		config.PriceDefaultF = assets.NewWeiI(42)
		config.LimitMultiplierF = float32(1.1)
		config.PriceMaxF = assets.NewWeiI(20)

		f := gas.NewFixedPriceEstimator(config, &blockHistoryConfig{}, logger.Test(t))
		gasPrice, gasLimit, err := f.GetLegacyGas(testutils.Context(t), nil, 100000, assets.NewWeiI(30))
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))
		assert.Equal(t, assets.NewWeiI(20), gasPrice)
	})

	t.Run("BumpLegacyGas calls BumpLegacyGasPriceOnly", func(t *testing.T) {
		config := &gas.MockGasEstimatorConfig{}
		config.PriceDefaultF = assets.NewWeiI(42)
		config.LimitMultiplierF = float32(1.1)
		config.PriceMaxF = maxGasPrice
		config.BumpPercentF = uint16(10)
		config.BumpMinF = assets.NewWeiI(150)

		lggr := logger.TestSugared(t)
		f := gas.NewFixedPriceEstimator(config, &blockHistoryConfig{}, lggr)

		gasPrice, gasLimit, err := f.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, maxGasPrice, nil)
		require.NoError(t, err)

		expectedGasPrice, expectedGasLimit, err := gas.BumpLegacyGasPriceOnly(config, lggr, nil, assets.NewWeiI(42), 100000, maxGasPrice)
		require.NoError(t, err)

		assert.Equal(t, expectedGasLimit, gasLimit)
		assert.Equal(t, expectedGasPrice, gasPrice)
	})

	t.Run("GetDynamicFee returns defaults from config, with multiplier applied", func(t *testing.T) {
		config := &gas.MockGasEstimatorConfig{}
		config.LimitMultiplierF = float32(1.1)
		config.PriceMaxF = maxGasPrice
		config.TipCapDefaultF = assets.NewWeiI(52)
		config.FeeCapDefaultF = assets.NewWeiI(100)
		config.BumpThresholdF = uint64(3)

		lggr := logger.Test(t)
		f := gas.NewFixedPriceEstimator(config, &blockHistoryConfig{}, lggr)

		fee, gasLimit, err := f.GetDynamicFee(testutils.Context(t), 100000, maxGasPrice)
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))

		assert.Equal(t, assets.NewWeiI(52), fee.TipCap)
		assert.Equal(t, assets.NewWeiI(100), fee.FeeCap)

		// Gas bumping disabled
		config.BumpThresholdF = uint64(0)

		fee, gasLimit, err = f.GetDynamicFee(testutils.Context(t), 100000, maxGasPrice)
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))

		assert.Equal(t, assets.NewWeiI(52), fee.TipCap)
		assert.Equal(t, maxGasPrice, fee.FeeCap)

		// override max gas price
		fee, gasLimit, err = f.GetDynamicFee(testutils.Context(t), 100000, assets.NewWeiI(10))
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))

		assert.Equal(t, assets.NewWeiI(52), fee.TipCap)
		assert.Equal(t, assets.NewWeiI(10), fee.FeeCap)
	})

	t.Run("BumpDynamicFee calls BumpDynamicFeeOnly", func(t *testing.T) {
		config := &gas.MockGasEstimatorConfig{}
		config.LimitMultiplierF = float32(1.1)
		config.PriceMaxF = maxGasPrice
		config.TipCapDefaultF = assets.NewWeiI(52)
		config.BumpMinF = assets.NewWeiI(150)
		config.BumpPercentF = uint16(10)

		lggr := logger.TestSugared(t)
		f := gas.NewFixedPriceEstimator(config, &blockHistoryConfig{}, lggr)

		originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(25)}
		fee, gasLimit, err := f.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
		require.NoError(t, err)

		expectedFee, expectedGasLimit, err := gas.BumpDynamicFeeOnly(config, 0, lggr, nil, nil, originalFee, 100000, maxGasPrice)
		require.NoError(t, err)

		assert.Equal(t, expectedGasLimit, gasLimit)
		assert.Equal(t, expectedFee, fee)
	})
}
