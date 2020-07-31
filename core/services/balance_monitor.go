package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// BalanceMonitor checks the balance for each key on every new head
type BalanceMonitor interface {
	store.HeadTrackable
	GetEthBalance(gethCommon.Address) *assets.Eth
}

type balanceMonitor struct {
	store             *store.Store
	gethClientWrapper store.GethClientWrapper
	ethBalances       map[gethCommon.Address]*assets.Eth
	ethBalancesMtx    *sync.RWMutex
}

// NewBalanceMonitor returns a new balanceMonitor
func NewBalanceMonitor(store *store.Store, gethClientWrapper store.GethClientWrapper) BalanceMonitor {
	return &balanceMonitor{
		store:             store,
		gethClientWrapper: gethClientWrapper,
		ethBalances:       make(map[gethCommon.Address]*assets.Eth),
		ethBalancesMtx:    new(sync.RWMutex),
	}
}

// Connect complies with HeadTrackable
func (bm *balanceMonitor) Connect(_ *models.Head) error {
	// Connect head can be out of date, so always query the latest balance
	bm.checkBalance(nil)
	return nil
}

// Disconnect complies with HeadTrackable
func (bm *balanceMonitor) Disconnect() {}

const ethFetchTimeout = 2 * time.Second

// OnNewLongestChain checks the balance for each key
func (bm *balanceMonitor) OnNewLongestChain(head models.Head) {
	bm.checkBalance(&head)
}

func (bm *balanceMonitor) checkBalance(head *models.Head) {
	keys, err := bm.store.Keys()
	if err != nil {
		logger.Error("BalanceMonitor: error getting keys", err)
	}

	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)

		go func(k models.Key) {
			ctx, cancel := context.WithTimeout(context.Background(), ethFetchTimeout)
			defer cancel()

			var headNum *big.Int

			err := bm.gethClientWrapper.GethClient(func(c eth.GethClient) error {
				var err error
				var bal *big.Int
				if head == nil {
					headNum = nil
					bal, err = c.BalanceAt(ctx, k.Address.Address(), nil)
				} else {
					headNum = big.NewInt(head.Number)
					bal, err = c.BalanceAt(ctx, k.Address.Address(), headNum)
				}
				if err != nil {
					return errors.Wrap(err, "error getting balance")
				}
				if bal == nil {
					return errors.New("BalanceMonitor: invariant violation, bal may not be nil")
				}
				ethBal := assets.Eth(*bal)

				bm.updateBalance(ethBal, k.Address.Address(), headNum)
				return nil
			})

			if err != nil {
				logger.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s: %s", k.Address.Hex(), err.Error()), "err", err, "address", k.Address, "headNum", headNum)
			}

			wg.Done()
		}(key)
	}

	wg.Wait()
}

func (bm *balanceMonitor) updateBalance(ethBal assets.Eth, address gethCommon.Address, headNum *big.Int) {
	store.PromUpdateEthBalance(&ethBal, address)

	bm.ethBalancesMtx.Lock()
	oldBal := bm.ethBalances[address]
	bm.ethBalances[address] = &ethBal
	bm.ethBalancesMtx.Unlock()

	loggerFields := []interface{}{
		"address", address.Hex(),
		"headNum", headNum,
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
