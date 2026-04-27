package promwrapper

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type functionType string

// OCR3_1 exposes the same Query/Observation/Reports/ShouldAccept/
// ShouldTransmit surface as OCR3 but adds StateTransition (replacing
// Outcome), Committed, and explicit ObservationQuorum. Metric labels
// track each separately so dashboards can distinguish the new critical-
// path methods.
const (
	query               functionType = "query"
	observation         functionType = "observation"
	validateObservation functionType = "validateObservation"
	observationQuorum   functionType = "observationQuorum"
	stateTransition     functionType = "stateTransition"
	committed           functionType = "committed"
	reports             functionType = "reports"
	shouldAccept        functionType = "shouldAccept"
	shouldTransmit      functionType = "shouldTransmit"
)

var (
	buckets = []float64{
		float64(10 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(700 * time.Millisecond),
		float64(time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
		float64(20 * time.Second),
		float64(30 * time.Second),
	}

	// Distinct metric names from the OCR3 promwrapper so dashboards can
	// slice OCR3 vs OCR3_1 DONs independently during the coexistence
	// window. Same label dimensions as OCR3 for easy template reuse.
	promOCR3_1ReportsGenerated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocr3_1_reporting_plugin_reports_processed",
			Help: "Tracks number of reports processed/generated within by different OCR3_1 functions",
		},
		[]string{"chainFamily", "chainID", "plugin", "function"},
	)
	promOCR3_1Durations = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr3_1_reporting_plugin_duration",
			Help:    "The amount of time elapsed during the OCR3_1 plugin's function",
			Buckets: buckets,
		},
		[]string{"chainFamily", "chainID", "plugin", "function", "success"},
	)
	promOCR3_1Sizes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocr3_1_reporting_plugin_data_sizes",
			Help: "Tracks the size of the data produced by OCR3_1 plugin in bytes (observations, precursors, reports)",
		},
		[]string{"chainFamily", "chainID", "plugin", "function"},
	)
	promOCR3_1PluginStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocr3_1_reporting_plugin_status",
			Help: "Gauge indicating whether the OCR3_1 plugin is up and running",
		},
		[]string{"chainFamily", "chainID", "plugin", "configDigest"},
	)
)
