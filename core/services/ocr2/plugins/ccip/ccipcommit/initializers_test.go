package ccipcommit

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	legacyEvmORMMocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
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
			chainSet := &legacyEvmORMMocks.LegacyChainContainer{}

			if tc.spec != nil {
				if chainID, ok := tc.spec.RelayConfig["chainID"]; ok {
					chainIdStr := strconv.FormatInt(int64(chainID.(float64)), 10)
					chainSet.On("Get", chainIdStr).
						Return(nil, fmt.Errorf("chain %d not found", chainID))
				}
			}

			err := UnregisterCommitPluginLpFilters(context.Background(), lggr, job.Job{OCR2OracleSpec: tc.spec}, chainSet)
			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			chainSet.AssertExpectations(t)
		})
	}

}
