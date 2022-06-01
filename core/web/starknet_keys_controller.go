package web

import (
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewStarknetKeysController(app chainlink.Application) KeysController {
	return NewKeysController[starkkey.Key, presenters.StarknetKeyResource](app.GetKeyStore().Starknet(), app.GetLogger(),
		"starknetKey", presenters.NewStarknetKeyResource, presenters.NewStarknetKeyResources)
}
