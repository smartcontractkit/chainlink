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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	gasmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/mocks"
	rollupmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/rollups/mocks"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
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
			invalidAbiConfig := modifyChainWriterConfig(baseConfig, func(cfg *config.ChainWriterConfig) {
				cfg.Contracts["forwarder"].ContractABI = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, ge, invalidAbiConfig, nil)
			require.Error(t, err)
		})

		t.Run("Fails with invalid method names", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidMethodNameConfig := modifyChainWriterConfig(baseConfig, func(cfg *config.ChainWriterConfig) {
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
			status types.TransactionStatus
		}{
			{uuid.NewString(), types.Unknown},
			{uuid.NewString(), types.Pending},
			{uuid.NewString(), types.Unconfirmed},
			{uuid.NewString(), types.Finalized},
			{uuid.NewString(), types.Failed},
			{uuid.NewString(), types.Fatal},
		}

		for _, tx := range txs {
			txm.On("GetTransactionStatus", mock.Anything, tx.txid).Return(tx.status, nil).Once()
		}

		for _, tx := range txs {
			var status types.TransactionStatus
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
	})
}

// Helper functions to remove redundant creation of configs
func newBaseChainWriterConfig() config.ChainWriterConfig {
	return config.ChainWriterConfig{
		Contracts: map[string]*config.ContractConfig{
			"forwarder": {
				// TODO: Use generic ABI / test contract rather than a keystone specific one
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*config.ChainWriterDefinition{
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

func modifyChainWriterConfig(baseConfig config.ChainWriterConfig, modifyFn func(*config.ChainWriterConfig)) config.ChainWriterConfig {
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
