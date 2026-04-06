package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	capStreams "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type testLocalCapabilities struct {
	cfgs map[string]config.CapabilityNodeConfig
}

func (t testLocalCapabilities) RegistryBasedLaunchAllowlist() []string {
	return nil
}

func (t testLocalCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	return t.cfgs
}

func (t testLocalCapabilities) IsAllowlisted(string) bool {
	return false
}

func (t testLocalCapabilities) GetCapabilityConfig(capabilityID string) config.CapabilityNodeConfig {
	if t.cfgs == nil {
		return nil
	}

	return t.cfgs[capabilityID]
}

type testCapabilityNodeConfig struct{}

func (testCapabilityNodeConfig) BinaryPathOverride() string {
	return ""
}

func (testCapabilityNodeConfig) Config() map[string]string {
	return nil
}

func TestNewLocalTestMetadataRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		localCfg  config.LocalCapabilities
		expectedF uint8
	}{
		{
			name:      "default workflow DON fault tolerance",
			localCfg:  nil,
			expectedF: 0,
		},
		{
			name: "mock trigger opt-in uses workflow DON fault tolerance one",
			localCfg: testLocalCapabilities{
				cfgs: map[string]config.CapabilityNodeConfig{
					capStreams.MockTriggerCapabilityID: testCapabilityNodeConfig{},
				},
			},
			expectedF: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := newLocalTestMetadataRegistry(tt.localCfg)
			require.Equal(t, tt.expectedF, registry.WorkflowDONF)
		})
	}
}
