package events_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

const (
	testOwner  = "1100000000000000000000000000000000000000"
	testExecID = "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"

	entityFinishedV2 = "workflows.v2." + events.WorkflowExecutionFinished
	entityFinishedV1 = "workflows.v1." + events.WorkflowExecutionFinished
	entityStartedV2  = "workflows.v2." + events.WorkflowExecutionStarted
)

func newTestInjector(t *testing.T, cfg events.FaultInjectorConfig) *events.FaultInjector {
	t.Helper()
	inj, err := events.NewFaultInjector(cfg, logger.Test(t))
	require.NoError(t, err)
	return inj
}

func TestNewFaultInjector_Validation(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	for _, tc := range []struct {
		name string
		cfg  events.FaultInjectorConfig
	}{
		{"rate below zero", events.FaultInjectorConfig{RateBps: -1, Level: 1}},
		{"rate above 10000", events.FaultInjectorConfig{RateBps: 10001, Level: 1}},
		{"invalid level 0", events.FaultInjectorConfig{RateBps: 100, Level: 0}},
		{"invalid level 3", events.FaultInjectorConfig{RateBps: 100, Level: 3}},
		{"blank allowlist entry", events.FaultInjectorConfig{RateBps: 100, Level: 1, OwnerAllowlist: []string{"  "}}},
		{"0x-only allowlist entry", events.FaultInjectorConfig{RateBps: 100, Level: 1, OwnerAllowlist: []string{"0x"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := events.NewFaultInjector(tc.cfg, lggr)
			require.Error(t, err)
		})
	}

	t.Run("nil logger", func(t *testing.T) {
		t.Parallel()
		_, err := events.NewFaultInjector(events.FaultInjectorConfig{RateBps: 100, Level: 1}, nil)
		require.Error(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		_, err := events.NewFaultInjector(events.FaultInjectorConfig{
			Enabled:        true,
			OwnerAllowlist: []string{"0xAABB00000000000000000000000000000000CCDD"},
			RateBps:        100,
			Seed:           "epoch-1",
			Level:          events.DropLevelAllLifecycle,
		}, lggr)
		require.NoError(t, err)
	})
}

func TestFaultInjector_ShouldDrop(t *testing.T) {
	t.Parallel()

	// rate 10000 => every allowlisted execution is selected, so the only
	// gates under test are the fail-closed identity/allowlist/level checks.
	base := events.FaultInjectorConfig{
		Enabled:        true,
		OwnerAllowlist: []string{testOwner},
		RateBps:        10000,
		Seed:           "seed",
		Level:          events.DropLevelFinishedOnly,
	}
	id := events.DropIdentity{WorkflowOwner: testOwner, WorkflowExecutionID: testExecID}

	for _, tc := range []struct {
		name      string
		cfg       func(events.FaultInjectorConfig) events.FaultInjectorConfig
		id        events.DropIdentity
		eventType string
		wantDrop  bool
	}{
		// fail-closed cases
		{
			name:      "disabled never drops",
			cfg:       func(c events.FaultInjectorConfig) events.FaultInjectorConfig { c.Enabled = false; return c },
			id:        id,
			eventType: entityFinishedV2,
		},
		{
			name:      "empty owner never drops",
			id:        events.DropIdentity{WorkflowOwner: "", WorkflowExecutionID: testExecID},
			eventType: entityFinishedV2,
		},
		{
			name:      "empty executionID never drops",
			id:        events.DropIdentity{WorkflowOwner: testOwner, WorkflowExecutionID: ""},
			eventType: entityFinishedV2,
		},
		{
			name:      "owner not allowlisted never drops",
			id:        events.DropIdentity{WorkflowOwner: "9900000000000000000000000000000000000000", WorkflowExecutionID: testExecID},
			eventType: entityFinishedV2,
		},
		{
			name:      "empty allowlist never drops",
			cfg:       func(c events.FaultInjectorConfig) events.FaultInjectorConfig { c.OwnerAllowlist = nil; return c },
			id:        id,
			eventType: entityFinishedV2,
		},
		{
			name:      "rate zero never drops",
			cfg:       func(c events.FaultInjectorConfig) events.FaultInjectorConfig { c.RateBps = 0; return c },
			id:        id,
			eventType: entityFinishedV2,
		},
		// level 1: WorkflowExecutionFinished only, v1 and v2 alike
		{name: "level1 drops finished v2", id: id, eventType: entityFinishedV2, wantDrop: true},
		{name: "level1 drops finished v1", id: id, eventType: entityFinishedV1, wantDrop: true},
		{name: "level1 skips started", id: id, eventType: entityStartedV2},
		{name: "level1 skips capability finished", id: id, eventType: "workflows.v2." + events.CapabilityExecutionFinished},
		// owner normalization: allowlist match is exact after 0x-strip + lowercase
		{
			name:      "owner match is case-insensitive with 0x prefix",
			id:        events.DropIdentity{WorkflowOwner: "0x" + "1100000000000000000000000000000000000000", WorkflowExecutionID: testExecID},
			eventType: entityFinishedV2,
			wantDrop:  true,
		},
		// level 2: all lifecycle events, never metering/deployment/user events
		{
			name: "level2 drops started",
			cfg: func(c events.FaultInjectorConfig) events.FaultInjectorConfig {
				c.Level = events.DropLevelAllLifecycle
				return c
			},
			id:        id,
			eventType: entityStartedV2,
			wantDrop:  true,
		},
		{
			name: "level2 drops trigger started",
			cfg: func(c events.FaultInjectorConfig) events.FaultInjectorConfig {
				c.Level = events.DropLevelAllLifecycle
				return c
			},
			id:        id,
			eventType: "workflows.v2." + events.TriggerExecutionStarted,
			wantDrop:  true,
		},
		{
			name: "level2 skips metering report",
			cfg: func(c events.FaultInjectorConfig) events.FaultInjectorConfig {
				c.Level = events.DropLevelAllLifecycle
				return c
			},
			id:        id,
			eventType: "workflows.v1." + events.MeteringReportEntity,
		},
		{
			name: "level2 skips user logs",
			cfg: func(c events.FaultInjectorConfig) events.FaultInjectorConfig {
				c.Level = events.DropLevelAllLifecycle
				return c
			},
			id:        id,
			eventType: "workflows.v2." + events.WorkflowUserLog,
		},
		{
			name: "level2 skips workflow status changed",
			cfg: func(c events.FaultInjectorConfig) events.FaultInjectorConfig {
				c.Level = events.DropLevelAllLifecycle
				return c
			},
			id:        id,
			eventType: "workflows.v1." + events.WorkflowStatusChanged,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg := base
			if tc.cfg != nil {
				cfg = tc.cfg(base)
			}
			inj := newTestInjector(t, cfg)
			assert.Equal(t, tc.wantDrop, inj.ShouldDrop(t.Context(), tc.id, tc.eventType))
		})
	}

	t.Run("nil injector never drops", func(t *testing.T) {
		t.Parallel()
		var inj *events.FaultInjector
		assert.False(t, inj.ShouldDrop(t.Context(), id, entityFinishedV2))
	})
}

// TestFaultInjector_Determinism pins the property the whole design rests on:
// the drop decision is a pure function of (seed, executionID), so independent
// injector instances (i.e. independent DON nodes) always agree.
func TestFaultInjector_Determinism(t *testing.T) {
	t.Parallel()

	cfg := events.FaultInjectorConfig{
		Enabled:        true,
		OwnerAllowlist: []string{testOwner},
		RateBps:        5000,
		Seed:           "epoch-1",
		Level:          events.DropLevelFinishedOnly,
	}
	nodeA := newTestInjector(t, cfg)
	nodeB := newTestInjector(t, cfg)

	otherSeed := cfg
	otherSeed.Seed = "epoch-2"
	nodeC := newTestInjector(t, otherSeed)

	sameAsOtherSeed := 0
	dropped := 0
	const n = 2000
	for i := range n {
		id := events.DropIdentity{
			WorkflowOwner:       testOwner,
			WorkflowExecutionID: fmt.Sprintf("execution-%d", i),
		}
		a := nodeA.ShouldDrop(t.Context(), id, entityFinishedV2)
		b := nodeB.ShouldDrop(t.Context(), id, entityFinishedV2)
		require.Equal(t, a, b, "independent instances disagreed for %s", id.WorkflowExecutionID)
		// repeated calls on the same instance are stable too
		require.Equal(t, a, nodeA.ShouldDrop(t.Context(), id, entityFinishedV2))
		if a {
			dropped++
		}
		if a == nodeC.ShouldDrop(t.Context(), id, entityFinishedV2) {
			sameAsOtherSeed++
		}
	}
	// at 5000 bps ~half the corpus is selected
	assert.InDelta(t, n/2, dropped, n/10, "drop rate should be ~rateBps/10000")
	// rotating the seed selects a different cohort
	assert.Less(t, sameAsOtherSeed, n, "a different seed must change at least one decision")
}

func TestFaultInjector_RateBoundaries(t *testing.T) {
	t.Parallel()

	mkID := func(i int) events.DropIdentity {
		return events.DropIdentity{
			WorkflowOwner:       testOwner,
			WorkflowExecutionID: fmt.Sprintf("execution-%d", i),
		}
	}
	mk := func(rateBps int) *events.FaultInjector {
		return newTestInjector(t, events.FaultInjectorConfig{
			Enabled:        true,
			OwnerAllowlist: []string{testOwner},
			RateBps:        rateBps,
			Seed:           "seed",
			Level:          events.DropLevelFinishedOnly,
		})
	}

	const n = 2000

	t.Run("rate 0 drops nothing", func(t *testing.T) {
		t.Parallel()
		inj := mk(0)
		for i := range n {
			require.False(t, inj.ShouldDrop(t.Context(), mkID(i), entityFinishedV2))
		}
	})

	t.Run("rate 10000 drops everything allowlisted", func(t *testing.T) {
		t.Parallel()
		inj := mk(10000)
		for i := range n {
			require.True(t, inj.ShouldDrop(t.Context(), mkID(i), entityFinishedV2))
		}
	})

	t.Run("drop set grows monotonically with rate", func(t *testing.T) {
		t.Parallel()
		low, high := mk(100), mk(1000)
		lowDrops, highDrops := 0, 0
		for i := range n {
			l := low.ShouldDrop(t.Context(), mkID(i), entityFinishedV2)
			h := high.ShouldDrop(t.Context(), mkID(i), entityFinishedV2)
			if l {
				lowDrops++
				require.True(t, h, "an execution dropped at 100 bps must also drop at 1000 bps")
			}
			if h {
				highDrops++
			}
		}
		assert.Positive(t, highDrops)
		assert.Less(t, lowDrops, highDrops)
	})
}

// TestFaultInjector_InjectionRecordContract pins the injection record: an
// external detector parses this exact message and these exact field keys to
// reconcile injected-vs-detected. Renaming any of them must break CI.
func TestFaultInjector_InjectionRecordContract(t *testing.T) {
	t.Parallel()

	lggr, observed := logger.TestObserved(t, zapcore.InfoLevel)
	inj, err := events.NewFaultInjector(events.FaultInjectorConfig{
		Enabled:        true,
		OwnerAllowlist: []string{testOwner},
		RateBps:        10000,
		Seed:           "seed",
		Level:          events.DropLevelFinishedOnly,
	}, lggr)
	require.NoError(t, err)

	id := events.DropIdentity{WorkflowOwner: testOwner, WorkflowExecutionID: testExecID}
	require.True(t, inj.ShouldDrop(t.Context(), id, entityFinishedV2))

	entries := observed.FilterMessage("Fault injection: dropped workflow event").All()
	require.Len(t, entries, 1, "expected exactly one injection record")

	fields := entries[0].ContextMap()
	assert.Equal(t, testExecID, fields["executionID"])
	assert.Equal(t, entityFinishedV2, fields["eventType"])
	assert.Equal(t, testOwner, fields["workflowOwner"])

	// no record when nothing is dropped
	require.False(t, inj.ShouldDrop(t.Context(), events.DropIdentity{}, entityFinishedV2))
	require.Len(t, observed.FilterMessage("Fault injection: dropped workflow event").All(), 1)
}

// TestEmit_FaultInjectionSeam verifies the seam itself: an installed decider
// suppresses the event on BOTH transports (beholder observer sees nothing),
// non-matching events still flow, and removing the decider restores default
// behavior.
func TestEmit_FaultInjectionSeam(t *testing.T) { //nolint:paralleltest // mutates the package-global drop decider and beholder client
	beholderObserver := beholdertest.NewObserver(t)
	lggr, observed := logger.TestObserved(t, zapcore.InfoLevel)

	inj, err := events.NewFaultInjector(events.FaultInjectorConfig{
		Enabled:        true,
		OwnerAllowlist: []string{testOwner},
		RateBps:        10000,
		Seed:           "seed",
		Level:          events.DropLevelFinishedOnly,
	}, lggr)
	require.NoError(t, err)
	events.SetDropDecider(inj)
	t.Cleanup(func() { events.SetDropDecider(nil) })

	labels := map[string]string{platform.KeyWorkflowOwner: testOwner}

	// Finished events (v1 + v2) are dropped before any transport.
	require.NoError(t, events.EmitExecutionFinishedEvent(t.Context(), labels, "completed", testExecID, nil, events.ErrorClassificationUnspecified, nil))
	assert.Empty(t, beholderObserver.Messages(t, "beholder_entity", entityFinishedV1))
	assert.Empty(t, beholderObserver.Messages(t, "beholder_entity", entityFinishedV2))

	// One injection record per dropped event (v1 + v2 => two records).
	records := observed.FilterMessage("Fault injection: dropped workflow event").All()
	require.Len(t, records, 2)
	gotEventTypes := []string{records[0].ContextMap()["eventType"].(string), records[1].ContextMap()["eventType"].(string)}
	assert.ElementsMatch(t, []string{entityFinishedV1, entityFinishedV2}, gotEventTypes)

	// Started events do not match level 1 and still flow.
	require.NoError(t, events.EmitExecutionStartedEvent(t.Context(), labels, "trigger-1", testExecID))
	assert.Len(t, beholderObserver.Messages(t, "beholder_entity", entityStartedV2), 1)

	// Non-allowlisted owners are untouched.
	otherLabels := map[string]string{platform.KeyWorkflowOwner: "9900000000000000000000000000000000000000"}
	require.NoError(t, events.EmitExecutionFinishedEvent(t.Context(), otherLabels, "completed", testExecID, nil, events.ErrorClassificationUnspecified, nil))
	assert.Len(t, beholderObserver.Messages(t, "beholder_entity", entityFinishedV2), 1)

	// Removing the decider restores passthrough (nil-safe default).
	events.SetDropDecider(nil)
	require.NoError(t, events.EmitExecutionFinishedEvent(t.Context(), labels, "completed", testExecID, nil, events.ErrorClassificationUnspecified, nil))
	assert.Len(t, beholderObserver.Messages(t, "beholder_entity", entityFinishedV2), 2)
}
