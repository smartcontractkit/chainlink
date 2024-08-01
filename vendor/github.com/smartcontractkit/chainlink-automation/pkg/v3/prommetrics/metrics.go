package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// NamespaceAutomation is the namespace for all Automation related metrics
const NamespaceAutomation = "automation"

// Plugin error types
const (
	PluginErrorTypeInvalidOracleObservation = "invalid_oracle_observation"
	PluginErrorTypeDecodeOutcome            = "decode_outcome"
	PluginErrorTypeEncodeReport             = "encode_report"
)

// Plugin steps
const (
	PluginStepResultStore = "result_store"
	PluginStepObservation = "observation"
	PluginStepOutcome     = "outcome"
	PluginStepReports     = "reports"
)

// Automation metrics
var (
	AutomationPluginPerformables = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NamespaceAutomation,
		Name:      "plugin_performables",
		Help:      "How many performables were present at a given step in the plugin flow",
	}, []string{
		"step",
	})
	AutomationPluginError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceAutomation,
		Name:      "plugin_error",
		Help:      "Count of how many errors were encountered in the plugin by label",
	}, []string{
		"step",
		"error",
	})
)
