package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

// BalanceMonitor checks the balance for each key on every new head
type BalanceMonitor interface {
	store.HeadTrackable
	GetEthBalance(gethCommon.Address) *assets.Eth
	Stop()
}

type balanceMonitor struct {
	store          *store.Store
	ethBalances    map[gethCommon.Address]*assets.Eth
	ethBalancesMtx *sync.RWMutex
	sleeperTask    SleeperTask
}

// NewBalanceMonitor returns a new balanceMonitor
func NewBalanceMonitor(store *store.Store) BalanceMonitor {
	bm := &balanceMonitor{
		store:          store,
		ethBalances:    make(map[gethCommon.Address]*assets.Eth),
		ethBalancesMtx: new(sync.RWMutex),
	}
	bm.sleeperTask = NewSleeperTask(&worker{bm: bm})
	return bm
}

// Connect complies with HeadTrackable
func (bm *balanceMonitor) Connect(_ *models.Head) error {
	// Connect head can be out of date, so always query the latest balance
	bm.checkBalance(nil)
	return nil
}

// Stop shuts down the BalanceMonitor, should not be used after this
func (bm *balanceMonitor) Stop() {
	bm.sleeperTask.Stop()
}

// Disconnect complies with HeadTrackable
func (bm *balanceMonitor) Disconnect() {}

// OnNewLongestChain checks the balance for each key
func (bm *balanceMonitor) OnNewLongestChain(_ context.Context, head models.Head) {
	bm.checkBalance(&head)
}

func (bm *balanceMonitor) checkBalance(head *models.Head) {
	logger.Debugw("BalanceMonitor: signalling balance worker")
	bm.sleeperTask.WakeUp()
}

func (bm *balanceMonitor) updateBalance(ethBal assets.Eth, address gethCommon.Address) {
	store.PromUpdateEthBalance(&ethBal, address)

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
		logger.Infow(fmt.Sprintf("ETH balance for %s: %s", address.Hex(), ethBal.String()), loggerFields...)
		return
	}

	if ethBal.Cmp(oldBal) != 0 {
		logger.Infow(fmt.Sprintf("New ETH balance for %s: %s", address.Hex(), ethBal.String()), loggerFields...)
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
	keys, err := w.bm.store.Keys()
	if err != nil {
		logger.Error("BalanceMonitor: error getting keys", err)
	}

	var wg sync.WaitGroup

	wg.Add(len(keys))
	for _, key := range keys {
		go func(k models.Key) {
			w.checkAccountBalance(k)
			wg.Done()
		}(key)
	}
	wg.Wait()
}

// Approximately ETH block time
const ethFetchTimeout = 15 * time.Second

func (w *worker) checkAccountBalance(k models.Key) {
	ctx, cancel := context.WithTimeout(context.Background(), ethFetchTimeout)
	defer cancel()

	bal, err := w.bm.store.EthClient.BalanceAt(ctx, k.Address.Address(), nil)
	if err != nil {
		logger.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s", k.Address.Hex()),
			"error", err,
			"address", k.Address,
		)
	} else if bal == nil {
		logger.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s: invariant violation, bal may not be nil", k.Address.Hex()),
			"error", err,
			"address", k.Address,
		)
	} else {
		ethBal := assets.Eth(*bal)
		w.bm.updateBalance(ethBal, k.Address.Address())
	}
}
