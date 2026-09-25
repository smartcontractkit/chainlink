package localcapmgr

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	keyCapabilityID   = "capability_id"
	keyDivergenceKind = "kind"
)

type metrics struct {
	launchesTotal  metric.Int64Counter
	stopsTotal     metric.Int64Counter
	configUpdates  metric.Int64Counter
	launchDuration metric.Float64Histogram
	runningGauge   metric.Int64Gauge

	// Offchain capabilities registry (Phase 2 parallel-run telemetry).
	offchainVersion     metric.Int64Gauge
	offchainMatchedCaps metric.Int64Gauge
	offchainDivergences metric.Int64Counter
}

func newMetrics() (*metrics, error) {
	meter := beholder.GetMeter()

	launchesTotal, err := meter.Int64Counter("platform_capability_launches_total",
		metric.WithDescription("Total capabilities started via registry"))
	if err != nil {
		return nil, err
	}

	stopsTotal, err := meter.Int64Counter("platform_capability_stops_total",
		metric.WithDescription("Total capabilities stopped via registry"))
	if err != nil {
		return nil, err
	}

	configUpdates, err := meter.Int64Counter("platform_capability_config_updates_total",
		metric.WithDescription("Total config updates received"))
	if err != nil {
		return nil, err
	}

	launchDuration, err := meter.Float64Histogram("platform_capability_launch_duration_seconds",
		metric.WithDescription("Time to start a capability"))
	if err != nil {
		return nil, err
	}

	runningGauge, err := meter.Int64Gauge("platform_capability_running",
		metric.WithDescription("Currently running capabilities"))
	if err != nil {
		return nil, err
	}

	offchainVersion, err := meter.Int64Gauge("platform_offchain_registry_version",
		metric.WithDescription("Applied offchain capabilities registry version (0 if none applied)"))
	if err != nil {
		return nil, err
	}

	offchainMatchedCaps, err := meter.Int64Gauge("platform_offchain_registry_matched_capabilities",
		metric.WithDescription("Allowlisted capabilities present in both the offchain and on-chain registries"))
	if err != nil {
		return nil, err
	}

	offchainDivergences, err := meter.Int64Counter("platform_offchain_registry_divergences_total",
		metric.WithDescription("Cross-validation divergences between the offchain and on-chain registries, by kind"))
	if err != nil {
		return nil, err
	}

	return &metrics{
		launchesTotal:       launchesTotal,
		stopsTotal:          stopsTotal,
		configUpdates:       configUpdates,
		launchDuration:      launchDuration,
		runningGauge:        runningGauge,
		offchainVersion:     offchainVersion,
		offchainMatchedCaps: offchainMatchedCaps,
		offchainDivergences: offchainDivergences,
	}, nil
}

func (m *metrics) recordLaunch(ctx context.Context, capID string, duration time.Duration) {
	attrs := metric.WithAttributes(attribute.String(keyCapabilityID, capID))
	m.launchesTotal.Add(ctx, 1, attrs)
	m.launchDuration.Record(ctx, duration.Seconds(), attrs)
}

func (m *metrics) recordStop(ctx context.Context, capID string) {
	attrs := metric.WithAttributes(attribute.String(keyCapabilityID, capID))
	m.stopsTotal.Add(ctx, 1, attrs)
}

func (m *metrics) recordConfigUpdate(ctx context.Context, capID string) {
	attrs := metric.WithAttributes(attribute.String(keyCapabilityID, capID))
	m.configUpdates.Add(ctx, 1, attrs)
}

func (m *metrics) recordRunning(ctx context.Context, count int64) {
	m.runningGauge.Record(ctx, count)
}

// recordOffchainCheck emits the result of an offchain-vs-onchain cross-validation pass.
func (m *localCapabilityManager) recordOffchainCheck(ctx context.Context, c offchainCrossCheck) {
	m.metrics.offchainVersion.Record(ctx, int64(c.version))
	if c.offchainEmpty {
		m.metrics.offchainMatchedCaps.Record(ctx, 0)
		return
	}
	m.metrics.offchainMatchedCaps.Record(ctx, c.matchedCaps)
	for kind, count := range c.divergences {
		if count == 0 {
			continue
		}
		m.metrics.offchainDivergences.Add(ctx, count, metric.WithAttributes(attribute.String(keyDivergenceKind, kind)))
	}
}
