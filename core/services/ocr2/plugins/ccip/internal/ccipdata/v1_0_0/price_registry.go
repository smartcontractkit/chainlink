package v1_0_0

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
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
			Retention: ccipdata.PriceUpdatesLogsRetention,
		},
		{
			Name:      logpoller.FilterName(ccipdata.FEE_TOKEN_ADDED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{feeTokenAdded},
			Addresses: []common.Address{priceRegistryAddr},
			Retention: ccipdata.CacheEvictionLogsRetention,
		},
		{
			Name:      logpoller.FilterName(ccipdata.FEE_TOKEN_REMOVED, priceRegistryAddr.String()),
			EventSigs: []common.Hash{feeTokenRemoved},
			Addresses: []common.Address{priceRegistryAddr},
			Retention: ccipdata.CacheEvictionLogsRetention,
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
			rpclib.DefaultMaxParallelRpcCalls,
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

func (p *PriceRegistry) GetTokenPrices(ctx context.Context, wantedTokens []cciptypes.Address) ([]cciptypes.TokenPriceUpdate, error) {
	evmAddrs, err := ccipcalc.GenericAddrsToEvm(wantedTokens...)
	if err != nil {
		return nil, err
	}

	tps, err := p.priceRegistry.GetTokenPrices(&bind.CallOpts{Context: ctx}, evmAddrs)
	if err != nil {
		return nil, err
	}
	var tpu []cciptypes.TokenPriceUpdate
	for i, tp := range tps {
		tpu = append(tpu, cciptypes.TokenPriceUpdate{
			TokenPrice: cciptypes.TokenPrice{
				Token: cciptypes.Address(evmAddrs[i].String()),
				Value: tp.Value,
			},
			TimestampUnixSec: big.NewInt(int64(tp.Timestamp)),
		})
	}
	return tpu, nil
}

func (p *PriceRegistry) Address(ctx context.Context) (cciptypes.Address, error) {
	return cciptypes.Address(p.address.String()), nil
}

func (p *PriceRegistry) GetFeeTokens(ctx context.Context) ([]cciptypes.Address, error) {
	feeTokens, err := p.feeTokensCache.Get(ctx, func(ctx context.Context) ([]common.Address, error) {
		return p.priceRegistry.GetFeeTokens(&bind.CallOpts{Context: ctx})
	})
	if err != nil {
		return nil, fmt.Errorf("get fee tokens: %w", err)
	}

	return ccipcalc.EvmAddrsToGeneric(feeTokens...), nil
}

func (p *PriceRegistry) Close() error {
	return logpollerutil.UnregisterLpFilters(p.lp, p.filters)
}

func (p *PriceRegistry) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]cciptypes.TokenPriceUpdateWithTxMeta, error) {
	logs, err := p.lp.LogsCreatedAfter(
		ctx,
		p.tokenUpdated,
		p.address,
		ts,
		evmtypes.Confirmations(confs),
	)
	if err != nil {
		return nil, err
	}

	parsedLogs, err := ccipdata.ParseLogs[cciptypes.TokenPriceUpdate](
		logs,
		p.lggr,
		func(log types.Log) (*cciptypes.TokenPriceUpdate, error) {
			tp, err1 := p.priceRegistry.ParseUsdPerTokenUpdated(log)
			if err1 != nil {
				return nil, err1
			}
			return &cciptypes.TokenPriceUpdate{
				TokenPrice: cciptypes.TokenPrice{
					Token: cciptypes.Address(tp.Token.String()),
					Value: tp.Value,
				},
				TimestampUnixSec: tp.Timestamp,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	res := make([]cciptypes.TokenPriceUpdateWithTxMeta, 0, len(parsedLogs))
	for _, log := range parsedLogs {
		res = append(res, cciptypes.TokenPriceUpdateWithTxMeta{
			TxMeta:           log.TxMeta,
			TokenPriceUpdate: log.Data,
		})
	}
	return res, nil
}

func (p *PriceRegistry) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	logs, err := p.lp.IndexedLogsCreatedAfter(
		ctx,
		p.gasUpdated,
		p.address,
		1,
		[]common.Hash{abihelpers.EvmWord(chainSelector)},
		ts,
		evmtypes.Confirmations(confs),
	)
	if err != nil {
		return nil, err
	}
	return p.parseGasPriceUpdatesLogs(logs)
}

func (p *PriceRegistry) GetAllGasPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	logs, err := p.lp.LogsCreatedAfter(
		ctx,
		p.gasUpdated,
		p.address,
		ts,
		evmtypes.Confirmations(confs),
	)
	if err != nil {
		return nil, err
	}
	return p.parseGasPriceUpdatesLogs(logs)
}

func (p *PriceRegistry) parseGasPriceUpdatesLogs(logs []logpoller.Log) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	parsedLogs, err := ccipdata.ParseLogs[cciptypes.GasPriceUpdate](
		logs,
		p.lggr,
		func(log types.Log) (*cciptypes.GasPriceUpdate, error) {
			p, err1 := p.priceRegistry.ParseUsdPerUnitGasUpdated(log)
			if err1 != nil {
				return nil, err1
			}
			return &cciptypes.GasPriceUpdate{
				GasPrice: cciptypes.GasPrice{
					DestChainSelector: p.DestChain,
					Value:             p.Value,
				},
				TimestampUnixSec: p.Timestamp,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	res := make([]cciptypes.GasPriceUpdateWithTxMeta, 0, len(parsedLogs))
	for _, log := range parsedLogs {
		res = append(res, cciptypes.GasPriceUpdateWithTxMeta{
			TxMeta:         log.TxMeta,
			GasPriceUpdate: log.Data,
		})
	}
	return res, nil
}

func (p *PriceRegistry) GetTokensDecimals(ctx context.Context, tokenAddresses []cciptypes.Address) ([]uint8, error) {
	evmAddrs, err := ccipcalc.GenericAddrsToEvm(tokenAddresses...)
	if err != nil {
		return nil, err
	}

	found := make(map[common.Address]bool)
	tokenDecimals := make([]uint8, len(evmAddrs))
	for i, tokenAddress := range evmAddrs {
		if v, ok := p.tokenDecimalsCache.Load(tokenAddress); ok {
			if decimals, isUint8 := v.(uint8); isUint8 {
				tokenDecimals[i] = decimals
				found[tokenAddress] = true
			} else {
				p.lggr.Errorf("token decimals cache contains invalid type %T", v)
			}
		}
	}
	if len(found) == len(evmAddrs) {
		return tokenDecimals, nil
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(evmAddrs))
	for _, tokenAddress := range evmAddrs {
		if !found[tokenAddress] {
			evmCalls = append(evmCalls, rpclib.NewEvmCall(abiERC20, "decimals", tokenAddress))
		}
	}

	results, err := p.evmBatchCaller.BatchCall(ctx, 0, evmCalls)
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
	for i, tokenAddress := range evmAddrs {
		if !found[tokenAddress] {
			tokenDecimals[i] = decimals[j]
			p.tokenDecimalsCache.Store(tokenAddress, tokenDecimals[i])
			j++
		}
	}
	return tokenDecimals, nil
}
