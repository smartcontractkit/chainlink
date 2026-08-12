package web

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStellarKeysController(app chainlink.Application) KeysController {
	return NewKeysController[stellarkey.Key, presenters.StellarKeyResource](app.GetKeyStore().Stellar(), app.GetLogger(), app.GetAuditLogger(),
		"stellarKey", presenters.NewStellarKeyResource, presenters.NewStellarKeyResources)
}
