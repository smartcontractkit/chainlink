package arbiter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// requestsTotal counts all gRPC requests by endpoint and status.
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "arbiter_requests_total",
			Help: "Total number of requests by endpoint and status",
		},
		[]string{"endpoint", "status"},
	)

	// currentReplicas tracks the current number of replicas observed.
	currentReplicas = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_current_replicas",
			Help: "Current number of replicas",
		},
	)

	// desiredReplicas tracks the number of replicas KEDA wants.
	desiredReplicas = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_desired_replicas",
			Help: "Desired number of replicas",
		},
	)

	// approvedReplicas tracks the number of replicas the Arbiter approved.
	approvedReplicas = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_approved_replicas",
			Help: "Approved number of replicas",
		},
	)

	// onChainMaxReplicas tracks the on-chain governance limit.
	onChainMaxReplicas = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_onchain_max_replicas",
			Help: "On-chain maximum replicas from ShardConfig contract",
		},
	)
)

// RecordRequest increments the request counter for the given endpoint and status.
func RecordRequest(endpoint, status string) {
	requestsTotal.WithLabelValues(endpoint, status).Inc()
}

// SetCurrentReplicas sets the current replica count gauge.
func SetCurrentReplicas(count int) {
	currentReplicas.Set(float64(count))
}

// SetDesiredReplicas sets the desired replica count gauge.
func SetDesiredReplicas(count int) {
	desiredReplicas.Set(float64(count))
}

// SetApprovedReplicas sets the approved replica count gauge.
func SetApprovedReplicas(count int) {
	approvedReplicas.Set(float64(count))
}

// SetOnChainMaxReplicas sets the on-chain max replica count gauge.
func SetOnChainMaxReplicas(count uint64) {
	onChainMaxReplicas.Set(float64(count))
}
