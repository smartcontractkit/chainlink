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
		float64(200 * time.Millisecond),
		float64(300 * time.Millisecond),
		float64(400 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(750 * time.Millisecond),
		float64(1 * time.Second),
		float64(2 * time.Second),
		float64(3 * time.Second),
		float64(4 * time.Second),
	}
	labels          = []string{"evmChainID", "plugin", "reader", "function", "success"}
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
	pluginName          string
	readerName          string
	chainId             int64
}

func withObservedInteraction[T any](metric metricDetails, function string, f func() (T, error)) (T, error) {
	contractExecutionStarted := time.Now()
	value, err := f()
	metric.interactionDuration.
		WithLabelValues(
			strconv.FormatInt(metric.chainId, 10),
			metric.pluginName,
			metric.readerName,
			function,
			strconv.FormatBool(err == nil),
		).
		Observe(float64(time.Since(contractExecutionStarted)))
	return value, err
}

func withObservedInteractionAndResults[T any](metric metricDetails, function string, f func() ([]T, error)) ([]T, error) {
	results, err := withObservedInteraction(metric, function, f)
	if err == nil {
		metric.resultSetSize.WithLabelValues(
			strconv.FormatInt(metric.chainId, 10),
			metric.pluginName,
			metric.readerName,
			function,
			strconv.FormatBool(err == nil),
		).Set(float64(len(results)))
	}
	return results, err
}
