package ccipdata

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	_ PriceRegistryReader = &PriceRegistryV1_0_0{}
)

const ExecPluginLabel = "exec"

type PriceRegistryV1_0_0 struct {
	priceRegistry price_registry.PriceRegistryInterface
	address       common.Address
	lp            logpoller.LogPoller
	lggr          logger.Logger
	filters       []logpoller.Filter
	tokenUpdated  common.Hash
	gasUpdated    common.Hash
}

func (p *PriceRegistryV1_0_0) FeeTokenEvents() []common.Hash {
	priceRegABI := abihelpers.MustParseABI(price_registry.PriceRegistryABI)
	return []common.Hash{abihelpers.MustGetEventID("FeeTokenRemoved", priceRegABI), abihelpers.MustGetEventID("FeeTokenAdded", priceRegABI)}
}

func (p *PriceRegistryV1_0_0) GetTokenPrices(ctx context.Context, wantedTokens []common.Address) ([]TokenPriceUpdate, error) {
	tps, err := p.priceRegistry.GetTokenPrices(&bind.CallOpts{Context: ctx}, wantedTokens)
	if err != nil {
		return nil, err
	}
	var tpu []TokenPriceUpdate
	for i, tp := range tps {
		tpu = append(tpu, TokenPriceUpdate{
			TokenPrice: TokenPrice{
				Token: wantedTokens[i],
				Value: tp.Value,
			},
			Timestamp: big.NewInt(int64(tp.Timestamp)), // TODO: valid conversion
		})
	}
	return tpu, nil
}

func (p *PriceRegistryV1_0_0) Address() common.Address {
	return p.address
}

func (p *PriceRegistryV1_0_0) GetFeeTokens(ctx context.Context) ([]common.Address, error) {
	return p.priceRegistry.GetFeeTokens(&bind.CallOpts{Context: ctx})
}

func (p *PriceRegistryV1_0_0) Close(opts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(p.lp, p.filters, opts...)
}

func (p *PriceRegistryV1_0_0) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]Event[TokenPriceUpdate], error) {
	logs, err := p.lp.LogsCreatedAfter(
		p.tokenUpdated,
		p.address,
		ts,
		confs,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return parseLogs[TokenPriceUpdate](
		logs,
		p.lggr,
		func(log types.Log) (*TokenPriceUpdate, error) {
			tp, err := p.priceRegistry.ParseUsdPerTokenUpdated(log)
			if err != nil {
				return nil, err
			}
			return &TokenPriceUpdate{
				TokenPrice: TokenPrice{
					Token: tp.Token,
					Value: tp.Value,
				},
				Timestamp: tp.Timestamp,
			}, nil
		},
	)
}

func (p *PriceRegistryV1_0_0) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]Event[GasPriceUpdate], error) {
	logs, err := p.lp.IndexedLogsCreatedAfter(
		p.gasUpdated,
		p.address,
		1,
		[]common.Hash{abihelpers.EvmWord(chainSelector)},
		ts,
		confs,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return parseLogs[GasPriceUpdate](
		logs,
		p.lggr,
		func(log types.Log) (*GasPriceUpdate, error) {
			p, err := p.priceRegistry.ParseUsdPerUnitGasUpdated(log)
			if err != nil {
				return nil, err
			}
			return &GasPriceUpdate{
				GasPrice: GasPrice{
					DestChainSelector: p.DestChain,
					Value:             p.Value,
				},
				Timestamp: p.Timestamp,
			}, nil
		},
	)
}

func NewPriceRegistryV1_0_0(lggr logger.Logger, priceRegistryAddr common.Address, lp logpoller.LogPoller, ec client.Client) (*PriceRegistryV1_0_0, error) {
	// TODO pass label
	priceRegistry, err := price_registry.NewPriceRegistry(priceRegistryAddr, ec)
	if err != nil {
		return nil, err
	}
	priceRegistryABI := abihelpers.MustParseABI(price_registry.PriceRegistryABI)
	// TODO: clean up strings
	tokenUpdated := abihelpers.MustGetEventID("UsdPerTokenUpdated", priceRegistryABI)
	var filters = []logpoller.Filter{{
		Name:      logpoller.FilterName(COMMIT_PRICE_UPDATES, priceRegistryAddr.String()),
		EventSigs: []common.Hash{UsdPerUnitGasUpdatedV1_0_0, tokenUpdated},
		Addresses: []common.Address{priceRegistryAddr},
	},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_ADDED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{abihelpers.MustGetEventID("FeeTokenAdded", priceRegistryABI)},
			Addresses: []common.Address{priceRegistryAddr},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_REMOVED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{abihelpers.MustGetEventID("FeeTokenAdded", priceRegistryABI)},
			Addresses: []common.Address{priceRegistryAddr},
		}}
	err = logpollerutil.RegisterLpFilters(lp, filters)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryV1_0_0{
		priceRegistry: priceRegistry,
		address:       priceRegistryAddr,
		lp:            lp,
		lggr:          lggr,
		gasUpdated:    UsdPerUnitGasUpdatedV1_0_0,
		tokenUpdated:  tokenUpdated,
		filters:       filters,
	}, nil
}
