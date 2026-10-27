package monitoring

import (
	"context"
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

const (
	ShardResolverSourcePerOwner      = "per_owner"
	ShardResolverSourcePerOrg        = "per_org"
	ShardResolverSourceStaticDefault = "static_default"
	ShardResolverSourceHashed        = "hashed"
)

const (
	ShardAssignmentUpdateApplied    = "applied"
	ShardAssignmentUpdateParseError = "parse_error"
)

type ShardAssignmentMetrics struct {
	modeGauge         metric.Int64Gauge
	resolverRequests  metric.Int64Counter
	resolverDuration  metric.Int64Histogram
	settingsHashGauge metric.Int64Gauge
	updated           metric.Int64Counter

	mu                 sync.Mutex
	appliedSeen        bool
	lastAppliedHash    uint32
	parseErrSeen       bool
	lastParseErrorHash uint32
}

var (
	globalShardAssignmentMetrics     *ShardAssignmentMetrics
	errGlobalShardAssignmentMetrics  error
	globalShardAssignmentMetricsOnce sync.Once
)

func GlobalShardAssignmentMetrics() (*ShardAssignmentMetrics, error) {
	globalShardAssignmentMetricsOnce.Do(func() {
		globalShardAssignmentMetrics, errGlobalShardAssignmentMetrics = NewShardAssignmentMetricsWithMeter(beholder.GetMeter())
	})
	return globalShardAssignmentMetrics, errGlobalShardAssignmentMetrics
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
	m.resolverDuration, err = meter.Int64Histogram("platform_shard_resolver_duration_seconds",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(1, 2, 5, 10, 30, 60, 120, 300),
	)
	if err != nil {
		return nil, err
	}
	m.settingsHashGauge, err = meter.Int64Gauge("platform_shard_assignment_settings_hash")
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
	m.resolverDuration.Record(ctx, int64(duration/time.Second), metric.WithAttributes(attribute.String("resolver", resolver)))
}

// RecordSettingsApplied records the hash of a successfully parsed settings
// version and returns true only on the first observation of that hash, so one
// applied update is counted per actual settings change.
func (m *ShardAssignmentMetrics) RecordSettingsApplied(ctx context.Context, settingsHash uint32) bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	changed := !m.appliedSeen || m.lastAppliedHash != settingsHash
	m.appliedSeen = true
	m.lastAppliedHash = settingsHash
	m.mu.Unlock()
	m.settingsHashGauge.Record(ctx, int64(settingsHash))
	if changed {
		m.updated.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", ShardAssignmentUpdateApplied)))
	}
	return changed
}

// RecordSettingsParseError records a settings version that failed to parse,
// keyed by the raw blob hash so one broken version counts once. The settings
// hash gauge is left untouched: the previously applied assignment stays live.
func (m *ShardAssignmentMetrics) RecordSettingsParseError(ctx context.Context, rawHash uint32) bool {
	if m == nil {
		return false
	}
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
