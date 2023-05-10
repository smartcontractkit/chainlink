package wsrpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type reqStatus string

const (
	statusSuccess reqStatus = "success"
	statusFailed  reqStatus = "failed"
)

var (
	aliveMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "mercury",
		Name:      "wsrpc_connection_alive",
		Help:      "Total time spent connected to the WSRPC server",
	})
	requestsStatusMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "mercury",
		Name:      "wsrpc_requests_status_count",
		Help:      "Number of request status made to the WSRPC server",
	}, []string{"status"})

	requestLatencyMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "mercury",
		Name:      "wsrpc_request_latency",
		Help:      "Latency of requests made to the WSRPC server",
	})
)

func setLivenessMetric(live bool) {
	if live {
		aliveMetric.Set(1)
	} else {
		aliveMetric.Set(0)
	}
}

func incRequestStatusMetric(status reqStatus) {
	requestsStatusMetric.WithLabelValues(string(status)).Inc()
}

func setRequestLatencyMetric(latency float64) {
	requestLatencyMetric.Set(latency)
}
