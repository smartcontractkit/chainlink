package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	return newChainsController(
		relay.NetworkSolana,
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.NetworkSolana)),
		ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
