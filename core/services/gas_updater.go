package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// maxStartTime is the maximum amount of time we are allowed to spend
	// trying to fill initial data on start. This must be capped because it can
	// block the application from starting.
	maxStartTime = 10 * time.Second
	// maxEthNodeRequestTime is the worst case time we will wait for a response
	// from the eth node before we consider it to be an error
	maxEthNodeRequestTime = 30 * time.Second
)

var (
	promGasUpdaterAllPercentiles = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_all_gas_percetiles",
		Help: "Gas price at given percentile",
	},
		[]string{"percentile"},
	)

	promGasUpdaterSetGasPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_updater_set_gas_price",
		Help: "Gas updater set gas price (in Wei)",
	},
		[]string{"percentile"},
	)
)

// GasUpdater listens for new heads and updates the base gas price dynamically
// based on the configured percentile of gas prices in that block
type GasUpdater interface {
	store.HeadTrackable
	Start() error
	Close() error
}

type gasUpdater struct {
	utils.StartStopOnce
	store                   *store.Store
	rollingBlockHistory     []Block
	rollingBlockHistorySize int64
	// HACK: blockDelay is the number of blocks that the gas updater trails behind head.
	// E.g. if this is set to 3, and we receive block 10, gas updater will
	// fetch block 7.
	// This is necessary because geth/parity send heads as soon as they get
	// them and often the actual block is not available until later. Fetching
	// it too early results in an empty block.
	blockDelay int64
	percentile int
	mb         *utils.Mailbox
	wg         *sync.WaitGroup
	ctx        context.Context
	ctxCancel  context.CancelFunc

	logger *logger.Logger
}

type blockInternal struct {
	Number       string
	Hash         common.Hash
	ParentHash   common.Hash
	Transactions []types.Transaction
}

func toBlockNumStr(n int64) string {
	return hexutil.EncodeBig(big.NewInt(n))
}

// Block represents an ethereum block
type Block struct {
	Number       int64
	Hash         common.Hash
	ParentHash   common.Hash
	Transactions []types.Transaction
}

// MarshalJSON implements json marshalling for Block
func (b Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(blockInternal{
		toBlockNumStr(b.Number),
		b.Hash,
		b.ParentHash,
		b.Transactions,
	})
}

// UnmarshalJSON unmarshals to a Block
func (b *Block) UnmarshalJSON(data []byte) error {
	bi := blockInternal{}
	if err := json.Unmarshal(data, &bi); err != nil {
		return err
	}
	n, err := hexutil.DecodeBig(bi.Number)
	if err != nil {
		return err
	}
	*b = Block{
		n.Int64(),
		bi.Hash,
		bi.ParentHash,
		bi.Transactions,
	}
	return nil
}

// NewGasUpdater returns a new gas updater.
func NewGasUpdater(store *store.Store) GasUpdater {
	ctx, cancel := context.WithCancel(context.Background())
	gu := &gasUpdater{
		StartStopOnce:           utils.StartStopOnce{},
		store:                   store,
		rollingBlockHistory:     make([]Block, 0),
		rollingBlockHistorySize: int64(store.Config.GasUpdaterBlockHistorySize()),
		blockDelay:              int64(store.Config.GasUpdaterBlockDelay()),
		percentile:              int(store.Config.GasUpdaterTransactionPercentile()),
		mb:                      utils.NewMailbox(1),
		wg:                      new(sync.WaitGroup),
		ctx:                     ctx,
		ctxCancel:               cancel,
		logger:                  logger.CreateLogger(logger.Default.With("id", "gas_updater")),
	}

	return gu
}

func (gu *gasUpdater) Connect(bn *models.Head) error {
	return nil
}

func (gu *gasUpdater) Disconnect() {
}

// OnNewLongestChain recalculates and sets global gas price if a new head comes
// in and we are not currently fetching
func (gu *gasUpdater) OnNewLongestChain(ctx context.Context, head models.Head) {
	// Bail out as early as possible if the gas updater is disabled so we avoid
	// any potential undesired side effects. Note that in a future iteration
	// the GasUpdaterEnabled setting could be modifiable at runtime
	if gu.State() != utils.StartStopOnce_Started {
		panic("OnNewLongestChain called on unstarted gas updater")
	}
	gu.mb.Deliver(head)
}

func (gu *gasUpdater) Start() error {
	if !gu.OkayToStart() {
		return errors.New("GasUpdater has already been started")
	}
	gu.logger.Debugw("GasUpdater: starting")
	if gu.rollingBlockHistorySize > int64(gu.store.Config.EthFinalityDepth()) {
		gu.logger.Warnf("GasUpdater: GAS_UPDATER_BLOCK_HISTORY_SIZE=%v is greater than ETH_FINALITY_DEPTH=%v, blocks deeper than finality depth will be refetched on every gas updater cycle, causing unnecessary load on the eth node. Consider decreasing GAS_UPDATER_BLOCK_HISTORY_SIZE or increasing ETH_FINALITY_DEPTH", gu.rollingBlockHistorySize, gu.store.Config.EthFinalityDepth())
	}
	ctx, cancel := context.WithTimeout(gu.ctx, maxStartTime)
	defer cancel()
	latestHead, err := gu.store.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		logger.Warnw("GasUpdater: initial check for latest head failed", "err", err)
	} else {
		gu.logger.Debugw("GasUpdater: got latest head", "number", latestHead.Number, "blockHash", latestHead.Hash.Hex())
		gu.FetchBlocksAndRecalculate(ctx, *latestHead)
	}
	gu.wg.Add(1)
	go gu.runLoop()
	gu.logger.Debugw("GasUpdater: started")
	return nil
}

func (gu *gasUpdater) Close() error {
	if !gu.OkayToStop() {
		return errors.New("GasUpdater has already been stopped")
	}
	gu.ctxCancel()
	gu.wg.Wait()
	return nil
}

func (gu *gasUpdater) runLoop() {
	defer gu.wg.Done()
	for {
		select {
		case <-gu.ctx.Done():
			return
		case <-gu.mb.Notify():
			head := gu.mb.Retrieve()
			h, is := head.(models.Head)
			if !is {
				panic(fmt.Sprintf("invariant violation, expected %T but got %T", models.Head{}, head))
			}
			gu.FetchBlocksAndRecalculate(gu.ctx, h)
		}
	}
}

func (gu *gasUpdater) FetchBlocksAndRecalculate(ctx context.Context, head models.Head) {
	ctx, cancel := context.WithTimeout(ctx, maxEthNodeRequestTime)
	defer cancel()

	if err := gu.FetchBlocks(ctx, head); err != nil {
		gu.logger.Warnw("GasUpdater: error fetching blocks", "head", head, "err", err)
		return
	}

	gu.Recalculate(head)
}

// FetchHeadsAndRecalculate adds the given heads to the history and recalculates gas price
func (gu *gasUpdater) Recalculate(head models.Head) {
	if len(gu.rollingBlockHistory) == 0 {
		return
	}

	percentileGasPrice, err := gu.percentileGasPrice()
	if err != nil {
		logger.Errorw("GasUpdater: error getting percentile gas price", "err", err)
		return
	}
	gasPriceGwei := fmt.Sprintf("%.2f", float64(percentileGasPrice)/1000000000)

	var numsInHistory []int64
	for _, b := range gu.rollingBlockHistory {
		numsInHistory = append(numsInHistory, b.Number)
	}
	gu.logger.Debugw(fmt.Sprintf("GasUpdater: setting new default gas price: %v Gwei", gasPriceGwei),
		"gasPriceWei", percentileGasPrice,
		"gasPriceGWei", gasPriceGwei,
		"headNum", head.Number,
		"blocks", numsInHistory,
	)
	if err := gu.setPercentileGasPrice(percentileGasPrice); err != nil {
		gu.logger.Errorw("GasUpdater: error setting gas price", "err", err)
		return
	}
	promGasUpdaterSetGasPrice.WithLabelValues(fmt.Sprintf("%v%%", gu.percentile)).Set(float64(percentileGasPrice))
}

func (gu *gasUpdater) FetchBlocks(ctx context.Context, head models.Head) error {
	highestBlockToFetch := head.Number - gu.blockDelay
	if highestBlockToFetch < 0 {
		return errors.Errorf("GasUpdater: cannot fetch, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY of %v", head.Number, gu.blockDelay)
	}
	lowestBlockToFetch := head.Number - gu.rollingBlockHistorySize - gu.blockDelay + 1
	if lowestBlockToFetch < 0 {
		lowestBlockToFetch = 0
	}

	blocks := make(map[int64]Block)
	for _, block := range gu.rollingBlockHistory {
		// Make a best-effort to be re-org resistant using the head
		// chain, refetch blocks that got re-org'd out.
		// NOTE: Any blocks older than the oldest block in the provided chain
		// will be also be refetched.
		if head.IsInChain(block.Hash) {
			blocks[block.Number] = block
		}
	}

	var reqs []rpc.BatchElem
	for i := lowestBlockToFetch; i <= highestBlockToFetch; i++ {
		// NOTE: To save rpc calls, don't fetch blocks we already have in the history
		if _, exists := blocks[i]; exists {
			continue
		}

		req := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{toBlockNumStr(i), true},
			Result: &Block{}, // FIXME: Can this be a *((*Block)(nil)) instead?
		}
		reqs = append(reqs, req)
	}

	gu.logger.Debugf("GasUpdater: fetching %v blocks (%v in local history)", len(reqs), len(blocks))
	err := gu.store.EthClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return errors.Wrap(err, "GasUpdater#fetchBlocks error fetching blocks with BatchCallContext")
	}

	for _, req := range reqs {
		result, err := req.Result, req.Error
		if err != nil {
			gu.logger.Warnw("GasUpdater#fetchBlocks ignoring errored block", "err", err, "blockNum", head.Number)
			continue
		}

		b, is := result.(*Block)
		if !is {
			return errors.Errorf("expected result to be a %T, got %T", &Block{}, result)
		}
		if b == nil {
			return errors.New("invariant violation: got nil block")
		}
		if b.Hash == (common.Hash{}) {
			gu.logger.Warnw("GasUpdater#fetchBlocks block was missing hash", "block", b, "blockNum", head.Number, "erroredBlockNum", b.Number)
			continue
		}

		blocks[b.Number] = *b
	}

	newBlockHistory := make([]Block, 0)
	for _, block := range blocks {
		newBlockHistory = append(newBlockHistory, block)
	}
	sort.Slice(newBlockHistory, func(i, j int) bool {
		return newBlockHistory[i].Number < newBlockHistory[j].Number
	})

	start := len(newBlockHistory) - int(gu.rollingBlockHistorySize)
	if start < 0 {
		gu.logger.Infow(fmt.Sprintf("GasUpdater: using fewer blocks than the specified history size: %v/%v", len(newBlockHistory), gu.rollingBlockHistorySize), "rollingBlockHistorySize", gu.rollingBlockHistorySize, "headNum", head.Number, "blocksAvailable", len(newBlockHistory))
		start = 0
	}

	gu.rollingBlockHistory = newBlockHistory[start:]

	return nil
}

func (gu *gasUpdater) percentileGasPrice() (int64, error) {
	gasPrices := make([]int64, 0)
	for _, block := range gu.rollingBlockHistory {
		for _, tx := range block.Transactions {
			gasPrices = append(gasPrices, tx.GasPrice().Int64())
		}
	}
	if len(gasPrices) == 0 {
		return 0, errors.New("no transactions")
	}
	sort.Slice(gasPrices, func(i, j int) bool { return gasPrices[i] < gasPrices[j] })
	idx := ((len(gasPrices) - 1) * gu.percentile) / 100
	for i := 0; i <= 100; i += 5 {
		jdx := ((len(gasPrices) - 1) * i) / 100
		promGasUpdaterAllPercentiles.WithLabelValues(fmt.Sprintf("%v%%", i)).Set(float64(gasPrices[jdx]))
	}
	return gasPrices[idx], nil
}

func (gu *gasUpdater) setPercentileGasPrice(gasPrice int64) error {
	bigGasPrice := big.NewInt(gasPrice)
	if bigGasPrice.Cmp(gu.store.Config.EthMaxGasPriceWei()) > 0 {
		return fmt.Errorf("cannot set gas price %s because it exceeds EthMaxGasPriceWei %s", bigGasPrice.String(), gu.store.Config.EthMaxGasPriceWei().String())
	}
	return gu.store.Config.SetEthGasPriceDefault(bigGasPrice)
}

func (gu *gasUpdater) RollingBlockHistory() []Block {
	return gu.rollingBlockHistory
}
