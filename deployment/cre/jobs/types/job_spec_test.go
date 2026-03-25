package job_types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

func mustMarshal(t *testing.T, v any) job_types.JobSpecInput {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return job_types.JobSpecInput{RawMessage: b}
}

func TestJobSpecInput_ToStandardCapabilityJob(t *testing.T) {
	t.Parallel()

	jobName := "test-job"

	t.Run("successful conversion", func(t *testing.T) {
		input := mustMarshal(t, map[string]any{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": map[string]any{
				"enabled":            true,
				"bootstrapPeers":     []string{"peer1", "peer2"},
				"ocrContractAddress": "0x123",
				"ocrKeyBundleID":     "bundle-id",
				"chainID":            "chain-id",
				"transmitterID":      "transmitter-id",
				"onchainSigningStrategy": map[string]any{
					"strategyName": "strategy-name",
					"config":       map[string]string{"key": "value"},
				},
			},
		})

		job, err := input.ToStandardCapabilityJob(jobName, false)
		require.NoError(t, err)
		assert.Equal(t, jobName, job.JobName)
		assert.Equal(t, "run", job.Command)
		assert.Equal(t, "param=value", job.Config)
		assert.Equal(t, "123", job.ExternalJobID)
		assert.True(t, job.OracleFactory.Enabled)
		assert.Equal(t, []string{"peer1", "peer2"}, job.OracleFactory.BootstrapPeers)
		assert.Equal(t, "0x123", job.OracleFactory.OCRContractAddress)
		assert.Equal(t, "bundle-id", job.OracleFactory.OCRKeyBundleID)
		assert.Equal(t, "chain-id", job.OracleFactory.ChainID)
		assert.Equal(t, "transmitter-id", job.OracleFactory.TransmitterID)
		assert.Equal(t, "strategy-name", job.OracleFactory.OnchainSigningStrategy.StrategyName)
		assert.Equal(t, map[string]string{"key": "value"}, job.OracleFactory.OnchainSigningStrategy.Config)
	})

	t.Run("missing command", func(t *testing.T) {
		input := mustMarshal(t, map[string]any{
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		})
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command is required")
	})

	t.Run("missing command field entirely", func(t *testing.T) {
		input := mustMarshal(t, map[string]any{
			"config":        "param=value",
			"externalJobID": "123",
		})
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})

	t.Run("config is optional", func(t *testing.T) {
		input := mustMarshal(t, map[string]any{
			"command":       "run",
			"config":        "",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		})
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.NoError(t, err)
	})
}
