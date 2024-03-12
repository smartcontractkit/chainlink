package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// AutomationNamespace is the namespace for all Automation related metrics
const AutomationLogTriggerNamespace = "automation_log_trigger"

// Automation metrics
var (
	AutomationLogsInLogBuffer = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: AutomationLogTriggerNamespace,
		Name:      "num_logs_in_log_buffer",
		Help:      "The total number of logs currently being stored in the log buffer",
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
)
