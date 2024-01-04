package liquiditymanager

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type EvmLiquidityManager struct{}

func NewEvmLiquidityManager(address models.Address) *EvmLiquidityManager {
	return &EvmLiquidityManager{}
}

func (e EvmLiquidityManager) MoveLiquidity(ctx context.Context, chainID models.NetworkID, amount *big.Int) error {
	return nil
}

func (e EvmLiquidityManager) GetLiquidityManagers(ctx context.Context) (map[models.NetworkID]models.Address, error) {
	return nil, nil
}

func (e EvmLiquidityManager) GetBalance(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (e EvmLiquidityManager) GetPendingTransfers(ctx context.Context) ([]models.PendingTransfer, error) {
	return nil, nil
}

func (e EvmLiquidityManager) Close(ctx context.Context) error {
	return nil
}
