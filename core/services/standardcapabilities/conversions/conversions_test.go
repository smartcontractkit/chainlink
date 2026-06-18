package conversions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetCapabilityIDFromCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		config   string
		expected string
	}{
		{
			name:     "consensus command",
			command:  "consensus",
			config:   "",
			expected: "consensus@1.0.0-alpha",
		},
		{
			name:     "evm command with valid config - mainnet",
			command:  "/usr/local/bin/evm",
			config:   `{"chainId": 1}`,
			expected: "evm:ChainSelector:5009297550715157269@1.0.0",
		},
		{
			name:     "evm command with valid config - sepolia",
			command:  "/usr/local/bin/evm",
			config:   `{"chainId": 11155111}`,
			expected: "evm:ChainSelector:16015286601757825753@1.0.0",
		},
		{
			name:     "evm command with valid config - arbitrum",
			command:  "/usr/local/bin/evm",
			config:   `{"chainId": 42161}`,
			expected: "evm:ChainSelector:4949039107694359620@1.0.0",
		},
		{
			name:     "evm command with additional config fields",
			command:  "/usr/local/bin/evm",
			config:   `{"chainId": 1, "network": "mainnet", "otherField": "value"}`,
			expected: "evm:ChainSelector:5009297550715157269@1.0.0",
		},
		{
			name:     "aptos command with valid config - localnet",
			command:  "/usr/local/bin/aptos",
			config:   `{"chainId":"4","network":"aptos"}`,
			expected: "aptos:ChainSelector:4457093679053095497@1.0.0",
		},
		{
			name:     "aptos command with invalid chainId",
			command:  "/usr/local/bin/aptos",
			config:   `{"chainId":"not-a-number","network":"aptos"}`,
			expected: "",
		},
		{
			name:     "aptos command with unknown chainId",
			command:  "/usr/local/bin/aptos",
			config:   `{"chainId":"999999","network":"aptos"}`,
			expected: "",
		},
		{
			name:     "evm command with invalid JSON",
			command:  "/usr/local/bin/evm",
			config:   `{invalid json}`,
			expected: "",
		},
		{
			name:     "evm command with missing chainId",
			command:  "/usr/local/bin/evm",
			config:   `{"network": "mainnet"}`,
			expected: "", // chainId defaults to 0, which is invalid
		},
		{
			name:     "evm command with zero chain ID",
			command:  "/usr/local/bin/evm",
			config:   `{"chainId": 0}`,
			expected: "",
		},
		{
			name:     "evm command with empty config",
			command:  "/usr/local/bin/evm",
			config:   "",
			expected: "",
		},
		{
			name:     "cron command",
			command:  "cron",
			config:   "",
			expected: "cron-trigger@1.0.0",
		},
		{
			name:     "http_trigger command",
			command:  "/opt/bin/http_trigger",
			config:   "",
			expected: "http-trigger@1.0.0-alpha",
		},
		{
			name:     "http_action command",
			command:  "http_action",
			config:   "",
			expected: "http-actions@1.0.0-alpha",
		},
		{
			name:     "unknown command",
			command:  "/usr/local/bin/unknown",
			config:   "",
			expected: "",
		},
		{
			name:     "empty command",
			command:  "",
			config:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCapabilityIDFromCommand(tt.command, tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_GetCommandFromCapabilityID(t *testing.T) {
	tests := []struct {
		name         string
		capabilityID string
		expected     string
	}{
		{
			name:         "consensus capability - exact",
			capabilityID: "consensus@1.0.0-alpha",
			expected:     "consensus",
		},
		{
			name:         "consensus capability - different version",
			capabilityID: "consensus@2.0.0",
			expected:     "consensus",
		},
		{
			name:         "cron trigger capability - exact",
			capabilityID: "cron-trigger@1.0.0",
			expected:     "cron",
		},
		{
			name:         "cron trigger capability - different version",
			capabilityID: "cron-trigger@2.0.0",
			expected:     "cron",
		},
		{
			name:         "http trigger capability - exact",
			capabilityID: "http-trigger@1.0.0-alpha",
			expected:     "http_trigger",
		},
		{
			name:         "http trigger capability - different version",
			capabilityID: "http-trigger@3.0.0",
			expected:     "http_trigger",
		},
		{
			name:         "http action capability - exact",
			capabilityID: "http-actions@1.0.0-alpha",
			expected:     "http_action",
		},
		{
			name:         "http action capability - different version",
			capabilityID: "http-actions@2.0.0",
			expected:     "http_action",
		},
		{
			name:         "evm mainnet capability",
			capabilityID: "evm:ChainSelector:5009297550715157269@1.0.0",
			expected:     "evm",
		},
		{
			name:         "evm sepolia capability",
			capabilityID: "evm:ChainSelector:16015286601757825753@1.0.0",
			expected:     "evm",
		},
		{
			name:         "evm arbitrum capability",
			capabilityID: "evm:ChainSelector:4949039107694359620@1.0.0",
			expected:     "evm",
		},
		{
			name:         "evm capability - different version",
			capabilityID: "evm:ChainSelector:5009297550715157269@2.0.0",
			expected:     "evm",
		},
		{
			name:         "aptos localnet capability",
			capabilityID: "aptos:ChainSelector:4457093679053095497@1.0.0",
			expected:     "aptos",
		},
		{
			name:         "aptos capability - different version",
			capabilityID: "aptos:ChainSelector:4457093679053095497@2.0.0",
			expected:     "aptos",
		},
		{
			name:         "unknown capability",
			capabilityID: "unknown@1.0.0",
			expected:     "",
		},
		{
			name:         "empty capability ID",
			capabilityID: "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCommandFromCapabilityID(tt.capabilityID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_roundTrip(t *testing.T) {
	commands := []string{"consensus", "cron", "http_trigger", "http_action"}
	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			capID := GetCapabilityIDFromCommand(cmd, "")
			assert.NotEmpty(t, capID)
			got := GetCommandFromCapabilityID(capID)
			assert.Equal(t, cmd, got)
		})
	}

	// EVM round-trip: command base name is preserved
	evmCapID := GetCapabilityIDFromCommand("/usr/local/bin/evm", `{"chainId": 1}`)
	assert.Equal(t, "evm", GetCommandFromCapabilityID(evmCapID))

	// Aptos round-trip: command base name is preserved
	aptosCapID := GetCapabilityIDFromCommand("/usr/local/bin/aptos", `{"chainId":"4","network":"aptos"}`)
	assert.Equal(t, "aptos", GetCommandFromCapabilityID(aptosCapID))
}
