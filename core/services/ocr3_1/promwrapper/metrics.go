// Package promwrapper instruments an OCR3.1 ReportingPlugin with prometheus
// metrics. It is the OCR3.1 counterpart of core/services/ocr3/promwrapper and
// mirrors that package's structure (see also the ocr3 / ocr3_1 split of
// beholderwrapper).
//
// It emits the same ocr3_reporting_plugin_* metric series as the OCR3.0
// wrapper — the two versions share one metric surface, differentiated by the
// "function" label (which includes the OCR3.1-only phases observationQuorum,
// stateTransition and committed). To keep that surface identical without
// importing the OCR3.0 package, the metric collectors is registered with
// registerOrExisting to avoid runtime issues.
package promwrapper

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type functionType string

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

	promOCR3ReportsGenerated = registerOrExisting(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocr3_reporting_plugin_reports_processed",
			Help: "Tracks number of reports processed/generated within by different OCR3 functions",
		},
		[]string{"chainFamily", "chainID", "plugin", "function"},
	))
	promOCR3Durations = registerOrExisting(prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr3_reporting_plugin_duration",
			Help:    "The amount of time elapsed during the OCR3 plugin's function",
			Buckets: buckets,
		},
		[]string{"chainFamily", "chainID", "plugin", "function", "success"},
	))
	promOCR3Sizes = registerOrExisting(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocr3_reporting_plugin_data_sizes",
			Help: "Tracks the size of the data produced by OCR3 plugin in bytes (e.g. reports, observations etc.)",
		},
		[]string{"chainFamily", "chainID", "plugin", "function"},
	))
	promOCR3PluginStatus = registerOrExisting(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocr3_reporting_plugin_status",
			Help: "Gauge indicating whether plugin is up and running or not",
		},
		[]string{"chainFamily", "chainID", "plugin", "configDigest"},
	))
)

// registerOrExisting registers c on the default registerer, or if an equal
// collector is already registered (e.g. because core/services/ocr3/promwrapper
// is also linked into this binary), returns that existing collector so both
// packages write to a single shared metric series.
func registerOrExisting[C prometheus.Collector](c C) C {
	if err := prometheus.DefaultRegisterer.Register(c); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if existing, ok := are.ExistingCollector.(C); ok {
				return existing
			}
		}
		panic(err)
	}
	return c
}

func boolToInt(arg bool) int {
	if arg {
		return 1
	}
	return 0
}
