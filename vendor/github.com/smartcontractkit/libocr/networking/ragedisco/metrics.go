package ragedisco

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/ragep2p/types"
)

type discoveryProtocolMetrics struct {
	registerer      prometheus.Registerer
	registeredPeers prometheus.Gauge
	discoveredPeers prometheus.Gauge
	bootstrappers   prometheus.Gauge
}

func newDiscoveryProtocolMetrics(registerer prometheus.Registerer, logger commontypes.Logger, peerID types.PeerID) *discoveryProtocolMetrics {
	labels := map[string]string{"peer_id": peerID.String()}

	registeredPeers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragedisco_registered_peers",
		Help:        "The number of registered peers in peer discovery",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, registeredPeers, "ragedisco_registered_peers")

	discoveredPeers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ragedisco_discovered_peers",
		Help: "The number of discovered peers in peer discovery. A peer " +
			"is considered discovered if we have obtained a valid announcement for it. " +
			"Note that a valid announcement may still contain an incorrect address.",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, discoveredPeers, "ragedisco_discovered_peers")

	bootstrappers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragedisco_bootstappers",
		Help:        "The number of bootstrappers in peer discovery",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, bootstrappers, "ragedisco_bootstappers")

	return &discoveryProtocolMetrics{
		registerer,
		registeredPeers,
		discoveredPeers,
		bootstrappers,
	}
}

func (dpm *discoveryProtocolMetrics) Close() {
	dpm.registerer.Unregister(dpm.registeredPeers)
	dpm.registerer.Unregister(dpm.bootstrappers)
	dpm.registerer.Unregister(dpm.discoveredPeers)
}
