package pricegetter

import (
	"context"
	"math/big"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type PriceGetter interface {
	cciptypes.PriceGetter
}

type AllTokensPriceGetter interface {
	PriceGetter
	// GetJobSpecTokenPricesUSD returns all token prices defined in the jobspec.
	GetJobSpecTokenPricesUSD(ctx context.Context) (map[cciptypes.Address]*big.Int, error)
}
