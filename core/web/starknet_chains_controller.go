package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStarkNetChainsController(app chainlink.Application) ChainsController {
	return newChainsController("starknet", app.GetChains().StarkNet, ErrStarkNetNotEnabled,
		presenters.NewStarkNetChainResource, app.GetLogger(), app.GetAuditLogger())
}
