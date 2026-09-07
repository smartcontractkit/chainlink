package web

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/tonkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewTONKeysController(app chainlink.Application) KeysController {
	return NewKeysController[tonkey.Key, presenters.TONKeyResource](app.GetKeyStore().TON(), app.GetLogger(), app.GetAuditLogger(),
		"tonKey", presenters.NewTONKeyResource, presenters.NewTONKeyResources)
}
