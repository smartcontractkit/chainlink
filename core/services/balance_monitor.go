package services

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

type (
	// BalanceMonitor checks the balance for each key on every new head
	BalanceMonitor interface {
		httypes.HeadTrackable
		GetEthBalance(gethCommon.Address) *assets.Eth
		service.Service
	}

	balanceMonitor struct {
		logger         *logger.Logger
		db             *gorm.DB
		ethClient      eth.Client
		ethKeyStore    keystore.Eth
		ethBalances    map[gethCommon.Address]*assets.Eth
		ethBalancesMtx *sync.RWMutex
		sleeperTask    utils.SleeperTask
	}

	NullBalanceMonitor struct{}
)

// NewBalanceMonitor returns a new balanceMonitor
func NewBalanceMonitor(db *gorm.DB, ethClient eth.Client, ethKeyStore keystore.Eth, logger *logger.Logger) BalanceMonitor {
	bm := &balanceMonitor{
		logger,
		db,
		ethClient,
		ethKeyStore,
		make(map[gethCommon.Address]*assets.Eth),
		new(sync.RWMutex),
		nil,
	}
	bm.sleeperTask = utils.NewSleeperTask(&worker{bm: bm})
	return bm
}

func (bm *balanceMonitor) Start() error {
	// Always query latest balance on start
	(&worker{bm}).Work()
	return nil
}

// Close shuts down the BalanceMonitor, should not be used after this
func (bm *balanceMonitor) Close() error {
	return bm.sleeperTask.Stop()
}

func (bm *balanceMonitor) Ready() error {
	return nil
}

func (bm *balanceMonitor) Healthy() error {
	return nil
}

// OnNewLongestChain checks the balance for each key
func (bm *balanceMonitor) OnNewLongestChain(_ context.Context, head models.Head) {
	bm.checkBalance(&head)
}

func (bm *balanceMonitor) checkBalance(head *models.Head) {
	bm.logger.Debugw("BalanceMonitor: signalling balance worker")
	bm.sleeperTask.WakeUp()
}

func (bm *balanceMonitor) updateBalance(ethBal assets.Eth, address gethCommon.Address) {
	bm.promUpdateEthBalance(&ethBal, address)

	bm.ethBalancesMtx.Lock()
	oldBal := bm.ethBalances[address]
	bm.ethBalances[address] = &ethBal
	bm.ethBalancesMtx.Unlock()

	loggerFields := []interface{}{
		"address", address.Hex(),
		"ethBalance", ethBal.String(),
		"weiBalance", ethBal.ToInt(),
		"id", "balance_log",
	}

	if oldBal == nil {
		bm.logger.Infow(fmt.Sprintf("ETH balance for %s: %s", address.Hex(), ethBal.String()), loggerFields...)
		return
	}

	if ethBal.Cmp(oldBal) != 0 {
		bm.logger.Infow(fmt.Sprintf("New ETH balance for %s: %s", address.Hex(), ethBal.String()), loggerFields...)
	}
}

func (bm *balanceMonitor) GetEthBalance(address gethCommon.Address) *assets.Eth {
	bm.ethBalancesMtx.RLock()
	defer bm.ethBalancesMtx.RUnlock()
	return bm.ethBalances[address]
}

type worker struct {
	bm *balanceMonitor
}

func (w *worker) Work() {
	keys, err := w.bm.ethKeyStore.SendingKeys()
	if err != nil {
		w.bm.logger.Error("BalanceMonitor: error getting keys", err)
	}

	var wg sync.WaitGroup

	wg.Add(len(keys))
	for _, key := range keys {
		go func(k ethkey.KeyV2) {
			defer wg.Done()
			w.checkAccountBalance(k)
		}(key)
	}
	wg.Wait()
}

// Approximately ETH block time
const ethFetchTimeout = 15 * time.Second

func (w *worker) checkAccountBalance(k ethkey.KeyV2) {
	ctx, cancel := context.WithTimeout(context.Background(), ethFetchTimeout)
	defer cancel()

	bal, err := w.bm.ethClient.BalanceAt(ctx, k.Address.Address(), nil)
	if err != nil {
		w.bm.logger.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s", k.Address.Hex()),
			"error", err,
			"address", k.Address,
		)
	} else if bal == nil {
		w.bm.logger.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s: invariant violation, bal may not be nil", k.Address.Hex()),
			"error", err,
			"address", k.Address,
		)
	} else {
		ethBal := assets.Eth(*bal)
		w.bm.updateBalance(ethBal, k.Address.Address())
	}
}

func (*NullBalanceMonitor) GetEthBalance(gethCommon.Address) *assets.Eth {
	return nil
}
func (*NullBalanceMonitor) Start() error                                            { return nil }
func (*NullBalanceMonitor) Close() error                                            { return nil }
func (*NullBalanceMonitor) Ready() error                                            { return nil }
func (*NullBalanceMonitor) Healthy() error                                          { return nil }
func (*NullBalanceMonitor) OnNewLongestChain(ctx context.Context, head models.Head) {}

var promETHBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "eth_balance",
		Help: "Each Ethereum account's balance",
	},
	[]string{"account"},
)

func (bm *balanceMonitor) promUpdateEthBalance(balance *assets.Eth, from gethCommon.Address) {
	balanceFloat, err := ApproximateFloat64(balance)

	if err != nil {
		bm.logger.Error(fmt.Errorf("updatePrometheusEthBalance: %v", err))
		return
	}

	promETHBalance.WithLabelValues(from.Hex()).Set(balanceFloat)
}

func ApproximateFloat64(e *assets.Eth) (float64, error) {
	ef := new(big.Float).SetInt(e.ToInt())
	weif := new(big.Float).SetInt(models.WeiPerEth)
	bf := new(big.Float).Quo(ef, weif)
	f64, _ := bf.Float64()
	if f64 == math.Inf(1) || f64 == math.Inf(-1) {
		return math.Inf(1), errors.New("assets.Eth.Float64: Could not approximate Eth value into float")
	}
	return f64, nil
}
