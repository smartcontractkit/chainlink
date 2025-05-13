package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewTonKeysController(app chainlink.Application) KeysController {
	return NewKeysController[tonkey.Key, presenters.TonKeyResource](app.GetKeyStore().Ton(), app.GetLogger(), app.GetAuditLogger(),
		"tonKey", presenters.NewTonKeyResource, presenters.NewTonKeyResources)
}
