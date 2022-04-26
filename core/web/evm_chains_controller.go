package web

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

var ErrEVMNotEnabled = errors.New("EVM is disabled. Set EVM_ENABLED=true to enable.")

func NewEVMChainsController(app chainlink.Application) ChainsController {
	return NewChainsController[utils.Big, types.ChainCfg, presenters.EVMChainResource](
		"evm", app.GetChains().EVM, ErrEVMNotEnabled, func(s string) (id utils.Big, err error) {
			err = id.UnmarshalText([]byte(s))
			return
		}, presenters.NewEVMChainResource)
}
