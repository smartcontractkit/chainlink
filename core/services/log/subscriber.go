package log

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tevino/abool"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type subscriber struct {
	orm           ORM
	ethClient     eth.Client
	relayer       *relayer
	backfillDepth uint64
	connected     *abool.AtomicBool

	contracts   map[common.Address]uint64
	addContract *utils.Mailbox
	rmContract  *utils.Mailbox

	utils.StartStopOnce
	utils.DependentAwaiter
	chStop chan struct{}
	chDone chan struct{}
}

func newSubscriber(
	orm ORM,
	ethClient eth.Client,
	relayer *relayer,
	backfillDepth uint64,
	dependentAwaiter utils.DependentAwaiter,
) *subscriber {
	return &subscriber{
		orm:              orm,
		ethClient:        ethClient,
		relayer:          relayer,
		backfillDepth:    backfillDepth,
		connected:        abool.New(),
		contracts:        make(map[common.Address]uint64),
		addContract:      utils.NewMailbox(0),
		rmContract:       utils.NewMailbox(0),
		DependentAwaiter: dependentAwaiter,
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
}

func (s *subscriber) Start() error {
	return s.StartOnce("Log subscriber", func() (err error) {
		go s.awaitInitialSubscribers()
		return nil
	})
}

func (s *subscriber) awaitInitialSubscribers() {
	for {
		select {
		case <-s.addContract.Notify():
			s.onAddContracts()

		case <-s.rmContract.Notify():
			s.onRmContracts()

		case <-s.DependentAwaiter.AwaitDependents():
			go s.startResubscribeLoop()
			return

		case <-s.chStop:
			close(s.chDone)
			return
		}
	}
}

func (s *subscriber) Stop() error {
	return s.StopOnce("Log subscriber", func() (err error) {
		close(s.chStop)
		<-s.chDone
		return nil
	})
}

func (s *subscriber) IsConnected() bool {
	return s.connected.IsSet()
}

func (s *subscriber) NotifyAddContract(address common.Address) {
	s.addContract.Deliver(address)
}

func (s *subscriber) NotifyRemoveContract(address common.Address) {
	s.rmContract.Deliver(address)
}

// The subscription is closed in two cases:
//   - intentionally, when the set of contracts we're listening to changes
//   - on a connection error
//
// This method recreates the subscription in both cases.  In the event of a connection
// error, it attempts to reconnect.  Any time there's a change in connection state, it
// notifies its subscribers.
func (s *subscriber) startResubscribeLoop() {
	defer close(s.chDone)

	var subscription managedSubscription = newNoopSubscription()
	defer func() { subscription.Unsubscribe() }()

	var chRawLogs chan types.Log
	for {
		newSubscription, abort := s.createSubscription()
		if abort {
			return
		}

		chBackfilledLogs, abort := s.backfillLogs()
		if abort {
			return
		}

		// Each time this loop runs, chRawLogs is reconstituted as:
		//     remaining logs from last subscription <- backfilled logs <- logs from new subscription
		// There will be duplicated logs in this channel.  It is the responsibility of subscribers
		// to account for this using the helpers on the Broadcast type.
		chRawLogs = s.appendLogChannel(chRawLogs, chBackfilledLogs)
		chRawLogs = s.appendLogChannel(chRawLogs, newSubscription.Logs())
		subscription.Unsubscribe()
		subscription = newSubscription

		s.connected.Set()
		s.relayer.NotifyConnected()

		shouldResubscribe, err := s.process(chRawLogs, subscription.Err())
		if err != nil {
			logger.Error(err)
			s.connected.UnSet()
			s.relayer.NotifyDisconnected()
			continue
		} else if !shouldResubscribe {
			s.connected.UnSet()
			s.relayer.NotifyDisconnected()
			return
		}
	}
}

func (s *subscriber) backfillLogs() (chBackfilledLogs chan types.Log, abort bool) {
	if len(s.contracts) == 0 {
		ch := make(chan types.Log)
		close(ch)
		return ch, false
	}

	ctx, cancel := utils.ContextFromChan(s.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		latestBlock, err := s.ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			logger.Errorw("Log subscriber backfill: could not fetch latest block header", "error", err)
			return true
		} else if latestBlock == nil {
			logger.Warn("got nil block header")
			return true
		}
		currentHeight := uint64(latestBlock.Number)

		// Backfill from `backfillDepth` blocks ago.  It's up to the subscribers to
		// filter out logs they've already dealt with.
		fromBlock := currentHeight - s.backfillDepth
		if fromBlock > currentHeight {
			fromBlock = 0 // Overflow protection
		}

		q := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			Addresses: s.addresses(),
		}

		logs, err := s.ethClient.FilterLogs(ctx, q)
		if err != nil {
			logger.Errorw("Log subscriber backfill: could not fetch logs", "error", err)
			return true
		}

		chBackfilledLogs = make(chan types.Log)
		go s.deliverBackfilledLogs(logs, chBackfilledLogs)

		return false
	})
	select {
	case <-s.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

func (s *subscriber) deliverBackfilledLogs(logs []types.Log, chBackfilledLogs chan<- types.Log) {
	defer close(chBackfilledLogs)
	for _, log := range logs {
		select {
		case chBackfilledLogs <- log:
		case <-s.chStop:
			return
		}
	}
}

func (s *subscriber) upsertOrDeleteLogs(logs ...types.Log) {
	for _, log := range logs {
		loggerFields := []interface{}{
			"contract", log.Address,
			"block", log.BlockHash,
			"blockNumber", log.BlockNumber,
			"tx", log.TxHash,
			"logIndex", log.Index,
			"removed", log.Removed,
		}

		if _, exists := s.contracts[log.Address]; !exists {
			logger.Warnw("Log subscriber got log from unknown contract", loggerFields...)
			continue
		}

		if log.Removed {
			err := s.orm.DeleteLogAndBroadcasts(log.BlockHash, log.Index)
			if err != nil {
				loggerFields = append(loggerFields, "error", err)
				logger.Errorw("Log subscriber could not delete reorged log", loggerFields...)
				continue
			}

		} else {
			err := s.orm.UpsertLog(log)
			if err != nil {
				loggerFields = append(loggerFields, "error", err)
				logger.Errorw("Log subscriber could not upsert log", loggerFields...)
				continue
			}
			s.relayer.NotifyNewLog(log)
		}
	}
}

func (s *subscriber) process(chRawLogs <-chan types.Log, chErr <-chan error) (shouldResubscribe bool, _ error) {
	// We debounce requests to subscribe and unsubscribe to avoid making too many
	// RPC calls to the Ethereum node, particularly on startup.
	var needsResubscribe bool
	debounceResubscribe := time.NewTicker(1 * time.Second)
	defer debounceResubscribe.Stop()

	for {
		select {
		case rawLog := <-chRawLogs:
			s.upsertOrDeleteLogs(rawLog)

		case err := <-chErr:
			return true, err

		case <-s.addContract.Notify():
			needsResubscribe = s.onAddContracts() || needsResubscribe

		case <-s.rmContract.Notify():
			needsResubscribe = s.onRmContracts() || needsResubscribe

		case <-debounceResubscribe.C:
			if needsResubscribe {
				return true, nil
			}

		case <-s.chStop:
			return false, nil
		}
	}
}

func (s *subscriber) onAddContracts() (needsResubscribe bool) {
	for {
		x := s.addContract.Retrieve()
		if x == nil {
			break
		}
		addr := x.(common.Address)
		needsResubscribe = s.contracts[addr] == 0 || needsResubscribe
		s.contracts[addr]++
	}
	return
}

func (s *subscriber) onRmContracts() (needsResubscribe bool) {
	for {
		x := s.rmContract.Retrieve()
		if x == nil {
			break
		}
		addr := x.(common.Address)
		if s.contracts[addr] == 0 {
			panic("cannot unsubscribe")
		}
		needsResubscribe = s.contracts[addr] == 1 || needsResubscribe
		s.contracts[addr]--
		if s.contracts[addr] == 0 {
			delete(s.contracts, addr)
		}
	}
	return
}

// createSubscription creates a new log subscription starting at the current block.  If previous logs
// are needed, they must be obtained through backfilling, as subscriptions can only be started from
// the current head.
func (s *subscriber) createSubscription() (sub managedSubscription, abort bool) {
	if len(s.contracts) == 0 {
		return newNoopSubscription(), false
	}

	ctx, cancel := utils.ContextFromChan(s.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		filterQuery := ethereum.FilterQuery{
			Addresses: s.addresses(),
		}
		chRawLogs := make(chan types.Log)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		innerSub, err := s.ethClient.SubscribeFilterLogs(ctx, filterQuery, chRawLogs)
		if err != nil {
			logger.Errorw("Log subscriber could not create subscription to Ethereum node", "error", err)
			return true
		}

		sub = managedSubscriptionImpl{
			subscription: innerSub,
			chRawLogs:    chRawLogs,
		}
		return false
	})
	select {
	case <-s.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

func (s *subscriber) addresses() []common.Address {
	var addresses []common.Address
	for address := range s.contracts {
		addresses = append(addresses, address)
	}
	return addresses
}

func (s *subscriber) appendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	if ch1 == nil && ch2 == nil {
		return nil
	}

	chCombined := make(chan types.Log)

	go func() {
		defer close(chCombined)
		if ch1 != nil {
			for rawLog := range ch1 {
				select {
				case chCombined <- rawLog:
				case <-s.chStop:
					return
				}
			}
		}
		if ch2 != nil {
			for rawLog := range ch2 {
				select {
				case chCombined <- rawLog:
				case <-s.chStop:
					return
				}
			}
		}
	}()

	return chCombined
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

func (s noopSubscription) Err() <-chan error    { return nil }
func (s noopSubscription) Logs() chan types.Log { return s.chRawLogs }
func (s noopSubscription) Unsubscribe()         { close(s.chRawLogs) }
