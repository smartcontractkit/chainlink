package localcapmgr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func specConfigCap(t *testing.T, kv map[string]any) *capabilitiespb.CapabilityConfig {
	t.Helper()
	vm, err := values.NewMap(kv)
	require.NoError(t, err)
	return &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(vm)}
}

func TestOffchainOverrides(t *testing.T) {
	t.Parallel()

	m := newCrossCheckMgr(t)
	reg := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {
			"cron@1.0.0":      specConfigCap(t, map[string]any{"interval": "30"}),
			"consensus@1.0.0": {}, // present, but no spec_config
		},
	})
	dons := m.offchainDONsByID(reg, []registry.DON{onchainDON(7), onchainDON(99)})

	t.Run("returns offchain spec_config for the DON", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, map[string]any{"interval": "30"}, m.offchainOverrides(dons[7], "cron@1.0.0", 7))
	})

	t.Run("nil for missing DON", func(t *testing.T) {
		t.Parallel()
		assert.NotContains(t, dons, uint32(99))
		assert.Nil(t, m.offchainOverrides(dons[99], "cron@1.0.0", 99))
	})

	t.Run("DON without a usable name is left out", func(t *testing.T) {
		t.Parallel()
		assert.Empty(t, m.offchainDONsByID(reg, []registry.DON{{ID: 7}}))
	})

	t.Run("nil for missing capability", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, m.offchainOverrides(dons[7], "missing@1.0.0", 7))
	})

	t.Run("nil for capability without spec_config", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, m.offchainOverrides(dons[7], "consensus@1.0.0", 7))
	})

	t.Run("nil snapshot", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, m.offchainDONsByID(nil, []registry.DON{onchainDON(7)}))
	})

	t.Run("malformed spec_config is skipped, not partially applied", func(t *testing.T) {
		t.Parallel()
		// Unreachable from an applied payload (globalconfig.Validate rejects it at ingestion).
		bad := &capabilitiespb.DONConfig{Capabilities: map[string]*capabilitiespb.CapabilityConfig{
			"cron@1.0.0": {SpecConfig: &valuespb.Map{Fields: map[string]*valuespb.Value{
				"ok":     valuespb.NewStringValue("x"),
				"broken": {}, // a Value with no kind cannot be converted
			}}},
		}}
		assert.Nil(t, m.offchainOverrides(bad, "cron@1.0.0", 7))
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
