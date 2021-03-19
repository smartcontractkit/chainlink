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
		registrations            map[common.Address]map[common.Hash]map[Listener]ListenerOpts // contractAddress => logTopic => Listener
		decoders                 map[common.Address]AbigenContract
		confirmationFilter       confFilter
		smallestNumConfirmations uint32
		largestNumConfirmations  uint32
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

	confFilter struct {
		confirmationNumbers      map[confirmationInfo]struct{}
		smallestNumConfirmations uint32
		largestNumConfirmations  uint32
	}

	confirmationInfo struct {
		address          common.Address
		numConfirmations uint32
	}
)

func newConfFilter() confFilter {
	return confFilter{
		confirmationNumbers:      make(map[confirmationInfo]struct{}, 0),
		smallestNumConfirmations: 0,
		largestNumConfirmations:  0,
	}
}

func (c confFilter) includeNumConfirmations(address common.Address, numConfirmations uint32) {
	info := confirmationInfo{
		address,
		numConfirmations,
	}
	c.confirmationNumbers[info] = struct{}{}
	if numConfirmations < c.smallestNumConfirmations {
		c.smallestNumConfirmations = numConfirmations
	}
	if numConfirmations > c.largestNumConfirmations {
		c.largestNumConfirmations = numConfirmations
	}
}

func (c confFilter) removeNumConfirmations(address common.Address, numConfirmations uint32) {
	//info := confirmationInfo{
	//	address,
	//	numConfirmations,
	//}
	//delete(c.confirmationNumbers[info])

	//c.confirmationNumbers[info] = struct{}{}
	//if numConfirmations < c.smallestNumConfirmations {
	//	c.smallestNumConfirmations = numConfirmations
	//}
	//if numConfirmations > c.largestNumConfirmations {
	//	c.largestNumConfirmations = numConfirmations
	//}
}

//func (c confFilter) isConfirmedFor(numConfirmations uint32) bool {
//	_, exists := c.confirmationNumbers[numConfirmations]
//	return exists
//}

func newRegistrations() registrations {
	return registrations{
		registrations:      make(map[common.Address]map[common.Hash]map[Listener]ListenerOpts),
		decoders:           make(map[common.Address]AbigenContract),
		confirmationFilter: newConfFilter(),
	}
}

func (r registrations) addSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract.Address()
	r.decoders[addr] = reg.opts.Contract
	//r.confirmationFilter.includeNumConfirmations(addr, reg.opts.NumConfirmations)

	//if reg.opts.NumConfirmations < r.smallestNumConfirmations {
	//	c.smallestNumConfirmations = numConfirmations
	//}
	//if numConfirmations > c.largestNumConfirmations {
	//	c.largestNumConfirmations = numConfirmations
	//}

	if _, exists := r.registrations[addr]; !exists {
		r.registrations[addr] = make(map[common.Hash]map[Listener]ListenerOpts)
	}

	topics := make([]common.Hash, len(reg.opts.Logs))
	for i, log := range reg.opts.Logs {
		topic := log.Topic()
		topics[i] = topic

		if _, exists := r.registrations[addr][topic]; !exists {
			r.registrations[addr][topic] = make(map[Listener]ListenerOpts)
			needsResubscribe = true
		}
		r.registrations[addr][topic][reg.listener] = reg.opts
	}
	return
}

func (r registrations) removeSubscriber(reg registration) (needsResubscribe bool) {
	addr := reg.opts.Contract.Address()
	//r.confirmationFilter.removeNumConfirmations(addr, reg.opts.NumConfirmations)

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
	return
}

func (r registrations) requiredConfirmations() ([]common.Address, []common.Hash) {
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

func (r registrations) addressesAndTopics() ([]common.Address, []common.Hash) {
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

func (r registrations) isAddressRegistered(address common.Address) bool {
	_, exists := r.registrations[address]
	return exists
}

func (r registrations) getDistinctConfirmationDepths() map[uint64]struct{} {
	depths := make(map[uint64]struct{}, 0)
	for _, perAddress := range r.registrations {
		for _, perTopic := range perAddress {
			for _, opts := range perTopic {
				depths[opts.NumConfirmations] = struct{}{}
			}
		}
	}
	return depths
}

func (r registrations) sendLog(log models.Log, orm ORM, latestBlockNumber int64) {
	var wg sync.WaitGroup
	for listener, opts := range r.registrations[log.Address][log.Topics[0]] {
		listener := listener

		// we attempt the send multiple times per log (depending on distinct num of confirmations of listeners),
		// so here we need to see if this particular listener actually should receive it at this depth
		requiredDepth := uint64(latestBlockNumber) - opts.NumConfirmations
		if log.BlockNumber != requiredDepth {
			continue
		}

		logCopy := copyLog(log)
		decodedLog, err := r.decoders[log.Address].ParseLog(logCopy)
		if err != nil {
			logger.Errorw("could not parse contract log", "error", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			listener.HandleLog(&broadcast{
				orm:        orm,
				rawLog:     logCopy,
				decodedLog: decodedLog,
				jobID:      listener.JobID(),
				jobIDV2:    listener.JobIDV2(),
				isV2:       listener.IsV2Job(),
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
