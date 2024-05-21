package web

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewCosmosChainsController(app chainlink.Application) ChainsController {
	return newChainsController[presenters.CosmosChainResource](
		types.NetworkCosmos,
		app.GetRelayers().List(chainlink.FilterRelayersByType(types.NetworkCosmos)),
		ErrCosmosNotEnabled,
		presenters.NewCosmosChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
