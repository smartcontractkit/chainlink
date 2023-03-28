package web

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when Solana.Enabled is not true.
var ErrSolanaNotEnabled = errChainDisabled{name: "Solana", envVar: "Solana.Enabled"}

func NewSolanaNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.SolanaNodeResource](
		app.GetChains().Solana, ErrSolanaNotEnabled, parse, presenters.NewSolanaNodeResource, app.GetAuditLogger())
}
