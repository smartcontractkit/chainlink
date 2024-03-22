// promwrapper wraps another OCR2 reporting plugin and provides standardized prometheus metrics
// for each of the OCR2 phases (Query, Observation, Report, ShouldAcceptFinalizedReport,
// ShouldTransmitAcceptedReport, and Close).
package promwrapper

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Type assertions, buckets and labels.
var (
	_       types.ReportingPlugin = &promPlugin{}
	_       PrometheusBackend     = &defaultPrometheusBackend{}
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
	labels          = []string{"chainType", "chainID", "plugin", "oracleID", "configDigest"}
	getLabelsValues = func(p *promPlugin, t types.ReportTimestamp) []string {
		return []string{
			p.chainType,                         // chainType
			p.chainID.String(),                  // chainID
			p.name,                              // plugin
			p.oracleID,                          // oracleID
			common.Bytes2Hex(t.ConfigDigest[:]), // configDigest
		}
	}
)

// Prometheus queries.
var (
	promQuery = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_query_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Query() method",
			Buckets: buckets,
		},
		labels,
	)
	promObservation = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_observation_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Observation() method",
			Buckets: buckets,
		},
		labels,
	)
	promReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_report_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Report() method",
			Buckets: buckets,
		},
		labels,
	)
	promShouldAcceptFinalizedReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_should_accept_finalized_report_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's ShouldAcceptFinalizedReport() method",
			Buckets: buckets,
		},
		labels,
	)
	promShouldTransmitAcceptedReport = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_should_transmit_accepted_report_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's ShouldTransmitAcceptedReport() method",
			Buckets: buckets,
		},
		labels,
	)
	promClose = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_close_time",
			Help:    "The amount of time elapsed during the OCR2 plugin's Close() method",
			Buckets: buckets,
		},
		[]string{"chainType", "chainID", "plugin", "oracleID", "configDigest"},
	)
	promQueryToObservationLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_end_query_to_begin_observation",
			Help:    "The amount of time elapsed after the OCR2 node's Query() method and before its Observation() method",
			Buckets: buckets,
		},
		labels,
	)
	promObservationToReportLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_end_observation_to_begin_report_time",
			Help:    "The amount of time elapsed after the OCR2 node's Observation() method and before its Report() method",
			Buckets: buckets,
		},
		labels,
	)
	promReportToAcceptFinalizedReportLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_end_report_to_begin_accept_finalized_report",
			Help:    "The amount of time elapsed after the OCR2 node's Report() method and before its ShouldAcceptFinalizedReport() method",
			Buckets: buckets,
		},
		labels,
	)
	promAcceptFinalizedReportToTransmitAcceptedReportLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocr2_reporting_plugin_end_accept_finalized_report_to_begin_transmit_accepted_report",
			Help:    "The amount of time elapsed after the OCR2 node's ShouldAcceptFinalizedReport() method and before its ShouldTransmitAcceptedReport() method",
			Buckets: buckets,
		},
		labels,
	)
)

//go:generate mockery --quiet --name PrometheusBackend --output ./mocks/ --case=underscore
type (
	// Contains interface for logging OCR telemetry.
	PrometheusBackend interface {
		// Intra-phase latency.
		SetQueryDuration([]string, float64)
		SetObservationDuration([]string, float64)
		SetReportDuration([]string, float64)
		SetShouldAcceptFinalizedReportDuration([]string, float64)
		SetShouldTransmitAcceptedReportDuration([]string, float64)
		SetCloseDuration([]string, float64)

		// Inter-phase latency.
		SetQueryToObservationLatency([]string, float64)
		SetObservationToReportLatency([]string, float64)
		SetReportToAcceptFinalizedReportLatency([]string, float64)
		SetAcceptFinalizedReportToTransmitAcceptedReportLatency([]string, float64)
	}

	defaultPrometheusBackend struct{} // implements PrometheusBackend

	// promPlugin consumes a report plugin and wraps its core functions e.g Report(), Observe()...
	promPlugin struct {
		wrapped                       types.ReportingPlugin
		name                          string
		chainType                     string
		chainID                       *big.Int
		oracleID                      string
		configDigest                  string
		queryEndTimes                 sync.Map
		observationEndTimes           sync.Map
		reportEndTimes                sync.Map
		acceptFinalizedReportEndTimes sync.Map
		prometheusBackend             PrometheusBackend
	}
)

func (*defaultPrometheusBackend) SetQueryDuration(labelValues []string, duration float64) {
	promQuery.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPrometheusBackend) SetObservationDuration(labelValues []string, duration float64) {
	promObservation.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPrometheusBackend) SetReportDuration(labelValues []string, duration float64) {
	promReport.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPrometheusBackend) SetShouldAcceptFinalizedReportDuration(labelValues []string, duration float64) {
	promShouldAcceptFinalizedReport.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPrometheusBackend) SetShouldTransmitAcceptedReportDuration(labelValues []string, duration float64) {
	promShouldTransmitAcceptedReport.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPrometheusBackend) SetCloseDuration(labelValues []string, duration float64) {
	promClose.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPrometheusBackend) SetQueryToObservationLatency(labelValues []string, latency float64) {
	promQueryToObservationLatency.WithLabelValues(labelValues...).Observe(latency)
}

func (*defaultPrometheusBackend) SetObservationToReportLatency(labelValues []string, latency float64) {
	promObservationToReportLatency.WithLabelValues(labelValues...).Observe(latency)
}

func (*defaultPrometheusBackend) SetReportToAcceptFinalizedReportLatency(labelValues []string, latency float64) {
	promReportToAcceptFinalizedReportLatency.WithLabelValues(labelValues...).Observe(latency)
}

func (*defaultPrometheusBackend) SetAcceptFinalizedReportToTransmitAcceptedReportLatency(labelValues []string, latency float64) {
	promAcceptFinalizedReportToTransmitAcceptedReportLatency.WithLabelValues(labelValues...).Observe(latency)
}

func New(
	plugin types.ReportingPlugin,
	name string,
	chainType string,
	chainID *big.Int,
	config types.ReportingPluginConfig,
	backend PrometheusBackend,
) types.ReportingPlugin {
	// Apply passed-in Prometheus backend if one is given.
	var prometheusBackend PrometheusBackend = &defaultPrometheusBackend{}
	if backend != nil {
		prometheusBackend = backend
	}

	return &promPlugin{
		wrapped:           plugin,
		name:              name,
		chainType:         chainType,
		chainID:           chainID,
		oracleID:          fmt.Sprintf("%d", config.OracleID),
		configDigest:      common.Bytes2Hex(config.ConfigDigest[:]),
		prometheusBackend: prometheusBackend,
	}
}

func (p *promPlugin) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetQueryDuration(getLabelsValues(p, timestamp), duration)
		p.queryEndTimes.Store(timestamp, time.Now().UTC()) // note time at end of Query()
	}()

	return p.wrapped.Query(ctx, timestamp)
}

func (p *promPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	start := time.Now().UTC()

	// Report latency between Query() and Observation().
	labelValues := getLabelsValues(p, timestamp)
	if queryEndTime, ok := p.queryEndTimes.Load(timestamp); ok {
		latency := float64(start.Sub(queryEndTime.(time.Time)))
		p.prometheusBackend.SetQueryToObservationLatency(labelValues, latency)
		p.queryEndTimes.Delete(timestamp)
	}

	// Report latency for Observation() at end of call.
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetObservationDuration(labelValues, duration)
		p.observationEndTimes.Store(timestamp, time.Now().UTC()) // note time at end of Observe()
	}()

	return p.wrapped.Observation(ctx, timestamp, query)
}

func (p *promPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	start := time.Now().UTC()

	// Report latency between Observation() and Report().
	labelValues := getLabelsValues(p, timestamp)
	if observationEndTime, ok := p.observationEndTimes.Load(timestamp); ok {
		latency := float64(start.Sub(observationEndTime.(time.Time)))
		p.prometheusBackend.SetObservationToReportLatency(labelValues, latency)
		p.observationEndTimes.Delete(timestamp)
	}

	// Report latency for Report() at end of call.
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetReportDuration(labelValues, duration)
		p.reportEndTimes.Store(timestamp, time.Now().UTC()) // note time at end of Report()
	}()

	return p.wrapped.Report(ctx, timestamp, query, observations)
}

func (p *promPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()

	// Report latency between Report() and ShouldAcceptFinalizedReport().
	labelValues := getLabelsValues(p, timestamp)
	if reportEndTime, ok := p.reportEndTimes.Load(timestamp); ok {
		latency := float64(start.Sub(reportEndTime.(time.Time)))
		p.prometheusBackend.SetReportToAcceptFinalizedReportLatency(labelValues, latency)
		p.reportEndTimes.Delete(timestamp)
	}

	// Report latency for ShouldAcceptFinalizedReport() at end of call.
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetShouldAcceptFinalizedReportDuration(labelValues, duration)
		p.acceptFinalizedReportEndTimes.Store(timestamp, time.Now().UTC()) // note time at end of ShouldAcceptFinalizedReport()
	}()

	return p.wrapped.ShouldAcceptFinalizedReport(ctx, timestamp, report)
}

func (p *promPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()

	// Report latency between ShouldAcceptFinalizedReport() and ShouldTransmitAcceptedReport().
	labelValues := getLabelsValues(p, timestamp)
	if acceptFinalizedReportEndTime, ok := p.acceptFinalizedReportEndTimes.Load(timestamp); ok {
		latency := float64(start.Sub(acceptFinalizedReportEndTime.(time.Time)))
		p.prometheusBackend.SetAcceptFinalizedReportToTransmitAcceptedReportLatency(labelValues, latency)
		p.acceptFinalizedReportEndTimes.Delete(timestamp)
	}

	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetShouldTransmitAcceptedReportDuration(labelValues, duration)
	}()

	return p.wrapped.ShouldTransmitAcceptedReport(ctx, timestamp, report)
}

// Note: the 'Close' method does not have access to a report timestamp, as it is not part of report generation.
func (p *promPlugin) Close() error {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		labelValues := []string{
			p.chainType,        // chainType
			p.chainID.String(), // chainID
			p.name,             // plugin
			p.oracleID,         // oracleID
			p.configDigest,     // configDigest
		}
		p.prometheusBackend.SetCloseDuration(labelValues, duration)
	}()

	return p.wrapped.Close()
}
