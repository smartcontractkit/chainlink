package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewEVMNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (id utils.Big, err error) {
		err = id.UnmarshalText([]byte(s))
		return
	}
	return newNodesController[utils.Big, types.Node, presenters.EVMNodeResource](
		app.GetChains().EVM, ErrEVMNotEnabled, parse, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
