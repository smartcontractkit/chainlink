package promwrapper

import (
	"context"
	"math/big"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.ReportingPlugin = &promPlugin{}

var (
	promQuery = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_query_duration",
			Help:    "The amount of time elapsed during the OCR2 plugin's Query() method",
			Buckets: []float64{}, // TODO: figure out buckets
		},
		[]string{"evmChainID", "pluginName"},
	)
	promObservation = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_observation_duration",
			Help:    "The amount of time elapsed during the OCR2 plugin's Observation() method",
			Buckets: []float64{}, // TODO: figure out buckets
		},
		[]string{"evmChainID", "pluginName"},
	)
	promReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_report_duration",
			Help:    "The amount of time elapsed during the OCR2 plugin's Report() method",
			Buckets: []float64{}, // TODO: figure out buckets
		},
		[]string{"evmChainID", "pluginName"},
	)
	promShouldAcceptFinalizedReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_should_accept_finalized_report_duration",
			Help:    "The amount of time elapsed during the OCR2 plugin's ShouldAcceptFinalizedReport() method",
			Buckets: []float64{}, // TODO: figure out buckets
		},
		[]string{"evmChainID", "pluginName"},
	)
	promShouldTransmitAcceptedReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_should_transmit_accepted_report_duration",
			Help:    "The amount of time elapsed during the OCR2 plugin's ShouldTransmitAcceptedReport() method",
			Buckets: []float64{}, // TODO: figure out buckets
		},
		[]string{"evmChainID", "pluginName"},
	)
	promClose = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_close_duration",
			Help:    "The amount of time elapsed during the OCR2 plugin's Close() method",
			Buckets: []float64{}, // TODO: figure out buckets
		},
		[]string{"evmChainID", "pluginName"},
	)
)

// promPlugin wraps another OCR2 reporting plugin and provides standardized prometheus metrics
// for each of the OCR2 phases (Query, Observation, Report, ShouldAcceptFinalizedReport,
// ShouldTransmitAcceptedReport, and Close).
type promPlugin struct {
	wrapped    types.ReportingPlugin
	pluginName string
	evmChainID *big.Int
}

func New(plugin types.ReportingPlugin, pluginName string, evmChainID *big.Int) types.ReportingPlugin {
	return &promPlugin{
		wrapped:    plugin,
		pluginName: pluginName,
		evmChainID: evmChainID,
	}
}

func (p *promPlugin) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promQuery.WithLabelValues(p.evmChainID.String(), p.pluginName).Observe(duration)
	}()

	return p.wrapped.Query(ctx, timestamp)
}

func (p *promPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promObservation.WithLabelValues(p.evmChainID.String(), p.pluginName).Observe(duration)
	}()

	return p.wrapped.Observation(ctx, timestamp, query)
}

func (p *promPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promReport.WithLabelValues(p.evmChainID.String(), p.pluginName).Observe(duration)
	}()

	return p.wrapped.Report(ctx, timestamp, query, observations)
}

func (p *promPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promShouldAcceptFinalizedReport.WithLabelValues(p.evmChainID.String(), p.pluginName).Observe(duration)
	}()

	return p.wrapped.ShouldAcceptFinalizedReport(ctx, timestamp, report)
}

func (p *promPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promShouldTransmitAcceptedReport.WithLabelValues(p.evmChainID.String(), p.pluginName).Observe(duration)
	}()

	return p.wrapped.ShouldTransmitAcceptedReport(ctx, timestamp, report)
}

func (p *promPlugin) Close() error {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		promClose.WithLabelValues(p.evmChainID.String(), p.pluginName).Observe(duration)
	}()

	return p.wrapped.Close()
}
