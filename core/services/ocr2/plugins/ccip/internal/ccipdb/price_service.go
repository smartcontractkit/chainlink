package db

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	cciporm "github.com/smartcontractkit/chainlink/v2/core/services/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// PriceService manages DB access for gas and token price data.
// In the background, PriceService periodically inserts latest gas and token prices into the DB.
// During `Observation` phase, Commit plugin calls PriceService to fetch the latest prices from DB.
// This enables all lanes connected to a chain to feed price data to the leader lane's Commit plugin for that chain.
type PriceService interface {
	job.ServiceCtx

	// UpdateDynamicConfig updates gasPriceEstimator and destPriceRegistryReader during Commit plugin dynamic config change.
	UpdateDynamicConfig(ctx context.Context, gasPriceEstimator prices.GasPriceEstimatorCommit, destPriceRegistryReader ccipdata.PriceRegistryReader) error

	// GetGasAndTokenPrices fetches source chain gas prices and relevant token prices from all lanes that touch the given dest chain.
	// The prices have been written into the DB by each lane's PriceService in the background. The prices are denoted in USD.
	GetGasAndTokenPrices(ctx context.Context, destChainSelector uint64) (map[uint64]*big.Int, map[cciptypes.Address]*big.Int, error)
}

var _ PriceService = (*priceService)(nil)

const (
	// Gas prices are refreshed every 1 minute, they are sufficiently accurate, and consistent with Commit OCR round time.
	gasPriceUpdateInterval = 1 * time.Minute
	// Token prices are refreshed every 10 minutes, we only report prices for blue chip tokens, DS&A simulation show
	// their prices are stable, 10-minute resolution is accurate enough.
	tokenPriceUpdateInterval = 10 * time.Minute
)

type priceService struct {
	gasUpdateInterval   time.Duration
	tokenUpdateInterval time.Duration

	lggr              logger.Logger
	orm               cciporm.ORM
	jobId             int32
	destChainSelector uint64

	sourceChainSelector     uint64
	sourceNative            cciptypes.Address
	priceGetter             pricegetter.AllTokensPriceGetter
	offRampReader           ccipdata.OffRampReader
	gasPriceEstimator       prices.GasPriceEstimatorCommit
	destPriceRegistryReader ccipdata.PriceRegistryReader

	services.StateMachine
	wg               *sync.WaitGroup
	backgroundCtx    context.Context //nolint:containedctx
	backgroundCancel context.CancelFunc
	dynamicConfigMu  *sync.RWMutex
}

func NewPriceService(
	lggr logger.Logger,
	orm cciporm.ORM,
	jobId int32,
	destChainSelector uint64,
	sourceChainSelector uint64,

	sourceNative cciptypes.Address,
	priceGetter pricegetter.AllTokensPriceGetter,
	offRampReader ccipdata.OffRampReader,
) PriceService {
	ctx, cancel := context.WithCancel(context.Background())

	pw := &priceService{
		gasUpdateInterval:   gasPriceUpdateInterval,
		tokenUpdateInterval: tokenPriceUpdateInterval,

		lggr:              lggr,
		orm:               orm,
		jobId:             jobId,
		destChainSelector: destChainSelector,

		sourceChainSelector: sourceChainSelector,
		sourceNative:        sourceNative,
		priceGetter:         priceGetter,
		offRampReader:       offRampReader,

		wg:               new(sync.WaitGroup),
		backgroundCtx:    ctx,
		backgroundCancel: cancel,
		dynamicConfigMu:  &sync.RWMutex{},
	}
	return pw
}

func (p *priceService) Start(context.Context) error {
	return p.StateMachine.StartOnce("PriceService", func() error {
		p.lggr.Info("Starting PriceService")
		p.wg.Add(1)
		p.run()
		return nil
	})
}

func (p *priceService) Close() error {
	return p.StateMachine.StopOnce("PriceService", func() error {
		p.lggr.Info("Closing PriceService")
		p.backgroundCancel()
		p.wg.Wait()
		return nil
	})
}

func (p *priceService) run() {
	gasUpdateTicker := time.NewTicker(utils.WithJitter(p.gasUpdateInterval))
	tokenUpdateTicker := time.NewTicker(utils.WithJitter(p.tokenUpdateInterval))

	go func() {
		defer p.wg.Done()
		defer gasUpdateTicker.Stop()
		defer tokenUpdateTicker.Stop()

		for {
			select {
			case <-p.backgroundCtx.Done():
				return
			case <-gasUpdateTicker.C:
				err := p.runGasPriceUpdate(p.backgroundCtx)
				if err != nil {
					p.lggr.Errorw("Error when updating gas prices in the background", "err", err)
				}
			case <-tokenUpdateTicker.C:
				err := p.runTokenPriceUpdate(p.backgroundCtx)
				if err != nil {
					p.lggr.Errorw("Error when updating token prices in the background", "err", err)
				}
			}
		}
	}()
}

func (p *priceService) UpdateDynamicConfig(ctx context.Context, gasPriceEstimator prices.GasPriceEstimatorCommit, destPriceRegistryReader ccipdata.PriceRegistryReader) error {
	p.dynamicConfigMu.Lock()
	p.gasPriceEstimator = gasPriceEstimator
	p.destPriceRegistryReader = destPriceRegistryReader
	p.dynamicConfigMu.Unlock()

	// Config update may substantially change the prices, refresh the prices immediately, this also makes testing easier
	// for not having to wait to the full update interval.
	if err := p.runGasPriceUpdate(ctx); err != nil {
		p.lggr.Errorw("Error when updating gas prices after dynamic config update", "err", err)
	}
	if err := p.runTokenPriceUpdate(ctx); err != nil {
		p.lggr.Errorw("Error when updating token prices after dynamic config update", "err", err)
	}

	return nil
}

func (p *priceService) GetGasAndTokenPrices(ctx context.Context, destChainSelector uint64) (map[uint64]*big.Int, map[cciptypes.Address]*big.Int, error) {
	eg := new(errgroup.Group)

	var gasPricesInDB []cciporm.GasPrice
	var tokenPricesInDB []cciporm.TokenPrice

	eg.Go(func() error {
		gasPrices, err := p.orm.GetGasPricesByDestChain(ctx, destChainSelector)
		if err != nil {
			return fmt.Errorf("failed to get gas prices from db: %w", err)
		}
		gasPricesInDB = gasPrices
		return nil
	})

	eg.Go(func() error {
		tokenPrices, err := p.orm.GetTokenPricesByDestChain(ctx, destChainSelector)
		if err != nil {
			return fmt.Errorf("failed to get token prices from db: %w", err)
		}
		tokenPricesInDB = tokenPrices
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	gasPrices := make(map[uint64]*big.Int, len(gasPricesInDB))
	tokenPrices := make(map[cciptypes.Address]*big.Int, len(tokenPricesInDB))

	for _, gasPrice := range gasPricesInDB {
		if gasPrice.GasPrice != nil {
			gasPrices[gasPrice.SourceChainSelector] = gasPrice.GasPrice.ToInt()
		}
	}

	for _, tokenPrice := range tokenPricesInDB {
		if tokenPrice.TokenPrice != nil {
			tokenPrices[cciptypes.Address(tokenPrice.TokenAddr)] = tokenPrice.TokenPrice.ToInt()
		}
	}

	return gasPrices, tokenPrices, nil
}

func (p *priceService) runGasPriceUpdate(ctx context.Context) error {
	// Protect against concurrent updates of `gasPriceEstimator` and `destPriceRegistryReader`
	// Price updates happen infrequently - once every `gasPriceUpdateInterval` seconds.
	// It does not happen on any code path that is performance sensitive.
	// We can afford to have non-performant unlocks here that is simple and safe.
	p.dynamicConfigMu.RLock()
	defer p.dynamicConfigMu.RUnlock()

	// There may be a period of time between service is started and dynamic config is updated
	if p.gasPriceEstimator == nil {
		p.lggr.Info("Skipping gas price update due to gasPriceEstimator not ready")
		return nil
	}

	sourceGasPriceUSD, err := p.observeGasPriceUpdates(ctx, p.lggr)
	if err != nil {
		return fmt.Errorf("failed to observe gas price updates: %w", err)
	}

	err = p.writeGasPricesToDB(ctx, sourceGasPriceUSD)
	if err != nil {
		return fmt.Errorf("failed to write gas prices to db: %w", err)
	}

	return nil
}

func (p *priceService) runTokenPriceUpdate(ctx context.Context) error {
	// Protect against concurrent updates of `tokenPriceEstimator` and `destPriceRegistryReader`
	// Price updates happen infrequently - once every `tokenPriceUpdateInterval` seconds.
	p.dynamicConfigMu.RLock()
	defer p.dynamicConfigMu.RUnlock()

	// There may be a period of time between service is started and dynamic config is updated
	if p.destPriceRegistryReader == nil {
		p.lggr.Info("Skipping token price update due to destPriceRegistry not ready")
		return nil
	}

	tokenPricesUSD, err := p.observeTokenPriceUpdates(ctx, p.lggr)
	if err != nil {
		return fmt.Errorf("failed to observe token price updates: %w", err)
	}

	err = p.writeTokenPricesToDB(ctx, tokenPricesUSD)
	if err != nil {
		return fmt.Errorf("failed to write token prices to db: %w", err)
	}

	return nil
}

func (p *priceService) observeGasPriceUpdates(
	ctx context.Context,
	lggr logger.Logger,
) (sourceGasPriceUSD *big.Int, err error) {
	if p.gasPriceEstimator == nil {
		return nil, fmt.Errorf("gasPriceEstimator is not set yet")
	}

	// Include wrapped native to identify the source native USD price, notice USD is in 1e18 scale, i.e. $1 = 1e18
	rawTokenPricesUSD, err := p.priceGetter.TokenPricesUSD(ctx, []cciptypes.Address{p.sourceNative})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch source native price (%s): %w", p.sourceNative, err)
	}

	sourceNativePriceUSD, exists := rawTokenPricesUSD[p.sourceNative]
	if !exists {
		return nil, fmt.Errorf("missing source native (%s) price", p.sourceNative)
	}

	sourceGasPrice, err := p.gasPriceEstimator.GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	if sourceGasPrice == nil {
		return nil, fmt.Errorf("missing gas price")
	}
	sourceGasPriceUSD, err = p.gasPriceEstimator.DenoteInUSD(sourceGasPrice, sourceNativePriceUSD)
	if err != nil {
		return nil, err
	}

	lggr.Infow("PriceService observed latest gas price",
		"sourceChainSelector", p.sourceChainSelector,
		"destChainSelector", p.destChainSelector,
		"sourceNative", p.sourceNative,
		"gasPriceWei", sourceGasPrice,
		"sourceNativePriceUSD", sourceNativePriceUSD,
		"sourceGasPriceUSD", sourceGasPriceUSD,
	)
	return sourceGasPriceUSD, nil
}

// All prices are USD ($1=1e18) denominated. All prices must be not nil.
// Jobspec should have the destination tokens (Aggregate Rate Limit, Bps) and 1 source token (source native).
// Not respecting this will error out as we need to fetch the token decimals for all tokens expect sourceNative.
// destTokens is only used to check if sourceNative has the same address as one of the dest tokens.
// Return token prices should contain the exact same tokens as in tokenDecimals.
func (p *priceService) observeTokenPriceUpdates(
	ctx context.Context,
	lggr logger.Logger,
) (tokenPricesUSD map[cciptypes.Address]*big.Int, err error) {
	if p.destPriceRegistryReader == nil {
		return nil, fmt.Errorf("destPriceRegistry is not set yet")
	}
	rawTokenPricesUSD, err := p.priceGetter.GetJobSpecTokenPricesUSD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token prices: %w", err)
	}

	// Verify no price returned by price getter is nil
	for token, price := range rawTokenPricesUSD {
		if price == nil {
			return nil, fmt.Errorf("Token price is nil for token %s", token)
		}
	}

	lggr.Infow("Raw token prices", "rawTokenPrices", rawTokenPricesUSD)

	sourceNativeEvmAddr, err := ccipcalc.GenericAddrToEvm(p.sourceNative)
	if err != nil {
		return nil, fmt.Errorf("failed to convert source native to EVM address: %w", err)
	}

	// Filter out source native token only if source native not in dest tokens
	var finalDestTokens []cciptypes.Address
	for token := range rawTokenPricesUSD {
		tokenEvmAddr, err2 := ccipcalc.GenericAddrToEvm(token)
		if err2 != nil {
			return nil, fmt.Errorf("failed to convert token to EVM address: %w", err)
		}

		if tokenEvmAddr != sourceNativeEvmAddr {
			finalDestTokens = append(finalDestTokens, token)
		}
	}

	fee, bridged, err := ccipcommon.GetDestinationTokens(ctx, p.offRampReader, p.destPriceRegistryReader)
	if err != nil {
		return nil, fmt.Errorf("get destination tokens: %w", err)
	}
	onchainDestTokens := ccipcommon.FlattenedAndSortedTokens(fee, bridged)
	lggr.Debugw("Destination tokens", "destTokens", onchainDestTokens)

	onchainTokensEvmAddr, err := ccipcalc.GenericAddrsToEvm(onchainDestTokens...)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sorted lane tokens to EVM addresses: %w", err)
	}
	// Check for case where sourceNative has same address as one of the dest tokens (example: WETH in Base and Optimism)
	hasSameDestAddress := slices.Contains(onchainTokensEvmAddr, sourceNativeEvmAddr)

	if hasSameDestAddress {
		finalDestTokens = append(finalDestTokens, p.sourceNative)
	}

	// Sort tokens to make the order deterministic, easier for testing and debugging
	sort.Slice(finalDestTokens, func(i, j int) bool {
		return finalDestTokens[i] < finalDestTokens[j]
	})

	destTokensDecimals, err := p.destPriceRegistryReader.GetTokensDecimals(ctx, finalDestTokens)
	if err != nil {
		return nil, fmt.Errorf("get tokens decimals: %w", err)
	}

	if len(destTokensDecimals) != len(finalDestTokens) {
		return nil, fmt.Errorf("mismatched token decimals and tokens")
	}

	tokenPricesUSD = make(map[cciptypes.Address]*big.Int, len(rawTokenPricesUSD))
	for i, token := range finalDestTokens {
		tokenPricesUSD[token] = calculateUsdPer1e18TokenAmount(rawTokenPricesUSD[token], destTokensDecimals[i])
	}

	lggr.Infow("PriceService observed latest token prices",
		"sourceChainSelector", p.sourceChainSelector,
		"destChainSelector", p.destChainSelector,
		"tokenPricesUSD", tokenPricesUSD,
	)
	return tokenPricesUSD, nil
}

func (p *priceService) writeGasPricesToDB(ctx context.Context, sourceGasPriceUSD *big.Int) error {
	if sourceGasPriceUSD == nil {
		return nil
	}

	_, err := p.orm.UpsertGasPricesForDestChain(ctx, p.destChainSelector, []cciporm.GasPrice{
		{
			SourceChainSelector: p.sourceChainSelector,
			GasPrice:            assets.NewWei(sourceGasPriceUSD),
		},
	})
	return err
}

func (p *priceService) writeTokenPricesToDB(ctx context.Context, tokenPricesUSD map[cciptypes.Address]*big.Int) error {
	if tokenPricesUSD == nil {
		return nil
	}

	var tokenPrices []cciporm.TokenPrice

	for token, price := range tokenPricesUSD {
		tokenPrices = append(tokenPrices, cciporm.TokenPrice{
			TokenAddr:  string(token),
			TokenPrice: assets.NewWei(price),
		})
	}

	// Sort token by addr to make price updates ordering deterministic, easier for testing and debugging
	sort.Slice(tokenPrices, func(i, j int) bool {
		return tokenPrices[i].TokenAddr < tokenPrices[j].TokenAddr
	})

	_, err := p.orm.UpsertTokenPricesForDestChain(ctx, p.destChainSelector, tokenPrices, p.tokenUpdateInterval)
	return err
}

// Input price is USD per full token, with 18 decimal precision
// Result price is USD per 1e18 of smallest token denomination, with 18 decimal precision
// Example: 1 USDC = 1.00 USD per full token, each full token is 6 decimals -> 1 * 1e18 * 1e18 / 1e6 = 1e30
func calculateUsdPer1e18TokenAmount(price *big.Int, decimals uint8) *big.Int {
	tmp := big.NewInt(0).Mul(price, big.NewInt(1e18))
	return tmp.Div(tmp, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
}
