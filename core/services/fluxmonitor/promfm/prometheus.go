package promfm

import (
	"math/big"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shopspring/decimal"
)

var (
	ReportedValue = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "flux_monitor_reported_value",
			Help: "Flux monitor's last reported price",
		},
		[]string{"job_spec_id"},
	)

	IndividualReportedValue = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "flux_monitor_individual_reported_value",
			Help: "Flux monitor's last reported price for each individual endpoint",
		},
		[]string{"url"},
	)

	SeenValue = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "flux_monitor_seen_value",
			Help: "Flux monitor's last observed value from target",
		},
		[]string{"job_spec_id"},
	)

	ReportedRound = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "flux_monitor_reported_round",
			Help: "Flux monitor's last reported round",
		},
		[]string{"job_spec_id"},
	)

	SeenRound = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "flux_monitor_seen_round",
			Help: "Last seen round by other node operators",
		},
		[]string{"job_spec_id"},
	)

	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "flux_monitor_request_duration_seconds",
			Help:    "Flux monitor's histogram of request latencies",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	ResponseSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "flux_monitor_response_size_bytes",
			Help:    "Flux monitor's last response body size",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// SetDecimal sets a decimal metric
func SetDecimal(gauge prometheus.Gauge, arg decimal.Decimal) {
	val, _ := arg.Float64()
	gauge.Set(val)
}

// SetBigInt sets a big.Int metric
func SetBigInt(gauge prometheus.Gauge, arg *big.Int) {
	gauge.Set(float64(arg.Int64()))
}

// SetUint32 sets a uint32 metric
func SetUint32(gauge prometheus.Gauge, arg uint32) {
	gauge.Set(float64(arg))
}

func InstrumentRoundTripperReponseSize(
	obs prometheus.Observer,
	next http.RoundTripper,
) promhttp.RoundTripperFunc {
	return promhttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(r)
		if err == nil && resp.ContentLength >= 0 {
			obs.Observe(float64(resp.ContentLength))
		}
		return resp, err
	})
}
