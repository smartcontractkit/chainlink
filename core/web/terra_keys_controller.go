package web

import (
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewTerraKeysController(app chainlink.Application) KeysController {
	return NewKeysController[terrakey.Key, presenters.TerraKeyResource](app.GetKeyStore().Terra(), app.GetLogger(), app.GetAuditLogger(),
		"terraKey", presenters.NewTerraKeyResource, presenters.NewTerraKeyResources)
}
