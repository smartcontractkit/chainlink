package workflows

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	monutils "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

var registerTriggerFailureCounter metric.Int64Counter
var triggerWorkflowStarterErrorCounter metric.Int64Counter
var workflowsRunningGauge metric.Int64Gauge
var capabilityInvocationCounter metric.Int64Counter
var capabilityFailureCounter metric.Int64Counter
var workflowRegisteredCounter metric.Int64Counter
var workflowUnregisteredCounter metric.Int64Counter
var workflowExecutionLatencyGauge metric.Int64Gauge // ms
var workflowStepErrorCounter metric.Int64Counter
var workflowInitializationCounter metric.Int64Counter
var engineHeartbeatCounter metric.Int64Counter
var workflowCompletedDurationSeconds metric.Int64Histogram
var workflowEarlyExitDurationSeconds metric.Int64Histogram
var workflowErrorDurationSeconds metric.Int64Histogram
var workflowTimeoutDurationSeconds metric.Int64Histogram

func initMonitoringResources() (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("platform_engine_registertrigger_failures")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure counter: %w", err)
	}

	triggerWorkflowStarterErrorCounter, err = beholder.GetMeter().Int64Counter("platform_engine_triggerworkflow_starter_errors")
	if err != nil {
		return fmt.Errorf("failed to register trigger workflow starter error counter: %w", err)
	}

	workflowsRunningGauge, err = beholder.GetMeter().Int64Gauge("platform_engine_workflow_count")
	if err != nil {
		return fmt.Errorf("failed to register workflows running gauge: %w", err)
	}

	capabilityInvocationCounter, err = beholder.GetMeter().Int64Counter("platform_engine_capabilities_count")
	if err != nil {
		return fmt.Errorf("failed to register capability invocation counter: %w", err)
	}

	capabilityFailureCounter, err = beholder.GetMeter().Int64Counter("platform_engine_capabilities_failures")
	if err != nil {
		return fmt.Errorf("failed to register capability failure counter: %w", err)
	}

	workflowRegisteredCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_registered_count")
	if err != nil {
		return fmt.Errorf("failed to register workflow registered counter: %w", err)
	}

	workflowUnregisteredCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_unregistered_count")
	if err != nil {
		return fmt.Errorf("failed to register workflow unregistered counter: %w", err)
	}

	workflowExecutionLatencyGauge, err = beholder.GetMeter().Int64Gauge("platform_engine_workflow_time")
	if err != nil {
		return fmt.Errorf("failed to register workflow execution latency gauge: %w", err)
	}

	workflowInitializationCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_initializations")
	if err != nil {
		return fmt.Errorf("failed to register workflow initialization counter: %w", err)
	}

	workflowStepErrorCounter, err = beholder.GetMeter().Int64Counter("platform_engine_workflow_errors")
	if err != nil {
		return fmt.Errorf("failed to register workflow step error counter: %w", err)
	}

	engineHeartbeatCounter, err = beholder.GetMeter().Int64Counter("platform_engine_heartbeat")
	if err != nil {
		return fmt.Errorf("failed to register engine heartbeat counter: %w", err)
	}

	workflowCompletedDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_completed_time_seconds",
		metric.WithDescription("Distribution of completed execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return fmt.Errorf("failed to register completed duration histogram: %w", err)
	}

	workflowEarlyExitDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_earlyexit_time_seconds",
		metric.WithDescription("Distribution of earlyexit execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return fmt.Errorf("failed to register early exit duration histogram: %w", err)
	}

	workflowErrorDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_error_time_seconds",
		metric.WithDescription("Distribution of error execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return fmt.Errorf("failed to register error duration histogram: %w", err)
	}

	workflowTimeoutDurationSeconds, err = beholder.GetMeter().Int64Histogram(
		"platform_engine_workflow_timeout_time_seconds",
		metric.WithDescription("Distribution of timeout execution latencies"),
		metric.WithUnit("seconds"))
	if err != nil {
		return fmt.Errorf("failed to register timeout duration histogram: %w", err)
	}

	return nil
}

// workflowsMetricLabeler wraps monitoring.MetricsLabeler to provide workflow specific utilities
// for monitoring resources
type workflowsMetricLabeler struct {
	metrics.Labeler
}

func (c workflowsMetricLabeler) with(keyValues ...string) workflowsMetricLabeler {
	return workflowsMetricLabeler{c.With(keyValues...)}
}

func (c workflowsMetricLabeler) incrementRegisterTriggerFailureCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	registerTriggerFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementTriggerWorkflowStarterErrorCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	triggerWorkflowStarterErrorCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementCapabilityInvocationCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	capabilityInvocationCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateWorkflowExecutionLatencyGauge(ctx context.Context, val int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowExecutionLatencyGauge.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementTotalWorkflowStepErrorsCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowStepErrorCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateTotalWorkflowsGauge(ctx context.Context, val int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowsRunningGauge.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementEngineHeartbeatCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	engineHeartbeatCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementCapabilityFailureCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	capabilityFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementWorkflowRegisteredCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowRegisteredCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementWorkflowUnregisteredCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowUnregisteredCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementWorkflowInitializationCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowInitializationCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateWorkflowCompletedDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowCompletedDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateWorkflowEarlyExitDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowEarlyExitDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateWorkflowErrorDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowErrorDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateWorkflowTimeoutDurationHistogram(ctx context.Context, duration int64) {
	otelLabels := monutils.KvMapToOtelAttributes(c.Labels)
	workflowTimeoutDurationSeconds.Record(ctx, duration, metric.WithAttributes(otelLabels...))
}
