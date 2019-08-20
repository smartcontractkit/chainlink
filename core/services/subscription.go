package services

import (
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
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
func StartJobSubscription(job models.JobSpec, head *models.Head, store *strpkg.Store) (JobSubscription, error) {
	var merr error
	var unsubscribers []Unsubscriber

	initrs := job.InitiatorsFor(
		models.InitiatorEthLog,
		models.InitiatorRunLog,
		models.InitiatorServiceAgreementExecutionLog,
	)

	for _, initr := range initrs {
		unsubscriber, err := NewInitiatorSubscription(initr, job, store, head, ReceiveLogRequest)
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
	Job       models.JobSpec
	Initiator models.Initiator
	store     *strpkg.Store
	callback  func(*strpkg.Store, models.LogRequest)
}

// NewInitiatorSubscription creates a new InitiatorSubscription that feeds received
// logs to the callback func parameter.
func NewInitiatorSubscription(
	initr models.Initiator,
	job models.JobSpec,
	store *strpkg.Store,
	from *models.Head,
	callback func(*strpkg.Store, models.LogRequest),
) (InitiatorSubscription, error) {
	filter, err := models.FilterQueryFactory(initr, from.NextInt()) // Exclude current block from subscription
	if err != nil {
		return InitiatorSubscription{}, errors.Wrap(err, "NewInitiatorSubscription#FilterQueryFactory")
	}

	sub := InitiatorSubscription{
		Job:       job,
		Initiator: initr,
		store:     store,
		callback:  callback,
	}

	managedSub, err := NewManagedSubscription(store, filter, sub.dispatchLog)
	if err != nil {
		return sub, errors.Wrap(err, "NewInitiatorSubscription#NewManagedSubscription")
	}

	sub.ManagedSubscription = managedSub
	loggerLogListening(initr, filter.FromBlock)
	return sub, nil
}

func (sub InitiatorSubscription) dispatchLog(log models.Log) {
	logger.Debugw(fmt.Sprintf("Log for %v initiator for job %v", sub.Initiator.Type, sub.Job.ID),
		"txHash", log.TxHash.Hex(), "logIndex", log.Index, "blockNumber", log.BlockNumber, "job", sub.Job.ID)

	base := models.InitiatorLogEvent{
		JobSpec:   sub.Job,
		Initiator: sub.Initiator,
		Log:       log,
	}
	sub.callback(sub.store, base.LogRequest())
}

func loggerLogListening(initr models.Initiator, blockNumber *big.Int) {
	msg := fmt.Sprintf(
		"Listening for %v from block %v for address %v for job %v",
		initr.Type,
		presenters.FriendlyBigInt(blockNumber),
		utils.LogListeningAddress(initr.Address),
		initr.JobSpecID)
	logger.Infow(msg)
}

// ReceiveLogRequest parses the log and runs the job indicated by a RunLog or
// ServiceAgreementExecutionLog. (Both log events have the same format.)
func ReceiveLogRequest(store *strpkg.Store, le models.LogRequest) {
	if !le.Validate() {
		return
	}

	if le.GetLog().Removed {
		logger.Debugw("Skipping run for removed log", "log", le.GetLog(), "jobId", le.GetJobSpec().ID)
		return
	}

	le.ToDebug()
	data, err := le.JSON()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	runJob(store, le, data)
}

func runJob(store *strpkg.Store, le models.LogRequest, data models.JSON) {
	payment, err := le.ContractPayment()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return
	}

	input := models.RunResult{
		Data:   data,
		Amount: payment,
	}
	if err := le.ValidateRequester(); err != nil {
		input.SetError(err)
		logger.Errorw(err.Error(), le.ForLogger()...)
	}
	rr, err := le.RunRequest()
	if err != nil {
		input.SetError(err)
		logger.Errorw(err.Error(), le.ForLogger()...)
	}

	_, err = ExecuteJobWithRunRequest(
		le.GetJobSpec(),
		le.GetInitiator(),
		input,
		le.BlockNumber(),
		store.Unscoped(),
		rr,
	)
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
	}
}

// ManagedSubscription encapsulates the connecting, backfilling, and clean up of an
// ethereum node subscription.
type ManagedSubscription struct {
	store           *strpkg.Store
	logs            chan models.Log
	ethSubscription models.EthSubscription
	callback        func(models.Log)
}

// NewManagedSubscription subscribes to the ethereum node with the passed filter
// and delegates incoming logs to callback.
func NewManagedSubscription(
	store *strpkg.Store,
	filter ethereum.FilterQuery,
	callback func(models.Log),
) (*ManagedSubscription, error) {
	logs := make(chan models.Log)
	es, err := store.TxManager.SubscribeToLogs(logs, filter)
	if err != nil {
		return nil, err
	}

	sub := &ManagedSubscription{
		store:           store,
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
// https://github.com/smartcontractkit/chainlink/pull/600#issuecomment-426320971
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

// Manually retrieve old logs since SubscribeToLogs(logs, filter) only returns newly
// imported blocks: https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB#logs
// Therefore TxManager.GetLogs does a one time retrieval of old logs.
func (sub ManagedSubscription) backfillLogs(q ethereum.FilterQuery) map[string]bool {
	backfilledSet := map[string]bool{}
	if q.FromBlock == nil {
		return backfilledSet
	}

	logs, err := sub.store.TxManager.GetLogs(q)
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
