package monitoring

import (
	"context"
	"hash/fnv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	ShardAssignmentModeManualOnly int64 = iota
	ShardAssignmentModeRingOCROnly
	ShardAssignmentModeRingOCRWithOverrides
)

const (
	ShardResolverManual   = "manual"
	ShardResolverRing     = "ring"
	ShardResolverOverride = "override"
)

const (
	ShardResolverOutcomeFound    = "found"
	ShardResolverOutcomeNotFound = "not_found"
	ShardResolverOutcomeError    = "error"
)

// ShardResolverSource* are the bounded values of the `source` label: which
// branch of the manual-assignment precedence decided a resolution. Resolvers
// with no source concept (ring, override) record the empty string, which the
// OTel→Prometheus translation drops from the series, keeping the label at
// these five values.
const (
	ShardResolverSourcePerOwner      = "per_owner"
	ShardResolverSourcePerOrg        = "per_org"
	ShardResolverSourceStaticDefault = "static_default"
	ShardResolverSourceHashed        = "hashed"
	ShardResolverSourceNone          = "none"
)

const (
	ShardAssignmentUpdateApplied    = "applied"
	ShardAssignmentUpdateParseError = "parse_error"
)

type ShardAssignmentMetrics struct {
	modeGauge         metric.Int64Gauge
	resolverRequests  metric.Int64Counter
	resolverDuration  metric.Float64Histogram
	settingsHashGauge metric.Int64Gauge
	configUpdateCount metric.Int64Gauge
	configStaleness   metric.Float64Histogram
	updated           metric.Int64Counter

	mu                 sync.Mutex
	appliedSeen        bool
	lastAppliedHash    uint32
	parseErrSeen       bool
	lastParseErrorHash uint32
	appliedUpdateCount int64
	lastArrivalRawHash uint32
	lastArrivalAt      time.Time
}

var (
	globalShardAssignmentMetricsMu     sync.Mutex
	globalShardAssignmentMetrics       *ShardAssignmentMetrics
	globalShardAssignmentMetricsClient *beholder.Client
)

// GlobalShardAssignmentMetrics binds the shard assignment instruments to the
// meter of the current global beholder client. beholder starts out with a
// noop client and only receives the real one via SetClient during boot, so a
// binding made too early must not be cached for the life of the process
// (every metric would be silently dropped). The binding is therefore retried
// whenever the global client changes, mirroring how the rest of the codebase
// reads beholder.GetMeter() at call time instead of caching it. Failures are
// not cached either: a later call retries the registration.
func GlobalShardAssignmentMetrics() (*ShardAssignmentMetrics, error) {
	client := beholder.GetClient()
	globalShardAssignmentMetricsMu.Lock()
	defer globalShardAssignmentMetricsMu.Unlock()
	if globalShardAssignmentMetrics != nil && globalShardAssignmentMetricsClient == client {
		return globalShardAssignmentMetrics, nil
	}
	m, err := NewShardAssignmentMetricsWithMeter(client.Meter)
	if err != nil {
		return nil, err
	}
	globalShardAssignmentMetrics, globalShardAssignmentMetricsClient = m, client
	return m, nil
}

// NewShardAssignmentMetricsWithMeter registers the shard assignment instruments
// on the given meter instead of the global beholder meter; exposed for tests.
func NewShardAssignmentMetricsWithMeter(meter metric.Meter) (*ShardAssignmentMetrics, error) {
	m := &ShardAssignmentMetrics{}
	var err error
	m.modeGauge, err = meter.Int64Gauge("platform_shard_assignment_mode")
	if err != nil {
		return nil, err
	}
	m.resolverRequests, err = meter.Int64Counter("platform_shard_resolver_requests")
	if err != nil {
		return nil, err
	}
	m.resolverDuration, err = meter.Float64Histogram("platform_shard_resolver_duration_seconds",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 30),
	)
	if err != nil {
		return nil, err
	}
	m.settingsHashGauge, err = meter.Int64Gauge("platform_shard_assignment_settings_hash")
	if err != nil {
		return nil, err
	}
	m.configUpdateCount, err = meter.Int64Gauge("platform_shard_assignment_config_update_count",
		metric.WithDescription("Node-local monotonic count of applied shard assignment settings updates (resets on restart)"))
	if err != nil {
		return nil, err
	}
	m.configStaleness, err = meter.Float64Histogram("platform_shard_assignment_config_staleness_seconds",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.1, 0.5, 1, 2, 5, 10, 30, 60, 120, 300),
	)
	if err != nil {
		return nil, err
	}
	m.updated, err = meter.Int64Counter("platform_shard_assignment_updated")
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *ShardAssignmentMetrics) RecordAssignmentMode(ctx context.Context, mode int64) {
	if m == nil {
		return
	}
	m.modeGauge.Record(ctx, mode)
}

// RecordResolverRequest counts one resolved workflow (or one errored call for
// batches) and records the request duration for single-shot resolver methods.
func (m *ShardAssignmentMetrics) RecordResolverRequest(ctx context.Context, resolver string, outcome string, source string, duration time.Duration) {
	if m == nil {
		return
	}
	m.AddResolverRequest(ctx, resolver, outcome, source)
	m.RecordResolverDuration(ctx, resolver, duration)
}

func (m *ShardAssignmentMetrics) AddResolverRequest(ctx context.Context, resolver string, outcome string, source string) {
	if m == nil {
		return
	}
	m.resolverRequests.Add(ctx, 1, metric.WithAttributes(
		attribute.String("resolver", resolver),
		attribute.String("outcome", outcome),
		attribute.String("source", source),
	))
}

func (m *ShardAssignmentMetrics) RecordResolverDuration(ctx context.Context, resolver string, duration time.Duration) {
	if m == nil {
		return
	}
	m.resolverDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attribute.String("resolver", resolver)))
}

// RecordSettingsArrival records when a settings update landed at this node,
// keyed by the raw blob. The next RecordSettingsApplied observing the same
// raw blob consumes the slot to measure arrival→in-effect staleness;
// single-slot state keeps arrivals bounded by design.
func (m *ShardAssignmentMetrics) RecordSettingsArrival(raw string, at time.Time) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.lastArrivalRawHash = settingsStringHash(raw)
	m.lastArrivalAt = at
	m.mu.Unlock()
}

// RecordSettingsApplied records the hash of a successfully parsed settings
// version and returns true only on the first observation of that hash, so one
// applied update is counted per actual settings change. On a change it also
// bumps the node-local config update count gauge, and when the last observed
// arrival matches the applied raw blob it records the arrival→in-effect
// staleness (superseded or unobserved arrivals produce no sample). A
// successful apply also resets the parse-error dedup so a broken version that
// recurs after a good one is counted and logged again.
func (m *ShardAssignmentMetrics) RecordSettingsApplied(ctx context.Context, settingsHash uint32, raw string) bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	changed := !m.appliedSeen || m.lastAppliedHash != settingsHash
	m.appliedSeen = true
	m.lastAppliedHash = settingsHash
	m.parseErrSeen = false
	m.lastParseErrorHash = 0
	var updateCount int64
	var staleness time.Duration
	arrivalMatched := false
	if changed {
		m.appliedUpdateCount++
		updateCount = m.appliedUpdateCount
		if m.lastArrivalRawHash == settingsStringHash(raw) {
			arrivalMatched = true
			staleness = time.Since(m.lastArrivalAt)
			m.lastArrivalRawHash = 0
		}
	}
	m.mu.Unlock()
	m.settingsHashGauge.Record(ctx, int64(settingsHash))
	if changed {
		m.configUpdateCount.Record(ctx, updateCount)
		m.updated.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", ShardAssignmentUpdateApplied)))
		if arrivalMatched {
			m.configStaleness.Record(ctx, staleness.Seconds())
		}
	}
	return changed
}

// RecordSettingsParseError records a settings version that failed to parse,
// keyed by the raw blob hash so one broken version counts once (until a good
// version is applied again). The settings hash and config version gauges are
// left at the last successfully applied version: a broken version is not
// applied.
func (m *ShardAssignmentMetrics) RecordSettingsParseError(ctx context.Context, raw string) bool {
	if m == nil {
		return false
	}
	rawHash := settingsStringHash(raw)
	m.mu.Lock()
	changed := !m.parseErrSeen || m.lastParseErrorHash != rawHash
	m.parseErrSeen = true
	m.lastParseErrorHash = rawHash
	m.mu.Unlock()
	if changed {
		m.updated.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", ShardAssignmentUpdateParseError)))
	}
	return changed
}

// settingsStringHash hashes a raw settings blob; it identifies the blob for
// the arrival slot and the parse-error dedup, which key on the raw bytes
// rather than the parsed config.
func settingsStringHash(raw string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(raw))
	return h.Sum32()
}
