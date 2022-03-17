package vrf

import (
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
