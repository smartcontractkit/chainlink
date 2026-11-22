package queue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// readerMetrics returns metrics recorded to a reader, which reads back what was recorded.
func readerMetrics(t *testing.T) (*pluginMetrics, *sdkmetric.ManualReader) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })
	m, err := newPluginMetrics(mp.Meter("test"), "digest-1")
	require.NoError(t, err)
	return m, reader
}

func matches(set attribute.Set, want map[string]string) bool {
	for k, v := range want {
		if got, ok := set.Value(attribute.Key(k)); !ok || got.AsString() != v {
			return false
		}
	}
	return true
}

func collectMetric(t *testing.T, reader *sdkmetric.ManualReader, name string) metricdata.Aggregation {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == MetricPrefix+"_"+name {
				return m.Data
			}
		}
	}
	return nil
}

// total is the sum of a counter's series that have the attributes.
func total(t *testing.T, reader *sdkmetric.ManualReader, name string, attrs map[string]string) (n int64) {
	t.Helper()
	data := collectMetric(t, reader, name)
	if data == nil {
		return 0
	}
	sum, ok := data.(metricdata.Sum[int64])
	require.True(t, ok, "%s is %T", name, data)
	for _, dp := range sum.DataPoints {
		if matches(dp.Attributes, attrs) {
			n += dp.Value
		}
	}
	return n
}

// last is the value of the gauge series that has the attributes, and whether there is one.
func last(t *testing.T, reader *sdkmetric.ManualReader, name string, attrs map[string]string) (int64, bool) {
	t.Helper()
	data := collectMetric(t, reader, name)
	if data == nil {
		return 0, false
	}
	g, ok := data.(metricdata.Gauge[int64])
	require.True(t, ok, "%s is %T", name, data)
	for _, dp := range g.DataPoints {
		if matches(dp.Attributes, attrs) {
			return dp.Value, true
		}
	}
	return 0, false
}

func orgsOf(orgs map[string]string) func(string) (string, bool) {
	return func(workflowID string) (string, bool) { org, ok := orgs[workflowID]; return org, ok }
}

func TestPluginMetrics_Observation(t *testing.T) {
	t.Parallel()
	m, reader := readerMetrics(t)

	m.recordObservation(t.Context(), 7)
	m.recordObservation(t.Context(), 0)
	m.recordObservation(t.Context(), 3)

	assert.Equal(t, int64(10), total(t, reader, "observed_events_total", map[string]string{"configDigest": "digest-1"}))
	hist, ok := collectMetric(t, reader, "observation_events").(metricdata.Histogram[int64])
	require.True(t, ok)
	require.Len(t, hist.DataPoints, 1)
	assert.Equal(t, uint64(3), hist.DataPoints[0].Count, "one per observation, including the empty one")
	assert.Equal(t, int64(10), hist.DataPoints[0].Sum)
}

func TestPluginMetrics_Transition(t *testing.T) {
	t.Parallel()
	m, reader := readerMetrics(t)
	ctx := t.Context()
	orgs := orgsOf(map[string]string{"wf-a": "org-1", "wf-b": "org-1", "wf-c": "org-2"})

	m.recordTransition(ctx, transition{
		admitted:  3,
		consensus: 2,
		stats:     kvStats{expiredQueued: 1, expiredObserved: 2, prunedObserved: 3, prunedQueued: 4},
		observed:  5,
		queuedByWF: map[string]int{
			"wf-a": 2, "wf-b": 3, // org-1
			"wf-c":      7, // org-2
			"wf-gone":   4, // not registered on this node
			"wf-no-org": 1, // registered, but with no org: counts as unknown too
		},
		stateBytes: map[string]int{stateObserved: 100, stateQueued: 200, stateConsensus: 300, stateBlobs: 400},
	}, func(wf string) (string, bool) {
		if wf == "wf-no-org" {
			return "", true
		}
		return orgs(wf)
	})

	assert.Equal(t, int64(3), total(t, reader, "admitted_events_total", nil))
	assert.Equal(t, int64(2), total(t, reader, "consensus_events_total", nil))
	for reason, want := range map[string]int64{dropExpired: 1, dropExpiredObserved: 2, dropPrunedObserved: 3, dropPrunedQueued: 4} {
		assert.Equal(t, want, total(t, reader, "dropped_events_total", map[string]string{"reason": reason}), reason)
	}

	depth := func(org string) int64 {
		v, ok := last(t, reader, "queue_depth", map[string]string{platform.KeyOrganizationID: org})
		require.True(t, ok, "depth for %s", org)
		return v
	}
	assert.Equal(t, int64(5), depth("org-1"))
	assert.Equal(t, int64(7), depth("org-2"))
	assert.Equal(t, int64(5), depth(unknownOrg), "workflows without an org here are together")
	assert.Equal(t, int64(5+7+5), depth("org-1")+depth("org-2")+depth(unknownOrg), "so that the depths add up to the queue")

	observed, ok := last(t, reader, "observed_depth", nil)
	require.True(t, ok)
	assert.Equal(t, int64(5), observed)
	for key, want := range map[string]int64{stateObserved: 100, stateQueued: 200, stateConsensus: 300, stateBlobs: 400} {
		got, ok := last(t, reader, "state_bytes", map[string]string{"key": key})
		require.True(t, ok, key)
		assert.Equal(t, want, got, key)
	}
}

func TestPluginMetrics_NothingDroppedIsNotRecorded(t *testing.T) {
	t.Parallel()
	m, reader := readerMetrics(t)
	m.recordTransition(t.Context(), transition{admitted: 1}, orgsOf(nil))
	assert.Nil(t, collectMetric(t, reader, "dropped_events_total"), "no series for reasons that did not happen")
}

// A gauge keeps its last value, so an org whose queue empties must be recorded as zero once.
func TestPluginMetrics_DepthOfAnOrgThatEmpties(t *testing.T) {
	t.Parallel()
	m, reader := readerMetrics(t)
	orgs := orgsOf(map[string]string{"wf-a": "org-1", "wf-b": "org-2"})
	depth := func(org string) int64 {
		v, ok := last(t, reader, "queue_depth", map[string]string{platform.KeyOrganizationID: org})
		require.True(t, ok, org)
		return v
	}

	m.recordTransition(t.Context(), transition{queuedByWF: map[string]int{"wf-a": 3, "wf-b": 4}}, orgs)
	assert.Equal(t, int64(3), depth("org-1"))
	assert.Equal(t, int64(4), depth("org-2"))

	m.recordTransition(t.Context(), transition{queuedByWF: map[string]int{"wf-a": 1}}, orgs)
	assert.Equal(t, int64(1), depth("org-1"))
	assert.Equal(t, int64(0), depth("org-2"), "emptied, so recorded as zero rather than left at 4")

	m.recordTransition(t.Context(), transition{}, orgs)
	assert.Equal(t, int64(0), depth("org-1"))
	assert.Equal(t, int64(0), depth("org-2"))

	m.mu.Lock()
	assert.Empty(t, m.depthOrgs, "orgs at zero are not tracked, so the set is bounded by the orgs with events")
	m.mu.Unlock()
}

func TestMetricViews(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, MetricViews())
}

// orgCoordinator knows the org of some workflows.
type orgCoordinator struct {
	noopCoordinator
	orgs map[string]string
}

func (o orgCoordinator) WorkflowOrg(workflowID string) (string, bool) {
	org, ok := o.orgs[workflowID]
	return org, ok
}

var _ triggers.Coordinator = orgCoordinator{}

// The metrics follow a plugin through rounds: the events it observes, those that reach consensus and wait in the queue
// (by org), and those it then takes from the queue to start.
func TestPluginMetrics_ThroughRounds(t *testing.T) {
	t.Parallel()

	c := newSTPlugin(t)
	m, reader := readerMetrics(t)
	c.metrics = m
	c.triggerCoordinator = orgCoordinator{orgs: map[string]string{testWF: "org-1"}}
	kv := memKV{}
	events := mkEvents(t, 3, 2000)
	allNodes := func(o nodeOpts) []ocrtypes.AttributedObservation {
		return []ocrtypes.AttributedObservation{ao(t, 0, nodeObs(t, o, events...)), ao(t, 1, nodeObs(t, o, events...)), ao(t, 2, nodeObs(t, o, events...))}
	}

	// no capacity: three events reach consensus, and wait in the queue
	noCapacity := nodeOpts{capacity: 0, wfCaps: map[string]uint32{testWF: 0}, expiration: 100}
	require.NoError(t, stateTransition(t, c, kv, allNodes(noCapacity)...))
	assert.Equal(t, int64(3), total(t, reader, "admitted_events_total", nil))
	assert.Zero(t, total(t, reader, "consensus_events_total", nil))
	depth, ok := last(t, reader, "queue_depth", map[string]string{platform.KeyOrganizationID: "org-1"})
	require.True(t, ok)
	assert.Equal(t, int64(3), depth, "the queue's depth, for the org of the workflow")
	for key, size := range map[string]int{stateQueued: len(kv[queuedEventsKey]), stateObserved: len(kv[observedEventsKey]), stateBlobs: len(kv[blobHandlesKey])} {
		got, _ := last(t, reader, "state_bytes", map[string]string{"key": key})
		assert.Equal(t, int64(size), got, "%s is what was written", key)
	}

	// capacity for two: two are taken from the queue to be started
	someCapacity := nodeOpts{capacity: 2, wfCaps: map[string]uint32{testWF: 2}, expiration: 100}
	require.NoError(t, stateTransition(t, c, kv, allNodes(someCapacity)...))
	assert.Equal(t, int64(3), total(t, reader, "admitted_events_total", nil), "nothing new was admitted")
	assert.Equal(t, int64(2), total(t, reader, "consensus_events_total", nil))
	depth, _ = last(t, reader, "queue_depth", map[string]string{platform.KeyOrganizationID: "org-1"})
	assert.Equal(t, int64(1), depth)

	// expired: the one left is dropped, not started. (No node observes anything, as the two that were started would
	// be admitted again, which a pass through the engine's duplicate check is what prevents from being run.)
	expired := nodeOpts{capacity: 5, wfCaps: map[string]uint32{testWF: 5}, expiration: 3000}
	require.NoError(t, stateTransition(t, c, kv,
		ao(t, 0, nodeObs(t, expired)), ao(t, 1, nodeObs(t, expired)), ao(t, 2, nodeObs(t, expired))))
	assert.Equal(t, int64(2), total(t, reader, "consensus_events_total", nil))
	assert.Equal(t, int64(1), total(t, reader, "dropped_events_total", map[string]string{"reason": dropExpired}))
	depth, _ = last(t, reader, "queue_depth", map[string]string{platform.KeyOrganizationID: "org-1"})
	assert.Zero(t, depth, "the queue is empty, which is recorded as zero")
}

// What an observation reports is counted when it is made.
func TestPluginMetrics_ObservationsAreCounted(t *testing.T) {
	t.Parallel()

	events := mkEvents(t, 4, int64(time.Now().Add(time.Hour).Unix()))
	blobs := &memBlobs{}
	c := testPlugin(t, sliceQueue{events: events}, blobs, 1<<19)
	m, reader := readerMetrics(t)
	c.metrics = m

	obs, err := observe(t, c, memKV{}, blobs)
	require.NoError(t, err)
	require.Len(t, obs.Events, 4)
	assert.Equal(t, int64(4), total(t, reader, "observed_events_total", nil))
}
