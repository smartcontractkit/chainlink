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
	t.Run("EstimateGas returns EthGasPriceDefault from config, with multiplier applied", func(t *testing.T) {
		config := new(mocks.Config)
		f := gas.NewFixedPriceEstimator(config)

		config.On("EthGasPriceDefault").Return(big.NewInt(42))
		config.On("EthGasLimitMultiplier").Return(float32(1.1))

		gasPrice, gasLimit, err := f.EstimateGas(nil, 100000)
		require.NoError(t, err)
		assert.Equal(t, 110000, int(gasLimit))
		assert.Equal(t, big.NewInt(42), gasPrice)

		config.AssertExpectations(t)
	})

	t.Run("BumpGas calls BumpGasPriceOnly", func(t *testing.T) {
		config := new(mocks.Config)
		f := gas.NewFixedPriceEstimator(config)

		config.On("EthGasPriceDefault").Return(big.NewInt(42))
		config.On("EthGasBumpPercent").Return(uint16(10))
		config.On("EthGasBumpWei").Return(big.NewInt(150))
		config.On("EthMaxGasPriceWei").Return(big.NewInt(1000000))
		config.On("EthGasLimitMultiplier").Return(float32(1.1))

		gasPrice, gasLimit, err := f.BumpGas(big.NewInt(42), 100000)
		require.NoError(t, err)

		expectedGasPrice, expectedGasLimit, err := gas.BumpGasPriceOnly(config, big.NewInt(42), 100000)
		require.NoError(t, err)

		assert.Equal(t, expectedGasLimit, gasLimit)
		assert.Equal(t, expectedGasPrice, gasPrice)

		config.AssertExpectations(t)
	})
}
