package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	chainSet := app.GetChains().Solana
	return newChainsController("solana", chainSet, ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource, app.GetLogger(), app.GetAuditLogger())
}
