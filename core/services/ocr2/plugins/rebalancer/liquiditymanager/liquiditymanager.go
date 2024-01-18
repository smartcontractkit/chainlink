package liquiditymanager

import (
	"context"
	"math/big"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// Rebalancer is an abstraction of the rebalancer contract.
//
//go:generate mockery --quiet --name Rebalancer --output ./mocks --filename rebalancer_mock.go --case=underscore
type Rebalancer interface {
	// GetRebalancers returns a mapping that contains the rebalancers for each destination chain.
	GetRebalancers(ctx context.Context) (map[models.NetworkSelector]models.Address, error)

	// GetBalance returns the current token/liquidity balance.
	GetBalance(ctx context.Context) (*big.Int, error)

	// GetPendingTransfers returns the pending liquidity transfers.
	GetPendingTransfers(ctx context.Context, since time.Time) ([]models.PendingTransfer, error)

	// Discover discovers other rebalancers
	Discover(ctx context.Context, lmFactory Factory) (*Registry, liquiditygraph.LiquidityGraph, error)

	// Close releases any resources.
	Close(ctx context.Context) error

	// ConfigDigest returns the OCR config digest for the rebalancer.
	ConfigDigest(ctx context.Context) (ocrtypes.ConfigDigest, error)
}
