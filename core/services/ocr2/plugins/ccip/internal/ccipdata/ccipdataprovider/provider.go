package ccipdataprovider

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
)

type PriceRegistry interface {
	NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error)
}

type EvmPriceRegistry struct {
	lp          logpoller.LogPoller
	ec          client.Client
	lggr        logger.Logger
	pluginLabel string
}

func NewEvmPriceRegistry(lp logpoller.LogPoller, ec client.Client, lggr logger.Logger, pluginLabel string) *EvmPriceRegistry {
	return &EvmPriceRegistry{
		lp:          lp,
		ec:          ec,
		lggr:        lggr,
		pluginLabel: pluginLabel,
	}
}

func (p *EvmPriceRegistry) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (cciptypes.PriceRegistryReader, error) {
	destPriceRegistryReader, err := factory.NewPriceRegistryReader(ctx, p.lggr, factory.NewEvmVersionFinder(), addr, p.lp, p.ec)
	if err != nil {
		return nil, err
	}
	return observability.NewPriceRegistryReader(destPriceRegistryReader, p.ec.ConfiguredChainID().Int64(), p.pluginLabel), nil
}
