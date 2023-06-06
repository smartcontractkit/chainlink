package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewEVMNodesController(app chainlink.Application) NodesController {
	return newNodesController[presenters.EVMNodeResource](
		app.GetChains().EVM, ErrEVMNotEnabled, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
