package log

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Broadcast --output ./mocks/ --case=underscore --structname Broadcast --filename broadcast.go

type (
	// The Broadcast type wraps a types.Log but provides additional functionality
	// for determining whether or not the log has been consumed and for marking
	// the log as consumed
	Broadcast interface {
		DecodedLog() interface{}
		RawLog() types.Log
		SetDecodedLog(interface{})
		String() string
		LatestBlockNumber() uint64
		LatestBlockHash() common.Hash
		JobID() interface{}
	}

	broadcast struct {
		latestBlockNumber uint64
		latestBlockHash   common.Hash
		decodedLog        interface{}
		rawLog            types.Log
		jobID             models.JobID
		jobIDV2           int32
		isV2              bool
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

func (b *broadcast) RawLog() types.Log {
	return b.rawLog
}

func (b *broadcast) SetDecodedLog(newLog interface{}) {
	b.decodedLog = newLog
}

func (b broadcast) JobID() interface{} {
	if b.isV2 {
		return b.jobIDV2
	}
	return b.jobID
}

func (b *broadcast) String() string {
	jobId := b.jobID.String()
	if b.isV2 {
		jobId = fmt.Sprintf("%v", b.jobIDV2)
	}
	return fmt.Sprintf("Broadcast(JobID:%v,LogAddress:%v,Topics(%d):%v)", jobId, b.rawLog.Address, len(b.rawLog.Topics), b.rawLog.Topics)
}

func NewLogBroadcast() Broadcast {
	return &broadcast{
		latestBlockNumber: 0,
		latestBlockHash:   common.Hash{},
		decodedLog:        nil,
		rawLog:            types.Log{},
		jobID:             models.JobID{},
		jobIDV2:           0,
		isV2:              false,
	}
}
