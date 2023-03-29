package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStarknetChainsController(app chainlink.Application) ChainsController {
	return newChainsController[string]("starknet", app.GetChains().Starknet, ErrStarknetNotEnabled,
		func(s string) (string, error) { return s, nil }, presenters.NewStarknetChainResource, app.GetLogger(), app.GetAuditLogger())
}
