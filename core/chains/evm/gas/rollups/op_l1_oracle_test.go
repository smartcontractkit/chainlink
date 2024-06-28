package rollups

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
)

func TestOPL1Oracle_ReadV1GasPrice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		isEcotoneError bool
		isFjordError   bool
	}{
		{
			name:           "calling isEcotone and isFjord returns false, fetches l1BaseFee",
			isEcotoneError: false,
			isFjordError:   false,
		},
		{
			name:           "calling isEcotone returns false and IsFjord errors when chain has not made Fjord upgrade, fetches l1BaseFee",
			isEcotoneError: false,
			isFjordError:   true,
		},
		{
			name:           "calling isEcotone and isFjord when chain has not made Ecotone upgrade, fetches l1BaseFee",
			isEcotoneError: true,
			isFjordError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l1BaseFee := big.NewInt(100)
			oracleAddress := common.HexToAddress("0x1234").String()

			l1BaseFeeMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
			require.NoError(t, err)
			l1BaseFeeCalldata, err := l1BaseFeeMethodAbi.Pack(l1BaseFeeMethod)
			require.NoError(t, err)

			// IsFjord calldata
			isFjordMethodAbi, err := abi.JSON(strings.NewReader(OPIsFjordAbiString))
			require.NoError(t, err)
			isFjordCalldata, err := isFjordMethodAbi.Pack(isFjordMethod)
			require.NoError(t, err)

			// IsEcotone calldata
			isEcotoneMethodAbi, err := abi.JSON(strings.NewReader(OPIsEcotoneAbiString))
			require.NoError(t, err)
			isEcotoneCalldata, err := isEcotoneMethodAbi.Pack(isEcotoneMethod)
			require.NoError(t, err)

			ethClient := mocks.NewL1OracleClient(t)
			ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
				rpcElements := args.Get(1).([]rpc.BatchElem)
				require.Equal(t, 2, len(rpcElements))
				for _, rE := range rpcElements {
					require.Equal(t, "eth_call", rE.Method)
					require.Equal(t, oracleAddress, rE.Args[0].(map[string]interface{})["to"])
					require.Equal(t, "latest", rE.Args[1])
				}
				require.Equal(t, hexutil.Bytes(isFjordCalldata), rpcElements[0].Args[0].(map[string]interface{})["data"])
				require.Equal(t, hexutil.Bytes(isEcotoneCalldata), rpcElements[1].Args[0].(map[string]interface{})["data"])
				isUpgraded := "0x0000000000000000000000000000000000000000000000000000000000000000"
				if tc.isFjordError {
					rpcElements[0].Error = fmt.Errorf("test error")
				} else {
					rpcElements[0].Result = &isUpgraded
				}
				if tc.isEcotoneError {
					rpcElements[1].Error = fmt.Errorf("test error")
				} else {
					rpcElements[1].Result = &isUpgraded
				}
			}).Return(nil).Once()

			ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
				callMsg := args.Get(1).(ethereum.CallMsg)
				blockNumber := args.Get(2).(*big.Int)
				require.Equal(t, l1BaseFeeCalldata, callMsg.Data)
				require.Equal(t, oracleAddress, callMsg.To.String())
				assert.Nil(t, blockNumber)
			}).Return(common.BigToHash(l1BaseFee).Bytes(), nil).Once()

			oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
			require.NoError(t, err)
			gasPrice, err := oracle.GetDAGasPrice(tests.Context(t))

			require.NoError(t, err)
			assert.Equal(t, l1BaseFee, gasPrice)
		})
	}
}

func setupUpgradeCheck(t *testing.T, oracleAddress string, isFjord, isEcotone bool) *mocks.L1OracleClient {
	trueHex := "0x0000000000000000000000000000000000000000000000000000000000000001"
	falseHex := "0x0000000000000000000000000000000000000000000000000000000000000000"
	boolToHexMap := map[bool]*string{
		true:  &trueHex,
		false: &falseHex,
	}
	// IsFjord calldata
	isFjordMethodAbi, err := abi.JSON(strings.NewReader(OPIsFjordAbiString))
	require.NoError(t, err)
	isFjordCalldata, err := isFjordMethodAbi.Pack(isFjordMethod)
	require.NoError(t, err)

	// IsEcotone calldata
	isEcotoneMethodAbi, err := abi.JSON(strings.NewReader(OPIsEcotoneAbiString))
	require.NoError(t, err)
	isEcotoneCalldata, err := isEcotoneMethodAbi.Pack(isEcotoneMethod)
	require.NoError(t, err)

	ethClient := mocks.NewL1OracleClient(t)
	ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
		rpcElements := args.Get(1).([]rpc.BatchElem)
		require.Equal(t, 2, len(rpcElements))
		for _, rE := range rpcElements {
			require.Equal(t, "eth_call", rE.Method)
			require.Equal(t, oracleAddress, rE.Args[0].(map[string]interface{})["to"])
			require.Equal(t, "latest", rE.Args[1])
		}
		require.Equal(t, hexutil.Bytes(isFjordCalldata), rpcElements[0].Args[0].(map[string]interface{})["data"])
		require.Equal(t, hexutil.Bytes(isEcotoneCalldata), rpcElements[1].Args[0].(map[string]interface{})["data"])

		rpcElements[0].Result = boolToHexMap[isFjord]
		rpcElements[1].Result = boolToHexMap[isEcotone]
	}).Return(nil).Once()

	return ethClient
}

func mockBatchContractCall(t *testing.T, ethClient *mocks.L1OracleClient, oracleAddress string, baseFeeVal, baseFeeScalarVal, blobBaseFeeVal, blobBaseFeeScalarVal, decimalsVal *big.Int) {
	// L1 base fee calldata
	l1BaseFeeMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
	require.NoError(t, err)
	l1BaseFeeCalldata, err := l1BaseFeeMethodAbi.Pack(l1BaseFeeMethod)
	require.NoError(t, err)

	// L1 base fee scalar calldata
	l1BaseFeeScalarMethodAbi, err := abi.JSON(strings.NewReader(OPBaseFeeScalarAbiString))
	require.NoError(t, err)
	l1BaseFeeScalarCalldata, err := l1BaseFeeScalarMethodAbi.Pack(baseFeeScalarMethod)
	require.NoError(t, err)

	// Blob base fee calldata
	blobBaseFeeMethodAbi, err := abi.JSON(strings.NewReader(OPBlobBaseFeeAbiString))
	require.NoError(t, err)
	blobBaseFeeCalldata, err := blobBaseFeeMethodAbi.Pack(blobBaseFeeMethod)
	require.NoError(t, err)

	// Blob base fee scalar calldata
	blobBaseFeeScalarMethodAbi, err := abi.JSON(strings.NewReader(OPBlobBaseFeeScalarAbiString))
	require.NoError(t, err)
	blobBaseFeeScalarCalldata, err := blobBaseFeeScalarMethodAbi.Pack(blobBaseFeeScalarMethod)
	require.NoError(t, err)

	// Decimals calldata
	decimalsMethodAbi, err := abi.JSON(strings.NewReader(OPDecimalsAbiString))
	require.NoError(t, err)
	decimalsCalldata, err := decimalsMethodAbi.Pack(decimalsMethod)
	require.NoError(t, err)

	ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
		rpcElements := args.Get(1).([]rpc.BatchElem)
		require.Equal(t, 5, len(rpcElements))

		for _, rE := range rpcElements {
			require.Equal(t, "eth_call", rE.Method)
			require.Equal(t, oracleAddress, rE.Args[0].(map[string]interface{})["to"])
			require.Equal(t, "latest", rE.Args[1])
		}

		require.Equal(t, hexutil.Bytes(l1BaseFeeCalldata), rpcElements[0].Args[0].(map[string]interface{})["data"])
		require.Equal(t, hexutil.Bytes(l1BaseFeeScalarCalldata), rpcElements[1].Args[0].(map[string]interface{})["data"])
		require.Equal(t, hexutil.Bytes(blobBaseFeeCalldata), rpcElements[2].Args[0].(map[string]interface{})["data"])
		require.Equal(t, hexutil.Bytes(blobBaseFeeScalarCalldata), rpcElements[3].Args[0].(map[string]interface{})["data"])
		require.Equal(t, hexutil.Bytes(decimalsCalldata), rpcElements[4].Args[0].(map[string]interface{})["data"])

		res1 := common.BigToHash(baseFeeVal).Hex()
		res2 := common.BigToHash(baseFeeScalarVal).Hex()
		res3 := common.BigToHash(blobBaseFeeVal).Hex()
		res4 := common.BigToHash(blobBaseFeeScalarVal).Hex()
		res5 := common.BigToHash(decimalsVal).Hex()
		rpcElements[0].Result = &res1
		rpcElements[1].Result = &res2
		rpcElements[2].Result = &res3
		rpcElements[3].Result = &res4
		rpcElements[4].Result = &res5
	}).Return(nil).Once()
}

func TestOPL1Oracle_CalculateEcotoneGasPrice(t *testing.T) {
	baseFee := big.NewInt(100000000)
	blobBaseFee := big.NewInt(25000000)
	baseFeeScalar := big.NewInt(10)
	blobBaseFeeScalar := big.NewInt(5)
	decimals := big.NewInt(6)
	oracleAddress := common.HexToAddress("0x1234").String()

	t.Parallel()

	t.Run("correctly fetches weighted gas price if chain has upgraded to Ecotone", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, false, true)
		mockBatchContractCall(t, ethClient, oracleAddress, baseFee, baseFeeScalar, blobBaseFee, blobBaseFeeScalar, decimals)

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		gasPrice, err := oracle.GetDAGasPrice(tests.Context(t))
		require.NoError(t, err)
		scaledGasPrice := big.NewInt(16125000000) // baseFee * scalar * 16 + blobBaseFee * scalar
		scale := big.NewInt(16000000)             // Scaled by 16 * 10 ^ decimals
		expectedGasPrice := new(big.Int).Div(scaledGasPrice, scale)
		assert.Equal(t, expectedGasPrice, gasPrice)
	})

	t.Run("fetching Ecotone price but rpc returns bad data", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, false, true)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			var badData = "zzz"
			rpcElements[0].Result = &badData
			rpcElements[1].Result = &badData
		}).Return(nil).Once()

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		_, err = oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})

	t.Run("fetching Ecotone price but rpc parent call errors", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, false, true)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Return(fmt.Errorf("revert")).Once()

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		_, err = oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})

	t.Run("fetching Ecotone price but one of the sub rpc call errors", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, false, true)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			res := common.BigToHash(baseFee).Hex()
			rpcElements[0].Result = &res
			rpcElements[1].Error = fmt.Errorf("revert")
		}).Return(nil).Once()

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		_, err = oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})
}

func TestOPL1Oracle_CalculateFjordGasPrice(t *testing.T) {
	baseFee := big.NewInt(100000000)
	blobBaseFee := big.NewInt(25000000)
	baseFeeScalar := big.NewInt(10)
	blobBaseFeeScalar := big.NewInt(5)
	decimals := big.NewInt(6)
	oracleAddress := common.HexToAddress("0x1234").String()

	t.Parallel()

	t.Run("correctly fetches gas price if chain has upgraded to Fjord", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, true, true)
		mockBatchContractCall(t, ethClient, oracleAddress, baseFee, baseFeeScalar, blobBaseFee, blobBaseFeeScalar, decimals)

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		gasPrice, err := oracle.GetDAGasPrice(tests.Context(t))
		require.NoError(t, err)
		scaledGasPrice := big.NewInt(16125000000) // baseFee * scalar * 16 + blobBaseFee * scalar
		scale := big.NewInt(16000000)             // Scaled by 16 * 10 ^ decimals
		expectedGasPrice := new(big.Int).Div(scaledGasPrice, scale)
		assert.Equal(t, expectedGasPrice, gasPrice)
	})

	t.Run("fetching Fjord price but rpc returns bad data", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, true, true)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			var badData = "zzz"
			rpcElements[0].Result = &badData
			rpcElements[1].Result = &badData
		}).Return(nil).Once()

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		_, err = oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})

	t.Run("fetching Fjord price but rpc parent call errors", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, true, true)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Return(fmt.Errorf("revert")).Once()

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		_, err = oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})

	t.Run("fetching Fjord price but one of the sub rpc call errors", func(t *testing.T) {
		ethClient := setupUpgradeCheck(t, oracleAddress, true, true)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			res := common.BigToHash(baseFee).Hex()
			rpcElements[0].Result = &res
			rpcElements[1].Error = fmt.Errorf("revert")
		}).Return(nil).Once()

		oracle, err := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		require.NoError(t, err)
		_, err = oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})
}
