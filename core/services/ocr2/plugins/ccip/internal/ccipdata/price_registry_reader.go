package ccipdata

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	COMMIT_PRICE_UPDATES = "Commit price updates"
	FEE_TOKEN_ADDED      = "Fee token added"
	FEE_TOKEN_REMOVED    = "Fee token removed"
	ExecPluginLabel      = "exec"
)

type TokenPrice struct {
	Token common.Address
	Value *big.Int
}

type TokenPriceUpdate struct {
	TokenPrice
	TimestampUnixSec *big.Int
}

type GasPrice struct {
	DestChainSelector uint64
	Value             *big.Int
}

type GasPriceUpdate struct {
	GasPrice
	TimestampUnixSec *big.Int
}

//go:generate mockery --quiet --name PriceRegistryReader --filename price_registry_reader_mock.go --case=underscore
type PriceRegistryReader interface {
	// GetTokenPriceUpdatesCreatedAfter returns all the token price updates that happened after the provided timestamp.
	// The returned updates are sorted by timestamp in ascending order.
	GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]Event[TokenPriceUpdate], error)
	// GetGasPriceUpdatesCreatedAfter returns all the gas price updates that happened after the provided timestamp.
	// The returned updates are sorted by timestamp in ascending order.
	GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]Event[GasPriceUpdate], error)
	Address() common.Address
	GetFeeTokens(ctx context.Context) ([]common.Address, error)
	GetTokenPrices(ctx context.Context, wantedTokens []common.Address) ([]TokenPriceUpdate, error)
	// TODO: consider moving this method to a different interface since it's not related to the price registry
	GetTokensDecimals(ctx context.Context, tokenAddresses []common.Address) ([]uint8, error)
	Close() error
}
