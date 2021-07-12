package log

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	ethSubscriber struct {
		ethClient eth.Client
		config    Config
		chStop    chan struct{}
	}
)

func newEthSubscriber(ethClient eth.Client, config Config, chStop chan struct{}) *ethSubscriber {
	return &ethSubscriber{
		ethClient: ethClient,
		config:    config,
		chStop:    chStop,
	}
}

func (sub *ethSubscriber) backfillLogs(fromBlockOverride null.Int64, addresses []common.Address, topics []common.Hash) (chBackfilledLogs chan types.Log, abort bool) {
	if len(addresses) == 0 {
		logger.Debug("LogBroadcaster: No addresses to backfill for, returning")
		ch := make(chan types.Log)
		close(ch)
		return ch, false
	}

	ctx, cancel := utils.ContextFromChan(sub.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		latestBlock, err := sub.ethClient.HeadByNumber(ctx, nil)
		if err != nil {
			logger.Errorw("LogBroadcaster: backfill - could not fetch latest block header, will retry", "err", err)
			return true
		} else if latestBlock == nil {
			logger.Warn("LogBroadcaster: got nil block header, will retry")
			return true
		}
		currentHeight := uint64(latestBlock.Number)

		// Backfill from `backfillDepth` blocks ago.  It's up to the subscribers to
		// filter out logs they've already dealt with.
		fromBlock := currentHeight - sub.config.BlockBackfillDepth()
		if fromBlock > currentHeight {
			fromBlock = 0 // Overflow protection
		}

		if fromBlockOverride.Valid {
			logger.Infow("LogBroadcaster: Using the override a limit of backfill", "blockNumber", fromBlockOverride.Int64)
			fromBlock = uint64(fromBlockOverride.Int64)
		}

		logger.Infow("LogBroadcaster: Backfilling logs from", "blockNumber", fromBlock)

		q := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}

		logs := make([]types.Log, 0)
		start := time.Now()
		// If we are significantly behind the latest head, there could be a very large (1000s)
		// of blocks to check for logs. We read the blocks in batches to avoid hitting the websocket
		// request data limit.
		// On matic its 5MB [https://github.com/maticnetwork/bor/blob/3de2110886522ab17e0b45f3c4a6722da72b7519/rpc/http.go#L35]
		// On ethereum its 15MB [https://github.com/ethereum/go-ethereum/blob/master/rpc/websocket.go#L40]
		batchSize := int64(sub.config.EthLogBackfillBatchSize())
		for i := q.FromBlock.Int64(); i <= int64(currentHeight); i += batchSize {

			untilIncluded := i + batchSize - 1
			if untilIncluded > int64(currentHeight) {
				untilIncluded = int64(currentHeight)
			}
			q.FromBlock = big.NewInt(i)
			q.ToBlock = big.NewInt(untilIncluded)
			batchLogs, err := sub.ethClient.FilterLogs(ctx, q)
			if err != nil {
				if ctx.Err() != nil {
					logger.Errorw("LogBroadcaster: Deadline exceeded, unable to backfill a batch of logs", "err", err, "elapsed", time.Since(start), "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
				} else {
					logger.Errorw("LogBroadcaster: Unable to backfill a batch of logs", "err", err, "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
				}
				return true
			}

			select {
			case <-sub.chStop:
				return
			default:
				logs = append(logs, batchLogs...)
			}
		}

		logger.Infof("LogBroadcaster: Finished getting %v logs for backfill", len(logs))

		// unbufferred channel, as it will be filled in the goroutine,
		// while the broadcaster's eventLoop is reading from it
		chBackfilledLogs = make(chan types.Log)
		go func() {
			defer close(chBackfilledLogs)
			for _, log := range logs {
				select {
				case chBackfilledLogs <- log:
				case <-sub.chStop:
					return
				}
			}
			logger.Infof("LogBroadcaster: Finished async backfill of %v logs", len(logs))
		}()

		return false
	})
	select {
	case <-sub.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

// createSubscription creates a new log subscription starting at the current block.  If previous logs
// are needed, they must be obtained through backfilling, as subscriptions can only be started from
// the current head.
func (sub *ethSubscriber) createSubscription(addresses []common.Address, topics []common.Hash) (subscr managedSubscription, abort bool) {
	if len(addresses) == 0 {
		return newNoopSubscription(), false
	}

	ctx, cancel := utils.ContextFromChan(sub.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {

		filterQuery := ethereum.FilterQuery{
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}
		chRawLogs := make(chan types.Log)

		ctx2, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		innerSub, err := sub.ethClient.SubscribeFilterLogs(ctx2, filterQuery, chRawLogs)
		if err != nil {
			logger.Errorw("Log subscriber could not create subscription to Ethereum node", "err", err)
			return true
		}

		subscr = managedSubscriptionImpl{
			subscription: innerSub,
			chRawLogs:    chRawLogs,
		}
		return false
	})
	select {
	case <-sub.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

// A managedSubscription acts as wrapper for the Subscription. Specifically, the
// managedSubscription closes the log channel as soon as the unsubscribe request is made
type managedSubscription interface {
	Err() <-chan error
	Logs() chan types.Log
	Unsubscribe()
}

type managedSubscriptionImpl struct {
	subscription ethereum.Subscription
	chRawLogs    chan types.Log
}

func (sub managedSubscriptionImpl) Err() <-chan error {
	return sub.subscription.Err()
}

func (sub managedSubscriptionImpl) Logs() chan types.Log {
	return sub.chRawLogs
}

func (sub managedSubscriptionImpl) Unsubscribe() {
	sub.subscription.Unsubscribe()
	close(sub.chRawLogs)
}

type noopSubscription struct {
	chRawLogs chan types.Log
}

func newNoopSubscription() noopSubscription {
	return noopSubscription{make(chan types.Log)}
}

func (b noopSubscription) Err() <-chan error    { return nil }
func (b noopSubscription) Logs() chan types.Log { return b.chRawLogs }
func (b noopSubscription) Unsubscribe()         { close(b.chRawLogs) }
