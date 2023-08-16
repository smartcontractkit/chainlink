package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStarkNetChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		relay.StarkNet,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.StarkNet)),
		ErrStarkNetNotEnabled,
		presenters.NewStarkNetChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
