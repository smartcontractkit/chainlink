package promwrapper

import (
	"context"
	"math/big"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var (
	_       types.ReportingPlugin = &promPlugin{}
	buckets                       = []float64{
		float64(1 * time.Millisecond),
		float64(5 * time.Millisecond),
		float64(10 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
		float64(30 * time.Second),
		float64(time.Minute),
		float64(2 * time.Minute),
		float64(5 * time.Minute),
		float64(10 * time.Minute),
	}
)

var (
	promQuery = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_query_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Query() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin"},
	)
	promObservation = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_observation_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Observation() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin"},
	)
	promReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_report_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Report() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin"},
	)
	promShouldAcceptFinalizedReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_should_accept_finalized_report_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's ShouldAcceptFinalizedReport() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin"},
	)
	promShouldTransmitAcceptedReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_should_transmit_accepted_report_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's ShouldTransmitAcceptedReport() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin"},
	)
	promClose = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_close_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Close() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin"},
	)
)

// promPlugin wraps another OCR2 reporting plugin and provides standardized prometheus metrics
// for each of the OCR2 phases (Query, Observation, Report, ShouldAcceptFinalizedReport,
// ShouldTransmitAcceptedReport, and Close).
type promPlugin struct {
	wrapped   types.ReportingPlugin
	name      string
	chainType string
	chainID   *big.Int
}

func New(plugin types.ReportingPlugin, name string, chainType string, chainID *big.Int) types.ReportingPlugin {
	return &promPlugin{
		wrapped:   plugin,
		name:      name,
		chainType: chainType,
		chainID:   chainID,
	}
}

func (p *promPlugin) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promQuery.WithLabelValues(string(p.chainType), p.chainID.String(), p.name).Observe(duration)
	}()

	return p.wrapped.Query(ctx, timestamp)
}

func (p *promPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promObservation.WithLabelValues(string(p.chainType), p.chainID.String(), p.name).Observe(duration)
	}()

	return p.wrapped.Observation(ctx, timestamp, query)
}

func (p *promPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promReport.WithLabelValues(string(p.chainType), p.chainID.String(), p.name).Observe(duration)
	}()

	return p.wrapped.Report(ctx, timestamp, query, observations)
}

func (p *promPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promShouldAcceptFinalizedReport.WithLabelValues(string(p.chainType), p.chainID.String(), p.name).Observe(duration)
	}()

	return p.wrapped.ShouldAcceptFinalizedReport(ctx, timestamp, report)
}

func (p *promPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promShouldTransmitAcceptedReport.WithLabelValues(string(p.chainType), p.chainID.String(), p.name).Observe(duration)
	}()

	return p.wrapped.ShouldTransmitAcceptedReport(ctx, timestamp, report)
}

func (p *promPlugin) Close() error {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promClose.WithLabelValues(string(p.chainType), p.chainID.String(), p.name).Observe(duration)
	}()

	return p.wrapped.Close()
}
