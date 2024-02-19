package cciptypes

import (
	"context"
	"math/big"
)

type PriceGetter interface {
	// TokenPricesUSD returns token prices in USD.
	// Note: The result might contain tokens that are not passed with the 'tokens' param.
	//       The opposite cannot happen, an error will be returned if a token price was not found.
	TokenPricesUSD(ctx context.Context, tokens []Address) (map[Address]*big.Int, error)
}
