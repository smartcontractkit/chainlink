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
}

type balanceMonitor struct {
	store          *store.Store
	ethBalances    map[gethCommon.Address]*assets.Eth
	ethBalancesMtx *sync.RWMutex
	headQueue      chan *models.Head
}

// NewBalanceMonitor returns a new balanceMonitor
func NewBalanceMonitor(store *store.Store) BalanceMonitor {
	bm := &balanceMonitor{
		store:          store,
		ethBalances:    make(map[gethCommon.Address]*assets.Eth),
		ethBalancesMtx: new(sync.RWMutex),
		headQueue:      make(chan *models.Head, 1),
	}
	go balanceWorker(bm)
	return bm
}

// Connect complies with HeadTrackable
func (bm *balanceMonitor) Connect(_ *models.Head) error {
	// Connect head can be out of date, so always query the latest balance
	bm.checkBalance(nil)
	return nil
}

// Disconnect complies with HeadTrackable
func (bm *balanceMonitor) Disconnect() {
	close(bm.headQueue)
}

// OnNewLongestChain checks the balance for each key
func (bm *balanceMonitor) OnNewLongestChain(_ context.Context, head models.Head) {
	bm.checkBalance(&head)
}

func (bm *balanceMonitor) checkBalance(head *models.Head) {
	// Non blocking send on headQueue, if the queue is full, hits the default
	// block and returns immediately
	select {
	case bm.headQueue <- head:
	default:
	}
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

func balanceWorker(bm *balanceMonitor) {
	for {
		select {
		case head, open := <-bm.headQueue:
			if !open {
				return
			}

			keys, err := bm.store.Keys()
			if err != nil {
				logger.Error("BalanceMonitor: error getting keys", err)
			}

			var wg sync.WaitGroup

			wg.Add(len(keys))
			for _, key := range keys {
				go func(k models.Key) {
					checkAccountBalance(bm, head, k)
					wg.Done()
				}(key)
			}
			wg.Wait()
		}
	}
}

const ethFetchTimeout = 2 * time.Second

func checkAccountBalance(bm *balanceMonitor, head *models.Head, k models.Key) {
	ctx, cancel := context.WithTimeout(context.TODO(), ethFetchTimeout)
	defer cancel()

	bal, err := bm.store.EthClient.BalanceAt(ctx, k.Address.Address(), nil)
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
		bm.updateBalance(ethBal, k.Address.Address())
	}
}
