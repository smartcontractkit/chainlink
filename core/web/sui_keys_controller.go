package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewSuiKeysController(app chainlink.Application) KeysController {
	return NewKeysController[suikey.Key, presenters.SuiKeyResource](app.GetKeyStore().Sui(), app.GetLogger(), app.GetAuditLogger(),
		"suiKey", presenters.NewSuiKeyResource, presenters.NewSuiKeyResources)
}
