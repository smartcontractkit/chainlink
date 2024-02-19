package cciptypes

import (
	"context"
	"math/big"
	"time"
)

type PriceRegistryReader interface {
	// GetTokenPriceUpdatesCreatedAfter returns all the token price updates that happened after the provided timestamp.
	// The returned updates are sorted by timestamp in ascending order.
	GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confirmations int) ([]TokenPriceUpdateWithTxMeta, error)

	// GetGasPriceUpdatesCreatedAfter returns all the gas price updates that happened after the provided timestamp.
	// The returned updates are sorted by timestamp in ascending order.
	GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confirmations int) ([]GasPriceUpdateWithTxMeta, error)

	Address() Address

	GetFeeTokens(ctx context.Context) ([]Address, error)

	GetTokenPrices(ctx context.Context, wantedTokens []Address) ([]TokenPriceUpdate, error)

	GetTokensDecimals(ctx context.Context, tokenAddresses []Address) ([]uint8, error)

	Close() error
}

type TokenPriceUpdateWithTxMeta struct {
	TxMeta
	TokenPriceUpdate
}

type TokenPriceUpdate struct {
	TokenPrice
	TimestampUnixSec *big.Int
}

type GasPriceUpdateWithTxMeta struct {
	TxMeta
	GasPriceUpdate
}

type GasPriceUpdate struct {
	GasPrice
	TimestampUnixSec *big.Int
}
