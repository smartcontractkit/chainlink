package ccipexec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	legacyEvmORMMocks "github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
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
		chainSet := &legacyEvmORMMocks.LegacyChainContainer{}
		t.Run(tc.description, func(t *testing.T) {
			err := UnregisterExecPluginLpFilters(context.Background(), logger.TestLogger(t), job.Job{OCR2OracleSpec: tc.spec}, chainSet)
			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
