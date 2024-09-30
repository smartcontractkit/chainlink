package gas

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"

	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// MaxStartTime is the maximum amount of time we are allowed to spend
// trying to fill initial data on start. This must be capped because it can
// block the application from starting.
var MaxStartTime = 10 * time.Second

var (
	promBlockHistoryEstimatorAllGasPricePercentiles = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_all_gas_price_percentiles",
		Help: "Gas price at given percentile",
	},
		[]string{"percentile", "evmChainID"},
	)
	promBlockHistoryEstimatorAllTipCapPercentiles = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_all_tip_cap_percentiles",
		Help: "Tip cap at given percentile",
	},
		[]string{"percentile", "evmChainID"},
	)
	promBlockHistoryEstimatorSetGasPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_set_gas_price",
		Help: "Gas updater set gas price (in Wei)",
	},
		[]string{"percentile", "evmChainID"},
	)
	promBlockHistoryEstimatorSetTipCap = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_set_tip_cap",
		Help: "Gas updater set gas tip cap (in Wei)",
	},
		[]string{"percentile", "evmChainID"},
	)
	promBlockHistoryEstimatorCurrentBaseFee = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_current_base_fee",
		Help: "Gas updater current block base fee in Wei",
	},
		[]string{"evmChainID"},
	)
	promBlockHistoryEstimatorConnectivityFailureCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "block_history_estimator_connectivity_failure_count",
		Help: "Counter is incremented every time a gas bump is prevented due to a detected network propagation/connectivity issue",
	},
		[]string{"evmChainID", "mode"},
	)
)

const BumpingHaltedLabel = "Tx gas bumping halted since price exceeds current block prices by significant margin; tx will continue to be rebroadcasted but your node, RPC, or the chain might be experiencing connectivity issues; please investigate and fix ASAP"

var _ EvmEstimator = &BlockHistoryEstimator{}

type chainConfig interface {
	ChainType() chaintype.ChainType
}

type estimatorGasEstimatorConfig interface {
	EIP1559DynamicFees() bool
	BumpThreshold() uint64
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMax() *assets.Wei
	PriceMin() *assets.Wei
	bumpConfig
}

type BlockHistoryEstimator struct {
	services.StateMachine
	ethClient feeEstimatorClient
	chainID   *big.Int
	config    chainConfig
	eConfig   estimatorGasEstimatorConfig
	bhConfig  BlockHistoryConfig
	// NOTE: it is assumed that blocks will be kept sorted by
	// block number ascending
	blocks   []evmtypes.Block
	blocksMu sync.RWMutex
	size     int64
	mb       *mailbox.Mailbox[*evmtypes.Head]
	wg       *sync.WaitGroup
	stopCh   services.StopChan

	gasPrice              *assets.Wei
	tipCap                *assets.Wei
	priceMu               sync.RWMutex
	maxPercentileGasPrice *assets.Wei
	maxPercentileTipCap   *assets.Wei
	maxPriceMu            sync.RWMutex
	latest                *evmtypes.Head
	latestMu              sync.RWMutex
	initialFetch          atomic.Bool

	logger logger.SugaredLogger

	l1Oracle rollups.L1Oracle
}

// NewBlockHistoryEstimator returns a new BlockHistoryEstimator that listens
// for new heads and updates the base gas price dynamically based on the
// configured percentile of gas prices in that block
func NewBlockHistoryEstimator(lggr logger.Logger, ethClient feeEstimatorClient, cfg chainConfig, eCfg estimatorGasEstimatorConfig, bhCfg BlockHistoryConfig, chainID *big.Int, l1Oracle rollups.L1Oracle) EvmEstimator {
	return &BlockHistoryEstimator{
		ethClient: ethClient,
		chainID:   chainID,
		config:    cfg,
		eConfig:   eCfg,
		bhConfig:  bhCfg,
		blocks:    make([]evmtypes.Block, 0),
		// Must have enough blocks for both estimator and connectivity checker
		size:     int64(mathutil.Max(bhCfg.BlockHistorySize(), bhCfg.CheckInclusionBlocks())),
		mb:       mailbox.NewSingle[*evmtypes.Head](),
		wg:       new(sync.WaitGroup),
		stopCh:   make(chan struct{}),
		logger:   logger.Sugared(logger.Named(lggr, "BlockHistoryEstimator")),
		l1Oracle: l1Oracle,
	}
}

// OnNewLongestChain recalculates and sets global gas price if a sampled new head comes
// in and we are not currently fetching
func (b *BlockHistoryEstimator) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	// set latest base fee here to avoid potential lag introduced by block delay
	// it is really important that base fee be as up-to-date as possible
	b.setLatest(head)
	b.mb.Deliver(head)
}

// setLatest assumes that head won't be mutated
func (b *BlockHistoryEstimator) setLatest(head *evmtypes.Head) {
	// Non-eip1559 blocks don't include base fee
	if baseFee := head.BaseFeePerGas; baseFee != nil {
		promBlockHistoryEstimatorCurrentBaseFee.WithLabelValues(b.chainID.String()).Set(float64(baseFee.Int64()))
	}
	b.logger.Debugw("Set latest block", "blockNum", head.Number, "blockHash", head.Hash, "baseFee", head.BaseFeePerGas, "baseFeeWei", head.BaseFeePerGas.ToInt())
	b.latestMu.Lock()
	defer b.latestMu.Unlock()
	b.latest = head
}

func (b *BlockHistoryEstimator) getCurrentBaseFee() *assets.Wei {
	b.latestMu.RLock()
	defer b.latestMu.RUnlock()
	if b.latest == nil {
		return nil
	}
	return b.latest.BaseFeePerGas
}

func (b *BlockHistoryEstimator) getBlocks() []evmtypes.Block {
	b.blocksMu.RLock()
	defer b.blocksMu.RUnlock()
	return b.blocks
}

// Start starts BlockHistoryEstimator service.
// The provided context can be used to terminate Start sequence.
func (b *BlockHistoryEstimator) Start(ctx context.Context) error {
	return b.StartOnce("BlockHistoryEstimator", func() error {
		b.logger.Trace("Starting")

		if b.bhConfig.CheckInclusionBlocks() > 0 {
			b.logger.Infof("Inclusion checking enabled, bumping will be prevented on transactions that have been priced above the %d percentile of transactions in the latest %d blocks", b.bhConfig.CheckInclusionPercentile(), b.bhConfig.CheckInclusionBlocks())
		}
		if b.bhConfig.BlockHistorySize() == 0 {
			return errors.New("BlockHistorySize must be set to a value greater than 0")
		}

		fetchCtx, cancel := context.WithTimeout(ctx, MaxStartTime)
		defer cancel()
		latestHead, err := b.ethClient.HeadByNumber(fetchCtx, nil)
		if err != nil {
			b.logger.Warnw("Initial check for latest head failed", "err", err)
		} else if latestHead == nil {
			b.logger.Warnw("initial check for latest head failed, head was unexpectedly nil")
		} else {
			b.logger.Debugw("Got latest head", "number", latestHead.Number, "blockHash", latestHead.Hash.Hex())
			b.setLatest(latestHead)
			b.FetchBlocksAndRecalculate(fetchCtx, latestHead)
		}

		// NOTE: This only checks the start context, not the fetch context
		if ctx.Err() != nil {
			return fmt.Errorf("failed to start BlockHistoryEstimator due to main context error: %w", ctx.Err())
		}

		b.wg.Add(1)
		go b.runLoop()

		b.logger.Trace("Started")
		return nil
	})
}

func (b *BlockHistoryEstimator) L1Oracle() rollups.L1Oracle {
	return b.l1Oracle
}

func (b *BlockHistoryEstimator) Close() error {
	return b.StopOnce("BlockHistoryEstimator", func() error {
		close(b.stopCh)
		b.wg.Wait()
		return nil
	})
}

func (b *BlockHistoryEstimator) Name() string {
	return b.logger.Name()
}
func (b *BlockHistoryEstimator) HealthReport() map[string]error {
	return map[string]error{b.Name(): b.Healthy()}
}

func (b *BlockHistoryEstimator) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint64, maxGasPriceWei *assets.Wei, _ ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	ok := b.IfStarted(func() {
		gasPrice = b.getGasPrice()
	})
	if !ok {
		return nil, 0, errors.New("BlockHistoryEstimator is not started; cannot estimate gas")
	}
	if gasPrice == nil {
		if !b.initialFetch.Load() {
			return nil, 0, errors.New("BlockHistoryEstimator has not finished the first gas estimation yet, likely because a failure on start")
		}
		b.logger.Warnw("Failed to estimate gas price. This is likely because there aren't any valid transactions to estimate from."+
			"Using Evm.GasEstimator.PriceDefault as fallback.", "blocks", b.getBlockHistoryNumbers())
		gasPrice = b.eConfig.PriceDefault()
	}
	gasPrice = capGasPrice(gasPrice, maxGasPriceWei, b.eConfig.PriceMax())
	chainSpecificGasLimit = gasLimit
	return
}

func (b *BlockHistoryEstimator) getGasPrice() *assets.Wei {
	b.priceMu.RLock()
	defer b.priceMu.RUnlock()
	return b.gasPrice
}

func (b *BlockHistoryEstimator) getMaxPercentileGasPrice() *assets.Wei {
	b.maxPriceMu.RLock()
	defer b.maxPriceMu.RUnlock()
	return b.maxPercentileGasPrice
}

func (b *BlockHistoryEstimator) setMaxPercentileGasPrice(gasPrice *assets.Wei) {
	b.maxPriceMu.Lock()
	defer b.maxPriceMu.Unlock()
	b.maxPercentileGasPrice = gasPrice
}

func (b *BlockHistoryEstimator) getBlockHistoryNumbers() (numsInHistory []int64) {
	b.blocksMu.RLock()
	defer b.blocksMu.RUnlock()
	for _, b := range b.blocks {
		numsInHistory = append(numsInHistory, b.Number)
	}
	return
}

func (b *BlockHistoryEstimator) getTipCap() *assets.Wei {
	b.priceMu.RLock()
	defer b.priceMu.RUnlock()
	return b.tipCap
}

func (b *BlockHistoryEstimator) getMaxPercentileTipCap() *assets.Wei {
	b.maxPriceMu.RLock()
	defer b.maxPriceMu.RUnlock()
	return b.maxPercentileTipCap
}

func (b *BlockHistoryEstimator) setMaxPercentileTipCap(tipCap *assets.Wei) {
	b.maxPriceMu.Lock()
	defer b.maxPriceMu.Unlock()
	b.maxPercentileTipCap = tipCap
}

func (b *BlockHistoryEstimator) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, gasLimit uint64, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumpedGasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	if b.bhConfig.CheckInclusionBlocks() > 0 {
		if err = b.haltBumping(attempts); err != nil {
			if errors.Is(err, commonfee.ErrConnectivity) {
				b.logger.Criticalw(BumpingHaltedLabel, "err", err)
				b.SvcErrBuffer.Append(err)
				promBlockHistoryEstimatorConnectivityFailureCount.WithLabelValues(b.chainID.String(), "legacy").Inc()
			}
			return nil, 0, err
		}
	}
	bumpedGasPrice, err = BumpLegacyGasPriceOnly(b.eConfig, b.logger, b.getGasPrice(), originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, 0, err
	}
	return bumpedGasPrice, gasLimit, err
}

// haltBumping prevents transactions from excessively bumping if an existing attempt's price is above a configurable percentile
// This check is required in case the transaction is not being included due to some kind of mempool propagation
// or connectivity issue rather than insufficiently high pricing
func (b *BlockHistoryEstimator) haltBumping(attempts []EvmPriorAttempt) error {
	percentile := int(b.bhConfig.CheckInclusionPercentile())
	// Get latest CheckInclusionPercentile prices to use in the checks below
	var maxGasPrice *assets.Wei
	var maxTipCap *assets.Wei
	ok := b.IfStarted(func() {
		maxGasPrice = b.getMaxPercentileGasPrice()
		maxTipCap = b.getMaxPercentileTipCap()
	})
	if !ok {
		return errors.New("BlockHistoryEstimator is not started; do not have max gas to allow bumping")
	}
	if !b.initialFetch.Load() {
		return errors.New("BlockHistoryEstimator has not finished the first gas estimation yet, likely because a failure on start")
	}
	// Return error to prevent bumping if gas price is nil or if EIP1559 is enabled and tip cap is nil
	if maxGasPrice == nil || (b.eConfig.EIP1559DynamicFees() && maxTipCap == nil) {
		errorMsg := fmt.Sprintf("%d percentile price is not set. This is likely because there aren't any valid transactions to estimate from. Preventing bumping until valid price is available to compare", percentile)
		b.logger.Debugf(errorMsg)
		return errors.New(errorMsg)
	}
	// Get the latest CheckInclusionBlocks from block history for fee cap check below
	blockHistory := b.getBlocks()
	blockRange := mathutil.Min(len(blockHistory), int(b.bhConfig.CheckInclusionBlocks()))
	startIdx := len(blockHistory) - blockRange
	checkInclusionBlocks := blockHistory[startIdx:]
	// Check each attempt for any with a gas price or tip cap (if EIP1559 type) exceeds the latest CheckInclusionPercentile prices
	for _, attempt := range attempts {
		if attempt.BroadcastBeforeBlockNum == nil {
			// this shouldn't happen; any broadcast attempt ought to have a
			// BroadcastBeforeBlockNum otherwise its an assumption violation
			return fmt.Errorf("BroadcastBeforeBlockNum was unexpectedly nil for attempt %s", attempt.TxHash)
		}
		var attemptEip1559 bool
		switch attempt.TxType {
		case 0x0, 0x1:
			attemptEip1559 = false
		case 0x2:
			attemptEip1559 = true
		default:
			return fmt.Errorf("attempt %s has unknown transaction type 0x%d", attempt.TxHash, attempt.TxType)
		}
		if !attemptEip1559 {
			if attempt.GasPrice.Cmp(maxGasPrice) > 0 {
				return fmt.Errorf("transaction %s has gas price of %s, which is above percentile=%d%% (percentile price: %s): %w", attempt.TxHash, attempt.GasPrice, percentile, maxGasPrice, commonfee.ErrConnectivity)
			}
			continue
		}
		sufficientFeeCap := true
		for _, b := range checkInclusionBlocks {
			// feecap must >= tipcap+basefee for the block, otherwise there
			// is no way this could have been included, and we must bail
			// out of the check
			attemptFeeCap := attempt.DynamicFee.FeeCap
			attemptTipCap := attempt.DynamicFee.TipCap
			if attemptFeeCap.Cmp(attemptTipCap.Add(b.BaseFeePerGas)) < 0 {
				sufficientFeeCap = false
				break
			}
		}
		if sufficientFeeCap && attempt.DynamicFee.TipCap.Cmp(maxTipCap) > 0 {
			return fmt.Errorf("transaction %s has tip cap of %s, which is above percentile=%d%% (percentile tip cap: %s): %w", attempt.TxHash, attempt.DynamicFee.TipCap, percentile, maxTipCap, commonfee.ErrConnectivity)
		}
	}
	return nil
}

func (b *BlockHistoryEstimator) GetDynamicFee(_ context.Context, maxGasPriceWei *assets.Wei) (fee DynamicFee, err error) {
	if !b.eConfig.EIP1559DynamicFees() {
		return fee, errors.New("can't get dynamic fee, EIP1559 is disabled")
	}

	var feeCap *assets.Wei
	var tipCap *assets.Wei
	ok := b.IfStarted(func() {
		b.priceMu.RLock()
		defer b.priceMu.RUnlock()
		tipCap = b.tipCap
		if tipCap == nil {
			if !b.initialFetch.Load() {
				err = errors.New("BlockHistoryEstimator has not finished the first gas estimation yet, likely because a failure on start")
				return
			}
			b.logger.Warnw("Failed to estimate gas price. This is likely because there aren't any valid transactions to estimate from."+
				"Using Evm.GasEstimator.TipCapDefault as fallback.", "blocks", b.getBlockHistoryNumbers())
			tipCap = b.eConfig.TipCapDefault()
		}
		maxGasPrice := getMaxGasPrice(maxGasPriceWei, b.eConfig.PriceMax())
		if b.eConfig.BumpThreshold() == 0 {
			// just use the max gas price if gas bumping is disabled
			feeCap = maxGasPrice
		} else if b.getCurrentBaseFee() != nil {
			// HACK: due to a flaw of how EIP-1559 is implemented we have to
			// set a much lower FeeCap than the actual maximum we are willing
			// to pay in order to give ourselves headroom for bumping
			// See: https://github.com/ethereum/go-ethereum/issues/24284
			feeCap = calcFeeCap(b.getCurrentBaseFee(), int(b.bhConfig.EIP1559FeeCapBufferBlocks()), tipCap, maxGasPrice)
		} else {
			// This shouldn't happen on EIP-1559 blocks, since if the tip cap
			// is set, Start must have succeeded and we would expect an initial
			// base fee to be set as well
			err = errors.New("BlockHistoryEstimator: no value for latest block base fee; cannot estimate EIP-1559 base fee. Are you trying to run with EIP1559 enabled on a non-EIP1559 chain?")
			return
		}
	})
	if !ok {
		return fee, errors.New("BlockHistoryEstimator is not started; cannot estimate gas")
	}
	if err != nil {
		return fee, err
	}
	fee.FeeCap = feeCap
	fee.TipCap = tipCap
	return
}

func calcFeeCap(latestAvailableBaseFeePerGas *assets.Wei, bufferBlocks int, tipCap *assets.Wei, maxGasPriceWei *assets.Wei) (feeCap *assets.Wei) {
	const maxBaseFeeIncreasePerBlock float64 = 1.125

	baseFee := new(big.Float)
	baseFee.SetInt(latestAvailableBaseFeePerGas.ToInt())
	// Find out the worst case base fee before we should bump
	multiplier := big.NewFloat(maxBaseFeeIncreasePerBlock)
	for i := 0; i < bufferBlocks; i++ {
		baseFee.Mul(baseFee, multiplier)
	}

	baseFeeInt, _ := baseFee.Int(nil)
	feeCap = assets.NewWei(baseFeeInt.Add(baseFeeInt, tipCap.ToInt()))

	if feeCap.Cmp(maxGasPriceWei) > 0 {
		return maxGasPriceWei
	}
	return feeCap
}

func (b *BlockHistoryEstimator) BumpDynamicFee(_ context.Context, originalFee DynamicFee, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumped DynamicFee, err error) {
	if b.bhConfig.CheckInclusionBlocks() > 0 {
		if err = b.haltBumping(attempts); err != nil {
			if errors.Is(err, commonfee.ErrConnectivity) {
				b.logger.Criticalw(BumpingHaltedLabel, "err", err)
				b.SvcErrBuffer.Append(err)
				promBlockHistoryEstimatorConnectivityFailureCount.WithLabelValues(b.chainID.String(), "eip1559").Inc()
			}
			return bumped, err
		}
	}
	return BumpDynamicFeeOnly(b.eConfig, b.bhConfig.EIP1559FeeCapBufferBlocks(), b.logger, b.getTipCap(), b.getCurrentBaseFee(), originalFee, maxGasPriceWei)
}

func (b *BlockHistoryEstimator) runLoop() {
	defer b.wg.Done()
	ctx, cancel := b.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-b.mb.Notify():
			head, exists := b.mb.Retrieve()
			if !exists {
				b.logger.Debug("No head to retrieve")
				continue
			}
			b.FetchBlocksAndRecalculate(ctx, head)
		}
	}
}

// FetchBlocksAndRecalculate fetches block history leading up to head and recalculates gas price.
func (b *BlockHistoryEstimator) FetchBlocksAndRecalculate(ctx context.Context, head *evmtypes.Head) {
	if err := b.FetchBlocks(ctx, head); err != nil {
		b.logger.Warnw("Error fetching blocks", "head", head, "err", err)
		return
	}
	b.initialFetch.Store(true)
	b.Recalculate(head)
}

// Recalculate adds the given heads to the history and recalculates gas price.
func (b *BlockHistoryEstimator) Recalculate(head *evmtypes.Head) {
	lggr := b.logger.With("head", head)

	blockHistory := b.getBlocks()
	if len(blockHistory) == 0 {
		lggr.Debug("No blocks in history, cannot set gas price")
		return
	}

	// Calculate and set the TransactionPercentile gas price and tip cap to use during gas estimation
	b.calculateGasPriceTipCap(lggr, blockHistory, head)

	// Calculate and set the CheckInclusionPercentile gas price and tip cap to halt excessive bumping
	b.calculateMaxPercentileGasPriceTipCap(lggr, blockHistory, head)
}

func (b *BlockHistoryEstimator) calculateGasPriceTipCap(lggr logger.SugaredLogger, blockHistory []evmtypes.Block, head *evmtypes.Head) {
	percentile := int(b.bhConfig.TransactionPercentile())
	eip1559 := b.eConfig.EIP1559DynamicFees()

	blockRange := mathutil.Min(len(blockHistory), int(b.bhConfig.BlockHistorySize()))
	startIdx := len(blockHistory) - blockRange
	blocks := blockHistory[startIdx:]

	percentileGasPrice, percentileTipCap, err := b.calculatePercentilePrices(blocks, percentile, eip1559,
		func(gasPrices []*assets.Wei) {
			for i := 0; i <= 100; i += 5 {
				jdx := ((len(gasPrices) - 1) * i) / 100
				promBlockHistoryEstimatorAllGasPricePercentiles.WithLabelValues(fmt.Sprintf("%v%%", i), b.chainID.String()).Set(float64(gasPrices[jdx].Int64()))
			}
		}, func(tipCaps []*assets.Wei) {
			for i := 0; i <= 100; i += 5 {
				jdx := ((len(tipCaps) - 1) * i) / 100
				promBlockHistoryEstimatorAllTipCapPercentiles.WithLabelValues(fmt.Sprintf("%v%%", i), b.chainID.String()).Set(float64(tipCaps[jdx].Int64()))
			}
		})
	if err != nil {
		if errors.Is(err, ErrNoSuitableTransactions) {
			lggr.Debug("No suitable transactions, skipping")
		} else {
			lggr.Warnw("Cannot calculate percentile prices", "err", err)
		}
		return
	}

	var numsForPrice []int64
	for _, b := range blocks {
		numsForPrice = append(numsForPrice, b.Number)
	}

	gasPriceGwei := percentileGasPrice.Text("gwei")
	lggrFields := []interface{}{
		"gasPriceWei", percentileGasPrice,
		"gasPriceGWei", gasPriceGwei,
		"maxGasPriceWei", b.eConfig.PriceMax(),
		"headNum", head.Number,
		"priceBlocks", numsForPrice,
	}
	b.setPercentileGasPrice(percentileGasPrice)
	promBlockHistoryEstimatorSetGasPrice.WithLabelValues(fmt.Sprintf("%v%%", percentile), b.chainID.String()).Set(float64(percentileGasPrice.Int64()))

	if !eip1559 {
		lggr.Debugw(fmt.Sprintf("Setting new default GasPrice: %v Gwei", gasPriceGwei), lggrFields...)
		return
	}

	tipCapGwei := percentileTipCap.Text("gwei")
	lggrFields = append(lggrFields, []interface{}{
		"tipCapWei", percentileTipCap,
		"tipCapGwei", tipCapGwei,
	}...)
	lggr.Debugw(fmt.Sprintf("Setting new default prices, GasPrice: %v Gwei, TipCap: %v Gwei", gasPriceGwei, tipCapGwei), lggrFields...)
	b.setPercentileTipCap(percentileTipCap)
	promBlockHistoryEstimatorSetTipCap.WithLabelValues(fmt.Sprintf("%v%%", percentile), b.chainID.String()).Set(float64(percentileTipCap.Int64()))
}

func (b *BlockHistoryEstimator) calculateMaxPercentileGasPriceTipCap(lggr logger.SugaredLogger, blockHistory []evmtypes.Block, head *evmtypes.Head) {
	if b.bhConfig.CheckInclusionBlocks() <= 0 {
		return
	}
	percentile := int(b.bhConfig.CheckInclusionPercentile())
	eip1559 := b.eConfig.EIP1559DynamicFees()

	blockRange := mathutil.Min(len(blockHistory), int(b.bhConfig.CheckInclusionBlocks()))
	startIdx := len(blockHistory) - blockRange
	checkInclusionBlocks := blockHistory[startIdx:]

	maxPercentileGasPrice, maxPercentileTipCap, err := b.calculatePercentilePrices(checkInclusionBlocks, percentile, eip1559, nil, nil)
	if err != nil {
		if errors.Is(err, ErrNoSuitableTransactions) {
			lggr.Debug("No suitable transactions found to calculate the max percentile prices, skipping")
		} else {
			lggr.Warnw("Cannot calculate max percentile prices", "err", err)
		}
		return
	}

	var numsForMaxPrices []int64
	for _, b := range checkInclusionBlocks {
		numsForMaxPrices = append(numsForMaxPrices, b.Number)
	}

	maxPercentileGasPriceGwei := maxPercentileGasPrice.Text("gwei")
	lggrFields := []interface{}{
		"maxPercentileGasPriceWei", maxPercentileGasPrice,
		"maxPercentileGasPriceGwei", maxPercentileGasPriceGwei,
		"headNum", head.Number,
		"maxPercentileBlocks", numsForMaxPrices,
	}

	b.setMaxPercentileGasPrice(maxPercentileGasPrice)

	if !eip1559 {
		lggr.Debugw(fmt.Sprintf("Setting new max percentile GasPrice: %v Gwei", maxPercentileGasPriceGwei), lggrFields...)
		return
	}

	maxPercentileTipCapGwei := maxPercentileTipCap.Text("gwei")
	lggrFields = append(lggrFields, []interface{}{
		"maxPercentileTipCapWei", maxPercentileTipCap,
		"maxPercentileTipCapGwei", maxPercentileTipCapGwei,
	}...)
	lggr.Debugw(fmt.Sprintf("Setting new default prices, max percentile GasPrice: %v Gwei, max percentile TipCap: %v Gwei", maxPercentileGasPriceGwei, maxPercentileTipCapGwei), lggrFields...)
	b.setMaxPercentileTipCap(maxPercentileTipCap)
}

// FetchBlocks fetches block history leading up to the given head.
func (b *BlockHistoryEstimator) FetchBlocks(ctx context.Context, head *evmtypes.Head) error {
	// HACK: blockDelay is the number of blocks that the block history estimator trails behind head.
	// E.g. if this is set to 3, and we receive block 10, block history estimator will
	// fetch block 7.
	// This is necessary because geth/parity send heads as soon as they get
	// them and often the actual block is not available until later. Fetching
	// it too early results in an empty block.
	blockDelay := int64(b.bhConfig.BlockDelay())
	historySize := b.size

	if historySize <= 0 {
		return fmt.Errorf("BlockHistoryEstimator: history size must be > 0, got: %d", historySize)
	}

	highestBlockToFetch := head.Number - blockDelay
	if highestBlockToFetch < 0 {
		return fmt.Errorf("BlockHistoryEstimator: cannot fetch, current block height %v is lower than EVM.RPCBlockQueryDelay=%v", head.Number, blockDelay)
	}
	lowestBlockToFetch := head.Number - historySize - blockDelay + 1
	if lowestBlockToFetch < 0 {
		lowestBlockToFetch = 0
	}

	blocks := make(map[int64]evmtypes.Block)
	earliestInChain := head.EarliestInChain()
	for _, block := range b.getBlocks() {
		// Make a best-effort to be re-org resistant using the head
		// chain, refetch blocks that got re-org'd out.
		// NOTE: Any blocks in the history that are older than the oldest block
		// in the provided chain will be assumed final.
		if block.Number < earliestInChain.BlockNumber() {
			blocks[block.Number] = block
		} else if head.IsInChain(block.Hash) {
			blocks[block.Number] = block
		}
	}

	var reqs []rpc.BatchElem
	// Fetch blocks in reverse order so if it times out halfway through we bias
	// towards more recent blocks
	for i := highestBlockToFetch; i >= lowestBlockToFetch; i-- {
		// NOTE: To save rpc calls, don't fetch blocks we already have in the history
		if _, exists := blocks[i]; exists {
			continue
		}

		req := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{Int64ToHex(i), true},
			Result: &evmtypes.Block{},
		}
		reqs = append(reqs, req)
	}

	lggr := b.logger.With("head", head)

	lggr.Tracew(fmt.Sprintf("Fetching %v blocks (%v in local history)", len(reqs), len(blocks)), "n", len(reqs), "inHistory", len(blocks), "blockNum", head.Number)
	if err := b.batchFetch(ctx, reqs); err != nil {
		return err
	}

	var missingBlocks []int64
	for _, req := range reqs {
		result, err := req.Result, req.Error
		if err != nil {
			if errors.Is(err, evmtypes.ErrMissingBlock) {
				num := HexToInt64(req.Args[0])
				missingBlocks = append(missingBlocks, num)
				lggr.Debugw(
					fmt.Sprintf("Failed to fetch block: RPC node returned a missing block on query for block number %d even though the WS subscription already sent us this block. It might help to increase EVM.RPCBlockQueryDelay (currently %d)",
						num, blockDelay,
					),
					"err", err, "blockNum", num, "headNum", head.Number)
			} else {
				lggr.Warnw("Failed to fetch block", "err", err, "blockNum", HexToInt64(req.Args[0]), "headNum", head.Number)
			}
			continue
		}

		block, is := result.(*evmtypes.Block)
		if !is {
			return fmt.Errorf("expected result to be a %T, got %T", &evmtypes.Block{}, result)
		}
		if block == nil {
			return errors.New("invariant violation: got nil block")
		}
		if block.Hash == (common.Hash{}) {
			lggr.Warnw("Block was missing hash", "block", b, "headNum", head.Number, "blockNum", block.Number)
			continue
		}

		blocks[block.Number] = *block
	}

	if len(missingBlocks) > 1 {
		lggr.Errorw(
			fmt.Sprintf("RPC node returned multiple missing blocks on query for block numbers %v even though the WS subscription already sent us these blocks. It might help to increase EVM.RPCBlockQueryDelay (currently %d)",
				missingBlocks, blockDelay,
			),
			"blockNums", missingBlocks, "headNum", head.Number)
	}

	newBlockHistory := make([]evmtypes.Block, 0)

	for _, block := range blocks {
		newBlockHistory = append(newBlockHistory, block)
	}
	sort.Slice(newBlockHistory, func(i, j int) bool {
		return newBlockHistory[i].Number < newBlockHistory[j].Number
	})

	start := len(newBlockHistory) - int(historySize)
	if start < 0 {
		lggr.Debugw(fmt.Sprintf("Using fewer blocks than the specified history size: %v/%v", len(newBlockHistory), historySize), "blocksSize", historySize, "headNum", head.Number, "blocksAvailable", len(newBlockHistory))
		start = 0
	}

	b.blocksMu.Lock()
	b.blocks = newBlockHistory[start:]
	b.blocksMu.Unlock()

	return nil
}

func (b *BlockHistoryEstimator) batchFetch(ctx context.Context, reqs []rpc.BatchElem) error {
	batchSize := int(b.bhConfig.BatchSize())

	if batchSize == 0 {
		batchSize = len(reqs)
	}

	for i := 0; i < len(reqs); i += batchSize {
		j := i + batchSize
		if j > len(reqs) {
			j = len(reqs)
		}

		b.logger.Tracew(fmt.Sprintf("Batch fetching blocks %v thru %v", HexToInt64(reqs[i].Args[0]), HexToInt64(reqs[j-1].Args[0])))

		err := b.ethClient.BatchCallContext(ctx, reqs[i:j])
		if errors.Is(err, context.DeadlineExceeded) {
			// We ran out of time, return what we have
			b.logger.Warnf("Batch fetching timed out; loaded %d/%d results: %v", i, len(reqs), err)
			for k := i; k < len(reqs); k++ {
				if k < j {
					reqs[k].Error = fmt.Errorf("request failed: %w", err)
				} else {
					reqs[k].Error = fmt.Errorf("request skipped; previous request exceeded deadline: %w", err)
				}
			}
			return nil
		} else if err != nil {
			return fmt.Errorf("BlockHistoryEstimator#fetchBlocks error fetching blocks with BatchCallContext: %w", err)
		}
	}
	return nil
}

var (
	ErrNoSuitableTransactions = errors.New("no suitable transactions")
)

func (b *BlockHistoryEstimator) calculatePercentilePrices(blocks []evmtypes.Block, percentile int, eip1559 bool, f func(gasPrices []*assets.Wei), f2 func(tipCaps []*assets.Wei)) (gasPrice, tipCap *assets.Wei, err error) {
	gasPrices, tipCaps := b.getPricesFromBlocks(blocks, eip1559)
	if len(gasPrices) == 0 {
		return nil, nil, ErrNoSuitableTransactions
	}
	sort.Slice(gasPrices, func(i, j int) bool { return gasPrices[i].Cmp(gasPrices[j]) < 0 })
	if f != nil {
		f(gasPrices)
	}
	gasPrice = gasPrices[((len(gasPrices)-1)*percentile)/100]

	if !eip1559 {
		return
	}
	if len(tipCaps) == 0 {
		return nil, nil, ErrNoSuitableTransactions
	}
	sort.Slice(tipCaps, func(i, j int) bool { return tipCaps[i].Cmp(tipCaps[j]) < 0 })
	if f2 != nil {
		f2(tipCaps)
	}
	tipCap = tipCaps[((len(tipCaps)-1)*percentile)/100]

	return
}

func (b *BlockHistoryEstimator) getPricesFromBlocks(blocks []evmtypes.Block, eip1559 bool) (gasPrices, tipCaps []*assets.Wei) {
	gasPrices = make([]*assets.Wei, 0)
	tipCaps = make([]*assets.Wei, 0)
	for _, block := range blocks {
		if err := verifyBlock(block, eip1559); err != nil {
			b.logger.Warnw(fmt.Sprintf("Block %v is not usable, %s", block.Number, err.Error()), "block", block, "err", err)
		}
		for _, tx := range block.Transactions {
			if b.IsUsable(tx, block, b.config.ChainType(), b.eConfig.PriceMin(), b.logger) {
				gp := b.EffectiveGasPrice(block, tx)
				if gp == nil {
					b.logger.Warnw("Unable to get gas price for tx", "tx", tx, "block", block)
					continue
				}
				gasPrices = append(gasPrices, gp)
				if !eip1559 {
					continue
				}
				tc := b.EffectiveTipCap(block, tx)
				if tc == nil {
					b.logger.Warnw("Unable to get tip cap for tx", "tx", tx, "block", block)
					continue
				}
				tipCaps = append(tipCaps, tc)
			}
		}
	}
	return
}

func verifyBlock(block evmtypes.Block, eip1559 bool) error {
	if eip1559 && block.BaseFeePerGas == nil {
		return errors.New("EIP-1559 mode was enabled, but block was missing baseFeePerGas")
	}
	return nil
}

func (b *BlockHistoryEstimator) setPercentileTipCap(tipCap *assets.Wei) {
	max := b.eConfig.PriceMax()
	min := b.eConfig.TipCapMin()

	b.priceMu.Lock()
	defer b.priceMu.Unlock()
	if tipCap.Cmp(max) > 0 {
		b.logger.Warnw(fmt.Sprintf("Calculated gas tip cap of %s exceeds EVM.GasEstimator.PriceMax=%[2]s, setting gas tip cap to the maximum allowed value of %[2]s instead", tipCap.String(), max.String()), "tipCapWei", tipCap, "minTipCapWei", min, "maxTipCapWei", max)
		b.tipCap = max
	} else if tipCap.Cmp(min) < 0 {
		b.logger.Warnw(fmt.Sprintf("Calculated gas tip cap of %s falls below EVM.GasEstimator.TipCapMin=%[2]s, setting gas tip cap to the minimum allowed value of %[2]s instead", tipCap.String(), min.String()), "tipCapWei", tipCap, "minTipCapWei", min, "maxTipCapWei", max)
		b.tipCap = min
	} else {
		b.tipCap = tipCap
	}
}

func (b *BlockHistoryEstimator) setPercentileGasPrice(gasPrice *assets.Wei) {
	max := b.eConfig.PriceMax()
	min := b.eConfig.PriceMin()
	eip1559 := b.eConfig.EIP1559DynamicFees()
	var warn string

	b.priceMu.Lock()
	defer b.priceMu.Unlock()
	if gasPrice.Cmp(max) > 0 {
		warn = fmt.Sprintf("Calculated gas price of %s exceeds EVM.GasEstimator.PriceMax=%[2]s, setting gas price to the maximum allowed value of %[2]s instead", gasPrice.String(), max.String())
		b.gasPrice = max
	} else if gasPrice.Cmp(min) < 0 {
		warn = fmt.Sprintf("Calculated gas price of %s falls below EVM.GasEstimator.PriceMin=%[2]s, setting gas price to the minimum allowed value of %[2]s instead", gasPrice.String(), min.String())
		b.gasPrice = min
	} else {
		b.gasPrice = gasPrice
	}

	if !eip1559 && len(warn) > 0 {
		b.logger.Warnw(warn, "gasPriceWei", gasPrice, "maxGasPriceWei", max, "minGasPriceWei", min)
	}
}

// isUsable returns true if the tx is usable both generally and specifically for
// this Config.
func (b *BlockHistoryEstimator) IsUsable(tx evmtypes.Transaction, block evmtypes.Block, chainType chaintype.ChainType, minGasPrice *assets.Wei, lggr logger.Logger) bool {
	// GasLimit 0 is impossible on Ethereum official, but IS possible
	// on forks/clones such as RSK. We should ignore these transactions
	// if they come up on any chain since they are not normal.
	if tx.GasLimit == 0 {
		return false
	}
	// NOTE: This really shouldn't be possible, but at least one node op has
	// reported it happening on mainnet so we need to handle this case
	if tx.GasPrice == nil && tx.Type == 0x0 {
		lggr.Debugw("Ignoring transaction that was unexpectedly missing gas price", "tx", tx)
		return false
	}
	return chainSpecificIsUsable(tx, block.BaseFeePerGas, chainType, minGasPrice)
}

func (b *BlockHistoryEstimator) EffectiveGasPrice(block evmtypes.Block, tx evmtypes.Transaction) *assets.Wei {
	switch tx.Type {
	case 0x0, 0x1:
		return tx.GasPrice
	case 0x2, 0x3:
		return b.getEffectiveGasPrice(block, tx)
	default:
		b.logger.Debugw(fmt.Sprintf("Ignoring unknown transaction type %v", tx.Type), "block", block, "tx", tx)
		return nil
	}
}

func (b *BlockHistoryEstimator) getEffectiveGasPrice(block evmtypes.Block, tx evmtypes.Transaction) *assets.Wei {
	if block.BaseFeePerGas == nil || tx.MaxPriorityFeePerGas == nil || tx.MaxFeePerGas == nil {
		b.logger.Warnw(fmt.Sprintf("Got transaction type %v but one of the required EIP1559 fields was missing, falling back to gasPrice", tx.Type), "block", block, "tx", tx)
		return tx.GasPrice
	}
	if tx.GasPrice != nil {
		// Always use the gas price if provided
		return tx.GasPrice
	}
	if tx.MaxFeePerGas.Cmp(block.BaseFeePerGas) < 0 {
		b.logger.AssumptionViolationw("MaxFeePerGas >= BaseFeePerGas", "block", block, "tx", tx)
		return nil
	}
	if tx.MaxFeePerGas.Cmp(tx.MaxPriorityFeePerGas) < 0 {
		b.logger.AssumptionViolationw("MaxFeePerGas >= MaxPriorityFeePerGas", "block", block, "tx", tx)
		return nil
	}

	// From: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md
	priorityFeePerGas := tx.MaxPriorityFeePerGas
	maxFeeMinusBaseFee := tx.MaxFeePerGas.Sub(block.BaseFeePerGas)
	if maxFeeMinusBaseFee.Cmp(priorityFeePerGas) < 0 {
		priorityFeePerGas = maxFeeMinusBaseFee
	}

	return priorityFeePerGas.Add(block.BaseFeePerGas)
}

func (b *BlockHistoryEstimator) EffectiveTipCap(block evmtypes.Block, tx evmtypes.Transaction) *assets.Wei {
	switch tx.Type {
	case 0x2, 0x3:
		return tx.MaxPriorityFeePerGas
	case 0x0, 0x1:
		if tx.GasPrice == nil {
			return nil
		}
		if block.BaseFeePerGas == nil {
			return nil
		}
		effectiveTipCap := tx.GasPrice.Sub(block.BaseFeePerGas)
		if effectiveTipCap.IsNegative() {
			b.logger.AssumptionViolationw("GasPrice - BaseFeePerGas may not be negative", "block", block, "tx", tx)
			return nil
		}
		return effectiveTipCap
	default:
		b.logger.Debugw(fmt.Sprintf("Ignoring unknown transaction type %v", tx.Type), "block", block, "tx", tx)
		return nil
	}
}
