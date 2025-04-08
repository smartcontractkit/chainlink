package pricegetter

import (
	"context"
	"io"
	"math/big"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
)

type PriceGetter interface {
	cciptypes.PriceGetter
}

type AllTokensPriceGetter interface {
	io.Closer

	// GetJobSpecTokenPricesUSD returns all token prices defined in the jobspec.
	GetJobSpecTokenPricesUSD(ctx context.Context) (map[ccipcommon.TokenID]*big.Int, error)
	// GetTokenPricesUSD returns the prices of the provided tokens in USD.
	GetTokenPricesUSD(ctx context.Context, tokens []ccipcommon.TokenID) (map[ccipcommon.TokenID]*big.Int, error)
}
