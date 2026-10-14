package localcapmgr

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// launch records one newServicesFn call.
type launch struct {
	capID      string
	donID      uint32
	command    string
	configJSON string
	ocr3       *ocrtypes.ContractConfig
	svc        *syncService
}

// launchRecorder is a NewServicesFn that records every launch, safe for concurrent use.
type launchRecorder struct {
	mu       sync.Mutex
	launches []launch
}

func (r *launchRecorder) newServices(_ context.Context, capID string, donID uint32, command string, configJSON string, ocr3 *ocrtypes.ContractConfig) ([]job.ServiceCtx, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc := &syncService{}
	r.launches = append(r.launches, launch{capID: capID, donID: donID, command: command, configJSON: configJSON, ocr3: ocr3, svc: svc})
	return []job.ServiceCtx{svc}, nil
}

func (r *launchRecorder) all() []launch {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]launch(nil), r.launches...)
}

func (r *launchRecorder) last(t *testing.T) launch {
	t.Helper()
	all := r.all()
	require.NotEmpty(t, all)
	return all[len(all)-1]
}

func (r *launchRecorder) configFor(t *testing.T, capID string, donID uint32) map[string]any {
	t.Helper()
	for _, l := range slices.Backward(r.all()) {
		if l.capID == capID && l.donID == donID {
			var got map[string]any
			require.NoError(t, json.Unmarshal([]byte(l.configJSON), &got))
			return got
		}
	}
	require.Failf(t, "not launched", "%s on DON %d", capID, donID)
	return nil
}

type syncService struct {
	mu     sync.Mutex
	closed bool
}

func (s *syncService) Start(context.Context) error { return nil }
func (s *syncService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}
func (s *syncService) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func storeOffchain(t *testing.T, gc *globalconfig.GlobalConfig, version uint64, dons map[uint32]map[string]*capabilitiespb.CapabilityConfig) {
	t.Helper()
	raw, err := marshalOffchainRegistry(offchainReg(version, dons))
	require.NoError(t, err)
	require.NoError(t, gc.Store(globalconfig.Update{Raw: raw}))
}

// cutoverManager reconciles against the current snapshot of gc, as the Launcher does.
type cutoverManager struct {
	*localCapabilityManager
	gc *globalconfig.GlobalConfig
}

func (c cutoverManager) reconcile(ctx context.Context, dons []registry.DON) error {
	reg, _ := c.gc.LoadParsed()
	return c.Reconcile(ctx, dons, reg)
}

func newCutoverManager(t *testing.T, localCfg *testLocalCapabilities, gc *globalconfig.GlobalConfig, gate bool) (cutoverManager, *launchRecorder) {
	t.Helper()
	rec := &launchRecorder{}
	m, err := NewLocalCapabilityManager(testLogger(t), localCfg, rec.newServices, gate)
	require.NoError(t, err)
	return cutoverManager{localCapabilityManager: m.(*localCapabilityManager), gc: gc}, rec
}

func onchainDONWith(id uint32, caps map[string]registry.CapabilityConfiguration) registry.DON {
	return registry.DON{ID: id, Name: testDONName(id), CapabilityConfigurations: caps}
}

// TestReconcile_EffectiveSpecConfigPrecedence drives the full launch path (Reconcile ->
// buildDesiredState -> startCapability -> newServicesFn) and asserts the config the capability
// is actually launched with: on-chain only with the gate off, offchain only with the gate on.
func TestReconcile_EffectiveSpecConfigPrecedence(t *testing.T) {
	t.Parallel()

	localCfg := func() *testLocalCapabilities {
		return &testLocalCapabilities{
			allowlisted: map[string]bool{"cron@1.0.0": true, "consensus@1.0.0": true},
		}
	}
	onchain := onchainSpecConfig(t, map[string]any{"interval": "20", "onchainOnly": "oc"})
	dons := []registry.DON{
		onchainDONWith(7, map[string]registry.CapabilityConfiguration{"cron@1.0.0": onchain, "consensus@1.0.0": onchain}),
		onchainDONWith(8, map[string]registry.CapabilityConfiguration{"cron@1.0.0": onchain}),
	}
	newGC := func(t *testing.T) *globalconfig.GlobalConfig {
		gc := globalconfig.New()
		storeOffchain(t, gc, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			// DON 7 overrides cron; consensus is absent. DON 8 is absent entirely.
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30", "offchainOnly": "add"})},
		})
		return gc
	}
	legacy := map[string]any{"interval": "20", "onchainOnly": "oc"}
	empty := map[string]any{}

	t.Run("gate off: offchain ignored, legacy behavior preserved", func(t *testing.T) {
		t.Parallel()
		m, rec := newCutoverManager(t, localCfg(), newGC(t), false)
		require.NoError(t, m.reconcile(t.Context(), dons))
		assert.Equal(t, legacy, rec.configFor(t, "cron@1.0.0", 7))
		assert.Equal(t, legacy, rec.configFor(t, "cron@1.0.0", 8))
	})

	t.Run("gate on: offchain only, on-chain keys are not merged", func(t *testing.T) {
		t.Parallel()
		m, rec := newCutoverManager(t, localCfg(), newGC(t), true)
		require.NoError(t, m.reconcile(t.Context(), dons))
		assert.Equal(t, map[string]any{"interval": "30", "offchainOnly": "add"}, rec.configFor(t, "cron@1.0.0", 7))
		// missing DON or capability in the offchain payload -> empty config
		assert.Equal(t, empty, rec.configFor(t, "cron@1.0.0", 8))
		assert.Equal(t, empty, rec.configFor(t, "consensus@1.0.0", 7))
	})

	t.Run("gate on, nothing applied yet: empty config", func(t *testing.T) {
		t.Parallel()
		m, rec := newCutoverManager(t, localCfg(), globalconfig.New(), true)
		require.NoError(t, m.reconcile(t.Context(), dons))
		assert.Equal(t, empty, rec.configFor(t, "cron@1.0.0", 7))
	})

	t.Run("gate on: launch allowlist stays node-local", func(t *testing.T) {
		t.Parallel()
		// Offchain configures a capability that is on-chain but NOT allowlisted on this node.
		gc := globalconfig.New()
		storeOffchain(t, gc, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"other@1.0.0": specConfigCap(t, map[string]any{"k": "v"})},
		})
		m, rec := newCutoverManager(t, localCfg(), gc, true)
		require.NoError(t, m.reconcile(t.Context(), []registry.DON{
			onchainDONWith(7, map[string]registry.CapabilityConfiguration{"other@1.0.0": onchain}),
		}))
		assert.Empty(t, rec.all())
	})
}

// TestReconcile_OCR3StaysOnchain asserts that, for this slice, the OCR3 config handed to the
// capability always comes from the on-chain registry; ocr3_configs in the offchain payload are
// not consumed.
func TestReconcile_OCR3StaysOnchain(t *testing.T) {
	t.Parallel()

	onchainRaw, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: {Signers: [][]byte{{0x01}}, Transmitters: [][]byte{{0xaa}}, F: 1},
		},
	})
	require.NoError(t, err)
	gc := globalconfig.New()
	storeOffchain(t, gc, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		7: {"consensus@1.0.0": {
			Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
				capabilitiespb.OCR3ConfigDefaultKey: {Signers: [][]byte{{0x09}}, Transmitters: [][]byte{{0xbb}}, F: 2},
			},
		}},
	})
	m, rec := newCutoverManager(t, &testLocalCapabilities{allowlisted: map[string]bool{"consensus@1.0.0": true}}, gc, true)
	require.NoError(t, m.reconcile(t.Context(), []registry.DON{
		onchainDONWith(7, map[string]registry.CapabilityConfiguration{"consensus@1.0.0": {Config: onchainRaw}}),
	}))

	l := rec.last(t)
	require.NotNil(t, l.ocr3)
	assert.Equal(t, uint8(1), l.ocr3.F)
	assert.Equal(t, ocrtypes.OnchainPublicKey{0x01}, l.ocr3.Signers[0])
}

// TestReconcile_OffchainOnlyUpdateRestartsCapability asserts that a change confined to the
// offchain payload (on-chain unchanged) restarts exactly the affected capability with the new
// effective config, and that re-reconciling an unchanged state is a no-op.
func TestReconcile_OffchainOnlyUpdateRestartsCapability(t *testing.T) {
	t.Parallel()

	localCfg := &testLocalCapabilities{allowlisted: map[string]bool{"cron@1.0.0": true, "consensus@1.0.0": true}}
	onchain := onchainSpecConfig(t, map[string]any{"interval": "20"})
	dons := []registry.DON{onchainDONWith(7, map[string]registry.CapabilityConfiguration{
		"cron@1.0.0": onchain, "consensus@1.0.0": onchain,
	})}

	t.Run("gate on", func(t *testing.T) {
		t.Parallel()
		gc := globalconfig.New()
		storeOffchain(t, gc, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})},
		})
		m, rec := newCutoverManager(t, localCfg, gc, true)
		require.NoError(t, m.reconcile(t.Context(), dons))
		require.Len(t, rec.all(), 2)
		first := map[string]launch{}
		for _, l := range rec.all() {
			first[l.capID] = l
		}

		// Unchanged on-chain + unchanged offchain: nothing restarts.
		require.NoError(t, m.reconcile(t.Context(), dons))
		require.Len(t, rec.all(), 2)

		// Offchain-only update for cron.
		storeOffchain(t, gc, 2, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "40"})},
		})
		require.NoError(t, m.reconcile(t.Context(), dons))

		all := rec.all()
		require.Len(t, all, 3, "only the affected capability is restarted")
		assert.Equal(t, "cron@1.0.0", all[2].capID)
		assert.Equal(t, map[string]any{"interval": "40"}, rec.configFor(t, "cron@1.0.0", 7))
		assert.True(t, first["cron@1.0.0"].svc.isClosed(), "old cron instance must be stopped")
		assert.False(t, first["consensus@1.0.0"].svc.isClosed(), "unaffected capability keeps running")

		// A new payload version that does not change cron's effective config does not restart it.
		storeOffchain(t, gc, 3, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "40"})},
		})
		require.NoError(t, m.reconcile(t.Context(), dons))
		require.Len(t, rec.all(), 3)

		// Removing the offchain entry restarts with an empty config (no on-chain fallback).
		storeOffchain(t, gc, 4, map[uint32]map[string]*capabilitiespb.CapabilityConfig{})
		require.NoError(t, m.reconcile(t.Context(), dons))
		require.Len(t, rec.all(), 4)
		assert.Equal(t, map[string]any{}, rec.configFor(t, "cron@1.0.0", 7))
	})

	t.Run("gate off: offchain-only update does not restart", func(t *testing.T) {
		t.Parallel()
		gc := globalconfig.New()
		storeOffchain(t, gc, 1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})},
		})
		m, rec := newCutoverManager(t, localCfg, gc, false)
		require.NoError(t, m.reconcile(t.Context(), dons))
		storeOffchain(t, gc, 2, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "40"})},
		})
		require.NoError(t, m.reconcile(t.Context(), dons))
		require.Len(t, rec.all(), 2)
	})
}

// TestReconcile_ConcurrentOffchainUpdates exercises reconciles racing offchain stores. Run with
// -race. After the dust settles the running
// capability must reflect the latest payload.
func TestReconcile_ConcurrentOffchainUpdates(t *testing.T) {
	t.Parallel()

	localCfg := &testLocalCapabilities{allowlisted: map[string]bool{"cron@1.0.0": true}}
	dons := []registry.DON{onchainDONWith(7, map[string]registry.CapabilityConfiguration{
		"cron@1.0.0": onchainSpecConfig(t, map[string]any{"interval": "20"}),
	})}
	gc := globalconfig.New()
	m, rec := newCutoverManager(t, localCfg, gc, true)
	require.NoError(t, m.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, m.Close()) })

	const n = 30
	payloads := make([]string, 0, n)
	for v := uint64(1); v <= n; v++ {
		raw, err := marshalOffchainRegistry(offchainReg(v, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
			7: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": strconv.FormatUint(v, 10)})},
		}))
		require.NoError(t, err)
		payloads = append(payloads, raw)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, raw := range payloads {
			assert.NoError(t, gc.Store(globalconfig.Update{Raw: raw}))
		}
	}()
	go func() {
		defer wg.Done()
		for range n {
			assert.NoError(t, m.reconcile(t.Context(), dons))
		}
	}()
	wg.Wait()

	// A final reconcile is deterministic regardless of interleaving.
	require.NoError(t, m.reconcile(t.Context(), dons))
	assert.Equal(t, map[string]any{"interval": strconv.Itoa(n)}, rec.configFor(t, "cron@1.0.0", 7))

	m.mu.RLock()
	defer m.mu.RUnlock()
	require.Len(t, m.runningCapabilities, 1)
}

// TestReconcile_OffchainMatchesByDONName covers the design's name-keyed DON map: offchain config
// applies to the on-chain DON with that name, regardless of its ID; DONs without a usable name
// (unnamed, or a name shared by two of this node's DONs) get no offchain config.
func TestReconcile_OffchainMatchesByDONName(t *testing.T) {
	t.Parallel()

	onchain := onchainSpecConfig(t, map[string]any{"interval": "20"})
	gc := globalconfig.New()
	raw, err := marshalOffchainRegistry(&capabilitiespb.OffchainCapabilitiesRegistry{
		Domain: "cre", Env: "test", Version: 1,
		Dons: map[string]*capabilitiespb.DONConfig{
			"workflow_1_zone-a": {Capabilities: map[string]*capabilitiespb.CapabilityConfig{"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})}},
			"shared":            {Capabilities: map[string]*capabilitiespb.CapabilityConfig{"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "99"})}},
		},
	})
	require.NoError(t, err)
	require.NoError(t, gc.Store(globalconfig.Update{Raw: raw}))

	named := func(id uint32, name string) registry.DON {
		return registry.DON{ID: id, Name: name, CapabilityConfigurations: map[string]registry.CapabilityConfiguration{"cron@1.0.0": onchain}}
	}
	m, rec := newCutoverManager(t, &testLocalCapabilities{allowlisted: map[string]bool{"cron@1.0.0": true}}, gc, true)
	require.NoError(t, m.reconcile(t.Context(), []registry.DON{
		named(42, "workflow_1_zone-a"), // matched by name; the ID is irrelevant
		named(43, ""),                  // unnamed (e.g. registry without DON names)
		named(44, "shared"),            // ambiguous: two of this node's DONs share the name
		named(45, "shared"),
	}))

	assert.Equal(t, map[string]any{"interval": "30"}, rec.configFor(t, "cron@1.0.0", 42))
	for _, id := range []uint32{43, 44, 45} {
		assert.Equal(t, map[string]any{}, rec.configFor(t, "cron@1.0.0", id), "DON %d must get no offchain config", id)
	}

	check := m.computeOffchainCrossCheck(func() *capabilitiespb.OffchainCapabilitiesRegistry { r, _ := gc.LoadParsed(); return r }(), 1, []registry.DON{
		named(42, "workflow_1_zone-a"), named(43, ""), named(44, "shared"), named(45, "shared"),
	})
	assert.Equal(t, int64(1), check.matchedCaps)
	assert.Equal(t, int64(3), check.divergences[divergenceMissingDON], "unnamed and ambiguous DONs cannot be matched")
	assert.Equal(t, int64(1), check.divergences[divergenceExtraDON], `offchain "shared" matches no unambiguous DON`)
}
