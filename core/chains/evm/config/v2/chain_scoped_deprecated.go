package v2

import (
	"fmt"
	"math/big"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
)

func (c *ChainScoped) SetEvmGasPriceDefault(_ *big.Int) error {
	panic(fmt.Errorf("cannot reconfigure gas price: %v", config.ErrUnsupported))
}

func (c *ChainScoped) Configure(_ evmtypes.ChainCfg) {
	panic(fmt.Errorf("cannot reconfigure chain: %v", config.ErrUnsupported))
}

func (c *ChainScoped) PersistedConfig() evmtypes.ChainCfg {
	panic(fmt.Errorf("cannot get persisted config: %v", config.ErrUnsupported))
}
