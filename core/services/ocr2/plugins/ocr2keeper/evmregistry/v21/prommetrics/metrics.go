package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// AutomationNamespace is the namespace for all Automation related metrics
const AutomationNamespace = "automation"

// Automation metrics
var (
	AutomationLogsInLogBuffer = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: AutomationNamespace,
		Name:      "num_logs_in_log_buffer",
		Help:      "The total number of logs currently being stored in the log buffer",
	})
)
