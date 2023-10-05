package ccipdata

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	COMMIT_PRICE_UPDATES = "Commit price updates"
	FEE_TOKEN_ADDED      = "Fee token added"
	FEE_TOKEN_REMOVED    = "Fee token removed"
)

var (
	UsdPerUnitGasUpdatedV1_0_0 = abihelpers.MustGetEventID("UsdPerUnitGasUpdated", abihelpers.MustParseABI(price_registry.PriceRegistryABI))
)

type TokenPrice struct {
	Token common.Address
	Value *big.Int
}

type TokenPriceUpdate struct {
	TokenPrice
	Timestamp *big.Int
}

type GasPrice struct {
	DestChainSelector uint64
	Value             *big.Int
}

type GasPriceUpdate struct {
	GasPrice
	Timestamp *big.Int
}

//go:generate mockery --quiet --name PriceRegistryReader --output . --filename price_registry_reader_mock.go --inpackage --case=underscore
type PriceRegistryReader interface {
	Close(qopts ...pg.QOpt) error
	// GetTokenPriceUpdatesCreatedAfter returns all the token price updates that happened after the provided timestamp.
	GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]Event[TokenPriceUpdate], error)
	// GetGasPriceUpdatesCreatedAfter returns all the gas price updates that happened after the provided timestamp.
	GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]Event[GasPriceUpdate], error)
	Address() common.Address
	FeeTokenEvents() []common.Hash
	GetFeeTokens(ctx context.Context) ([]common.Address, error)
	GetTokenPrices(ctx context.Context, wantedTokens []common.Address) ([]TokenPriceUpdate, error)
}

// NewPriceRegistryReader determines the appropriate version of the price registry and returns a reader for it.
func NewPriceRegistryReader(lggr logger.Logger, priceRegistryAddress common.Address, lp logpoller.LogPoller, cl client.Client) (PriceRegistryReader, error) {
	_, version, err := ccipconfig.TypeAndVersion(priceRegistryAddress, cl)
	if err != nil {
		// TODO: would this always through a method not found?
		// Unfortunately the v1 price registry doesn't have a method to get the version so assume if it errors
		// its v1.
		return NewPriceRegistryV1_0_0(lggr, priceRegistryAddress, lp, cl)
	}
	switch version.String() {
	case v1_2_0:
		// TODO: ABI is same now but will break shortly with multigas price updates
		return NewPriceRegistryV1_0_0(lggr, priceRegistryAddress, lp, cl)
	default:
		return nil, errors.Errorf("got unexpected version %v", version.String())
	}
}
