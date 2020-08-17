package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
		replayFromBlockBN := big.NewInt(replayFromBlock)
		nextHead = replayFromBlockBN
	}

	for _, initr := range initrs {
		unsubscriber, err := NewInitiatorSubscription(initr, store.EthClient, runManager, nextHead, ReceiveLogRequest)
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
	for _, sub := range js.unsubscribers {
		sub.Unsubscribe()
	}
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
	callback func(RunManager, models.LogRequest),
) (InitiatorSubscription, error) {

	filter, err := models.FilterQueryFactory(initr, nextHead)
	if err != nil {
		return InitiatorSubscription{}, errors.Wrap(err, "NewInitiatorSubscription#FilterQueryFactory")
	}

	sub := InitiatorSubscription{
		runManager: runManager,
		Initiator:  initr,
		callback:   callback,
	}

	managedSub, err := NewManagedSubscription(client, filter, sub.dispatchLog)
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
	logSubscriber   eth.Client
	logs            chan models.Log
	ethSubscription ethereum.Subscription
	callback        func(models.Log)
}

// NewManagedSubscription subscribes to the ethereum node with the passed filter
// and delegates incoming logs to callback.
func NewManagedSubscription(
	logSubscriber eth.Client,
	filter ethereum.FilterQuery,
	callback func(models.Log),
) (*ManagedSubscription, error) {
	ctx := context.Background()
	logs := make(chan models.Log)
	es, err := logSubscriber.SubscribeFilterLogs(ctx, filter, logs)
	if err != nil {
		return nil, err
	}

	sub := &ManagedSubscription{
		logSubscriber:   logSubscriber,
		callback:        callback,
		logs:            logs,
		ethSubscription: es,
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
	backfilledSet := sub.backfillLogs(q)
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
				logger.Errorw(fmt.Sprintf("Error in log subscription: %s", err.Error()), "err", err)
			}
		}
	}
}

// Manually retrieve old logs since SubscribeFilterLogs(ctx, filter, chLogs) only returns newly
// imported blocks: https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB#logs
// Therefore TxManager.FilterLogs does a one time retrieval of old logs.
func (sub ManagedSubscription) backfillLogs(q ethereum.FilterQuery) map[string]bool {
	backfilledSet := map[string]bool{}
	if q.FromBlock == nil {
		return backfilledSet
	}

	logs, err := sub.logSubscriber.FilterLogs(context.TODO(), q)
	if err != nil {
		logger.Errorw("Unable to backfill logs", "err", err, "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
		return backfilledSet
	}

	for _, log := range logs {
		backfilledSet[log.BlockHash.String()] = true
		sub.callback(log)
	}
	return backfilledSet
}
