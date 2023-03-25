package web

import (
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewCosmosKeysController(app chainlink.Application) KeysController {
	return NewKeysController[cosmoskey.Key, presenters.CosmosKeyResource](app.GetKeyStore().Cosmos(), app.GetLogger(), app.GetAuditLogger(),
		"cosmosKey", presenters.NewCosmosKeyResource, presenters.NewCosmosKeyResources)
}
