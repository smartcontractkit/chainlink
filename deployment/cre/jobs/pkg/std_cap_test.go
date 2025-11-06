package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

func TestStdCap_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		job   pkg.StandardCapabilityJob
		error string
	}{
		{
			name:  "must contain name",
			job:   pkg.StandardCapabilityJob{},
			error: pkg.ErrorEmptyJobName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.ErrorContains(t, tt.job.Validate(), tt.error)
		})
	}
}

const (
	expectedCRONSpec = `type = "standardcapabilities"
schemaVersion = 1
name = "cron-capabilities"
externalJobID = "14d3a547-5e4d-5f22-bfd7-9940cc6cefe2"
forwardingAllowed = false
command = "cron"
config = """"""
`
	expectedComputeSpec = `type = "standardcapabilities"
schemaVersion = 1
name = "compute-capabilities"
externalJobID = "fe41c282-0393-5559-9e6c-4ce592b7f9ac"
forwardingAllowed = false
command = "__builtin_custom-compute-action"
config = """NumWorkers = 3
[rateLimiter]
globalRPS = 20.0
globalBurst = 30
perSenderRPS = 1.0
perSenderBurst = 5"""
`
	expectedHTTPTriggerSpec = `type = "standardcapabilities"
schemaVersion = 1
name = "http-trigger-capabilities"
externalJobID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
forwardingAllowed = false
command = "http_trigger"
config = """{}"""
`
	expectedHTTPActionSpec = `type = "standardcapabilities"
schemaVersion = 1
name = "http-action-capabilities"
externalJobID = "b2c3d4e5-f6g7-8901-bcde-f23456789012"
forwardingAllowed = false
command = "http_action"
config = """{"proxyMode": "direct"}"""
`
)

func TestStdCap_Resolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expectedTOML  string
		command       string
		config        string
		externalJobID string
	}{
		{
			name:          "cron-capabilities",
			expectedTOML:  expectedCRONSpec,
			command:       "cron",
			config:        "",
			externalJobID: "14d3a547-5e4d-5f22-bfd7-9940cc6cefe2",
		},
		{
			name:         "compute-capabilities",
			expectedTOML: expectedComputeSpec,
			command:      "__builtin_custom-compute-action",
			config: `NumWorkers = 3
[rateLimiter]
globalRPS = 20.0
globalBurst = 30
perSenderRPS = 1.0
perSenderBurst = 5`,
			externalJobID: "fe41c282-0393-5559-9e6c-4ce592b7f9ac",
		},
		{
			name:          "http-trigger-capabilities",
			expectedTOML:  expectedHTTPTriggerSpec,
			command:       "http_trigger",
			config:        "{}",
			externalJobID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		},
		{
			name:          "http-action-capabilities",
			expectedTOML:  expectedHTTPActionSpec,
			command:       "http_action",
			config:        `{"proxyMode": "direct"}`,
			externalJobID: "b2c3d4e5-f6g7-8901-bcde-f23456789012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := pkg.StandardCapabilityJob{
				JobName:       tt.name,
				Command:       tt.command,
				Config:        tt.config,
				ExternalJobID: tt.externalJobID,
			}

			spec, err := s.Resolve()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedTOML, spec)
		})
	}
}
