package web

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

var ErrEVMNotEnabled = errChainDisabled{name: "EVM", tomlKey: "EVM.Enabled"}

func NewEVMChainsController(app chainlink.Application) ChainsController {
	return newChainsController[presenters.EVMChainResource](
		"evm", app.GetChains().EVM, ErrEVMNotEnabled, presenters.NewEVMChainResource, app.GetLogger(), app.GetAuditLogger())
}
