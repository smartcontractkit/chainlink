package logprovider

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

type upkeepFilter struct {
	addr []byte
	// selector is the filter selector in log trigger config
	selector uint8
	topics   []common.Hash
	upkeepID *big.Int
	// configUpdateBlock is the block number the filter was last updated at
	configUpdateBlock uint64
	// lastPollBlock is the last block number the logs were fetched for this upkeep
	// used by log event provider.
	lastPollBlock int64
	// lastRePollBlock is the last block number the logs were recovered for this upkeep
	// used by log recoverer.
	lastRePollBlock int64
}

func (f upkeepFilter) Clone() upkeepFilter {
	topics := make([]common.Hash, len(f.topics))
	copy(topics, f.topics)
	addr := make([]byte, len(f.addr))
	copy(addr, f.addr)
	return upkeepFilter{
		upkeepID:          f.upkeepID,
		selector:          f.selector,
		topics:            topics,
		addr:              addr,
		configUpdateBlock: f.configUpdateBlock,
		lastPollBlock:     f.lastPollBlock,
		lastRePollBlock:   f.lastRePollBlock,
	}
}

// Select returns a slice of logs which match the upkeep filter.
func (f upkeepFilter) Select(logs ...logpoller.Log) []logpoller.Log {
	var selected []logpoller.Log
	for _, log := range logs {
		if f.match(log) {
			selected = append(selected, log)
		}
	}
	return selected
}

// match returns a bool indicating if the log's topics data matches selector and indexed topics in upkeep filter.
func (f upkeepFilter) match(log logpoller.Log) bool {
	filters := f.topics[1:]
	selector := f.selector

	if selector == 0 {
		// no filters
		return true
	}

	for i, filter := range filters {
		// bitwise AND the selector with the index to check
		// if the filter is needed
		mask := uint8(1 << uint8(i))
		if selector&mask == uint8(0) {
			continue
		}
		if len(log.Topics) <= i+1 {
			// log doesn't have enough topics
			return false
		}
		if !bytes.Equal(filter.Bytes(), log.Topics[i+1]) {
			return false
		}
	}
	return true
}
