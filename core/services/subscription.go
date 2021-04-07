package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/jpillora/backoff"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// Unsubscriber is the interface for all subscriptions, allowing one to unsubscribe.
type Unsubscriber interface {
	Unsubscribe()
}

// JobSubscription listens to event logs being pushed from the Ethereum Node to a job.
type JobSubscription struct {
	Job           models.JobSpec
	unsubscribers []Unsubscriber
}

// StartJobSubscription constructs a JobSubscription which listens for and
// tracks event logs corresponding to the specified job. Ignores any errors if
// there is at least one successful subscription to an initiator log.
func StartJobSubscription(job models.JobSpec, head *models.Head, store *strpkg.Store, runManager RunManager) (JobSubscription, error) {
	var merr error
	var unsubscribers []Unsubscriber

	initrs := job.InitiatorsFor(models.LogBasedChainlinkJobInitiators...)

	nextHead := head.NextInt() // Exclude current block from subscription
	if replayFromBlock := store.Config.ReplayFromBlock(); replayFromBlock >= 0 {
		if replayFromBlock >= nextHead.Int64() {
			logger.Infof("StartJobSubscription: Next head was supposed to be %v but ReplayFromBlock flag manually overrides to %v, will subscribe from blocknum %v", nextHead, replayFromBlock, replayFromBlock)
			replayFromBlockBN := big.NewInt(replayFromBlock)
			nextHead = replayFromBlockBN
		}
		logger.Warnf("StartJobSubscription: ReplayFromBlock was set to %v which is older than the next head of %v, will subscribe from blocknum %v", replayFromBlock, nextHead, nextHead)
	}

	for _, initr := range initrs {
		unsubscriber, err := NewInitiatorSubscription(initr, store.EthClient, runManager, nextHead, store.Config, ReceiveLogRequest)
		if err == nil {
			unsubscribers = append(unsubscribers, unsubscriber)
		} else {
			merr = multierr.Append(merr, err)
		}
	}

	if len(unsubscribers) == 0 {
		return JobSubscription{}, multierr.Append(
			merr, errors.New(
				"unable to subscribe to any logs, check earlier errors in this message, and the initiator types"))
	}
	return JobSubscription{Job: job, unsubscribers: unsubscribers}, merr
}

// Unsubscribe stops the subscription and cleans up associated resources.
func (js JobSubscription) Unsubscribe() {
	waiter := sync.WaitGroup{}
	waiter.Add(len(js.unsubscribers))
	for _, sub := range js.unsubscribers {
		go func(sub Unsubscriber) {
			sub.Unsubscribe()
			waiter.Done()
		}(sub)
	}
	waiter.Wait()
}

// InitiatorSubscription encapsulates all functionality needed to wrap an ethereum subscription
// for use with a Chainlink Initiator. Initiator specific functionality is delegated
// to the callback.
type InitiatorSubscription struct {
	*ManagedSubscription
	runManager RunManager
	Initiator  models.Initiator
	callback   func(RunManager, models.LogRequest)
}

// NewInitiatorSubscription creates a new InitiatorSubscription that feeds received
// logs to the callback func parameter.
func NewInitiatorSubscription(
	initr models.Initiator,
	client eth.Client,
	runManager RunManager,
	nextHead *big.Int,
	config orm.ConfigReader,
	callback func(RunManager, models.LogRequest),
) (InitiatorSubscription, error) {

	filter, err := models.FilterQueryFactory(initr, nextHead, config.OperatorContractAddress())
	if err != nil {
		return InitiatorSubscription{}, errors.Wrap(err, "NewInitiatorSubscription#FilterQueryFactory")
	}

	sub := InitiatorSubscription{
		runManager: runManager,
		Initiator:  initr,
		callback:   callback,
	}

	managedSub, err := NewManagedSubscription(client, filter, sub.dispatchLog, config.EthLogBackfillBatchSize())
	if err != nil {
		return sub, errors.Wrap(err, "NewInitiatorSubscription#NewManagedSubscription")
	}

	sub.ManagedSubscription = managedSub
	loggerLogListening(initr, filter.FromBlock)
	return sub, nil
}

func (sub InitiatorSubscription) dispatchLog(log models.Log) {
	logger.Debugw(fmt.Sprintf("Log for %v initiator for job %s", sub.Initiator.Type, sub.Initiator.JobSpecID.String()),
		"txHash", log.TxHash.Hex(), "logIndex", log.Index, "blockNumber", log.BlockNumber, "job", sub.Initiator.JobSpecID.String())

	base := models.InitiatorLogEvent{
		Initiator: sub.Initiator,
		Log:       log,
	}
	sub.callback(sub.runManager, base.LogRequest())
}

func loggerLogListening(initr models.Initiator, blockNumber *big.Int) {
	msg := fmt.Sprintf("Listening for %v from block %v", initr.Type, presenters.FriendlyBigInt(blockNumber))
	logger.Infow(msg, "address", utils.LogListeningAddress(initr.Address), "jobID", initr.JobSpecID.String())
}

// ReceiveLogRequest parses the log and runs the job it indicated by its
// GetJobSpecID method
func ReceiveLogRequest(runManager RunManager, le models.LogRequest) {
	if !le.Validate() {
		logger.Debugw("discarding INVALID EVENT LOG", "log", le.GetLog())
		return
	}

	if le.GetLog().Removed {
		logger.Debugw("Skipping run for removed log", "log", le.GetLog(), "jobId", le.GetJobSpecID().String())
		return
	}

	le.ToDebug()

	runJob(runManager, le)

}

func runJob(runManager RunManager, le models.LogRequest) {
	jobSpecID := le.GetJobSpecID()
	initiator := le.GetInitiator()

	if err := le.ValidateRequester(); err != nil {
		if _, e := runManager.CreateErrored(jobSpecID, initiator, err); e != nil {
			logger.Errorw(e.Error())
		}
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	rr, err := le.RunRequest()
	if err != nil {
		if _, e := runManager.CreateErrored(jobSpecID, initiator, err); e != nil {
			logger.Errorw(e.Error())
		}
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	_, err = runManager.Create(jobSpecID, &initiator, le.BlockNumber(), &rr)
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
	}
}

// ManagedSubscription encapsulates the connecting, backfilling, and clean up of an
// ethereum node subscription.
type ManagedSubscription struct {
	logSubscriber     eth.Client
	logs              chan models.Log
	ethSubscription   ethereum.Subscription
	callback          func(models.Log)
	backfillBatchSize uint32
}

// NewManagedSubscription subscribes to the ethereum node with the passed filter
// and delegates incoming logs to callback.
func NewManagedSubscription(logSubscriber eth.Client, filter ethereum.FilterQuery, callback func(models.Log), backfillBatchSize uint32) (*ManagedSubscription, error) {
	ctx := context.Background()
	logs := make(chan models.Log)
	es, err := logSubscriber.SubscribeFilterLogs(ctx, filter, logs)
	if err != nil {
		return nil, err
	}

	sub := &ManagedSubscription{
		logSubscriber:     logSubscriber,
		callback:          callback,
		logs:              logs,
		ethSubscription:   es,
		backfillBatchSize: backfillBatchSize,
	}
	go sub.listenToLogs(filter)
	return sub, nil
}

// Unsubscribe closes channels and cleans up resources.
func (sub ManagedSubscription) Unsubscribe() {
	if sub.ethSubscription != nil {
		timedUnsubscribe(sub.ethSubscription)
	}
	close(sub.logs)
}

// timedUnsubscribe attempts to unsubscribe but aborts abruptly after a time delay
// unblocking the application. This is an effort to mitigate the occasional
// indefinite block described here from go-ethereum:
// https://chainlink/pull/600#issuecomment-426320971
func timedUnsubscribe(unsubscriber Unsubscriber) {
	unsubscribed := make(chan struct{})
	go func() {
		unsubscriber.Unsubscribe()
		close(unsubscribed)
	}()
	select {
	case <-unsubscribed:
	case <-time.After(100 * time.Millisecond):
		logger.Warnf("Subscription %T Unsubscribe timed out.", unsubscriber)
	}
}

func (sub ManagedSubscription) listenToLogs(q ethereum.FilterQuery) {
	// If we spend too long backfilling without processing
	// logs from our subscription, geth will consider the client dead
	// and drop the subscription, so we set an upper bound on backlog processing time.
	// https://github.com/ethereum/go-ethereum/blob/2e5d14170846ae72adc47467a1129e41d6800349/rpc/client.go#L430
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	backfilledSet := sub.backfillLogs(ctx, q)
	for {
		select {
		case log, open := <-sub.logs:
			if !open {
				return
			}
			if _, present := backfilledSet[log.BlockHash.String()]; !present {
				sub.callback(log)
			}
		case err, ok := <-sub.ethSubscription.Err():
			if ok {
				logger.Warnw("Error in log subscription. Attempting to reconnect to eth node", "err", err)
				b := &backoff.Backoff{
					Min:    100 * time.Millisecond,
					Max:    10 * time.Second,
					Factor: 2,
					Jitter: false,
				}
				for {
					newLogs := make(chan models.Log)
					newSub, err := sub.logSubscriber.SubscribeFilterLogs(context.Background(), q, newLogs)
					if err != nil {
						logger.Warnw(fmt.Sprintf("Failed to reconnect to eth node. Trying again in %v", b.Duration()), "err", err.Error())
						time.Sleep(b.Duration())
						continue
					}
					sub.ethSubscription = newSub
					sub.logs = newLogs
					logger.Infow("Successfully reconnected to eth node.")
					break
				}
			}
		}
	}
}

// Manually retrieve old logs since SubscribeFilterLogs(ctx, filter, chLogs) only returns newly
// imported blocks: https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB#logs
// Therefore TxManager.FilterLogs does a one time retrieval of old logs.
func (sub ManagedSubscription) backfillLogs(ctx context.Context, q ethereum.FilterQuery) map[string]bool {
	start := time.Now()
	backfilledSet := map[string]bool{}
	if q.FromBlock == nil {
		return backfilledSet
	}
	b, err := sub.logSubscriber.BlockByNumber(ctx, nil)
	if err != nil {
		logger.Errorw("Unable to backfill logs: couldn't read latest block", "err", err)
		return backfilledSet
	}

	// If we are significantly behind the latest head, there could be a very large (1000s)
	// of blocks to check for logs. We read the blocks in batches to avoid hitting the websocket
	// request data limit.
	// On matic its 5MB [https://github.com/maticnetwork/bor/blob/3de2110886522ab17e0b45f3c4a6722da72b7519/rpc/http.go#L35]
	// On ethereum its 15MB [https://github.com/ethereum/go-ethereum/blob/master/rpc/websocket.go#L40]
	latest := b.Number()
	batchSize := int64(sub.backfillBatchSize)
	for i := q.FromBlock.Int64(); i < latest.Int64(); i += batchSize {
		q.FromBlock = big.NewInt(i)
		to := utils.BigIntSlice{big.NewInt(i + batchSize - 1), latest}
		q.ToBlock = to.Min()
		batchLogs, err := sub.logSubscriber.FilterLogs(ctx, q)
		if err != nil {
			if ctx.Err() != nil {
				logger.Errorw("Deadline exceeded, unable to backfill logs", "err", err, "elapsed", time.Since(start), "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
			} else {
				logger.Errorw("Unable to backfill logs", "err", err, "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
			}
			return backfilledSet
		}
		for _, log := range batchLogs {
			select {
			case <-ctx.Done():
				logger.Errorw("Deadline exceeded, unable to backfill logs", "elapsed", time.Since(start), "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
				return backfilledSet
			default:
				backfilledSet[log.BlockHash.String()] = true
				sub.callback(log)
			}
		}
	}
	return backfilledSet
}
