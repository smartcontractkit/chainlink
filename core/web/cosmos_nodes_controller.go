package web

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrCosmosNotEnabled is returned when COSMOS_ENABLED is not true.
var ErrCosmosNotEnabled = errChainDisabled{name: "Cosmos", envVar: "COSMOS_ENABLED"}

func NewCosmosNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.CosmosNodeResource](
		app.GetChains().Cosmos, ErrCosmosNotEnabled, parse, presenters.NewCosmosNodeResource, app.GetAuditLogger(),
	)
}
