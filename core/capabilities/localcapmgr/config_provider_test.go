package localcapmgr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
)

func TestTomlCapabilityConfigProvider(t *testing.T) {
	t.Parallel()

	t.Run("nil localCfg returns nil", func(t *testing.T) {
		t.Parallel()
		p := tomlCapabilityConfigProvider{}
		assert.Nil(t, p.LocalConfigOverrides("cron@1.0.0"))
	})

	t.Run("returns capability config map", func(t *testing.T) {
		t.Parallel()
		p := tomlCapabilityConfigProvider{localCfg: &testLocalCapabilities{
			configs: map[string]*testCapabilityNodeConfig{
				"cron@1.0.0": {cfg: map[string]string{"interval": "60"}},
			},
		}}
		assert.Equal(t, map[string]string{"interval": "60"}, p.LocalConfigOverrides("cron@1.0.0"))
	})

	t.Run("unknown capability returns nil", func(t *testing.T) {
		t.Parallel()
		p := tomlCapabilityConfigProvider{localCfg: &testLocalCapabilities{}}
		assert.Nil(t, p.LocalConfigOverrides("missing@1.0.0"))
	})
}

// stubConfigProvider lets tests drive buildConfigJSON through the seam directly.
type stubConfigProvider struct {
	overrides map[string]map[string]string
}

func (s stubConfigProvider) LocalConfigOverrides(capID string) map[string]string {
	return s.overrides[capID]
}

func TestBuildConfigJSON_UsesConfigProvider(t *testing.T) {
	t.Parallel()

	mgr := &localCapabilityManager{
		lggr: testLogger(t),
		configProvider: stubConfigProvider{overrides: map[string]map[string]string{
			"cron@1.0.0": {"interval": "60"},
		}},
	}

	got, err := mgr.buildConfigJSON(&capabilityInfo{
		capID:  "cron@1.0.0",
		config: registry.CapabilityConfiguration{},
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{"interval":"60"}`, got)
}
