package ccip

import (
	"context"
	"io"
	"math/big"
)

type PriceGetter interface {
	// FilterConfiguredTokens filters a list of token addresses
	// for only those that are configured to be able to get a price and those that aren't
	FilterConfiguredTokens(ctx context.Context, tokens []Address) (configured []Address, unconfigured []Address, err error)
	// TokenPricesUSD returns token prices in USD.
	// Note: The result might contain tokens that are not passed with the 'tokens' param.
	//       The opposite cannot happen, an error will be returned if a token price was not found.
	TokenPricesUSD(ctx context.Context, tokens []Address) (map[Address]*big.Int, error)
	io.Closer
}
