package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewTronChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		relay.NetworkTron,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.NetworkTron)),
		ErrTronNotEnabled,
		presenters.NewTronChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
