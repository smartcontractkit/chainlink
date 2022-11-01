package v2

import (
	"fmt"
	"math/big"

	config "github.com/smartcontractkit/chainlink/core/config/v2"
)

func (c *ChainScoped) SetEvmGasPriceDefault(_ *big.Int) error {
	panic(fmt.Errorf("cannot reconfigure gas price: %v", config.ErrUnsupported))
}
