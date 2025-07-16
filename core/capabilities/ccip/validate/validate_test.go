package validate_test

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestNewCCIPSpecToml(t *testing.T) {
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		want     string
		wantErr  bool
	}{
		{
			name: "valid spec args",
			specArgs: validate.SpecArgs{
				P2PV2Bootstrappers:     []string{"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:9000"},
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				OCRKeyBundleIDs:        map[string]string{"evm": "test-key-bundle-id"},
				P2PKeyID:               "test-p2p-key-id",
				RelayConfigs:           map[string]any{"evm": map[string]any{"chainReader": map[string]any{}}},
				PluginConfig:           map[string]any{"tokenPrices": "test-pipeline"},
			},
			wantErr: false,
		},
		{
			name: "empty capability version",
			specArgs: validate.SpecArgs{
				P2PV2Bootstrappers:     []string{"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:9000"},
				CapabilityVersion:      "",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false, // NewCCIPSpecToml doesn't validate, only generates TOML
		},
		{
			name: "minimal valid spec",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate.NewCCIPSpecToml(tt.specArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify the generated TOML contains expected fields
				require.Contains(t, got, `type = "ccip"`)
				require.Contains(t, got, `schemaVersion = 1`)
				require.Contains(t, got, fmt.Sprintf(`capabilityVersion = "%s"`, tt.specArgs.CapabilityVersion))
				require.Contains(t, got, fmt.Sprintf(`capabilityLabelledName = "%s"`, tt.specArgs.CapabilityLabelledName))
				require.Contains(t, got, fmt.Sprintf(`p2pKeyID = "%s"`, tt.specArgs.P2PKeyID))
				// Verify external job ID is generated
				require.Contains(t, got, "externalJobID")
				require.Contains(t, got, "name")
			}
		})
	}
}

func TestValidatedCCIPSpec(t *testing.T) {
	type args struct {
		tomlString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errorMsg string
	}{
		{
			name: "valid CCIP spec",
			args: args{
				tomlString: `
type = "ccip"
schemaVersion = 1
name = "test-ccip-job"
externalJobID = "550e8400-e29b-41d4-a716-446655440000"
capabilityVersion = "v1.0.0"
capabilityLabelledName = "ccip"
p2pKeyID = "test-p2p-key-id"
p2pV2Bootstrappers = ["12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:9000"]
[ocrKeyBundleIDs]
evm = "test-key-bundle-id"
[relayConfigs]
[pluginConfig]
`,
			},
			wantErr: false,
		},
		{
			name: "invalid TOML syntax",
			args: args{
				tomlString: `invalid toml [[[`,
			},
			wantErr: true,
			errorMsg: "toml error on load",
		},
		{
			name: "wrong job type",
			args: args{
				tomlString: `
type = "webhook"
schemaVersion = 1
capabilityVersion = "v1.0.0"
capabilityLabelledName = "ccip"
p2pKeyID = "test-p2p-key-id"
`,
			},
			wantErr: true,
			errorMsg: "the only supported type is currently 'ccip'",
		},
		{
			name: "missing capabilityLabelledName",
			args: args{
				tomlString: `
type = "ccip"
schemaVersion = 1
capabilityVersion = "v1.0.0"
p2pKeyID = "test-p2p-key-id"
`,
			},
			wantErr: true,
			errorMsg: "capabilityLabelledName must be set",
		},
		{
			name: "missing capabilityVersion",
			args: args{
				tomlString: `
type = "ccip"
schemaVersion = 1
capabilityLabelledName = "ccip"
p2pKeyID = "test-p2p-key-id"
`,
			},
			wantErr: true,
			errorMsg: "capabilityVersion must be set",
		},
		{
			name: "missing p2pKeyID",
			args: args{
				tomlString: `
type = "ccip"
schemaVersion = 1
capabilityVersion = "v1.0.0"
capabilityLabelledName = "ccip"
`,
			},
			wantErr: true,
			errorMsg: "p2pKeyID must be set",
		},
		{
			name: "invalid bootstrapper format",
			args: args{
				tomlString: `
type = "ccip"
schemaVersion = 1
capabilityVersion = "v1.0.0"
capabilityLabelledName = "ccip"
p2pKeyID = "test-p2p-key-id"
p2pV2Bootstrappers = ["invalid-bootstrapper-format"]
`,
			},
			wantErr: true,
			errorMsg: "p2p v2 bootstrapper locator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJb, err := validate.ValidatedCCIPSpec(tt.args.tomlString)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				// Verify the job was properly constructed
				require.Equal(t, job.CCIP, gotJb.Type)
				require.NotNil(t, gotJb.CCIPSpec)
				require.Equal(t, "v1.0.0", gotJb.CCIPSpec.CapabilityVersion)
				require.Equal(t, "ccip", gotJb.CCIPSpec.CapabilityLabelledName)
				require.Equal(t, "test-p2p-key-id", gotJb.CCIPSpec.P2PKeyID)
			}
		})
	}
}
