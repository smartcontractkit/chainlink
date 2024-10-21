package rollups

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func TestL1Oracle(t *testing.T) {
	t.Parallel()

	t.Run("Unsupported ChainType returns nil", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracle := CreateTestDAOracle(t, toml.DAOracleOPStack, utils.RandomAddress().String(), "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainCelo, daOracle)
		require.NoError(t, err)
		assert.Nil(t, oracle)
	})

	t.Run("DAOracle config is nil, falls back to using chainType", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainArbitrum, nil)
		require.NoError(t, err)
		assert.NotNil(t, oracle)
	})

	t.Run("DAOracle config is not nil, but OracleType is empty, falls back to using chainType arbitrum", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracle := CreateTestDAOracle(t, "", utils.RandomAddress().String(), "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainArbitrum, daOracle)
		require.NoError(t, err)
		assert.NotNil(t, oracle)
		assert.Equal(t, oracle.Name(), "L1GasOracle(arbitrum)")
	})

	t.Run("DAOracle config is not nil, but OracleType is empty, falls back to using chainType ZKSync", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracle := CreateTestDAOracle(t, "", utils.RandomAddress().String(), "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainZkSync, daOracle)
		require.NoError(t, err)
		assert.NotNil(t, oracle)
		assert.Equal(t, oracle.Name(), "L1GasOracle(zkSync)")
	})
}

func TestL1Oracle_GasPrice(t *testing.T) {
	t.Parallel()

	t.Run("Calling GasPrice on unstarted L1Oracle returns error", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		daOracle := CreateTestDAOracle(t, toml.DAOracleOPStack, utils.RandomAddress().String(), "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, daOracle)
		require.NoError(t, err)

		_, err = oracle.GasPrice(tests.Context(t))
		assert.EqualError(t, err, "L1GasOracle is not started; cannot estimate gas")
	})

	t.Run("Calling GasPrice on started Arbitrum L1Oracle returns Arbitrum l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)
		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(GetL1BaseFeeEstimateAbiString))
		require.NoError(t, err)

		ethClient := mocks.NewL1OracleClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("getL1BaseFeeEstimate")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		daOracle := CreateTestDAOracle(t, toml.DAOracleArbitrum, "0x0000000000000000000000000000000000000000", "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainArbitrum, daOracle)
		require.NoError(t, err)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started Kroma L1Oracle returns Kroma l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)

		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
		require.NoError(t, err)

		oracleAddress := utils.RandomAddress().String()
		ethClient := setupUpgradeCheck(t, oracleAddress, false, false) // Ecotone, Fjord disabled

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		daOracle := CreateTestDAOracle(t, toml.DAOracleOPStack, oracleAddress, "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainKroma, daOracle)
		require.NoError(t, err)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started OPStack L1Oracle returns OPStack l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(100)

		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
		require.NoError(t, err)

		oracleAddress := utils.RandomAddress().String()
		ethClient := setupUpgradeCheck(t, oracleAddress, false, false) // Ecotone, Fjord disabled

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		daOracle := CreateTestDAOracle(t, toml.DAOracleOPStack, oracleAddress, "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, daOracle)
		require.NoError(t, err)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started Scroll L1Oracle returns Scroll l1GasPrice", func(t *testing.T) {
		l1BaseFee := big.NewInt(200)
		l1GasPriceMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
		require.NoError(t, err)

		oracleAddress := utils.RandomAddress().String()
		ethClient := setupUpgradeCheck(t, oracleAddress, false, false) // Ecotone, Fjord disabled

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		daOracle := CreateTestDAOracle(t, toml.DAOracleOPStack, oracleAddress, "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainScroll, daOracle)
		require.NoError(t, err)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(l1BaseFee), gasPrice)
	})

	t.Run("Calling GasPrice on started zkSync L1Oracle returns ZkSync l1GasPrice", func(t *testing.T) {
		gasPerPubByteL2 := big.NewInt(1100)
		gasPriceL2 := big.NewInt(25000000)
		ZksyncGasInfo_getGasPriceL2 := "0xfe173b97"
		ZksyncGasInfo_getGasPerPubdataByteL2 := "0x7cb9357e"
		ethClient := mocks.NewL1OracleClient(t)

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err := hex.DecodeString(ZksyncGasInfo_getGasPriceL2[2:])
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(gasPriceL2).Bytes(), nil).Once()

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err := hex.DecodeString(ZksyncGasInfo_getGasPerPubdataByteL2[2:])
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(gasPerPubByteL2).Bytes(), nil)

		daOracle := CreateTestDAOracle(t, toml.DAOracleZKSync, "0x0000000000000000000000000000000000000000", "")
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainZkSync, daOracle)
		require.NoError(t, err)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(new(big.Int).Mul(gasPriceL2, gasPerPubByteL2)), gasPrice)
	})
}
