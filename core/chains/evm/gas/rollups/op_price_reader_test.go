package rollups

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestDAPriceReader_ReadV1GasPrice(t *testing.T) {
	testCases := []struct {
		name           string
		isEcotoneError bool
		returnBadData  bool
	}{
		{
			name:           "calling isEcotone returns false",
			isEcotoneError: false,
		},
		{
			name:           "calling isEcotone when chain has not made Ecotone upgrade",
			isEcotoneError: true,
		},
		{
			name:           "calling isEcotone returns bad data",
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

			ethClient := mocks.NewETHClient(t)
			call := ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
				callMsg := args.Get(1).(ethereum.CallMsg)
				blockNumber := args.Get(2).(*big.Int)
				require.Equal(t, isEcotoneCalldata, callMsg.Data)
				require.Equal(t, oracleAddress, callMsg.To.Hex())
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
					require.Equal(t, oracleAddress, callMsg.To.Hex())
					assert.Nil(t, blockNumber)
				}).Return(common.BigToHash(l1BaseFee).Bytes(), nil).Once()
			}

			oracle := newOPPriceReader(logger.Test(t), ethClient, config.ChainOptimismBedrock, oracleAddress)
			gasPrice, err := oracle.GetDAGasPrice(testutils.Context(t))

			if tc.returnBadData {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, l1BaseFee, gasPrice)
			}
		})
	}

	t.Run("Calling l1BaseFee when chain has not made Ecotone upgrade", func(t *testing.T) {
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

		ethClient := mocks.NewETHClient(t)
		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			require.Equal(t, isEcotoneCalldata, callMsg.Data)
			require.Equal(t, oracleAddress, callMsg.To.Hex())
			assert.Nil(t, blockNumber)
		}).Return(isEcotoneMethodAbi.Methods["isEcotone"].Outputs.Pack(false)).Once()

		ethClient.On("CallContract", mock.Anything, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).Run(func(args mock.Arguments) {
			callMsg := args.Get(1).(ethereum.CallMsg)
			blockNumber := args.Get(2).(*big.Int)
			require.Equal(t, l1BaseFeeCalldata, callMsg.Data)
			require.Equal(t, oracleAddress, callMsg.To.Hex())
			assert.Nil(t, blockNumber)
		}).Return(common.BigToHash(l1BaseFee).Bytes(), nil).Once()

		oracle := newOPPriceReader(logger.Test(t), ethClient, config.ChainOptimismBedrock, oracleAddress)

		gasPrice, err := oracle.GetDAGasPrice(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, l1BaseFee, gasPrice)
	})
}
