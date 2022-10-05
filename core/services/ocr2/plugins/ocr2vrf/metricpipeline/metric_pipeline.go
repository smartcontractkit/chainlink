package metricpipeline

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metriPipeline struct {
	jobName                 string
	externalJobID           string
	stageDurationMetric     *prometheus.HistogramVec
	transmissionCountMetric *prometheus.CounterVec
}

// NewMetricPipeline provides a pipeline by which OCR job metrics can be reported.
func NewMetricPipeline(
	initialBucket time.Duration,
	bucketSize time.Duration,
	bucketCount int64,
	bucketsOverride []float64,
	jobName string,
	externalJobID string,
	metricPrefix string,
) *metriPipeline {

	// Construct buckets.
	var buckets []float64
	for i := int64(0); i < bucketCount; i++ {
		buckets = append(buckets, float64(i*int64(bucketSize)+int64(initialBucket)))
	}

	// Use override for buckets if given.
	if len(bucketsOverride) > 0 {
		buckets = bucketsOverride
	}

	// Construct a histogram metric to log stage duration results.
	stageDurationMetric := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    fmt.Sprintf("%s_ocr_stage_duration", metricPrefix),
		Help:    "Duration of the local execution for an OCR stage.",
		Buckets: buckets,
	}, []string{"job_name", "external_job_id"})

	// Construct a counter metric to log transmissions.
	transmissionCountMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_ocr_transmission_count", metricPrefix),
		Help: "Counter for each OCR transmission.",
	}, []string{"job_name", "external_job_id"})

	return &metriPipeline{
		jobName:                 jobName,
		externalJobID:           externalJobID,
		stageDurationMetric:     stageDurationMetric,
		transmissionCountMetric: transmissionCountMetric,
	}
}

// ReportStageDuration reports the duration of the local execution for an OCR stage.
func (m *metriPipeline) ReportStageDuration(duration float64) {
	m.stageDurationMetric.WithLabelValues(m.jobName, m.externalJobID).Observe(duration)
}

// ReportTransmission reports a new transmission event for an OCR job.
func (m *metriPipeline) ReportTransmission() {
	m.transmissionCountMetric.WithLabelValues(m.jobName, m.externalJobID).Inc()
}
