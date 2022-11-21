package web

import (
	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewStarkNetKeysController(app chainlink.Application) KeysController {
	return NewKeysController[starkkey.Key, presenters.StarkNetKeyResource](app.GetKeyStore().StarkNet(), app.GetLogger(), app.GetAuditLogger(),
		"starknetKey", presenters.NewStarkNetKeyResource, presenters.NewStarkNetKeyResources)
}
