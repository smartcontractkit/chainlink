package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewCosmosChainsController(app chainlink.Application) ChainsController {
	return newChainsController[presenters.CosmosChainResource](
		relay.Cosmos,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.Cosmos)),
		ErrCosmosNotEnabled,
		presenters.NewCosmosChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
