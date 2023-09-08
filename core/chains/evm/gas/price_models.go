package gas

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

type PriceType string

const (
	GAS_PRICE   PriceType = "GAS_PRICE"
	L1_BASE_FEE PriceType = "L1_BASE_FEE"
)

type PriceComponent struct {
	Price     *assets.Wei
	PriceType PriceType
}

// PriceComponentGetter provides interface for implementing chain-specific price components
// On L1 chains like Ethereum or Avax, the only component is the gas price.
// On Optimistic Rollups, there are two components - the L2 gas price, and L1 base fee for data availability.
// On future chains, there could be more or differing price components.
type PriceComponentGetter interface {
	RefreshComponents(ctx context.Context) error
	// GetPriceComponents The first component in prices should always be the passed-in gasPrice
	GetPriceComponents(ctx context.Context, gasPrice *assets.Wei) (prices []PriceComponent, err error)
}
