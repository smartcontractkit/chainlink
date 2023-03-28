package web

import (
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewEVMNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (id utils.Big, err error) {
		err = id.UnmarshalText([]byte(s))
		return
	}
	return newNodesController[utils.Big, types.Node, presenters.EVMNodeResource](
		app.GetChains().EVM, ErrEVMNotEnabled, parse, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
