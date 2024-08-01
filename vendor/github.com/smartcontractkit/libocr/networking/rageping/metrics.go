package rageping

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type latencyMetrics struct {
	// Be careful about the field order here and during initialization.
	registerer prometheus.Registerer

	// Captures the round-trip time of a PING/PONG message pair between this host and the remote peer.
	roundTripLatencySeconds prometheus.Histogram

	// Counts the number of outgoing PING messages sent to the remote peer.
	sentRequestsTotal prometheus.Counter

	// Counts the number of valid incoming PING messages received from the remote peer.
	receivedRequestsTotal prometheus.Counter

	// Counts the number of PING messages for which no valid PONG message was received in time from the remote peer.
	timedOutRequestsTotal prometheus.Counter

	// Counts the number of other invalid messages received from the remote peer.
	// An invalid message could be of invalid size, have an invalid message type, or be late PONG message.
	invalidMessagesReceivedTotal prometheus.Counter
}

func newLatencyMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	peerID ragetypes.PeerID,
	remotePeerID ragetypes.PeerID,
	config *LatencyMetricsServiceConfig,
) *latencyMetrics {
	constLabels := prometheus.Labels{
		"peer_id":        peerID.String(),
		"remote_peer_id": remotePeerID.String(),
		"ping_size":      fmt.Sprint(config.PingSize),
		"min_period":     fmt.Sprint(config.MinPeriod),
		"max_period":     fmt.Sprint(config.MaxPeriod),
		"timeout":        fmt.Sprint(config.Timeout),
	}

	roundTripLatencyBuckets := config.Buckets
	if config.Buckets == nil {
		roundTripLatencyBuckets = DefaultLatencyBuckets()
	}
	roundTripLatencySeconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        "rageping_round_trip_latency_seconds",
		Help:        "The round trip latency, i.e., the time between sending a PING to receiving the corresponding PONG.",
		ConstLabels: constLabels,
		Buckets:     roundTripLatencyBuckets,
	})
	metricshelper.RegisterOrLogError(logger, registerer, roundTripLatencySeconds, "rageping_round_trip_latency_seconds")

	sentRequestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "rageping_sent_requests_total",
		Help:        "The number of PING requests sent to the remote peer.",
		ConstLabels: constLabels,
	})
	metricshelper.RegisterOrLogError(logger, registerer, sentRequestsTotal, "rageping_sent_requests_total")

	receivedRequestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "rageping_received_requests_total",
		Help:        "The number of PING requests received from the remote remote peer.",
		ConstLabels: constLabels,
	})
	metricshelper.RegisterOrLogError(logger, registerer, receivedRequestsTotal, "rageping_received_requests_total")

	timedOutRequestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "rageping_timed_out_requests_total",
		Help:        "The number of PING requests sent to the remote peer, for which no (valid) response was received before the configured timeout.",
		ConstLabels: constLabels,
	})
	metricshelper.RegisterOrLogError(logger, registerer, timedOutRequestsTotal, "rageping_timed_out_requests_total")

	invalidMessagesReceivedTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "rageping_invalid_messages_received_total",
		Help:        "The number of invalid messages received from the remote peer. These are expected rarely due to restarts of the underlying network connection.",
		ConstLabels: constLabels,
	})
	metricshelper.RegisterOrLogError(logger, registerer, invalidMessagesReceivedTotal, "rageping_invalid_messages_received_total")

	// Be careful about the initialization order. peer_test.go tests that the metrics below are indeed exposed and
	// updated. Be sure to update peer_test.go when changing the metric names here.
	return &latencyMetrics{
		registerer,
		roundTripLatencySeconds,
		sentRequestsTotal,
		receivedRequestsTotal,
		timedOutRequestsTotal,
		invalidMessagesReceivedTotal,
	}
}

func (m *latencyMetrics) Close() {
	m.registerer.Unregister(m.roundTripLatencySeconds)
	m.registerer.Unregister(m.sentRequestsTotal)
	m.registerer.Unregister(m.receivedRequestsTotal)
	m.registerer.Unregister(m.timedOutRequestsTotal)
	m.registerer.Unregister(m.invalidMessagesReceivedTotal)
}
