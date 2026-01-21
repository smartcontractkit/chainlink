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

	// currentShardCount tracks the current number of shards observed.
	currentShardCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_current_shard_count",
			Help: "Current number of shards",
		},
	)

	// desiredShardCount tracks the number of shards KEDA wants.
	desiredShardCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_desired_shard_count",
			Help: "Desired number of shards",
		},
	)

	// approvedShardCount tracks the number of shards the Arbiter approved.
	approvedShardCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_approved_shard_count",
			Help: "Approved number of shards",
		},
	)

	// onChainShardNumber tracks the on-chain governance limit.
	onChainShardNumber = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "arbiter_onchain_shard_number",
			Help: "On-chain shard number from ShardConfig contract",
		},
	)
)

// RecordRequest increments the request counter for the given endpoint and status.
func RecordRequest(endpoint, status string) {
	requestsTotal.WithLabelValues(endpoint, status).Inc()
}

// SetCurrentShardCount sets the current shard count gauge.
func SetCurrentShardCount(count int) {
	currentShardCount.Set(float64(count))
}

// SetDesiredShardCount sets the desired shard count gauge.
func SetDesiredShardCount(count int) {
	desiredShardCount.Set(float64(count))
}

// SetApprovedShardCount sets the approved shard count gauge.
func SetApprovedShardCount(count int) {
	approvedShardCount.Set(float64(count))
}

// SetOnChainShardNumber sets the on-chain shard number gauge.
func SetOnChainShardNumber(count uint64) {
	onChainShardNumber.Set(float64(count))
}
