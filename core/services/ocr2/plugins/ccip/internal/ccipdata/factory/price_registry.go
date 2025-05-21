package factory

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

// NewPriceRegistryReader determines the appropriate version of the price registry and returns a reader for it.
func NewPriceRegistryReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, priceRegistryAddress cciptypes.Address, lp logpoller.LogPoller, cl client.Client) (ccipdata.PriceRegistryReader, error) {
	return initOrClosePriceRegistryReader(ctx, lggr, versionFinder, priceRegistryAddress, lp, cl, false)
}

func ClosePriceRegistryReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, priceRegistryAddress cciptypes.Address, lp logpoller.LogPoller, cl client.Client) error {
	_, err := initOrClosePriceRegistryReader(ctx, lggr, versionFinder, priceRegistryAddress, lp, cl, true)
	return err
}

func initOrClosePriceRegistryReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, priceRegistryAddress cciptypes.Address, lp logpoller.LogPoller, cl client.Client, closeReader bool) (ccipdata.PriceRegistryReader, error) {
	registerFilters := !closeReader

	priceRegistryEvmAddr, err := ccipcalc.GenericAddrToEvm(priceRegistryAddress)
	if err != nil {
		return nil, err
	}

	contractType, version, err := versionFinder.TypeAndVersion(priceRegistryAddress, cl)
	if err != nil {
		return nil, err
	}
	if contractType != ccipconfig.PriceRegistry {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.PriceRegistry, contractType)
	}
	switch version.String() {
	case ccipdata.V1_2_0:
		pr, err := v1_2_0.NewPriceRegistry(ctx, lggr, priceRegistryEvmAddr, lp, cl, registerFilters)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, pr.Close()
		}
		return pr, nil
	case ccipdata.V1_6_0:
		pr, err := v1_2_0.NewPriceRegistry(ctx, lggr, priceRegistryEvmAddr, lp, cl, registerFilters)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, pr.Close()
		}
		return pr, nil
	default:
		return nil, errors.Errorf("unsupported price registry version %v", version.String())
	}
}
