package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewCosmosKeysController(app chainlink.Application) KeysController {
	return NewKeysController[cosmoskey.Key, presenters.CosmosKeyResource](app.GetKeyStore().Cosmos(), app.GetLogger(), app.GetAuditLogger(),
		"cosmosKey", presenters.NewCosmosKeyResource, presenters.NewCosmosKeyResources)
}
