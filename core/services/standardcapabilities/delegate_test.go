package standardcapabilities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
)

func Test_ValidatedStandardCapabilitiesSpec(t *testing.T) {
	type testCase struct {
		name          string
		tomlString    string
		expectedError string
		expectedSpec  *job.StandardCapabilitiesSpec
	}

	testCases := []testCase{
		{
			name:          "invalid TOML string",
			tomlString:    `[[]`,
			expectedError: "toml error on load standard capabilities",
		},
		{
			name: "incorrect job type",
			tomlString: `
			type="nonstandardcapabilities"
			`,
			expectedError: "standard capabilities unsupported job type",
		},
		{
			name: "command unset",
			tomlString: `
			type="standardcapabilities"
			`,
			expectedError: "standard capabilities command must be set",
		},
		{
			name: "invalid oracle config: malformed peer",
			tomlString: `
			type="standardcapabilities"
			command="path/to/binary"

			[oracleFactory]
			enabled=true
			bootstrapPeers = [
				"invalid_p2p_id@invalid_ip:1111"
			]
			`,
			expectedError: "failed to parse bootstrap peers",
		},
		{
			name: "invalid oracle config: missing bootstrap peers",
			tomlString: `
			type="standardcapabilities"
			command="path/to/binary"

			[oracleFactory]
			enabled=true
			`,
			expectedError: "no bootstrap peers found",
		},
		{
			name: "valid spec",
			tomlString: `
			type="standardcapabilities"
			command="path/to/binary"
			`,
		},
		{
			name: "valid spec with oracle config",
			tomlString: `
			type="standardcapabilities"
			command="path/to/binary"

			[oracleFactory]
			enabled=true
			bootstrapPeers = [
				"12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690"
			]
			`,
			expectedSpec: &job.StandardCapabilitiesSpec{
				Command: "path/to/binary",
				OracleFactory: job.JSONConfig{
					"enabled": true,
					"bootstrapPeers": []interface{}{
						"12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jobSpec, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(tc.tomlString)

			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tc.expectedSpec != nil {
				assert.EqualValues(t, tc.expectedSpec, jobSpec.StandardCapabilitiesSpec)
			}
		})
	}
}
