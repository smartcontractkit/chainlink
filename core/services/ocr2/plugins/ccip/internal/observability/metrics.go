package observability

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	latencyBuckets = []float64{
		float64(30 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(300 * time.Millisecond),
		float64(1 * time.Second),
		float64(3 * time.Second),
	}
	labels          = []string{"evmChainID", "reader", "function"}
	readerHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_reader_duration",
		Help:    "Duration of calls to Reader instance",
		Buckets: latencyBuckets,
	}, labels)
	readerDatasetSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_reader_dataset_size",
		Help: "Size of the dataset returned from the Reader instance",
	}, labels)
)

type metricDetails struct {
	interactionDuration *prometheus.HistogramVec
	resultSetSize       *prometheus.GaugeVec
	readerName          string
	chainId             int64
}

func withObservedInteraction[T any](metric metricDetails, function string, f func() (T, error)) (T, error) {
	contractExecutionStarted := time.Now()
	value, err := f()
	metric.interactionDuration.
		WithLabelValues(
			strconv.FormatInt(metric.chainId, 10),
			metric.readerName,
			function,
		).
		Observe(float64(time.Since(contractExecutionStarted)))
	return value, err
}

func withObservedInteractionAndResults[T any](metric metricDetails, function string, f func() ([]T, error)) ([]T, error) {
	results, err := withObservedInteraction(metric, function, f)
	if err == nil {
		metric.resultSetSize.WithLabelValues(
			strconv.FormatInt(metric.chainId, 10),
			metric.readerName,
			function,
		).Set(float64(len(results)))
	}
	return results, err
}
