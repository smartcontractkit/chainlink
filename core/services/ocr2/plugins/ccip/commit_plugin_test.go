package ccip

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
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
