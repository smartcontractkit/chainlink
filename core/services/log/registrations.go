package log

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// How it works in general:
// 1. Each listener being registered can specify a custom NumConfirmations - number of block confirmations required for any log being sent to it.
// 2. Adding and removing listeners updates the highestNumConfirmations - a number tracking what's the current highest NumConfirmations globally
//
// 3. All received logs are kept in an array and deleted ONLY after they are outside the confirmation range for all subscribers
// (when given log height is lower than (latest height - max(highestNumConfirmations, ETH_FINALITY_DEPTH)) ) -> see: pool.go
//
// 4. The logs are attempted to be sent after every new head arrival:
// 		Each stored log is then checked against every matched listener and is sent unless:
//    A) is too young for that listener
//    B) the corresponding block height is known to be already processed for that listener
// In the normal case, each log will be only processed once, and then its corresponding head will be remembered, so it's not double-sent
//
// After processing the whole batch of logs considered for sending, the per-listener metadata is updated in applyListenerInfoUpdates.
// If a re-org happens, the stored lowestAllowedBlockNumber (per-listener) is re-set,
// so the logs from that chain are then considered unprocessed, and will be sent again.
//
type (
	registrations struct {
		subscribers map[uint64]*subscribers
		decoders    map[common.Address]ParseLogFunc

		// highest 'NumConfirmations' per all listeners, used to decide about deleting older logs if it's higher than EthFinalityDepth
		// it's: max(listeners.map(l => l.num_confirmations)
		highestNumConfirmations uint64
	}

	subscribers struct {
		handlers map[common.Address]map[common.Hash]map[Listener]*listenerMetadata // contractAddress => logTopic => Listener
	}

	// The Listener responds to log events through HandleLog.
	Listener interface {
		HandleLog(b Broadcast)
		JobID() models.JobID
		JobIDV2() int32
		IsV2Job() bool
	}

	// metadata structure maintained per listener, used to avoid double-sends of logs
	listenerMetadata struct {
		opts    ListenerOpts
		filters [][]Topic
	}
)

func newRegistrations() *registrations {
	return &registrations{
		subscribers: make(map[uint64]*subscribers),
		decoders:    make(map[common.Address]ParseLogFunc),
	}
}

func (r *registrations) addSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract
	r.decoders[addr] = reg.opts.ParseLog

	if reg.opts.NumConfirmations <= 0 {
		reg.opts.NumConfirmations = 1
	}

	if _, exists := r.subscribers[reg.opts.NumConfirmations]; !exists {
		r.subscribers[reg.opts.NumConfirmations] = newSubscribers()
	}

	needsResubscribe = r.subscribers[reg.opts.NumConfirmations].addSubscriber(reg)

	r.maybeIncreaseHighestNumConfirmations(reg.opts.NumConfirmations)
	return
}

func (r *registrations) removeSubscriber(reg registration) (needsResubscribe bool) {
	l, exists := r.subscribers[reg.opts.NumConfirmations]
	if !exists {
		return
	}

	needsResubscribe = l.removeSubscriber(reg)
	r.resetHighestNumConfirmations()
	return
}

// increase the highestNumConfirmations stored if the new listener has a higher value
func (r *registrations) maybeIncreaseHighestNumConfirmations(newNumConfirmations uint64) {
	if newNumConfirmations > r.highestNumConfirmations {
		r.highestNumConfirmations = newNumConfirmations
	}
}

// reset the highest confirmation number per all current listeners
func (r *registrations) resetHighestNumConfirmations() {
	highestNumConfirmations := uint64(0)
	for numConf := range r.subscribers {
		if numConf > highestNumConfirmations {
			highestNumConfirmations = numConf
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

func (r *registrations) sendLogs(logsToSend []logsOnBlock, latestHead models.Head, broadcasts []LogBroadcast) {
	broadcastsExisting := make(map[LogBroadcastAsKey]struct{})
	for _, b := range broadcasts {

		broadcastsExisting[b.AsKey()] = struct{}{}
	}

	latestBlockNumber := uint64(latestHead.Number)

	for _, logsPerBlock := range logsToSend {
		for numConfirmations, subscribers := range r.subscribers {
			if latestBlockNumber < numConfirmations {
				// Skipping send because not enough height to send
				continue
			}
			// We attempt the send multiple times per log (depending on distinct num of confirmations of listeners),
			// even if the logs are too young
			// so here we need to see if this particular listener actually should receive it at this depth
			isOldEnough := (logsPerBlock.BlockNumber + numConfirmations - 1) <= latestBlockNumber
			if !isOldEnough {
				continue
			}

			for _, log := range logsPerBlock.Logs {
				subscribers.sendLog(log, latestHead, broadcastsExisting, r.decoders)
			}
		}
	}
}

// Returns true if there is at least one filter value (or no filters) that matches an actual received value for every index i, or false otherwise
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

func newSubscribers() *subscribers {
	return &subscribers{
		handlers: make(map[common.Address]map[common.Hash]map[Listener]*listenerMetadata),
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
		if _, exists := r.handlers[addr][topic]; !exists {
			continue
		}

		delete(r.handlers[addr][topic], reg.listener)

		if len(r.handlers[addr][topic]) == 0 {
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

func (r *subscribers) sendLog(log types.Log, latestHead models.Head, broadcasts map[LogBroadcastAsKey]struct{}, decoders map[common.Address]ParseLogFunc) {
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

		logCopy := gethwrappers.CopyLog(log)

		var decodedLog generated.AbigenLog
		var err error
		if parseLog := decoders[log.Address]; parseLog != nil {
			decodedLog, err = parseLog(logCopy)
			if err != nil {
				logger.Errorw("Could not parse contract log", "error", err)
				continue
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			listener.HandleLog(&broadcast{
				latestBlockNumber: latestBlockNumber,
				latestBlockHash:   latestHead.Hash,
				rawLog:            logCopy,
				decodedLog:        decodedLog,
				jobID:             NewJobIdFromListener(listener),
			})
		}()
	}
	wg.Wait()
}
