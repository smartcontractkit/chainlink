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
			type = "standardcapabilities"
			schemaVersion = 1
			name = "consensus-capabilities"
			externalJobID = "aea7103f-6e87-5c01-b644-a0b4aeaed3eb"
			forwardingAllowed = false
			command = "path/to/binary"
			config = """"""
			
			[oracle_factory]
			enabled = true
			bootstrap_peers = ["12D3KooWBAzThfs9pD4WcsFKCi68EUz2fZgZskDBT6JcJRndPss5@cl-keystone-two-bt-0:5001"]
			ocr_contract_address = "0x2C84cff4cd5fA5a0c17dbc710fcCb8FC6A03dEEd"
			ocr_key_bundle_id = "5fbb7d5dc1e592142a979b7014552e07a78cb89b1a8626c6412f12f2adfcb240"
			chain_id = "11155111"
			transmitter_id = "0x60042fBB756f736744C334c463BeBE1A72Add04F"
			[oracle_factory.onchainSigningStrategy]
			strategyName = "multi-chain"
			[oracle_factory.onchainSigningStrategy.config]
			aptos = "7c2df2e806306383f9aa2bc7a3198cf0e1c626f873799992b2841240c6931733"
			evm = "5fbb7d5dc1e592142a979b7014552e07a78cb89b1a8626c6412f12f2adfcb240"
			`,
			expectedSpec: &job.StandardCapabilitiesSpec{
				Command: "path/to/binary",
				OracleFactory: job.OracleFactoryConfig{
					Enabled: true,
					BootstrapPeers: []string{
						"12D3KooWBAzThfs9pD4WcsFKCi68EUz2fZgZskDBT6JcJRndPss5@cl-keystone-two-bt-0:5001",
					},
					OCRContractAddress: "0x2C84cff4cd5fA5a0c17dbc710fcCb8FC6A03dEEd",
					OCRKeyBundleID:     "5fbb7d5dc1e592142a979b7014552e07a78cb89b1a8626c6412f12f2adfcb240",
					ChainID:            "11155111",
					TransmitterID:      "0x60042fBB756f736744C334c463BeBE1A72Add04F",
					OnchainSigning: job.OnchainSigningStrategy{
						StrategyName: "multi-chain",
						Config: map[string]string{
							"aptos": "7c2df2e806306383f9aa2bc7a3198cf0e1c626f873799992b2841240c6931733",
							"evm":   "5fbb7d5dc1e592142a979b7014552e07a78cb89b1a8626c6412f12f2adfcb240",
						},
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
