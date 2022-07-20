package web

import (
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewStarkNetKeysController(app chainlink.Application) KeysController {
	return NewKeysController[starkkey.Key, presenters.StarkNetKeyResource](app.GetKeyStore().StarkNet(), app.GetLogger(),
		"starknetKey", presenters.NewStarkNetKeyResource, presenters.NewStarkNetKeyResources)
}
