package evm

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	gasmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/mocks"
	rollupmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/rollups/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestChainWriter(t *testing.T) {
	lggr := logger.Test(t)
	ctx := testutils.Context(t)

	txm := txmmocks.NewMockEvmTxManager(t)
	client := clienttest.NewClient(t)
	ge := gasmocks.NewEvmFeeEstimator(t)
	l1Oracle := rollupmocks.NewL1Oracle(t)

	chainWriterConfig := newBaseChainWriterConfig()
	cw, err := NewChainWriterService(lggr, client, txm, ge, chainWriterConfig, nil)
	require.NoError(t, err)

	t.Run("Initialization", func(t *testing.T) {
		t.Run("Fails with invalid ABI", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidAbiConfig := modifyChainWriterConfig(baseConfig, func(cfg *relayevmtypes.ChainWriterConfig) {
				cfg.Contracts["forwarder"].ContractABI = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, ge, invalidAbiConfig, nil)
			require.Error(t, err)
		})

		t.Run("Fails with invalid method names", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidMethodNameConfig := modifyChainWriterConfig(baseConfig, func(cfg *relayevmtypes.ChainWriterConfig) {
				cfg.Contracts["forwarder"].Configs["report"].ChainSpecificName = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, ge, invalidMethodNameConfig, nil)
			require.Error(t, err)
		})
	})

	t.Run("SubmitTransaction", func(t *testing.T) {
		// TODO: implement
	})

	t.Run("GetTransactionStatus", func(t *testing.T) {
		txs := []struct {
			txid   string
			status commontypes.TransactionStatus
		}{
			{uuid.NewString(), commontypes.Unknown},
			{uuid.NewString(), commontypes.Pending},
			{uuid.NewString(), commontypes.Unconfirmed},
			{uuid.NewString(), commontypes.Finalized},
			{uuid.NewString(), commontypes.Failed},
			{uuid.NewString(), commontypes.Fatal},
		}

		for _, tx := range txs {
			txm.On("GetTransactionStatus", mock.Anything, tx.txid).Return(tx.status, nil).Once()
		}

		for _, tx := range txs {
			var status commontypes.TransactionStatus
			status, err = cw.GetTransactionStatus(ctx, tx.txid)
			require.NoError(t, err)
			require.Equal(t, tx.status, status)
		}
	})

	t.Run("GetFeeComponents", func(t *testing.T) {
		ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{
			GasPrice:   assets.NewWei(big.NewInt(1000000001)),
			DynamicFee: gas.DynamicFee{GasFeeCap: assets.NewWei(big.NewInt(1000000002)), GasTipCap: assets.NewWei(big.NewInt(1000000003))},
		}, uint64(0), nil).Twice()

		l1Oracle.On("GasPrice", mock.Anything).Return(assets.NewWei(big.NewInt(1000000004)), nil).Once()
		ge.On("L1Oracle", mock.Anything).Return(l1Oracle).Once()
		var feeComponents *types.ChainFeeComponents
		t.Run("Returns valid FeeComponents", func(t *testing.T) {
			feeComponents, err = cw.GetFeeComponents(ctx)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(1000000002), feeComponents.ExecutionFee)
			assert.Equal(t, big.NewInt(1000000004), feeComponents.DataAvailabilityFee)
		})

		ge.On("L1Oracle", mock.Anything).Return(nil).Twice()

		t.Run("Returns valid FeeComponents with no L1Oracle", func(t *testing.T) {
			feeComponents, err = cw.GetFeeComponents(ctx)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(1000000002), feeComponents.ExecutionFee)
			assert.Equal(t, big.NewInt(0), feeComponents.DataAvailabilityFee)
		})

		t.Run("Returns Legacy Fee in absence of Dynamic Fee", func(t *testing.T) {
			ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{
				GasPrice:   assets.NewWei(big.NewInt(1000000001)),
				DynamicFee: gas.DynamicFee{GasFeeCap: nil, GasTipCap: assets.NewWei(big.NewInt(1000000003))},
			}, uint64(0), nil).Once()
			feeComponents, err = cw.GetFeeComponents(ctx)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(1000000001), feeComponents.ExecutionFee)
			assert.Equal(t, big.NewInt(0), feeComponents.DataAvailabilityFee)
		})

		t.Run("Fails when neither legacy or dynamic fee is available", func(t *testing.T) {
			ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{
				GasPrice:   nil,
				DynamicFee: gas.DynamicFee{},
			}, uint64(0), nil).Once()

			_, err = cw.GetFeeComponents(ctx)
			require.Error(t, err)
		})

		t.Run("Fails when GetFee returns an error", func(t *testing.T) {
			expectedErr := fmt.Errorf("GetFee error")
			ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{
				GasPrice:   nil,
				DynamicFee: gas.DynamicFee{},
			}, uint64(0), expectedErr).Once()
			_, err = cw.GetFeeComponents(ctx)
			require.Equal(t, expectedErr, err)
		})

		t.Run("Fails when L1Oracle returns error", func(t *testing.T) {
			ge.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{
				GasPrice:   assets.NewWei(big.NewInt(1000000001)),
				DynamicFee: gas.DynamicFee{GasFeeCap: assets.NewWei(big.NewInt(1000000002)), GasTipCap: assets.NewWei(big.NewInt(1000000003))},
			}, uint64(0), nil).Once()
			ge.On("L1Oracle", mock.Anything).Return(l1Oracle).Once()

			expectedErr := fmt.Errorf("l1Oracle error")
			l1Oracle.On("GasPrice", mock.Anything).Return(nil, expectedErr).Once()
			_, err = cw.GetFeeComponents(ctx)
			require.Equal(t, expectedErr, err)
		})
	})

	t.Run("Tron Conversion Methods", func(t *testing.T) {
		t.Run("buildMethodSignature", func(t *testing.T) {
			t.Run("Single parameter method", func(t *testing.T) {
				abiMethod := createTestABIMethod("transfer", []string{"address"})
				signature := cw.(*chainWriter).buildMethodSignature(abiMethod)
				assert.Equal(t, "transfer(address)", signature)
			})

			t.Run("Multiple parameter method", func(t *testing.T) {
				abiMethod := createTestABIMethod("mint", []string{"address", "uint256"})
				signature := cw.(*chainWriter).buildMethodSignature(abiMethod)
				assert.Equal(t, "mint(address,uint256)", signature)
			})

			t.Run("No parameter method", func(t *testing.T) {
				abiMethod := createTestABIMethod("pause", []string{})
				signature := cw.(*chainWriter).buildMethodSignature(abiMethod)
				assert.Equal(t, "pause()", signature)
			})

			t.Run("Complex types method", func(t *testing.T) {
				abiMethod := createTestABIMethod("complexMethod", []string{"bytes32", "bool", "uint256[]"})
				signature := cw.(*chainWriter).buildMethodSignature(abiMethod)
				assert.Equal(t, "complexMethod(bytes32,bool,uint256[])", signature)
			})
		})

		t.Run("convertArgsToTronParams", func(t *testing.T) {
			t.Run("Slice arguments", func(t *testing.T) {
				abiMethod := createTestABIMethod("mint", []string{"address", "uint256"})
				args := []any{
					"0x1234567890123456789012345678901234567890",
					big.NewInt(1000),
				}

				params, err := cw.(*chainWriter).convertArgsToTronParams(abiMethod, args)
				require.NoError(t, err)

				expected := []any{"address", "0x1234567890123456789012345678901234567890", "uint256", "1000"}
				assert.Equal(t, expected, params)
			})

			t.Run("Struct arguments", func(t *testing.T) {
				abiMethod := createTestABIMethod("transfer", []string{"address", "uint256", "bool"})
				args := struct {
					To     string
					Amount *big.Int
					Active bool
				}{
					To:     "0x1234567890123456789012345678901234567890",
					Amount: big.NewInt(500),
					Active: true,
				}

				params, err := cw.(*chainWriter).convertArgsToTronParams(abiMethod, args)
				require.NoError(t, err)

				expected := []any{"address", "0x1234567890123456789012345678901234567890", "uint256", "500", "bool", "true"}
				assert.Equal(t, expected, params)
			})

			t.Run("Argument count mismatch in slice", func(t *testing.T) {
				abiMethod := createTestABIMethod("mint", []string{"address", "uint256"})
				args := []any{"0x1234567890123456789012345678901234567890"}

				_, err := cw.(*chainWriter).convertArgsToTronParams(abiMethod, args)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "argument count mismatch")
			})

			t.Run("Field count mismatch in struct", func(t *testing.T) {
				abiMethod := createTestABIMethod("mint", []string{"address", "uint256"})
				args := struct {
					To     string
					Amount *big.Int
					Extra  bool
				}{
					To:     "0x1234567890123456789012345678901234567890",
					Amount: big.NewInt(500),
					Extra:  true,
				}

				_, err := cw.(*chainWriter).convertArgsToTronParams(abiMethod, args)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "struct field count mismatch")
			})

			t.Run("Unsupported argument type", func(t *testing.T) {
				abiMethod := createTestABIMethod("mint", []string{"address"})
				args := "invalid"

				_, err := cw.(*chainWriter).convertArgsToTronParams(abiMethod, args)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported args type")
			})
		})

		t.Run("convertValueToString", func(t *testing.T) {
			t.Run("Address types", func(t *testing.T) {
				addressType := createTestABIType("address")

				result, err := cw.(*chainWriter).convertValueToString(common.HexToAddress("0x1234567890123456789012345678901234567890"), addressType)
				require.NoError(t, err)
				assert.Equal(t, "0x1234567890123456789012345678901234567890", result)

				result, err = cw.(*chainWriter).convertValueToString("0x1234567890123456789012345678901234567890", addressType)
				require.NoError(t, err)
				assert.Equal(t, "0x1234567890123456789012345678901234567890", result)

				_, err = cw.(*chainWriter).convertValueToString(123, addressType)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid address type")
			})

			t.Run("Integer types", func(t *testing.T) {
				uint256Type := createTestABIType("uint256")

				result, err := cw.(*chainWriter).convertValueToString(big.NewInt(12345), uint256Type)
				require.NoError(t, err)
				assert.Equal(t, "12345", result)

				result, err = cw.(*chainWriter).convertValueToString(int64(999), uint256Type)
				require.NoError(t, err)
				assert.Equal(t, "999", result)

				result, err = cw.(*chainWriter).convertValueToString(uint32(777), uint256Type)
				require.NoError(t, err)
				assert.Equal(t, "777", result)
			})

			t.Run("String types", func(t *testing.T) {
				stringType := createTestABIType("string")

				result, err := cw.(*chainWriter).convertValueToString("hello world", stringType)
				require.NoError(t, err)
				assert.Equal(t, "hello world", result)

				result, err = cw.(*chainWriter).convertValueToString(123, stringType)
				require.NoError(t, err)
				assert.Equal(t, "123", result)
			})

			t.Run("Boolean types", func(t *testing.T) {
				boolType := createTestABIType("bool")

				result, err := cw.(*chainWriter).convertValueToString(true, boolType)
				require.NoError(t, err)
				assert.Equal(t, "true", result)

				result, err = cw.(*chainWriter).convertValueToString(false, boolType)
				require.NoError(t, err)
				assert.Equal(t, "false", result)

				result, err = cw.(*chainWriter).convertValueToString("true", boolType)
				require.NoError(t, err)
				assert.Equal(t, "true", result)
			})

			t.Run("Bytes types", func(t *testing.T) {
				bytesType := createTestABIType("bytes")

				result, err := cw.(*chainWriter).convertValueToString([]byte{0x12, 0x34, 0xab}, bytesType)
				require.NoError(t, err)
				assert.Equal(t, "0x1234ab", result)

				result, err = cw.(*chainWriter).convertValueToString("0x1234ab", bytesType)
				require.NoError(t, err)
				assert.Equal(t, "0x1234ab", result)
			})

			t.Run("Unsupported types fallback", func(t *testing.T) {
				functionType := createTestABIType("function")

				result, err := cw.(*chainWriter).convertValueToString("some_function", functionType)
				require.NoError(t, err)
				assert.Equal(t, "some_function", result)
			})
		})
	})
}

// Helper functions to remove redundant creation of configs
func newBaseChainWriterConfig() relayevmtypes.ChainWriterConfig {
	return relayevmtypes.ChainWriterConfig{
		Contracts: map[string]*relayevmtypes.ContractConfig{
			"forwarder": {
				// TODO: Use generic ABI / test contract rather than a keystone specific one
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainWriterDefinition{
					"report": {
						ChainSpecificName: "report",
						Checker:           "simulate",
						FromAddress:       testutils.NewAddress(),
						GasLimit:          200_000,
					},
				},
			},
		},
		MaxGasPrice: assets.NewWeiI(1000000000000),
	}
}

func modifyChainWriterConfig(baseConfig relayevmtypes.ChainWriterConfig, modifyFn func(*relayevmtypes.ChainWriterConfig)) relayevmtypes.ChainWriterConfig {
	modifiedConfig := baseConfig
	modifyFn(&modifiedConfig)
	return modifiedConfig
}

func createTestABIMethod(name string, params []string) abi.Method {
	var inputs abi.Arguments
	for _, param := range params {
		abiType, _ := abi.NewType(param, "", nil)
		inputs = append(inputs, abi.Argument{
			Type: abiType,
		})
	}

	return abi.Method{
		Name:   name,
		Inputs: inputs,
	}
}

func createTestABIType(typeStr string) abi.Type {
	abiType, _ := abi.NewType(typeStr, "", nil)
	return abiType
}
