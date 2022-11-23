package web

import (
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewDKGEncryptKeysController(app chainlink.Application) KeysController {
	return NewKeysController[dkgencryptkey.Key, presenters.DKGEncryptKeyResource](
		app.GetKeyStore().DKGEncrypt(),
		app.GetLogger(),
		app.GetAuditLogger(),
		"dkgencryptKey",
		presenters.NewDKGEncryptKeyResource,
		presenters.NewDKGEncryptKeyResources)
}
