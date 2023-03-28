package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	return newChainsController[string]("solana", app.GetChains().Solana, ErrSolanaNotEnabled,
		func(s string) (string, error) { return s, nil }, presenters.NewSolanaChainResource, app.GetLogger(), app.GetAuditLogger())
}
