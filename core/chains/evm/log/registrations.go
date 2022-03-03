package log

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

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
		// Map only used for invariant checking:
		// registeredSubs is used to sanity check adding/removing the exact same subscriber twice
		registeredSubs map[*subscriber]struct{}
		// Map only used for invariant checking:
		// jobIDAddr enforces that no two listeners can share the same jobID and contract address
		// This is because log_broadcasts table can only be consumed once and
		// assumes one listener per job per log event
		jobIDAddrs map[int32]map[common.Address]struct{}

		// handlersByConfs maps numConfirmations => *handler
		handlersByConfs map[uint32]*handler
		logger          logger.Logger
		evmChainID      big.Int

		// highest 'NumConfirmations' per all listeners, used to decide about deleting older logs if it's higher than EvmFinalityDepth
		// it's: max(listeners.map(l => l.num_confirmations)
		highestNumConfirmations uint32
	}

	handler struct {
		lookupSubs map[common.Address]map[common.Hash]subscribers // contractAddress => logTopic => *subscriber => topicValueFilters
		evmChainID big.Int
		logger     logger.Logger
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

	// subscribers type for convenience and readability
	subscribers map[*subscriber][][]Topic
)

func newRegistrations(logger logger.Logger, evmChainID big.Int) *registrations {
	return &registrations{
		registeredSubs:  make(map[*subscriber]struct{}),
		jobIDAddrs:      make(map[int32]map[common.Address]struct{}),
		handlersByConfs: make(map[uint32]*handler),
		evmChainID:      evmChainID,
		logger:          logger.Named("Registrations"),
	}
}

func (r *registrations) addSubscriber(sub *subscriber) (needsResubscribe bool) {
	if err := r.checkAddSubscriber(sub); err != nil {
		r.logger.Panicw(err.Error(), "err", err, "addr", sub.opts.Contract.Hex(), "jobID", sub.listener.JobID())
	}

	r.logger.Tracef("Added subscription %p with job ID %v", sub, sub.listener.JobID())

	handler, exists := r.handlersByConfs[sub.opts.MinIncomingConfirmations]
	if !exists {
		handler = newHandler(r.logger, r.evmChainID)
		r.handlersByConfs[sub.opts.MinIncomingConfirmations] = handler
	}

	needsResubscribe = handler.addSubscriber(sub, r.handlersWithGreaterConfs(sub.opts.MinIncomingConfirmations))

	// increase the variable for highest number of confirmations among all subscribers,
	// if the new subscriber has a higher value
	if sub.opts.MinIncomingConfirmations > r.highestNumConfirmations {
		r.highestNumConfirmations = sub.opts.MinIncomingConfirmations
	}
	return
}

// handlersWithGreaterConfs allows for an optimisation - in the case that we
// are already listening on this topic for a handler with a GREATER
// MinIncomingConfirmations, it is not necessary to subscribe again
func (r *registrations) handlersWithGreaterConfs(confs uint32) (handlersWithGreaterConfs []*handler) {
	for hConfs, handler := range r.handlersByConfs {
		if hConfs > confs {
			handlersWithGreaterConfs = append(handlersWithGreaterConfs, handler)
		}
	}
	return
}

// checkAddSubscriber registers the subsciber and makes sure we aren't violating any assumptions
// maps modified are only used for checks
func (r *registrations) checkAddSubscriber(sub *subscriber) error {
	if sub.opts.MinIncomingConfirmations <= 0 {
		return errors.Errorf("LogBroadcaster requires that MinIncomingConfirmations must be at least 1 (got %v). Logs must have been confirmed in at least 1 block, it does not support reading logs from the mempool before they have been mined.", sub.opts.MinIncomingConfirmations)
	}

	jobID := sub.listener.JobID()
	if _, exists := r.registeredSubs[sub]; exists {
		return errors.Errorf("Cannot add subscriber %p for job ID %v: already added", sub, jobID)
	}
	r.registeredSubs[sub] = struct{}{}
	addrs, exists := r.jobIDAddrs[jobID]
	if !exists {
		r.jobIDAddrs[jobID] = make(map[common.Address]struct{})
	}
	if _, exists := addrs[sub.opts.Contract]; exists {
		return errors.Errorf("Cannot add subscriber %p: only one subscription is allowed per jobID/contract address. There is already a subscription with job ID %v listening on %s", sub, jobID, sub.opts.Contract.Hex())
	}
	r.jobIDAddrs[jobID][sub.opts.Contract] = struct{}{}
	return nil
}

func (r *registrations) removeSubscriber(sub *subscriber) (needsResubscribe bool) {
	if err := r.checkRemoveSubscriber(sub); err != nil {
		r.logger.Panicw(err.Error(), "err", err, "addr", sub.opts.Contract.Hex(), "jobID", sub.listener.JobID())
	}
	r.logger.Tracef("Removed subscription %p with job ID %v", sub, sub.listener.JobID())

	handlers, exists := r.handlersByConfs[sub.opts.MinIncomingConfirmations]
	if !exists {
		return
	}

	needsResubscribe = handlers.removeSubscriber(sub, r.handlersByConfs)

	if len(r.handlersByConfs[sub.opts.MinIncomingConfirmations].lookupSubs) == 0 {
		delete(r.handlersByConfs, sub.opts.MinIncomingConfirmations)
		r.resetHighestNumConfirmationsValue()
	}

	return
}

// checkRemoveSubscriber deregisters the subscriber and validates we aren't
// violating any assumptions
// maps modified are only used for checks
func (r *registrations) checkRemoveSubscriber(sub *subscriber) error {
	jobID := sub.listener.JobID()
	if _, exists := r.registeredSubs[sub]; !exists {
		return errors.Errorf("Cannot remove subscriber %p for job ID %v: not registered", sub, jobID)
	}
	delete(r.registeredSubs, sub)
	addrs, exists := r.jobIDAddrs[jobID]
	if !exists {
		return errors.Errorf("Cannot remove subscriber %p: jobIDAddrs was missing job ID %v", sub, jobID)
	}
	_, exists = addrs[sub.opts.Contract]
	if !exists {
		return errors.Errorf("Cannot remove subscriber %p: jobIDAddrs was missing address %s", sub, sub.opts.Contract.Hex())
	}
	delete(r.jobIDAddrs[jobID], sub.opts.Contract)
	if len(r.jobIDAddrs[jobID]) == 0 {
		delete(r.jobIDAddrs, jobID)
	}
	return nil
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

func newHandler(lggr logger.Logger, evmChainID big.Int) *handler {
	return &handler{
		lookupSubs: make(map[common.Address]map[common.Hash]subscribers),
		evmChainID: evmChainID,
		logger:     lggr,
	}
}

func (r *handler) addSubscriber(sub *subscriber, handlersWithGreaterConfs []*handler) (needsResubscribe bool) {
	addr := sub.opts.Contract

	if sub.opts.MinIncomingConfirmations <= 0 {
		r.logger.Panicw(fmt.Sprintf("LogBroadcaster requires that MinIncomingConfirmations must be at least 1 (got %v). Logs must have been confirmed in at least 1 block, it does not support reading logs from the mempool before they have been mined.", sub.opts.MinIncomingConfirmations), "addr", sub.opts.Contract.Hex(), "jobID", sub.listener.JobID())
	}

	if _, exists := r.lookupSubs[addr]; !exists {
		r.lookupSubs[addr] = make(map[common.Hash]subscribers)
	}

	for topic, topicValueFilters := range sub.opts.LogsWithTopics {
		if _, exists := r.lookupSubs[addr][topic]; !exists {
			r.logger.Tracef("No existing sub for addr %s and topic %s at this MinIncomingConfirmations of %v", addr.Hex(), topic.Hex(), sub.opts.MinIncomingConfirmations)
			r.lookupSubs[addr][topic] = make(subscribers)

			func() {
				if !needsResubscribe {
					// NOTE: This is an optimization; if we already have a
					// subscription to this addr/topic at a higher
					// MinIncomingConfirmations then we don't need to resubscribe
					// again since even the worst case lookback is already covered
					for _, existingHandler := range handlersWithGreaterConfs {
						if _, exists := existingHandler.lookupSubs[addr][topic]; exists {
							r.logger.Tracef("Sub already exists for addr %s and topic %s at greater than this MinIncomingConfirmations of %v. Resubscribe is not required", addr.Hex(), topic.Hex(), sub.opts.MinIncomingConfirmations)
							return
						}
					}
					r.logger.Tracef("No sub exists for addr %s and topic %s at this or greater MinIncomingConfirmations of %v. Resubscribe is required", addr.Hex(), topic.Hex(), sub.opts.MinIncomingConfirmations)
					needsResubscribe = true
				}
			}()
		}
		r.lookupSubs[addr][topic][sub] = topicValueFilters
	}
	return
}

func (r *handler) removeSubscriber(sub *subscriber, allHandlers map[uint32]*handler) (needsResubscribe bool) {
	addr := sub.opts.Contract

	for topic := range sub.opts.LogsWithTopics {
		// OK to panic on missing addr/topic here, since that would be an invariant violation:
		// Both addr and topic will always have been added on addSubscriber
		// LogsWithTopics should never be mutated
		// Only removeSubscriber should ever remove anything from this map
		addrTopics, exists := r.lookupSubs[addr]
		if !exists {
			r.logger.Panicf("AssumptionViolation: expected lookupSubs to contain addr %s for subscriber %p with job ID %v", addr.Hex(), sub, sub.listener.JobID())
		}
		topicMap, exists := addrTopics[topic]
		if !exists {
			r.logger.Panicf("AssumptionViolation: expected addrTopics to contain topic %v for subscriber %p with job ID %v", topic, sub, sub.listener.JobID())
		}
		if _, exists = topicMap[sub]; !exists {
			r.logger.Panicf("AssumptionViolation: expected topicMap to contain subscriber %p with job ID %v", sub, sub.listener.JobID())
		}
		delete(topicMap, sub)

		// cleanup and resubscribe if necessary
		if len(topicMap) == 0 {
			r.logger.Tracef("No subs left for addr %s and topic %s at this MinIncomingConfirmations of %v", addr.Hex(), topic.Hex(), sub.opts.MinIncomingConfirmations)

			func() {
				if !needsResubscribe {
					// NOTE: This is an optimization. Resub not necessary if there
					// are still any other handlers listening on this addr/topic.
					for confs, otherHandler := range allHandlers {
						if confs == sub.opts.MinIncomingConfirmations {
							// no need to check ourself, already did this above
							continue
						}
						if _, exists := otherHandler.lookupSubs[addr][topic]; exists {
							r.logger.Tracef("Sub still exists for addr %s and topic %s. Resubscribe will not be performed", addr.Hex(), topic.Hex())
							return
						}
					}

					r.logger.Tracef("No sub exists for addr %s and topic %s. Resubscribe will be performed", addr.Hex(), topic.Hex())
					needsResubscribe = true
				}
			}()
			delete(r.lookupSubs[addr], topic)
		}
		if len(r.lookupSubs[addr]) == 0 {
			delete(r.lookupSubs, addr)
		}
	}
	return
}

func (r *handler) addressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for addr := range r.lookupSubs {
		addresses = append(addresses, addr)
		for topic := range r.lookupSubs[addr] {
			topics = append(topics, topic)
		}
	}
	return addresses, topics
}

func (r *handler) isAddressRegistered(addr common.Address) bool {
	_, exists := r.lookupSubs[addr]
	return exists
}

var _ broadcastCreator = &orm{}

type broadcastCreator interface {
	CreateBroadcast(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, pqOpts ...pg.QOpt) error
}

func (r *handler) sendLog(log types.Log, latestHead evmtypes.Head,
	broadcasts map[LogBroadcastAsKey]bool,
	bc broadcastCreator,
	logger logger.Logger) {

	topic := log.Topics[0]

	latestBlockNumber := uint64(latestHead.Number)
	var wg sync.WaitGroup
	for sub, filters := range r.lookupSubs[log.Address][topic] {
		currentBroadcast := NewLogBroadcastAsKey(log, sub.listener)
		consumed, exists := broadcasts[currentBroadcast]
		if exists && consumed {
			continue
		}

		if len(filters) > 0 && len(log.Topics) > 1 {
			topicValues := log.Topics[1:]
			if !filtersContainValues(topicValues, filters) {
				continue
			}
		}

		logCopy := gethwrappers.DeepCopyLog(log)

		var decodedLog generated.AbigenLog
		var err error
		decodedLog, err = sub.opts.ParseLog(logCopy)
		if err != nil {
			logger.Errorw("Could not parse contract log", "err", err)
			continue
		}

		jobID := sub.listener.JobID()
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

		// must copy function pointer here since range pointer (sub) may not be
		// used in goroutine below
		handleLog := sub.listener.HandleLog
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleLog(&broadcast{
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
