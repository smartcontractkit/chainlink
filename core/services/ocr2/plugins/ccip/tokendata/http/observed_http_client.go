package http

import (
	"context"
	"net/http"
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
		float64(3 * time.Second),
		float64(4 * time.Second),
		float64(5 * time.Second),
	}
	usdcClientHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_usdc_client_request_total",
		Help:    "Latency of calls to the USDC client",
		Buckets: latencyBuckets,
	}, []string{"status", "success"})
	lbtcClientHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_lbtc_client_request_total",
		Help:    "Latency of calls to the LBTC client",
		Buckets: latencyBuckets,
	}, []string{"status", "success"})
)

type ObservedIHttpClient struct {
	IHttpClient
	histogram *prometheus.HistogramVec
}

// NewObservedUsdcIHttpClient Create a new ObservedIHttpClient with the USDC client metric.
func NewObservedUsdcIHttpClient(origin IHttpClient) *ObservedIHttpClient {
	return NewObservedIHttpClientWithMetric(origin, usdcClientHistogram)
}

// NewObservedLbtcIHttpClient Create a new ObservedIHttpClient with the LBTC client metric.
func NewObservedLbtcIHttpClient(origin IHttpClient) *ObservedIHttpClient {
	return NewObservedIHttpClientWithMetric(origin, lbtcClientHistogram)
}

func NewObservedIHttpClientWithMetric(origin IHttpClient, histogram *prometheus.HistogramVec) *ObservedIHttpClient {
	return &ObservedIHttpClient{
		IHttpClient: origin,
		histogram:   histogram,
	}
}

func (o *ObservedIHttpClient) Get(ctx context.Context, url string, timeout time.Duration) ([]byte, int, http.Header, error) {
	return withObservedHttpClient(o.histogram, func() ([]byte, int, http.Header, error) {
		return o.IHttpClient.Get(ctx, url, timeout)
	})
}

func withObservedHttpClient[T any](histogram *prometheus.HistogramVec, contract func() (T, int, http.Header, error)) (T, int, http.Header, error) {
	contractExecutionStarted := time.Now()
	value, status, headers, err := contract()
	histogram.
		WithLabelValues(
			strconv.FormatInt(int64(status), 10),
			strconv.FormatBool(err == nil),
		).
		Observe(float64(time.Since(contractExecutionStarted)))
	return value, status, headers, err
}
