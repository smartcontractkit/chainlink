package promwrapper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type functionType string

const (
	query               functionType = "query"
	observation         functionType = "observation"
	validateObservation functionType = "validateObservation"
	outcome             functionType = "outcome"
	reports             functionType = "reports"
	shouldAccept        functionType = "shouldAccept"
	shouldTransmit      functionType = "shouldTransmit"
)

var (
	buckets = []float64{
		float64(1 * time.Millisecond),
		float64(5 * time.Millisecond),
		float64(10 * time.Millisecond),
		float64(25 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(75 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(700 * time.Millisecond),
		float64(time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
	}

	promOCR3ReportsGenerated = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocr3_reports_generated",
			Help: "Tracks number of reports generated withing a single OCR3's Reports step",
		},
		[]string{"chainID", "plugin"},
	)
	promOCR3Durations = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr3_phase_duration",
			Help:    "The amount of time elapsed during the OCR3 plugin's function",
			Buckets: buckets,
		},
		[]string{"chainID", "plugin", "function", "success"},
	)
)
