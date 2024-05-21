package web

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		types.NetworkSolana,
		app.GetRelayers().List(chainlink.FilterRelayersByType(types.NetworkSolana)),
		ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
