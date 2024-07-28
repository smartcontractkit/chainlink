package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewAptosKeysController(app chainlink.Application) KeysController {
	return NewKeysController[aptoskey.Key, presenters.AptosKeyResource](app.GetKeyStore().Aptos(), app.GetLogger(), app.GetAuditLogger(),
		"aptosKey", presenters.NewAptosKeyResource, presenters.NewAptosKeyResources)
}
