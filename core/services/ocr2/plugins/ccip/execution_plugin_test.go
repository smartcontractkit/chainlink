package ccip

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocklp "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
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
			err := UnregisterExecPluginLpFilters(context.Background(), tc.spec, chainSet)
			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetExecutionPluginFilterNames(t *testing.T) {
	specContractID := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc1") // off-ramp addr
	onRampAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc2")
	commitStoreAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc3")
	srcPriceRegAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc9")
	dstPriceRegAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98b19")

	mockOffRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	mockOffRamp.On("Address").Return(specContractID)
	mockOffRamp.On("GetDynamicConfig", mock.Anything).Return(
		evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig{
			PriceRegistry: dstPriceRegAddr,
		}, nil)

	mockOnRamp := mock_contracts.NewEVM2EVMOnRampInterface(t)
	mockOnRamp.On("TypeAndVersion", mock.Anything).Return(fmt.Sprintf("%s %s", ccipconfig.EVM2EVMOnRamp, "1.2.0"), nil)
	mockOnRamp.On("GetDynamicConfig", mock.Anything).Return(
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			PriceRegistry: srcPriceRegAddr,
		}, nil)

	srcLP := mocklp.NewLogPoller(t)
	srcFilters := []string{
		"Exec ccip sends - 0xdafea492D9c6733aE3d56B7ED1aDb60692C98bc2",
		"Fee token added - 0xdAFea492D9c6733aE3d56B7ed1ADb60692c98bC9",
		"Fee token removed - 0xdAFea492D9c6733aE3d56B7ed1ADb60692c98bC9",
	}
	for _, f := range srcFilters {
		srcLP.On("UnregisterFilter", f, mock.Anything).Return(nil)
	}

	dstLP := mocklp.NewLogPoller(t)
	dstFilters := []string{
		"Exec report accepts - 0xdafEa492d9C6733aE3D56b7eD1aDb60692c98bc3",
		"Exec execution state changes - 0xdafeA492d9c6733Ae3d56B7ed1AdB60692C98bC1",
		"Token pool added - 0xdafeA492d9c6733Ae3d56B7ed1AdB60692C98bC1",
		"Token pool removed - 0xdafeA492d9c6733Ae3d56B7ed1AdB60692C98bC1",
		"Fee token added - 0xdaFEa492D9C6733Ae3D56b7ed1adB60692C98b19",
		"Fee token removed - 0xdaFEa492D9C6733Ae3D56b7ed1adB60692C98b19",
	}
	for _, f := range dstFilters {
		dstLP.On("UnregisterFilter", f, mock.Anything).Return(nil)
	}

	err := unregisterExecutionPluginLpFilters(
		context.Background(),
		srcLP,
		dstLP,
		mockOffRamp,
		evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore: commitStoreAddr,
			OnRamp:      onRampAddr,
		},
		mockOnRamp,
		nil,
	)
	assert.NoError(t, err)

	srcLP.AssertExpectations(t)
	dstLP.AssertExpectations(t)
}
