package evm

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonLggr "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	rollupmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type MockGasEstimatorConfig struct {
	EIP1559DynamicFeesF bool
	BumpPercentF        uint16
	BumpThresholdF      uint64
	BumpMinF            *assets.Wei
	LimitMultiplierF    float32
	TipCapDefaultF      *assets.Wei
	TipCapMinF          *assets.Wei
	PriceMaxF           *assets.Wei
	PriceMinF           *assets.Wei
	PriceDefaultF       *assets.Wei
	FeeCapDefaultF      *assets.Wei
	LimitMaxF           uint64
	ModeF               string
}

func TestChainWriter(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	txm := txmmocks.NewMockEvmTxManager(t)
	client := evmclimocks.NewClient(t)
	ge := gasmocks.NewEvmEstimator(t)

	geCfg := gas.NewMockGasConfig()

	getEst := func(commonLggr.Logger) gas.EvmEstimator { return ge }

	feeEstimator := gas.NewEvmFeeEstimator(commonLggr.Test(t), getEst, true, geCfg)
	l1Oracle := rollupmocks.NewL1Oracle(t)

	chainWriterConfig := newBaseChainWriterConfig()
	cw, err := NewChainWriterService(lggr, client, txm, feeEstimator, chainWriterConfig)

	require.NoError(t, err)

	t.Run("Initialization", func(t *testing.T) {
		t.Run("Fails with invalid ABI", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidAbiConfig := modifyChainWriterConfig(baseConfig, func(cfg *relayevmtypes.ChainWriterConfig) {
				cfg.Contracts["forwarder"].ContractABI = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, feeEstimator, invalidAbiConfig)
			require.Error(t, err)
		})

		t.Run("Fails with invalid method names", func(t *testing.T) {
			baseConfig := newBaseChainWriterConfig()
			invalidMethodNameConfig := modifyChainWriterConfig(baseConfig, func(cfg *relayevmtypes.ChainWriterConfig) {
				cfg.Contracts["forwarder"].Configs["report"].ChainSpecificName = ""
			})
			_, err = NewChainWriterService(lggr, client, txm, feeEstimator, invalidMethodNameConfig)
			require.Error(t, err)
		})
	})

	t.Run("SubmitTransaction", func(t *testing.T) {
		// TODO: implement
	})

	t.Run("GetFeeComponents", func(t *testing.T) {
		ge.On("GetDynamicFee", mock.Anything, mock.Anything).Return(gas.DynamicFee{
			FeeCap: assets.NewWei(big.NewInt(1000000002)),
			TipCap: assets.NewWei(big.NewInt(1000000003)),
		}, nil).Twice()

		l1Oracle.On("GasPrice", mock.Anything).Return(assets.NewWei(big.NewInt(1000000004)), nil).Once()
		ge.On("L1Oracle", mock.Anything).Return(l1Oracle).Once()
		var feeComponents *types.ChainFeeComponents
		t.Run("Returns valid FeeComponents", func(t *testing.T) {
			feeComponents, err = cw.GetFeeComponents(ctx)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(1000000002), &feeComponents.ExecutionFee)
			assert.Equal(t, big.NewInt(1000000004), &feeComponents.DataAvailabilityFee)
		})

		ge.On("L1Oracle", mock.Anything).Return(nil).Twice()

		t.Run("Returns valid FeeComponents with no L1Oracle", func(t *testing.T) {
			feeComponents, err = cw.GetFeeComponents(ctx)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(1000000002), &feeComponents.ExecutionFee)
			assert.Equal(t, big.NewInt(0), &feeComponents.DataAvailabilityFee)
		})

		t.Run("Returns Legacy Fee for non-EIP1559 enabled gas estimator", func(t *testing.T) {
			noDynamicFeeEstimator := gas.NewEvmFeeEstimator(commonLggr.Test(t), getEst, false, geCfg)
			var noDyanmicCW ChainWriterService
			noDyanmicCW, err = NewChainWriterService(lggr, client, txm, noDynamicFeeEstimator, chainWriterConfig)
			ge.On("GetLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assets.NewWei(big.NewInt(1000000001)), uint64(0), nil).Once()

			feeComponents, err = noDyanmicCW.GetFeeComponents(ctx)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(1000000001), &feeComponents.ExecutionFee)
			assert.Equal(t, big.NewInt(0), &feeComponents.DataAvailabilityFee)
		})

		t.Run("Fails when neither legacy or dynamic fee is available", func(t *testing.T) {
			ge.On("GetDynamicFee", mock.Anything, mock.Anything).Return(gas.DynamicFee{
				FeeCap: nil,
				TipCap: nil,
			}, nil).Once()

			_, err = cw.GetFeeComponents(ctx)
			require.Error(t, err)
		})

		t.Run("Fails when GetFee returns an error", func(t *testing.T) {
			expectedErr := fmt.Errorf("GetFee error")
			ge.On("GetDynamicFee", mock.Anything, mock.Anything).Return(gas.DynamicFee{
				FeeCap: nil,
				TipCap: nil,
			}, expectedErr).Once()
			_, err = cw.GetFeeComponents(ctx)
			require.Equal(t, expectedErr, err)
		})

		t.Run("Fails when L1Oracle returns error", func(t *testing.T) {
			ge.On("GetDynamicFee", mock.Anything, mock.Anything).Return(gas.DynamicFee{
				FeeCap: assets.NewWei(big.NewInt(1000000002)),
				TipCap: assets.NewWei(big.NewInt(1000000003)),
			}, nil).Once()

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
		MaxGasPrice: big.NewInt(1000000000000),
	}
}

func modifyChainWriterConfig(baseConfig relayevmtypes.ChainWriterConfig, modifyFn func(*relayevmtypes.ChainWriterConfig)) relayevmtypes.ChainWriterConfig {
	modifiedConfig := baseConfig
	modifyFn(&modifiedConfig)
	return modifiedConfig
}
