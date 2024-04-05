package ccip

import (
	"context"
	"io"
	"math/big"
)

type PriceGetter interface {
	// IsTokenConfigured returns if a token address is configured to be able to get a price
	IsTokenConfigured(ctx context.Context, token Address) (bool, error)
	// TokenPricesUSD returns token prices in USD.
	// Note: The result might contain tokens that are not passed with the 'tokens' param.
	//       The opposite cannot happen, an error will be returned if a token price was not found.
	TokenPricesUSD(ctx context.Context, tokens []Address) (map[Address]*big.Int, error)
	io.Closer
}
