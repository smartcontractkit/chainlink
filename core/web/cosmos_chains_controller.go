package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewCosmosChainsController(app chainlink.Application) ChainsController {
	return newChainsController[presenters.CosmosChainResource](
		"cosmos", app.GetChains().Cosmos, ErrCosmosNotEnabled, presenters.NewCosmosChainResource, app.GetLogger(), app.GetAuditLogger())
}
