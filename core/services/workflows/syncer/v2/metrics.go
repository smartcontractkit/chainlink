package v2

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type metrics struct {
	handleDuration    metric.Int64Histogram
	fetchedWorkflows  metric.Int64Gauge
	runningWorkflows  metric.Int64Gauge
	drainingWorkflows metric.Int64Gauge
	completedSyncs    metric.Int64Counter
	drainStarted      metric.Int64Counter
	drainCompleted    metric.Int64Counter
	drainDuration     metric.Int64Histogram
	deleteDeferred    metric.Int64Counter

	// Per-source metrics for multi-source observability
	sourceHealth        metric.Int64Gauge     // 1=healthy, 0=unhealthy per source
	workflowsPerSource  metric.Int64Gauge     // workflows fetched per source
	sourceFetchDuration metric.Int64Histogram // fetch latency per source
	sourceFetchErrors   metric.Int64Counter   // error count per source

	// Per-tick reconciliation metrics
	reconcileEventsDispatched metric.Int64Histogram // events dispatched per source per tick
	reconcileDuration         metric.Int64Histogram // wall-clock ms for parallel event processing
	reconcileEventsBackoff    metric.Int64Counter   // events skipped due to backoff

	// On-disk WASM cache write (sync); duration tails indicate IO contention vs typical skew.
	moduleStoreDuration metric.Int64Histogram
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

func (m *metrics) recordModuleStore(ctx context.Context, d time.Duration, success bool) {
	if m == nil {
		return
	}
	m.moduleStoreDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", strconv.FormatBool(success)),
	))
}

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
}

func (cm *CacheMetrics) recordEviction(ctx context.Context, count int) {
	if cm == nil {
		return
	}
	cm.evictionTotal.Add(ctx, int64(count))
}

func (cm *CacheMetrics) recordLoaded(ctx context.Context, count int) {
	if cm == nil {
		return
	}
	cm.loadedGauge.Record(ctx, int64(count))
}

func (cm *CacheMetrics) recordMemorySaved(ctx context.Context, bytes int64) {
	if cm == nil {
		return
	}
	cm.memorySaved.Record(ctx, bytes)
}

func (cm *CacheMetrics) recordVersionMismatch(ctx context.Context) {
	if cm == nil {
		return
	}
	cm.versionMismatch.Add(ctx, 1)
}

func (cm *CacheMetrics) recordPinExhausted(ctx context.Context) {
	if cm == nil {
		return
	}
	cm.pinExhausted.Add(ctx, 1)
	if cachePinExhaustedHook != nil {
		cachePinExhaustedHook()
	}
}

func (cm *CacheMetrics) recordTryAcquireExhausted(ctx context.Context) {
	if cm == nil {
		return
	}
	cm.tryAcquireExhausted.Add(ctx, 1)
	if cacheTryAcquireExhaustedHook != nil {
		cacheTryAcquireExhaustedHook()
	}
}

func NewCacheMetrics() (*CacheMetrics, error) {
	reloadSource, err := beholder.GetMeter().Int64Counter("platform_workflow_module_cache_reload_total")
	if err != nil {
		return nil, err
	}
	evictionTotal, err := beholder.GetMeter().Int64Counter("platform_workflow_module_cache_eviction_total")
	if err != nil {
		return nil, err
	}
	loadedGauge, err := beholder.GetMeter().Int64Gauge("platform_workflow_module_cache_loaded")
	if err != nil {
		return nil, err
	}
	memorySaved, err := beholder.GetMeter().Int64Gauge("platform_workflow_module_cache_memory_saved_bytes")
	if err != nil {
		return nil, err
	}
	versionMismatch, err := beholder.GetMeter().Int64Counter("platform_workflow_module_cache_version_mismatch_total")
	if err != nil {
		return nil, err
	}
	pinExhausted, err := beholder.GetMeter().Int64Counter("platform_workflow_module_cache_pin_exhausted_total")
	if err != nil {
		return nil, err
	}
	tryAcquireExhausted, err := beholder.GetMeter().Int64Counter("platform_workflow_module_cache_try_acquire_exhausted_total")
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

func newMetrics() (*metrics, error) {
	handleDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	fetchedWorkflows, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_fetched_workflows")
	if err != nil {
		return nil, err
	}

	runningWorkflows, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_running_workflows")
	if err != nil {
		return nil, err
	}

	drainingWorkflows, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_draining_workflows")
	if err != nil {
		return nil, err
	}

	completedSyncs, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_completed_syncs_total")
	if err != nil {
		return nil, err
	}

	drainStarted, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_drain_started_total")
	if err != nil {
		return nil, err
	}

	drainCompleted, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_drain_completed_total")
	if err != nil {
		return nil, err
	}

	drainDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_drain_duration_ms")
	if err != nil {
		return nil, err
	}

	deleteDeferred, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_delete_deferred_total")
	if err != nil {
		return nil, err
	}

	// Per-source metrics
	sourceHealth, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_source_health")
	if err != nil {
		return nil, err
	}

	workflowsPerSource, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_workflows_per_source")
	if err != nil {
		return nil, err
	}

	sourceFetchDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_source_fetch_duration_ms")
	if err != nil {
		return nil, err
	}

	sourceFetchErrors, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_source_fetch_errors_total")
	if err != nil {
		return nil, err
	}

	reconcileEventsDispatched, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_reconcile_events_dispatched")
	if err != nil {
		return nil, err
	}

	reconcileDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_reconcile_duration_ms")
	if err != nil {
		return nil, err
	}

	reconcileEventsBackoff, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_reconcile_events_backoff_total")
	if err != nil {
		return nil, err
	}

	moduleStoreDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_module_store_duration_ms")
	if err != nil {
		return nil, err
	}

	return &metrics{
		handleDuration:            handleDuration,
		fetchedWorkflows:          fetchedWorkflows,
		runningWorkflows:          runningWorkflows,
		drainingWorkflows:         drainingWorkflows,
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
		moduleStoreDuration:       moduleStoreDuration,
	}, nil
}
