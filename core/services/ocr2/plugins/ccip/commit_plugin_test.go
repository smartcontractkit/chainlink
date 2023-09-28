package ccip

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	mocklp "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	pipelinemocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestGetCommitPluginFilterNamesFromSpec(t *testing.T) {
	lggr := logger.TestLogger(t)
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
				ContractID:   utils.ZeroAddress.String(),
				PluginConfig: map[string]interface{}{},
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
		{
			description: "valid config",
			spec: &job.OCR2OracleSpec{
				ContractID:   utils.ZeroAddress.String(),
				PluginConfig: map[string]interface{}{},
				RelayConfig: map[string]interface{}{
					"chainID": 1234.0,
				},
			},
			expectingErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			chainSet := &evmmocks.LegacyChainContainer{}
			prMock := &pipelinemocks.Runner{}

			if tc.spec != nil {
				if chainID, ok := tc.spec.RelayConfig["chainID"]; ok {
					chainIdStr := strconv.FormatInt(int64(chainID.(float64)), 10)
					chainSet.On("Get", chainIdStr).
						Return(nil, fmt.Errorf("chain %d not found", chainID))
				}
			}

			err := UnregisterCommitPluginLpFilters(context.Background(), lggr, job.Job{OCR2OracleSpec: tc.spec}, prMock, chainSet)
			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			chainSet.AssertExpectations(t)
		})
	}

}

func TestGetCommitPluginFilterNames(t *testing.T) {
	onRampAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc2")
	priceRegAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc3")
	offRampAddr := common.HexToAddress("0xDAFeA492D9c6733Ae3D56b7eD1AdB60692C98BC4")

	mockCommitStore, _ := testhelpers.NewFakeCommitStore(t, 1)
	mockCommitStore.SetStaticConfig(commit_store.CommitStoreStaticConfig{OnRamp: onRampAddr})
	mockCommitStore.SetDynamicConfig(commit_store.CommitStoreDynamicConfig{PriceRegistry: priceRegAddr})

	srcLP := mocklp.NewLogPoller(t)
	dstLP := mocklp.NewLogPoller(t)

	dstLP.On("UnregisterFilter", "Commit price updates - 0xdafEa492d9C6733aE3D56b7eD1aDb60692c98bc3", mock.Anything).Return(nil)
	dstLP.On("UnregisterFilter", "Fee token added - 0xdafEa492d9C6733aE3D56b7eD1aDb60692c98bc3", mock.Anything).Return(nil)
	dstLP.On("UnregisterFilter", "Fee token removed - 0xdafEa492d9C6733aE3D56b7eD1aDb60692c98bc3", mock.Anything).Return(nil)
	dstLP.On("UnregisterFilter", "Token pool added - 0xDAFeA492D9c6733Ae3D56b7eD1AdB60692C98BC4", mock.Anything).Return(nil)
	dstLP.On("UnregisterFilter", "Token pool removed - 0xDAFeA492D9c6733Ae3D56b7eD1AdB60692C98BC4", mock.Anything).Return(nil)

	err := unregisterCommitPluginFilters(context.Background(), dstLP, mockCommitStore, offRampAddr)
	assert.NoError(t, err)

	srcLP.AssertExpectations(t)
	dstLP.AssertExpectations(t)
}

func Test_updateCommitPluginLogPollerFilters(t *testing.T) {
	srcLP := &mocklp.LogPoller{}
	dstLP := &mocklp.LogPoller{}

	priceRegAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc3")
	offRampAddr := common.HexToAddress("0xDAFeA492D9c6733Ae3D56b7eD1AdB60692C98BC4")
	offRamp := &mock_contracts.EVM2EVMOffRampInterface{}
	offRamp.On("Address").Return(offRampAddr)

	newDestFilters := getCommitPluginDestLpFilters(priceRegAddr, offRampAddr)

	rf := &CommitReportingPluginFactory{
		config: CommitPluginConfig{
			destLP:  dstLP,
			offRamp: offRamp,
		},
		destChainFilters: []logpoller.Filter{
			{Name: "a"},
			{Name: "b"},
		},
		filtersMu: &sync.Mutex{},
	}

	// make sure existing filters get unregistered
	for _, f := range rf.destChainFilters {
		dstLP.On("UnregisterFilter", f.Name, mock.Anything).Return(nil)
	}
	// make sure new filters are registered
	for _, f := range newDestFilters {
		dstLP.On("RegisterFilter", f).Return(nil)
	}
	err := rf.UpdateLogPollerFilters(priceRegAddr)
	assert.NoError(t, err)

	srcLP.AssertExpectations(t)
	dstLP.AssertExpectations(t)
}
