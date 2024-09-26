package cache

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var (
	laneHealthStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_lane_healthcheck_status",
		Help: "Keep track of the chain healthcheck calls for each lane and plugin",
	}, []string{"plugin", "source", "dest", "onramp"})
)

type ObservedChainHealthcheck struct {
	ChainHealthcheck

	sourceChain string
	destChain   string
	plugin      string
	// onrampAddress is used to distinguish between 1.0/2.0 lanes or blue/green lanes during deployment
	// This changes very rarely, so it's not a performance concern for Prometheus
	onrampAddress    string
	laneHealthStatus *prometheus.GaugeVec
}

func NewObservedChainHealthCheck(
	chainHealthcheck ChainHealthcheck,
	plugin string,
	sourceChain int64,
	destChain int64,
	onrampAddress cciptypes.Address,
) *ObservedChainHealthcheck {
	return &ObservedChainHealthcheck{
		ChainHealthcheck: chainHealthcheck,
		sourceChain:      strconv.FormatInt(sourceChain, 10),
		destChain:        strconv.FormatInt(destChain, 10),
		plugin:           plugin,
		laneHealthStatus: laneHealthStatus,
		onrampAddress:    string(onrampAddress),
	}
}

func (o *ObservedChainHealthcheck) IsHealthy(ctx context.Context) (bool, error) {
	healthy, err := o.ChainHealthcheck.IsHealthy(ctx)
	o.trackState(healthy, err)
	return healthy, err
}

func (o *ObservedChainHealthcheck) trackState(healthy bool, err error) {
	if err != nil {
		// Don't report errors as unhealthy, as they are not necessarily indicative of the chain's health
		// Could be RPC issues, etc.
		return
	}

	status := 0
	if healthy {
		status = 1
	}

	o.laneHealthStatus.
		WithLabelValues(o.plugin, o.sourceChain, o.destChain, o.onrampAddress).
		Set(float64(status))
}
