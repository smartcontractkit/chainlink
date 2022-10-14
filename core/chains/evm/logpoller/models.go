package logpoller

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// LogPollerBlock represents an unfinalized block
// used for reorg detection when polling.
type LogPollerBlock struct {
	EvmChainId *utils.Big
	BlockHash  common.Hash
	// Note geth uses int64 internally https://github.com/ethereum/go-ethereum/blob/f66f1a16b3c480d3a43ac7e8a09ab3e362e96ae4/eth/filters/api.go#L340
	BlockNumber int64
	CreatedAt   time.Time
}

// Log represents an EVM log.
type Log struct {
	EvmChainId  *utils.Big
	LogIndex    int64
	BlockHash   common.Hash
	BlockNumber int64
	Topics      pq.ByteaArray
	EventSig    common.Hash
	Address     common.Address
	TxHash      common.Hash
	Data        []byte
	CreatedAt   time.Time
}

func (l *Log) GetTopics() []common.Hash {
	var tps []common.Hash
	for _, topic := range l.Topics {
		tps = append(tps, common.BytesToHash(topic))
	}
	return tps
}

func (l *Log) ToGethLog() types.Log {
	return types.Log{
		Data:        l.Data,
		Address:     l.Address,
		BlockHash:   l.BlockHash,
		BlockNumber: uint64(l.BlockNumber),
		Topics:      l.GetTopics(),
		TxHash:      l.TxHash,
		Index:       uint(l.LogIndex),
	}
}
