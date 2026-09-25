package localcapmgr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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

// onchainSpecConfig builds a registry.CapabilityConfiguration whose SpecConfig unwraps to kv,
// for exercising buildConfigJSON's on-chain layer.
func onchainSpecConfig(t *testing.T, kv map[string]any) registry.CapabilityConfiguration {
	t.Helper()
	vm, err := values.NewMap(kv)
	require.NoError(t, err)
	cc := &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(vm)}
	b, err := proto.Marshal(cc)
	require.NoError(t, err)
	return registry.CapabilityConfiguration{Config: b}
}

// TestBuildConfigJSON_Precedence verifies the backwards-compatible cutover layering:
// TOML (base) < on-chain SpecConfig < offchain SpecConfig (only when the gate is on).
func TestBuildConfigJSON_Precedence(t *testing.T) {
	t.Parallel()

	localCfg := &testLocalCapabilities{configs: map[string]*testCapabilityNodeConfig{
		"cron@1.0.0": {cfg: map[string]string{"interval": "10", "tomlOnly": "keep"}},
	}}
	onchain := onchainSpecConfig(t, map[string]any{"interval": "20", "onchainOnly": "oc"})
	gc := storedRegistry(t, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30", "offchainOnly": "add"})},
	})
	noop := func(_ context.Context, _ string, _ uint32, _ string, _ string, _ *ocrtypes.ContractConfig) ([]job.ServiceCtx, error) {
		return nil, nil
	}
	info := &capabilityInfo{capID: "cron@1.0.0", donID: 7, config: onchain}

	t.Run("gate off: on-chain wins over TOML, offchain ignored", func(t *testing.T) {
		t.Parallel()
		m, err := NewLocalCapabilityManager(testLogger(t), localCfg, noop, gc, false)
		require.NoError(t, err)
		got, err := m.(*localCapabilityManager).buildConfigJSON(info)
		require.NoError(t, err)
		assert.JSONEq(t, `{"interval":"20","tomlOnly":"keep","onchainOnly":"oc"}`, got)
	})

	t.Run("gate on: offchain wins, absent keys fall back to on-chain/TOML", func(t *testing.T) {
		t.Parallel()
		m, err := NewLocalCapabilityManager(testLogger(t), localCfg, noop, gc, true)
		require.NoError(t, err)
		got, err := m.(*localCapabilityManager).buildConfigJSON(info)
		require.NoError(t, err)
		assert.JSONEq(t, `{"interval":"30","offchainOnly":"add","onchainOnly":"oc","tomlOnly":"keep"}`, got)
	})

	t.Run("gate on but offchain has no entry for this DON: falls back to on-chain/TOML", func(t *testing.T) {
		t.Parallel()
		m, err := NewLocalCapabilityManager(testLogger(t), localCfg, noop, gc, true)
		require.NoError(t, err)
		info99 := &capabilityInfo{capID: "cron@1.0.0", donID: 99, config: onchain}
		got, err := m.(*localCapabilityManager).buildConfigJSON(info99)
		require.NoError(t, err)
		assert.JSONEq(t, `{"interval":"20","tomlOnly":"keep","onchainOnly":"oc"}`, got)
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
