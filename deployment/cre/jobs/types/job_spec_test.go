package job_types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

func TestJobSpecInput_ToStandardCapabilityJob(t *testing.T) {
	t.Parallel()

	jobName := "test-job"

	t.Run("successful conversion", func(t *testing.T) {
		input := job_types.JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{
				Enabled:            true,
				BootstrapPeers:     []string{"peer1", "peer2"},
				OCRContractAddress: "0x123",
				OCRKeyBundleID:     "bundle-id",
				ChainID:            "chain-id",
				TransmitterID:      "transmitter-id",
				OnchainSigningStrategy: pkg.OnchainSigningStrategy{
					StrategyName: "strategy-name",
					Config:       map[string]string{"key": "value"},
				},
			},
		}

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
		input := job_types.JobSpecInput{
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command is required")
	})

	t.Run("invalid command type", func(t *testing.T) {
		input := job_types.JobSpecInput{
			"command":       nil,
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})

	t.Run("config is optional", func(t *testing.T) {
		input := job_types.JobSpecInput{
			"command":       "run",
			"config":        "",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.NoError(t, err)
	})

	t.Run("invalid config type", func(t *testing.T) {
		input := job_types.JobSpecInput{
			"command":       "run",
			"config":        struct{}{},
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unmarshal !!map into string")
	})

	t.Run("invalid externalJobID type", func(t *testing.T) {
		input := job_types.JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": struct{}{},
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unmarshal !!map into string")
	})

	t.Run("invalid oracleFactory type", func(t *testing.T) {
		input := job_types.JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": "not a factory",
		}
		_, err := input.ToStandardCapabilityJob(jobName, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unmarshal !!str")
	})
}
