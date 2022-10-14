package web

import (
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

var ErrEVMNotEnabled = errChainDisabled{name: "EVM", envVar: "EVM_ENABLED"}

func NewEVMChainsController(app chainlink.Application) ChainsController {
	parse := func(s string) (id utils.Big, err error) {
		err = id.UnmarshalText([]byte(s))
		return
	}
	return newChainsController[utils.Big, *types.ChainCfg, presenters.EVMChainResource](
		"evm", app.GetChains().EVM, ErrEVMNotEnabled, parse, presenters.NewEVMChainResource, app.GetLogger(), app.GetAuditLogger())
}
