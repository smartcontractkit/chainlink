package evmliquiditymanager

import (
	"context"
	"math/big"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// LiquidityManager is an abstraction of the rebalancer contract.
// TODO: extract the common interface which is chain dependent.
type LiquidityManager interface {
	// GetRebalancers returns a mapping that contains the rebalancers for each destination chain.
	GetRebalancers(ctx context.Context) (map[models.NetworkSelector]models.Address, error)

	// GetBalance returns the current token/liquidity balance.
	GetBalance(ctx context.Context) (*big.Int, error)

	// Close releases any resources.
	Close(ctx context.Context) error

	// ConfigDigest returns the OCR config digest for the rebalancer.
	ConfigDigest(ctx context.Context) (ocrtypes.ConfigDigest, error)

	// GetTokenAddress returns the token address of the rebalancer.
	GetTokenAddress(ctx context.Context) (models.Address, error)

	// GetLatestSequenceNumber returns the latest sequence number of the rebalancer.
	GetLatestSequenceNumber(ctx context.Context) (uint64, error)
}
