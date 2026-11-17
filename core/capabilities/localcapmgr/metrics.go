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
	keyConfigSource   = "config_source"

	configSourceOnchain  = "onchain"
	configSourceOffchain = "offchain"
)

type metrics struct {
	launchesTotal  metric.Int64Counter
	stopsTotal     metric.Int64Counter
	configUpdates  metric.Int64Counter
	launchDuration metric.Float64Histogram
	runningGauge   metric.Int64Gauge

	// The applied version is emitted by the cresettings delegate.
	offchainDivergences metric.Int64Gauge

	gateEnabled      metric.Int64Gauge
	runningBySource  metric.Int64Gauge
	reconcileSkipped metric.Int64Counter
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

	offchainDivergences, err := meter.Int64Gauge("platform_offchain_registry_divergences",
		metric.WithDescription("Allowlisted capabilities whose offchain config is missing or differs from on-chain, by kind"))
	if err != nil {
		return nil, err
	}

	gateEnabled, err := meter.Int64Gauge("platform_offchain_registry_gate_enabled",
		metric.WithDescription("1 when [Capabilities.Local] UseOffchainRegistry is enabled on this node, else 0"))
	if err != nil {
		return nil, err
	}

	runningBySource, err := meter.Int64Gauge("platform_capability_running_by_source",
		metric.WithDescription("Currently running capabilities by spec config source (onchain or offchain)"))
	if err != nil {
		return nil, err
	}

	reconcileSkipped, err := meter.Int64Counter("platform_offchain_registry_reconcile_skipped_total",
		metric.WithDescription("Reconciles skipped because the offchain cutover is enabled but no payload is applied"))
	if err != nil {
		return nil, err
	}

	return &metrics{
		launchesTotal:       launchesTotal,
		stopsTotal:          stopsTotal,
		configUpdates:       configUpdates,
		launchDuration:      launchDuration,
		runningGauge:        runningGauge,
		offchainDivergences: offchainDivergences,
		gateEnabled:         gateEnabled,
		runningBySource:     runningBySource,
		reconcileSkipped:    reconcileSkipped,
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

func (m *metrics) recordGate(ctx context.Context, enabled bool) {
	var v int64
	if enabled {
		v = 1
	}
	m.gateEnabled.Record(ctx, v)
}

func (m *metrics) recordReconcileSkipped(ctx context.Context) {
	m.reconcileSkipped.Add(ctx, 1)
}

// recordRunningBySource records how many running capabilities were launched from each source.
func (m *metrics) recordRunningBySource(ctx context.Context, onchain, offchain int64) {
	m.runningBySource.Record(ctx, onchain, metric.WithAttributes(attribute.String(keyConfigSource, configSourceOnchain)))
	m.runningBySource.Record(ctx, offchain, metric.WithAttributes(attribute.String(keyConfigSource, configSourceOffchain)))
}

// recordOffchainDivergences records the current divergence count per kind.
func (m *localCapabilityManager) recordOffchainDivergences(ctx context.Context, divergences map[string]int64) {
	for kind, count := range divergences {
		m.metrics.offchainDivergences.Record(ctx, count, metric.WithAttributes(attribute.String(keyDivergenceKind, kind)))
	}
}
