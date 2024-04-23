package web

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewStarkNetChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		types.NetworkStarkNet,
		app.GetRelayers().List(chainlink.FilterRelayersByType(types.NetworkStarkNet)),
		ErrStarkNetNotEnabled,
		presenters.NewStarkNetChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
