package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrCosmosNotEnabled is returned when COSMOS_ENABLED is not true.
var ErrCosmosNotEnabled = errChainDisabled{name: "Cosmos", tomlKey: "Cosmos.Enabled"}

func NewCosmosNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkCosmos)

	return newNodesController[presenters.CosmosNodeResource](
		scopedNodeStatuser, ErrCosmosNotEnabled, presenters.NewCosmosNodeResource, app.GetAuditLogger(),
	)
}
