package gas

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// metrics are thread safe
var (
	promFeeHistoryEstimatorGasPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_price_updater",
		Help: "Sets latest gas price (in Wei)",
	},
		[]string{"evmChainID"},
	)
	promFeeHistoryEstimatorBaseFee = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "base_fee_updater",
		Help: "Sets latest BaseFee (in Wei)",
	},
		[]string{"evmChainID"},
	)
	promFeeHistoryEstimatorMaxPriorityFeePerGas = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_priority_fee_per_gas_updater",
		Help: "Sets latest MaxPriorityFeePerGas (in Wei)",
	},
		[]string{"evmChainID"},
	)
	promFeeHistoryEstimatorMaxFeePerGas = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_fee_per_gas_updater",
		Help: "Sets latest MaxFeePerGas (in Wei)",
	},
		[]string{"evmChainID"},
	)
)

const (
	MinimumBumpPercentage   = 10 // based on geth's spec
	ConnectivityPercentile  = 85
	BaseFeeBufferPercentage = 40
)

type FeeHistoryEstimatorConfig struct {
	BumpPercent  uint16
	CacheTimeout time.Duration

	EIP1559          bool
	BlockHistorySize uint64
	RewardPercentile float64
}

type feeHistoryEstimatorClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	FeeHistory(ctx context.Context, blockCount uint64, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory, err error)
}

type FeeHistoryEstimator struct {
	services.StateMachine

	client  feeHistoryEstimatorClient
	logger  logger.Logger
	config  FeeHistoryEstimatorConfig
	chainID *big.Int

	gasPriceMu sync.RWMutex
	gasPrice   *assets.Wei

	dynamicPriceMu sync.RWMutex
	dynamicPrice   DynamicFee

	priorityFeeThresholdMu sync.RWMutex
	priorityFeeThreshold   *assets.Wei

	l1Oracle rollups.L1Oracle

	wg        *sync.WaitGroup
	stopCh    services.StopChan
	refreshCh chan struct{}
}

func NewFeeHistoryEstimator(lggr logger.Logger, client feeHistoryEstimatorClient, cfg FeeHistoryEstimatorConfig, chainID *big.Int, l1Oracle rollups.L1Oracle) *FeeHistoryEstimator {
	return &FeeHistoryEstimator{
		client:    client,
		logger:    logger.Named(lggr, "FeeHistoryEstimator"),
		config:    cfg,
		chainID:   chainID,
		l1Oracle:  l1Oracle,
		wg:        new(sync.WaitGroup),
		stopCh:    make(chan struct{}),
		refreshCh: make(chan struct{}),
	}
}

func (f *FeeHistoryEstimator) Start(context.Context) error {
	return f.StartOnce("FeeHistoryEstimator", func() error {
		if f.config.BumpPercent < MinimumBumpPercentage {
			return fmt.Errorf("BumpPercent: %s is less than minimum allowed percentage: %s",
				strconv.FormatUint(uint64(f.config.BumpPercent), 10), strconv.Itoa(MinimumBumpPercentage))
		}
		if f.config.EIP1559 && f.config.RewardPercentile > ConnectivityPercentile {
			return fmt.Errorf("RewardPercentile: %s is greater than maximum allowed percentile: %s",
				strconv.FormatUint(uint64(f.config.RewardPercentile), 10), strconv.Itoa(ConnectivityPercentile))
		}
		f.wg.Add(1)
		go f.run()

		return nil
	})
}

func (f *FeeHistoryEstimator) Close() error {
	return f.StopOnce("FeeHistoryEstimator", func() error {
		close(f.stopCh)
		f.wg.Wait()
		return nil
	})
}

func (f *FeeHistoryEstimator) run() {
	defer f.wg.Done()

	t := services.TickerConfig{
		JitterPct: services.DefaultJitter,
	}.NewTicker(f.config.CacheTimeout)

	for {
		select {
		case <-f.stopCh:
			return
		case <-f.refreshCh:
			t.Reset()
		case <-t.C:
			if f.config.EIP1559 {
				if err := f.RefreshDynamicPrice(); err != nil {
					f.logger.Error(err)
				}
			} else {
				if _, err := f.RefreshGasPrice(); err != nil {
					f.logger.Error(err)
				}
			}
		}
	}
}

// GetLegacyGas will fetch the cached gas price value.
func (f *FeeHistoryEstimator) GetLegacyGas(ctx context.Context, _ []byte, gasLimit uint64, maxPrice *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	chainSpecificGasLimit = gasLimit
	if gasPrice, err = f.getGasPrice(); err != nil {
		return
	}

	if gasPrice.Cmp(maxPrice) > 0 {
		f.logger.Warnf("estimated gas price: %s is greater than the maximum gas price configured: %s, returning the maximum price instead.", gasPrice, maxPrice)
		return maxPrice, chainSpecificGasLimit, nil
	}
	return
}

// RefreshGasPrice will use eth_gasPrice to fetch and cache the latest gas price from the RPC.
func (f *FeeHistoryEstimator) RefreshGasPrice() (*assets.Wei, error) {
	ctx, cancel := f.stopCh.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	gasPrice, err := f.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	promFeeHistoryEstimatorGasPrice.WithLabelValues(f.chainID.String()).Set(float64(gasPrice.Int64()))

	gasPriceWei := assets.NewWei(gasPrice)

	f.logger.Debugf("fetched new gas price: %v", gasPriceWei)

	f.gasPriceMu.Lock()
	defer f.gasPriceMu.Unlock()
	f.gasPrice = gasPriceWei
	return f.gasPrice, nil
}

func (f *FeeHistoryEstimator) getGasPrice() (*assets.Wei, error) {
	f.gasPriceMu.RLock()
	defer f.gasPriceMu.RUnlock()
	if f.gasPrice == nil {
		return f.gasPrice, fmt.Errorf("gas price not set")
	}
	return f.gasPrice, nil
}

// GetDynamicFee will fetch the cached dynamic prices.
func (f *FeeHistoryEstimator) GetDynamicFee(ctx context.Context, maxPrice *assets.Wei) (fee DynamicFee, err error) {
	if fee, err = f.getDynamicPrice(); err != nil {
		return
	}

	if fee.FeeCap.Cmp(maxPrice) > 0 {
		f.logger.Warnf("estimated maxFeePerGas: %v is greater than the maximum price configured: %v, returning the maximum price instead.",
			fee.FeeCap, maxPrice)
		fee.FeeCap = maxPrice
		if fee.TipCap.Cmp(maxPrice) > 0 {
			f.logger.Warnf("estimated maxPriorityFeePerGas: %v is greater than the maximum price configured: %v, returning the maximum price instead.",
				fee.TipCap, maxPrice)
			fee.TipCap = maxPrice
		}
	}

	return
}

// RefreshDynamicPrice uses eth_feeHistory to fetch the baseFee of the next block and the Nth maxPriorityFeePerGas percentiles
// of the past X blocks. It also fetches the highest 85th maxPriorityFeePerGas percentile of the past X blocks, which represents
// the highest percentile we're willing to pay. A buffer is added on top of the latest baseFee to catch fluctuations in the next
// blocks. On Ethereum the increase is baseFee * 1.125 per block, however in some chains that may vary.
func (f *FeeHistoryEstimator) RefreshDynamicPrice() error {
	ctx, cancel := f.stopCh.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	// RewardPercentile will be used for maxPriorityFeePerGas estimations and connectivityPercentile to set the highest threshold for bumping.
	feeHistory, err := f.client.FeeHistory(ctx, max(f.config.BlockHistorySize, 1), []float64{f.config.RewardPercentile, ConnectivityPercentile})
	if err != nil {
		return err
	}

	// eth_feeHistory doesn't return the latest baseFee of the range but rather the latest + 1, because it can be derived from the existing
	// values. Source: https://github.com/ethereum/go-ethereum/blob/b0f66e34ca2a4ea7ae23475224451c8c9a569826/eth/gasprice/feehistory.go#L235
	// nextBlock is the latest returned + 1 to be aligned with the base fee value.
	nextBaseFee := assets.NewWei(feeHistory.BaseFee[len(feeHistory.BaseFee)-1])
	nextBlock := big.NewInt(0).Add(feeHistory.OldestBlock, big.NewInt(int64(f.config.BlockHistorySize)))

	// If BlockHistorySize is 0 it means priority fees will be ignored from the calculations, so we set them to 0.
	// If it's not we exclude 0 priced priority fees from the RPC response, even though some networks allow them. For empty blocks, eth_feeHistory
	// returns priority fees with 0 values so it's safer to discard them in order to pick values from a more representative sample.
	maxPriorityFeePerGas := assets.NewWeiI(0)
	priorityFeeThresholdWei := assets.NewWeiI(0)
	if f.config.BlockHistorySize > 0 {
		var nonZeroRewardsLen int64
		priorityFee := big.NewInt(0)
		priorityFeeThreshold := big.NewInt(0)
		for _, reward := range feeHistory.Reward {
			// reward needs to have values for two percentiles
			if len(reward) < 2 {
				return fmt.Errorf("reward size incorrect: %d", len(reward))
			}
			// We'll calculate the average of non-zero priority fees
			if reward[0].Cmp(big.NewInt(0)) > 0 {
				priorityFee = priorityFee.Add(priorityFee, reward[0])
				nonZeroRewardsLen++
			}
			// We take the max value for the bumping threshold
			if reward[1].Cmp(big.NewInt(0)) > 0 {
				priorityFeeThreshold = bigmath.Max(priorityFeeThreshold, reward[1])
			}
		}

		if nonZeroRewardsLen == 0 || priorityFeeThreshold.Cmp(big.NewInt(0)) == 0 {
			return nil
		}
		priorityFeeThresholdWei = assets.NewWei(priorityFeeThreshold)
		maxPriorityFeePerGas = assets.NewWei(priorityFee.Div(priorityFee, big.NewInt(nonZeroRewardsLen)))
	}
	// BaseFeeBufferPercentage is used as a safety to catch any fluctuations in the Base Fee during the next blocks.
	maxFeePerGas := nextBaseFee.AddPercentage(BaseFeeBufferPercentage).Add(maxPriorityFeePerGas)

	promFeeHistoryEstimatorBaseFee.WithLabelValues(f.chainID.String()).Set(float64(nextBaseFee.Int64()))
	promFeeHistoryEstimatorMaxPriorityFeePerGas.WithLabelValues(f.chainID.String()).Set(float64(maxPriorityFeePerGas.Int64()))
	promFeeHistoryEstimatorMaxFeePerGas.WithLabelValues(f.chainID.String()).Set(float64(maxFeePerGas.Int64()))

	f.logger.Debugf("Fetched new dynamic prices, nextBlock#: %v - oldestBlock#: %v - nextBaseFee: %v - maxFeePerGas: %v - maxPriorityFeePerGas: %v - maxPriorityFeeThreshold: %v",
		nextBlock, feeHistory.OldestBlock, nextBaseFee, maxFeePerGas, maxPriorityFeePerGas, priorityFeeThresholdWei)

	f.priorityFeeThresholdMu.Lock()
	f.priorityFeeThreshold = priorityFeeThresholdWei
	f.priorityFeeThresholdMu.Unlock()

	f.dynamicPriceMu.Lock()
	defer f.dynamicPriceMu.Unlock()
	f.dynamicPrice.FeeCap = maxFeePerGas
	f.dynamicPrice.TipCap = maxPriorityFeePerGas
	return nil
}

func (f *FeeHistoryEstimator) getDynamicPrice() (fee DynamicFee, err error) {
	f.dynamicPriceMu.RLock()
	defer f.dynamicPriceMu.RUnlock()
	if f.dynamicPrice.FeeCap == nil || f.dynamicPrice.TipCap == nil {
		return fee, fmt.Errorf("dynamic price not set")
	}
	return f.dynamicPrice, nil
}

// BumpLegacyGas provides a bumped gas price value by bumping the previous one by BumpPercent.
// If the original value is higher than the max price it returns an error as there is no room for bumping.
// It aggregates the market, bumped, and max gas price to provide a correct value.
func (f *FeeHistoryEstimator) BumpLegacyGas(ctx context.Context, originalGasPrice *assets.Wei, gasLimit uint64, maxPrice *assets.Wei, _ []EvmPriorAttempt) (*assets.Wei, uint64, error) {
	// Sanitize original fee input
	if originalGasPrice == nil || originalGasPrice.Cmp(maxPrice) >= 0 {
		return nil, 0, fmt.Errorf("%w: error while retrieving original gas price: originalGasPrice: %s. Maximum price configured: %s",
			commonfee.ErrBump, originalGasPrice, maxPrice)
	}

	currentGasPrice, err := f.RefreshGasPrice()
	if err != nil {
		return nil, 0, err
	}
	f.IfStarted(func() { f.refreshCh <- struct{}{} })

	bumpedGasPrice := originalGasPrice.AddPercentage(f.config.BumpPercent)
	bumpedGasPrice, err = LimitBumpedFee(originalGasPrice, currentGasPrice, bumpedGasPrice, maxPrice)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to limit gas price: %w", err)
	}

	f.logger.Debugw("bumped gas price", "originalGasPrice", originalGasPrice, "marketGasPrice", currentGasPrice, "bumpedGasPrice", bumpedGasPrice)

	return bumpedGasPrice, gasLimit, nil
}

// BumpDynamicFee provides a bumped dynamic fee by bumping the previous one by BumpPercent.
// If the original values are higher than the max price it returns an error as there is no room for bumping. If maxPriorityFeePerGas is bumped
// above the priority fee threshold then there is a good chance there is a connectivity issue and we shouldn't bump.
// Both maxFeePerGas as well as maxPriorityFeePerGas need to be bumped otherwise the RPC won't accept the transaction and throw an error.
// See: https://github.com/ethereum/go-ethereum/issues/24284
// It aggregates the market, bumped, and max price to provide a correct value, for both maxFeePerGas as well as maxPriorityFerPergas.
func (f *FeeHistoryEstimator) BumpDynamicFee(ctx context.Context, originalFee DynamicFee, maxPrice *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, err error) {
	// For chains that don't have a mempool there is no concept of gas bumping so we force-call RefreshDynamicPrice to update the underlying base fee value
	if f.config.BlockHistorySize == 0 {
		if !f.IfStarted(func() {
			if refreshErr := f.RefreshDynamicPrice(); refreshErr != nil {
				err = refreshErr
				return
			}
			f.refreshCh <- struct{}{}
			bumped, err = f.GetDynamicFee(ctx, maxPrice)
		}) {
			return bumped, fmt.Errorf("estimator not started")
		}
		return bumped, err
	}

	// Sanitize original fee input
	// According to geth's spec we need to bump both maxFeePerGas and maxPriorityFeePerGas for the new attempt to be accepted by the RPC
	if originalFee.FeeCap == nil ||
		originalFee.TipCap == nil ||
		((originalFee.TipCap.Cmp(originalFee.FeeCap)) > 0) ||
		(originalFee.FeeCap.Cmp(maxPrice) >= 0) {
		return bumped, fmt.Errorf("%w: error while retrieving original dynamic fees: (originalFeePerGas: %s - originalPriorityFeePerGas: %s). Maximum price configured: %s",
			commonfee.ErrBump, originalFee.FeeCap, originalFee.TipCap, maxPrice)
	}

	currentDynamicPrice, err := f.getDynamicPrice()
	if err != nil {
		return
	}

	bumpedMaxPriorityFeePerGas := originalFee.TipCap.AddPercentage(f.config.BumpPercent)
	bumpedMaxFeePerGas := originalFee.FeeCap.AddPercentage(f.config.BumpPercent)

	bumpedMaxPriorityFeePerGas, err = LimitBumpedFee(originalFee.TipCap, currentDynamicPrice.TipCap, bumpedMaxPriorityFeePerGas, maxPrice)
	if err != nil {
		return bumped, fmt.Errorf("failed to limit maxPriorityFeePerGas: %w", err)
	}

	priorityFeeThreshold, e := f.getPriorityFeeThreshold()
	if e != nil {
		return bumped, e
	}

	if bumpedMaxPriorityFeePerGas.Cmp(priorityFeeThreshold) > 0 {
		return bumped, fmt.Errorf("bumpedMaxPriorityFeePerGas: %s is above market's %sth percentile: %s, bumping is halted",
			bumpedMaxPriorityFeePerGas, strconv.Itoa(ConnectivityPercentile), priorityFeeThreshold)
	}

	bumpedMaxFeePerGas, err = LimitBumpedFee(originalFee.FeeCap, currentDynamicPrice.FeeCap, bumpedMaxFeePerGas, maxPrice)
	if err != nil {
		return bumped, fmt.Errorf("failed to limit maxFeePerGas: %w", err)
	}

	bumpedFee := DynamicFee{FeeCap: bumpedMaxFeePerGas, TipCap: bumpedMaxPriorityFeePerGas}
	f.logger.Debugw("bumped dynamic fee", "originalFee", originalFee, "marketFee", currentDynamicPrice, "bumpedFee", bumpedFee)

	return bumpedFee, nil
}

// LimitBumpedFee selects the maximum value between the bumped attempt and the current fee, if there is one. If the result is higher than the max price it gets capped.
// Geth's implementation has a hard 10% minimum limit for the bumped values, otherwise it rejects the transaction with an error.
// See: https://github.com/ethereum/go-ethereum/blob/bff330335b94af3643ac2fb809793f77de3069d4/core/tx_list.go#L298
//
// Note: for chains that support EIP-1559 but we still choose to send Legacy transactions to them, the limit is still enforcable due to the fact that Legacy transactions
// are treated the same way as Dynamic transactions under the hood. For chains that don't support EIP-1559 at all, the limit isn't enforcable but a 10% minimum bump percentage
// makes sense anyway.
func LimitBumpedFee(originalFee *assets.Wei, currentFee *assets.Wei, bumpedFee *assets.Wei, maxPrice *assets.Wei) (*assets.Wei, error) {
	if currentFee != nil {
		bumpedFee = assets.WeiMax(currentFee, bumpedFee)
	}
	bumpedFee = assets.WeiMin(bumpedFee, maxPrice)

	// The first check is added for the following edge case:
	// If originalFee is below 10 wei, then adding the minimum bump percentage won't have any effect on the final value because of rounding down.
	// Similarly for bumpedFee, it can have the exact same value as the originalFee, even if we bumped, given an originalFee of less than 10 wei
	// and a small enough BumpPercent.
	if bumpedFee.Cmp(originalFee) == 0 ||
		bumpedFee.Cmp(originalFee.AddPercentage(MinimumBumpPercentage)) < 0 {
		return nil, fmt.Errorf("%w: %s is bumped less than minimum allowed percentage(%s) from originalFee: %s - maxPrice: %s",
			commonfee.ErrBump, bumpedFee, strconv.Itoa(MinimumBumpPercentage), originalFee, maxPrice)
	}
	return bumpedFee, nil
}

func (f *FeeHistoryEstimator) getPriorityFeeThreshold() (*assets.Wei, error) {
	f.priorityFeeThresholdMu.RLock()
	defer f.priorityFeeThresholdMu.RUnlock()
	if f.priorityFeeThreshold == nil {
		return f.priorityFeeThreshold, fmt.Errorf("priorityFeeThreshold not set")
	}
	return f.priorityFeeThreshold, nil
}

func (f *FeeHistoryEstimator) Name() string                                      { return f.logger.Name() }
func (f *FeeHistoryEstimator) L1Oracle() rollups.L1Oracle                        { return f.l1Oracle }
func (f *FeeHistoryEstimator) HealthReport() map[string]error                    { return map[string]error{f.Name(): nil} }
func (f *FeeHistoryEstimator) OnNewLongestChain(context.Context, *evmtypes.Head) {}
