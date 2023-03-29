package web

import (
	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStarknetKeysController(app chainlink.Application) KeysController {
	return NewKeysController[starkkey.Key, presenters.StarknetKeyResource](app.GetKeyStore().Starknet(), app.GetLogger(), app.GetAuditLogger(),
		"starknetKey", presenters.NewStarknetKeyResource, presenters.NewStarknetKeyResources)
}
