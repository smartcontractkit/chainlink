package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"

	capStreams "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type stubLocalCapabilities struct {
	cfgs map[string]config.CapabilityNodeConfig
}

func (s stubLocalCapabilities) RegistryBasedLaunchAllowlist() []string {
	return nil
}

func (s stubLocalCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	return s.cfgs
}

func (s stubLocalCapabilities) IsAllowlisted(string) bool {
	return false
}

func (s stubLocalCapabilities) GetCapabilityConfig(capabilityID string) config.CapabilityNodeConfig {
	if s.cfgs == nil {
		return nil
	}
	return s.cfgs[capabilityID]
}

type stubCapabilityNodeConfig struct{}

func (stubCapabilityNodeConfig) BinaryPathOverride() string {
	return ""
}

func (stubCapabilityNodeConfig) Config() map[string]string {
	return nil
}

func TestShouldRegisterMockStreamsTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		local config.LocalCapabilities
		want  bool
	}{
		{name: "nil config", local: nil, want: false},
		{
			name: "different local capability",
			local: stubLocalCapabilities{
				cfgs: map[string]config.CapabilityNodeConfig{
					"cron@1.0.0": stubCapabilityNodeConfig{},
				},
			},
			want: false,
		},
		{
			name: "mock trigger opted in",
			local: stubLocalCapabilities{
				cfgs: map[string]config.CapabilityNodeConfig{
					capStreams.MockTriggerCapabilityID: stubCapabilityNodeConfig{},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, shouldRegisterMockStreamsTrigger(tt.local))
		})
	}
}
