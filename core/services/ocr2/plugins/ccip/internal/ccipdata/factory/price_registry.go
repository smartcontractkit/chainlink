package factory

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

// NewPriceRegistryReader determines the appropriate version of the price registry and returns a reader for it.
func NewPriceRegistryReader(lggr logger.Logger, priceRegistryAddress common.Address, lp logpoller.LogPoller, cl client.Client) (ccipdata.PriceRegistryReader, error) {
	contractType, version, err := ccipconfig.TypeAndVersion(priceRegistryAddress, cl)
	if err != nil {
		if strings.Contains(err.Error(), "execution reverted") {
			lggr.Infof("Assuming %v is 1.0.0 price registry, got %v", priceRegistryAddress.String(), err)
			// Unfortunately the v1 price registry doesn't have a method to get the version so assume if it reverts
			// its v1.
			return v1_0_0.NewPriceRegistry(lggr, priceRegistryAddress, lp, cl)
		}
		return nil, errors.Wrapf(err, "unable to read type and version")
	}
	if contractType != ccipconfig.PriceRegistry {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.PriceRegistry, contractType)
	}
	switch version.String() {
	case ccipdata.V1_2_0:
		return v1_2_0.NewPriceRegistry(lggr, priceRegistryAddress, lp, cl)
	default:
		return nil, errors.Errorf("unsupported price registry version %v", version.String())
	}
}
