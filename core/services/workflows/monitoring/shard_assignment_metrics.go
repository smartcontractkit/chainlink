package monitoring

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
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
	updated           metric.Int64Counter

	mu                 sync.Mutex
	appliedSeen        bool
	lastAppliedHash    uint32
	lastAppliedCfg     *cresettings.ShardAssignmentConfig
	parseErrSeen       bool
	lastParseErrorHash uint32
}

// NewShardAssignmentMetricsWithMeter registers the shard assignment instruments
// on meter.
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

// RecordSettingsApplied records the hash of a successfully parsed settings
// version and returns true only on the first observation of that hash, so one
// applied update is counted per actual settings change. The previous config is
// returned as the diff baseline for the caller's applied log; it is kept here
// rather than per-resolver so the baseline always matches the global dedup
// state (a second resolver observing the change logs the real diff, not
// "initial"). A successful apply also resets the parse-error dedup so a broken
// version that recurs after a good one is counted and logged again.
func (m *ShardAssignmentMetrics) RecordSettingsApplied(ctx context.Context, settingsHash uint32, cfg *cresettings.ShardAssignmentConfig) (bool, *cresettings.ShardAssignmentConfig) {
	if m == nil {
		return false, nil
	}
	m.mu.Lock()
	changed := !m.appliedSeen || m.lastAppliedHash != settingsHash
	prev := m.lastAppliedCfg
	m.appliedSeen = true
	m.lastAppliedHash = settingsHash
	m.lastAppliedCfg = cfg
	m.parseErrSeen = false
	m.lastParseErrorHash = 0
	m.mu.Unlock()
	m.settingsHashGauge.Record(ctx, int64(settingsHash))
	if changed {
		m.updated.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", ShardAssignmentUpdateApplied)))
	}
	return changed, prev
}

// RecordSettingsParseError records a settings version that failed to parse,
// keyed by the raw blob hash so one broken version counts once (until a good
// version is applied again). The settings hash gauge is left at the last
// successfully applied version.
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
