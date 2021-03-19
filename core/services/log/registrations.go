package log

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	registrations struct {
		registrations           map[common.Address]map[common.Hash]map[Listener]*listenerInfo // contractAddress => logTopic => Listener
		decoders                map[common.Address]AbigenContract
		highestNumConfirmations uint64
	}

	// The Listener responds to log events through HandleLog, and contains setup/tear-down
	// callbacks in the On* functions.
	Listener interface {
		OnConnect()
		OnDisconnect()
		HandleLog(b Broadcast)
		JobID() models.JobID
		JobIDV2() int32
		IsV2Job() bool
	}

	listenerInfo struct {
		opts                     ListenerOpts
		lowestAllowedBlockNumber uint64
		lastSeenChain            *models.Head
	}

	listenerInfoUpdate struct {
		toUpdate                    *listenerInfo
		newLowestAllowedBlockNumber uint64
	}
)

func newRegistrations() *registrations {
	return &registrations{
		registrations: make(map[common.Address]map[common.Hash]map[Listener]*listenerInfo),
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
		r.registrations[addr] = make(map[common.Hash]map[Listener]*listenerInfo)
	}

	topics := make([]common.Hash, len(reg.opts.Logs))
	for i, log := range reg.opts.Logs {
		topic := log.Topic()
		topics[i] = topic

		// is it ok that there can be only one listener per contract address and topic combination?
		if _, exists := r.registrations[addr][topic]; !exists {
			r.registrations[addr][topic] = make(map[Listener]*listenerInfo)
			needsResubscribe = true
		}
		r.registrations[addr][topic][reg.listener] = &listenerInfo{
			opts:                     reg.opts,
			lowestAllowedBlockNumber: uint64(0),
		}
	}

	if reg.opts.NumConfirmations > r.highestNumConfirmations {
		r.highestNumConfirmations = reg.opts.NumConfirmations
	}
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

	highestNumConfirmations := uint64(0)
	// reset the highest confirmation number stored.
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
	return
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

func (r *registrations) sendLogs(logs []models.Log, orm ORM, latestHead *models.Head) {
	updates := make([]listenerInfoUpdate, 0)
	for _, log := range logs {
		logger.Tracef("Sending log at block num: %v", log.BlockNumber)
		r.sendLog(log, orm, latestHead, &updates)
	}

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

func (r *registrations) sendLog(log models.Log, orm ORM, latestHead *models.Head, updates *[]listenerInfoUpdate) {
	latestBlockNumber := uint64(latestHead.Number)
	var wg sync.WaitGroup
	for listener, info := range r.registrations[log.Address][log.Topics[0]] {
		listener := listener
		numConfirmations := info.opts.NumConfirmations

		if latestBlockNumber < numConfirmations {
			logger.Tracef("Skipping send because not enough height to send: %v - num confirmations: %v", latestBlockNumber, numConfirmations)
			continue
		}

		// we attempt the send multiple times per log (depending on distinct num of confirmations of listeners),
		// so here we need to see if this particular listener actually should receive it at this depth
		isOldEnough := (log.BlockNumber + numConfirmations - 1) <= latestBlockNumber
		if !isOldEnough {
			continue
		}

		// all logs for blocks below lowestAllowedBlockNumber were already sent to this listener, so we skip them
		if log.BlockNumber < info.lowestAllowedBlockNumber && info.lastSeenChain != nil && info.lastSeenChain.IsInChain(log.BlockHash) {
			logger.Tracef("Skipping send because height %v is below lowest unprocessed: %v in current chain (ending at %v - %v)",
				log.BlockNumber, info.lowestAllowedBlockNumber, info.lastSeenChain.Number, info.lastSeenChain.Hash)
			continue
		} else {
			logger.Tracef("height %v, lowest unprocessed: %v in current chain (chain %v)",
				log.BlockNumber, info.lowestAllowedBlockNumber, info.lastSeenChain)
		}
		*updates = append(*updates, listenerInfoUpdate{info, log.BlockNumber + 1})

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
