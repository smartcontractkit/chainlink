package rollups

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func TestL1Oracle(t *testing.T) {
	t.Parallel()

	t.Run("Unsupported ChainType returns nil", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainCelo)
		require.NoError(t, err)
		assert.Nil(t, oracle)
	})
}

func TestL1Oracle_GasPrice(t *testing.T) {
	t.Parallel()

	t.Run("Calling GasPrice on unstarted L1Oracle returns error", func(t *testing.T) {
		ethClient := mocks.NewL1OracleClient(t)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock)
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

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainArbitrum)
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

		ethClient := setupUpgradeCheck(t, KromaGasOracleAddress, false, false) // Ecotone, Fjord disabled

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainKroma)
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

		ethClient := setupUpgradeCheck(t, OPGasOracleAddress, false, false) // Ecotone, Fjord disabled

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock)
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

		ethClient := setupUpgradeCheck(t, ScrollGasOracleAddress, false, false) // Ecotone, Fjord disabled

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasPriceMethodAbi.Pack("l1BaseFee")
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainScroll)
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

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainZkSync)
		require.NoError(t, err)
		servicetest.RunHealthy(t, oracle)

		gasPrice, err := oracle.GasPrice(tests.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWei(new(big.Int).Mul(gasPriceL2, gasPerPubByteL2)), gasPrice)
	})
}

func TestL1Oracle_GetGasCost(t *testing.T) {
	t.Parallel()

	t.Run("Calling GetGasCost on started Arbitrum L1Oracle returns Arbitrum getL1Fee", func(t *testing.T) {
		l1GasCost := big.NewInt(100)
		baseFee := utils.Uint256ToBytes32(big.NewInt(1000))
		l1BaseFeeEstimate := utils.Uint256ToBytes32(big.NewInt(500))
		blockNum := big.NewInt(1000)
		toAddress := utils.RandomAddress()
		callData := []byte{1, 2, 3, 4, 5, 6, 7}
		l1GasCostMethodAbi, err := abi.JSON(strings.NewReader(GasEstimateL1ComponentAbiString))
		require.NoError(t, err)

		tx := types.NewTx(&types.LegacyTx{
			Nonce: 42,
			To:    &toAddress,
			Data:  callData,
		})
		result := common.LeftPadBytes(l1GasCost.Bytes(), 8)
		result = append(result, baseFee...)
		result = append(result, l1BaseFeeEstimate...)

		ethClient := mocks.NewL1OracleClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasCostMethodAbi.Pack("gasEstimateL1Component", toAddress, false, callData)
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			require.Equal(t, blockNum, blockNumber)
		}).Return(result, nil)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainArbitrum)
		require.NoError(t, err)

		gasCost, err := oracle.GetGasCost(tests.Context(t), tx, blockNum)
		require.NoError(t, err)
		require.Equal(t, assets.NewWei(l1GasCost), gasCost)
	})

	t.Run("Calling GetGasCost on started Kroma L1Oracle returns error", func(t *testing.T) {
		blockNum := big.NewInt(1000)
		tx := types.NewTx(&types.LegacyTx{})

		ethClient := mocks.NewL1OracleClient(t)
		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainKroma)
		require.NoError(t, err)

		_, err = oracle.GetGasCost(tests.Context(t), tx, blockNum)
		require.Error(t, err, "L1 gas cost not supported for this chain: kroma")
	})

	t.Run("Calling GetGasCost on started OPStack L1Oracle returns OPStack getL1Fee", func(t *testing.T) {
		l1GasCost := big.NewInt(100)
		blockNum := big.NewInt(1000)
		toAddress := utils.RandomAddress()
		callData := []byte{1, 2, 3}
		l1GasCostMethodAbi, err := abi.JSON(strings.NewReader(GetL1FeeAbiString))
		require.NoError(t, err)

		tx := types.NewTx(&types.LegacyTx{
			Nonce: 42,
			To:    &toAddress,
			Data:  callData,
		})

		encodedTx, err := tx.MarshalBinary()
		require.NoError(t, err)

		ethClient := mocks.NewL1OracleClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasCostMethodAbi.Pack("getL1Fee", encodedTx)
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			require.Equal(t, blockNum, blockNumber)
		}).Return(common.BigToHash(l1GasCost).Bytes(), nil)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock)
		require.NoError(t, err)

		gasCost, err := oracle.GetGasCost(tests.Context(t), tx, blockNum)
		require.NoError(t, err)
		require.Equal(t, assets.NewWei(l1GasCost), gasCost)
	})

	t.Run("Calling GetGasCost on started Scroll L1Oracle returns Scroll getL1Fee", func(t *testing.T) {
		l1GasCost := big.NewInt(100)
		blockNum := big.NewInt(1000)
		toAddress := utils.RandomAddress()
		callData := []byte{1, 2, 3}
		l1GasCostMethodAbi, err := abi.JSON(strings.NewReader(GetL1FeeAbiString))
		require.NoError(t, err)

		tx := types.NewTx(&types.LegacyTx{
			Nonce: 42,
			To:    &toAddress,
			Data:  callData,
		})

		encodedTx, err := tx.MarshalBinary()
		require.NoError(t, err)

		ethClient := mocks.NewL1OracleClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			var payload []byte
			payload, err = l1GasCostMethodAbi.Pack("getL1Fee", encodedTx)
			require.NoError(t, err)
			require.Equal(t, payload, callMsg.Data)
			require.Equal(t, blockNum, blockNumber)
		}).Return(common.BigToHash(l1GasCost).Bytes(), nil)

		oracle, err := NewL1GasOracle(logger.Test(t), ethClient, chaintype.ChainScroll)
		require.NoError(t, err)

		gasCost, err := oracle.GetGasCost(tests.Context(t), tx, blockNum)
		require.NoError(t, err)
		require.Equal(t, assets.NewWei(l1GasCost), gasCost)
	})
}
