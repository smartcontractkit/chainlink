package log

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

type (
	// The Broadcast type wraps a types.Log but provides additional functionality
	// for determining whether or not the log has been consumed and for marking
	// the log as consumed
	Broadcast interface {
		DecodedLog() interface{}
		RawLog() types.Log
		String() string
		LatestBlockNumber() uint64
		LatestBlockHash() common.Hash
		ReceiptsRoot() common.Hash
		TransactionsRoot() common.Hash
		StateRoot() common.Hash
		JobID() int32
		EVMChainID() big.Int
	}

	broadcast struct {
		latestBlockNumber uint64
		latestBlockHash   common.Hash
		receiptsRoot      common.Hash
		transactionsRoot  common.Hash
		stateRoot         common.Hash
		decodedLog        interface{}
		rawLog            types.Log
		jobID             int32
		evmChainID        big.Int
	}
)

func (b *broadcast) DecodedLog() interface{} {
	return b.decodedLog
}

func (b *broadcast) LatestBlockNumber() uint64 {
	return b.latestBlockNumber
}

func (b *broadcast) LatestBlockHash() common.Hash {
	return b.latestBlockHash
}

func (b *broadcast) ReceiptsRoot() common.Hash {
	return b.receiptsRoot
}

func (b *broadcast) TransactionsRoot() common.Hash {
	return b.transactionsRoot
}

func (b *broadcast) StateRoot() common.Hash {
	return b.stateRoot
}

func (b *broadcast) RawLog() types.Log {
	return b.rawLog
}

func (b *broadcast) SetDecodedLog(newLog interface{}) {
	b.decodedLog = newLog
}

func (b *broadcast) JobID() int32 {
	return b.jobID
}

func (b *broadcast) String() string {
	return fmt.Sprintf("Broadcast(JobID:%v,LogAddress:%v,Topics(%d):%v)", b.jobID, b.rawLog.Address, len(b.rawLog.Topics), b.rawLog.Topics)
}

func (b *broadcast) EVMChainID() big.Int {
	return b.evmChainID
}

func NewLogBroadcast(rawLog types.Log, evmChainID big.Int, decodedLog interface{}) Broadcast {
	return &broadcast{
		latestBlockNumber: 0,
		latestBlockHash:   common.Hash{},
		receiptsRoot:      common.Hash{},
		transactionsRoot:  common.Hash{},
		stateRoot:         common.Hash{},
		decodedLog:        decodedLog,
		rawLog:            rawLog,
		jobID:             0,
		evmChainID:        evmChainID,
	}
}

type AbigenContract interface {
	Address() common.Address
	ParseLog(log types.Log) (generated.AbigenLog, error)
}
