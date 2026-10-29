package v2

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// engineType label values for the running-engines gauge.
const (
	engineTypeCoordinated = "coordinated"
	engineTypeLegacy      = "legacy"
)

type metrics struct {
	handleDuration    metric.Int64Histogram
	fetchedWorkflows  metric.Int64Gauge
	runningWorkflows  metric.Int64Gauge
	drainingWorkflows metric.Int64Gauge
	// runningEngines counts the engines currently in the registry, split by
	// engineType ("coordinated" or "legacy"). Both series are always reported,
	// including the zeroes, so an absent series means the node stopped
	// reporting rather than "no engines of that type".
	runningEngines metric.Int64Gauge
	completedSyncs metric.Int64Counter
	drainStarted   metric.Int64Counter
	drainCompleted metric.Int64Counter
	drainDuration  metric.Int64Histogram
	deleteDeferred metric.Int64Counter

	// Per-source metrics for multi-source observability
	sourceHealth        metric.Int64Gauge     // 1=healthy, 0=unhealthy per source
	workflowsPerSource  metric.Int64Gauge     // workflows fetched per source
	sourceFetchDuration metric.Int64Histogram // fetch latency per source
	sourceFetchErrors   metric.Int64Counter   // error count per source

	// Per-tick reconciliation metrics
	reconcileEventsDispatched metric.Int64Histogram // events dispatched per source per tick
	reconcileDuration         metric.Int64Histogram // wall-clock ms for parallel event processing
	reconcileEventsBackoff    metric.Int64Counter   // events skipped due to backoff
	activationDropped         metric.Int64Counter   // activations dropped after max load/init retries
	activationAbandoned       metric.Int64Counter   // activations abandoned with reason label

	// On-disk WASM cache write (sync); duration tails indicate IO contention vs typical skew.
	moduleStoreDuration metric.Int64Histogram

	// Backfill of workflow_specs_v2.workflow_tag; "reason" is
	// "filled_empty" (row had an empty tag) or "refreshed" (row had a stale tag).
	workflowTagBackfill metric.Int64Counter
}

func (m *metrics) recordHandleDuration(ctx context.Context, d time.Duration, event string, success bool) {
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("eventType", event),
	))
}

func (m *metrics) recordFetchedWorkflows(ctx context.Context, count int) {
	m.fetchedWorkflows.Record(ctx, int64(count))
}

func (m *metrics) recordRunningWorkflows(ctx context.Context, count int) {
	m.runningWorkflows.Record(ctx, int64(count))
}

func (m *metrics) recordDrainingWorkflows(ctx context.Context, count int) {
	m.drainingWorkflows.Record(ctx, int64(count))
}

// recordRunningEngines reports the engine-type split of the registry. Called
// once per reconciliation tick, with both counts, so the two series stay in
// lockstep and "coordinated == 0" is an assertion the alert can rely on rather
// than the absence of a sample.
func (m *metrics) recordRunningEngines(ctx context.Context, coordinated, legacy int) {
	m.runningEngines.Record(ctx, int64(coordinated), metric.WithAttributes(
		attribute.String("engineType", engineTypeCoordinated),
	))
	m.runningEngines.Record(ctx, int64(legacy), metric.WithAttributes(
		attribute.String("engineType", engineTypeLegacy),
	))
}

func (m *metrics) incrementCompletedSyncs(ctx context.Context) {
	m.completedSyncs.Add(ctx, 1)
}

func (m *metrics) incrementDrainStarted(ctx context.Context) {
	m.drainStarted.Add(ctx, 1)
}

func (m *metrics) recordDrainCompleted(ctx context.Context, duration time.Duration) {
	m.drainCompleted.Add(ctx, 1)
	m.drainDuration.Record(ctx, duration.Milliseconds())
}

func (m *metrics) incrementDeleteDeferred(ctx context.Context, reason string) {
	m.deleteDeferred.Add(ctx, 1, metric.WithAttributes(
		attribute.String("reason", reason),
	))
}

// recordSourceFetch records metrics for a source fetch operation.
func (m *metrics) recordSourceFetch(ctx context.Context, sourceName string, workflowCount int, duration time.Duration, err error) {
	attrs := metric.WithAttributes(attribute.String("source", sourceName))

	// Record fetch duration
	m.sourceFetchDuration.Record(ctx, duration.Milliseconds(), attrs)

	// Record workflow count per source
	m.workflowsPerSource.Record(ctx, int64(workflowCount), attrs)

	// Record health status (1=healthy, 0=unhealthy)
	if err != nil {
		m.sourceHealth.Record(ctx, 0, attrs)
		m.sourceFetchErrors.Add(ctx, 1, attrs)
	} else {
		m.sourceHealth.Record(ctx, 1, attrs)
	}
}

func (m *metrics) recordReconcileBatch(ctx context.Context, source string, dispatched int, duration time.Duration) {
	attrs := metric.WithAttributes(attribute.String("source", source))
	m.reconcileEventsDispatched.Record(ctx, int64(dispatched), attrs)
	m.reconcileDuration.Record(ctx, duration.Milliseconds(), attrs)
}

func (m *metrics) recordReconcileBackoff(ctx context.Context, source string, count int) {
	m.reconcileEventsBackoff.Add(ctx, int64(count), metric.WithAttributes(
		attribute.String("source", source),
	))
}

func (m *metrics) incrementActivationDropped(ctx context.Context, source, reason string) {
	m.activationDropped.Add(ctx, 1, metric.WithAttributes(
		attribute.String("source", source),
		attribute.String("reason", reason),
	))
}

func (m *metrics) incrementActivationAbandoned(ctx context.Context, source, reason string) {
	m.activationAbandoned.Add(ctx, 1, metric.WithAttributes(
		attribute.String("source", source),
		attribute.String("reason", reason),
	))
}

func (m *metrics) recordModuleStore(ctx context.Context, d time.Duration, success bool) {
	if m == nil {
		return
	}
	m.moduleStoreDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", strconv.FormatBool(success)),
	))
}

func (m *metrics) incrementWorkflowTagBackfill(ctx context.Context, reason string) {
	if m == nil {
		return
	}
	m.workflowTagBackfill.Add(ctx, 1, metric.WithAttributes(
		attribute.String("reason", reason),
	))
}

const (
	GaugeWorkflowModuleCacheDiskUsageBytes     = "platform_workflow_module_cache_disk_usage_bytes"
	WorkflowModuleCacheDiskMonitorTickInterval = time.Minute
)

type CacheMetrics struct {
	reloadSource        metric.Int64Counter // attribute "source": "weak_ref" | "disk"
	evictionTotal       metric.Int64Counter
	loadedGauge         metric.Int64Gauge
	memorySaved         metric.Int64Gauge   // bytes saved by evicting idle modules
	versionMismatch     metric.Int64Counter // cached binary rejected due to engine version mismatch
	pinExhausted        metric.Int64Counter // execute retries exhausted before a pin succeeds
	tryAcquireExhausted metric.Int64Counter // moduleEntry CAS attempts exhausted while pinning
}

// cachePinExhaustedHook is a test-only hook to observe pin-exhausted recordings.
var cachePinExhaustedHook func()

// cacheTryAcquireExhaustedHook is a test-only hook to observe tryAcquire-exhausted recordings.
var cacheTryAcquireExhaustedHook func()

func (cm *CacheMetrics) recordReload(ctx context.Context, source string) {
	if cm == nil {
		return
	}
	cm.reloadSource.Add(ctx, 1, metric.WithAttributes(attribute.String("source", source)))
	promModuleCacheReloadTotal.WithLabelValues(source).Inc()
}

func (cm *CacheMetrics) recordEviction(ctx context.Context, count int) {
	if cm == nil {
		return
	}
	cm.evictionTotal.Add(ctx, int64(count))
	promModuleCacheEvictionTotal.Add(float64(count))
}

func (cm *CacheMetrics) recordLoaded(ctx context.Context, count int) {
	if cm == nil {
		return
	}
	cm.loadedGauge.Record(ctx, int64(count))
	promModuleCacheLoaded.Set(float64(count))
}

func (cm *CacheMetrics) recordMemorySaved(ctx context.Context, bytes int64) {
	if cm == nil {
		return
	}
	cm.memorySaved.Record(ctx, bytes)
	promModuleCacheMemorySavedBytes.Set(float64(bytes))
}

func (cm *CacheMetrics) recordVersionMismatch(ctx context.Context) {
	if cm == nil {
		return
	}
	cm.versionMismatch.Add(ctx, 1)
	promModuleCacheVersionMismatchTotal.Inc()
}

func (cm *CacheMetrics) recordPinExhausted(ctx context.Context) {
	if cm == nil {
		return
	}
	cm.pinExhausted.Add(ctx, 1)
	promModuleCachePinExhaustedTotal.Inc()
	if cachePinExhaustedHook != nil {
		cachePinExhaustedHook()
	}
}

func (cm *CacheMetrics) recordTryAcquireExhausted(ctx context.Context) {
	if cm == nil {
		return
	}
	cm.tryAcquireExhausted.Add(ctx, 1)
	promModuleCacheTryAcquireExhaustedTotal.Inc()
	if cacheTryAcquireExhaustedHook != nil {
		cacheTryAcquireExhaustedHook()
	}
}

// NewCacheMetrics registers the module cache instruments on meter.
func NewCacheMetrics(meter metric.Meter) (*CacheMetrics, error) {
	reloadSource, err := meter.Int64Counter("platform_workflow_module_cache_reload_total")
	if err != nil {
		return nil, err
	}
	evictionTotal, err := meter.Int64Counter("platform_workflow_module_cache_eviction_total")
	if err != nil {
		return nil, err
	}
	loadedGauge, err := meter.Int64Gauge("platform_workflow_module_cache_loaded")
	if err != nil {
		return nil, err
	}
	memorySaved, err := meter.Int64Gauge("platform_workflow_module_cache_memory_saved_bytes")
	if err != nil {
		return nil, err
	}
	versionMismatch, err := meter.Int64Counter("platform_workflow_module_cache_version_mismatch_total")
	if err != nil {
		return nil, err
	}
	pinExhausted, err := meter.Int64Counter("platform_workflow_module_cache_pin_exhausted_total")
	if err != nil {
		return nil, err
	}
	tryAcquireExhausted, err := meter.Int64Counter("platform_workflow_module_cache_try_acquire_exhausted_total")
	if err != nil {
		return nil, err
	}
	return &CacheMetrics{
		reloadSource:        reloadSource,
		evictionTotal:       evictionTotal,
		loadedGauge:         loadedGauge,
		memorySaved:         memorySaved,
		versionMismatch:     versionMismatch,
		pinExhausted:        pinExhausted,
		tryAcquireExhausted: tryAcquireExhausted,
	}, nil
}

// newMetrics registers the syncer's instruments on meter.
func newMetrics(meter metric.Meter) (*metrics, error) {
	handleDuration, err := meter.Int64Histogram("platform_workflow_registry_syncer_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	fetchedWorkflows, err := meter.Int64Gauge("platform_workflow_registry_syncer_fetched_workflows")
	if err != nil {
		return nil, err
	}

	runningWorkflows, err := meter.Int64Gauge("platform_workflow_registry_syncer_running_workflows")
	if err != nil {
		return nil, err
	}

	drainingWorkflows, err := meter.Int64Gauge("platform_workflow_registry_syncer_draining_workflows")
	if err != nil {
		return nil, err
	}

	runningEngines, err := meter.Int64Gauge("platform_workflow_registry_syncer_running_engines")
	if err != nil {
		return nil, err
	}

	completedSyncs, err := meter.Int64Counter("platform_workflow_registry_syncer_completed_syncs_total")
	if err != nil {
		return nil, err
	}

	drainStarted, err := meter.Int64Counter("platform_workflow_registry_syncer_drain_started_total")
	if err != nil {
		return nil, err
	}

	drainCompleted, err := meter.Int64Counter("platform_workflow_registry_syncer_drain_completed_total")
	if err != nil {
		return nil, err
	}

	drainDuration, err := meter.Int64Histogram("platform_workflow_registry_syncer_drain_duration_ms")
	if err != nil {
		return nil, err
	}

	deleteDeferred, err := meter.Int64Counter("platform_workflow_registry_syncer_delete_deferred_total")
	if err != nil {
		return nil, err
	}

	// Per-source metrics
	sourceHealth, err := meter.Int64Gauge("platform_workflow_registry_syncer_source_health")
	if err != nil {
		return nil, err
	}

	workflowsPerSource, err := meter.Int64Gauge("platform_workflow_registry_syncer_workflows_per_source")
	if err != nil {
		return nil, err
	}

	sourceFetchDuration, err := meter.Int64Histogram("platform_workflow_registry_syncer_source_fetch_duration_ms")
	if err != nil {
		return nil, err
	}

	sourceFetchErrors, err := meter.Int64Counter("platform_workflow_registry_syncer_source_fetch_errors_total")
	if err != nil {
		return nil, err
	}

	reconcileEventsDispatched, err := meter.Int64Histogram("platform_workflow_registry_syncer_reconcile_events_dispatched")
	if err != nil {
		return nil, err
	}

	reconcileDuration, err := meter.Int64Histogram("platform_workflow_registry_syncer_reconcile_duration_ms")
	if err != nil {
		return nil, err
	}

	reconcileEventsBackoff, err := meter.Int64Counter("platform_workflow_registry_syncer_reconcile_events_backoff_total")
	if err != nil {
		return nil, err
	}

	activationDropped, err := meter.Int64Counter("platform_workflow_registry_syncer_activation_dropped_total")
	if err != nil {
		return nil, err
	}

	activationAbandoned, err := meter.Int64Counter("platform_workflow_registry_syncer_activation_abandoned_total")
	if err != nil {
		return nil, err
	}

	moduleStoreDuration, err := meter.Int64Histogram("platform_workflow_registry_syncer_module_store_duration_ms")
	if err != nil {
		return nil, err
	}

	workflowTagBackfill, err := meter.Int64Counter("platform_workflow_registry_syncer_workflow_tag_backfill_total")
	if err != nil {
		return nil, err
	}

	return &metrics{
		handleDuration:            handleDuration,
		fetchedWorkflows:          fetchedWorkflows,
		runningWorkflows:          runningWorkflows,
		drainingWorkflows:         drainingWorkflows,
		runningEngines:            runningEngines,
		completedSyncs:            completedSyncs,
		drainStarted:              drainStarted,
		drainCompleted:            drainCompleted,
		drainDuration:             drainDuration,
		deleteDeferred:            deleteDeferred,
		sourceHealth:              sourceHealth,
		workflowsPerSource:        workflowsPerSource,
		sourceFetchDuration:       sourceFetchDuration,
		sourceFetchErrors:         sourceFetchErrors,
		reconcileEventsDispatched: reconcileEventsDispatched,
		reconcileDuration:         reconcileDuration,
		reconcileEventsBackoff:    reconcileEventsBackoff,
		activationDropped:         activationDropped,
		activationAbandoned:       activationAbandoned,
		moduleStoreDuration:       moduleStoreDuration,
		workflowTagBackfill:       workflowTagBackfill,
	}, nil
}
