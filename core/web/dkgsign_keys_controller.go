package web

import (
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewDKGSignKeysController(app chainlink.Application) KeysController {
	return NewKeysController[dkgsignkey.Key, presenters.DKGSignKeyResource](
		app.GetKeyStore().DKGSign(),
		app.GetLogger(),
		app.GetAuditLogger(),
		"dkgsignKey",
		presenters.NewDKGSignKeyResource,
		presenters.NewDKGSignKeyResources)
}
