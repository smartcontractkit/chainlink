package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Namespaces
const AutomationLogTriggerNamespace = "automation_log_trigger"
const AutomationStreamsNamespace = "automation_streams"

// Streams steps
const StreamsLookupStepDoMercuryRequest = "do_mercury_request"
const StreamsLookupStepCheckErrorHandler = "check_error_handler"
const StreamsLookupStepCheckCallback = "check_callback"

// Streams error labels
const StreamsLookupErrorReasonNotReverted = "error_reason_not_target_check_reverted"
const StreamsLookupErrorDecodeRequestFailed = "decode_request_failed"
const StreamsLookupErrorCredentialsNotConfigured = "credentials_not_configured"
const StreamsLookupErrorDoMercuryRequest = "do_mercury_request"
const StreamsLookupErrorCodeNotNil = "err_code_not_nil"
const StreamsLookupErrorCheckCallback = "check_callback"
const StreamsLookupErrorPackUserCheckErrorHandler = "pack_user_check_error_handler"
const StreamsLookupErrorPackExecuteCallback = "pack_execute_callback"

// Metric labels
const (
	LogBufferFlowDirectionIngress = "ingress"
	LogBufferFlowDirectionEgress  = "egress"
	LogBufferFlowDirectionDropped = "dropped"
)

// Automation metrics
var (
	// Log Trigger metrics
	AutomationLogBufferFlow = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: AutomationLogTriggerNamespace,
		Name:      "num_logs_in_log_buffer",
		Help:      "The total number of logs currently being stored in the log buffer",
	}, []string{
		"direction",
	})
	AutomationRecovererMissedLogs = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: AutomationLogTriggerNamespace,
		Name:      "num_recoverer_missed_logs",
		Help:      "How many valid log triggers were identified as being missed by the recoverer",
	})
	AutomationRecovererPendingPayloads = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: AutomationLogTriggerNamespace,
		Name:      "num_recoverer_pending_payloads",
		Help:      "How many log trigger payloads are currently pending in the recoverer",
	})
	AutomationActiveUpkeeps = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: AutomationLogTriggerNamespace,
		Name:      "num_active_upkeeps",
		Help:      "How many log trigger upkeeps are currently active",
	})
	AutomationLogProviderLatestBlock = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: AutomationLogTriggerNamespace,
		Name:      "log_provider_latest_block",
		Help:      "The latest block number the log provider has seen",
	})

	// Streams metrics
	AutomationStreamsLookupStep = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: AutomationStreamsNamespace,
		Name:      "streams_lookup_step_count",
		Help:      "How many times individual steps of the streams lookup process run",
	}, []string{
		"step",
	})
	AutomationStreamsLookupError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: AutomationStreamsNamespace,
		Name:      "streams_lookup_error_count",
		Help:      "Errors occured during a streams lookup attempt",
	}, []string{
		"error",
	})
)
