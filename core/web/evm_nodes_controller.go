package web

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewEVMNodesController(app chainlink.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), types.NetworkEVM)

	return newNodesController[presenters.EVMNodeResource](
		scopedNodeStatuser, ErrEVMNotEnabled, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
