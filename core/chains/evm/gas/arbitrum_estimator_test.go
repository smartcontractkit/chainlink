package gas_test

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestArbitrumEstimator(t *testing.T) {
	t.Parallel()

	maxGasPrice := assets.NewWeiI(100)
	const maxGasLimit uint32 = 500_000
	calldata := []byte{0x00, 0x00, 0x01, 0x02, 0x03}
	const gasLimit uint32 = 80000
	const gasPriceBufferPercentage = 50

	t.Run("calling GetLegacyGas on unstarted estimator returns error", func(t *testing.T) {
		config := mocks.NewConfig(t)
		rpcClient := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, rpcClient, ethClient)
		_, _, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, maxGasPrice)
		assert.EqualError(t, err, "estimator is not started")
	})

	var zeros bytes.Buffer
	zeros.Write(common.BigToHash(big.NewInt(0)).Bytes())
	zeros.Write(common.BigToHash(big.NewInt(0)).Bytes())
	zeros.Write(common.BigToHash(big.NewInt(123455)).Bytes())
	t.Run("calling GetLegacyGas on started estimator returns estimates", func(t *testing.T) {
		config := mocks.NewConfig(t)
		config.On("EvmGasLimitMax").Return(maxGasLimit)
		rpcClient := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, gas.ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, gas.ArbGasInfo_getPricesInArbGas, fmt.Sprintf("%x", callMsg.Data))
			assert.Equal(t, big.NewInt(-1), blockNumber)
		}).Return(zeros.Bytes(), nil)

		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, rpcClient, ethClient)
		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, o.Close()) })
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, maxGasPrice)
		require.NoError(t, err)
		// Expected price for a standard l2_suggested_estimator would be 42, but we add a fixed gasPriceBufferPercentage.
		assert.Equal(t, assets.NewWeiI(42).AddPercentage(gasPriceBufferPercentage), gasPrice)
		assert.Equal(t, gasLimit, chainSpecificGasLimit)
	})

	t.Run("gas price is lower than user specified max gas price", func(t *testing.T) {
		client := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		config := mocks.NewConfig(t)
		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, client, ethClient)

		client.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, gas.ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, gas.ArbGasInfo_getPricesInArbGas, fmt.Sprintf("%x", callMsg.Data))
			assert.Equal(t, big.NewInt(-1), blockNumber)
		}).Return(zeros.Bytes(), nil)

		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, o.Close()) })
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, assets.NewWeiI(40))
		require.Error(t, err)
		assert.EqualError(t, err, "estimated gas price: 42 wei is greater than the maximum gas price configured: 40 wei")
		assert.Nil(t, gasPrice)
		assert.Equal(t, uint32(0), chainSpecificGasLimit)
	})

	t.Run("gas price is lower than global max gas price", func(t *testing.T) {
		ethClient := mocks.NewETHClient(t)
		config := mocks.NewConfig(t)
		client := mocks.NewRPCClient(t)
		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, client, ethClient)

		client.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(120)
		})
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, gas.ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, gas.ArbGasInfo_getPricesInArbGas, fmt.Sprintf("%x", callMsg.Data))
			assert.Equal(t, big.NewInt(-1), blockNumber)
		}).Return(zeros.Bytes(), nil)

		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, o.Close()) })
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, assets.NewWeiI(110))
		assert.EqualError(t, err, "estimated gas price: 120 wei is greater than the maximum gas price configured: 110 wei")
		assert.Nil(t, gasPrice)
		assert.Equal(t, uint32(0), chainSpecificGasLimit)
	})

	t.Run("calling BumpLegacyGas always returns error", func(t *testing.T) {
		config := mocks.NewConfig(t)
		rpcClient := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, rpcClient, ethClient)
		_, _, err := o.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), gasLimit, assets.NewWeiI(10), nil)
		assert.EqualError(t, err, "bump gas is not supported for this l2")
	})

	t.Run("calling GetLegacyGas on started estimator if initial call failed returns error", func(t *testing.T) {
		config := mocks.NewConfig(t)
		client := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, client, ethClient)

		client.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(errors.New("kaboom"))
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, gas.ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, gas.ArbGasInfo_getPricesInArbGas, fmt.Sprintf("%x", callMsg.Data))
			assert.Equal(t, big.NewInt(-1), blockNumber)
		}).Return(zeros.Bytes(), nil)

		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, o.Close()) })

		_, _, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, maxGasPrice)
		assert.EqualError(t, err, "failed to estimate l2 gas; gas price not set")
	})

	t.Run("limit computes", func(t *testing.T) {
		config := mocks.NewConfig(t)
		config.On("EvmGasLimitMax").Return(maxGasLimit)
		rpcClient := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})
		const (
			perL2Tx       = 50_000
			perL1Calldata = 10_000
		)
		var expLimit = gasLimit + perL2Tx + perL1Calldata*uint32(len(calldata))

		var b bytes.Buffer
		b.Write(common.BigToHash(big.NewInt(perL2Tx)).Bytes())
		b.Write(common.BigToHash(big.NewInt(perL1Calldata)).Bytes())
		b.Write(common.BigToHash(big.NewInt(123455)).Bytes())
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, gas.ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, gas.ArbGasInfo_getPricesInArbGas, fmt.Sprintf("%x", callMsg.Data))
			assert.Equal(t, big.NewInt(-1), blockNumber)
		}).Return(b.Bytes(), nil)

		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, rpcClient, ethClient)
		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, o.Close()) })
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, maxGasPrice)
		require.NoError(t, err)
		require.NotNil(t, gasPrice)
		// Again, a normal l2_suggested_estimator would return 42, but arbitrum_estimator adds a buffer.
		assert.Equal(t, "63 wei", gasPrice.String())
		assert.Equal(t, expLimit, chainSpecificGasLimit, "expected %d but got %d", expLimit, chainSpecificGasLimit)
	})

	t.Run("limit exceeds max", func(t *testing.T) {
		config := mocks.NewConfig(t)
		config.On("EvmGasLimitMax").Return(maxGasLimit)
		rpcClient := mocks.NewRPCClient(t)
		ethClient := mocks.NewETHClient(t)
		rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_gasPrice").Return(nil).Run(func(args mock.Arguments) {
			res := args.Get(1).(*hexutil.Big)
			(*big.Int)(res).SetInt64(42)
		})
		const (
			perL2Tx       = 500_000
			perL1Calldata = 100_000
		)

		var b bytes.Buffer
		b.Write(common.BigToHash(big.NewInt(perL2Tx)).Bytes())
		b.Write(common.BigToHash(big.NewInt(perL1Calldata)).Bytes())
		b.Write(common.BigToHash(big.NewInt(123455)).Bytes())
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, gas.ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, gas.ArbGasInfo_getPricesInArbGas, fmt.Sprintf("%x", callMsg.Data))
			assert.Equal(t, big.NewInt(-1), blockNumber)
		}).Return(b.Bytes(), nil)

		o := gas.NewArbitrumEstimator(logger.TestLogger(t), config, rpcClient, ethClient)
		require.NoError(t, o.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, o.Close()) })
		gasPrice, chainSpecificGasLimit, err := o.GetLegacyGas(testutils.Context(t), calldata, gasLimit, maxGasPrice)
		require.Error(t, err, "expected error but got (%s, %d)", gasPrice, chainSpecificGasLimit)
	})
}
