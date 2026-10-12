package standardcapabilities

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// capabilityDonIDUnresolved is 1 when the host could not resolve the capability
// DON ID for a capability at startup, so its events carry no authoritative DON ID.
var capabilityDonIDUnresolved = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "standard_capabilities_capability_don_id_unresolved",
		Help: "1 if the capability DON ID could not be resolved at startup, 0 otherwise",
	},
	[]string{"capability_id"},
)
