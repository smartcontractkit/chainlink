package queue

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

// MetricPrefix is the prefix of the consensus queue plugin's metrics. Those of its methods (durations, data sizes, the
// KV and blob operations it makes) come from the beholder wrapper, as for any plugin.
const MetricPrefix = "platform_consensus_queue"

// unknownOrg labels the queue depth of workflows that are not registered on this node, so the depth adds up.
const unknownOrg = "unknown"

// The names of the KV values, as the state_bytes metric labels them.
const (
	stateObserved  = "observed_events"
	stateQueued    = "queued_events"
	stateConsensus = "consensus_events"
	stateBlobs     = "blob_handles"
)

// Why events left the consensus queue without being executed, for the dropped events metric.
const (
	dropExpired         = "expired"          // the event's deadline, and the grace period, had passed
	dropExpiredObserved = "expired_observed" // every vote for the event expired before it had enough
	dropPrunedObserved  = "pruned_observed"  // dropped from the observed events to fit the state limits
	dropPrunedQueued    = "pruned_queued"    // dropped from the queued events to fit the state limits
)

// pluginMetrics are the consensus queue plugin's own instruments. They carry the plugin's config digest, and for the
// depth the org of the workflows, but never a workflow.
type pluginMetrics struct {
	configDigest string

	observedEvents    metric.Int64Counter   // events in this node's observations
	observationEvents metric.Int64Histogram // events in each of this node's observations
	admittedEvents    metric.Int64Counter   // events that reached consensus, and joined the queue
	consensusEvents   metric.Int64Counter   // events chosen to be started, from the queue
	droppedEvents     metric.Int64Counter   // events removed without being started, by reason
	queueDepth        metric.Int64Gauge     // events in the queue, by org
	observedDepth     metric.Int64Gauge     // events waiting for enough votes to reach consensus
	stateBytes        metric.Int64Gauge     // bytes in each KV value written by a state transition

	mu        sync.Mutex
	depthOrgs map[string]struct{} // the orgs a depth was reported for, so that it is reported again as zero when it empties
}

func newPluginMetrics(meter metric.Meter, configDigest string) (*pluginMetrics, error) {
	m := &pluginMetrics{configDigest: configDigest, depthOrgs: map[string]struct{}{}}
	var err error
	counter := func(name, description string) (c metric.Int64Counter) {
		if err == nil {
			c, err = meter.Int64Counter(MetricPrefix+"_"+name, metric.WithUnit("{event}"), metric.WithDescription(description))
		}
		return
	}
	gauge := func(name, description, unit string) (g metric.Int64Gauge) {
		if err == nil {
			g, err = meter.Int64Gauge(MetricPrefix+"_"+name, metric.WithUnit(unit), metric.WithDescription(description))
		}
		return
	}
	m.observedEvents = counter("observed_events_total", "Events in this node's observations")
	m.admittedEvents = counter("admitted_events_total", "Events that reached consensus and were added to the queue")
	m.consensusEvents = counter("consensus_events_total", "Events taken from the queue to be started, as decided by consensus")
	m.droppedEvents = counter("dropped_events_total", "Events removed from the consensus queue without being started, by reason")
	m.queueDepth = gauge("queue_depth", "Events in the consensus queue, by org (unknown for workflows not registered on this node)", "{event}")
	m.observedDepth = gauge("observed_depth", "Events observed by too few nodes to have reached consensus yet", "{event}")
	m.stateBytes = gauge("state_bytes", "Bytes in each of the KV values written by the last state transition", "By")
	if err != nil {
		return nil, fmt.Errorf("failed to create consensus queue metrics: %w", err)
	}

	m.observationEvents, err = meter.Int64Histogram(MetricPrefix+"_observation_events", metric.WithUnit("{event}"),
		metric.WithDescription("Events in each of this node's observations"))
	if err != nil {
		return nil, fmt.Errorf("failed to create consensus queue observation events histogram: %w", err)
	}
	return m, nil
}

func (m *pluginMetrics) attrs(extra ...attribute.KeyValue) metric.MeasurementOption {
	return metric.WithAttributes(append([]attribute.KeyValue{attribute.String("configDigest", m.configDigest)}, extra...)...)
}

// recordObservation records the events in an observation this node made.
func (m *pluginMetrics) recordObservation(ctx context.Context, events int) {
	m.observedEvents.Add(ctx, int64(events), m.attrs())
	m.observationEvents.Record(ctx, int64(events), m.attrs())
}

// transition is what a state transition did to the state, and what it left.
type transition struct {
	admitted   int            // events that reached consensus
	consensus  int            // events chosen to be started
	stats      kvStats        // events removed without being started
	observed   int            // events waiting for votes
	queuedByWF map[string]int // events in the queue, by workflow
	stateBytes map[string]int // bytes written, by KV value
}

// recordTransition records a committed state transition. org maps a workflow to its org, if it has one on this node.
func (m *pluginMetrics) recordTransition(ctx context.Context, t transition, org func(workflowID string) (string, bool)) {
	m.admittedEvents.Add(ctx, int64(t.admitted), m.attrs())
	m.consensusEvents.Add(ctx, int64(t.consensus), m.attrs())
	for reason, n := range map[string]int{
		dropExpired:         t.stats.expiredQueued,
		dropExpiredObserved: t.stats.expiredObserved,
		dropPrunedObserved:  t.stats.prunedObserved,
		dropPrunedQueued:    t.stats.prunedQueued,
	} {
		if n > 0 {
			m.droppedEvents.Add(ctx, int64(n), m.attrs(attribute.String("reason", reason)))
		}
	}

	m.observedDepth.Record(ctx, int64(t.observed), m.attrs())
	for key, n := range t.stateBytes {
		m.stateBytes.Record(ctx, int64(n), m.attrs(attribute.String("key", key)))
	}
	m.recordDepth(ctx, t.queuedByWF, org)
}

// recordDepth records the queue depth by org. An org that has no events now is recorded as zero once, as a gauge
// would otherwise keep its last value.
func (m *pluginMetrics) recordDepth(ctx context.Context, queuedByWF map[string]int, org func(workflowID string) (string, bool)) {
	byOrg := make(map[string]int)
	for workflowID, n := range queuedByWF {
		orgID, ok := org(workflowID)
		if !ok || orgID == "" {
			orgID = unknownOrg
		}
		byOrg[orgID] += n
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for orgID := range m.depthOrgs {
		if _, ok := byOrg[orgID]; !ok {
			byOrg[orgID] = 0
			delete(m.depthOrgs, orgID)
		}
	}
	for orgID, n := range byOrg {
		m.queueDepth.Record(ctx, int64(n), m.attrs(attribute.String(platform.KeyOrganizationID, orgID)))
		if n > 0 {
			m.depthOrgs[orgID] = struct{}{}
		}
	}
}

// MetricViews returns the buckets of the plugin's histograms. All of a histogram's buckets must be defined when the
// beholder client is created.
func MetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: MetricPrefix + "_observation_events"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 1, 5, 10, 50, 100, 250, 500, 1000},
			}},
		),
	}
}
