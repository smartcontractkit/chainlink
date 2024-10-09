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

			[oracle_factory]
			enabled=true
			bootstrap_peers = [
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

			[oracle_factory]
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

			[capabilities]
			target = "enabled"

			[oracle_factory]
			enabled=true
			bootstrap_peers = [
				"12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690"
			]
			network="evm"
			chain_id="31337"
			ocr_contract_address="0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"
			`,
			expectedSpec: &job.StandardCapabilitiesSpec{
				Command: "path/to/binary",
				OracleFactory: job.OracleFactoryConfig{
					Enabled: true,
					BootstrapPeers: []string{
						"12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690",
					},
					OCRContractAddress: "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6",
					ChainID:            "31337",
					Network:            "evm",
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
