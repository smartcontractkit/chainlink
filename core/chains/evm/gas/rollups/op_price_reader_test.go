package rollups

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestDAPriceReader_GasDAPrice(t *testing.T) {
	t.Parallel()

	t.Run("Calling GasPrice on unstarted L1Oracle returns error", func(t *testing.T) {
		ethClient := mocks.NewETHClient(t)

		oracle := NewL1GasOracle(logger.Test(t), ethClient, config.ChainOptimismBedrock)

		_, err := oracle.GasPrice(testutils.Context(t))
		assert.EqualError(t, err, "L1GasOracle is not started; cannot estimate gas")
	})

	t.Run("Calling GasPrice on started Arbitrum L1Oracle returns Arbitrum l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)
		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(GetL1BaseFeeEstimateAbiString))
		require.NoError(t, err)

		ethClient := mocks.NewETHClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("getL1BaseFeeEstimate")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle := NewL1GasOracle(logger.Test(t), ethClient, config.ChainArbitrum)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started Kroma L1Oracle returns Kroma l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)
		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
		require.NoError(t, err)

		ethClient := mocks.NewETHClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle := NewL1GasOracle(logger.Test(t), ethClient, config.ChainKroma)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started OPStack L1Oracle returns OPStack l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)
		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
		require.NoError(t, err)

		ethClient := mocks.NewETHClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle := NewL1GasOracle(logger.Test(t), ethClient, config.ChainOptimismBedrock)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})
}
