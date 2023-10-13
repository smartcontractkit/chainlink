package observability

import (
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
		float64(2 * time.Second),
	}
	labels                 = []string{"evmChainID", "plugin", "function", "success"}
	priceRegistryHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_price_registry_contract_duration",
		Help:    "Duration of calls to the Price Registry reader",
		Buckets: latencyBuckets,
	}, labels)
	commitStoreHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_commit_store_contract_duration",
		Help:    "Duration of calls to the Commit Store reader",
		Buckets: latencyBuckets,
	}, labels)
	onRampHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_onramp_contract_duration",
		Help:    "Duration of calls to the OnRamp reader",
		Buckets: latencyBuckets,
	}, labels)
	offRampHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_offramp_contract_duration",
		Help:    "Duration of calls to the OffRamp contract",
		Buckets: latencyBuckets,
	}, labels)
)

type metricDetails struct {
	histogram  *prometheus.HistogramVec
	pluginName string
	chainId    int64
}

func withObservedContract[T any](metric metricDetails, function string, contract func() (T, error)) (T, error) {
	contractExecutionStarted := time.Now()
	value, err := contract()
	metric.histogram.
		WithLabelValues(
			strconv.FormatInt(metric.chainId, 10),
			metric.pluginName,
			function,
			strconv.FormatBool(err == nil),
		).
		Observe(float64(time.Since(contractExecutionStarted)))
	return value, err
}
