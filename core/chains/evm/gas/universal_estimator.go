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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// metrics are thread safe
var (
	promUniversalEstimatorGasPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_price_updater",
		Help: "Sets latest gas price (in Wei)",
	},
		[]string{"evmChainID"},
	)
	promUniversalEstimatorBaseFee = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "base_fee_updater",
		Help: "Sets latest BaseFee (in Wei)",
	},
		[]string{"evmChainID"},
	)
	promUniversalEstimatorMaxPriorityFeePerGas = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_priority_fee_per_gas_updater",
		Help: "Sets latest MaxPriorityFeePerGas (in Wei)",
	},
		[]string{"evmChainID"},
	)
	promUniversalEstimatorMaxFeePerGas = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_fee_per_gas_updater",
		Help: "Sets latest MaxFeePerGas (in Wei)",
	},
		[]string{"evmChainID"},
	)
)

const (
	queryTimeout = 10 * time.Second

	MinimumBumpPercentage   = 10 // based on geth's spec
	ConnectivityPercentile  = 80
	BaseFeeBufferPercentage = 40
)

type UniversalEstimatorConfig struct {
	BumpPercent  uint16
	CacheTimeout time.Duration
	EIP1559      bool

	BlockHistorySize uint64
	RewardPercentile float64
	HasMempool       bool
}

type universalEstimatorClient interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	FeeHistory(ctx context.Context, blockCount uint64, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory, err error)
}

type UniversalEstimator struct {
	services.StateMachine

	client  universalEstimatorClient
	logger  logger.Logger
	config  UniversalEstimatorConfig
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

func NewUniversalEstimator(lggr logger.Logger, client universalEstimatorClient, cfg UniversalEstimatorConfig, chainID *big.Int, l1Oracle rollups.L1Oracle) *UniversalEstimator {
	return &UniversalEstimator{
		client:   client,
		logger:   logger.Named(lggr, "UniversalEstimator"),
		config:   cfg,
		chainID:  chainID,
		l1Oracle: l1Oracle,
		wg:       new(sync.WaitGroup),
		stopCh:   make(chan struct{}),
	}
}

func (u *UniversalEstimator) Start(context.Context) error {
	return u.StartOnce("UniversalEstimator", func() error {
		if u.config.BumpPercent < MinimumBumpPercentage {
			return fmt.Errorf("BumpPercent: %s is less than minimum allowed percentage: %s",
				strconv.FormatUint(uint64(u.config.BumpPercent), 10), strconv.Itoa(MinimumBumpPercentage))
		}
		if u.config.EIP1559 && u.config.RewardPercentile > ConnectivityPercentile {
			return fmt.Errorf("RewardPercentile: %s is greater than maximum allowed connectivity percentage: %s",
				strconv.FormatUint(uint64(u.config.RewardPercentile), 10), strconv.Itoa(ConnectivityPercentile))
		}
		if u.config.EIP1559 && u.config.BlockHistorySize == 0 {
			return fmt.Errorf("BlockHistorySize is set to 0 and EIP1559 is enabled")
		}
		u.wg.Add(1)
		go u.run()

		return nil
	})
}

func (u *UniversalEstimator) Close() error {
	return u.StopOnce("UniversalEstimator", func() error {
		close(u.stopCh)
		u.wg.Wait()
		return nil
	})
}

func (u *UniversalEstimator) run() {
	defer u.wg.Done()

	ctx, cancel := u.stopCh.NewCtx()
	defer cancel()

	t := services.NewTicker(u.config.CacheTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if u.config.EIP1559 {
				if _, err := u.FetchDynamicPrice(ctx); err != nil {
					u.logger.Error(err)
				}
			} else {
				if _, err := u.FetchGasPrice(ctx); err != nil {
					u.logger.Error(err)
				}
			}
		}
	}
}

// GetLegacyGas will fetch the cached gas price value.
func (u *UniversalEstimator) GetLegacyGas(ctx context.Context, _ []byte, gasLimit uint64, maxPrice *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	chainSpecificGasLimit = gasLimit
	if gasPrice, err = u.getGasPrice(); err != nil {
		return
	}

	if gasPrice.Cmp(maxPrice) > 0 {
		u.logger.Warnf("estimated gas price: %s is greater than the maximum gas price configured: %s, returning the maximum price instead.", gasPrice, maxPrice)
		return maxPrice, chainSpecificGasLimit, nil
	}
	return
}

// FetchGasPrice will use eth_gasPrice to fetch and cache the latest gas price from the RPC.
func (u *UniversalEstimator) FetchGasPrice(parentCtx context.Context) (*assets.Wei, error) {
	ctx, cancel := context.WithTimeout(parentCtx, queryTimeout)
	defer cancel()

	gasPrice, err := u.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gas price: %s", err)
	}

	promUniversalEstimatorGasPrice.WithLabelValues(u.chainID.String()).Set(float64(gasPrice.Int64()))

	gasPriceWei := (*assets.Wei)(gasPrice)

	u.logger.Debugf("fetched new gas price: %v", gasPriceWei)

	u.gasPriceMu.Lock()
	defer u.gasPriceMu.Unlock()
	u.gasPrice = gasPriceWei
	return u.gasPrice, nil
}

func (u *UniversalEstimator) getGasPrice() (*assets.Wei, error) {
	u.gasPriceMu.RLock()
	defer u.gasPriceMu.RUnlock()
	if u.gasPrice == nil {
		return u.gasPrice, fmt.Errorf("gas price not set")
	}
	return u.gasPrice, nil
}

// GetDynamicFee will fetch the cached dynamic prices.
func (u *UniversalEstimator) GetDynamicFee(ctx context.Context, maxPrice *assets.Wei) (fee DynamicFee, err error) {
	if fee, err = u.getDynamicPrice(); err != nil {
		return
	}

	if fee.FeeCap == nil || fee.TipCap == nil {
		return fee, fmt.Errorf("dynamic price not set")
	}
	if fee.FeeCap.Cmp(maxPrice) > 0 {
		u.logger.Warnf("estimated maxFeePerGas: %v is greater than the maximum price configured: %v, returning the maximum price instead.",
			fee.FeeCap, maxPrice)
		fee.FeeCap = maxPrice
		if fee.TipCap.Cmp(maxPrice) > 0 {
			u.logger.Warnf("estimated maxPriorityFeePerGas: %v is greater than the maximum price configured: %v, returning the maximum price instead. There won't be any room for base fee!",
				fee.TipCap, maxPrice)
			fee.TipCap = maxPrice
		}
	}

	return
}

// FetchDynamicPrice uses eth_feeHistory to fetch the basFee of the latest block and the Nth maxPriorityFeePerGas percentiles
// of the past X blocks. It also fetches the highest Zth maxPriorityFeePerGas percentile of the past X blocks. Z is configurable
// and it represents the highest percentile we're willing to pay.
// A buffer is added on top of the latest basFee to catch fluctuations in the next blocks. On Ethereum the increase is baseFee*1.125 per block.
func (u *UniversalEstimator) FetchDynamicPrice(parentCtx context.Context) (fee DynamicFee, err error) {
	ctx, cancel := context.WithTimeout(parentCtx, queryTimeout)
	defer cancel()

	if u.config.BlockHistorySize == 0 {
		return fee, fmt.Errorf("BlockHistorySize cannot be 0")
	}
	// RewardPercentile will be used for maxPriorityFeePerGas estimations and connectivityPercentile to set the highest threshold for bumping.
	feeHistory, err := u.client.FeeHistory(ctx, u.config.BlockHistorySize, []float64{u.config.RewardPercentile, ConnectivityPercentile})
	if err != nil {
		return fee, fmt.Errorf("failed to fetch dynamic prices: %s", err)
	}

	// Latest base fee
	baseFee := (*assets.Wei)(feeHistory.BaseFee[len(feeHistory.BaseFee)-1])
	latestBlock := big.NewInt(0).Add(feeHistory.OldestBlock, big.NewInt(int64(u.config.BlockHistorySize)))

	priorityFee := big.NewInt(0)
	priorityFeeThreshold := big.NewInt(0)
	for _, fee := range feeHistory.Reward {
		priorityFee = priorityFee.Add(priorityFee, fee[0])
		// We don't need an average, we need the max value
		priorityFeeThreshold = bigmath.Max(priorityFeeThreshold, fee[1])
	}

	u.priorityFeeThresholdMu.Lock()
	u.priorityFeeThreshold = (*assets.Wei)(priorityFeeThreshold)
	u.priorityFeeThresholdMu.Unlock()

	maxPriorityFeePerGas := (*assets.Wei)(priorityFee.Div(priorityFee, big.NewInt(int64(u.config.BlockHistorySize))))
	// baseFeeBufferPercentage is used as a safety to catch fluctuations in the next block.
	maxFeePerGas := baseFee.AddPercentage(BaseFeeBufferPercentage).Add((maxPriorityFeePerGas))

	promUniversalEstimatorBaseFee.WithLabelValues(u.chainID.String()).Set(float64(baseFee.Int64()))
	promUniversalEstimatorMaxPriorityFeePerGas.WithLabelValues(u.chainID.String()).Set(float64(maxPriorityFeePerGas.Int64()))
	promUniversalEstimatorMaxFeePerGas.WithLabelValues(u.chainID.String()).Set(float64(maxFeePerGas.Int64()))

	u.logger.Debugf("Fetched new dynamic prices, block#: %v - maxFeePerGas: %v - maxPriorityFeePerGas: %v - maxPriorityFeeThreshold: %v",
		latestBlock, maxFeePerGas, maxPriorityFeePerGas, (*assets.Wei)(priorityFeeThreshold))

	u.dynamicPriceMu.Lock()
	defer u.dynamicPriceMu.Unlock()
	u.dynamicPrice.FeeCap = maxFeePerGas
	u.dynamicPrice.TipCap = maxPriorityFeePerGas
	return u.dynamicPrice, nil
}

func (u *UniversalEstimator) getDynamicPrice() (fee DynamicFee, err error) {
	u.dynamicPriceMu.RLock()
	defer u.dynamicPriceMu.RUnlock()
	if u.dynamicPrice.FeeCap == nil || u.dynamicPrice.TipCap == nil {
		return fee, fmt.Errorf("dynamic price not set")
	}
	return u.dynamicPrice, nil
}

// BumpLegacyGas provides a bumped gas price value by bumping the previous one by BumpPercent.
// If the original value is higher than the max price it returns an error as there is no room for bumping.
// It aggregates the market, bumped, and max gas price to provide a correct value.
func (u *UniversalEstimator) BumpLegacyGas(ctx context.Context, originalGasPrice *assets.Wei, gasLimit uint64, maxPrice *assets.Wei, _ []EvmPriorAttempt) (*assets.Wei, uint64, error) {
	// Sanitize original fee input
	if originalGasPrice == nil || originalGasPrice.Cmp(maxPrice) >= 0 {
		return nil, 0, fmt.Errorf("%w: error while retrieving original gas price: originalGasPrice: %s. Maximum price configured: %s",
			commonfee.ErrBump, originalGasPrice, maxPrice)
	}

	currentGasPrice, err := u.getGasPrice()
	if err != nil {
		return nil, 0, err
	}

	bumpedGasPrice := originalGasPrice.AddPercentage(u.config.BumpPercent)
	bumpedGasPrice, err = LimitBumpedFee(originalGasPrice, currentGasPrice, bumpedGasPrice, maxPrice)
	if err != nil {
		return nil, 0, fmt.Errorf("gas price error: %s", err.Error())
	}

	u.logger.Debugw("bumped gas price", "originalGasPrice", originalGasPrice, "bumpedGasPrice", bumpedGasPrice)

	return bumpedGasPrice, gasLimit, nil
}

// BumpDynamicFee provides a bumped dynamic fee by bumping the previous one by BumpPercent.
// If the original values are higher than the max price it returns an error as there is no room for bumping. Both maxFeePerGas
// as well as maxPriorityFerPergas need to be bumped otherwise the RPC won't accept the transaction and throw an error.
// See: https://github.com/ethereum/go-ethereum/issues/24284
// It aggregates the market, bumped, and max price to provide a correct value, for both maxFeePerGas as well as maxPriorityFerPergas.
func (u *UniversalEstimator) BumpDynamicFee(ctx context.Context, originalFee DynamicFee, maxPrice *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, err error) {
	// Sanitize original fee input
	// According to geth's spec we need to bump both maxFeePerGas and maxPriorityFeePerGas for the new attempt to be accepted by the RPC
	if originalFee.FeeCap == nil ||
		originalFee.TipCap == nil ||
		((originalFee.TipCap.Cmp(originalFee.FeeCap)) > 0) ||
		(originalFee.FeeCap.Cmp(maxPrice) >= 0) {
		return bumped, fmt.Errorf("%w: error while retrieving original dynamic fees: (originalFeePerGas: %s - originalPriorityFeePerGas: %s). Maximum price configured: %s",
			commonfee.ErrBump, originalFee.FeeCap, originalFee.TipCap, maxPrice)
	}

	currentDynamicPrice, err := u.getDynamicPrice()
	if err != nil {
		return
	}

	bumpedMaxPriorityFeePerGas := originalFee.TipCap.AddPercentage(u.config.BumpPercent)
	bumpedMaxFeePerGas := originalFee.FeeCap.AddPercentage(u.config.BumpPercent)

	if u.config.HasMempool {
		bumpedMaxPriorityFeePerGas, err = LimitBumpedFee(originalFee.TipCap, currentDynamicPrice.TipCap, bumpedMaxPriorityFeePerGas, maxPrice)
		if err != nil {
			return bumped, fmt.Errorf("maxPriorityFeePerGas error: %s", err.Error())
		}

		priorityFeeThreshold, e := u.getPriorityFeeThreshold()
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
	u.logger.Debugw("bumped dynamic fee", "originalFee", originalFee, "bumpedFee", bumpedFee)

	return bumpedFee, nil
}

// LimitBumpedFee selects the maximum value between the original fee and the bumped attempt. If the result is higher than the max price it gets capped.
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
	// and a small BumpPercent.
	if bumpedFee.Cmp(originalFee) == 0 ||
		bumpedFee.Cmp(originalFee.AddPercentage(MinimumBumpPercentage)) < 0 {
		return nil, fmt.Errorf("%w: %s is bumped less than minimum allowed percentage(%s) from originalFee: %s - maxPrice: %s",
			commonfee.ErrBump, bumpedFee, strconv.Itoa(MinimumBumpPercentage), originalFee, maxPrice)
	}
	return bumpedFee, nil
}

func (u *UniversalEstimator) getPriorityFeeThreshold() (*assets.Wei, error) {
	u.priorityFeeThresholdMu.RLock()
	defer u.priorityFeeThresholdMu.RUnlock()
	if u.priorityFeeThreshold == nil {
		return u.priorityFeeThreshold, fmt.Errorf("priorityFeeThreshold not set")
	}
	return u.priorityFeeThreshold, nil
}

func (u *UniversalEstimator) Name() string                                      { return u.logger.Name() }
func (u *UniversalEstimator) L1Oracle() rollups.L1Oracle                        { return u.l1Oracle }
func (u *UniversalEstimator) HealthReport() map[string]error                    { return map[string]error{u.Name(): nil} }
func (u *UniversalEstimator) OnNewLongestChain(context.Context, *evmtypes.Head) {}
