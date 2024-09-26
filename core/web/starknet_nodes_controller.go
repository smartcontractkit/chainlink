package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrStarkNetNotEnabled is returned when Starknet.Enabled is not true.
var ErrStarkNetNotEnabled = errChainDisabled{name: "StarkNet", tomlKey: "Starknet.Enabled"}

func NewStarkNetNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkStarkNet)

	return newNodesController[presenters.StarkNetNodeResource](
		scopedNodeStatuser, ErrStarkNetNotEnabled, presenters.NewStarkNetNodeResource, app.GetAuditLogger())
}
