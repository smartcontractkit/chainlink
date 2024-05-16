package reader

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TokenPrices interface {
	// GetTokenPricesUSD returns the prices of the provided tokens in USD.
	// The order of the returned prices corresponds to the order of the provided tokens.
	GetTokenPricesUSD(ctx context.Context, tokens []types.Account) ([]*big.Int, error)
}
