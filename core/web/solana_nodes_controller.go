package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when Solana.Enabled is not true.
var ErrSolanaNotEnabled = errChainDisabled{name: "Solana", tomlKey: "Solana.Enabled"}

func NewSolanaNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkSolana)

	return newNodesController[presenters.SolanaNodeResource](
		scopedNodeStatuser, ErrSolanaNotEnabled, presenters.NewSolanaNodeResource, app.GetAuditLogger())
}
