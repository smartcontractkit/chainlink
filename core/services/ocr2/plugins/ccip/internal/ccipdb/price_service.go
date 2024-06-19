package db

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	cciporm "github.com/smartcontractkit/chainlink/v2/core/services/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

// PriceService manages DB access for gas and token price data.
// In the background, PriceService periodically inserts latest gas and token prices into the DB.
// During `Observation` phase, Commit plugin calls PriceService to fetch the latest prices from DB.
// This enables all lanes connected to a chain to feed price data to the leader lane's Commit plugin for that chain.
//
//go:generate mockery --quiet --name PriceService --filename price_service_mock.go --case=underscore
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
	// Prices should expire after 10 minutes in DB. Prices should be fresh in the Commit plugin.
	// 10 min provides sufficient buffer for the Commit plugin to withstand transient price update outages, while
	// surfacing price update outages quickly enough.
	priceExpireSec = 600
	// Cleanups are called every 10 minutes. For a given job, on average we may expect 3 token prices and 1 gas price.
	// 10 minutes should result in 40 rows being cleaned up per job, it is not a heavy load on DB, so there is no need
	// to run cleanup more frequently. We shouldn't clean up less frequently than `priceExpireSec`.
	priceCleanupInterval = 600 * time.Second

	// Prices are refreshed every 1 minute, they are sufficiently accurate, and consistent with Commit OCR round time.
	priceUpdateInterval = 60 * time.Second
)

type priceService struct {
	priceExpireSec  int
	cleanupInterval time.Duration
	updateInterval  time.Duration

	lggr              logger.Logger
	orm               cciporm.ORM
	jobId             int32
	destChainSelector uint64

	sourceChainSelector     uint64
	sourceNative            cciptypes.Address
	priceGetter             pricegetter.PriceGetter
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
	priceGetter pricegetter.PriceGetter,
	offRampReader ccipdata.OffRampReader,
) PriceService {
	ctx, cancel := context.WithCancel(context.Background())

	pw := &priceService{
		priceExpireSec:  priceExpireSec,
		cleanupInterval: utils.WithJitter(priceCleanupInterval), // use WithJitter to avoid multiple services impacting DB at same time
		updateInterval:  utils.WithJitter(priceUpdateInterval),

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
	cleanupTicker := time.NewTicker(p.cleanupInterval)
	updateTicker := time.NewTicker(p.updateInterval)

	go func() {
		defer p.wg.Done()

		for {
			select {
			case <-p.backgroundCtx.Done():
				return
			case <-cleanupTicker.C:
				err := p.runCleanup(p.backgroundCtx)
				if err != nil {
					p.lggr.Errorw("Error when cleaning up in-db prices in the background", "err", err)
				}
			case <-updateTicker.C:
				err := p.runUpdate(p.backgroundCtx)
				if err != nil {
					p.lggr.Errorw("Error when updating prices in the background", "err", err)
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
	if err := p.runUpdate(ctx); err != nil {
		p.lggr.Errorw("Error when updating prices after dynamic config update", "err", err)
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

func (p *priceService) runCleanup(ctx context.Context) error {
	eg := new(errgroup.Group)

	eg.Go(func() error {
		err := p.orm.ClearGasPricesByDestChain(ctx, p.destChainSelector, p.priceExpireSec)
		if err != nil {
			return fmt.Errorf("error clearing gas prices: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		err := p.orm.ClearTokenPricesByDestChain(ctx, p.destChainSelector, p.priceExpireSec)
		if err != nil {
			return fmt.Errorf("error clearing token prices: %w", err)
		}
		return nil
	})

	return eg.Wait()
}

func (p *priceService) runUpdate(ctx context.Context) error {
	// Protect against concurrent updates of `gasPriceEstimator` and `destPriceRegistryReader`
	// Price updates happen infrequently - once every `priceUpdateInterval` seconds.
	// It does not happen on any code path that is performance sensitive.
	// We can afford to have non-performant unlocks here that is simple and safe.
	p.dynamicConfigMu.RLock()
	defer p.dynamicConfigMu.RUnlock()

	// There may be a period of time between service is started and dynamic config is updated
	if p.gasPriceEstimator == nil || p.destPriceRegistryReader == nil {
		p.lggr.Info("Skipping price update due to gasPriceEstimator and/or destPriceRegistry not ready")
		return nil
	}

	sourceGasPriceUSD, tokenPricesUSD, err := p.observePriceUpdates(ctx, p.lggr)
	if err != nil {
		return fmt.Errorf("failed to observe price updates: %w", err)
	}

	err = p.writePricesToDB(ctx, sourceGasPriceUSD, tokenPricesUSD)
	if err != nil {
		return fmt.Errorf("failed to write prices to db: %w", err)
	}

	return nil
}

func (p *priceService) observePriceUpdates(
	ctx context.Context,
	lggr logger.Logger,
) (sourceGasPriceUSD *big.Int, tokenPricesUSD map[cciptypes.Address]*big.Int, err error) {
	if p.gasPriceEstimator == nil || p.destPriceRegistryReader == nil {
		return nil, nil, fmt.Errorf("gasPriceEstimator and/or destPriceRegistry is not set yet")
	}

	sortedLaneTokens, filteredLaneTokens, err := ccipcommon.GetFilteredSortedLaneTokens(ctx, p.offRampReader, p.destPriceRegistryReader, p.priceGetter)

	lggr.Debugw("Filtered bridgeable tokens with no configured price getter", "filteredLaneTokens", filteredLaneTokens)

	if err != nil {
		return nil, nil, fmt.Errorf("get destination tokens: %w", err)
	}

	return p.generatePriceUpdates(ctx, lggr, sortedLaneTokens)
}

// All prices are USD ($1=1e18) denominated. All prices must be not nil.
// Return token prices should contain the exact same tokens as in tokenDecimals.
func (p *priceService) generatePriceUpdates(
	ctx context.Context,
	lggr logger.Logger,
	sortedLaneTokens []cciptypes.Address,
) (sourceGasPriceUSD *big.Int, tokenPricesUSD map[cciptypes.Address]*big.Int, err error) {
	// Include wrapped native in our token query as way to identify the source native USD price.
	// notice USD is in 1e18 scale, i.e. $1 = 1e18
	queryTokens := ccipcommon.FlattenUniqueSlice([]cciptypes.Address{p.sourceNative}, sortedLaneTokens)

	rawTokenPricesUSD, err := p.priceGetter.TokenPricesUSD(ctx, queryTokens)
	if err != nil {
		return nil, nil, err
	}
	lggr.Infow("Raw token prices", "rawTokenPrices", rawTokenPricesUSD)

	// make sure that we got prices for all the tokens of our query
	for _, token := range queryTokens {
		if rawTokenPricesUSD[token] == nil {
			return nil, nil, fmt.Errorf("missing token price: %+v", token)
		}
	}

	sourceNativePriceUSD, exists := rawTokenPricesUSD[p.sourceNative]
	if !exists {
		return nil, nil, fmt.Errorf("missing source native (%s) price", p.sourceNative)
	}

	destTokensDecimals, err := p.destPriceRegistryReader.GetTokensDecimals(ctx, sortedLaneTokens)
	if err != nil {
		return nil, nil, fmt.Errorf("get tokens decimals: %w", err)
	}

	tokenPricesUSD = make(map[cciptypes.Address]*big.Int, len(rawTokenPricesUSD))
	for i, token := range sortedLaneTokens {
		tokenPricesUSD[token] = calculateUsdPer1e18TokenAmount(rawTokenPricesUSD[token], destTokensDecimals[i])
	}

	sourceGasPrice, err := p.gasPriceEstimator.GetGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}
	if sourceGasPrice == nil {
		return nil, nil, fmt.Errorf("missing gas price")
	}
	sourceGasPriceUSD, err = p.gasPriceEstimator.DenoteInUSD(sourceGasPrice, sourceNativePriceUSD)
	if err != nil {
		return nil, nil, err
	}

	lggr.Infow("PriceService observed latest price",
		"sourceChainSelector", p.sourceChainSelector,
		"destChainSelector", p.destChainSelector,
		"gasPriceWei", sourceGasPrice,
		"sourceNativePriceUSD", sourceNativePriceUSD,
		"sourceGasPriceUSD", sourceGasPriceUSD,
		"tokenPricesUSD", tokenPricesUSD,
	)
	return sourceGasPriceUSD, tokenPricesUSD, nil
}

func (p *priceService) writePricesToDB(
	ctx context.Context,
	sourceGasPriceUSD *big.Int,
	tokenPricesUSD map[cciptypes.Address]*big.Int,
) (err error) {
	eg := new(errgroup.Group)

	if sourceGasPriceUSD != nil {
		eg.Go(func() error {
			return p.orm.InsertGasPricesForDestChain(ctx, p.destChainSelector, p.jobId, []cciporm.GasPriceUpdate{
				{
					SourceChainSelector: p.sourceChainSelector,
					GasPrice:            assets.NewWei(sourceGasPriceUSD),
				},
			})
		})
	}

	if tokenPricesUSD != nil {
		var tokenPrices []cciporm.TokenPriceUpdate

		for token, price := range tokenPricesUSD {
			tokenPrices = append(tokenPrices, cciporm.TokenPriceUpdate{
				TokenAddr:  string(token),
				TokenPrice: assets.NewWei(price),
			})
		}

		// Sort token by addr to make price updates ordering deterministic, easier to testing and debugging
		sort.Slice(tokenPrices, func(i, j int) bool {
			return tokenPrices[i].TokenAddr < tokenPrices[j].TokenAddr
		})

		eg.Go(func() error {
			return p.orm.InsertTokenPricesForDestChain(ctx, p.destChainSelector, p.jobId, tokenPrices)
		})
	}

	return eg.Wait()
}

// Input price is USD per full token, with 18 decimal precision
// Result price is USD per 1e18 of smallest token denomination, with 18 decimal precision
// Example: 1 USDC = 1.00 USD per full token, each full token is 6 decimals -> 1 * 1e18 * 1e18 / 1e6 = 1e30
func calculateUsdPer1e18TokenAmount(price *big.Int, decimals uint8) *big.Int {
	tmp := big.NewInt(0).Mul(price, big.NewInt(1e18))
	return tmp.Div(tmp, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
}
