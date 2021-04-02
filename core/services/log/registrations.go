package log

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// How it works in general:
// 1. All received logs are kept in a data structure and deleted ONLY after they are too old for all subscribers
// 2. The logs are attempted to be sent after every new head arrival:
// 		Each stored log is then checked against every matched listener and is sent unless:
//    A) is too young for that listener
//    B) the corresponding block height is known to be already processed for that listener
// So in the normal case, each log will be only processed once, and then its corresponding head will be remembered, so it's not double-sent
// In case of a re-org, the stored lowestAllowedBlockNumber is reset, so the logs from that chain are considered unprocessed at first.
type (
	registrations struct {
		registrations map[common.Address]map[common.Hash]map[Listener]*listenerMetadata // contractAddress => logTopic => Listener
		decoders      map[common.Address]AbigenContract

		// highest 'num_confirmations' per all listeners, used to know when to delete old logs
		highestNumConfirmations uint64
	}

	// The Listener responds to log events through HandleLog, and contains setup/tear-down
	// callbacks in the On* functions.
	Listener interface {
		HandleLog(b Broadcast)
		JobID() models.JobID
		JobIDV2() int32
		IsV2Job() bool
	}

	// metadata structure maintained per listener, used to avoid double-sends of logs
	listenerMetadata struct {
		opts                     ListenerOpts
		lowestAllowedBlockNumber uint64
		lastSeenChain            *models.Head
	}

	// an update to listener metadata structure
	listenerMetadataUpdate struct {
		toUpdate                    *listenerMetadata
		newLowestAllowedBlockNumber uint64
	}
)

func newRegistrations() *registrations {
	return &registrations{
		registrations: make(map[common.Address]map[common.Hash]map[Listener]*listenerMetadata),
		decoders:      make(map[common.Address]AbigenContract),
	}
}

func (r *registrations) addSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract.Address()
	r.decoders[addr] = reg.opts.Contract

	if reg.opts.NumConfirmations <= 0 {
		reg.opts.NumConfirmations = 1
	}

	if _, exists := r.registrations[addr]; !exists {
		r.registrations[addr] = make(map[common.Hash]map[Listener]*listenerMetadata)
	}

	topics := make([]common.Hash, len(reg.opts.Logs))
	for i, log := range reg.opts.Logs {
		topic := log.Topic()
		topics[i] = topic

		if _, exists := r.registrations[addr][topic]; !exists {
			r.registrations[addr][topic] = make(map[Listener]*listenerMetadata)
			needsResubscribe = true
		}
		r.registrations[addr][topic][reg.listener] = &listenerMetadata{
			opts:                     reg.opts,
			lowestAllowedBlockNumber: uint64(0),
		}
	}

	r.maybeIncreaseHighestNumConfirmations(reg.opts.NumConfirmations)
	return
}

func (r *registrations) removeSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract.Address()

	if _, exists := r.registrations[addr]; !exists {
		return
	}
	for _, logType := range reg.opts.Logs {
		topic := logType.Topic()

		if _, exists := r.registrations[addr][topic]; !exists {
			continue
		}

		delete(r.registrations[addr][topic], reg.listener)

		if len(r.registrations[addr][topic]) == 0 {
			needsResubscribe = true
			delete(r.registrations[addr], topic)
		}
		if len(r.registrations[addr]) == 0 {
			delete(r.registrations, addr)
		}
	}

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

	for _, perAddress := range r.registrations {
		for _, perTopic := range perAddress {
			for _, listener := range perTopic {
				if listener.opts.NumConfirmations > highestNumConfirmations {
					highestNumConfirmations = listener.opts.NumConfirmations
				}
			}
		}
	}
	r.highestNumConfirmations = highestNumConfirmations
}

func (r *registrations) addressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for addr := range r.registrations {
		addresses = append(addresses, addr)
		for topic := range r.registrations[addr] {
			topics = append(topics, topic)
		}
	}
	return addresses, topics
}

func (r *registrations) isAddressRegistered(address common.Address) bool {
	_, exists := r.registrations[address]
	return exists
}

func (r *registrations) sendLogs(logs []types.Log, orm ORM, latestHead *models.Head) {
	updates := make([]listenerMetadataUpdate, 0)
	for _, log := range logs {
		logger.Tracef("Sending log at block num: %v", log.BlockNumber)
		r.sendLog(log, orm, latestHead, &updates)
	}
	applyListenerInfoUpdates(updates, latestHead)
}

func (r *registrations) sendLog(log types.Log, orm ORM, latestHead *models.Head, updates *[]listenerMetadataUpdate) {
	latestBlockNumber := uint64(latestHead.Number)
	var wg sync.WaitGroup
	for listener, metadata := range r.registrations[log.Address][log.Topics[0]] {
		listener := listener
		numConfirmations := metadata.opts.NumConfirmations

		if latestBlockNumber < numConfirmations {
			logger.Tracef("Skipping send because not enough height to send: %v - num confirmations: %v", latestBlockNumber, numConfirmations)
			continue
		}

		// we attempt the send multiple times per log (depending on distinct num of confirmations of listeners),
		// even if the logs are too young
		// so here we need to see if this particular listener actually should receive it at this depth
		isOldEnough := (log.BlockNumber + numConfirmations - 1) <= latestBlockNumber
		if !isOldEnough {
			continue
		}

		// all logs for blocks below lowestAllowedBlockNumber were already sent to this listener, so we skip them
		if log.BlockNumber < metadata.lowestAllowedBlockNumber && metadata.lastSeenChain != nil && metadata.lastSeenChain.IsInChain(log.BlockHash) {
			logger.Tracef("Skipping send because height %v is below lowest unprocessed: %v in current chain (ending at %v - %v)",
				log.BlockNumber, metadata.lowestAllowedBlockNumber, metadata.lastSeenChain.Number, metadata.lastSeenChain.Hash)
			continue
		} else {
			logger.Tracef("height %v, lowest unprocessed: %v in current chain (chain %v)",
				log.BlockNumber, metadata.lowestAllowedBlockNumber, metadata.lastSeenChain)
		}

		// make sure that this log is not sent again on the next head by increasing the newLowestAllowedBlockNumber
		*updates = append(*updates, listenerMetadataUpdate{
			toUpdate:                    metadata,
			newLowestAllowedBlockNumber: log.BlockNumber + 1,
		})

		logCopy := copyLog(log)
		decodedLog, err := r.decoders[log.Address].ParseLog(logCopy)
		if err != nil {
			logger.Errorw("Could not parse contract log", "error", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			listener.HandleLog(&broadcast{
				orm:               orm,
				latestBlockNumber: latestBlockNumber,
				latestBlockHash:   latestHead.Hash,
				rawLog:            logCopy,
				decodedLog:        decodedLog,
				jobID:             listener.JobID(),
				jobIDV2:           listener.JobIDV2(),
				isV2:              listener.IsV2Job(),
			})
		}()
	}
	wg.Wait()
}

/**
	After processing the logs in this batch, the listenerMetadata structures that we touched, are updated
  with new information about the canonical chain and the lowestAllowedBlockNumber value (higher every time) that is used to guard against double-sends
  Note that the updates are applied only after all the logs for the (latest height - num_confirmations) head height were sent.
*/
func applyListenerInfoUpdates(updates []listenerMetadataUpdate, latestHead *models.Head) {
	for _, update := range updates {
		if update.toUpdate.lastSeenChain == nil || latestHead.IsInChain(update.toUpdate.lastSeenChain.Hash) {
			if update.toUpdate.lastSeenChain == nil {
				logger.Tracef("No chain saved for listener on address (%v), confirmations: %v", update.toUpdate.opts.Contract.Address(), update.toUpdate.opts.NumConfirmations)
			}
			if update.toUpdate.lowestAllowedBlockNumber < update.newLowestAllowedBlockNumber {
				update.toUpdate.lowestAllowedBlockNumber = update.newLowestAllowedBlockNumber
			}
		} else {
			logger.Debugf("Chain reorg on height %v, hash: %v is not in %v (%v)", latestHead.Number, latestHead.Hash, update.toUpdate.lastSeenChain.Number, update.toUpdate.lastSeenChain.Hash)
			// re-org situation: the chain was changed, so we can't use the number that tracked last unprocessed height of the previous chain
			update.toUpdate.lowestAllowedBlockNumber = 0
		}
		logger.Tracef("Setting (%v %v) as latest head", latestHead.Number, latestHead.Hash)

		update.toUpdate.lastSeenChain = latestHead
	}
}

func copyLog(l types.Log) types.Log {
	var cpy types.Log
	cpy.Address = l.Address
	if l.Topics != nil {
		cpy.Topics = make([]common.Hash, len(l.Topics))
		copy(cpy.Topics, l.Topics)
	}
	if l.Data != nil {
		cpy.Data = make([]byte, len(l.Data))
		copy(cpy.Data, l.Data)
	}
	cpy.BlockNumber = l.BlockNumber
	cpy.TxHash = l.TxHash
	cpy.TxIndex = l.TxIndex
	cpy.BlockHash = l.BlockHash
	cpy.Index = l.Index
	cpy.Removed = l.Removed
	return cpy
}
