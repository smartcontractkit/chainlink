package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// Attribute constants for consistent labeling
const (
	AttrNodeAddress = "node_address"
	AttrNodeName    = "node_name"
	AttrStatusCode  = "status_code"
	AttrErrorCode   = "error_code"
	AttrErrorString = "error_string"
	AttrMethodName  = "method_name"
)

// CommonMetrics contains shared metrics between action and trigger handlers
type CommonMetrics struct {
	capabilityNodeThrottled metric.Int64Counter
	globalThrottled         metric.Int64Counter
}

// ActionMetrics contains metrics for HTTP actions
type ActionMetrics struct {
	requestCount                   metric.Int64Counter
	requestFailures                metric.Int64Counter
	requestLatency                 metric.Int64Histogram
	customerEndpointRequestLatency metric.Int64Histogram
	customerEndpointResponseCount  metric.Int64Counter
	cacheReadCount                 metric.Int64Counter
	cacheHitCount                  metric.Int64Counter
	cacheCleanUpCount              metric.Int64Counter
	cacheSize                      metric.Int64Gauge
	capabilityRequestCount         metric.Int64Counter
	capabilityFailures             metric.Int64Counter
	blockedRequestCount            metric.Int64Counter
	httpSendErrorCount             metric.Int64Counter
	httpReadErrorCount             metric.Int64Counter
}

// TriggerMetrics contains metrics for HTTP triggers
type TriggerMetrics struct {
	requestCount                     metric.Int64Counter
	requestErrors                    metric.Int64Counter
	requestSuccess                   metric.Int64Counter
	workflowThrottled                metric.Int64Counter
	pendingRequestsCleanUpCount      metric.Int64Counter
	pendingRequestsCount             metric.Int64Gauge
	requestHandlerLatency            metric.Int64Histogram
	capabilityRequestCount           metric.Int64Counter
	capabilityRequestFailures        metric.Int64Counter
	metadataProcessingFailures       metric.Int64Counter
	metadataRequestCount             metric.Int64Counter
	metadataObservationsCleanUpCount metric.Int64Counter
	metadataObservationsCount        metric.Int64Gauge
	jwtCacheSize                     metric.Int64Gauge
	jwtCacheCleanUpCount             metric.Int64Counter
	metadataSyncStartupLatency       metric.Int64Histogram
	loadedMetadataSize               metric.Int64Gauge
}

// Metrics combines all gateway metrics for dependency injection
type Metrics struct {
	common                *CommonMetrics
	action                *ActionMetrics
	trigger               *TriggerMetrics
	nodeAddressToNodeName map[string]string
}

// NewMetrics creates a new instance of Metrics with all metrics initialized
func NewMetrics(donConfig *config.DONConfig) (*Metrics, error) {
	meter := beholder.GetMeter()

	common, err := newCommonMetrics(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create common metrics: %w", err)
	}

	action, err := newActionMetrics(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create action metrics: %w", err)
	}

	trigger, err := newTriggerMetrics(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create trigger metrics: %w", err)
	}

	nodeAddressToNodeName := make(map[string]string)
	if donConfig != nil {
		for _, member := range donConfig.Members {
			nodeAddressToNodeName[member.Address] = member.Name
		}
	}

	return &Metrics{
		common:                common,
		action:                action,
		trigger:               trigger,
		nodeAddressToNodeName: nodeAddressToNodeName,
	}, nil
}

// newCommonMetrics initializes common metrics
func newCommonMetrics(meter metric.Meter) (*CommonMetrics, error) {
	m := &CommonMetrics{}

	var err error
	m.capabilityNodeThrottled, err = meter.Int64Counter(
		"http_handler_capability_node_throttled",
		metric.WithDescription("Number of calls from the capability node to the gateway throttled due to per-capability-node rate limit"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP handler capability node throttled metric: %w", err)
	}

	m.globalThrottled, err = meter.Int64Counter(
		"http_handler_global_throttled",
		metric.WithDescription("Number of calls from the capability node to the gateway throttled due to global rate limit"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP handler global throttled metric: %w", err)
	}

	return m, nil
}

// newActionMetrics initializes action metrics
func newActionMetrics(meter metric.Meter) (*ActionMetrics, error) {
	m := &ActionMetrics{}

	var err error
	m.requestCount, err = meter.Int64Counter(
		"http_action_gateway_request_count",
		metric.WithDescription("Number of HTTP action requests received by the gateway"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action gateway request count metric: %w", err)
	}

	m.requestFailures, err = meter.Int64Counter(
		"http_action_gateway_request_failures",
		metric.WithDescription("Number of HTTP action request failures in the gateway"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action gateway request failures metric: %w", err)
	}

	m.requestLatency, err = meter.Int64Histogram(
		"http_action_gateway_request_latency_ms",
		metric.WithDescription("HTTP action request latency in milliseconds in the gateway"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action gateway request latency metric: %w", err)
	}

	m.customerEndpointRequestLatency, err = meter.Int64Histogram(
		"http_action_customer_endpoint_request_latency_ms",
		metric.WithDescription("Request latency while calling customer endpoint in milliseconds"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action customer endpoint request latency metric: %w", err)
	}

	m.customerEndpointResponseCount, err = meter.Int64Counter(
		"http_action_customer_endpoint_response_count",
		metric.WithDescription("Number of customer endpoint responses by status code"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action customer endpoint response count metric: %w", err)
	}

	m.cacheReadCount, err = meter.Int64Counter(
		"http_action_cache_read_count",
		metric.WithDescription("Number of HTTP action cache read operations"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action response cache read count metric: %w", err)
	}

	m.cacheHitCount, err = meter.Int64Counter(
		"http_action_cache_hit_count",
		metric.WithDescription("Number of HTTP action response cache hits"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action cache hit count metric: %w", err)
	}

	m.cacheCleanUpCount, err = meter.Int64Counter(
		"http_action_cache_cleanup_count",
		metric.WithDescription("Number of HTTP action response cache entries cleaned up"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action cache cleanup count metric: %w", err)
	}

	m.cacheSize, err = meter.Int64Gauge(
		"http_action_cache_size",
		metric.WithDescription("Current number of entries in HTTP action response cache"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action cache size metric: %w", err)
	}

	m.capabilityRequestCount, err = meter.Int64Counter(
		"http_action_gateway_capability_request_count",
		metric.WithDescription("Number of gateway responses to the capability nodes for HTTP action"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action gateway capability request count metric: %w", err)
	}

	m.capabilityFailures, err = meter.Int64Counter(
		"http_action_gateway_capability_failures",
		metric.WithDescription("Number of errors while responding to the capability nodes for HTTP action"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action gateway capability failures metric: %w", err)
	}

	m.blockedRequestCount, err = meter.Int64Counter(
		"http_action_blocked_request_count",
		metric.WithDescription("Number of HTTP action requests blocked due to invalid input (blocked IP, invalid method, etc)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action blocked request count metric: %w", err)
	}

	m.httpSendErrorCount, err = meter.Int64Counter(
		"http_action_http_send_error_count",
		metric.WithDescription("Number of HTTP send errors (network failures, connection errors, timeouts during request)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action HTTP send error count metric: %w", err)
	}

	m.httpReadErrorCount, err = meter.Int64Counter(
		"http_action_http_read_error_count",
		metric.WithDescription("Number of HTTP read errors (failures reading response body)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP action HTTP read error count metric: %w", err)
	}

	return m, nil
}

// newTriggerMetrics initializes trigger metrics
func newTriggerMetrics(meter metric.Meter) (*TriggerMetrics, error) {
	m := &TriggerMetrics{}

	var err error
	m.requestCount, err = meter.Int64Counter(
		"http_trigger_gateway_request_count",
		metric.WithDescription("Number of user HTTP trigger requests received by the gateway"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway request count metric: %w", err)
	}

	m.requestErrors, err = meter.Int64Counter(
		"http_trigger_gateway_request_errors",
		metric.WithDescription("Number of HTTP trigger gateway request errors"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway request errors metric: %w", err)
	}

	m.requestSuccess, err = meter.Int64Counter(
		"http_trigger_gateway_successful_requests",
		metric.WithDescription("Number of successful HTTP trigger gateway requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway successful requests metric: %w", err)
	}

	m.workflowThrottled, err = meter.Int64Counter(
		"http_trigger_gateway_workflow_throttled",
		metric.WithDescription("Number of HTTP trigger gateway requests throttled per workflow"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway workflow throttled metric: %w", err)
	}

	m.pendingRequestsCleanUpCount, err = meter.Int64Counter(
		"http_trigger_gateway_pending_requests_cleanup_count",
		metric.WithDescription("Number of pending HTTP trigger gateway requests cleaned up"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway pending requests cleanup count metric: %w", err)
	}

	m.pendingRequestsCount, err = meter.Int64Gauge(
		"http_trigger_gateway_pending_requests_count",
		metric.WithDescription("Current number of pending HTTP trigger gateway requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway pending requests count metric: %w", err)
	}

	m.requestHandlerLatency, err = meter.Int64Histogram(
		"http_trigger_gateway_request_handler_latency_ms",
		metric.WithDescription("HTTP trigger gateway request handler latency in milliseconds"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway request latency metric: %w", err)
	}

	m.capabilityRequestCount, err = meter.Int64Counter(
		"http_trigger_gateway_capability_request_count",
		metric.WithDescription("Number of HTTP trigger requests sent from gateway node to capability nodes"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway capability request count metric: %w", err)
	}

	m.capabilityRequestFailures, err = meter.Int64Counter(
		"http_trigger_gateway_capability_request_failures",
		metric.WithDescription("Number of errors while sending HTTP trigger requests from gateway node to capability nodes"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway capability request failures metric: %w", err)
	}

	m.metadataProcessingFailures, err = meter.Int64Counter(
		"http_trigger_gateway_metadata_processing_failures",
		metric.WithDescription("Number of HTTP trigger gateway metadata processing failures"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway metadata processing failures metric: %w", err)
	}

	m.metadataRequestCount, err = meter.Int64Counter(
		"http_trigger_gateway_metadata_request_count",
		metric.WithDescription("Number of HTTP trigger gateway metadata requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger gateway metadata request count metric: %w", err)
	}

	m.metadataObservationsCleanUpCount, err = meter.Int64Counter(
		"http_trigger_metadata_observations_clean_count",
		metric.WithDescription("Number of workflow metadata observations cleaned"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow metadata observations clean count metric: %w", err)
	}

	m.metadataObservationsCount, err = meter.Int64Gauge(
		"http_trigger_metadata_observations_count",
		metric.WithDescription("Current number of workflow metadata observations in memory"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow metadata observations count metric: %w", err)
	}

	m.jwtCacheSize, err = meter.Int64Gauge(
		"http_trigger_jwt_cache_size",
		metric.WithDescription("Current number of entries in JWT replay protection cache"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger JWT cache size metric: %w", err)
	}

	m.jwtCacheCleanUpCount, err = meter.Int64Counter(
		"http_trigger_jwt_cache_cleanup_count",
		metric.WithDescription("Number of JWT replay protection cache entries cleaned up"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger JWT cache cleanup count metric: %w", err)
	}

	m.metadataSyncStartupLatency, err = meter.Int64Histogram(
		"http_trigger_metadata_sync_startup_latency_ms",
		metric.WithDescription("Time in milliseconds from handler start to first successful workflow metadata sync (i.e. first workflow loaded)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger metadata sync startup latency metric: %w", err)
	}

	m.loadedMetadataSize, err = meter.Int64Gauge(
		"http_trigger_loaded_metadata_size",
		metric.WithDescription("Number of workflows loaded for authorization after f+1 identical metadata received from workflow nodes"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP trigger loaded metadata size metric: %w", err)
	}

	return m, nil
}

// Common Metrics Methods

func (m *Metrics) IncrementCapabilityNodeThrottled(ctx context.Context, nodeAddress string, lggr logger.Logger) {
	m.common.capabilityNodeThrottled.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
	))
}

func (m *Metrics) IncrementGlobalThrottled(ctx context.Context, lggr logger.Logger) {
	m.common.globalThrottled.Add(ctx, 1)
}

// Action Metrics Methods

func (m *Metrics) IncrementActionRequestCount(ctx context.Context, nodeAddress string, lggr logger.Logger) {
	m.action.requestCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
	))
}

func (m *Metrics) IncrementActionRequestFailures(ctx context.Context, nodeAddress string, lggr logger.Logger) {
	m.action.requestFailures.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
	))
}

func (m *Metrics) RecordActionRequestLatency(ctx context.Context, latencyMs int64, lggr logger.Logger) {
	m.action.requestLatency.Record(ctx, latencyMs)
}

func (m *Metrics) RecordCustomerEndpointRequestLatency(ctx context.Context, latencyMs int64, lggr logger.Logger) {
	m.action.customerEndpointRequestLatency.Record(ctx, latencyMs)
}

func (m *Metrics) IncrementCustomerEndpointResponseCount(ctx context.Context, statusCode string, lggr logger.Logger) {
	m.action.customerEndpointResponseCount.Add(ctx, 1, metric.WithAttributes(attribute.String(AttrStatusCode, statusCode)))
}

func (m *Metrics) IncrementCacheReadCount(ctx context.Context, lggr logger.Logger) {
	m.action.cacheReadCount.Add(ctx, 1)
}

func (m *Metrics) IncrementCacheHitCount(ctx context.Context, lggr logger.Logger) {
	m.action.cacheHitCount.Add(ctx, 1)
}

func (m *Metrics) IncrementCacheCleanUpCount(ctx context.Context, count int64, lggr logger.Logger) {
	m.action.cacheCleanUpCount.Add(ctx, count)
}

func (m *Metrics) RecordCacheSize(ctx context.Context, size int64, lggr logger.Logger) {
	m.action.cacheSize.Record(ctx, size)
}

func (m *Metrics) IncrementActionCapabilityRequestCount(ctx context.Context, nodeAddress string, lggr logger.Logger) {
	m.action.capabilityRequestCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
	))
}

func (m *Metrics) IncrementActionCapabilityFailures(ctx context.Context, nodeAddress string, lggr logger.Logger) {
	m.action.capabilityFailures.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
	))
}

func (m *Metrics) IncrementBlockedRequestCount(ctx context.Context, lggr logger.Logger) {
	m.action.blockedRequestCount.Add(ctx, 1)
}

func (m *Metrics) IncrementHTTPSendErrorCount(ctx context.Context, lggr logger.Logger) {
	m.action.httpSendErrorCount.Add(ctx, 1)
}

func (m *Metrics) IncrementHTTPReadErrorCount(ctx context.Context, lggr logger.Logger) {
	m.action.httpReadErrorCount.Add(ctx, 1)
}

// Trigger Metrics Methods

func (m *Metrics) IncrementTriggerRequestCount(ctx context.Context, lggr logger.Logger) {
	m.trigger.requestCount.Add(ctx, 1)
}

func (m *Metrics) IncrementRequestErrors(ctx context.Context, errorCode int64, lggr logger.Logger) {
	errorString := api.FromJSONRPCErrorCode(errorCode).String()
	m.trigger.requestErrors.Add(ctx, 1, metric.WithAttributes(
		attribute.Int64(AttrErrorCode, errorCode),
		attribute.String(AttrErrorString, errorString),
	))
}

func (m *Metrics) IncrementRequestSuccess(ctx context.Context, lggr logger.Logger) {
	m.trigger.requestSuccess.Add(ctx, 1)
}

func (m *Metrics) IncrementWorkflowThrottled(ctx context.Context, lggr logger.Logger) {
	m.trigger.workflowThrottled.Add(ctx, 1)
}

func (m *Metrics) IncrementPendingRequestsCleanUpCount(ctx context.Context, count int64, lggr logger.Logger) {
	m.trigger.pendingRequestsCleanUpCount.Add(ctx, count)
}

func (m *Metrics) RecordPendingRequestsCount(ctx context.Context, count int64, lggr logger.Logger) {
	m.trigger.pendingRequestsCount.Record(ctx, count)
}

func (m *Metrics) RecordRequestHandlerLatency(ctx context.Context, latencyMs int64, lggr logger.Logger) {
	m.trigger.requestHandlerLatency.Record(ctx, latencyMs)
}

func (m *Metrics) IncrementTriggerCapabilityRequestCount(ctx context.Context, nodeAddress string, methodName string, lggr logger.Logger) {
	m.trigger.capabilityRequestCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
		attribute.String(AttrMethodName, methodName),
	))
}

func (m *Metrics) IncrementTriggerCapabilityRequestFailures(ctx context.Context, nodeAddress string, methodName string, lggr logger.Logger) {
	m.trigger.capabilityRequestFailures.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
		attribute.String(AttrMethodName, methodName),
	))
}

func (m *Metrics) IncrementMetadataProcessingFailures(ctx context.Context, nodeAddress string, methodName string, lggr logger.Logger) {
	m.trigger.metadataProcessingFailures.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
		attribute.String(AttrMethodName, methodName),
	))
}

func (m *Metrics) IncrementMetadataRequestCount(ctx context.Context, nodeAddress string, methodName string, lggr logger.Logger) {
	m.trigger.metadataRequestCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String(AttrNodeAddress, nodeAddress),
		attribute.String(AttrNodeName, m.nodeAddressToNodeName[nodeAddress]),
		attribute.String(AttrMethodName, methodName),
	))
}

func (m *Metrics) IncrementMetadataObservationsCleanUpCount(ctx context.Context, count int64, lggr logger.Logger) {
	m.trigger.metadataObservationsCleanUpCount.Add(ctx, count)
}

func (m *Metrics) RecordMetadataObservationsCount(ctx context.Context, count int64, lggr logger.Logger) {
	m.trigger.metadataObservationsCount.Record(ctx, count)
}

func (m *Metrics) RecordJwtCacheSize(ctx context.Context, size int64, lggr logger.Logger) {
	m.trigger.jwtCacheSize.Record(ctx, size)
}

func (m *Metrics) IncrementJwtCacheCleanUpCount(ctx context.Context, count int64, lggr logger.Logger) {
	m.trigger.jwtCacheCleanUpCount.Add(ctx, count)
}

func (m *Metrics) RecordMetadataSyncStartupLatency(ctx context.Context, latencyMs int64, lggr logger.Logger) {
	m.trigger.metadataSyncStartupLatency.Record(ctx, latencyMs)
}

func (m *Metrics) RecordLoadedMetadataSize(ctx context.Context, size int64, lggr logger.Logger) {
	m.trigger.loadedMetadataSize.Record(ctx, size)
}
