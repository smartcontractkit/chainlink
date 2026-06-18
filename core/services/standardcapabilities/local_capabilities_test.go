package standardcapabilities

import (
	"testing"

	"github.com/stretchr/testify/require"

	capStreams "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type mockTriggerLocalCapabilities struct {
	cfgs map[string]config.CapabilityNodeConfig
}

func (s mockTriggerLocalCapabilities) RegistryBasedLaunchAllowlist() []string {
	return nil
}

func (s mockTriggerLocalCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	return s.cfgs
}

func (s mockTriggerLocalCapabilities) IsAllowlisted(string) bool {
	return false
}

func (s mockTriggerLocalCapabilities) GetCapabilityConfig(capabilityID string) config.CapabilityNodeConfig {
	if s.cfgs == nil {
		return nil
	}
	return s.cfgs[capabilityID]
}

type mockTriggerCapabilityNodeConfig struct{}

func (mockTriggerCapabilityNodeConfig) BinaryPathOverride() string {
	return ""
}

func (mockTriggerCapabilityNodeConfig) Config() map[string]string {
	return nil
}

func TestShouldRegisterMockStreamsTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		localCfg config.LocalCapabilities
		want     bool
	}{
		{name: "nil config", localCfg: nil, want: false},
		{
			name: "different local capability",
			localCfg: mockTriggerLocalCapabilities{
				cfgs: map[string]config.CapabilityNodeConfig{
					"cron@1.0.0": mockTriggerCapabilityNodeConfig{},
				},
			},
			want: false,
		},
		{
			name: "mock trigger opted in",
			localCfg: mockTriggerLocalCapabilities{
				cfgs: map[string]config.CapabilityNodeConfig{
					capStreams.MockTriggerCapabilityID: mockTriggerCapabilityNodeConfig{},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, shouldRegisterMockStreamsTrigger(tt.localCfg))
		})
	}
}
