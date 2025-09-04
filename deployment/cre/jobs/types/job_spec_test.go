package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

func TestJobSpecInput_ToStandardCapabilityJob(t *testing.T) {
	t.Parallel()

	jobName := "test-job"

	t.Run("successful conversion", func(t *testing.T) {
		input := JobSpecInput{
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

		job, err := input.ToStandardCapabilityJob(jobName)
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
		input := JobSpecInput{
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command is required")
	})

	t.Run("invalid command type", func(t *testing.T) {
		input := JobSpecInput{
			"command":       123,
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command is required and must be a string")
	})

	t.Run("missing config", func(t *testing.T) {
		input := JobSpecInput{
			"command":       "run",
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config is required")
	})

	t.Run("invalid config type", func(t *testing.T) {
		input := JobSpecInput{
			"command":       "run",
			"config":        123,
			"externalJobID": "123",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config is required and must be a string")
	})

	t.Run("missing externalJobID", func(t *testing.T) {
		input := JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "externalJobID is required")
	})

	t.Run("invalid externalJobID type", func(t *testing.T) {
		input := JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": 123,
			"oracleFactory": pkg.OracleFactory{},
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "externalJobID is required and must be a string")
	})

	t.Run("missing oracleFactory", func(t *testing.T) {
		input := JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": "123",
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "oracleFactory is required")
	})

	t.Run("invalid oracleFactory type", func(t *testing.T) {
		input := JobSpecInput{
			"command":       "run",
			"config":        "param=value",
			"externalJobID": "123",
			"oracleFactory": "not a factory",
		}
		_, err := input.ToStandardCapabilityJob(jobName)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "oracleFactory is required and must be of type OracleFactory")
	})
}
