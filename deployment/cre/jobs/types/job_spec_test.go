package job_types_test

import (
	"encoding/json"
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

func TestUnmarshalTo_JSONNumberIntFields(t *testing.T) {
	t.Parallel()

	type DON struct {
		Name     string   `yaml:"name"`
		F        int      `yaml:"f"`
		Handlers []string `yaml:"handlers"`
	}

	type GatewayInput struct {
		GatewayRequestTimeoutSec int    `yaml:"gatewayRequestTimeoutSec"`
		DONs                     []DON  `yaml:"dons"`
		SomeString               string `yaml:"someString"`
	}

	// Simulate the exact values that arrive after the YAML->JSON->UseNumber() pipeline.
	// In the real pipeline, YamlNodeToAny converts YAML integers to json.Number,
	// then json.Marshal -> env var -> json.Decoder.UseNumber() preserves them as json.Number.
	input := job_types.JobSpecInput{
		"gatewayRequestTimeoutSec": json.Number("70"),
		"someString":               "hello",
		"dons": []any{
			map[string]any{
				"name":     "some_don",
				"f":        json.Number("1"),
				"handlers": []any{"http-capabilities", "web-api-capabilities"},
			},
		},
	}

	var target GatewayInput
	err := input.UnmarshalTo(&target)
	require.NoError(t, err, "UnmarshalTo should handle json.Number values in int fields")

	assert.Equal(t, 70, target.GatewayRequestTimeoutSec)
	assert.Equal(t, "hello", target.SomeString)
	require.Len(t, target.DONs, 1)
	assert.Equal(t, "some_don", target.DONs[0].Name)
	assert.Equal(t, 1, target.DONs[0].F)
	assert.Equal(t, []string{"http-capabilities", "web-api-capabilities"}, target.DONs[0].Handlers)
}

// TestUnmarshalTo_NativeIntFields verifies the happy path where map values
// are already native Go ints (e.g. when constructed in Go code rather than
// through the JSON pipeline). This should always pass.
func TestUnmarshalTo_NativeIntFields(t *testing.T) {
	t.Parallel()

	type Target struct {
		Count   int    `yaml:"count"`
		Message string `yaml:"message"`
	}

	input := job_types.JobSpecInput{
		"count":   42,
		"message": "test",
	}

	var target Target
	err := input.UnmarshalTo(&target)
	require.NoError(t, err)
	assert.Equal(t, 42, target.Count)
	assert.Equal(t, "test", target.Message)
}
