package web

import (
	solkey "github.com/smartcontractkit/chainlink-solana/pkg/solana/keys"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewSolanaKeysController(app chainlink.Application) KeysController {
	return NewKeysController[solkey.Key, presenters.SolanaKeyResource](app.GetKeyStore().Solana(), app.GetLogger(), app.GetAuditLogger(),
		"solanaKey", presenters.NewSolanaKeyResource, presenters.NewSolanaKeyResources)
}
