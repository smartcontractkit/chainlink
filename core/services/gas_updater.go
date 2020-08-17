package services

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
		[]string{"percentile", "block_num"},
	)
)

// GasUpdater listens for new heads and updates the base gas price dynamically
// based on the configured percentile of gas prices in that block
type GasUpdater interface {
	store.HeadTrackable
	RollingBlockHistory() []*types.Block
}

type gasUpdater struct {
	store                   *store.Store
	rollingBlockHistory     []*types.Block
	rollingBlockHistorySize int
	// HACK: blockDelay is the number of blocks that the gas updater trails behind head.
	// E.g. if this is set to 3, and we receive block 10, gas updater will
	// fetch block 7.
	// This is necessary because geth/parity send heads as soon as they get
	// them and often the actual block is not available until later. Fetching
	// it too early results in an empty block.
	blockDelay int64
	percentile int
}

// NewGasUpdater returns a new gas updater.
func NewGasUpdater(store *store.Store) GasUpdater {
	gu := &gasUpdater{
		store:                   store,
		rollingBlockHistory:     make([]*types.Block, 0),
		rollingBlockHistorySize: int(store.Config.GasUpdaterBlockHistorySize()),
		blockDelay:              int64(store.Config.GasUpdaterBlockDelay()),
		percentile:              int(store.Config.GasUpdaterTransactionPercentile()),
	}
	return gu
}

func (gu *gasUpdater) Connect(bn *models.Head) error {
	if gu.store.Config.GasUpdaterEnabled() {
		logger.Debugw("GasUpdater: dynamic gas updates are enabled", "ethGasPriceDefault", gu.store.Config.EthGasPriceDefault())
	} else {
		logger.Debugw("GasUpdater: dynamic gas updating is disabled", "ethGasPriceDefault", gu.store.Config.EthGasPriceDefault())
	}
	return nil
}

func (gu *gasUpdater) Disconnect() {
}

// OnNewLongestChain recalculates and sets global gas price on every head
func (gu *gasUpdater) OnNewLongestChain(head models.Head) {
	// Bail out as early as possible if the gas updater is disabled so we avoid
	// any potential undesired side effects. Note that in a future iteration
	// the GasUpdaterEnabled setting could be modifiable at runtime
	if !gu.store.Config.GasUpdaterEnabled() {
		return
	}
	blockToFetch := head.Number - gu.blockDelay
	if blockToFetch < 0 {
		logger.Warnf("GasUpdater: skipping gas calculation, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY of %v", head.Number, gu.blockDelay)
		return
	}
	block, err := gu.store.EthClient.BlockByNumber(context.TODO(), big.NewInt(blockToFetch))
	if err != nil {
		logger.Error(err, fmt.Sprintf("GasUpdater: error retrieving block %v", blockToFetch))
		return
	}
	if len(block.Transactions()) > 0 {
		gu.rollingBlockHistory = append(gu.rollingBlockHistory, block)
		if len(gu.rollingBlockHistory) > gu.rollingBlockHistorySize {
			gu.rollingBlockHistory = gu.rollingBlockHistory[1:]
			percentileGasPrice := gu.percentileGasPrice()
			err := gu.setPercentileGasPrice(percentileGasPrice)
			if err != nil {
				logger.Error("GasUpdater error setting gas price: ", err)
				return
			}
			promGasUpdaterSetGasPrice.WithLabelValues(fmt.Sprintf("%v%%", gu.percentile), string(blockToFetch)).Set(float64(percentileGasPrice))
		} else {
			logger.Debugw(fmt.Sprintf("GasUpdater: waiting for blocks: %v/%v", len(gu.rollingBlockHistory), gu.rollingBlockHistorySize), "inHistory", len(gu.rollingBlockHistory), "required", gu.rollingBlockHistorySize)
		}
	} else {
		logger.Debugw(fmt.Sprintf("GasUpdater: skipping empty block: %v", blockToFetch), "blockNumber", blockToFetch)
	}
}

func (gu *gasUpdater) percentileGasPrice() int64 {
	gasPrices := make([]int64, 0)
	for _, block := range gu.rollingBlockHistory {
		for _, tx := range block.Transactions() {
			gasPrices = append(gasPrices, tx.GasPrice().Int64())
		}
	}
	sort.Slice(gasPrices, func(i, j int) bool { return gasPrices[i] < gasPrices[j] })
	idx := ((len(gasPrices) - 1) * gu.percentile) / 100
	for i := 0; i <= 100; i += 5 {
		jdx := ((len(gasPrices) - 1) * i) / 100
		promGasUpdaterAllPercentiles.WithLabelValues(fmt.Sprintf("%v%%", i)).Set(float64(gasPrices[jdx]))
	}
	return gasPrices[idx]
}

func (gu *gasUpdater) setPercentileGasPrice(gasPrice int64) error {
	gasPriceGwei := fmt.Sprintf("%.2f", float64(gasPrice)/1000000000)
	bigGasPrice := big.NewInt(gasPrice)
	if bigGasPrice.Cmp(gu.store.Config.EthMaxGasPriceWei()) > 0 {
		return fmt.Errorf("cannot set gas price %s because it exceeds EthMaxGasPriceWei %s", bigGasPrice.String(), gu.store.Config.EthMaxGasPriceWei().String())
	}
	logger.Debugw(fmt.Sprintf("GasUpdater: setting new default gas price: %v Gwei", gasPriceGwei), "gasPriceWei", gasPrice, "gasPriceGWei", gasPriceGwei)
	return gu.store.Config.SetEthGasPriceDefault(bigGasPrice)
}

func (gu *gasUpdater) RollingBlockHistory() []*types.Block {
	return gu.rollingBlockHistory
}
