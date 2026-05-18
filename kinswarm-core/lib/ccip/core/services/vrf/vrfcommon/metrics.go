package vrfcommon

import (
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// version describes a VRF version.
type Version string

const (
	V1     Version = "V1"
	V2     Version = "V2"
	V2Plus Version = "V2Plus"
)

// dropReason describes a reason why a VRF request is dropped from the queue.
type dropReason string

const (
	// ReasonMailboxSize describes when a VRF request is dropped due to the log mailbox being
	// over capacity.
	ReasonMailboxSize dropReason = "mailbox_size"

	// ReasonAge describes when a VRF request is dropped due to its age.
	ReasonAge dropReason = "age"
)

var (
	MetricQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vrf_request_queue_size",
		Help: "The number of VRF requests currently in the in-memory queue.",
	}, []string{"job_name", "external_job_id", "vrf_version"})

	MetricProcessedReqs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vrf_processed_request_count",
		Help: "The number of VRF requests processed.",
	}, []string{"job_name", "external_job_id", "vrf_version"})

	MetricDroppedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vrf_dropped_request_count",
		Help: "The number of VRF requests dropped due to reasons such as expiry or mailbox size.",
	}, []string{"job_name", "external_job_id", "vrf_version", "drop_reason"})

	MetricDupeRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vrf_duplicate_requests",
		Help: "The number of times the VRF listener receives duplicate requests, which could indicate a reorg.",
	}, []string{"job_name", "external_job_id", "vrf_version"})

	MetricTimeBetweenSims = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "vrf_request_time_between_sims",
		Help: "How long a VRF request sits in the in-memory queue in between simulation attempts.",
		Buckets: []float64{
			float64(time.Second),
			float64(30 * time.Second),
			float64(time.Minute),
			float64(2 * time.Minute),
			float64(5 * time.Minute),
		},
	}, []string{"job_name", "external_job_id", "vrf_version"})

	MetricTimeUntilInitialSim = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "vrf_request_time_until_initial_sim",
		Help: "How long a VRF request sits in the in-memory queue until it gets simulated for the first time.",
		Buckets: []float64{
			float64(time.Second),
			float64(30 * time.Second),
			float64(time.Minute),
			float64(2 * time.Minute),
			float64(5 * time.Minute),
		},
	}, []string{"job_name", "external_job_id", "vrf_version"})
)

func UpdateQueueSize(jobName string, extJobID uuid.UUID, vrfVersion Version, size int) {
	MetricQueueSize.WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).
		Set(float64(size))
}

func IncProcessedReqs(jobName string, extJobID uuid.UUID, vrfVersion Version) {
	MetricProcessedReqs.WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).Inc()
}

func IncDroppedReqs(jobName string, extJobID uuid.UUID, vrfVersion Version, reason dropReason) {
	MetricDroppedRequests.WithLabelValues(
		jobName, extJobID.String(), string(vrfVersion), string(reason)).Inc()
}

func IncDupeReqs(jobName string, extJobID uuid.UUID, vrfVersion Version) {
	MetricDupeRequests.WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).Inc()
}
