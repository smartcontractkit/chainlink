package rollups

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestCustomCalldataDAOracle_NewCustomCalldata(t *testing.T) {
	oracleAddress := utils.RandomAddress().String()
	t.Parallel()

	t.Run("throws error if oracle type is not custom_calldata", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)
		daOracleConfig := CreateTestDAOracle(t, toml.DAOracleOPStack, oracleAddress, "")
		_, err := NewCustomCalldataDAOracle(logger.Test(t), ethClient, chaintype.ChainArbitrum, daOracleConfig)
		require.Error(t, err)
	})

	t.Run("throws error if CustomGasPriceCalldata is empty", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracleConfig := CreateTestDAOracle(t, toml.DAOracleCustomCalldata, oracleAddress, "")
		_, err := NewCustomCalldataDAOracle(logger.Test(t), ethClient, chaintype.ChainCelo, daOracleConfig)
		require.Error(t, err)
	})

	t.Run("correctly creates custom calldata DA oracle", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)
		calldata := "0x0000000000000000000000000000000000001234"

		daOracleConfig := CreateTestDAOracle(t, toml.DAOracleCustomCalldata, oracleAddress, calldata)
		oracle, err := NewCustomCalldataDAOracle(logger.Test(t), ethClient, chaintype.ChainZkSync, daOracleConfig)
		require.NoError(t, err)
		require.NotNil(t, oracle)
	})
}

func TestCustomCalldataDAOracle_getCustomCalldataGasPrice(t *testing.T) {
	oracleAddress := utils.RandomAddress().String()
	t.Parallel()

	t.Run("correctly fetches gas price if DA oracle config has custom calldata", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)
		expectedPriceHex := "0x32" // 50

		daOracleConfig := CreateTestDAOracle(t, toml.DAOracleCustomCalldata, oracleAddress, "0x0000000000000000000000000000000000001234")
		oracle, err := NewCustomCalldataDAOracle(logger.Test(t), ethClient, chaintype.ChainZkSync, daOracleConfig)
		require.NoError(t, err)

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			require.NotNil(t, callMsg.To)
			require.Equal(t, oracleAddress, callMsg.To.String())
			require.Nil(t, blockNumber)
		}).Return(hexutil.MustDecode(expectedPriceHex), nil).Once()

		price, err := oracle.getCustomCalldataGasPrice(tests.Context(t))
		require.NoError(t, err)
		require.Equal(t, big.NewInt(50), price)
	})

	t.Run("throws error if custom calldata fails to decode", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracleConfig := CreateTestDAOracle(t, toml.DAOracleCustomCalldata, oracleAddress, "0xblahblahblah")
		oracle, err := NewCustomCalldataDAOracle(logger.Test(t), ethClient, chaintype.ChainCelo, daOracleConfig)
		require.NoError(t, err)

		_, err = oracle.getCustomCalldataGasPrice(tests.Context(t))
		require.Error(t, err)
	})

	t.Run("throws error if CallContract call fails", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracleConfig := CreateTestDAOracle(t, toml.DAOracleCustomCalldata, oracleAddress, "0x0000000000000000000000000000000000000000000000000000000000000032")
		oracle, err := NewCustomCalldataDAOracle(logger.Test(t), ethClient, chaintype.ChainCelo, daOracleConfig)
		require.NoError(t, err)

		ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("RPC failure")).Once()

		_, err = oracle.getCustomCalldataGasPrice(tests.Context(t))
		require.Error(t, err)
	})
}
