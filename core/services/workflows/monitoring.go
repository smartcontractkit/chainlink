package workflows

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	monutils "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

// wm AKA "workflow metrics" is to locally scope these instruments to avoid
// data races in testing
type workflowMetrics struct {
	totalWorkflows metric.Int64Gauge
}

// InitWorkflowMonitoringResources is exported to allow the syncer to report on
// totalWorkflows, as well as the delegate in this package
func InitWorkflowMonitoringResources() (l *WorkflowMetricLabeler, err error) {
	wm := workflowMetrics{}

	wm.totalWorkflows, err = beholder.GetMeter().Int64Gauge("platform_workflow_engines_total")
	if err != nil {
		return nil, fmt.Errorf("error initializing total workflows: %w", err)
	}

	l = &WorkflowMetricLabeler{
		metrics.NewLabeler(),
		wm,
	}
	return l, nil
}

// WorkflowMetricLabeler wraps a Labeler to provide workflow specific utilities
// for monitoring resources
type WorkflowMetricLabeler struct {
	metrics.Labeler
	wm workflowMetrics
}

func (l WorkflowMetricLabeler) UpdateTotalWorkflowsGauge(ctx context.Context, val int64) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.wm.totalWorkflows.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

// em AKA "engine metrics" is to locally scope these instruments to avoid
// data races in testing
type engineMetrics struct {
	registerTriggerFailureCounter      metric.Int64Counter
	triggerWorkflowStarterErrorCounter metric.Int64Counter
	capabilityInvocationCounter        metric.Int64Counter
	capabilityFailureCounter           metric.Int64Counter
	workflowRegisteredCounter          metric.Int64Counter
	// workflowUnregisteredCounter can also be viewed as workflow deletions
	workflowUnregisteredCounter              metric.Int64Counter
	workflowExecutionRateLimitGlobalCounter  metric.Int64Counter
	workflowExecutionRateLimitPerUserCounter metric.Int64Counter
	workflowLimitGlobalCounter               metric.Int64Counter
	workflowLimitPerOwnerCounter             metric.Int64Counter
	workflowStepErrorCounter                 metric.Int64Counter
	workflowInitializationCounter            metric.Int64Counter
	engineHeartbeatCounter                   metric.Int64Counter
	workflowCompletedDurationSeconds         metric.Int64Histogram
	workflowEarlyExitDurationSeconds         metric.Int64Histogram
	workflowErrorDurationSeconds             metric.Int64Histogram
	workflowTimeoutDurationSeconds           metric.Int64Histogram
	workflowStepDurationSeconds              metric.Int64Histogram
}

func initEngineMonitoringResources() (m *engineMetricLabeler, err error) {
	em := engineMetrics{}

	em.workflowExecutionRateLimitGlobalCounter, err = beholder.GetMeter().Int64Counter("platform_engine_execution_ratelimit_global")
	if err != nil {
		return nil, fmt.Errorf("failed to register execution rate limit global counter: %w", err)
	}

	em.workflowExecutionRateLimitPerUserCounter, err = beholder.GetMeter().Int64Counter("platform_engine_execution_ratelimit_peruser")
	if err != nil {
		return nil, fmt.Errorf("failed to register execution rate limit per user counter: %w", err)
	}

	em.workflowLimitGlobalCounter, err = beholder.GetMeter().Int64Counter("platform_engine_limit_global")
	if err != nil {
		return nil, fmt.Errorf("failed to register execution limit global counter: %w", err)
	}

	em.workflowLimitPerOwnerCounter, err = beholder.GetMeter().Int64Counter("platform_engine_limit_perowner")
	if err != nil {
		return nil, fmt.Errorf("failed to register execution limit per owner counter: %w", err)
	}

	em.registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("platform_engine_registertrigger_failures")
	if err != nil {
		return nil, fmt.Errorf("failed to register trigger failure counter: %w", err)
	}

	em.triggerWorkflowStarterErrorCounter, err = beholder.GetMeter().Int64Counter("platform_engine_triggerworkflow_starter_errors")
	if err != nil {
		return nil, fmt.Errorf("failed to register trigger workflow starter error counter: %w", err)
	}

	em.capabilityInvocationCounter, err = beholder.GetMeter().Int64Counter("platform_engine_capabilities_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register capability invocation counter: %w", err)
	}

	em.capabilityFailureCounter, err = beholder.GetMeter().Int64Counter("platform_engine_capabilities_failures")
	if err != nil {
		return nil, fmt.Errorf("failed to register capability failure counter: %w", err)
	}

	em.workflowRegisteredCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_registered_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register workflow registered counter: %w", err)
	}

	em.workflowUnregisteredCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_unregistered_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register workflow unregistered counter: %w", err)
	}

	em.workflowInitializationCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_initializations")
	if err != nil {
		return nil, fmt.Errorf("failed to register workflow initialization counter: %w", err)
	}

	em.workflowStepErrorCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_errors")
	if err != nil {
		return nil, fmt.Errorf("failed to register workflow step error counter: %w", err)
	}

	em.engineHeartbeatCounter, err = beholder.GetMeter().Int64Counter("platform_engine_heartbeat")
	if err != nil {
		return nil, fmt.Errorf("failed to register engine heartbeat counter: %w", err)
	}

	em.workflowCompletedDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_completed_time_seconds",
		metric.WithDescription("Distribution of completed execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return nil, fmt.Errorf("failed to register completed duration histogram: %w", err)
	}

	em.workflowEarlyExitDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_earlyexit_time_seconds",
		metric.WithDescription("Distribution of earlyexit execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return nil, fmt.Errorf("failed to register early exit duration histogram: %w", err)
	}

	em.workflowErrorDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_error_time_seconds",
		metric.WithDescription("Distribution of error execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return nil, fmt.Errorf("failed to register error duration histogram: %w", err)
	}

	em.workflowTimeoutDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_timeout_time_seconds",
		metric.WithDescription("Distribution of timeout execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return nil, fmt.Errorf("failed to register timeout duration histogram: %w", err)
	}

	em.workflowStepDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_step_time_seconds",
		metric.WithDescription("Distribution of step execution times"),
		metric.WithUnit("seconds"))
	if err != nil {
		return nil, fmt.Errorf("failed to register step execution time histogram: %w", err)
	}

	return &engineMetricLabeler{metrics.NewLabeler(), em}, nil
}

// Note: due to the OTEL specification, all histogram buckets
// Must be defined when the beholder client is created
func MetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_engine_workflow_earlyexit_time_seconds"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 1, 10, 30, 120},
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_engine_workflow_completed_time_seconds"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// increased granularity for the workflow execution latencies near expected values
				Boundaries: []float64{0, 10, 20, 40, 50, 70, 90, 120, 150, 180, 210, 300, 600, 900, 1200},
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_engine_workflow_error_time_seconds"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 30, 60, 120, 240, 600},
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_engine_workflow_step_time_seconds"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 20, 60, 120, 240},
			}},
		),
	}
}

// engineMetricLabeler wraps monitoring.MetricsLabeler to provide workflow specific utilities
// for monitoring resources
type engineMetricLabeler struct {
	metrics.Labeler
	em engineMetrics
}

func (c engineMetricLabeler) with(keyValues ...string) engineMetricLabeler {
	return engineMetricLabeler{c.With(keyValues...), c.em}
}

func (c engineMetricLabeler) incrementWorkflowExecutionRateLimitGlobalCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowExecutionRateLimitGlobalCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementWorkflowExecutionRateLimitPerUserCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowExecutionRateLimitPerUserCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementRegisterTriggerFailureCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.registerTriggerFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementWorkflowLimitGlobalCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowLimitGlobalCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementWorkflowLimitPerOwnerCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowLimitPerOwnerCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementTriggerWorkflowStarterErrorCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.triggerWorkflowStarterErrorCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementCapabilityInvocationCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.capabilityInvocationCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementTotalWorkflowStepErrorsCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowStepErrorCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementEngineHeartbeatCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.engineHeartbeatCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementCapabilityFailureCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.capabilityFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementWorkflowRegisteredCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowRegisteredCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementWorkflowUnregisteredCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowUnregisteredCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) incrementWorkflowInitializationCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowInitializationCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) updateWorkflowCompletedDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowCompletedDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) updateWorkflowEarlyExitDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowEarlyExitDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) updateWorkflowErrorDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowErrorDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) updateWorkflowTimeoutDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowTimeoutDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c engineMetricLabeler) updateWorkflowStepDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	c.em.workflowStepDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}
