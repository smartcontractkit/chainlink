package evm

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
	cw, err := NewChainWriterService(lggr, client, txm, ge, chainWriterConfig)
	require.NoError(t, err)

	t.Run("Initialization", func(t *testing.T) {
		t.Run("Fails with invalid ABI", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidAbiConfig := modifyChainWriterConfig(baseConfig, func(cfg *relayevmtypes.ChainWriterConfig) {
				cfg.Contracts["forwarder"].ContractABI = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, ge, invalidAbiConfig)
			require.Error(t, err)
		})

		t.Run("Fails with invalid method names", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidMethodNameConfig := modifyChainWriterConfig(baseConfig, func(cfg *relayevmtypes.ChainWriterConfig) {
				cfg.Contracts["forwarder"].Configs["report"].ChainSpecificName = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, ge, invalidMethodNameConfig)
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
