package gas_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
)

func TestWrappedEvmEstimator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// fee values
	gasLimit := uint32(10)
	legacyFee := assets.NewWeiI(10)
	dynamicFee := gas.DynamicFee{
		FeeCap: assets.NewWeiI(20),
		TipCap: assets.NewWeiI(21),
	}

	cfg := mocks.NewConfig(t)
	e := mocks.NewEvmEstimator(t)
	e.On("GetDynamicFee", mock.Anything, mock.Anything, mock.Anything).
		Return(dynamicFee, gasLimit, nil).Once()
	e.On("GetLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Once()
	e.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(dynamicFee, gasLimit, nil).Once()
	e.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(legacyFee, gasLimit, nil).Once()

	// GetFee returns gas estimation based on configuration value
	t.Run("GetFee", func(t *testing.T) {
		// expect legacy fee data
		cfg.On("EvmEIP1559DynamicFees").Return(false).Once()
		estimator := gas.NewWrappedEvmEstimator(e, cfg)
		fee, max, err := estimator.GetFee(ctx, nil, 0, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.Dynamic)

		// expect dynamic fee data
		cfg.On("EvmEIP1559DynamicFees").Return(true).Once()
		estimator = gas.NewWrappedEvmEstimator(e, cfg)
		fee, max, err = estimator.GetFee(ctx, nil, 0, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.Dynamic.FeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.Dynamic.TipCap))
		assert.Nil(t, fee.Legacy)
	})

	// BumpFee returns bumped fee type based on original fee calculation
	t.Run("BumpFee", func(t *testing.T) {
		cfg.On("EvmEIP1559DynamicFees").Return(false).Once().Maybe()
		estimator := gas.NewWrappedEvmEstimator(e, cfg)

		// expect legacy fee data
		fee, max, err := estimator.BumpFee(ctx, gas.EvmFee{Legacy: assets.NewWeiI(0)}, 0, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, legacyFee.Equal(fee.Legacy))
		assert.Nil(t, fee.Dynamic)

		// expect dynamic fee data
		var d gas.DynamicFee
		fee, max, err = estimator.BumpFee(ctx, gas.EvmFee{Dynamic: &d}, 0, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, gasLimit, max)
		assert.True(t, dynamicFee.FeeCap.Equal(fee.Dynamic.FeeCap))
		assert.True(t, dynamicFee.TipCap.Equal(fee.Dynamic.TipCap))
		assert.Nil(t, fee.Legacy)

		// expect error
		_, _, err = estimator.BumpFee(ctx, gas.EvmFee{}, 0, nil, nil)
		assert.Error(t, err)
		_, _, err = estimator.BumpFee(ctx, gas.EvmFee{
			Legacy:  legacyFee,
			Dynamic: &dynamicFee,
		}, 0, nil, nil)
		assert.Error(t, err)
	})
}
