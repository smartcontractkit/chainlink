package observability

import (
	"math/big"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	latencyBuckets = []float64{
		float64(10 * time.Millisecond),
		float64(25 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(75 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(250 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(750 * time.Millisecond),
		float64(1 * time.Second),
	}
	labels                 = []string{"evmChainID", "plugin", "function", "success"}
	priceRegistryHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_price_registry_contract_rpc_duration",
		Help:    "Duration of RPC calls to the Price Registry contract",
		Buckets: latencyBuckets,
	}, labels)
	commitStoreHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_commit_store_contract_rpc_duration",
		Help:    "Duration of RPC calls to the Commit Store contract",
		Buckets: latencyBuckets,
	}, labels)
	evm2evmOnRampHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_evm2evm_onramp_contract_rpc_duration",
		Help:    "Duration of RPC calls to the EVM2EVMOnRamp contract",
		Buckets: latencyBuckets,
	}, labels)
	evm2evmOffRampHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_evm2evm_offramp_contract_rpc_duration",
		Help:    "Duration of RPC calls to the EVM2EVMOffRamp contract",
		Buckets: latencyBuckets,
	}, labels)
)

type metricDetails struct {
	histogram  *prometheus.HistogramVec
	pluginName string
	chainId    *big.Int
}

func withObservedContract[T any](metric metricDetails, function string, contract func() (T, error)) (T, error) {
	contractExecutionStarted := time.Now()
	value, err := contract()
	metric.histogram.
		WithLabelValues(
			metric.chainId.String(),
			metric.pluginName,
			function,
			strconv.FormatBool(err == nil),
		).
		Observe(float64(time.Since(contractExecutionStarted)))
	return value, err
}
