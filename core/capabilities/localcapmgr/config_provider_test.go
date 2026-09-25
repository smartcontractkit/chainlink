package localcapmgr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func specConfigCap(t *testing.T, kv map[string]any) *capabilitiespb.CapabilityConfig {
	t.Helper()
	vm, err := values.NewMap(kv)
	require.NoError(t, err)
	return &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(vm)}
}

func storedRegistry(t *testing.T, version uint64, dons map[uint32]map[string]*capabilitiespb.CapabilityConfig) *globalconfig.GlobalConfig {
	t.Helper()
	raw, err := marshalOffchainRegistry(offchainReg(version, dons))
	require.NoError(t, err)
	gc := globalconfig.New()
	require.NoError(t, gc.Store(globalconfig.Update{Raw: raw, Hash: fmt.Sprintf("h%d", version)}))
	return gc
}

func TestTomlCapabilityConfigProvider(t *testing.T) {
	t.Parallel()

	t.Run("nil localCfg returns nil", func(t *testing.T) {
		t.Parallel()
		p := tomlCapabilityConfigProvider{}
		assert.Nil(t, p.LocalConfigOverrides("cron@1.0.0", 1))
	})

	t.Run("returns capability config map", func(t *testing.T) {
		t.Parallel()
		p := tomlCapabilityConfigProvider{localCfg: &testLocalCapabilities{
			configs: map[string]*testCapabilityNodeConfig{
				"cron@1.0.0": {cfg: map[string]string{"interval": "60"}},
			},
		}}
		assert.Equal(t, map[string]any{"interval": "60"}, p.LocalConfigOverrides("cron@1.0.0", 1))
	})

	t.Run("unknown capability returns nil", func(t *testing.T) {
		t.Parallel()
		p := tomlCapabilityConfigProvider{localCfg: &testLocalCapabilities{}}
		assert.Nil(t, p.LocalConfigOverrides("missing@1.0.0", 1))
	})
}

// stubConfigProvider lets tests drive buildConfigJSON through the seam directly.
type stubConfigProvider struct {
	overrides map[string]map[string]any
}

func (s stubConfigProvider) LocalConfigOverrides(capID string, _ uint32) map[string]any {
	return s.overrides[capID]
}

func TestOffchainCapabilityConfigProvider(t *testing.T) {
	t.Parallel()

	gc := storedRegistry(t, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})},
	})

	t.Run("returns offchain spec_config for the DON", func(t *testing.T) {
		t.Parallel()
		p := offchainCapabilityConfigProvider{registry: gc}
		assert.Equal(t, map[string]any{"interval": "30"}, p.LocalConfigOverrides("cron@1.0.0", 7))
	})

	t.Run("nil for wrong DON", func(t *testing.T) {
		t.Parallel()
		p := offchainCapabilityConfigProvider{registry: gc}
		assert.Nil(t, p.LocalConfigOverrides("cron@1.0.0", 99))
	})

	t.Run("nil for unknown capability", func(t *testing.T) {
		t.Parallel()
		p := offchainCapabilityConfigProvider{registry: gc}
		assert.Nil(t, p.LocalConfigOverrides("missing@1.0.0", 7))
	})

	t.Run("nil registry / nothing applied", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, offchainCapabilityConfigProvider{registry: nil}.LocalConfigOverrides("cron@1.0.0", 7))
		assert.Nil(t, offchainCapabilityConfigProvider{registry: globalconfig.New()}.LocalConfigOverrides("cron@1.0.0", 7))
	})
}

func TestLayeredCapabilityConfigProvider_OffchainWins(t *testing.T) {
	t.Parallel()

	toml := tomlCapabilityConfigProvider{localCfg: &testLocalCapabilities{
		configs: map[string]*testCapabilityNodeConfig{
			"cron@1.0.0": {cfg: map[string]string{"interval": "60", "tomlOnly": "keep"}},
		},
	}}
	offchain := offchainCapabilityConfigProvider{registry: storedRegistry(t, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30", "offchainOnly": "add"})},
	})}

	layered := layeredCapabilityConfigProvider{toml: toml, offchain: offchain}

	got := layered.LocalConfigOverrides("cron@1.0.0", 7)
	assert.Equal(t, map[string]any{
		"interval":     "30",   // offchain wins over TOML's "60"
		"tomlOnly":     "keep", // TOML-only key preserved
		"offchainOnly": "add",  // offchain-only key added
	}, got)

	// With no offchain entry for the DON, TOML is returned unchanged.
	assert.Equal(t, map[string]any{"interval": "60", "tomlOnly": "keep"}, layered.LocalConfigOverrides("cron@1.0.0", 99))
}

func TestNewLocalCapabilityManager_CutoverGate(t *testing.T) {
	t.Parallel()

	localCfg := &testLocalCapabilities{configs: map[string]*testCapabilityNodeConfig{
		"cron@1.0.0": {cfg: map[string]string{"interval": "60"}},
	}}
	gc := storedRegistry(t, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})},
	})
	noop := func(_ context.Context, _ string, _ uint32, _ string, _ string, _ *ocrtypes.ContractConfig) ([]job.ServiceCtx, error) {
		return nil, nil
	}

	t.Run("gate off uses TOML only", func(t *testing.T) {
		t.Parallel()
		m, err := NewLocalCapabilityManager(testLogger(t), localCfg, noop, gc, false)
		require.NoError(t, err)
		lcm := m.(*localCapabilityManager)
		assert.Equal(t, map[string]any{"interval": "60"}, lcm.overridesFor("cron@1.0.0", 7))
	})

	t.Run("gate on layers offchain over TOML", func(t *testing.T) {
		t.Parallel()
		m, err := NewLocalCapabilityManager(testLogger(t), localCfg, noop, gc, true)
		require.NoError(t, err)
		lcm := m.(*localCapabilityManager)
		assert.Equal(t, map[string]any{"interval": "30"}, lcm.overridesFor("cron@1.0.0", 7))
	})

	t.Run("gate on but nil registry falls back to TOML", func(t *testing.T) {
		t.Parallel()
		m, err := NewLocalCapabilityManager(testLogger(t), localCfg, noop, nil, true)
		require.NoError(t, err)
		lcm := m.(*localCapabilityManager)
		assert.Equal(t, map[string]any{"interval": "60"}, lcm.overridesFor("cron@1.0.0", 7))
	})
}

func TestBuildConfigJSON_UsesConfigProvider(t *testing.T) {
	t.Parallel()

	mgr := &localCapabilityManager{
		lggr: testLogger(t),
		configProvider: stubConfigProvider{overrides: map[string]map[string]any{
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
