package promfm

import (
	"math/big"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
