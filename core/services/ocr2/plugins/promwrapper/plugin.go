// promwrapper wraps another OCR2 reporting plugin and provides standardized prometheus metrics
// for each of the OCR2 phases (Query, Observation, Report, ShouldAcceptFinalizedReport,
// ShouldTransmitAcceptedReport, and Close).
package promwrapper

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Type assertions, buckets and labels.
var (
	_       types.ReportingPlugin = &promPlugin{}
	_       PrometheusBackend     = &defaultPometheusBackend{}
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
	labels          = []string{"chainType", "chainID", "plugin", "oracleID", "configDigest", "epoch", "round"}
	getLabelsValues = func(p *promPlugin, t types.ReportTimestamp) []string {
		return []string{
			string(p.chainType),                 // chainType
			p.chainID.String(),                  // chainID
			p.name,                              // plugin
			p.oracleID,                          // oracleID
			common.Bytes2Hex(t.ConfigDigest[:]), // configDigest
			fmt.Sprintf("%d", t.Epoch),          // epoch
			fmt.Sprintf("%d", t.Round),          // round
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

	defaultPometheusBackend struct{} // implements PrometheusBackend

	// promPlugin consumes a report plugin and wraps its core functions e.g Report(), Observe()...
	promPlugin struct {
		wrapped                       types.ReportingPlugin
		name                          string
		chainType                     string
		chainID                       *big.Int
		oracleID                      string
		configDigest                  string
		queryEndTimes                 map[types.ReportTimestamp]time.Time
		observationEndTimes           map[types.ReportTimestamp]time.Time
		reportEndTimes                map[types.ReportTimestamp]time.Time
		acceptFinalizedReportEndTimes map[types.ReportTimestamp]time.Time
		prometheusBackend             PrometheusBackend
	}
)

func (*defaultPometheusBackend) SetQueryDuration(labelValues []string, duration float64) {
	promQuery.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPometheusBackend) SetObservationDuration(labelValues []string, duration float64) {
	promObservation.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPometheusBackend) SetReportDuration(labelValues []string, duration float64) {
	promReport.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPometheusBackend) SetShouldAcceptFinalizedReportDuration(labelValues []string, duration float64) {
	promShouldAcceptFinalizedReport.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPometheusBackend) SetShouldTransmitAcceptedReportDuration(labelValues []string, duration float64) {
	promShouldTransmitAcceptedReport.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPometheusBackend) SetCloseDuration(labelValues []string, duration float64) {
	promClose.WithLabelValues(labelValues...).Observe(duration)
}

func (*defaultPometheusBackend) SetQueryToObservationLatency(labelValues []string, latency float64) {
	promQueryToObservationLatency.WithLabelValues(labelValues...).Observe(latency)
}

func (*defaultPometheusBackend) SetObservationToReportLatency(labelValues []string, latency float64) {
	promObservationToReportLatency.WithLabelValues(labelValues...).Observe(latency)
}

func (*defaultPometheusBackend) SetReportToAcceptFinalizedReportLatency(labelValues []string, latency float64) {
	promReportToAcceptFinalizedReportLatency.WithLabelValues(labelValues...).Observe(latency)
}

func (*defaultPometheusBackend) SetAcceptFinalizedReportToTransmitAcceptedReportLatency(labelValues []string, latency float64) {
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
	var prometheusBackend PrometheusBackend = &defaultPometheusBackend{}
	if backend != nil {
		prometheusBackend = backend
	}

	return &promPlugin{
		wrapped:                       plugin,
		name:                          name,
		chainType:                     chainType,
		chainID:                       chainID,
		oracleID:                      fmt.Sprintf("%d", config.OracleID),
		configDigest:                  common.Bytes2Hex(config.ConfigDigest[:]),
		queryEndTimes:                 make(map[types.ReportTimestamp]time.Time),
		observationEndTimes:           make(map[types.ReportTimestamp]time.Time),
		reportEndTimes:                make(map[types.ReportTimestamp]time.Time),
		acceptFinalizedReportEndTimes: make(map[types.ReportTimestamp]time.Time),
		prometheusBackend:             prometheusBackend,
	}
}

func (p *promPlugin) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	start := time.Now().UTC()
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetQueryDuration(getLabelsValues(p, timestamp), duration)
		p.queryEndTimes[timestamp] = time.Now() // note time at end of Query()
	}()

	return p.wrapped.Query(ctx, timestamp)
}

func (p *promPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	start := time.Now().UTC()

	// Report latency between Query() and Observation().
	labelValues := getLabelsValues(p, timestamp)
	if queryEndTime, ok := p.queryEndTimes[timestamp]; ok {
		latency := float64(start.Sub(queryEndTime))
		p.prometheusBackend.SetQueryToObservationLatency(labelValues, latency)
	}

	// Report latency for Observation() at end of call.
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetObservationDuration(labelValues, duration)
		p.observationEndTimes[timestamp] = time.Now() // note time at end of Observe()
	}()

	return p.wrapped.Observation(ctx, timestamp, query)
}

func (p *promPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	start := time.Now().UTC()

	// Report latency between Observation() and Report().
	labelValues := getLabelsValues(p, timestamp)
	if observationEndTime, ok := p.observationEndTimes[timestamp]; ok {
		latency := float64(start.Sub(observationEndTime))
		p.prometheusBackend.SetObservationToReportLatency(labelValues, latency)
	}

	// Report latency for Report() at end of call.
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetReportDuration(labelValues, duration)
		p.reportEndTimes[timestamp] = time.Now() // note time at end of Report()
	}()

	return p.wrapped.Report(ctx, timestamp, query, observations)
}

func (p *promPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()

	// Report latency between Report() and ShouldAcceptFinalizedReport().
	labelValues := getLabelsValues(p, timestamp)
	if reportEndTime, ok := p.reportEndTimes[timestamp]; ok {
		latency := float64(start.Sub(reportEndTime))
		p.prometheusBackend.SetReportToAcceptFinalizedReportLatency(labelValues, latency)
	}

	// Report latency for ShouldAcceptFinalizedReport() at end of call.
	defer func() {
		duration := float64(time.Now().UTC().Sub(start))
		p.prometheusBackend.SetShouldAcceptFinalizedReportDuration(labelValues, duration)
		p.acceptFinalizedReportEndTimes[timestamp] = time.Now() // note time at end of ShouldAcceptFinalizedReport()
	}()

	return p.wrapped.ShouldAcceptFinalizedReport(ctx, timestamp, report)
}

func (p *promPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	start := time.Now().UTC()

	// Report latency between ShouldAcceptFinalizedReport() and ShouldTransmitAcceptedReport().
	labelValues := getLabelsValues(p, timestamp)
	if acceptFinalizedReportEndTime, ok := p.acceptFinalizedReportEndTimes[timestamp]; ok {
		latency := float64(start.Sub(acceptFinalizedReportEndTime))
		p.prometheusBackend.SetAcceptFinalizedReportToTransmitAcceptedReportLatency(labelValues, latency)
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
			string(p.chainType), // chainType
			p.chainID.String(),  // chainID
			p.name,              // plugin
			p.oracleID,          // oracleID
			p.configDigest,      // configDigest
		}
		p.prometheusBackend.SetCloseDuration(labelValues, duration)
	}()

	return p.wrapped.Close()
}
