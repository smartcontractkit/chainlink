package liquidityrebalancer

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

//go:generate mockery --quiet --name Rebalancer --output ../rebalancermocks --filename rebalancer_mock.go --case=underscore
type Rebalancer interface {
	ComputeTransfersToBalance(
		g liquiditygraph.LiquidityGraph,
		inflightTransfers []models.PendingTransfer,
		medianLiquidityPerChain []models.NetworkLiquidity,
	) ([]models.Transfer, error)
}
