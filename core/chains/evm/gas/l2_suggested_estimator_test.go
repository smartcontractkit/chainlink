package gas_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestL2SuggestedEstimator(t *testing.T) {
	t.Parallel()

	config := new(mocks.Config)
	client := new(mocks.RPCClient)
	o := gas.NewL2SuggestedEstimator(logger.TestLogger(t), config, client)

	calldata := []byte{0x00, 0x00, 0x01, 0x02, 0x03}
	var gasLimit uint64 = 80000

	t.Run("calling EstimateGas on unstarted estimator returns error", func(t *testing.T) {
		_, _, err := o.GetLegacyGas(calldata, gasLimit)
		assert.EqualError(t, err, "estimator is not started")
	})

	t.Run("calling EstimateGas on started estimator returns prices", func(t *testing.T) {
		client.On("Call", mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(0).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})

		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { require.NoError(t, o.Close()) })
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(calldata, gasLimit)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(42), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("calling BumpGas always returns error", func(t *testing.T) {
		_, _, err := o.BumpLegacyGas(big.NewInt(42), gasLimit, big.NewInt(10))
		assert.EqualError(t, err, "bump gas is not supported for this l2")
	})

	t.Run("calling EstimateGas on started estimator if initial call failed returns error", func(t *testing.T) {
		config := new(mocks.Config)
		client := new(mocks.RPCClient)
		o := gas.NewL2SuggestedEstimator(logger.TestLogger(t), config, client)

		client.On("Call", mock.Anything, "eth_gasPrice").Return(errors.New("kaboom"))

		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { require.NoError(t, o.Close()) })

		_, _, err := o.GetLegacyGas(calldata, gasLimit)
		assert.EqualError(t, err, "failed to estimate l2 gas; gas price not set")
	})
}
