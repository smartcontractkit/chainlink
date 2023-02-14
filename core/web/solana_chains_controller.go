package web

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	_ = app.GetChains().SolanaRelayer //TODO support dynamic toml
	return newChainsController[string, *db.ChainCfg]("solana", nil, ErrSolanaNotEnabled,
		func(s string) (string, error) { return s, nil }, presenters.NewSolanaChainResource, app.GetLogger(), app.GetAuditLogger())
}
