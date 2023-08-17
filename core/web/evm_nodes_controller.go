package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewEVMNodesController(app chainlink.Application) NodesController {
	return newNodesController[presenters.EVMNodeResource](
		app.GetRelayers().List(chainlink.FilterRelayersByType(relay.EVM)), ErrEVMNotEnabled, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
