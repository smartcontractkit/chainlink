package gas_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/services/gas/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FixedPriceEstimator(t *testing.T) {
	t.Parallel()

	// TODO: Add this test for BlockHistoryEstimator also
	t.Run("GetLegacyGas returns EvmGasPriceDefault from config, with multiplier applied", func(t *testing.T) {
		config := new(mocks.Config)
		f := gas.NewFixedPriceEstimator(config)

		config.On("EvmGasPriceDefault").Return(big.NewInt(42))
		config.On("EvmGasLimitMultiplier").Return(float32(1.1))

		gasPrice, gasLimit, err := f.GetLegacyGas(nil, 100000)
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))
		assert.Equal(t, big.NewInt(42), gasPrice)

		config.AssertExpectations(t)
	})

	t.Run("BumpLegacyGas calls BumpLegacyGasPriceOnly", func(t *testing.T) {
		config := new(mocks.Config)
		f := gas.NewFixedPriceEstimator(config)

		config.On("EvmGasPriceDefault").Return(big.NewInt(42))
		config.On("EvmGasBumpPercent").Return(uint16(10))
		config.On("EvmGasBumpWei").Return(big.NewInt(150))
		config.On("EvmMaxGasPriceWei").Return(big.NewInt(1000000))
		config.On("EvmGasLimitMultiplier").Return(float32(1.1))

		gasPrice, gasLimit, err := f.BumpLegacyGas(big.NewInt(42), 100000)
		require.NoError(t, err)

		expectedGasPrice, expectedGasLimit, err := gas.BumpLegacyGasPriceOnly(config, big.NewInt(42), 100000)
		require.NoError(t, err)

		assert.Equal(t, expectedGasLimit, gasLimit)
		assert.Equal(t, expectedGasPrice, gasPrice)

		config.AssertExpectations(t)
	})

	t.Run("GetDynamicFee returns defaults from config, with multiplier applied", func(t *testing.T) {
		config := new(mocks.Config)
		f := gas.NewFixedPriceEstimator(config)

		config.On("EvmGasLimitMultiplier").Return(float32(1.1))
		config.On("EvmGasTipCapDefault").Return(big.NewInt(52))
		config.On("EvmGasFeeCap").Return(big.NewInt(100))

		fee, gasLimit, err := f.GetDynamicFee(100000)
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))

		assert.Equal(t, big.NewInt(52), fee.TipCap)
		assert.Equal(t, big.NewInt(100), fee.FeeCap)

		config.AssertExpectations(t)
	})

	t.Run("BumpDynamicFee calls BumpDynamicFeeOnly", func(t *testing.T) {
		config := new(mocks.Config)
		f := gas.NewFixedPriceEstimator(config)

		config.On("EvmGasBumpPercent").Return(uint16(10))
		config.On("EvmGasBumpWei").Return(big.NewInt(150))
		config.On("EvmMaxGasPriceWei").Return(big.NewInt(1000000))
		config.On("EvmGasLimitMultiplier").Return(float32(1.1))
		config.On("EvmGasTipCapDefault").Return(big.NewInt(52))

		originalFee := gas.DynamicFee{FeeCap: big.NewInt(100), TipCap: big.NewInt(25)}
		fee, gasLimit, err := f.BumpDynamicFee(originalFee, 100000)
		require.NoError(t, err)

		expectedFee, expectedGasLimit, err := gas.BumpDynamicFeeOnly(config, originalFee, 100000)
		require.NoError(t, err)

		assert.Equal(t, expectedGasLimit, gasLimit)
		assert.Equal(t, expectedFee, fee)

		config.AssertExpectations(t)
	})
}
