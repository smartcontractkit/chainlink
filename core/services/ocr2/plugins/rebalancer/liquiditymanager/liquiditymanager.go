package liquiditymanager

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// LiquidityManager is an abstraction of the liquidity manager contract.
//
//go:generate mockery --quiet --name LiquidityManager --output ../rebalancermocks --filename lm_mock.go --case=underscore
type LiquidityManager interface {
	// MoveLiquidity moves the target amount to the liquidity manager of the provided chain.
	// todo: consider passing some meta
	MoveLiquidity(ctx context.Context, chainID models.NetworkID, amount *big.Int) error

	// GetLiquidityManagers returns a mapping that contains the liquidity managers for each destination chain.
	GetLiquidityManagers(ctx context.Context) (map[models.NetworkID]models.Address, error)

	// GetBalance returns the current token/liquidity balance.
	GetBalance(ctx context.Context) (*big.Int, error)

	// GetPendingTransfers returns the pending liquidity transfers.
	GetPendingTransfers(ctx context.Context) ([]models.PendingTransfer, error)

	// Discover discovers other liquidity managers
	Discover(ctx context.Context, lmFactory Factory) (*Registry, liquiditygraph.LiquidityGraph, error)

	// Close releases any resources.
	Close(ctx context.Context) error
}
