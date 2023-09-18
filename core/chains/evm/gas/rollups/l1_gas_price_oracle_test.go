package rollups

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestL1GasPriceOracle(t *testing.T) {
	t.Parallel()

	t.Run("Unsupported ChainType returns nil", func(t *testing.T) {
		ethClient := mocks.NewETHClient(t)

		assert.Panicsf(t, func() { NewL1GasPriceOracle(logger.TestLogger(t), ethClient, config.ChainCelo) }, "Received unspported chaintype %s", config.ChainCelo)
	})

	t.Run("Calling L1GasPrice on unstarted L1Oracle returns error", func(t *testing.T) {
		ethClient := mocks.NewETHClient(t)

		oracle := NewL1GasPriceOracle(logger.TestLogger(t), ethClient, config.ChainOptimismBedrock)

		_, err := oracle.GasPrice(testutils.Context(t))
		assert.EqualError(t, err, "L1GasPriceOracle is not started; cannot estimate gas")
	})

	t.Run("Calling GasPrice on started Arbitrum L1Oracle returns Arbitrum l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)

		ethClient := mocks.NewETHClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, ArbGasInfoAddress, callMsg.To.String())
			assert.Equal(t, ArbGasInfo_getL1BaseFeeEstimate, fmt.Sprintf("%x", callMsg.Data))
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle := NewL1GasPriceOracle(logger.TestLogger(t), ethClient, config.ChainArbitrum)
		require.NoError(t, oracle.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, oracle.Close()) })

		gasPrice, err := oracle.GasPrice(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started OPStack L1Oracle returns OPStack l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(200)

		ethClient := mocks.NewETHClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			assert.Equal(t, OPGasOracleAddress, callMsg.To.String())
			assert.Equal(t, OPGasOracle_l1BaseFee, fmt.Sprintf("%x", callMsg.Data))
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle := NewL1GasPriceOracle(logger.TestLogger(t), ethClient, config.ChainOptimismBedrock)
		require.NoError(t, oracle.Start(testutils.Context(t)))
		t.Cleanup(func() { assert.NoError(t, oracle.Close()) })

		gasPrice, err := oracle.GasPrice(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})
}
