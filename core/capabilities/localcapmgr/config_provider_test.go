package localcapmgr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
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

func TestOffchainCapabilityConfigProvider(t *testing.T) {
	t.Parallel()

	gc := storedRegistry(t, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {
			"cron@1.0.0":      specConfigCap(t, map[string]any{"interval": "30"}),
			"consensus@1.0.0": {}, // present, but no spec_config
		},
	})
	reg, version := gc.LoadParsed()
	p := offchainCapabilityConfigProvider{reg: reg, donNames: map[uint32]string{7: testDONName(7), 99: testDONName(99)}, version: version}

	t.Run("returns offchain spec_config for the DON", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, map[string]any{"interval": "30"}, p.LocalConfigOverrides("cron@1.0.0", 7))
	})

	t.Run("nil for missing DON", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, p.LocalConfigOverrides("cron@1.0.0", 99))
	})

	t.Run("nil for DON without a usable name", func(t *testing.T) {
		t.Parallel()
		unnamed := offchainCapabilityConfigProvider{reg: reg, donNames: map[uint32]string{}, version: version}
		assert.Nil(t, unnamed.LocalConfigOverrides("cron@1.0.0", 7))
	})

	t.Run("nil for missing capability", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, p.LocalConfigOverrides("missing@1.0.0", 7))
	})

	t.Run("nil for capability without spec_config", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, p.LocalConfigOverrides("consensus@1.0.0", 7))
	})

	t.Run("nil snapshot", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, offchainCapabilityConfigProvider{}.LocalConfigOverrides("cron@1.0.0", 7))
	})

	t.Run("malformed spec_config is skipped, not partially applied", func(t *testing.T) {
		t.Parallel()
		// Unreachable from an applied payload (globalconfig.Validate rejects it at ingestion),
		// but the provider must still degrade to "no offchain override" rather than panic or
		// apply a partial map.
		bad := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": {SpecConfig: &valuespb.Map{Fields: map[string]*valuespb.Value{
				"ok":     valuespb.NewStringValue("x"),
				"broken": {}, // a Value with no kind cannot be converted
			}}}},
		})
		p := offchainCapabilityConfigProvider{reg: bad, donNames: map[uint32]string{7: testDONName(7)}, version: 1, lggr: testLogger(t)}
		assert.Nil(t, p.LocalConfigOverrides("cron@1.0.0", 7))
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
