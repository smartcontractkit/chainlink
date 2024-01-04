package liquiditymanager

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// Factory initializes a new liquidity manager instance.
//
//go:generate mockery --quiet --name Factory --output ../rebalancermocks --filename lm_factory_mock.go --case=underscore
type Factory interface {
	NewLiquidityManager(networkID models.NetworkID, address models.Address) (LiquidityManager, error)
}

type BaseLiquidityManagerFactory struct{}

func NewBaseLiquidityManagerFactory() *BaseLiquidityManagerFactory {
	return &BaseLiquidityManagerFactory{}
}

func (b *BaseLiquidityManagerFactory) NewLiquidityManager(networkID models.NetworkID, address models.Address) (LiquidityManager, error) {
	switch typ := networkID.Type(); typ {
	case models.NetworkTypeEvm:
		return NewEvmLiquidityManager(address), nil
	default:
		return nil, fmt.Errorf("liquidity manager of type %v is not supported", typ)
	}
}
