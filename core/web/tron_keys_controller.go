package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewTronKeysController(app chainlink.Application) KeysController {
	return NewKeysController[tronkey.Key, presenters.TronKeyResource](app.GetKeyStore().Tron(), app.GetLogger(), app.GetAuditLogger(),
		"tronKey", presenters.NewTronKeyResource, presenters.NewTronKeyResources)
}
