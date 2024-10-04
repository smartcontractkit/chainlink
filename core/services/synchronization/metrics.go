package synchronization

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TelemetryClientConnectionStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "telemetry_client_connection_status",
		Help: "Status of the connection to the telemetry ingress server",
	}, []string{"endpoint"})
)

var (
	TelemetryClientMessagesSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "telemetry_client_messages_sent",
		Help: "Number of telemetry messages sent to the telemetry ingress server",
	}, []string{"endpoint"})
)

var (
	TelemetryClientMessagesDropped = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "telemetry_client_messages_dropped",
		Help: "Number of telemetry messages dropped",
	}, []string{"endpoint", "worker_name"})
)

var (
	TelemetryClientWorkers = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "telemetry_client_workers",
		Help: "Number of telemetry workers",
	}, []string{"endpoint"})
)
