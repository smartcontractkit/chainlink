package web

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewCosmosChainsController(app chainlink.Application) ChainsController {
	parse := func(s string) (string, error) { return s, nil }
	return newChainsController[string, *db.ChainCfg, presenters.CosmosChainResource](
		"cosmos", app.GetChains().Cosmos, ErrCosmosNotEnabled, parse, presenters.NewCosmosChainResource, app.GetLogger(), app.GetAuditLogger())
}
