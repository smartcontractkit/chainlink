package ccipdata

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	Close(qopts ...pg.QOpt) error
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
}

// NewPriceRegistryReader determines the appropriate version of the price registry and returns a reader for it.
func NewPriceRegistryReader(lggr logger.Logger, priceRegistryAddress common.Address, lp logpoller.LogPoller, cl client.Client) (PriceRegistryReader, error) {
	_, version, err := ccipconfig.TypeAndVersion(priceRegistryAddress, cl)
	if err != nil {
		if strings.Contains(err.Error(), "execution reverted") {
			lggr.Infof("Assuming %v is 1.0.0 price registry, got %v", priceRegistryAddress.String(), err)
			// Unfortunately the v1 price registry doesn't have a method to get the version so assume if it reverts
			// its v1.
			return NewPriceRegistryV1_0_0(lggr, priceRegistryAddress, lp, cl)
		}
		return nil, err
	}
	switch version.String() {
	case V1_2_0:
		return NewPriceRegistryV1_2_0(lggr, priceRegistryAddress, lp, cl)
	default:
		return nil, errors.Errorf("got unexpected version %v", version.String())
	}
}
