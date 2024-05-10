package gas_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	rollupMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
)

func TestSuggestedPriceEstimator(t *testing.T) {
	t.Parallel()

	maxGasPrice := assets.NewWeiI(100)

	calldata := []byte{0x00, 0x00, 0x01, 0x02, 0x03}
	const gasLimit uint64 = 80000

	cfg := &gas.MockGasEstimatorConfig{BumpPercentF: 10, BumpMinF: assets.NewWei(big.NewInt(1)), BumpThresholdF: 1}

	t.Run("calling GetLegacyGas on unstarted estimator returns error", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		_, _, err := o.GetLegacyGas(tests.Context(t), calldata, gasLimit, maxGasPrice)
		assert.EqualError(t, err, "estimator is not started")
	})

	t.Run("calling GetLegacyGas on started estimator returns prices", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(tests.Context(t), calldata, gasLimit, maxGasPrice)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(42), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("gas price is lower than user specified max gas price", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})

		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(tests.Context(t), calldata, gasLimit, assets.NewWeiI(40))
		require.Error(t, err)
		assert.EqualError(t, err, "estimated gas price: 42 wei is greater than the maximum gas price configured: 40 wei")
		assert.Nil(t, gasPrice)
		assert.Equal(t, uint64(0), chainSpecificGasLimit)
	})

	t.Run("gas price is lower than global max gas price", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(120)
		})

		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(tests.Context(t), calldata, gasLimit, assets.NewWeiI(110))
		assert.EqualError(t, err, "estimated gas price: 120 wei is greater than the maximum gas price configured: 110 wei")
		assert.Nil(t, gasPrice)
		assert.Equal(t, uint64(0), chainSpecificGasLimit)
	})

	t.Run("calling GetLegacyGas on started estimator if initial call failed returns error", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(pkgerrors.New("kaboom"))

		servicetest.RunHealthy(t, o)

		_, _, err := o.GetLegacyGas(tests.Context(t), calldata, gasLimit, maxGasPrice)
		assert.EqualError(t, err, "failed to estimate gas; gas price not set")
	})

	t.Run("calling GetDynamicFee always returns error", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		_, err := o.GetDynamicFee(tests.Context(t), maxGasPrice)
		assert.EqualError(t, err, "dynamic fees are not implemented for this estimator")
	})

	t.Run("calling BumpLegacyGas on unstarted estimator returns error", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		_, _, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(42), gasLimit, maxGasPrice, nil)
		assert.EqualError(t, err, "estimator is not started")
	})

	t.Run("calling BumpDynamicFee always returns error", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		fee := gas.DynamicFee{
			FeeCap: assets.NewWeiI(42),
			TipCap: assets.NewWeiI(5),
		}
		_, err := o.BumpDynamicFee(tests.Context(t), fee, maxGasPrice, nil)
		assert.EqualError(t, err, "dynamic fees are not implemented for this estimator")
	})

	t.Run("calling BumpLegacyGas on started estimator returns new price buffered with bumpPercent", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(40)
		})

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, maxGasPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(44), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("calling BumpLegacyGas on started estimator returns new price buffered with bumpMin", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(40)
		})

		testCfg := &gas.MockGasEstimatorConfig{BumpPercentF: 1, BumpMinF: assets.NewWei(big.NewInt(1)), BumpThresholdF: 1, LimitMultiplierF: 1}
		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, testCfg, l1Oracle)
		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, maxGasPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(41), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("calling BumpLegacyGas on started estimator returns original price when lower than previous", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(5)
		})

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)
		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, maxGasPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(10), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("calling BumpLegacyGas on started estimator returns error, suggested gas price is higher than max gas price", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})

		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, assets.NewWeiI(40), nil)
		require.Error(t, err)
		assert.EqualError(t, err, "estimated gas price: 42 wei is greater than the maximum gas price configured: 40 wei")
		assert.Nil(t, gasPrice)
		assert.Equal(t, uint64(0), chainSpecificGasLimit)
	})

	t.Run("calling BumpLegacyGas on started estimator returns max gas price when suggested price under max but the buffer exceeds it", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(39)
		})

		servicetest.RunHealthy(t, o)
		gasPrice, chainSpecificGasLimit, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, assets.NewWeiI(40), nil)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(40), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("calling BumpLegacyGas on started estimator if initial call failed returns error", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(pkgerrors.New("kaboom"))

		servicetest.RunHealthy(t, o)

		_, _, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, maxGasPrice, nil)
		assert.EqualError(t, err, "failed to refresh and return gas; gas price not set")
	})

	t.Run("calling BumpLegacyGas on started estimator if refresh call failed returns price from previous update", func(t *testing.T) {
		feeEstimatorClient := mocks.NewFeeEstimatorClient(t)
		l1Oracle := rollupMocks.NewL1Oracle(t)

		o := gas.NewSuggestedPriceEstimator(logger.Test(t), feeEstimatorClient, cfg, l1Oracle)

		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(40)
		}).Once()
		feeEstimatorClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(pkgerrors.New("kaboom"))

		servicetest.RunHealthy(t, o)

		gasPrice, chainSpecificGasLimit, err := o.BumpLegacyGas(tests.Context(t), assets.NewWeiI(10), gasLimit, maxGasPrice, nil)
		require.NoError(t, err)
		assert.Equal(t, assets.NewWeiI(44), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})
}
