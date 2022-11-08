package vrf

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
)

// version describes a VRF version.
type version string

const (
	v1 version = "v1"
	v2 version = "v2"
)

// dropReason describes a reason why a VRF request is dropped from the queue.
type dropReason string

const (
	// reasonMailboxSize describes when a VRF request is dropped due to the log mailbox being
	// over capacity.
	reasonMailboxSize dropReason = "mailbox_size"

	// reasonAge describes when a VRF request is dropped due to its age.
	reasonAge dropReason = "age"
)

var (
	metricQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vrf_request_queue_size",
		Help: "The number of VRF requests currently in the in-memory queue.",
	}, []string{"job_name", "external_job_id", "vrf_version"})

	metricProcessedReqs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vrf_processed_request_count",
		Help: "The number of VRF requests processed.",
	}, []string{"job_name", "external_job_id", "vrf_version"})

	metricDroppedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vrf_dropped_request_count",
		Help: "The number of VRF requests dropped due to reasons such as expiry or mailbox size.",
	}, []string{"job_name", "external_job_id", "vrf_version", "drop_reason"})

	metricDupeRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vrf_duplicate_requests",
		Help: "The number of times the VRF listener receives duplicate requests, which could indicate a reorg.",
	}, []string{"job_name", "external_job_id", "vrf_version"})

	metricTimeBetweenSims = promauto.NewHistogramVec(prometheus.HistogramOpts{
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

	metricTimeUntilInitialSim = promauto.NewHistogramVec(prometheus.HistogramOpts{
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

func updateQueueSize(jobName string, extJobID uuid.UUID, vrfVersion version, size int) {
	metricQueueSize.WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).
		Set(float64(size))
}

func incProcessedReqs(jobName string, extJobID uuid.UUID, vrfVersion version) {
	metricProcessedReqs.WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).Inc()
}

func incDroppedReqs(jobName string, extJobID uuid.UUID, vrfVersion version, reason dropReason) {
	metricDroppedRequests.WithLabelValues(
		jobName, extJobID.String(), string(vrfVersion), string(reason)).Inc()
}

func incDupeReqs(jobName string, extJobID uuid.UUID, vrfVersion version) {
	metricDupeRequests.WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).Inc()
}

// observeRequestSimDuration records the time between the given requests simulations or
// the time until it's first simulation, whichever is applicable.
// Cases:
// 1. Never simulated: in this case, we want to observe the time until simulated
// on the utcTimestamp field of the pending request.
// 2. Simulated before: in this case, lastTry will be set to a non-zero time value,
// in which case we'd want to use that as a relative point from when we last tried
// the request.
func observeRequestSimDuration(jobName string, extJobID uuid.UUID, vrfVersion version, pendingReqs []pendingRequest) {
	now := time.Now().UTC()
	for _, request := range pendingReqs {
		// First time around lastTry will be zero because the request has not been
		// simulated yet. It will be updated every time the request is simulated (in the event
		// the request is simulated multiple times, due to it being underfunded).
		if request.lastTry.IsZero() {
			metricTimeUntilInitialSim.
				WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).
				Observe(float64(now.Sub(request.utcTimestamp)))
		} else {
			metricTimeBetweenSims.
				WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).
				Observe(float64(now.Sub(request.lastTry)))
		}
	}
}
