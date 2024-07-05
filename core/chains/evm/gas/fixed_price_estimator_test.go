package gas_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
)

func Test_FixedPriceEstimator(t *testing.T) {
	t.Parallel()
	maxPrice := assets.NewWeiI(1000000)

	t.Run("GetLegacyGas returns PriceDefault from config", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("PriceDefault").Return(assets.NewWeiI(42))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		gasPrice, _, err := f.GetLegacyGas(tests.Context(t), nil, 100000, maxPrice)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(42), gasPrice)
	})

	t.Run("GetLegacyGas returns user specified maximum gas price if default is higher", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("PriceDefault").Return(assets.NewWeiI(42))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		maxPrice := assets.NewWeiI(30)
		gasPrice, _, err := f.GetLegacyGas(tests.Context(t), nil, 100000, maxPrice)
		require.NoError(t, err)
		assert.Equal(t, maxPrice, gasPrice)
	})

	t.Run("BumpLegacyGas fails if original gas price is invalid", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		// original gas price is nil
		maxPrice := assets.NewWeiI(30)
		var originalGasPrice *assets.Wei
		_, _, err := f.BumpLegacyGas(tests.Context(t), originalGasPrice, 100000, maxPrice, nil)
		require.Error(t, err)

		// original gas price is higher than max
		originalGasPrice = assets.NewWeiI(40)
		_, _, err = f.BumpLegacyGas(tests.Context(t), originalGasPrice, 100000, maxPrice, nil)
		require.Error(t, err)
	})

	t.Run("BumpLegacyGas bumps original gas price by BumpPercent", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("BumpPercent").Return(uint16(20))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		// original gas price is nil
		maxPrice := assets.NewWeiI(100)
		originalGasPrice := assets.NewWeiI(40)
		bumpedGas, _, err := f.BumpLegacyGas(tests.Context(t), originalGasPrice, 100000, maxPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, originalGasPrice.AddPercentage(20), bumpedGas)

	})

	t.Run("BumpLegacyGas bumps original gas price by BumpPercent but caps on max price", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("BumpPercent").Return(uint16(20))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		// original gas price is nil
		maxPrice := assets.NewWeiI(41)
		originalGasPrice := assets.NewWeiI(40)
		bumpedGas, _, err := f.BumpLegacyGas(tests.Context(t), originalGasPrice, 100000, maxPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, maxPrice, bumpedGas)

	})

	t.Run("GetDynamicFee returns TipCapDefault and FeeCapDefault from config", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("TipCapDefault").Return(assets.NewWeiI(10))
		config.On("FeeCapDefault").Return(assets.NewWeiI(20))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		dynamicFee, err := f.GetDynamicFee(tests.Context(t), maxPrice)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), dynamicFee.TipCap)
		assert.Equal(t, assets.NewWeiI(20), dynamicFee.FeeCap)
	})

	t.Run("GetDynamicFee returns user specified maximum price if defaults are higher", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("TipCapDefault").Return(assets.NewWeiI(10))
		config.On("FeeCapDefault").Return(assets.NewWeiI(20))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		maxPrice := assets.NewWeiI(8)

		dynamicFee, err := f.GetDynamicFee(tests.Context(t), maxPrice)
		require.NoError(t, err)
		assert.Equal(t, maxPrice, dynamicFee.TipCap)
		assert.Equal(t, maxPrice, dynamicFee.FeeCap)
	})

	t.Run("BumpDynamicFee fails if original fee is invalid", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		maxPrice := assets.NewWeiI(8)

		// original fee is nil
		var dynamicFee gas.DynamicFee
		_, err := f.BumpDynamicFee(tests.Context(t), dynamicFee, maxPrice, nil)
		require.Error(t, err)

		// TipCap is higher than FeeCap
		dynamicFee.FeeCap = assets.NewWeiI(10)
		dynamicFee.TipCap = assets.NewWeiI(11)
		_, err = f.BumpDynamicFee(tests.Context(t), dynamicFee, maxPrice, nil)
		require.Error(t, err)

		// FeeCap is higher than max price
		dynamicFee.FeeCap = assets.NewWeiI(10)
		dynamicFee.TipCap = assets.NewWeiI(8)
		_, err = f.BumpDynamicFee(tests.Context(t), dynamicFee, maxPrice, nil)
		require.Error(t, err)
	})

	t.Run("BumpDynamicFee bumps original fee by BumpPercent", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("BumpPercent").Return(uint16(20))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		maxPrice := assets.NewWeiI(100)
		feeCap := assets.NewWeiI(20)
		tipCap := assets.NewWeiI(10)
		dynamicFee := gas.DynamicFee{FeeCap: feeCap, TipCap: tipCap}
		bumpedFee, err := f.BumpDynamicFee(tests.Context(t), dynamicFee, maxPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, feeCap.AddPercentage(20), bumpedFee.FeeCap)
		assert.Equal(t, tipCap.AddPercentage(20), bumpedFee.TipCap)
	})

	t.Run("BumpDynamicFee bumps original fee by BumpPercent but caps on max price", func(t *testing.T) {
		config := gasmocks.NewFixedPriceEstimatorConfig(t)
		config.On("BumpPercent").Return(uint16(20))
		l1Oracle := rollupMocks.NewL1Oracle(t)
		f := gas.NewFixedPriceEstimator(logger.Test(t), config, l1Oracle)

		maxPrice := assets.NewWeiI(22)
		feeCap := assets.NewWeiI(20)
		tipCap := assets.NewWeiI(19)
		dynamicFee := gas.DynamicFee{FeeCap: feeCap, TipCap: tipCap}
		bumpedFee, err := f.BumpDynamicFee(tests.Context(t), dynamicFee, maxPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, maxPrice, bumpedFee.FeeCap)
		assert.Equal(t, maxPrice, bumpedFee.TipCap)
	})
}
