package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrTronNotEnabled is returned when Starknet.Enabled is not true.
var ErrTronNotEnabled = errChainDisabled{name: "Tron", tomlKey: "Tron.Enabled"}

func NewTronNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkTron)

	return newNodesController[presenters.TronNodeResource](
		scopedNodeStatuser, ErrTronNotEnabled, presenters.NewTronNodeResource, app.GetAuditLogger())
}
