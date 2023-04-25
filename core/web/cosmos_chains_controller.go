package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewCosmosChainsController(app chainlink.Application) ChainsController {
	parse := func(s string) (string, error) { return s, nil }
	return newChainsController[string, presenters.CosmosChainResource](
		"cosmos", app.GetChains().Cosmos, ErrCosmosNotEnabled, parse, presenters.NewCosmosChainResource, app.GetLogger(), app.GetAuditLogger())
}
