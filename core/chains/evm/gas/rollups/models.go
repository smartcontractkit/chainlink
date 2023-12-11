package rollups

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// L1Oracle provides interface for fetching L1-specific fee components if the chain is an L2.
// For example, on Optimistic Rollups, this oracle can return rollup-specific l1BaseFee
//
//go:generate mockery --quiet --name L1Oracle --output ./mocks/ --case=underscore
type L1Oracle interface {
	services.ServiceCtx

	GasPrice(ctx context.Context) (*assets.Wei, error)
}
