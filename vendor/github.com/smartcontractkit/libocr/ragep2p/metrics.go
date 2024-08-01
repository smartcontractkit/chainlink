package ragep2p

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/ragep2p/types"
)

type hostMetrics struct {
	registerer        prometheus.Registerer
	inboundDialsTotal prometheus.Counter
}

func newHostMetrics(registerer prometheus.Registerer, logger commontypes.Logger, self types.PeerID) *hostMetrics {
	labels := map[string]string{"peer_id": self.String()}

	inboundDialsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "ragep2p_host_inbound_dials_total",
		Help:        "The number of inbound dial attempts received by the host",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, inboundDialsTotal, "ragep2p_host_inbound_dials_total")

	return &hostMetrics{
		registerer,
		inboundDialsTotal,
	}
}

func (m *hostMetrics) Close() {
	m.registerer.Unregister(m.inboundDialsTotal)
}

type peerMetrics struct {
	registerer                  prometheus.Registerer
	connEstablishedTotal        prometheus.Counter
	connEstablishedInboundTotal prometheus.Counter
	connReadProcessedBytesTotal prometheus.Counter
	connReadSkippedBytesTotal   prometheus.Counter
	connWrittenBytesTotal       prometheus.Counter
	rawconnReadBytesTotal       prometheus.Counter
	rawconnWrittenBytesTotal    prometheus.Counter
	rawconnRateLimitRate        prometheus.Gauge
	rawconnRateLimitCapacity    prometheus.Gauge
	messageBytes                prometheus.Histogram
}

func newPeerMetrics(registerer prometheus.Registerer, logger commontypes.Logger, self types.PeerID, other types.PeerID) *peerMetrics {
	labels := map[string]string{"peer_id": self.String(), "remote_peer_id": other.String()}

	connEstablishedTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "ragep2p_peer_conn_established_total",
		Help:        "The number of secure connections established with the remote peer. At most one connection can be active at any time.",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, connEstablishedTotal, "ragep2p_peer_conn_established_total")

	connEstablishedInboundTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "ragep2p_peer_conn_established_inbound_total",
		Help:        "The number of secure connections established with the remote peer from inbound dials. At most one connection can be active at any time.",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, connEstablishedInboundTotal, "ragep2p_peer_conn_established_inbound_total")

	connReadProcessedBytesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "ragep2p_peer_conn_read_processed_bytes_total",
		Help:        "The number of bytes read on secure connections with the remote peer for processing",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, connReadProcessedBytesTotal, "ragep2p_peer_conn_read_processed_bytes_total")

	connReadSkippedBytesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "ragep2p_peer_conn_read_skipped_bytes_total",
		Help:        "The number of bytes read on secure connections with the remote peer that have been skipped, e.g. due to rate limits being exceeded",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, connReadSkippedBytesTotal, "ragep2p_peer_conn_read_skipped_bytes_total")

	connWrittenBytesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "ragep2p_peer_conn_written_bytes_total",
		Help:        "The number of bytes written on secure connections with the remote peer",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, connWrittenBytesTotal, "ragep2p_peer_conn_written_bytes_total")

	rawconnReadBytesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ragep2p_peer_rawconn_read_bytes_total",
		Help: "The number of raw bytes read on raw (post-knock, tcp) connections with the remote " +
			"peer. Knocks are ~100 bytes and thus have negligible impact. This metric is useful for " +
			"tracking overall bandwidth usage.",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, rawconnReadBytesTotal, "ragep2p_peer_rawconn_read_bytes_total")

	rawconnWrittenBytesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ragep2p_peer_rawconn_written_bytes_total",
		Help: "The number of raw bytes written on raw (post-knock, tcp) connections with the remote " +
			"peer. Knocks are ~100 bytes and thus have negligible impact. This metric is useful for " +
			"tracking overall bandwidth usage.",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, rawconnWrittenBytesTotal, "ragep2p_peer_rawconn_written_bytes_total")

	rawconnRateLimitRate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragep2p_peer_rawconn_rate_limit_rate",
		Help:        "The refill rate in bytes per second for the token bucket rate limiting reads from raw connections with the remote peer",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, rawconnRateLimitRate, "ragep2p_peer_rawconn_rate_limit_rate")

	rawconnRateLimitCapacity := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragep2p_peer_rawconn_rate_limit_capacity",
		Help:        "The capacity in bytes for the token bucket rate limiting reads from raw connections with the remote peer",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, rawconnRateLimitCapacity, "ragep2p_peer_rawconn_rate_limit_capacity")

	messageBytes := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        "ragep2p_experimental_peer_message_bytes",
		Help:        "The size of messages sent to the remote peer",
		ConstLabels: labels,
		Buckets:     []float64{1 << 8, 1 << 10, 1 << 12, 1 << 14, 1 << 16, 1 << 18, 1 << 20, 1 << 22, 1 << 24, 1 << 26, 1 << 28}, // 256 bytes, 1KiB, ..., 256MiB
	})

	metricshelper.RegisterOrLogError(logger, registerer, messageBytes, "ragep2p_experimental_peer_message_bytes")

	return &peerMetrics{
		registerer,
		connEstablishedTotal,
		connEstablishedInboundTotal,
		connReadProcessedBytesTotal,
		connReadSkippedBytesTotal,
		connWrittenBytesTotal,
		rawconnReadBytesTotal,
		rawconnWrittenBytesTotal,
		rawconnRateLimitRate,
		rawconnRateLimitCapacity,
		messageBytes,
	}
}

func (m *peerMetrics) Close() {
	m.registerer.Unregister(m.connEstablishedTotal)
	m.registerer.Unregister(m.connEstablishedInboundTotal)
	m.registerer.Unregister(m.connReadProcessedBytesTotal)
	m.registerer.Unregister(m.connReadSkippedBytesTotal)
	m.registerer.Unregister(m.connWrittenBytesTotal)
	m.registerer.Unregister(m.rawconnReadBytesTotal)
	m.registerer.Unregister(m.rawconnWrittenBytesTotal)
	m.registerer.Unregister(m.rawconnRateLimitRate)
	m.registerer.Unregister(m.rawconnRateLimitCapacity)
	m.registerer.Unregister(m.messageBytes)
}

func (m *peerMetrics) SetConnRateLimit(tokenBucketParams TokenBucketParams) {
	m.rawconnRateLimitRate.Set(tokenBucketParams.Rate)
	m.rawconnRateLimitCapacity.Set(float64(tokenBucketParams.Capacity))
}
