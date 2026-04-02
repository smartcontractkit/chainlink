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

	require.False(t, shouldRegisterMockStreamsTrigger(nil))

	require.False(t, shouldRegisterMockStreamsTrigger(stubLocalCapabilities{
		cfgs: map[string]config.CapabilityNodeConfig{
			"cron@1.0.0": stubCapabilityNodeConfig{},
		},
	}))

	require.True(t, shouldRegisterMockStreamsTrigger(stubLocalCapabilities{
		cfgs: map[string]config.CapabilityNodeConfig{
			capStreams.MockTriggerCapabilityID: stubCapabilityNodeConfig{},
		},
	}))
}
