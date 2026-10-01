package localcapmgr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
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

func TestComputeOffchainCrossCheck(t *testing.T) {
	t.Parallel()

	emptyCap := map[string]*capabilitiespb.CapabilityConfig{"cron@1.0.0": {}, "consensus@1.0.0": {}}

	t.Run("nil offchain registry is empty", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t, "cron@1.0.0")
		got := mgr.computeOffchainCrossCheck(nil, 0, []registry.DON{onchainDON(1, "cron@1.0.0")})
		assert.True(t, got.offchainEmpty)
		assert.Equal(t, uint64(0), got.version)
	})

	t.Run("full match records matched caps and no divergence", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t, "cron@1.0.0", "consensus@1.0.0")
		reg := offchainReg(5, map[uint32]map[string]*capabilitiespb.CapabilityConfig{1: emptyCap})
		got := mgr.computeOffchainCrossCheck(reg, 5, []registry.DON{onchainDON(1, "cron@1.0.0", "consensus@1.0.0")})

		assert.False(t, got.offchainEmpty)
		assert.Equal(t, uint64(5), got.version)
		assert.Equal(t, int64(1), got.comparedDONs)
		assert.Equal(t, int64(2), got.matchedCaps)
		assert.Empty(t, got.divergences)
	})

	t.Run("only allowlisted caps are compared", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t, "cron@1.0.0") // consensus not allowlisted
		reg := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{1: {"cron@1.0.0": {}}})
		// on-chain DON offers both; only cron is allowlisted, and offchain has cron -> match.
		got := mgr.computeOffchainCrossCheck(reg, 1, []registry.DON{onchainDON(1, "cron@1.0.0", "consensus@1.0.0")})
		assert.Equal(t, int64(1), got.matchedCaps)
		assert.Empty(t, got.divergences)
	})

	t.Run("missing DON", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t, "cron@1.0.0")
		reg := offchainReg(1, nil) // no DONs offchain
		got := mgr.computeOffchainCrossCheck(reg, 1, []registry.DON{onchainDON(1, "cron@1.0.0")})
		assert.Equal(t, int64(1), got.divergences[divergenceMissingDON])
		assert.Equal(t, int64(0), got.matchedCaps)
	})

	t.Run("missing capability", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t, "cron@1.0.0", "consensus@1.0.0")
		reg := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{1: {"cron@1.0.0": {}}})
		got := mgr.computeOffchainCrossCheck(reg, 1, []registry.DON{onchainDON(1, "cron@1.0.0", "consensus@1.0.0")})
		assert.Equal(t, int64(1), got.matchedCaps)
		assert.Equal(t, int64(1), got.divergences[divergenceMissingCapability])
	})

	t.Run("extra DON and extra capability", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t, "cron@1.0.0")
		reg := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			1: {"cron@1.0.0": {}, "ghost@1.0.0": {}}, // ghost not on-chain
			9: {"cron@1.0.0": {}},                    // DON 9 not on-chain for this node
		})
		got := mgr.computeOffchainCrossCheck(reg, 1, []registry.DON{onchainDON(1, "cron@1.0.0")})
		assert.Equal(t, int64(1), got.matchedCaps)
		assert.Equal(t, int64(1), got.divergences[divergenceExtraDON])
		assert.Equal(t, int64(1), got.divergences[divergenceExtraCapability])
	})

	t.Run("DON with no allowlisted caps is skipped", func(t *testing.T) {
		t.Parallel()
		mgr := newCrossCheckMgr(t) // nothing allowlisted
		reg := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{1: emptyCap})
		got := mgr.computeOffchainCrossCheck(reg, 1, []registry.DON{onchainDON(1, "cron@1.0.0")})
		assert.Equal(t, int64(0), got.comparedDONs)
		assert.Equal(t, int64(0), got.matchedCaps)
		assert.Equal(t, int64(0), got.divergences[divergenceMissingDON])
	})
}

func TestCrossValidateOffchain_NilRegistryIsNoop(t *testing.T) {
	t.Parallel()
	mgr := newCrossCheckMgr(t, "cron@1.0.0")
	mgr.offchainRegistry = nil
	// Should not panic and should not require metrics to be set.
	mgr.crossValidateOffchain(t.Context(), []registry.DON{onchainDON(1, "cron@1.0.0")})
}

func TestCrossValidateOffchain_EndToEndFromStoredPayload(t *testing.T) {
	t.Parallel()

	gc := globalconfig.New()
	// Store a proto-JSON payload via the same path the delegate uses.
	raw, err := marshalOffchainRegistry(offchainReg(3, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		1: {"cron@1.0.0": {}},
	}))
	require.NoError(t, err)
	require.NoError(t, gc.Store(globalconfig.Update{Raw: raw, Hash: "h1"}))

	metrics, err := newMetrics()
	require.NoError(t, err)
	mgr := newCrossCheckMgr(t, "cron@1.0.0")
	mgr.metrics = metrics
	mgr.offchainRegistry = gc

	// Exercises LoadParsed -> compute -> record without panicking.
	mgr.crossValidateOffchain(t.Context(), []registry.DON{onchainDON(1, "cron@1.0.0")})

	reg, version := gc.LoadParsed()
	require.NotNil(t, reg)
	assert.Equal(t, uint64(3), version)
}

func TestComputeOffchainCrossCheck_ConfigMismatch(t *testing.T) {
	t.Parallel()

	mgr := &localCapabilityManager{
		lggr: testLogger(t),
		localCfg: &testLocalCapabilities{
			allowlisted: map[string]bool{"cron@1.0.0": true},
			configs: map[string]*testCapabilityNodeConfig{
				"cron@1.0.0": {cfg: map[string]string{"tomlOnly": "t"}},
			},
		},
	}
	onchain := []registry.DON{{
		ID:   1,
		Name: testDONName(1),
		CapabilityConfigurations: map[string]registry.CapabilityConfiguration{
			"cron@1.0.0": onchainSpecConfig(t, map[string]any{"interval": "20", "onchainOnly": "o"}),
		},
	}}
	check := func(offchainSpec map[string]any) offchainCrossCheck {
		cc := &capabilitiespb.CapabilityConfig{}
		if offchainSpec != nil {
			cc = specConfigCap(t, offchainSpec)
		}
		reg := offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{1: {"cron@1.0.0": cc}})
		return mgr.computeOffchainCrossCheck(reg, 1, onchain)
	}

	for _, tc := range []struct {
		name     string
		offchain map[string]any
		mismatch int64
	}{
		{"no spec_config", nil, 0},
		{"same value as on-chain", map[string]any{"interval": "20"}, 0},
		{"same value as TOML", map[string]any{"tomlOnly": "t"}, 0},
		{"overrides on-chain value", map[string]any{"interval": "30"}, 1},
		{"overrides TOML value", map[string]any{"tomlOnly": "x"}, 1},
		{"adds a key", map[string]any{"interval": "20", "extra": "e"}, 1},
		{"different type, same text", map[string]any{"interval": int64(20)}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := check(tc.offchain)
			assert.Equal(t, int64(1), got.matchedCaps)
			assert.Equal(t, tc.mismatch, got.divergences[divergenceConfigMismatch])
		})
	}
}
