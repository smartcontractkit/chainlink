package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStarkNetChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		relay.NetworkStarkNet,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.NetworkStarkNet)),
		ErrStarkNetNotEnabled,
		presenters.NewStarkNetChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
