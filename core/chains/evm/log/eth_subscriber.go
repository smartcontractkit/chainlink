package log

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

type (
	ethSubscriber struct {
		ethClient evmclient.Client
		config    Config
		logger    logger.Logger
		chStop    services.StopChan
	}
)

func newEthSubscriber(ethClient evmclient.Client, config Config, lggr logger.Logger, chStop chan struct{}) *ethSubscriber {
	return &ethSubscriber{
		ethClient: ethClient,
		config:    config,
		logger:    logger.Named(lggr, "EthSubscriber"),
		chStop:    chStop,
	}
}

// backfillLogs - fetches earlier logs either from a relatively recent block (latest minus BlockBackfillDepth) or from the given fromBlockOverride
// note that the whole operation has no timeout - it relies on BlockBackfillSkip (set outside) to optionally prevent very deep, long backfills
// Max runtime is: (10 sec + 1 min * numBlocks/batchSize) * 3 retries
func (sub *ethSubscriber) backfillLogs(fromBlockOverride sql.NullInt64, addresses []common.Address, topics []common.Hash) (chBackfilledLogs chan types.Log, abort bool) {
	sub.logger.Infow("backfilling logs", "from", fromBlockOverride, "addresses", addresses)
	if len(addresses) == 0 {
		sub.logger.Debug("LogBroadcaster: No addresses to backfill for, returning")
		ch := make(chan types.Log)
		close(ch)
		return ch, false
	}

	ctxParent, cancel := sub.chStop.NewCtx()
	defer cancel()

	var latestHeight int64 = -1
	retryCount := 0
	utils.RetryWithBackoff(ctxParent, func() (retry bool) {
		if retryCount > 3 {
			return false
		}
		retryCount++

		if latestHeight < 0 {
			latestBlock, err := sub.ethClient.HeadByNumber(ctxParent, nil)
			if err != nil {
				sub.logger.Warnw("LogBroadcaster: Backfill - could not fetch latest block header, will retry", "err", err)
				return true
			} else if latestBlock == nil {
				sub.logger.Warn("LogBroadcaster: Got nil block header, will retry")
				return true
			}
			latestHeight = latestBlock.Number
		}

		// Backfill from `backfillDepth` blocks ago.  It's up to the subscribers to
		// filter out logs they've already dealt with.
		fromBlock := uint64(latestHeight) - sub.config.BlockBackfillDepth()
		if fromBlock > uint64(latestHeight) {
			fromBlock = 0 // Overflow protection
		}

		if fromBlockOverride.Valid {
			fromBlock = uint64(fromBlockOverride.Int64)
		}

		if fromBlock <= uint64(latestHeight) {
			sub.logger.Infow(fmt.Sprintf("LogBroadcaster: Starting backfill of logs from %v blocks...", uint64(latestHeight)-fromBlock), "fromBlock", fromBlock, "latestHeight", latestHeight)
		} else {
			sub.logger.Infow("LogBroadcaster: Backfilling will be nop because fromBlock is above latestHeight",
				"fromBlock", fromBlock, "latestHeight", latestHeight)
		}

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
		batchSize := int64(sub.config.LogBackfillBatchSize())
		for from := q.FromBlock.Int64(); from <= latestHeight; from += batchSize {
			to := from + batchSize - 1
			if to > latestHeight {
				to = latestHeight
			}
			q.FromBlock = big.NewInt(from)
			q.ToBlock = big.NewInt(to)

			ctx, cancel := context.WithTimeout(ctxParent, time.Minute)
			batchLogs, err := sub.fetchLogBatch(ctx, q, start)
			cancel()

			elapsed := time.Since(start)

			var elapsedMessage string
			if elapsed > time.Minute {
				elapsedMessage = " (backfill is taking a long time, delaying processing of newest logs - if it's an issue, consider setting the EVM.BlockBackfillSkip configuration variable to \"true\")"
			}
			if err != nil {
				if ctx.Err() != nil {
					sub.logger.Errorw("LogBroadcaster: Deadline exceeded, unable to backfill a batch of logs. Consider setting EVM.LogBackfillBatchSize to a lower value", "err", err, "elapsed", elapsed, "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
				} else {
					sub.logger.Errorw("LogBroadcaster: Unable to backfill a batch of logs after retries", "err", err, "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
				}
				return true
			}

			sub.logger.Infow(fmt.Sprintf("LogBroadcaster: Fetched a batch of %v logs from %v to %v%s", len(batchLogs), from, to, elapsedMessage), "len", len(batchLogs), "fromBlock", from, "toBlock", to, "remaining", latestHeight-to)

			select {
			case <-sub.chStop:
				return false
			default:
				logs = append(logs, batchLogs...)
			}
		}

		sub.logger.Infof("LogBroadcaster: Fetched a total of %v logs for backfill", len(logs))

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
			sub.logger.Infof("LogBroadcaster: Finished async backfill of %v logs", len(logs))
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

func (sub *ethSubscriber) fetchLogBatch(ctx context.Context, query ethereum.FilterQuery, start time.Time) ([]types.Log, error) {
	var errOuter error
	var result []types.Log
	utils.RetryWithBackoff(ctx, func() (retry bool) {
		batchLogs, err := sub.ethClient.FilterLogs(ctx, query)

		errOuter = err

		if err != nil {
			if ctx.Err() != nil {
				sub.logger.Errorw("LogBroadcaster: Inner deadline exceeded, unable to backfill a batch of logs. Consider setting EVM.LogBackfillBatchSize to a lower value", "err", err, "elapsed", time.Since(start),
					"fromBlock", query.FromBlock.String(), "toBlock", query.ToBlock.String())
			} else {
				sub.logger.Errorw("LogBroadcaster: Unable to backfill a batch of logs", "err", err,
					"fromBlock", query.FromBlock.String(), "toBlock", query.ToBlock.String())
			}
			return true
		}
		result = batchLogs
		return false
	})
	return result, errOuter
}

// createSubscription creates a new log subscription starting at the current block.  If previous logs
// are needed, they must be obtained through backfilling, as subscriptions can only be started from
// the current head.
func (sub *ethSubscriber) createSubscription(addresses []common.Address, topics []common.Hash) (subscr managedSubscription, abort bool) {
	if len(addresses) == 0 {
		return newNoopSubscription(), false
	}

	ctx, cancel := sub.chStop.NewCtx()
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		filterQuery := ethereum.FilterQuery{
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}
		chRawLogs := make(chan types.Log)

		sub.logger.Debugw("Calling SubscribeFilterLogs with params", "addresses", addresses, "topics", topics)

		innerSub, err := sub.ethClient.SubscribeFilterLogs(ctx, filterQuery, chRawLogs)
		if err != nil {
			sub.logger.Errorw("Log subscriber could not create subscription to Ethereum node", "err", err)
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
	<-sub.Err() // ensure sending has stopped before closing the chan
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
