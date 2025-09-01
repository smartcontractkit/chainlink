package p2p

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type SharedPeerMetrics struct {
	// Measure current number of PeerGroups. On non-bootstrap nodes, the number of discovery groups
	// should be at least equal to number of remote DONs this node is expected to connect with.
	// The number of messaging groups should be equal to number of remote nodes this node is expected to connect with.
	discoveryGroups metric.Int64Gauge
	messagingGroups metric.Int64Gauge

	// Failure count and latency of a peer group update call.
	groupUpdateFailureCounter metric.Int64Counter
	groupUpdateDurationMs     metric.Int64Histogram
}

func initSharedPeerMetrics() (m *SharedPeerMetrics, err error) {
	m = &SharedPeerMetrics{}
	m.discoveryGroups, err = beholder.GetMeter().Int64Gauge("platform_don2don_discovery_groups")
	if err != nil {
		return nil, fmt.Errorf("failed to register platform_don2don_discovery_groups gauge: %w", err)
	}
	m.messagingGroups, err = beholder.GetMeter().Int64Gauge("platform_don2don_messaging_groups")
	if err != nil {
		return nil, fmt.Errorf("failed to register platform_don2don_messaging_groups gauge: %w", err)
	}
	m.groupUpdateFailureCounter, err = beholder.GetMeter().Int64Counter("platform_don2don_group_update_failure_total")
	if err != nil {
		return nil, fmt.Errorf("failed to register platform_don2don_group_update_failure_total counter: %w", err)
	}
	m.groupUpdateDurationMs, err = beholder.GetMeter().Int64Histogram(
		"platform_don2don_group_update_duration_ms",
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to register platform_don2don_group_update_duration_ms histogram: %w", err)
	}
	return m, nil
}
