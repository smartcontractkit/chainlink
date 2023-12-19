package ccip

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type skipReason string

const (
	// ReasonNotBlessed describes when a report is skipped due to not being blessed.
	ReasonNotBlessed skipReason = "not blessed"

	// ReasonAllExecuted describes when a report is skipped due to messages being all executed.
	ReasonAllExecuted skipReason = "all executed"
)

var (
	execPluginLabels          = []string{"configDigest"}
	execPluginDurationBuckets = []float64{
		float64(10 * time.Millisecond),
		float64(20 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(1 * time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
	}
	metricReportSkipped = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ccip_unexpired_report_skipped",
		Help: "Times report is skipped for the possible reasons",
	}, []string{"reason"})
	execPluginReportsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_execution_observation_reports_count",
		Help: "Number of reports that are being processed by Execution Plugin during single observation",
	}, execPluginLabels)
	execPluginObservationBuildDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_execution_observation_build_duration",
		Help:    "Duration of generating Observation in Execution Plugin",
		Buckets: execPluginDurationBuckets,
	}, execPluginLabels)
	execPluginBatchBuildDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_execution_build_single_batch",
		Help:    "Duration of building single batch in Execution Plugin",
		Buckets: execPluginDurationBuckets,
	}, execPluginLabels)

	//nolint unused
	execPluginReportsIterationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_execution_reports_iteration_build_batch",
		Help:    "Duration of iterating over all unexpired reports in Execution Plugin",
		Buckets: execPluginDurationBuckets,
	}, execPluginLabels)
)

func measureExecPluginDuration(histogram *prometheus.HistogramVec, timestamp types.ReportTimestamp, duration time.Duration) {
	histogram.
		WithLabelValues(timestampToLabels(timestamp)...).
		Observe(float64(duration))
}

func MeasureObservationBuildDuration(timestamp types.ReportTimestamp, duration time.Duration) {
	measureExecPluginDuration(execPluginObservationBuildDuration, timestamp, duration)
}

func MeasureBatchBuildDuration(timestamp types.ReportTimestamp, duration time.Duration) {
	measureExecPluginDuration(execPluginBatchBuildDuration, timestamp, duration)
}

// nolint unused
func measureReportsIterationDuration(timestamp types.ReportTimestamp, duration time.Duration) {
	measureExecPluginDuration(execPluginReportsIterationDuration, timestamp, duration)
}

func MeasureNumberOfReportsProcessed(timestamp types.ReportTimestamp, count int) {
	execPluginReportsCount.
		WithLabelValues(timestampToLabels(timestamp)...).
		Set(float64(count))
}

func IncSkippedRequests(reason skipReason) {
	metricReportSkipped.WithLabelValues(string(reason)).Inc()
}

func timestampToLabels(t types.ReportTimestamp) []string {
	return []string{t.ConfigDigest.Hex()}
}
