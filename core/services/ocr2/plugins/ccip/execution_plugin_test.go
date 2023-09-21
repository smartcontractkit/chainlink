package ccip

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocklp "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/usdc"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestGetExecutionPluginFilterNamesFromSpec(t *testing.T) {
	testCases := []struct {
		description  string
		spec         *job.OCR2OracleSpec
		expectingErr bool
	}{
		{
			description:  "should not panic with nil spec",
			spec:         nil,
			expectingErr: true,
		},
		{
			description: "invalid config",
			spec: &job.OCR2OracleSpec{
				PluginConfig: map[string]interface{}{},
			},
			expectingErr: true,
		},
		{
			description: "invalid off ramp address",
			spec: &job.OCR2OracleSpec{
				PluginConfig: map[string]interface{}{"offRamp": "123"},
			},
			expectingErr: true,
		},
		{
			description: "invalid contract id",
			spec: &job.OCR2OracleSpec{
				ContractID: "whatever...",
			},
			expectingErr: true,
		},
	}

	for _, tc := range testCases {
		chainSet := &mocks.LegacyChainContainer{}
		t.Run(tc.description, func(t *testing.T) {
			err := UnregisterExecPluginLpFilters(context.Background(), logger.TestLogger(t), tc.spec, chainSet)
			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetExecutionPluginFilterNames(t *testing.T) {
	commitStoreAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc3")
	srcPriceRegAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc9")
	dstPriceRegAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98b19")

	mockOffRamp, offRampAddr := testhelpers.NewFakeOffRamp(t)
	mockOffRamp.SetDynamicConfig(evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig{PriceRegistry: dstPriceRegAddr})

	mockOnRamp, onRampAddr := testhelpers.NewFakeOnRamp(t)
	mockOnRamp.SetDynamicCfg(evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{PriceRegistry: srcPriceRegAddr})

	pluginConfig := config.ExecutionPluginJobSpecConfig{
		USDCConfig: config.USDCConfig{
			SourceTokenAddress:              utils.RandomAddress(),
			SourceMessageTransmitterAddress: utils.RandomAddress(),
			AttestationAPI:                  "http://localhost:8080",
		},
	}

	srcLP := mocklp.NewLogPoller(t)
	srcFilters := []string{
		"Exec ccip sends - " + onRampAddr.String(),
		"Fee token added - 0xdAFea492D9c6733aE3d56B7ed1ADb60692c98bC9",
		"Fee token removed - 0xdAFea492D9c6733aE3d56B7ed1ADb60692c98bC9",
		usdc.MESSAGE_SENT_FILTER_NAME + " - " + pluginConfig.USDCConfig.SourceMessageTransmitterAddress.Hex(),
	}
	for _, f := range srcFilters {
		srcLP.On("UnregisterFilter", f, mock.Anything).Return(nil)
	}

	dstLP := mocklp.NewLogPoller(t)
	dstFilters := []string{
		"Exec report accepts - 0xdafEa492d9C6733aE3D56b7eD1aDb60692c98bc3",
		"Exec execution state changes - " + offRampAddr.String(),
		"Token pool added - " + offRampAddr.String(),
		"Token pool removed - " + offRampAddr.String(),
		"Fee token added - 0xdaFEa492D9C6733Ae3D56b7ed1adB60692C98b19",
		"Fee token removed - 0xdaFEa492D9C6733Ae3D56b7ed1adB60692C98b19",
	}
	for _, f := range dstFilters {
		dstLP.On("UnregisterFilter", f, mock.Anything).Return(nil)
	}

	err := unregisterExecutionPluginLpFilters(
		context.Background(),
		logger.TestLogger(t),
		srcLP,
		dstLP,
		mockOffRamp,
		evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore:         commitStoreAddr,
			OnRamp:              onRampAddr,
			SourceChainSelector: 5009297550715157269,
		},
		mockOnRamp,
		nil,
		pluginConfig,
	)
	assert.NoError(t, err)

	srcLP.AssertExpectations(t)
	dstLP.AssertExpectations(t)
}
