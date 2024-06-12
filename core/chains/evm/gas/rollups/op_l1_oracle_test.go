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

func TestDAPriceReader_ReadV1GasPrice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		isEcotoneError bool
		returnBadData  bool
	}{
		{
			name:           "calling isEcotone returns false, fetches l1BaseFee",
			isEcotoneError: false,
		},
		{
			name:           "calling isEcotone when chain has not made Ecotone upgrade, fetches l1BaseFee",
			isEcotoneError: true,
		},
		{
			name:           "calling isEcotone returns bad data, returns error",
			isEcotoneError: false,
			returnBadData:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l1BaseFee := big.NewInt(100)
			oracleAddress := common.HexToAddress("0x1234").String()

			l1BaseFeeMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
			require.NoError(t, err)
			l1BaseFeeCalldata, err := l1BaseFeeMethodAbi.Pack(OPStackGasOracle_l1BaseFee)
			require.NoError(t, err)

			isEcotoneMethodAbi, err := abi.JSON(strings.NewReader(OPIsEcotoneAbiString))
			require.NoError(t, err)
			isEcotoneCalldata, err := isEcotoneMethodAbi.Pack(OPStackGasOracle_isEcotone)
			require.NoError(t, err)

			ethClient := mocks.NewL1OracleClient(t)
			call := ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
				callMsg := args.Get(1).(ethereum.CallMsg)
				blockNumber := args.Get(2).(*big.Int)
				require.Equal(t, isEcotoneCalldata, callMsg.Data)
				require.Equal(t, oracleAddress, callMsg.To.String())
				assert.Nil(t, blockNumber)
			})

			if tc.returnBadData {
				call.Return([]byte{0x2, 0x2}, nil).Once()
			} else if tc.isEcotoneError {
				call.Return(nil, fmt.Errorf("test error")).Once()
			} else {
				call.Return(isEcotoneMethodAbi.Methods["isEcotone"].Outputs.Pack(false)).Once()
			}

			if !tc.returnBadData {
				ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
					callMsg := args.Get(1).(ethereum.CallMsg)
					blockNumber := args.Get(2).(*big.Int)
					require.Equal(t, l1BaseFeeCalldata, callMsg.Data)
					require.Equal(t, oracleAddress, callMsg.To.String())
					assert.Nil(t, blockNumber)
				}).Return(common.BigToHash(l1BaseFee).Bytes(), nil).Once()
			}

			oracle := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
			gasPrice, err := oracle.GetDAGasPrice(tests.Context(t))

			if tc.returnBadData {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, l1BaseFee, gasPrice)
			}
		})
	}
}

func setupIsEcotone(t *testing.T, oracleAddress string) *mocks.L1OracleClient {
	isEcotoneMethodAbi, err := abi.JSON(strings.NewReader(OPIsEcotoneAbiString))
	require.NoError(t, err)
	isEcotoneCalldata, err := isEcotoneMethodAbi.Pack(OPStackGasOracle_isEcotone)
	require.NoError(t, err)

	ethClient := mocks.NewL1OracleClient(t)
	ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
		callMsg := args.Get(1).(ethereum.CallMsg)
		blockNumber := args.Get(2).(*big.Int)
		require.Equal(t, isEcotoneCalldata, callMsg.Data)
		require.Equal(t, oracleAddress, callMsg.To.String())
		assert.Nil(t, blockNumber)
	}).Return(isEcotoneMethodAbi.Methods["isEcotone"].Outputs.Pack(true)).Once()

	return ethClient
}

func TestDAPriceReader_ReadEcotoneGasPrice(t *testing.T) {
	l1BaseFee := big.NewInt(100)
	oracleAddress := common.HexToAddress("0x1234").String()

	t.Parallel()

	t.Run("correctly fetches weighted gas price if chain has upgraded to Ecotone", func(t *testing.T) {
		ethClient := setupIsEcotone(t, oracleAddress)
		getL1GasUsedMethodAbi, err := abi.JSON(strings.NewReader(OPGetL1GasUsedAbiString))
		require.NoError(t, err)
		getL1GasUsedCalldata, err := getL1GasUsedMethodAbi.Pack(OPStackGasOracle_getL1GasUsed, []byte{0x1})
		require.NoError(t, err)

		getL1FeeMethodAbi, err := abi.JSON(strings.NewReader(GetL1FeeAbiString))
		require.NoError(t, err)
		getL1FeeCalldata, err := getL1FeeMethodAbi.Pack(OPStackGasOracle_getL1Fee, []byte{0x1})
		require.NoError(t, err)

		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			require.Equal(t, 2, len(rpcElements))

			for _, rE := range rpcElements {
				require.Equal(t, "eth_call", rE.Method)
				require.Equal(t, oracleAddress, rE.Args[0].(map[string]interface{})["to"])
				require.Equal(t, "latest", rE.Args[1])
			}

			require.Equal(t, hexutil.Bytes(getL1GasUsedCalldata), rpcElements[0].Args[0].(map[string]interface{})["data"])
			require.Equal(t, hexutil.Bytes(getL1FeeCalldata), rpcElements[1].Args[0].(map[string]interface{})["data"])

			res1 := common.BigToHash(big.NewInt(1)).Hex()
			res2 := common.BigToHash(l1BaseFee).Hex()
			rpcElements[0].Result = &res1
			rpcElements[1].Result = &res2
		}).Return(nil).Once()

		oracle := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		gasPrice, err := oracle.GetDAGasPrice(tests.Context(t))
		require.NoError(t, err)
		assert.Equal(t, l1BaseFee, gasPrice)
	})

	t.Run("fetching Ecotone price but rpc returns bad data", func(t *testing.T) {
		ethClient := setupIsEcotone(t, oracleAddress)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			var badData = "zzz"
			rpcElements[0].Result = &badData
			rpcElements[1].Result = &badData
		}).Return(nil).Once()

		oracle := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		_, err := oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})

	t.Run("fetching Ecotone price but rpc parent call errors", func(t *testing.T) {
		ethClient := setupIsEcotone(t, oracleAddress)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Return(fmt.Errorf("revert")).Once()

		oracle := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		_, err := oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})

	t.Run("fetching Ecotone price but one of the sub rpc call errors", func(t *testing.T) {
		ethClient := setupIsEcotone(t, oracleAddress)
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			res := common.BigToHash(l1BaseFee).Hex()
			rpcElements[0].Result = &res
			rpcElements[1].Error = fmt.Errorf("revert")
		}).Return(nil).Once()

		oracle := newOpStackL1GasOracle(logger.Test(t), ethClient, chaintype.ChainOptimismBedrock, oracleAddress)
		_, err := oracle.GetDAGasPrice(tests.Context(t))
		assert.Error(t, err)
	})
}
