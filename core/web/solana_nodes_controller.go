package web

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when SOLANA_ENABLED is not true.
var ErrSolanaNotEnabled = errChainDisabled{name: "Solana", envVar: "SOLANA_ENABLED"}

func NewSolanaNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	_ = app.GetChains().SolanaRelayer //TODO support dynamic toml https://smartcontract-it.atlassian.net/browse/BCF-2114
	return newNodesController[string, db.Node, presenters.SolanaNodeResource](
		nil, ErrSolanaNotEnabled, parse, presenters.NewSolanaNodeResource, app.GetAuditLogger())
}
