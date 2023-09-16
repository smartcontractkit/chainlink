package pricegetter

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockery --quiet --name PriceGetter --output . --filename mock.go --inpackage --case=underscore
type PriceGetter interface {
	// TokenPricesUSD returns token prices in USD.
	// Note: The result might contain tokens that are not passed with the 'tokens' param.
	//       The opposite cannot happen, an error will be returned if a token price was not found.
	TokenPricesUSD(ctx context.Context, tokens []common.Address) (map[common.Address]*big.Int, error)
}
