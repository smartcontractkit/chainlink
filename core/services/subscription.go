package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// JobSubscription listens to event logs being pushed from the Ethereum Node to a job.
type JobSubscription struct {
	Job                    models.JobSpec
	initiatorSubscriptions []InitiatorSubscription
}

// StartJobSubscription constructs a JobSubscription which listens for and
// tracks event logs corresponding to the specified job. Ignores any errors if
// there is at least one successful subscription to an initiator log.
func StartJobSubscription(job models.JobSpec, head *models.Head, store *strpkg.Store, runManager RunManager, ethClient eth.Client) (JobSubscription, error) {
	var merr error
	var initatorSubscriptions []InitiatorSubscription
	var nextHead *big.Int

	initrs := job.InitiatorsFor(models.LogBasedChainlinkJobInitiators...)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if head == nil {
		latestBlock, err := ethClient.BlockByNumber(ctx, nil)
		if err != nil {
			return JobSubscription{}, err
		}
		backfillDepth := new(big.Int).SetUint64(store.Config.BlockBackfillDepth())
		nextHead = new(big.Int).Sub(latestBlock.Number(), backfillDepth)
		if nextHead.Cmp(big.NewInt(0)) < 1 {
			nextHead = big.NewInt(1)
		}
	} else {
		nextHead = head.NextInt() // Exclude current block from subscription
	}

	if replayFromBlock := store.Config.ReplayFromBlock(); replayFromBlock >= 0 {
		if replayFromBlock >= nextHead.Int64() {
			logger.Infof("StartJobSubscription: next head was supposed to be %v but ReplayFromBlock flag manually overrides to %v, will subscribe from blocknum %v", nextHead, replayFromBlock, replayFromBlock)
			replayFromBlockBN := big.NewInt(replayFromBlock)
			nextHead = replayFromBlockBN
		}
		logger.Warnf("StartJobSubscription: replayFromBlock was set to %v which is older than the next head of %v, will subscribe from blocknum %v", replayFromBlock, nextHead, nextHead)
	}

	for _, initr := range initrs {
		filter, err := models.FilterQueryFactory(initr, nextHead, store.Config.OperatorContractAddress())
		if err != nil {
			merr = multierr.Append(merr, err)
			continue
		}
		is, err := NewInitiatorSubscription(initr, ethClient, runManager, filter, store.Config.EthLogBackfillBatchSize(), ProcessLogRequest)
		if err != nil {
			merr = multierr.Append(merr, err)
		} else {
			is.Start()
			initatorSubscriptions = append(initatorSubscriptions, *is)
		}
	}

	if len(initatorSubscriptions) == 0 {
		return JobSubscription{}, multierr.Append(
			merr, errors.New(
				"unable to subscribe to any logs, check earlier errors in this message, and the initiator types"))
	}
	return JobSubscription{Job: job, initiatorSubscriptions: initatorSubscriptions}, merr
}

// Unsubscribe stops the subscription and cleans up associated resources.
func (js JobSubscription) Unsubscribe() {
	var wg sync.WaitGroup
	wg.Add(len(js.initiatorSubscriptions))
	for _, sub := range js.initiatorSubscriptions {
		go func(s InitiatorSubscription) {
			defer wg.Done()
			logger.Debugw("JobSubscription: unsubscribing", "initiator", s.Initiator)
			s.Unsubscribe()
		}(sub)
	}
	wg.Wait()
}

// InitiatorSubscription encapsulates all functionality needed to wrap an ethereum subscription
// for use with a Chainlink Initiator. Initiator specific functionality is delegated
// to the callback.
type InitiatorSubscription struct {
	done              chan struct{}
	logSubscriber     eth.Client
	logs              chan types.Log
	ethSubscription   ethereum.Subscription
	filter            ethereum.FilterQuery
	backfillBatchSize uint32
	runManager        RunManager
	Initiator         models.Initiator
	callback          func(RunManager, models.LogRequest)
}

// NewInitiatorSubscription creates a new InitiatorSubscription that feeds received
// logs to the callback func parameter.
func NewInitiatorSubscription(
	initr models.Initiator,
	client eth.Client,
	runManager RunManager,
	filter ethereum.FilterQuery,
	backfillBatchSize uint32,
	callback func(RunManager, models.LogRequest),
) (*InitiatorSubscription, error) {
	logs := make(chan types.Log)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	es, err := client.SubscribeFilterLogs(ctx, filter, logs)
	if err != nil {
		return nil, err
	}
	return &InitiatorSubscription{
		Initiator:         initr,
		done:              make(chan struct{}),
		logSubscriber:     client,
		ethSubscription:   es,
		filter:            filter,
		logs:              logs,
		runManager:        runManager,
		backfillBatchSize: backfillBatchSize,
		callback:          callback,
	}, nil
}

func (sub *InitiatorSubscription) Start() {
	go gracefulpanic.WrapRecover(func() {
		sub.listenForLogs()
	})
}

func (sub *InitiatorSubscription) listenForLogs() {
	// If we spend too long backfilling without processing
	// logs from our subscription, geth will consider the client dead
	// and drop the subscription, so we set an upper bound on backlog processing time.
	// https://github.com/ethereum/go-ethereum/blob/2e5d14170846ae72adc47467a1129e41d6800349/rpc/client.go#L430
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	l := logger.CreateLogger(logger.Default.With(
		"jobID", sub.Initiator.JobSpecID,
		"initiatorID", sub.Initiator.ID,
		"type", sub.Initiator.Type,
		"topics", sub.Initiator.InitiatorParams.Topics,
	))
	backfilledSet := sub.backfillLogs(ctx, sub.filter, *l)
	l.Debugw("InitiatorSubscription: listening for logs",
		"fromBlock", sub.filter.FromBlock,
	)
	for {
		select {
		case <-sub.done:
			close(sub.logs)
			return
		case log := <-sub.logs:
			if _, present := backfilledSet[log.BlockHash.String()]; !present {
				logger.Infow("InitiatorSubscription: log received",
					"blockNumber", log.BlockNumber,
					"txHash", log.TxHash.Hex(),
					"logIndex", log.Index,
					"address", log.Address,
					"log", log,
				)
				sub.callback(sub.runManager, models.InitiatorLogEvent{
					Initiator: sub.Initiator,
					Log:       log,
				}.LogRequest())
			}
		case err, ok := <-sub.ethSubscription.Err():
			// If !ok, then we intentionally closed the subscription,
			// do not try and reconnect.
			if ok {
				l.Warnw("InitiatorSubscription: error in log subscription, attempting to reconnect to eth node", "err", err)
				b := &backoff.Backoff{
					Min:    100 * time.Millisecond,
					Max:    10 * time.Second,
					Factor: 2,
					Jitter: false,
				}
				for {
					newLogs := make(chan types.Log)
					newSub, err := sub.logSubscriber.SubscribeFilterLogs(context.Background(), sub.filter, newLogs)
					if err != nil {
						l.Warnw(fmt.Sprintf("InitiatorSubscription: failed to reconnect to eth node, trying again in %v", b.Duration()), "err", err.Error())
						time.Sleep(b.Duration())
						continue
					}
					sub.ethSubscription = newSub
					sub.logs = newLogs
					l.Infow("InitiatorSubscription: successfully reconnected to eth node")
					break
				}
			}
		}
	}
}

// Manually retrieve old logs since SubscribeFilterLogs(ctx, filter, chLogs) only returns newly
// imported blocks: https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB#logs
// Therefore TxManager.FilterLogs does a one time retrieval of old logs.
func (sub *InitiatorSubscription) backfillLogs(ctx context.Context, q ethereum.FilterQuery, l logger.Logger) map[string]bool {
	start := time.Now()
	backfilledSet := map[string]bool{}
	if q.FromBlock == nil {
		return backfilledSet
	}
	b, err := sub.logSubscriber.BlockByNumber(ctx, nil)
	if err != nil {
		l.Errorw("InitiatorSubscriber: unable to backfill logs couldn't read latest block", "err", err)
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
				l.Errorw("InitiatorSubscriber: deadline exceeded, unable to backfill logs", "err", err, "elapsed", time.Since(start), "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
			} else {
				l.Errorw("InitiatorSubscriber: unable to backfill logs", "err", err, "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
			}
			return backfilledSet
		}
		for _, log := range batchLogs {
			select {
			case <-ctx.Done():
				l.Errorw("InitiatorSubscriber: deadline exceeded, unable to backfill logs", "elapsed", time.Since(start), "fromBlock", q.FromBlock.String(), "toBlock", q.ToBlock.String())
				return backfilledSet
			default:
				backfilledSet[log.BlockHash.String()] = true
				logger.Infow("InitiatorSubscription: backfilled log received",
					"blockNumber", log.BlockNumber,
					"txHash", log.TxHash.Hex(),
					"logIndex", log.Index,
					"address", log.Address,
					"log", log,
				)
				sub.callback(sub.runManager, models.InitiatorLogEvent{
					Initiator: sub.Initiator,
					Log:       log,
				}.LogRequest())
			}
		}
	}
	return backfilledSet
}

// Unsubscribe closes channels and cleans up resources.
func (sub *InitiatorSubscription) Unsubscribe() {
	if sub.ethSubscription != nil {
		sub.ethSubscription.Unsubscribe()
	}
	close(sub.done)
}

// ReceiveLogRequest parses the log and runs the job it indicated by its
// GetJobSpecID method
func ProcessLogRequest(runManager RunManager, le models.LogRequest) {
	if !le.Validate() {
		logger.Debugw("InitiatorSubscription: discarding invalid event log", le.ForLogger()...)
		return
	}

	if le.GetLog().Removed {
		logger.Debugw("InitiatorSubscription: skipping run for removed log", le.ForLogger()...)
		return
	}

	le.ToDebug()
	jobSpecID := le.GetJobSpecID()
	initiator := le.GetInitiator()

	if err := le.ValidateRequester(); err != nil {
		if _, e := runManager.CreateErrored(jobSpecID, initiator, err); e != nil {
			logger.Errorw("InitiatorSubscription: invalid requester, error creating errored job", le.ForLogger("err", e.Error())...)
		}
		logger.Errorw("InitiatorSubscription: invalid requester, created errored job", le.ForLogger("err", err)...)
		return
	}

	rr, err := le.RunRequest()
	if err != nil {
		if _, e := runManager.CreateErrored(jobSpecID, initiator, err); e != nil {
			logger.Errorw("InitiatorSubscription: invalid run request, error creating errored job", le.ForLogger("err", e.Error())...)
		}
		logger.Errorw("InitiatorSubscription: invalid run request, created errored job", le.ForLogger("err", err)...)
		return
	}

	_, err = runManager.Create(jobSpecID, &initiator, le.BlockNumber(), &rr)
	if err != nil {
		logger.Errorw("InitiatorSubscription: error creating run from log", le.ForLogger("err", err)...)
	}
}
