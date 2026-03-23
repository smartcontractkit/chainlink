package standardcapabilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
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
			jobSpec, err := ValidatedStandardCapabilitiesSpec(tc.tomlString)

			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tc.expectedSpec != nil {
				assert.Equal(t, tc.expectedSpec, jobSpec.StandardCapabilitiesSpec)
			}
		})
	}
}

func Test_ServicesForSpec_AllowlistEnforcement(t *testing.T) {
	t.Run("allowlisted consensus capability is rejected", func(t *testing.T) {
		d := &Delegate{
			localCfg: &stubLocalCapabilities{allowlisted: map[string]bool{"consensus@1.0.0-alpha": true}},
		}
		spec := job.Job{
			ExternalJobID:            uuid.New(),
			StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{Command: "consensus"},
		}
		_, err := d.ServicesForSpec(context.Background(), spec)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistryBasedLaunchAllowlist")
		assert.Contains(t, err.Error(), "LocalCapabilityManager")
		assert.Contains(t, err.Error(), "consensus@1.0.0-alpha")
	})

	t.Run("allowlisted cron trigger is rejected", func(t *testing.T) {
		d := &Delegate{
			localCfg: &stubLocalCapabilities{allowlisted: map[string]bool{"cron-trigger@1.0.0": true}},
		}
		spec := job.Job{
			ExternalJobID:            uuid.New(),
			StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{Command: "cron"},
		}
		_, err := d.ServicesForSpec(context.Background(), spec)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistryBasedLaunchAllowlist")
		assert.Contains(t, err.Error(), "cron-trigger@1.0.0")
	})

	t.Run("allowlisted http trigger is rejected", func(t *testing.T) {
		d := &Delegate{
			localCfg: &stubLocalCapabilities{allowlisted: map[string]bool{"http-trigger@1.0.0-alpha": true}},
		}
		spec := job.Job{
			ExternalJobID:            uuid.New(),
			StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{Command: "http_trigger"},
		}
		_, err := d.ServicesForSpec(context.Background(), spec)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistryBasedLaunchAllowlist")
	})

	t.Run("non-allowlisted capability passes through", func(t *testing.T) {
		d := &Delegate{
			localCfg: &stubLocalCapabilities{allowlisted: map[string]bool{}},
		}
		spec := job.Job{
			ExternalJobID:            uuid.New(),
			StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{Command: "consensus"},
		}
		// The allowlist check should pass; the call will panic deeper in NewServices due to
		// nil dependencies. A panic (not an allowlist error) proves the check passed.
		assert.Panics(t, func() {
			_, _ = d.ServicesForSpec(context.Background(), spec)
		})
	})

	t.Run("nil localCfg allows all capabilities", func(t *testing.T) {
		d := &Delegate{}
		spec := job.Job{
			ExternalJobID:            uuid.New(),
			StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{Command: "consensus"},
		}
		assert.Panics(t, func() {
			_, _ = d.ServicesForSpec(context.Background(), spec)
		})
	})

	t.Run("unknown command bypasses allowlist check", func(t *testing.T) {
		d := &Delegate{
			localCfg: &stubLocalCapabilities{allowlisted: map[string]bool{"consensus@1.0.0-alpha": true}},
		}
		spec := job.Job{
			ExternalJobID:            uuid.New(),
			StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{Command: "unknown-binary"},
		}
		// Unknown commands have no capability ID mapping, so the allowlist check is skipped.
		assert.Panics(t, func() {
			_, _ = d.ServicesForSpec(context.Background(), spec)
		})
	})
}

// stubLocalCapabilities is a minimal test implementation of config.LocalCapabilities.
type stubLocalCapabilities struct {
	allowlisted map[string]bool
}

func (s *stubLocalCapabilities) RegistryBasedLaunchAllowlist() []string { return nil }
func (s *stubLocalCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	return nil
}
func (s *stubLocalCapabilities) IsAllowlisted(capabilityID string) bool {
	return s.allowlisted[capabilityID]
}
func (s *stubLocalCapabilities) GetCapabilityConfig(string) config.CapabilityNodeConfig {
	return nil
}
