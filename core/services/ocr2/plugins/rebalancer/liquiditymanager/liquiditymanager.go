package liquiditymanager

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// LiquidityManager is an abstraction of the liquidity manager contract.
//
//go:generate mockery --quiet --name LiquidityManager --output ../rebalancermocks --filename lm_mock.go --case=underscore
type LiquidityManager interface {
	// GetLiquidityManagers returns a mapping that contains the liquidity managers for each destination chain.
	GetLiquidityManagers(ctx context.Context) (map[models.NetworkSelector]models.Address, error)

	// GetBalance returns the current token/liquidity balance.
	GetBalance(ctx context.Context) (*big.Int, error)

	// GetPendingTransfers returns the pending liquidity transfers.
	GetPendingTransfers(ctx context.Context, since time.Time) ([]models.PendingTransfer, error)

	// Discover discovers other liquidity managers
	Discover(ctx context.Context, lmFactory Factory) (*Registry, liquiditygraph.LiquidityGraph, error)

	// Close releases any resources.
	Close(ctx context.Context) error
}
