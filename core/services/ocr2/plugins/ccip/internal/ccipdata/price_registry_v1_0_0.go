package ccipdata

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	_ PriceRegistryReader = &PriceRegistryV1_0_0{}
	// Exposed only for backwards compatibility with tests.
	UsdPerUnitGasUpdatedV1_0_0 = abihelpers.MustGetEventID("UsdPerUnitGasUpdated", abihelpers.MustParseABI(price_registry_1_0_0.PriceRegistryABI))
)

type PriceRegistryV1_0_0 struct {
	priceRegistry   price_registry_1_0_0.PriceRegistryInterface
	address         common.Address
	lp              logpoller.LogPoller
	lggr            logger.Logger
	filters         []logpoller.Filter
	tokenUpdated    common.Hash
	gasUpdated      common.Hash
	feeTokenAdded   common.Hash
	feeTokenRemoved common.Hash
}

func (p *PriceRegistryV1_0_0) FeeTokenEvents() []common.Hash {
	return []common.Hash{p.feeTokenAdded, p.feeTokenRemoved}
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
			TimestampUnixSec: big.NewInt(int64(tp.Timestamp)),
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
				TimestampUnixSec: tp.Timestamp,
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
				TimestampUnixSec: p.Timestamp,
			}, nil
		},
	)
}

func NewPriceRegistryV1_0_0(lggr logger.Logger, priceRegistryAddr common.Address, lp logpoller.LogPoller, ec client.Client) (*PriceRegistryV1_0_0, error) {
	priceRegistry, err := price_registry_1_0_0.NewPriceRegistry(priceRegistryAddr, ec)
	if err != nil {
		return nil, err
	}
	priceRegABI := abihelpers.MustParseABI(price_registry_1_0_0.PriceRegistryABI)
	usdPerTokenUpdated := abihelpers.MustGetEventID("UsdPerTokenUpdated", priceRegABI)
	feeTokenRemoved := abihelpers.MustGetEventID("FeeTokenRemoved", priceRegABI)
	feeTokenAdded := abihelpers.MustGetEventID("FeeTokenAdded", priceRegABI)
	var filters = []logpoller.Filter{
		{
			Name:      logpoller.FilterName(COMMIT_PRICE_UPDATES, priceRegistryAddr.String()),
			EventSigs: []common.Hash{UsdPerUnitGasUpdatedV1_0_0, usdPerTokenUpdated},
			Addresses: []common.Address{priceRegistryAddr},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_ADDED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{feeTokenAdded},
			Addresses: []common.Address{priceRegistryAddr},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_REMOVED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{feeTokenRemoved},
			Addresses: []common.Address{priceRegistryAddr},
		}}
	err = logpollerutil.RegisterLpFilters(lp, filters)
	if err != nil {
		return nil, err
	}
	return &PriceRegistryV1_0_0{
		priceRegistry:   priceRegistry,
		address:         priceRegistryAddr,
		lp:              lp,
		lggr:            lggr,
		gasUpdated:      UsdPerUnitGasUpdatedV1_0_0,
		tokenUpdated:    usdPerTokenUpdated,
		feeTokenRemoved: feeTokenRemoved,
		feeTokenAdded:   feeTokenAdded,
		filters:         filters,
	}, nil
}

// ApplyPriceRegistryUpdateV1_0_0 is a helper function used in tests only.
func ApplyPriceRegistryUpdateV1_0_0(t *testing.T, user *bind.TransactOpts, addr common.Address, ec client.Client, gasPrice []GasPrice, tokenPrices []TokenPrice) {
	require.True(t, len(gasPrice) <= 1)
	pr, err := price_registry_1_0_0.NewPriceRegistry(addr, ec)
	require.NoError(t, err)
	var tps []price_registry_1_0_0.InternalTokenPriceUpdate
	for _, tp := range tokenPrices {
		tps = append(tps, price_registry_1_0_0.InternalTokenPriceUpdate{
			SourceToken: tp.Token,
			UsdPerToken: tp.Value,
		})
	}
	dest := uint64(0)
	gas := big.NewInt(0)
	if len(gasPrice) == 1 {
		dest = gasPrice[0].DestChainSelector
		gas = gasPrice[0].Value
	}
	_, err = pr.UpdatePrices(user, price_registry_1_0_0.InternalPriceUpdates{
		TokenPriceUpdates: tps,
		DestChainSelector: dest,
		UsdPerUnitGas:     gas,
	})
	require.NoError(t, err)
}
