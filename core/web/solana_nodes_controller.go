package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when Solana.Enabled is not true.
var ErrSolanaNotEnabled = errChainDisabled{name: "Solana", tomlKey: "Solana.Enabled"}

func NewSolanaNodesController(app chainlink.Application) NodesController {
	return newNodesController[presenters.SolanaNodeResource](
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.Solana)), ErrSolanaNotEnabled, presenters.NewSolanaNodeResource, app.GetAuditLogger())
}
