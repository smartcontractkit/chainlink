package v1_0_0

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	abiERC20                              = abihelpers.MustParseABI(erc20.ERC20ABI)
	_        ccipdata.PriceRegistryReader = &PriceRegistry{}
	// Exposed only for backwards compatibility with tests.
	UsdPerUnitGasUpdated = abihelpers.MustGetEventID("UsdPerUnitGasUpdated", abihelpers.MustParseABI(price_registry_1_0_0.PriceRegistryABI))
)

type PriceRegistry struct {
	priceRegistry   price_registry_1_0_0.PriceRegistryInterface
	address         common.Address
	lp              logpoller.LogPoller
	evmBatchCaller  rpclib.EvmBatchCaller
	lggr            logger.Logger
	filters         []logpoller.Filter
	tokenUpdated    common.Hash
	gasUpdated      common.Hash
	feeTokenAdded   common.Hash
	feeTokenRemoved common.Hash

	feeTokensCache     cache.AutoSync[[]common.Address]
	tokenDecimalsCache sync.Map
}

func NewPriceRegistry(lggr logger.Logger, priceRegistryAddr common.Address, lp logpoller.LogPoller, ec client.Client, registerFilters bool) (*PriceRegistry, error) {
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
			Name:      logpoller.FilterName(ccipdata.COMMIT_PRICE_UPDATES, priceRegistryAddr.String()),
			EventSigs: []common.Hash{UsdPerUnitGasUpdated, usdPerTokenUpdated},
			Addresses: []common.Address{priceRegistryAddr},
		},
		{
			Name:      logpoller.FilterName(ccipdata.FEE_TOKEN_ADDED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{feeTokenAdded},
			Addresses: []common.Address{priceRegistryAddr},
		},
		{
			Name:      logpoller.FilterName(ccipdata.FEE_TOKEN_REMOVED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{feeTokenRemoved},
			Addresses: []common.Address{priceRegistryAddr},
		}}
	if registerFilters {
		err = logpollerutil.RegisterLpFilters(lp, filters)
		if err != nil {
			return nil, err
		}
	}
	return &PriceRegistry{
		priceRegistry: priceRegistry,
		address:       priceRegistryAddr,
		lp:            lp,
		evmBatchCaller: rpclib.NewDynamicLimitedBatchCaller(
			lggr,
			ec,
			rpclib.DefaultRpcBatchSizeLimit,
			rpclib.DefaultRpcBatchBackOffMultiplier,
		),
		lggr:            lggr,
		gasUpdated:      UsdPerUnitGasUpdated,
		tokenUpdated:    usdPerTokenUpdated,
		feeTokenRemoved: feeTokenRemoved,
		feeTokenAdded:   feeTokenAdded,
		filters:         filters,
		feeTokensCache: cache.NewLogpollerEventsBased[[]common.Address](
			lp,
			[]common.Hash{feeTokenAdded, feeTokenRemoved},
			priceRegistryAddr,
		),
	}, nil
}

func (p *PriceRegistry) GetTokenPrices(ctx context.Context, wantedTokens []common.Address) ([]ccipdata.TokenPriceUpdate, error) {
	tps, err := p.priceRegistry.GetTokenPrices(&bind.CallOpts{Context: ctx}, wantedTokens)
	if err != nil {
		return nil, err
	}
	var tpu []ccipdata.TokenPriceUpdate
	for i, tp := range tps {
		tpu = append(tpu, ccipdata.TokenPriceUpdate{
			TokenPrice: ccipdata.TokenPrice{
				Token: wantedTokens[i],
				Value: tp.Value,
			},
			TimestampUnixSec: big.NewInt(int64(tp.Timestamp)),
		})
	}
	return tpu, nil
}

func (p *PriceRegistry) Address() common.Address {
	return p.address
}

func (p *PriceRegistry) GetFeeTokens(ctx context.Context) ([]common.Address, error) {
	feeTokens, err := p.feeTokensCache.Get(ctx, func(ctx context.Context) ([]common.Address, error) {
		return p.priceRegistry.GetFeeTokens(&bind.CallOpts{Context: ctx})
	})

	if err != nil {
		return nil, fmt.Errorf("get fee tokens: %w", err)
	}
	return feeTokens, nil
}

func (p *PriceRegistry) Close() error {
	return logpollerutil.UnregisterLpFilters(p.lp, p.filters)
}

func (p *PriceRegistry) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]ccipdata.Event[ccipdata.TokenPriceUpdate], error) {
	logs, err := p.lp.LogsCreatedAfter(
		p.tokenUpdated,
		p.address,
		ts,
		logpoller.Confirmations(confs),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return ccipdata.ParseLogs[ccipdata.TokenPriceUpdate](
		logs,
		p.lggr,
		func(log types.Log) (*ccipdata.TokenPriceUpdate, error) {
			tp, err := p.priceRegistry.ParseUsdPerTokenUpdated(log)
			if err != nil {
				return nil, err
			}
			return &ccipdata.TokenPriceUpdate{
				TokenPrice: ccipdata.TokenPrice{
					Token: tp.Token,
					Value: tp.Value,
				},
				TimestampUnixSec: tp.Timestamp,
			}, nil
		},
	)
}

func (p *PriceRegistry) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]ccipdata.Event[ccipdata.GasPriceUpdate], error) {
	logs, err := p.lp.IndexedLogsCreatedAfter(
		p.gasUpdated,
		p.address,
		1,
		[]common.Hash{abihelpers.EvmWord(chainSelector)},
		ts,
		logpoller.Confirmations(confs),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return ccipdata.ParseLogs[ccipdata.GasPriceUpdate](
		logs,
		p.lggr,
		func(log types.Log) (*ccipdata.GasPriceUpdate, error) {
			p, err := p.priceRegistry.ParseUsdPerUnitGasUpdated(log)
			if err != nil {
				return nil, err
			}
			return &ccipdata.GasPriceUpdate{
				GasPrice: ccipdata.GasPrice{
					DestChainSelector: p.DestChain,
					Value:             p.Value,
				},
				TimestampUnixSec: p.Timestamp,
			}, nil
		},
	)
}

func (p *PriceRegistry) GetTokensDecimals(ctx context.Context, tokenAddresses []common.Address) ([]uint8, error) {
	found := make(map[common.Address]bool)
	tokenDecimals := make([]uint8, len(tokenAddresses))
	for i, tokenAddress := range tokenAddresses {
		if v, ok := p.tokenDecimalsCache.Load(tokenAddress); ok {
			if decimals, isUint8 := v.(uint8); isUint8 {
				tokenDecimals[i] = decimals
				found[tokenAddress] = true
			} else {
				p.lggr.Errorf("token decimals cache contains invalid type %T", v)
			}
		}
	}
	if len(found) == len(tokenAddresses) {
		return tokenDecimals, nil
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(tokenAddresses))
	for _, tokenAddress := range tokenAddresses {
		if !found[tokenAddress] {
			evmCalls = append(evmCalls, rpclib.NewEvmCall(abiERC20, "decimals", tokenAddress))
		}
	}

	latestBlock, err := p.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	results, err := p.evmBatchCaller.BatchCall(ctx, uint64(latestBlock.BlockNumber), evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	decimals, err := rpclib.ParseOutputs[uint8](results, func(d rpclib.DataAndErr) (uint8, error) {
		return rpclib.ParseOutput[uint8](d, 0)
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}

	j := 0
	for i, tokenAddress := range tokenAddresses {
		if !found[tokenAddress] {
			tokenDecimals[i] = decimals[j]
			p.tokenDecimalsCache.Store(tokenAddress, tokenDecimals[i])
			j++
		}
	}
	return tokenDecimals, nil
}
