package log

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// 1. Each listener being registered can specify a custom NumConfirmations - number of block confirmations required for any log being sent to it.
//
// 2. All received logs are kept in an array and deleted ONLY after they are outside the confirmation range for all subscribers
// (when given log height is lower than (latest height - max(highestNumConfirmations, ETH_FINALITY_DEPTH)) ) -> see: pool.go
//
// 3. Information about already consumed logs is fetched from the database and used as a filter
//
// 4. The logs are attempted to be sent after every new head arrival:
// 		Each stored log is checked against every matched listener and is sent unless:
//    A) is too young for that listener
//    B) matches a log already consumed (via the database information from log_broadcasts table)
//
// A log might be sent multiple times, if a consumer processes logs asynchronously (e.g. via a queue or a Mailbox), in which case the log
// may not be marked as consumed before the next sending operation. That's why customers must still check the state via WasAlreadyConsumed
// before processing the log.
//
// The registrations' methods are NOT thread-safe.
type (
	registrations struct {
		// jobID is not unique since one job can have multiple subscribers
		registeredSubs  map[*subscriber]struct{}
		handlersByConfs map[uint32]*handlers
		logger          logger.Logger
		evmChainID      big.Int

		// highest 'NumConfirmations' per all listeners, used to decide about deleting older logs if it's higher than EvmFinalityDepth
		// it's: max(listeners.map(l => l.num_confirmations)
		highestNumConfirmations uint32
	}

	handlers struct {
		handlersByAddr map[common.Address]map[common.Hash]map[Listener]*listenerMetadata // contractAddress => logTopic => Listener
		evmChainID     big.Int
	}

	// The Listener responds to log events through HandleLog.
	Listener interface {
		HandleLog(b Broadcast)
		JobID() int32
	}

	// Metadata structure maintained per listener
	listenerMetadata struct {
		opts    ListenerOpts
		filters [][]Topic
	}
)

func newRegistrations(logger logger.Logger, evmChainID big.Int) *registrations {
	return &registrations{
		registeredSubs:  make(map[*subscriber]struct{}),
		handlersByConfs: make(map[uint32]*handlers),
		evmChainID:      evmChainID,
		logger:          logger.Named("Registrations"),
	}
}

func (r *registrations) addSubscriber(sub *subscriber) (needsResubscribe bool) {
	jobID := sub.listener.JobID()
	if _, exists := r.registeredSubs[sub]; exists {
		r.logger.Panicf("Cannot add subscriber: subscription %p with job ID %v already added", sub, jobID)
	}
	r.registeredSubs[sub] = struct{}{}
	r.logger.Tracef("Removed subscription %p with job ID %v", sub, jobID)

	if _, exists := r.handlersByConfs[sub.opts.MinIncomingConfirmations]; !exists {
		r.handlersByConfs[sub.opts.MinIncomingConfirmations] = newHandlers(r.evmChainID)
	}

	needsResubscribe = r.handlersByConfs[sub.opts.MinIncomingConfirmations].addSubscriber(sub)

	// increase the variable for highest number of confirmations among all subscribers,
	// if the new subscriber has a higher value
	if sub.opts.MinIncomingConfirmations > r.highestNumConfirmations {
		r.highestNumConfirmations = sub.opts.MinIncomingConfirmations
	}
	return
}

func (r *registrations) removeSubscriber(sub *subscriber) (needsResubscribe bool) {
	jobID := sub.listener.JobID()
	if _, exists := r.registeredSubs[sub]; !exists {
		r.logger.Panicf("Cannot remove subscriber: subscription %p with job ID %v is not registered", sub, jobID)
	}
	delete(r.registeredSubs, sub)
	r.logger.Tracef("Added subscription %p with job ID %v", sub, jobID)

	handlers, exists := r.handlersByConfs[sub.opts.MinIncomingConfirmations]
	if !exists {
		return
	}

	needsResubscribe = handlers.removeSubscriber(sub)

	if len(r.handlersByConfs[sub.opts.MinIncomingConfirmations].handlersByAddr) == 0 {
		delete(r.handlersByConfs, sub.opts.MinIncomingConfirmations)
		r.resetHighestNumConfirmationsValue()
	}

	return
}

// reset the number tracking highest num confirmations among all subscribers
func (r *registrations) resetHighestNumConfirmationsValue() {
	highestNumConfirmations := uint32(0)

	for numConfirmations := range r.handlersByConfs {
		if numConfirmations > highestNumConfirmations {
			highestNumConfirmations = numConfirmations
		}
	}
	r.highestNumConfirmations = highestNumConfirmations
}

func (r *registrations) addressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for _, sub := range r.handlersByConfs {
		add, t := sub.addressesAndTopics()
		addresses = append(addresses, add...)
		topics = append(topics, t...)
	}
	return addresses, topics
}

func (r *registrations) isAddressRegistered(address common.Address) bool {
	for _, sub := range r.handlersByConfs {
		if sub.isAddressRegistered(address) {
			return true
		}
	}
	return false
}

func (r *registrations) sendLogs(logsToSend []logsOnBlock, latestHead evmtypes.Head, broadcasts []LogBroadcast, bc broadcastCreator) {
	broadcastsExisting := make(map[LogBroadcastAsKey]bool)
	for _, b := range broadcasts {
		broadcastsExisting[b.AsKey()] = b.Consumed
	}

	latestBlockNumber := uint64(latestHead.Number)

	for _, logsPerBlock := range logsToSend {
		for numConfirmations, handlers := range r.handlersByConfs {

			if numConfirmations != 0 && latestBlockNumber < uint64(numConfirmations) {
				// Skipping send because the block is definitely too young
				continue
			}

			// We attempt the send multiple times per log
			// so here we need to see if this particular listener actually should receive it at this depth
			isOldEnough := numConfirmations == 0 || (logsPerBlock.BlockNumber+uint64(numConfirmations)-1) <= latestBlockNumber
			if !isOldEnough {
				continue
			}

			for _, log := range logsPerBlock.Logs {
				handlers.sendLog(log, latestHead, broadcastsExisting, bc, r.logger)
			}
		}
	}
}

// Returns true if there is at least one filter value (or no filters at all) that matches an actual received value for every index i, or false otherwise
func filtersContainValues(topicValues []common.Hash, filters [][]Topic) bool {
	for i := 0; i < len(topicValues) && i < len(filters); i++ {
		filterValues := filters[i]
		valueFound := len(filterValues) == 0 // empty filter for given index means: all values allowed
		for _, filterValue := range filterValues {
			if common.Hash(filterValue) == topicValues[i] {
				valueFound = true
				break
			}
		}
		if !valueFound {
			return false
		}
	}
	return true
}

func newHandlers(evmChainID big.Int) *handlers {
	return &handlers{
		handlersByAddr: make(map[common.Address]map[common.Hash]map[Listener]*listenerMetadata),
		evmChainID:     evmChainID,
	}
}

func (r *handlers) addSubscriber(sub *subscriber) (needsResubscribe bool) {
	addr := sub.opts.Contract

	if sub.opts.MinIncomingConfirmations <= 0 {
		sub.opts.MinIncomingConfirmations = 1
	}

	if _, exists := r.handlersByAddr[addr]; !exists {
		r.handlersByAddr[addr] = make(map[common.Hash]map[Listener]*listenerMetadata)
	}

	for topic, topicValueFilters := range sub.opts.LogsWithTopics {
		if _, exists := r.handlersByAddr[addr][topic]; !exists {
			r.handlersByAddr[addr][topic] = make(map[Listener]*listenerMetadata)
			needsResubscribe = true
		}

		r.handlersByAddr[addr][topic][sub.listener] = &listenerMetadata{
			opts:    sub.opts,
			filters: topicValueFilters,
		}
	}
	return
}

func (r *handlers) removeSubscriber(sub *subscriber) (needsResubscribe bool) {
	addr := sub.opts.Contract

	// FIXME: What about the case where you remove/add a job with the same contract address?
	// addr is not good enough to be a unique key
	if _, exists := r.handlersByAddr[addr]; !exists {
		return
	}
	for topic := range sub.opts.LogsWithTopics {
		topicMap, exists := r.handlersByAddr[addr][topic]
		if !exists {
			continue
		}

		delete(topicMap, sub.listener)

		if len(topicMap) == 0 {
			needsResubscribe = true
			delete(r.handlersByAddr[addr], topic)
		}
		if len(r.handlersByAddr[addr]) == 0 {
			delete(r.handlersByAddr, addr)
		}
	}
	return
}

func (r *handlers) addressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for addr := range r.handlersByAddr {
		addresses = append(addresses, addr)
		for topic := range r.handlersByAddr[addr] {
			topics = append(topics, topic)
		}
	}
	return addresses, topics
}

func (r *handlers) isAddressRegistered(address common.Address) bool {
	_, exists := r.handlersByAddr[address]
	return exists
}

var _ broadcastCreator = &orm{}

type broadcastCreator interface {
	CreateBroadcast(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, pqOpts ...pg.QOpt) error
}

func (r *handlers) sendLog(log types.Log, latestHead evmtypes.Head,
	broadcasts map[LogBroadcastAsKey]bool,
	bc broadcastCreator,
	logger logger.Logger) {

	latestBlockNumber := uint64(latestHead.Number)
	var wg sync.WaitGroup
	for listener, metadata := range r.handlersByAddr[log.Address][log.Topics[0]] {
		listener := listener

		currentBroadcast := NewLogBroadcastAsKey(log, listener)
		consumed, exists := broadcasts[currentBroadcast]
		if exists && consumed {
			continue
		}

		if len(metadata.filters) > 0 && len(log.Topics) > 1 {
			topicValues := log.Topics[1:]
			if !filtersContainValues(topicValues, metadata.filters) {
				continue
			}
		}

		logCopy := gethwrappers.DeepCopyLog(log)

		var decodedLog generated.AbigenLog
		var err error
		decodedLog, err = metadata.opts.ParseLog(logCopy)
		if err != nil {
			logger.Errorw("Could not parse contract log", "err", err)
			continue
		}

		jobID := listener.JobID()
		if !exists {
			// Create unconsumed broadcast
			if err := bc.CreateBroadcast(log.BlockHash, log.BlockNumber, log.Index, jobID); err != nil {
				logger.Errorw("Could not create broadcast log", "blockNumber", log.BlockNumber,
					"blockHash", log.BlockHash, "address", log.Address, "jobID", jobID, "error", err)
				continue
			}
		}

		logger.Debugw("LogBroadcaster: Sending out log",
			"blockNumber", log.BlockNumber, "blockHash", log.BlockHash,
			"address", log.Address, "latestBlockNumber", latestBlockNumber, "jobID", jobID)

		wg.Add(1)
		go func() {
			defer wg.Done()
			listener.HandleLog(&broadcast{
				latestBlockNumber,
				latestHead.Hash,
				decodedLog,
				logCopy,
				jobID,
				r.evmChainID,
			})
		}()
	}
	wg.Wait()
}
