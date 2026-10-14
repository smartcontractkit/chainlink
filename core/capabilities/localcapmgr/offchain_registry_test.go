package localcapmgr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func marshalOffchainRegistry(reg *capabilitiespb.OffchainCapabilitiesRegistry) (string, error) {
	b, err := protojson.Marshal(reg)
	return string(b), err
}

// testDONName is the on-chain name test DONs are given; the offchain registry is keyed by it.
func testDONName(id uint32) string { return fmt.Sprintf("don-%d", id) }

// offchainReg builds an offchain registry keyed by testDONName of each DON ID.
func offchainReg(version uint64, dons map[uint32]map[string]*capabilitiespb.CapabilityConfig) *capabilitiespb.OffchainCapabilitiesRegistry {
	reg := &capabilitiespb.OffchainCapabilitiesRegistry{Domain: "cre", Env: "test", Version: version, Dons: map[string]*capabilitiespb.DONConfig{}}
	for donID, caps := range dons {
		reg.Dons[testDONName(donID)] = &capabilitiespb.DONConfig{Capabilities: caps}
	}
	return reg
}

func onchainDON(id uint32, capIDs ...string) registry.DON {
	cfgs := map[string]registry.CapabilityConfiguration{}
	for _, c := range capIDs {
		cfgs[c] = registry.CapabilityConfiguration{}
	}
	return registry.DON{ID: id, Name: testDONName(id), CapabilityConfigurations: cfgs}
}

func newCrossCheckMgr(t *testing.T, allowlisted ...string) *localCapabilityManager {
	al := map[string]bool{}
	for _, c := range allowlisted {
		al[c] = true
	}
	return &localCapabilityManager{
		lggr:     testLogger(t),
		localCfg: &testLocalCapabilities{allowlisted: al},
	}
}

func specConfigCap(t *testing.T, kv map[string]any) *capabilitiespb.CapabilityConfig {
	t.Helper()
	vm, err := values.NewMap(kv)
	require.NoError(t, err)
	return &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(vm)}
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

func TestOffchainDivergences(t *testing.T) {
	t.Parallel()

	m := newCrossCheckMgr(t, "cron@1.0.0", "consensus@1.0.0")
	onchain := onchainSpecConfig(t, map[string]any{"interval": "20"})
	don := func(id uint32) registry.DON {
		return onchainDONWith(id, map[string]registry.CapabilityConfiguration{
			"cron@1.0.0": onchain, "consensus@1.0.0": onchain, "other@1.0.0": onchain, // other is not allowlisted
		})
	}
	dons := []registry.DON{don(1), don(2)}

	for _, tc := range []struct {
		name     string
		offchain map[uint32]map[string]*capabilitiespb.CapabilityConfig
		want     map[string]int64
	}{
		{"no payload", nil, map[string]int64{divergenceMissing: 4, divergenceMismatch: 0}},
		{"same as on-chain", map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			1: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "20"}), "consensus@1.0.0": specConfigCap(t, map[string]any{"interval": "20"})},
			2: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "20"}), "consensus@1.0.0": specConfigCap(t, map[string]any{"interval": "20"})},
		}, map[string]int64{divergenceMissing: 0, divergenceMismatch: 0}},
		{"missing DON, missing capability and mismatches", map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			1: {
				"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"}), // different value
				"x@1.0.0":    {},                                                 // extra entries are ignored
			},
			3: {"cron@1.0.0": {}}, // DON this node is not in
		}, map[string]int64{divergenceMissing: 3, divergenceMismatch: 1}},
		{"omitted key and empty spec_config", map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			1: {"cron@1.0.0": specConfigCap(t, map[string]any{}), "consensus@1.0.0": {}},
			2: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": int64(20)}), "consensus@1.0.0": specConfigCap(t, map[string]any{"interval": "20", "extra": "e"})},
		}, map[string]int64{divergenceMissing: 0, divergenceMismatch: 4}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var reg *capabilitiespb.OffchainCapabilitiesRegistry
			if tc.offchain != nil {
				reg = offchainReg(1, tc.offchain)
			}
			assert.Equal(t, tc.want, m.offchainDivergences(m.offchainDONsByID(reg, dons), dons))
		})
	}
}

func TestOffchainSpecConfig_MalformedIsIgnored(t *testing.T) {
	t.Parallel()
	// Unreachable from an applied payload (globalconfig.Validate rejects it at ingestion).
	bad := &capabilitiespb.DONConfig{Capabilities: map[string]*capabilitiespb.CapabilityConfig{
		"cron@1.0.0": {SpecConfig: &valuespb.Map{Fields: map[string]*valuespb.Value{
			"ok":     valuespb.NewStringValue("x"),
			"broken": {}, // a Value with no kind cannot be converted
		}}},
	}}
	assert.Nil(t, newCrossCheckMgr(t).offchainSpecConfig(bad, "cron@1.0.0", 7))
}
