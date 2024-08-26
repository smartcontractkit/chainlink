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
	EIP1559      bool

	BlockHistorySize uint64
	RewardPercentile float64
	HasMempool       bool
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

	wg     *sync.WaitGroup
	stopCh services.StopChan
}

func NewFeeHistoryEstimator(lggr logger.Logger, client feeHistoryEstimatorClient, cfg FeeHistoryEstimatorConfig, chainID *big.Int, l1Oracle rollups.L1Oracle) *FeeHistoryEstimator {
	return &FeeHistoryEstimator{
		client:   client,
		logger:   logger.Named(lggr, "FeeHistoryEstimator"),
		config:   cfg,
		chainID:  chainID,
		l1Oracle: l1Oracle,
		wg:       new(sync.WaitGroup),
		stopCh:   make(chan struct{}),
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
		if f.config.EIP1559 && f.config.BlockHistorySize == 0 {
			return fmt.Errorf("BlockHistorySize is set to 0 and EIP1559 is enabled")
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

	t := services.NewTicker(f.config.CacheTimeout)
	for {
		select {
		case <-f.stopCh:
			return
		case <-t.C:
			if f.config.EIP1559 {
				if _, err := f.FetchDynamicPrice(); err != nil {
					f.logger.Error(err)
				}
			} else {
				if _, err := f.FetchGasPrice(); err != nil {
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

// FetchGasPrice will use eth_gasPrice to fetch and cache the latest gas price from the RPC.
func (f *FeeHistoryEstimator) FetchGasPrice() (*assets.Wei, error) {
	ctx, cancel := f.stopCh.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	gasPrice, err := f.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gas price: %s", err)
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

// FetchDynamicPrice uses eth_feeHistory to fetch the baseFee of the next block and the Nth maxPriorityFeePerGas percentiles
// of the past X blocks. It also fetches the highest 85th maxPriorityFeePerGas percentile of the past X blocks, which represents
// the highest percentile we're willing to pay. A buffer is added on top of the latest baseFee to catch fluctuations in the next
// blocks. On Ethereum the increase is baseFee * 1.125 per block, however in some chains that may vary.
func (f *FeeHistoryEstimator) FetchDynamicPrice() (fee DynamicFee, err error) {
	ctx, cancel := f.stopCh.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	if f.config.BlockHistorySize == 0 {
		return fee, fmt.Errorf("BlockHistorySize cannot be 0")
	}
	// RewardPercentile will be used for maxPriorityFeePerGas estimations and connectivityPercentile to set the highest threshold for bumping.
	feeHistory, err := f.client.FeeHistory(ctx, f.config.BlockHistorySize, []float64{f.config.RewardPercentile, ConnectivityPercentile})
	if err != nil {
		return fee, fmt.Errorf("failed to fetch dynamic prices: %s", err)
	}

	// eth_feeHistory doesn't return the latest baseFee of the range but rather the latest + 1, because it can be derived from the existing
	// values. Source: https://github.com/ethereum/go-ethereum/blob/b0f66e34ca2a4ea7ae23475224451c8c9a569826/eth/gasprice/feehistory.go#L235
	// nextBlock is the latest returned + 1 to be aligned with the base fee value.
	baseFee := assets.NewWei(feeHistory.BaseFee[len(feeHistory.BaseFee)-1])
	nextBlock := big.NewInt(0).Add(feeHistory.OldestBlock, big.NewInt(int64(f.config.BlockHistorySize)))

	priorityFee := big.NewInt(0)
	priorityFeeThreshold := big.NewInt(0)
	for _, fee := range feeHistory.Reward {
		priorityFee = priorityFee.Add(priorityFee, fee[0])
		// We don't need an average, we need the max value
		priorityFeeThreshold = bigmath.Max(priorityFeeThreshold, fee[1])
	}
	priorityFeeThresholdWei := assets.NewWei(priorityFeeThreshold)

	f.priorityFeeThresholdMu.Lock()
	f.priorityFeeThreshold = priorityFeeThresholdWei
	f.priorityFeeThresholdMu.Unlock()

	// eth_feeHistory may return less results than BlockHistorySize so we need to divide by the length of the result
	maxPriorityFeePerGas := assets.NewWei(priorityFee.Div(priorityFee, big.NewInt(int64(len(feeHistory.Reward)))))
	// baseFeeBufferPercentage is used as a safety to catch fluctuations in the next block.
	maxFeePerGas := baseFee.AddPercentage(BaseFeeBufferPercentage).Add(maxPriorityFeePerGas)

	promFeeHistoryEstimatorBaseFee.WithLabelValues(f.chainID.String()).Set(float64(baseFee.Int64()))
	promFeeHistoryEstimatorMaxPriorityFeePerGas.WithLabelValues(f.chainID.String()).Set(float64(maxPriorityFeePerGas.Int64()))
	promFeeHistoryEstimatorMaxFeePerGas.WithLabelValues(f.chainID.String()).Set(float64(maxFeePerGas.Int64()))

	f.logger.Debugf("Fetched new dynamic prices, nextBlock#: %v - oldestBlock#: %v - maxFeePerGas: %v - maxPriorityFeePerGas: %v - maxPriorityFeeThreshold: %v",
		nextBlock, feeHistory.OldestBlock, maxFeePerGas, maxPriorityFeePerGas, priorityFeeThresholdWei)

	f.dynamicPriceMu.Lock()
	defer f.dynamicPriceMu.Unlock()
	f.dynamicPrice.FeeCap = maxFeePerGas
	f.dynamicPrice.TipCap = maxPriorityFeePerGas
	return f.dynamicPrice, nil
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

	currentGasPrice, err := f.getGasPrice()
	if err != nil {
		return nil, 0, err
	}

	bumpedGasPrice := originalGasPrice.AddPercentage(f.config.BumpPercent)
	bumpedGasPrice, err = LimitBumpedFee(originalGasPrice, currentGasPrice, bumpedGasPrice, maxPrice)
	if err != nil {
		return nil, 0, fmt.Errorf("gas price error: %s", err.Error())
	}

	f.logger.Debugw("bumped gas price", "originalGasPrice", originalGasPrice, "bumpedGasPrice", bumpedGasPrice)

	return bumpedGasPrice, gasLimit, nil
}

// BumpDynamicFee provides a bumped dynamic fee by bumping the previous one by BumpPercent.
// If the original values are higher than the max price it returns an error as there is no room for bumping. If maxPriorityFeePerGas is bumped
// above the priority fee threshold then there is a good chance there is a connectivity issue and we shouldn't bump.
// Both maxFeePerGas as well as maxPriorityFeePerGas need to be bumped otherwise the RPC won't accept the transaction and throw an error.
// See: https://github.com/ethereum/go-ethereum/issues/24284
// It aggregates the market, bumped, and max price to provide a correct value, for both maxFeePerGas as well as maxPriorityFerPergas.
func (f *FeeHistoryEstimator) BumpDynamicFee(ctx context.Context, originalFee DynamicFee, maxPrice *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, err error) {
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

	if f.config.HasMempool {
		bumpedMaxPriorityFeePerGas, err = LimitBumpedFee(originalFee.TipCap, currentDynamicPrice.TipCap, bumpedMaxPriorityFeePerGas, maxPrice)
		if err != nil {
			return bumped, fmt.Errorf("maxPriorityFeePerGas error: %s", err.Error())
		}

		priorityFeeThreshold, e := f.getPriorityFeeThreshold()
		if e != nil {
			err = e
			return
		}

		// If either of these two values are 0 it could be that the network has extremely low priority fees. We should skip the
		// connectivity check because we're only going to be charged for the base fee, which is mandatory.
		if (priorityFeeThreshold.Cmp(assets.NewWeiI(0)) > 0) && (bumpedMaxPriorityFeePerGas.Cmp(assets.NewWeiI(0)) > 0) &&
			bumpedMaxPriorityFeePerGas.Cmp(priorityFeeThreshold) > 0 {
			return bumped, fmt.Errorf("bumpedMaxPriorityFeePergas: %s is above market's %sth percentile: %s, bumping is halted",
				bumpedMaxPriorityFeePerGas, strconv.Itoa(ConnectivityPercentile), priorityFeeThreshold)
		}
	} else {
		// If the network doesn't have a mempool then transactions are processed in a FCFS manner and maxPriorityFeePerGas value is irrelevant.
		// We just need to cap the value at maxPrice in case maxFeePerGas also gets capped.
		bumpedMaxPriorityFeePerGas = assets.WeiMin(bumpedMaxPriorityFeePerGas, maxPrice)
	}

	bumpedMaxFeePerGas, err = LimitBumpedFee(originalFee.FeeCap, currentDynamicPrice.FeeCap, bumpedMaxFeePerGas, maxPrice)
	if err != nil {
		return bumped, fmt.Errorf("maxFeePerGas error: %s", err.Error())
	}

	bumpedFee := DynamicFee{FeeCap: bumpedMaxFeePerGas, TipCap: bumpedMaxPriorityFeePerGas}
	f.logger.Debugw("bumped dynamic fee", "originalFee", originalFee, "bumpedFee", bumpedFee)

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
