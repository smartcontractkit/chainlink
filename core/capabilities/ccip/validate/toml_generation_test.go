package validate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCCIPSpecToml_BasicGeneration tests the basic TOML generation functionality
func TestNewCCIPSpecToml_BasicGeneration(t *testing.T) {
	tests := []struct {
		name     string
		args     SpecArgs
		expected []string // Expected strings to be present in TOML
	}{
		{
			name: "basic_configuration",
			args: SpecArgs{
				CapabilityLabelledName: "ccip-capability",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-key-123",
				OCRKeyBundleIDs: map[string]string{
					"evm.1": "bundle-1",
					"evm.2": "bundle-2",
				},
				RelayConfigs: map[string]any{
					"evm.1": map[string]any{
						"chainID": "1",
						"rpcURL":  "https://eth-mainnet.example.com",
					},
				},
				PluginConfig: map[string]any{
					"maxGasPrice": "1000000000",
					"timeout":     "30s",
				},
				P2PV2Bootstrappers: []string{
					"12D3KooWExample1@192.168.1.1:9000",
					"12D3KooWExample2@192.168.1.2:9001",
				},
			},
			expected: []string{
				`type = "ccip"`,
				`capabilityLabelledName = "ccip-capability"`,
				`capabilityVersion = "v1.0.0"`,
				`p2pKeyID = "test-key-123"`,
				`[ocrKeyBundleIDs]`,
				`"evm.1" = "bundle-1"`,
				`"evm.2" = "bundle-2"`,
				`[relayConfigs]`,
				`[pluginConfig]`,
				`maxGasPrice = "1000000000"`,
				`timeout = "30s"`,
				`p2pV2Bootstrappers = ["12D3KooWExample1@192.168.1.1:9000", "12D3KooWExample2@192.168.1.2:9001"]`,
			},
		},
		{
			name: "minimal_configuration",
			args: SpecArgs{
				CapabilityLabelledName: "minimal-ccip",
				CapabilityVersion:      "v0.1.0",
				P2PKeyID:               "minimal-key",
				OCRKeyBundleIDs:        map[string]string{},
				RelayConfigs:           map[string]any{},
				PluginConfig:           map[string]any{},
				P2PV2Bootstrappers:     []string{},
			},
			expected: []string{
				`type = "ccip"`,
				`capabilityLabelledName = "minimal-ccip"`,
				`capabilityVersion = "v0.1.0"`,
				`p2pKeyID = "minimal-key"`,
				`p2pV2Bootstrappers = []`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCCIPSpecToml(tt.args)
			require.NoError(t, err, "NewCCIPSpecToml should not return an error")
			require.NotEmpty(t, result, "Generated TOML should not be empty")

			// Check that all expected strings are present in the generated TOML
			for _, expected := range tt.expected {
				assert.Contains(t, result, expected, "Generated TOML should contain: %s", expected)
			}

			// Verify TOML structure basics
			assert.True(t, strings.Contains(result, "[ocrKeyBundleIDs]") || len(tt.args.OCRKeyBundleIDs) == 0,
				"TOML should contain ocrKeyBundleIDs section when not empty")
			assert.True(t, strings.Contains(result, "[relayConfigs]") || len(tt.args.RelayConfigs) == 0,
				"TOML should contain relayConfigs section when not empty")
			assert.True(t, strings.Contains(result, "[pluginConfig]") || len(tt.args.PluginConfig) == 0,
				"TOML should contain pluginConfig section when not empty")
		})
	}
}

// TestNewCCIPSpecToml_SpecialCharacters tests TOML generation with special characters
func TestNewCCIPSpecToml_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		args SpecArgs
	}{
		{
			name: "special_characters_in_strings",
			args: SpecArgs{
				CapabilityLabelledName: "ccip-test-with-special-chars-@#$%",
				CapabilityVersion:      "v1.0.0-beta+build.123",
				P2PKeyID:               "key-with-underscores_and_dashes-123",
				OCRKeyBundleIDs: map[string]string{
					"evm.chain-1": "bundle-with-special@chars",
				},
				RelayConfigs: map[string]any{
					"evm.mainnet": map[string]any{
						"rpcURL": "https://api.example.com/v1/rpc?key=abc123&format=json",
					},
				},
				PluginConfig: map[string]any{
					"description": "Config with quotes \"and\" backslashes\\",
				},
				P2PV2Bootstrappers: []string{
					"12D3KooWSpecialChars@example.com:8080",
				},
			},
		},
		{
			name: "unicode_characters",
			args: SpecArgs{
				CapabilityLabelledName: "ccip-测试-capability",
				CapabilityVersion:      "v1.0.0-α",
				P2PKeyID:               "key-with-émojis-🔗",
				OCRKeyBundleIDs: map[string]string{
					"evm.测试": "bundle-with-unicode-字符",
				},
				RelayConfigs:       map[string]any{},
				PluginConfig:       map[string]any{},
				P2PV2Bootstrappers: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCCIPSpecToml(tt.args)
			require.NoError(t, err, "NewCCIPSpecToml should handle special characters without error")
			require.NotEmpty(t, result, "Generated TOML should not be empty")

			// Verify that the TOML contains the basic structure
			assert.Contains(t, result, `type = "ccip"`, "TOML should contain type field")
			assert.Contains(t, result, "capabilityLabelledName", "TOML should contain capabilityLabelledName field")
			assert.Contains(t, result, "capabilityVersion", "TOML should contain capabilityVersion field")
			assert.Contains(t, result, "p2pKeyID", "TOML should contain p2pKeyID field")
		})
	}
}

// TestNewCCIPSpecToml_ComplexNestedStructures tests TOML generation with complex nested data
func TestNewCCIPSpecToml_ComplexNestedStructures(t *testing.T) {
	args := SpecArgs{
		CapabilityLabelledName: "complex-ccip",
		CapabilityVersion:      "v2.0.0",
		P2PKeyID:               "complex-key",
		OCRKeyBundleIDs: map[string]string{
			"evm.1":      "bundle-1",
			"evm.137":    "bundle-polygon",
			"evm.42161":  "bundle-arbitrum",
			"evm.10":     "bundle-optimism",
			"evm.43114":  "bundle-avalanche",
		},
		RelayConfigs: map[string]any{
			"evm.1": map[string]any{
				"chainID":           "1",
				"rpcURL":            "https://mainnet.infura.io/v3/key",
				"gasEstimatorMode":  "BlockHistory",
				"blockConfirmations": 12,
				"linkContractAddress": "0x514910771AF9Ca656af840dff83E8264EcF986CA",
			},
			"evm.137": map[string]any{
				"chainID":           "137",
				"rpcURL":            "https://polygon-rpc.com",
				"gasEstimatorMode":  "FixedPrice",
				"blockConfirmations": 200,
				"finalityDepth":     500,
			},
		},
		PluginConfig: map[string]any{
			"maxGasPrice":     "100000000000",
			"gasLimitDefault": 500000,
			"timeout":         "60s",
			"retryConfig": map[string]any{
				"maxRetries":    3,
				"retryInterval": "5s",
				"backoffFactor": 2.0,
			},
			"monitoring": map[string]any{
				"enabled":        true,
				"metricsPort":    8080,
				"healthCheckURL": "/health",
			},
		},
		P2PV2Bootstrappers: []string{
			"12D3KooWBootstrap1@bootstrap1.example.com:9000",
			"12D3KooWBootstrap2@bootstrap2.example.com:9001",
			"12D3KooWBootstrap3@bootstrap3.example.com:9002",
			"12D3KooWBootstrap4@bootstrap4.example.com:9003",
		},
	}

	result, err := NewCCIPSpecToml(args)
	require.NoError(t, err, "NewCCIPSpecToml should handle complex nested structures")
	require.NotEmpty(t, result, "Generated TOML should not be empty")

	// Verify main sections are present
	assert.Contains(t, result, `type = "ccip"`, "TOML should contain type field")
	assert.Contains(t, result, `capabilityLabelledName = "complex-ccip"`, "TOML should contain capabilityLabelledName")
	assert.Contains(t, result, `capabilityVersion = "v2.0.0"`, "TOML should contain capabilityVersion")
	assert.Contains(t, result, `p2pKeyID = "complex-key"`, "TOML should contain p2pKeyID")

	// Verify sections are present
	assert.Contains(t, result, "[ocrKeyBundleIDs]", "TOML should contain ocrKeyBundleIDs section")
	assert.Contains(t, result, "[relayConfigs]", "TOML should contain relayConfigs section")
	assert.Contains(t, result, "[pluginConfig]", "TOML should contain pluginConfig section")

	// Verify some specific nested values
	assert.Contains(t, result, "bundle-polygon", "TOML should contain polygon bundle")
	assert.Contains(t, result, "bootstrap1.example.com", "TOML should contain bootstrap node")
	assert.Contains(t, result, "maxGasPrice", "TOML should contain maxGasPrice config")

	// Verify array formatting
	assert.Contains(t, result, "p2pV2Bootstrappers = [", "TOML should contain bootstrappers array")

	// Check that the TOML is reasonably structured (contains multiple lines)
	lines := strings.Split(result, "\n")
	assert.Greater(t, len(lines), 10, "Generated TOML should have multiple lines")
}

// TestNewCCIPSpecToml_EdgeCases tests edge cases and boundary conditions
func TestNewCCIPSpecToml_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		args        SpecArgs
		expectError bool
		description string
	}{
		{
			name: "empty_strings",
			args: SpecArgs{
				CapabilityLabelledName: "",
				CapabilityVersion:      "",
				P2PKeyID:               "",
				OCRKeyBundleIDs:        map[string]string{},
				RelayConfigs:           map[string]any{},
				PluginConfig:           map[string]any{},
				P2PV2Bootstrappers:     []string{},
			},
			expectError: false,
			description: "Should handle empty strings gracefully",
		},
		{
			name: "nil_maps",
			args: SpecArgs{
				CapabilityLabelledName: "test",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-key",
				OCRKeyBundleIDs:        nil,
				RelayConfigs:           nil,
				PluginConfig:           nil,
				P2PV2Bootstrappers:     nil,
			},
			expectError: false,
			description: "Should handle nil maps and slices gracefully",
		},
		{
			name: "very_long_values",
			args: SpecArgs{
				CapabilityLabelledName: strings.Repeat("a", 1000),
				CapabilityVersion:      strings.Repeat("v", 500) + ".0.0",
				P2PKeyID:               strings.Repeat("k", 2000),
				OCRKeyBundleIDs: map[string]string{
					strings.Repeat("chain", 100): strings.Repeat("bundle", 200),
				},
				RelayConfigs:       map[string]any{},
				PluginConfig:       map[string]any{},
				P2PV2Bootstrappers: []string{},
			},
			expectError: false,
			description: "Should handle very long string values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCCIPSpecToml(tt.args)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotEmpty(t, result, "Generated TOML should not be empty")
				// Verify basic TOML structure is maintained
				assert.Contains(t, result, `type = "ccip"`, "TOML should always contain type field")
			}
		})
	}
}

// TestNewCCIPSpecToml_DataTypes tests various data types in TOML generation
func TestNewCCIPSpecToml_DataTypes(t *testing.T) {
	args := SpecArgs{
		CapabilityLabelledName: "datatype-test",
		CapabilityVersion:      "v1.0.0",
		P2PKeyID:               "datatype-key",
		OCRKeyBundleIDs: map[string]string{
			"string_value": "test-string",
			"int_value":    "123",
			"float_value":  "45.67",
			"bool_value":   "true",
		},
		RelayConfigs: map[string]any{
			"evm.1": map[string]any{
				"stringField":  "value",
				"intField":     42,
				"floatField":   3.14159,
				"boolField":    false,
				"arrayField":   []string{"item1", "item2", "item3"},
				"nestedObject": map[string]any{
					"nestedString": "nested_value",
					"nestedInt":    999,
				},
			},
		},
		PluginConfig: map[string]any{
			"mixedTypes": map[string]any{
				"timeout":       "30s",
				"maxRetries":    5,
				"enableLogging": true,
				"threshold":     0.95,
				"endpoints":     []string{"http://api1.com", "http://api2.com"},
			},
		},
		P2PV2Bootstrappers: []string{
			"12D3KooWDataType1@192.168.1.100:9000",
			"12D3KooWDataType2@192.168.1.101:9001",
		},
	}

	result, err := NewCCIPSpecToml(args)
	require.NoError(t, err, "NewCCIPSpecToml should handle various data types")
	require.NotEmpty(t, result, "Generated TOML should not be empty")

	// Verify that different data types are properly represented
	assert.Contains(t, result, "123", "TOML should contain integer values")
	assert.Contains(t, result, "45.67", "TOML should contain float values")
	assert.Contains(t, result, "true", "TOML should contain boolean values")
	assert.Contains(t, result, "false", "TOML should contain boolean values")
	assert.Contains(t, result, "test-string", "TOML should contain string values")

	// Verify array representation
	assert.Contains(t, result, "[", "TOML should contain array brackets")
	assert.Contains(t, result, "item1", "TOML should contain array items")

	// Verify nested structure representation
	assert.Contains(t, result, "nested_value", "TOML should contain nested values")
	assert.Contains(t, result, "999", "TOML should contain nested integer values")

	// Verify main structure
	assert.Contains(t, result, `type = "ccip"`, "TOML should contain type field")
	assert.Contains(t, result, "[ocrKeyBundleIDs]", "TOML should contain ocrKeyBundleIDs section")
	assert.Contains(t, result, "[relayConfigs]", "TOML should contain relayConfigs section")
	assert.Contains(t, result, "[pluginConfig]", "TOML should contain pluginConfig section")
}