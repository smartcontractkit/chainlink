package gas

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
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
)

var _ Estimator = &BlockHistoryEstimator{}

//go:generate mockery --name Config --output ./mocks/ --case=underscore
type (
	BlockHistoryEstimator struct {
		utils.StartStopOnce
		ethClient           evmclient.Client
		chainID             big.Int
		config              Config
		rollingBlockHistory []Block
		mb                  *utils.Mailbox[*evmtypes.Head]
		wg                  *sync.WaitGroup
		ctx                 context.Context
		ctxCancel           context.CancelFunc

		gasPrice      *big.Int
		tipCap        *big.Int
		latestBaseFee *big.Int
		mu            sync.RWMutex

		logger logger.SugaredLogger
	}
)

// NewBlockHistoryEstimator returns a new BlockHistoryEstimator that listens
// for new heads and updates the base gas price dynamically based on the
// configured percentile of gas prices in that block
func NewBlockHistoryEstimator(lggr logger.Logger, ethClient evmclient.Client, cfg Config, chainID big.Int) Estimator {
	ctx, cancel := context.WithCancel(context.Background())
	b := &BlockHistoryEstimator{
		utils.StartStopOnce{},
		ethClient,
		chainID,
		cfg,
		make([]Block, 0),
		utils.NewMailbox[*evmtypes.Head](1),
		new(sync.WaitGroup),
		ctx,
		cancel,
		nil,
		nil,
		nil,
		sync.RWMutex{},
		logger.Sugared(lggr.Named("BlockHistoryEstimator")),
	}

	return b
}

// OnNewLongestChain recalculates and sets global gas price if a sampled new head comes
// in and we are not currently fetching
func (b *BlockHistoryEstimator) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	// set latest base fee here to avoid potential lag introduced by block delay
	// it is really important that base fee be as up-to-date as possible
	b.setLatestBaseFee(head.BaseFeePerGas)
	b.mb.Deliver(head)
}

func (b *BlockHistoryEstimator) setLatestBaseFee(baseFee *utils.Big) {
	// Non-eip1559 blocks don't include base fee; just ignore
	if baseFee == nil {
		return
	}
	promBlockHistoryEstimatorCurrentBaseFee.WithLabelValues(b.chainID.String()).Set(float64(baseFee.ToInt().Int64()))
	b.mu.Lock()
	defer b.mu.Unlock()
	b.latestBaseFee = new(big.Int)
	b.latestBaseFee.Set(baseFee.ToInt())
}
func (b *BlockHistoryEstimator) getCurrentBaseFee() *big.Int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.latestBaseFee
}

// Start starts BlockHistoryEstimator service.
// The provided context can be used to terminate Start sequence.
func (b *BlockHistoryEstimator) Start(ctx context.Context) error {
	return b.StartOnce("BlockHistoryEstimator", func() error {
		b.logger.Trace("Starting")

		fetchCtx, cancel := context.WithTimeout(ctx, MaxStartTime)
		defer cancel()
		latestHead, err := b.ethClient.HeadByNumber(fetchCtx, nil)
		if err != nil {
			b.logger.Warnw("Initial check for latest head failed", "err", err)
		} else if latestHead == nil {
			b.logger.Warnw("initial check for latest head failed, head was unexpectedly nil")
		} else {
			b.logger.Debugw("Got latest head", "number", latestHead.Number, "blockHash", latestHead.Hash.Hex())
			b.setLatestBaseFee(latestHead.BaseFeePerGas)
			b.FetchBlocksAndRecalculate(fetchCtx, latestHead)
		}

		// NOTE: This only checks the start context, not the fetch context
		if ctx.Err() != nil {
			return errors.Wrap(ctx.Err(), "failed to start BlockHistoryEstimator due to main context error")
		}

		b.wg.Add(1)
		go b.runLoop()

		b.logger.Trace("Started")
		return nil
	})
}

func (b *BlockHistoryEstimator) Close() error {
	return b.StopOnce("BlockHistoryEstimator", func() error {
		b.ctxCancel()
		b.wg.Wait()
		return nil
	})
}

func (b *BlockHistoryEstimator) GetLegacyGas(_ []byte, gasLimit uint64, _ ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	ok := b.IfStarted(func() {
		chainSpecificGasLimit = applyMultiplier(gasLimit, b.config.EvmGasLimitMultiplier())
		gasPrice = b.getGasPrice()
	})
	if !ok {
		return nil, 0, errors.New("BlockHistoryEstimator is not started; cannot estimate gas")
	}
	if gasPrice == nil {
		return nil, 0, errors.New("BlockHistoryEstimator has not finished the first gas estimation yet, likely because a failure on start")
	}
	return
}

func (b *BlockHistoryEstimator) getGasPrice() *big.Int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.gasPrice
}
func (b *BlockHistoryEstimator) getTipCap() *big.Int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.tipCap
}

func (b *BlockHistoryEstimator) BumpLegacyGas(originalGasPrice *big.Int, gasLimit uint64) (bumpedGasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	return BumpLegacyGasPriceOnly(b.config, b.logger, b.getGasPrice(), originalGasPrice, gasLimit)
}

func (b *BlockHistoryEstimator) GetDynamicFee(gasLimit uint64) (fee DynamicFee, chainSpecificGasLimit uint64, err error) {
	if !b.config.EvmEIP1559DynamicFees() {
		return fee, 0, errors.New("Can't get dynamic fee, EIP1559 is disabled")
	}

	var feeCap *big.Int
	var tipCap *big.Int
	ok := b.IfStarted(func() {
		chainSpecificGasLimit = applyMultiplier(gasLimit, b.config.EvmGasLimitMultiplier())
		b.mu.RLock()
		defer b.mu.RUnlock()
		tipCap = b.tipCap
		if tipCap == nil {
			err = errors.New("BlockHistoryEstimator has not finished the first gas estimation yet, likely because a failure on start")
			return
		}
		if b.config.EvmGasBumpThreshold() == 0 {
			// just use the max gas price if gas bumping is disabled
			feeCap = b.config.EvmMaxGasPriceWei()
		} else if b.latestBaseFee != nil {
			// HACK: due to a flaw of how EIP-1559 is implemented we have to
			// set a much lower FeeCap than the actual maximum we are willing
			// to pay in order to give ourselves headroom for bumping
			// See: https://github.com/ethereum/go-ethereum/issues/24284
			feeCap = calcFeeCap(b.latestBaseFee, b.config, tipCap)
		} else {
			// This shouldn't happen on EIP-1559 blocks, since if the tip cap
			// is set, Start must have succeeded and we would expect an initial
			// base fee to be set as well
			err = errors.New("BlockHistoryEstimator: no value for latest block base fee; cannot estimate EIP-1559 base fee. Are you trying to run with EIP1559 enabled on a non-EIP1559 chain?")
			return
		}
	})
	if !ok {
		return fee, 0, errors.New("BlockHistoryEstimator is not started; cannot estimate gas")
	}
	if err != nil {
		return fee, 0, err
	}
	fee.FeeCap = feeCap
	fee.TipCap = tipCap
	return
}

func calcFeeCap(latestAvailableBaseFeePerGas *big.Int, cfg Config, tipCap *big.Int) (feeCap *big.Int) {
	const maxBaseFeeIncreasePerBlock float64 = 1.125

	bufferBlocks := int(cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks())

	baseFee := new(big.Float)
	baseFee.SetInt(latestAvailableBaseFeePerGas)
	// Find out the worst case base fee before we should bump
	multiplier := big.NewFloat(maxBaseFeeIncreasePerBlock)
	for i := 0; i < bufferBlocks; i++ {
		baseFee.Mul(baseFee, multiplier)
	}

	baseFeeInt, _ := baseFee.Int(nil)
	feeCap = new(big.Int).Add(baseFeeInt, tipCap)

	if feeCap.Cmp(cfg.EvmMaxGasPriceWei()) > 0 {
		return cfg.EvmMaxGasPriceWei()
	}
	return feeCap
}

func (b *BlockHistoryEstimator) BumpDynamicFee(originalFee DynamicFee, originalGasLimit uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	return BumpDynamicFeeOnly(b.config, b.logger, b.getTipCap(), b.getCurrentBaseFee(), originalFee, originalGasLimit)
}

func (b *BlockHistoryEstimator) runLoop() {
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
func (b *BlockHistoryEstimator) FetchBlocksAndRecalculate(ctx context.Context, head *evmtypes.Head) {
	if err := b.FetchBlocks(ctx, head); err != nil {
		b.logger.Warnw("Error fetching blocks", "head", head, "err", err)
		return
	}

	b.Recalculate(head)
}

// Recalculate adds the given heads to the history and recalculates gas price.
func (b *BlockHistoryEstimator) Recalculate(head *evmtypes.Head) {
	enableEIP1559 := b.config.EvmEIP1559DynamicFees()

	percentile := int(b.config.BlockHistoryEstimatorTransactionPercentile())

	lggr := b.logger.With("head", head)

	if len(b.rollingBlockHistory) == 0 {
		lggr.Debug("No blocks in history, cannot set gas price")
		return
	}

	percentileGasPrice, percentileTipCap, err := b.percentilePrices(percentile, enableEIP1559)
	if err != nil {
		if errors.Is(err, ErrNoSuitableTransactions) {
			lggr.Debug("No suitable transactions, skipping")
		} else {
			lggr.Warnw("Cannot calculate percentile prices", "err", err)
		}
		return
	}

	var numsInHistory []int64
	for _, b := range b.rollingBlockHistory {
		numsInHistory = append(numsInHistory, b.Number)
	}

	float := new(big.Float).SetInt(percentileGasPrice)
	gwei, _ := big.NewFloat(0).Quo(float, big.NewFloat(1000000000)).Float64()
	gasPriceGwei := fmt.Sprintf("%.2f", gwei)

	lggrFields := []interface{}{
		"gasPriceWei", percentileGasPrice,
		"gasPriceGWei", gasPriceGwei,
		"maxGasPriceWei", b.config.EvmMaxGasPriceWei(),
		"headNum", head.Number,
		"blocks", numsInHistory,
	}
	b.setPercentileGasPrice(percentileGasPrice)
	promBlockHistoryEstimatorSetGasPrice.WithLabelValues(fmt.Sprintf("%v%%", percentile), b.chainID.String()).Set(float64(percentileGasPrice.Int64()))

	if enableEIP1559 {
		float = new(big.Float).SetInt(percentileTipCap)
		gwei, _ = big.NewFloat(0).Quo(float, big.NewFloat(1000000000)).Float64()
		tipCapGwei := fmt.Sprintf("%.2f", gwei)
		lggrFields = append(lggrFields, []interface{}{
			"tipCapWei", percentileTipCap,
			"tipCapGwei", tipCapGwei,
		}...)
		lggr.Debugw(fmt.Sprintf("Setting new default prices, GasPrice: %v Gwei, TipCap: %v Gwei", gasPriceGwei, tipCapGwei), lggrFields...)
		b.setPercentileTipCap(percentileTipCap)
		promBlockHistoryEstimatorSetTipCap.WithLabelValues(fmt.Sprintf("%v%%", percentile), b.chainID.String()).Set(float64(percentileTipCap.Int64()))
	} else {
		lggr.Debugw(fmt.Sprintf("Setting new default gas price: %v Gwei", gasPriceGwei), lggrFields...)
	}
}

// FetchBlocks fetches block history leading up to the given head.
func (b *BlockHistoryEstimator) FetchBlocks(ctx context.Context, head *evmtypes.Head) error {
	// HACK: blockDelay is the number of blocks that the block history estimator trails behind head.
	// E.g. if this is set to 3, and we receive block 10, block history estimator will
	// fetch block 7.
	// This is necessary because geth/parity send heads as soon as they get
	// them and often the actual block is not available until later. Fetching
	// it too early results in an empty block.
	blockDelay := int64(b.config.BlockHistoryEstimatorBlockDelay())
	historySize := int64(b.config.BlockHistoryEstimatorBlockHistorySize())

	if historySize <= 0 {
		return errors.Errorf("BlockHistoryEstimator: history size must be > 0, got: %d", historySize)
	}

	highestBlockToFetch := head.Number - blockDelay
	if highestBlockToFetch < 0 {
		return errors.Errorf("BlockHistoryEstimator: cannot fetch, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY=%v", head.Number, blockDelay)
	}
	lowestBlockToFetch := head.Number - historySize - blockDelay + 1
	if lowestBlockToFetch < 0 {
		lowestBlockToFetch = 0
	}

	blocks := make(map[int64]Block)
	for _, block := range b.rollingBlockHistory {
		// Make a best-effort to be re-org resistant using the head
		// chain, refetch blocks that got re-org'd out.
		// NOTE: Any blocks in the history that are older than the oldest block
		// in the provided chain will be assumed final.
		if block.Number < head.EarliestInChain().Number {
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
			Result: &Block{},
		}
		reqs = append(reqs, req)
	}

	lggr := b.logger.With("head", head)

	lggr.Tracew(fmt.Sprintf("Fetching %v blocks (%v in local history)", len(reqs), len(blocks)), "n", len(reqs), "inHistory", len(blocks), "blockNum", head.Number)
	if err := b.batchFetch(ctx, reqs); err != nil {
		return err
	}

	for _, req := range reqs {
		result, err := req.Result, req.Error
		if err != nil {
			if strings.Contains(err.Error(), "failed to decode block number while unmarshalling block") {
				lggr.Errorw(
					fmt.Sprintf("Failed to fetch block: RPC node returned an empty block on query for block number %d even though the WS subscription already sent us this block. It might help to increase BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY (currently %d)",
						HexToInt64(req.Args[0]), blockDelay,
					),
					"err", err, "blockNum", HexToInt64(req.Args[0]), "headNum", head.Number)
			} else {
				lggr.Warnw("Failed to fetch block", "err", err, "blockNum", HexToInt64(req.Args[0]), "headNum", head.Number)
			}
			continue
		}

		block, is := result.(*Block)
		if !is {
			return errors.Errorf("expected result to be a %T, got %T", &Block{}, result)
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

	newBlockHistory := make([]Block, 0)
	for _, block := range blocks {
		newBlockHistory = append(newBlockHistory, block)
	}
	sort.Slice(newBlockHistory, func(i, j int) bool {
		return newBlockHistory[i].Number < newBlockHistory[j].Number
	})

	start := len(newBlockHistory) - int(historySize)
	if start < 0 {
		lggr.Debugw(fmt.Sprintf("Using fewer blocks than the specified history size: %v/%v", len(newBlockHistory), historySize), "rollingBlockHistorySize", historySize, "headNum", head.Number, "blocksAvailable", len(newBlockHistory))
		start = 0
	}

	b.rollingBlockHistory = newBlockHistory[start:]

	return nil
}

func (b *BlockHistoryEstimator) batchFetch(ctx context.Context, reqs []rpc.BatchElem) error {
	batchSize := int(b.config.BlockHistoryEstimatorBatchSize())

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
			b.logger.Warnf("Batch fetching timed out; loaded %d/%d results", i, len(reqs))
			for k := i; k < len(reqs); k++ {
				if k < j {
					reqs[k].Error = errors.Wrap(err, "request failed")
				} else {
					reqs[k].Error = errors.Wrap(err, "request skipped; previous request exceeded deadline")
				}
			}
			return nil
		} else if err != nil {
			return errors.Wrap(err, "BlockHistoryEstimator#fetchBlocks error fetching blocks with BatchCallContext")
		}
	}
	return nil
}

var (
	ErrNoSuitableTransactions = errors.New("no suitable transactions")
)

func (b *BlockHistoryEstimator) percentilePrices(percentile int, eip1559 bool) (gasPrice, tipCap *big.Int, err error) {
	gasPrices := make([]*big.Int, 0)
	tipCaps := make([]*big.Int, 0)
	for _, block := range b.rollingBlockHistory {
		if err := verifyBlock(block, eip1559); err != nil {
			b.logger.Warnw(fmt.Sprintf("Block %v is not usable, %s", block.Number, err.Error()), "block", block, "err", err)
		}
		for _, tx := range block.Transactions {
			if tx.isUsable(b.config, b.logger) {
				gp := b.EffectiveGasPrice(block, tx)
				if gp != nil {
					gasPrices = append(gasPrices, gp)
				} else {
					b.logger.Warnw("Unable to get gas price for tx", "tx", tx, "block", block)
					continue
				}
				if eip1559 {
					tc := b.EffectiveTipCap(block, tx)
					if tc != nil {
						tipCaps = append(tipCaps, tc)
					} else {
						b.logger.Warnw("Unable to get tip cap for tx", "tx", tx, "block", block)
						continue
					}
				}
			}
		}
	}
	if len(gasPrices) == 0 {
		return nil, nil, ErrNoSuitableTransactions
	}
	sort.Slice(gasPrices, func(i, j int) bool { return gasPrices[i].Cmp(gasPrices[j]) < 0 })
	idx := ((len(gasPrices) - 1) * percentile) / 100
	for i := 0; i <= 100; i += 5 {
		jdx := ((len(gasPrices) - 1) * i) / 100
		promBlockHistoryEstimatorAllGasPricePercentiles.WithLabelValues(fmt.Sprintf("%v%%", i), b.chainID.String()).Set(float64(gasPrices[jdx].Int64()))
	}
	gasPrice = gasPrices[idx]

	if eip1559 && len(tipCaps) == 0 {
		return nil, nil, ErrNoSuitableTransactions
	} else if eip1559 {
		sort.Slice(tipCaps, func(i, j int) bool { return tipCaps[i].Cmp(tipCaps[j]) < 0 })
		idx := ((len(tipCaps) - 1) * percentile) / 100
		for i := 0; i <= 100; i += 5 {
			jdx := ((len(tipCaps) - 1) * i) / 100
			promBlockHistoryEstimatorAllTipCapPercentiles.WithLabelValues(fmt.Sprintf("%v%%", i), b.chainID.String()).Set(float64(tipCaps[jdx].Int64()))
		}
		tipCap = tipCaps[idx]
	}

	return
}

func verifyBlock(block Block, eip1559 bool) error {
	if eip1559 && block.BaseFeePerGas == nil {
		return errors.New("EIP-1559 mode was enabled, but block was missing baseFeePerGas")
	}
	return nil
}

func (b *BlockHistoryEstimator) setPercentileTipCap(tipCap *big.Int) {
	min := b.config.EvmGasTipCapMinimum()

	b.mu.Lock()
	defer b.mu.Unlock()
	if tipCap.Cmp(min) < 0 {
		b.logger.Warnw(fmt.Sprintf("Calculated gas tip cap of %s Wei falls below EVM_GAS_TIP_CAP_MINIMUM=%[2]s, setting gas tip cap to the minimum allowed value of %[2]s Wei instead", tipCap.String(), min.String()), "tipCapWei", tipCap, "minTipCapWei", min)
		b.tipCap = min
	} else {
		b.tipCap = tipCap
	}
}

func (b *BlockHistoryEstimator) setPercentileGasPrice(gasPrice *big.Int) {
	max := b.config.EvmMaxGasPriceWei()
	min := b.config.EvmMinGasPriceWei()

	b.mu.Lock()
	defer b.mu.Unlock()
	if gasPrice.Cmp(max) > 0 {
		b.logger.Warnw(fmt.Sprintf("Calculated gas price of %s Wei exceeds ETH_MAX_GAS_PRICE_WEI=%[2]s, setting gas price to the maximum allowed value of %[2]s Wei instead", gasPrice.String(), max.String()), "gasPriceWei", gasPrice, "maxGasPriceWei", max)
		b.gasPrice = max
	} else if gasPrice.Cmp(min) < 0 {
		b.logger.Warnw(fmt.Sprintf("Calculated gas price of %s Wei falls below ETH_MIN_GAS_PRICE_WEI=%[2]s, setting gas price to the minimum allowed value of %[2]s Wei instead", gasPrice.String(), min.String()), "gasPriceWei", gasPrice, "minGasPriceWei", min)
		b.gasPrice = min
	} else {
		b.gasPrice = gasPrice
	}
}

func (b *BlockHistoryEstimator) RollingBlockHistory() []Block {
	return b.rollingBlockHistory
}

// isUsable returns true if the tx is usable both generally and specifically for
// this Config.
func (tx *Transaction) isUsable(cfg Config, lggr logger.Logger) bool {
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
	return tx.chainSpecificIsUsable(cfg)
}

func (b *BlockHistoryEstimator) EffectiveGasPrice(block Block, tx Transaction) *big.Int {
	switch tx.Type {
	case 0x0, 0x1:
		return tx.GasPrice
	case 0x2:
		if block.BaseFeePerGas == nil || tx.MaxPriorityFeePerGas == nil || tx.MaxFeePerGas == nil {
			b.logger.Warnw("Got transaction type 0x2 but one of the required EIP1559 fields was missing, falling back to gasPrice", "block", block, "tx", tx)
			return tx.GasPrice
		}
		if tx.MaxFeePerGas.Cmp(block.BaseFeePerGas) < 0 {
			// This should not pass config validation
			b.logger.AssumptionViolationw("MaxFeePerGas >= BaseFeePerGas", "block", block, "tx", tx)
			return nil
		}
		if tx.MaxFeePerGas.Cmp(tx.MaxPriorityFeePerGas) < 0 {
			// This should not pass config validation
			b.logger.AssumptionViolationw("MaxFeePerGas >= MaxPriorityFeePerGas", "block", block, "tx", tx)
			return nil
		}
		if tx.GasPrice != nil {
			// Always use the gas price if provided
			return tx.GasPrice
		}

		// From: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md
		priorityFeePerGas := tx.MaxPriorityFeePerGas
		maxFeeMinusBaseFee := big.NewInt(0).Sub(tx.MaxFeePerGas, block.BaseFeePerGas)
		if maxFeeMinusBaseFee.Cmp(priorityFeePerGas) < 0 {
			priorityFeePerGas = maxFeeMinusBaseFee
		}

		effectiveGasPrice := big.NewInt(0).Add(priorityFeePerGas, block.BaseFeePerGas)
		return effectiveGasPrice
	default:
		b.logger.Warnw(fmt.Sprintf("Ignoring unknown transaction type %v", tx.Type), "block", block, "tx", tx)
		return nil
	}
}

func (b *BlockHistoryEstimator) EffectiveTipCap(block Block, tx Transaction) *big.Int {
	switch tx.Type {
	case 0x2:
		return tx.MaxPriorityFeePerGas
	case 0x0, 0x1:
		if tx.GasPrice == nil {
			return nil
		}
		if block.BaseFeePerGas == nil {
			return nil
		}
		effectiveTipCap := big.NewInt(0).Sub(tx.GasPrice, block.BaseFeePerGas)
		if effectiveTipCap.Cmp(big.NewInt(0)) < 0 {
			b.logger.AssumptionViolationw("GasPrice - BaseFeePerGas >= 0", "block", block, "tx", tx)
			return nil
		}
		return effectiveTipCap
	default:
		b.logger.Warnw(fmt.Sprintf("Ignoring unknown transaction type %v", tx.Type), "block", block, "tx", tx)
		return nil
	}
}
