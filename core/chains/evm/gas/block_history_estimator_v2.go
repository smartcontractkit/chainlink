package gas

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var (
	promBlockHistoryEstimatorV2AllGasCostPercentiles = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_all_gas_cost_percentiles",
		Help: "Gas price at given percentile",
	},
		[]string{"percentile", "evmChainID"},
	)

	promBlockHistoryEstimatorV2SetGasCost = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_set_gas_cost",
		Help: "Gas updater set gas cost (in Wei)",
	},
		[]string{"percentile", "evmChainID"},
	)
)

var _ EvmEstimator = &BlockHistoryEstimatorV2{}

type bheV2ChainConfig interface {
	ChainType() config.ChainType
}

type bheV2GasEstimatorConfig interface {
	EIP1559DynamicFees() bool
	BumpThreshold() uint64
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMax() *assets.Wei
	PriceMin() *assets.Wei
	bumpConfig
}

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore
type BlockHistoryEstimatorV2 struct {
	services.StateMachine
	ethClient feeEstimatorClient
	chainID   *big.Int
	config    bheV2ChainConfig
	eConfig   bheV2GasEstimatorConfig
	bhConfig  BlockHistoryConfig
	// NOTE: it is assumed that blocks will be kept sorted by
	// block number ascending
	blocks    []evmtypes.Block
	blocksMu  sync.RWMutex
	size      int64
	mb        *mailbox.Mailbox[*evmtypes.Head]
	wg        *sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc

	gasCost      *assets.Wei
	costMu       sync.RWMutex
	latest       *evmtypes.Head
	latestMu     sync.RWMutex
	initialFetch atomic.Bool

	logger logger.SugaredLogger

	l1Oracle rollups.L1Oracle
}

// NewBlockHistoryEstimator returns a new BlockHistoryEstimator that listens
// for new heads and updates the base gas price dynamically based on the
// configured percentile of gas prices in that block
func NewBlockHistoryEstimatorV2(lggr logger.Logger, ethClient feeEstimatorClient, cfg chainConfig, eCfg estimatorGasEstimatorConfig, bhCfg BlockHistoryConfig, chainID *big.Int, l1Oracle rollups.L1Oracle) EvmEstimator {
	ctx, cancel := context.WithCancel(context.Background())

	b := &BlockHistoryEstimatorV2{
		ethClient: ethClient,
		chainID:   chainID,
		config:    cfg,
		eConfig:   eCfg,
		bhConfig:  bhCfg,
		blocks:    make([]evmtypes.Block, 0),
		// Must have enough blocks for both estimator and connectivity checker
		size:      int64(mathutil.Max(bhCfg.BlockHistorySize(), bhCfg.CheckInclusionBlocks())),
		mb:        mailbox.NewSingle[*evmtypes.Head](),
		wg:        new(sync.WaitGroup),
		ctx:       ctx,
		ctxCancel: cancel,
		logger:    logger.Sugared(logger.Named(lggr, "BlockHistoryEstimatorV2")),
		l1Oracle:  l1Oracle,
	}

	return b
}

// OnNewLongestChain recalculates and sets global gas price if a sampled new head comes
// in and we are not currently fetching
func (b *BlockHistoryEstimatorV2) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	// set latest base fee here to avoid potential lag introduced by block delay
	// it is really important that base fee be as up-to-date as possible
	b.setLatest(head)
	b.mb.Deliver(head)
}

// setLatest assumes that head won't be mutated
func (b *BlockHistoryEstimatorV2) setLatest(head *evmtypes.Head) {
	// Non-eip1559 blocks don't include base fee
	if baseFee := head.BaseFeePerGas; baseFee != nil {
		promBlockHistoryEstimatorCurrentBaseFee.WithLabelValues(b.chainID.String()).Set(float64(baseFee.Int64()))
	}
	b.logger.Debugw("Set latest block", "blockNum", head.Number, "blockHash", head.Hash, "baseFee", head.BaseFeePerGas, "baseFeeWei", head.BaseFeePerGas.ToInt())
	b.latestMu.Lock()
	defer b.latestMu.Unlock()
	b.latest = head
}

func (b *BlockHistoryEstimatorV2) getCurrentBaseFee() *assets.Wei {
	b.latestMu.RLock()
	defer b.latestMu.RUnlock()
	if b.latest == nil {
		return nil
	}
	return b.latest.BaseFeePerGas
}

func (b *BlockHistoryEstimatorV2) getCurrentBlockNum() *int64 {
	b.latestMu.RLock()
	defer b.latestMu.RUnlock()
	if b.latest == nil {
		return nil
	}
	return &b.latest.Number
}

func (b *BlockHistoryEstimatorV2) getBlocks() []evmtypes.Block {
	b.blocksMu.RLock()
	defer b.blocksMu.RUnlock()
	return b.blocks
}

// Start starts BlockHistoryEstimatorV2 service.
// The provided context can be used to terminate Start sequence.
func (b *BlockHistoryEstimatorV2) Start(ctx context.Context) error {
	return b.StartOnce("BlockHistoryEstimatorV2", func() error {
		b.logger.Trace("Starting")

		if b.bhConfig.CheckInclusionBlocks() > 0 {
			b.logger.Infof("Inclusion checking enabled, bumping will be prevented on transactions that have been priced above the %d percentile for %d blocks", b.bhConfig.CheckInclusionPercentile(), b.bhConfig.CheckInclusionBlocks())
		}
		if b.bhConfig.BlockHistorySize() == 0 {
			return pkgerrors.New("BlockHistorySize must be set to a value greater than 0")
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
			return pkgerrors.Wrap(ctx.Err(), "failed to start BlockHistoryEstimatorV2 due to main context error")
		}

		b.wg.Add(1)
		go b.runLoop()

		b.logger.Trace("Started")
		return nil
	})
}

func (b *BlockHistoryEstimatorV2) L1Oracle() rollups.L1Oracle {
	return b.l1Oracle
}

func (b *BlockHistoryEstimatorV2) Close() error {
	return b.StopOnce("BlockHistoryEstimatorV2", func() error {
		b.ctxCancel()
		b.wg.Wait()
		return nil
	})
}

func (b *BlockHistoryEstimatorV2) Name() string {
	return b.logger.Name()
}
func (b *BlockHistoryEstimatorV2) HealthReport() map[string]error {
	return map[string]error{b.Name(): b.Healthy()}
}

func (b *BlockHistoryEstimatorV2) GetLegacyGas(_ context.Context, _ []byte, gasLimit uint64, maxGasPriceWei *assets.Wei, _ ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	var gasCost *assets.Wei
	ok := b.IfStarted(func() {
		gasCost = b.getGasCost()
	})
	if !ok {
		return nil, 0, pkgerrors.New("BlockHistoryEstimatorV2 is not started; cannot estimate gas")
	}
	if gasCost == nil {
		if !b.initialFetch.Load() {
			return nil, 0, pkgerrors.New("BlockHistoryEstimatorV2 has not finished the first gas estimation yet, likely because a failure on start")
		}
		b.logger.Warnw("Failed to estimate gas price. This is likely because there aren't any valid transactions to estimate from."+
			"Using Evm.GasEstimator.PriceDefault as fallback.", "blocks", b.getBlockHistoryNumbers())
		gasPrice = b.eConfig.PriceDefault()
	} else {
		// current gas price = cost / limit
		gasPrice = assets.NewWei(new(big.Int).Div(gasCost.ToInt(), new(big.Int).SetUint64(gasLimit)))
		// bound the gas price with the min and max configs
		gasPrice = b.boundFee(gasPrice, b.eConfig.PriceMin(), maxGasPriceWei, b.eConfig.PriceMax())
	}
	chainSpecificGasLimit = gasLimit
	return
}

func (b *BlockHistoryEstimatorV2) getGasCost() *assets.Wei {
	b.costMu.RLock()
	defer b.costMu.RUnlock()
	return b.gasCost
}

func (b *BlockHistoryEstimatorV2) getBlockHistoryNumbers() (numsInHistory []int64) {
	for _, b := range b.blocks {
		numsInHistory = append(numsInHistory, b.Number)
	}
	return
}

func (b *BlockHistoryEstimatorV2) BumpLegacyGas(_ context.Context, originalGasPrice *assets.Wei, gasLimit uint64, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumpedGasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	if b.bhConfig.CheckInclusionBlocks() > 0 {
		if err = b.checkConnectivity(attempts); err != nil {
			if pkgerrors.Is(err, commonfee.ErrConnectivity) {
				b.logger.Criticalw(BumpingHaltedLabel, "err", err)
				b.SvcErrBuffer.Append(err)
				promBlockHistoryEstimatorConnectivityFailureCount.WithLabelValues(b.chainID.String(), "legacy").Inc()
			}
			return nil, 0, err
		}
	}
	gasCost := b.getGasCost()
	var currentGasPrice *assets.Wei
	if gasCost != nil {
		// current gas price = cost / limit
		currentGasPrice = assets.NewWei(new(big.Int).Div(gasCost.ToInt(), new(big.Int).SetUint64(gasLimit)))
		// bound the gas price with the min and max configs
		currentGasPrice = b.boundFee(currentGasPrice, b.eConfig.PriceMin(), maxGasPriceWei, b.eConfig.PriceMax())
	}
	bumpedGasPrice, err = BumpLegacyGasPriceOnly(b.eConfig, b.logger, currentGasPrice, originalGasPrice, maxGasPriceWei)
	if err != nil {
		return nil, 0, err
	}
	return bumpedGasPrice, gasLimit, err
}

// checkConnectivity detects if the transaction is not being included due to
// some kind of mempool propagation or connectivity issue rather than
// insufficiently high pricing and returns error if so
func (b *BlockHistoryEstimatorV2) checkConnectivity(attempts []EvmPriorAttempt) error {
	percentile := int(b.bhConfig.CheckInclusionPercentile())
	// how many blocks since broadcast?
	latestBlockNum := b.getCurrentBlockNum()
	if latestBlockNum == nil {
		b.logger.Warn("Latest block is unknown; skipping inclusion check")
		// can't determine anything if we don't have/know latest block num yet
		return nil
	}
	expectInclusionWithinBlocks := int(b.bhConfig.CheckInclusionBlocks())
	blockHistory := b.getBlocks()
	if len(blockHistory) < expectInclusionWithinBlocks {
		b.logger.Warnf("Block history in memory with length %d is insufficient to determine whether transaction should have been included within the past %d blocks", len(blockHistory), b.bhConfig.CheckInclusionBlocks())
		return nil
	}
	for _, attempt := range attempts {
		if attempt.BroadcastBeforeBlockNum == nil {
			// this shouldn't happen; any broadcast attempt ought to have a
			// BroadcastBeforeBlockNum otherwise its an assumption violation
			return pkgerrors.Errorf("BroadcastBeforeBlockNum was unexpectedly nil for attempt %s", attempt.TxHash)
		}
		broadcastBeforeBlockNum := *attempt.BroadcastBeforeBlockNum
		blocksSinceBroadcast := *latestBlockNum - broadcastBeforeBlockNum
		if blocksSinceBroadcast < int64(expectInclusionWithinBlocks) {
			// only check attempts that have been waiting around longer than
			// CheckInclusionBlocks
			continue
		}
		// has not been included for at least the required number of blocks
		b.logger.Debugw(fmt.Sprintf("transaction %s has been pending inclusion for %d blocks which equals or exceeds expected specified check inclusion blocks of %d", attempt.TxHash, blocksSinceBroadcast, expectInclusionWithinBlocks), "broadcastBeforeBlockNum", broadcastBeforeBlockNum, "latestBlockNum", *latestBlockNum)
		// is the price in the right percentile for all of these blocks?
		var blocks []evmtypes.Block
		l := expectInclusionWithinBlocks
		// reverse order since we want to go highest -> lowest block number and bail out early
		for i := l - 1; i >= 0; i-- {
			block := blockHistory[i]
			if block.Number < broadcastBeforeBlockNum {
				break
			}
			blocks = append(blocks, block)
		}
		var eip1559 bool
		switch attempt.TxType {
		case 0x0, 0x1:
			eip1559 = false
		case 0x2:
			eip1559 = true
		default:
			return pkgerrors.Errorf("attempt %s has unknown transaction type 0x%d", attempt.TxHash, attempt.TxType)
		}
		gasCost, err := b.calculatePercentileCosts(blocks, percentile, eip1559, nil)
		if err != nil {
			if pkgerrors.Is(err, ErrNoSuitableTransactions) {
				b.logger.Warnf("no suitable transactions found to verify if transaction %s has been included within expected inclusion blocks of %d", attempt.TxHash, expectInclusionWithinBlocks)
				return nil
			}
			b.logger.AssumptionViolationw("unexpected error while verifying transaction inclusion", "err", err, "txHash", attempt.TxHash.String())
			return nil
		}
		attemptLimit := new(big.Int).SetUint64(attempt.ChainSpecificFeeLimit)
		if !eip1559 {
			// Calculate the gas cost of the attempt to compare against the percentile gas cost
			attemptGasCost := attempt.GasPrice.Mul(attemptLimit)
			if attemptGasCost.Cmp(gasCost) > 0 {
				return pkgerrors.Wrapf(commonfee.ErrConnectivity, "transaction %s has gas price of %s and gas limit of %d (gas cost: %s), which is above percentile=%d%% (percentile cost: %s) for blocks %d thru %d (checking %d blocks)", attempt.TxHash, attempt.GasPrice, attempt.ChainSpecificFeeLimit, attemptGasCost.String(), percentile, gasCost, blockHistory[l-1].Number, blockHistory[0].Number, expectInclusionWithinBlocks)
			}
			continue
		}
		sufficientFeeCap := true
		for _, b := range blocks {
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
		// Calculate the gas cost of the attempt to compare against the percentile gas cost
		attemptGasCost := attempt.DynamicFee.FeeCap.Mul(attemptLimit)
		if sufficientFeeCap && attemptGasCost.Cmp(gasCost) > 0 {
			return pkgerrors.Wrapf(commonfee.ErrConnectivity, "transaction %s has fee cap of %s and gas limit of %d (gas cost: %s), which is above percentile=%d%% (percentile gas cost: %s) for blocks %d thru %d (checking %d blocks)", attempt.TxHash, attempt.DynamicFee.FeeCap.String(), attempt.ChainSpecificFeeLimit, attemptGasCost.String(), percentile, gasCost, blockHistory[l-1].Number, blockHistory[0].Number, expectInclusionWithinBlocks)
		}
	}
	return nil
}

func (b *BlockHistoryEstimatorV2) GetDynamicFee(_ context.Context, gasLimit uint64, maxGasPriceWei *assets.Wei) (fee DynamicFee, err error) {
	if !b.eConfig.EIP1559DynamicFees() {
		return fee, pkgerrors.New("Can't get dynamic fee, EIP1559 is disabled")
	}

	var feeCap *assets.Wei
	var tipCap *assets.Wei
	ok := b.IfStarted(func() {
		b.costMu.RLock()
		defer b.costMu.RUnlock()
		gasCost := b.gasCost
		if gasCost == nil {
			if !b.initialFetch.Load() {
				err = pkgerrors.New("BlockHistoryEstimatorV2 has not finished the first gas estimation yet, likely because a failure on start")
				return
			}
			b.logger.Warnw("Failed to estimate gas price. This is likely because there aren't any valid transactions to estimate from."+
				"Using Evm.GasEstimator.TipCapDefault as fallback.", "blocks", b.getBlockHistoryNumbers())
			tipCap = b.eConfig.TipCapDefault()
		} else if b.getCurrentBaseFee() != nil {
			// If gasCost is not nil, use the current base fee and limit to calculate a tip cap based on the target percentile gas cost
			// gasCost = (baseFee + tip) * limit => tip = (gasCost / limit) - baseFee
			basePlusTipFee := assets.NewWei(new(big.Int).Div(gasCost.ToInt(), new(big.Int).SetUint64(gasLimit)))
			tipCap = basePlusTipFee.Sub(b.getCurrentBaseFee())
			// bound the tip cap with the min and max configs
			tipCap = b.boundFee(tipCap, b.eConfig.TipCapMin(), maxGasPriceWei, b.eConfig.PriceMax())
		} else {
			// This shouldn't happen on EIP-1559 blocks, since if the gas cost
			// is set, Start must have succeeded and we would expect an initial
			// base fee to be set as well
			err = pkgerrors.New("BlockHistoryEstimatorV2: no value for latest block base fee; cannot estimate EIP-1559 base fee. Are you trying to run with EIP1559 enabled on a non-EIP1559 chain?")
			return
		}
		maxGasPrice := getMaxGasPrice(maxGasPriceWei, b.eConfig.PriceMax())
		if b.eConfig.BumpThreshold() == 0 {
			// just use the max gas price if gas bumping is disabled
			feeCap = maxGasPrice
		} else {
			// HACK: due to a flaw of how EIP-1559 is implemented we have to
			// set a much lower FeeCap than the actual maximum we are willing
			// to pay in order to give ourselves headroom for bumping
			// See: https://github.com/ethereum/go-ethereum/issues/24284
			feeCap = calcFeeCap(b.getCurrentBaseFee(), int(b.bhConfig.EIP1559FeeCapBufferBlocks()), tipCap, maxGasPrice)
		}
	})
	if !ok {
		return fee, pkgerrors.New("BlockHistoryEstimatorV2 is not started; cannot estimate gas")
	}
	if err != nil {
		return fee, err
	}
	fee.FeeCap = feeCap
	fee.TipCap = tipCap
	return
}

func (b *BlockHistoryEstimatorV2) BumpDynamicFee(_ context.Context, originalFee DynamicFee, gasLimit uint64, maxGasPriceWei *assets.Wei, attempts []EvmPriorAttempt) (bumped DynamicFee, err error) {
	if b.bhConfig.CheckInclusionBlocks() > 0 {
		if err = b.checkConnectivity(attempts); err != nil {
			if pkgerrors.Is(err, commonfee.ErrConnectivity) {
				b.logger.Criticalw(BumpingHaltedLabel, "err", err)
				b.SvcErrBuffer.Append(err)
				promBlockHistoryEstimatorConnectivityFailureCount.WithLabelValues(b.chainID.String(), "eip1559").Inc()
			}
			return bumped, err
		}
	}
	gasCost := b.getGasCost()
	baseFee := b.getCurrentBaseFee()
	var tipCap *assets.Wei
	if gasCost != nil && baseFee != nil {
		basePlusTipFee := assets.NewWei(new(big.Int).Div(gasCost.ToInt(), new(big.Int).SetUint64(gasLimit)))
		tipCap = basePlusTipFee.Sub(baseFee)
		// bound the tip cap with the min and max configs
		tipCap = b.boundFee(tipCap, b.eConfig.TipCapMin(), maxGasPriceWei, b.eConfig.PriceMax())
	}
	return BumpDynamicFeeOnly(b.eConfig, b.bhConfig.EIP1559FeeCapBufferBlocks(), b.logger, tipCap, baseFee, originalFee, maxGasPriceWei)
}

func (b *BlockHistoryEstimatorV2) runLoop() {
	defer b.wg.Done()
	for {
		select {
		case <-b.ctx.Done():
			return
		case <-b.mb.Notify():
			head, exists := b.mb.Retrieve()
			if !exists {
				b.logger.Debug("No head to retrieve")
				continue
			}
			b.FetchBlocksAndRecalculate(b.ctx, head)
		}
	}
}

// FetchBlocksAndRecalculate fetches block history leading up to head and recalculates gas price.
func (b *BlockHistoryEstimatorV2) FetchBlocksAndRecalculate(ctx context.Context, head *evmtypes.Head) {
	if err := b.FetchBlocks(ctx, head); err != nil {
		b.logger.Warnw("Error fetching blocks", "head", head, "err", err)
		return
	}
	b.initialFetch.Store(true)
	b.Recalculate(head)
}

// Recalculate adds the given heads to the history and recalculates gas price.
func (b *BlockHistoryEstimatorV2) Recalculate(head *evmtypes.Head) {
	percentile := int(b.bhConfig.TransactionPercentile())

	lggr := b.logger.With("head", head)

	blockHistory := b.getBlocks()
	if len(blockHistory) == 0 {
		lggr.Debug("No blocks in history, cannot set gas price")
		return
	}

	l := mathutil.Min(len(blockHistory), int(b.bhConfig.BlockHistorySize()))
	blocks := blockHistory[:l]

	eip1559 := b.eConfig.EIP1559DynamicFees()
	percentileGasCost, err := b.calculatePercentileCosts(blocks, percentile, eip1559,
		func(gasPrices []*assets.Wei) {
			for i := 0; i <= 100; i += 5 {
				jdx := ((len(gasPrices) - 1) * i) / 100
				promBlockHistoryEstimatorV2AllGasCostPercentiles.WithLabelValues(fmt.Sprintf("%v%%", i), b.chainID.String()).Set(float64(gasPrices[jdx].Int64()))
			}
		})
	if err != nil {
		if pkgerrors.Is(err, ErrNoSuitableTransactions) {
			lggr.Debug("No suitable transactions, skipping")
		} else {
			lggr.Warnw("Cannot calculate percentile prices", "err", err)
		}
		return
	}

	var numsInHistory []int64
	for _, b := range blockHistory {
		numsInHistory = append(numsInHistory, b.Number)
	}

	float := new(big.Float).SetInt(percentileGasCost.ToInt())
	gwei, _ := big.NewFloat(0).Quo(float, big.NewFloat(1000000000)).Float64()
	gasCostGwei := fmt.Sprintf("%.2f", gwei)

	lggrFields := []interface{}{
		"gasCostWei", percentileGasCost,
		"gasCostGWei", gasCostGwei,
		"headNum", head.Number,
		"blocks", numsInHistory,
	}
	b.setPercentileGasCost(percentileGasCost)
	promBlockHistoryEstimatorV2SetGasCost.WithLabelValues(fmt.Sprintf("%v%%", percentile), b.chainID.String()).Set(float64(percentileGasCost.Int64()))
	lggr.Debugw(fmt.Sprintf("Setting new default gas cost: %v Gwei", gasCostGwei), lggrFields...)
}

// FetchBlocks fetches block history leading up to the given head.
func (b *BlockHistoryEstimatorV2) FetchBlocks(ctx context.Context, head *evmtypes.Head) error {
	// HACK: blockDelay is the number of blocks that the block history estimator trails behind head.
	// E.g. if this is set to 3, and we receive block 10, block history estimator will
	// fetch block 7.
	// This is necessary because geth/parity send heads as soon as they get
	// them and often the actual block is not available until later. Fetching
	// it too early results in an empty block.
	blockDelay := int64(b.bhConfig.BlockDelay())
	historySize := b.size

	if historySize <= 0 {
		return pkgerrors.Errorf("BlockHistoryEstimatorV2: history size must be > 0, got: %d", historySize)
	}

	highestBlockToFetch := head.Number - blockDelay
	if highestBlockToFetch < 0 {
		return pkgerrors.Errorf("BlockHistoryEstimatorV2: cannot fetch, current block height %v is lower than EVM.RPCBlockQueryDelay=%v", head.Number, blockDelay)
	}
	lowestBlockToFetch := head.Number - historySize - blockDelay + 1
	if lowestBlockToFetch < 0 {
		lowestBlockToFetch = 0
	}

	blocks := make(map[int64]evmtypes.Block)
	for _, block := range b.getBlocks() {
		// Make a best-effort to be re-org resistant using the head
		// chain, refetch blocks that got re-org'd out.
		// NOTE: Any blocks in the history that are older than the oldest block
		// in the provided chain will be assumed final.
		if block.Number < head.EarliestInChain().BlockNumber() {
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
			if pkgerrors.Is(err, evmtypes.ErrMissingBlock) {
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
			return pkgerrors.Errorf("expected result to be a %T, got %T", &evmtypes.Block{}, result)
		}
		if block == nil {
			return pkgerrors.New("invariant violation: got nil block")
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

func (b *BlockHistoryEstimatorV2) batchFetch(ctx context.Context, reqs []rpc.BatchElem) error {
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
		if pkgerrors.Is(err, context.DeadlineExceeded) {
			// We ran out of time, return what we have
			b.logger.Warnf("Batch fetching timed out; loaded %d/%d results: %v", i, len(reqs), err)
			for k := i; k < len(reqs); k++ {
				if k < j {
					reqs[k].Error = pkgerrors.Wrap(err, "request failed")
				} else {
					reqs[k].Error = pkgerrors.Wrap(err, "request skipped; previous request exceeded deadline")
				}
			}
			return nil
		} else if err != nil {
			return pkgerrors.Wrap(err, "BlockHistoryEstimatorV2#fetchBlocks error fetching blocks with BatchCallContext")
		}
	}
	return nil
}

func (b *BlockHistoryEstimatorV2) calculatePercentileCosts(blocks []evmtypes.Block, percentile int, eip1559 bool, f func(gasPrices []*assets.Wei)) (gasCost *assets.Wei, err error) {
	gasCosts := b.getGasCostsFromBlocks(blocks)
	if len(gasCosts) == 0 {
		return nil, ErrNoSuitableTransactions
	}
	sort.Slice(gasCosts, func(i, j int) bool { return gasCosts[i].Cmp(gasCosts[j]) < 0 })
	if f != nil {
		f(gasCosts)
	}
	gasCost = gasCosts[((len(gasCosts)-1)*percentile)/100]

	if !eip1559 {
		return
	}

	return
}

func (b *BlockHistoryEstimatorV2) getGasCostsFromBlocks(blocks []evmtypes.Block) (gasCosts []*assets.Wei) {
	gasCosts = make([]*assets.Wei, 0)
	for _, block := range blocks {
		// Unlike BHEv1, the verify block check is not needed here. Using gas costs allows us to use transactions that have gas price set.
		// Checks for nil blockBaseFee is done downstream
		for _, tx := range block.Transactions {
			if b.IsUsable(tx, block, b.config.ChainType(), b.eConfig.PriceMin(), b.logger) {
				gasCost := b.EffectiveGasCost(block, tx)
				if gasCost == nil {
					b.logger.Warnw("Unable to get gas cost for tx", "tx", tx, "block", block)
					continue
				}
				gasCosts = append(gasCosts, gasCost)
			}
		}
	}
	return
}

func (b *BlockHistoryEstimatorV2) setPercentileGasCost(gasCost *assets.Wei) {
	b.costMu.Lock()
	defer b.costMu.Unlock()
	b.gasCost = gasCost
}

// isUsable returns true if the tx is usable both generally and specifically for
// this Config.
func (b *BlockHistoryEstimatorV2) IsUsable(tx evmtypes.Transaction, block evmtypes.Block, chainType config.ChainType, minGasPrice *assets.Wei, lggr logger.Logger) bool {
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

// Gets the effective gas cost for transactions of any type
func (b *BlockHistoryEstimatorV2) EffectiveGasCost(block evmtypes.Block, tx evmtypes.Transaction) *assets.Wei {
	switch tx.Type {
	case 0x0, 0x1:
		return b.getLegacyGasCost(tx)
	case 0x2, 0x3:
		return b.getDynamicGasCost(block, tx)
	default:
		b.logger.Debugw(fmt.Sprintf("Ignoring unknown transaction type %v", tx.Type), "block", block, "tx", tx)
		return nil
	}
}

// Gets the gas cost for dynamic transactions using EIP1559 fields. Fallsback to GasPrice if it exists.
func (b *BlockHistoryEstimatorV2) getLegacyGasCost(tx evmtypes.Transaction) *assets.Wei {
	if tx.GasPrice == nil {
		b.logger.Warnw(fmt.Sprintf("Got transaction type %v but gas price was missing", tx.Type), "tx", tx)
		return nil
	}
	return tx.GasPrice.Mul(big.NewInt(int64(tx.GasLimit)))
}

// Gets the gas cost for dynamic transactions using EIP1559 fields. Fallsback to GasPrice if it exists.
func (b *BlockHistoryEstimatorV2) getDynamicGasCost(block evmtypes.Block, tx evmtypes.Transaction) *assets.Wei {
	// Used as fallback if it exists and EIP-1559 fields are invalid
	var gasCostUsingPrice *assets.Wei
	if tx.GasPrice != nil {
		gasCostUsingPrice = tx.GasPrice.Mul(big.NewInt(int64(tx.GasLimit)))
	}
	if block.BaseFeePerGas == nil || tx.MaxPriorityFeePerGas == nil || tx.MaxFeePerGas == nil {
		b.logger.Warnw(fmt.Sprintf("Got transaction type %v but one of the required EIP1559 fields was missing", tx.Type), "block", block, "tx", tx)
		return gasCostUsingPrice
	}
	if tx.MaxFeePerGas.Cmp(block.BaseFeePerGas) < 0 {
		b.logger.AssumptionViolationw("Expect MaxFeePerGas >= BaseFeePerGas", "block", block, "tx", tx)
		return gasCostUsingPrice
	}
	if tx.MaxFeePerGas.Cmp(tx.MaxPriorityFeePerGas) < 0 {
		b.logger.AssumptionViolationw("Expect MaxFeePerGas >= MaxPriorityFeePerGas", "block", block, "tx", tx)
		return gasCostUsingPrice
	}

	// Use MaxFeePerGas if the sum of block BaseFeePerGas and tx MaxPriorityFeePerGas exceeds it
	feePerGas := assets.WeiMin(tx.MaxFeePerGas, block.BaseFeePerGas.Add(tx.MaxPriorityFeePerGas))
	return feePerGas.Mul(big.NewInt(int64(tx.GasLimit)))
}

func (b *BlockHistoryEstimatorV2) boundFee(calculatedFee, min, userSpecifiedMaxFee, maxFeeWei *assets.Wei) *assets.Wei {
	fieldName := "gas price"
	if b.eConfig.EIP1559DynamicFees() {
		fieldName = "tip cap"
	}
	max := getMaxGasPrice(userSpecifiedMaxFee, maxFeeWei)
	fee := calculatedFee
	if calculatedFee.Cmp(max) > 0 {
		fee = max
		b.logger.Warnf("Calculated %[1]s of %[2]s exceeds the configured max: %[3]s, setting %[1]s to the maximum allowed value of %[3]s instead", fieldName, fee.String(), max.String())
	} else if calculatedFee.Cmp(min) < 0 {
		fee = min
		b.logger.Warnf("Calculated %[1]s of %[2]s falls below the configured min: %[3]s, setting %[1]s to the minimum allowed value of %[3]s instead", fieldName, fee.String(), min.String())
	}
	return fee
}
