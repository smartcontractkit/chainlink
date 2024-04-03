package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Namespaces
const (
	NamespaceAutomationLogTrigger = "automation_log_trigger"
	NamespaceAutomationStreams    = "automation_streams"
)

// Streams steps
const (
	StreamsLookupStepDoMercuryRequest  = "do_mercury_request"
	StreamsLookupStepCheckErrorHandler = "check_error_handler"
	StreamsLookupStepCheckCallback     = "check_callback"
)

// Streams error labels
const (
	StreamsLookupErrorReasonNotReverted         = "reason_not_target_check_reverted"
	StreamsLookupErrorDecodeRequestFailed       = "decode_request_failed"
	StreamsLookupErrorCredentialsNotConfigured  = "credentials_not_configured"
	StreamsLookupErrorDoMercuryRequest          = "do_mercury_request"
	StreamsLookupErrorCodeNotNil                = "err_code_not_nil"
	StreamsLookupErrorCheckCallback             = "check_callback"
	StreamsLookupErrorPackUserCheckErrorHandler = "pack_user_check_error_handler"
	StreamsLookupErrorPackExecuteCallback       = "pack_execute_callback"
)

// Streams versions
const (
	StreamsVersion02 = "v02"
	StreamsVersion03 = "v03"
)

// Metric labels
const (
	LogBufferFlowDirectionIngress = "ingress"
	LogBufferFlowDirectionEgress  = "egress"
	LogBufferFlowDirectionDropped = "dropped"
	LogBufferFlowDirectionExpired = "expired"
)

// Automation metrics
var (
	// Log Trigger metrics
	AutomationLogBufferFlow = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceAutomationLogTrigger,
		Name:      "num_logs_in_log_buffer",
		Help:      "The total number of logs currently being stored in the log buffer",
	}, []string{
		"direction",
	})
	AutomationRecovererMissedLogs = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: NamespaceAutomationLogTrigger,
		Name:      "num_recoverer_missed_logs",
		Help:      "How many valid log triggers were identified as being missed by the recoverer",
	})
	AutomationRecovererPendingPayloads = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: NamespaceAutomationLogTrigger,
		Name:      "num_recoverer_pending_payloads",
		Help:      "How many log trigger payloads are currently pending in the recoverer",
	})
	AutomationActiveUpkeeps = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: NamespaceAutomationLogTrigger,
		Name:      "num_active_upkeeps",
		Help:      "How many log trigger upkeeps are currently active",
	})
	AutomationLogProviderLatestBlock = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: NamespaceAutomationLogTrigger,
		Name:      "log_provider_latest_block",
		Help:      "The latest block number the log provider has seen",
	})

	// Streams metrics
	AutomationStreamsLookupStep = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceAutomationStreams,
		Name:      "streams_lookup_step_count",
		Help:      "How many times individual steps of the streams lookup process run",
	}, []string{
		"step",
	})
	AutomationStreamsLookupError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceAutomationStreams,
		Name:      "streams_lookup_error_count",
		Help:      "Errors occurred during a streams lookup attempt",
	}, []string{
		"error",
	})
	AutomationStreamsRetries = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceAutomationStreams,
		Name:      "streams_retries",
		Help:      "Count of the times a streams lookup was retried",
	}, []string{
		"version",
	})
	AutomationStreamsResponses = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceAutomationStreams,
		Name:      "streams_responses",
		Help:      "Count of individual response codes from streams lookup",
	}, []string{
		"version",
		"status",
	})
)
