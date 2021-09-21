package log

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
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
		subscribers map[uint64]*subscribers
		decoders    map[common.Address]ParseLogFunc
		logger      logger.Logger
		evmChainID  big.Int

		// highest 'NumConfirmations' per all listeners, used to decide about deleting older logs if it's higher than EvmFinalityDepth
		// it's: max(listeners.map(l => l.num_confirmations)
		highestNumConfirmations uint64
	}

	subscribers struct {
		handlers   map[common.Address]map[common.Hash]map[Listener]*listenerMetadata // contractAddress => logTopic => Listener
		evmChainID big.Int
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
		subscribers: make(map[uint64]*subscribers),
		decoders:    make(map[common.Address]ParseLogFunc),
		evmChainID:  evmChainID,
		logger:      logger,
	}
}

func (r *registrations) addSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract
	r.decoders[addr] = reg.opts.ParseLog

	if _, exists := r.subscribers[reg.opts.NumConfirmations]; !exists {
		r.subscribers[reg.opts.NumConfirmations] = newSubscribers(r.evmChainID)
	}

	needsResubscribe = r.subscribers[reg.opts.NumConfirmations].addSubscriber(reg)

	// increase the variable for highest number of confirmations among all subscribers,
	// if the new subscriber has a higher value
	if reg.opts.NumConfirmations > r.highestNumConfirmations {
		r.highestNumConfirmations = reg.opts.NumConfirmations
	}
	return
}

func (r *registrations) removeSubscriber(reg registration) (needsResubscribe bool) {
	subscribers, exists := r.subscribers[reg.opts.NumConfirmations]
	if !exists {
		return
	}

	needsResubscribe = subscribers.removeSubscriber(reg)

	if len(r.subscribers[reg.opts.NumConfirmations].handlers) == 0 {
		delete(r.subscribers, reg.opts.NumConfirmations)
		r.resetHighestNumConfirmationsValue()
	}
	return
}

// reset the number tracking highest num confirmations among all subscribers
func (r *registrations) resetHighestNumConfirmationsValue() {
	highestNumConfirmations := uint64(0)

	for numConfirmations := range r.subscribers {
		if numConfirmations > highestNumConfirmations {
			highestNumConfirmations = numConfirmations
		}
	}
	r.highestNumConfirmations = highestNumConfirmations
}

func (r *registrations) addressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for _, sub := range r.subscribers {
		add, t := sub.addressesAndTopics()
		addresses = append(addresses, add...)
		topics = append(topics, t...)
	}
	return addresses, topics
}

func (r *registrations) isAddressRegistered(address common.Address) bool {
	for _, sub := range r.subscribers {
		if sub.isAddressRegistered(address) {
			return true
		}
	}
	return false
}

func (r *registrations) sendLogs(logsToSend []logsOnBlock, latestHead eth.Head, broadcasts []LogBroadcast) {
	broadcastsExisting := make(map[LogBroadcastAsKey]struct{})
	for _, b := range broadcasts {

		broadcastsExisting[b.AsKey()] = struct{}{}
	}

	latestBlockNumber := uint64(latestHead.Number)

	for _, logsPerBlock := range logsToSend {
		for numConfirmations, subscribers := range r.subscribers {

			if numConfirmations != 0 && latestBlockNumber < numConfirmations {
				// Skipping send because the block is definitely too young
				continue
			}

			// We attempt the send multiple times per log
			// so here we need to see if this particular listener actually should receive it at this depth
			isOldEnough := numConfirmations == 0 || (logsPerBlock.BlockNumber+numConfirmations-1) <= latestBlockNumber
			if !isOldEnough {
				continue
			}

			for _, log := range logsPerBlock.Logs {
				subscribers.sendLog(log, latestHead, broadcastsExisting, r.decoders, r.logger)
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

func newSubscribers(evmChainID big.Int) *subscribers {
	return &subscribers{
		handlers:   make(map[common.Address]map[common.Hash]map[Listener]*listenerMetadata),
		evmChainID: evmChainID,
	}
}

func (r *subscribers) addSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract

	if reg.opts.NumConfirmations <= 0 {
		reg.opts.NumConfirmations = 1
	}

	if _, exists := r.handlers[addr]; !exists {
		r.handlers[addr] = make(map[common.Hash]map[Listener]*listenerMetadata)
	}

	for topic, topicValueFilters := range reg.opts.LogsWithTopics {
		if _, exists := r.handlers[addr][topic]; !exists {
			r.handlers[addr][topic] = make(map[Listener]*listenerMetadata)
			needsResubscribe = true
		}

		r.handlers[addr][topic][reg.listener] = &listenerMetadata{
			opts:    reg.opts,
			filters: topicValueFilters,
		}
	}
	return
}

func (r *subscribers) removeSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract

	if _, exists := r.handlers[addr]; !exists {
		return
	}
	for topic := range reg.opts.LogsWithTopics {
		topicMap, exists := r.handlers[addr][topic]
		if !exists {
			continue
		}

		delete(topicMap, reg.listener)

		if len(topicMap) == 0 {
			needsResubscribe = true
			delete(r.handlers[addr], topic)
		}
		if len(r.handlers[addr]) == 0 {
			delete(r.handlers, addr)
		}
	}
	return
}

func (r *subscribers) addressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for addr := range r.handlers {
		addresses = append(addresses, addr)
		for topic := range r.handlers[addr] {
			topics = append(topics, topic)
		}
	}
	return addresses, topics
}

func (r *subscribers) isAddressRegistered(address common.Address) bool {
	_, exists := r.handlers[address]
	return exists
}

func (r *subscribers) sendLog(log types.Log, latestHead eth.Head,
	broadcasts map[LogBroadcastAsKey]struct{},
	decoders map[common.Address]ParseLogFunc,
	logger logger.Logger) {

	latestBlockNumber := uint64(latestHead.Number)
	var wg sync.WaitGroup
	for listener, metadata := range r.handlers[log.Address][log.Topics[0]] {
		listener := listener

		currentBroadcast := NewLogBroadcastAsKey(log, listener)
		_, exists := broadcasts[currentBroadcast]
		if exists {
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
		if parseLog := decoders[log.Address]; parseLog != nil {
			decodedLog, err = parseLog(logCopy)
			if err != nil {
				logger.Errorw("Could not parse contract log", "error", err)
				continue
			}
		}

		logger.Debugw("LogBroadcaster: Sending out log",
			"blockNumber", log.BlockNumber, "blockHash", log.BlockHash,
			"address", log.Address, "latestBlockNumber", latestBlockNumber)

		wg.Add(1)
		go func() {
			defer wg.Done()
			listener.HandleLog(&broadcast{
				latestBlockNumber,
				latestHead.Hash,
				decodedLog,
				logCopy,
				listener.JobID(),
				r.evmChainID,
			})
		}()
	}
	wg.Wait()
}
