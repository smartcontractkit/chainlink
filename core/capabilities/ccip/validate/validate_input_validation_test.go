package validate_test

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/stretchr/testify/require"
)

func TestSpecArgs_OCRKeyBundleIDs_Validation(t *testing.T) {
	tests := []struct {
		name            string
		ocrKeyBundleIDs map[string]string
		wantErr         bool
		errorContains   string
	}{
		{
			name: "valid OCR key bundle IDs",
			ocrKeyBundleIDs: map[string]string{
				"evm":    "test-key-bundle-id-evm",
				"solana": "test-key-bundle-id-solana",
			},
			wantErr: false,
		},
		{
			name:            "empty OCR key bundle IDs map",
			ocrKeyBundleIDs: map[string]string{},
			wantErr:         false, // Empty map should be allowed
		},
		{
			name:            "nil OCR key bundle IDs",
			ocrKeyBundleIDs: nil,
			wantErr:         false, // Nil map should be allowed
		},
		{
			name: "empty key in OCR key bundle IDs",
			ocrKeyBundleIDs: map[string]string{
				"": "test-key-bundle-id",
			},
			wantErr:       true, // Empty key causes TOML syntax error
			errorContains: "toml error",
		},
		{
			name: "empty value in OCR key bundle IDs",
			ocrKeyBundleIDs: map[string]string{
				"evm": "",
			},
			wantErr: false, // Empty value should be allowed
		},
		{
			name: "special characters in OCR key bundle IDs",
			ocrKeyBundleIDs: map[string]string{
				"evm-mainnet":    "key-bundle-123-abc",
				"polygon_mumbai": "key_bundle_456_def",
				"arbitrum.one":   "key.bundle.789.ghi",
			},
			wantErr: false,
		},
		{
			name: "very long OCR key bundle ID",
			ocrKeyBundleIDs: map[string]string{
				"evm": strings.Repeat("a", 1000),
			},
			wantErr: false, // Should handle long strings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specArgs := validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
				OCRKeyBundleIDs:        tt.ocrKeyBundleIDs,
			}

			tomlString, err := validate.NewCCIPSpecToml(specArgs)
			if tt.wantErr {
				// Error might occur in TOML generation or validation
				if err != nil {
					if tt.errorContains != "" {
						require.Contains(t, err.Error(), tt.errorContains)
					}
					return
				}
				// If TOML generation succeeded, validation should fail
				_, err = validate.ValidatedCCIPSpec(tomlString)
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tomlString)

				// Validate that the generated TOML can be parsed back
				_, err = validate.ValidatedCCIPSpec(tomlString)
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecArgs_RelayConfigs_Validation(t *testing.T) {
	tests := []struct {
		name          string
		relayConfigs  map[string]any
		wantErr       bool
		errorContains string
	}{
		{
			name: "valid simple relay configs",
			relayConfigs: map[string]any{
				"evm": map[string]any{
					"chainReader": map[string]any{},
				},
			},
			wantErr: false,
		},
		{
			name: "complex nested relay configs",
			relayConfigs: map[string]any{
				"evm": map[string]any{
					"chainReader": map[string]any{
						"contracts": map[string]any{
							"ccipHome": map[string]any{
								"address": "0x1234567890123456789012345678901234567890",
								"abi":     "[{\"type\":\"function\"}]",
							},
						},
					},
					"chainWriter": map[string]any{
						"gasLimit": 500000,
						"gasPrice": "20000000000",
					},
				},
				"solana": map[string]any{
					"commitment": "confirmed",
					"timeout":    "30s",
				},
			},
			wantErr: false,
		},
		{
			name:         "empty relay configs",
			relayConfigs: map[string]any{},
			wantErr:      false,
		},
		{
			name:         "nil relay configs",
			relayConfigs: nil,
			wantErr:      false,
		},
		{
			name: "relay configs with various data types",
			relayConfigs: map[string]any{
				"evm": map[string]any{
					"stringValue": "test",
					"intValue":    123,
					"floatValue":  123.456,
					"boolValue":   true,
					"arrayValue":  []string{"a", "b", "c"},
					"mapValue":    map[string]string{"key": "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "deeply nested relay configs",
			relayConfigs: map[string]any{
				"evm": map[string]any{
					"level1": map[string]any{
						"level2": map[string]any{
							"level3": map[string]any{
								"level4": map[string]any{
									"level5": "deep_value",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specArgs := validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
				RelayConfigs:           tt.relayConfigs,
			}

			tomlString, err := validate.NewCCIPSpecToml(specArgs)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tomlString)

				// Validate that the generated TOML can be parsed back
				_, err = validate.ValidatedCCIPSpec(tomlString)
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecArgs_PluginConfig_Validation(t *testing.T) {
	tests := []struct {
		name          string
		pluginConfig  map[string]any
		wantErr       bool
		errorContains string
	}{
		{
			name: "valid simple plugin config",
			pluginConfig: map[string]any{
				"tokenPrices": "test-pipeline",
			},
			wantErr: false,
		},
		{
			name: "complex plugin config",
			pluginConfig: map[string]any{
				"tokenPrices": "test-pipeline",
				"gasEstimator": map[string]any{
					"mode":     "optimistic",
					"gasLimit": 500000,
				},
				"rmnConfig": map[string]any{
					"enabled":   true,
					"threshold": 2,
					"nodes": []map[string]any{
						{"url": "http://rmn1.example.com", "weight": 1},
						{"url": "http://rmn2.example.com", "weight": 1},
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "empty plugin config",
			pluginConfig: map[string]any{},
			wantErr:      false,
		},
		{
			name:         "nil plugin config",
			pluginConfig: nil,
			wantErr:      false,
		},
		{
			name: "plugin config with special characters",
			pluginConfig: map[string]any{
				"token-prices":    "test_pipeline",
				"gas.estimator":   "optimistic",
				"rmn@config":      "enabled",
				"config[special]": "value",
			},
			wantErr: false,
		},
		{
			name: "plugin config with large values",
			pluginConfig: map[string]any{
				"largeString": strings.Repeat("x", 10000),
				"largeNumber": 999999999999999,
				"largeArray":  make([]string, 1000),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specArgs := validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
				PluginConfig:           tt.pluginConfig,
			}

			tomlString, err := validate.NewCCIPSpecToml(specArgs)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tomlString)

				// Validate that the generated TOML can be parsed back
				_, err = validate.ValidatedCCIPSpec(tomlString)
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecArgs_P2PV2Bootstrappers_InvalidFormats(t *testing.T) {
	tests := []struct {
		name          string
		bootstrappers []string
		wantErr       bool
		errorContains string
	}{
		{
			name: "valid bootstrapper format",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:9000",
			},
			wantErr: false,
		},
		{
			name: "multiple valid bootstrappers",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:9000",
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@10.0.0.1:9001",
			},
			wantErr: false,
		},
		{
			name:          "empty bootstrappers list",
			bootstrappers: []string{},
			wantErr:       false, // Empty list should be allowed
		},
		{
			name:          "nil bootstrappers",
			bootstrappers: nil,
			wantErr:       false, // Nil should be allowed
		},
		{
			name: "missing @ separator",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N192.168.1.1:9000",
			},
			wantErr:       true,
			errorContains: "p2p v2 bootstrapper locator",
		},
		{
			name: "missing port",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1",
			},
			wantErr:       true,
			errorContains: "p2p v2 bootstrapper locator",
		},
		{
			name: "port 0 is allowed",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:0",
			},
			wantErr: false, // Port 0 is actually allowed
		},
		{
			name: "port 65536 is allowed",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:65536",
			},
			wantErr: false,
		},
		{
			name: "invalid IP format is allowed by parser",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@999.999.999.999:9000",
			},
			wantErr: false, // Parser may allow invalid IP formats
		},
		{
			name: "invalid peer ID format",
			bootstrappers: []string{
				"invalid-peer-id@192.168.1.1:9000",
			},
			wantErr:       true,
			errorContains: "p2p v2 bootstrapper locator",
		},
		{
			name: "empty string in bootstrappers",
			bootstrappers: []string{
				"",
			},
			wantErr:       true,
			errorContains: "p2p v2 bootstrapper locator",
		},
		{
			name: "IPv6 address format",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@[2001:db8::1]:9000",
			},
			wantErr: false, // IPv6 should be supported
		},
		{
			name: "hostname instead of IP",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@bootstrap.example.com:9000",
			},
			wantErr: false, // Hostname should be supported
		},
		{
			name: "mixed valid and invalid bootstrappers",
			bootstrappers: []string{
				"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N@192.168.1.1:9000",
				"invalid-format",
			},
			wantErr:       true,
			errorContains: "p2p v2 bootstrapper locator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specArgs := validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
				P2PV2Bootstrappers:     tt.bootstrappers,
			}

			tomlString, err := validate.NewCCIPSpecToml(specArgs)
			require.NoError(t, err) // TOML generation should always succeed

			// The validation happens in ValidatedCCIPSpec
			_, err = validate.ValidatedCCIPSpec(tomlString)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecArgs_RequiredFields_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		specArgs      validate.SpecArgs
		wantErr       bool
		errorContains string
	}{
		{
			name: "all required fields present",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "capability version with special characters",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      "v1.0.0-beta.1+build.123",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "capability labelled name with special characters",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip-v2",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "p2p key ID with special characters",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      "v1.0.0",
				CapabilityLabelledName: "ccip",
				P2PKeyID:               "test-p2p-key-id_123",
			},
			wantErr: false,
		},
		{
			name: "very long field values",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      strings.Repeat("v", 100) + "1.0.0",
				CapabilityLabelledName: strings.Repeat("ccip", 50),
				P2PKeyID:               strings.Repeat("key", 100),
			},
			wantErr: false,
		},
		{
			name: "unicode characters in fields",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      "v1.0.0-测试",
				CapabilityLabelledName: "ccip-测试",
				P2PKeyID:               "test-p2p-key-id-测试",
			},
			wantErr: false,
		},
		{
			name: "whitespace in field values",
			specArgs: validate.SpecArgs{
				CapabilityVersion:      " v1.0.0 ",
				CapabilityLabelledName: " ccip ",
				P2PKeyID:               " test-p2p-key-id ",
			},
			wantErr: false, // Whitespace should be preserved
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tomlString, err := validate.NewCCIPSpecToml(tt.specArgs)
			require.NoError(t, err) // TOML generation should always succeed

			// The validation happens in ValidatedCCIPSpec
			_, err = validate.ValidatedCCIPSpec(tomlString)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
