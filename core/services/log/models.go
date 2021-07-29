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
		String() string
		LatestBlockNumber() uint64
		LatestBlockHash() common.Hash
		JobID() JobIdSelect
	}

	broadcast struct {
		latestBlockNumber uint64
		latestBlockHash   common.Hash
		decodedLog        interface{}
		rawLog            types.Log
		jobID             JobIdSelect
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

func (b *broadcast) JobID() JobIdSelect {
	return b.jobID
}

func (b *broadcast) String() string {
	jobId := b.jobID.String()
	return fmt.Sprintf("Broadcast(JobID:%v,LogAddress:%v,Topics(%d):%v)", jobId, b.rawLog.Address, len(b.rawLog.Topics), b.rawLog.Topics)
}

func NewLogBroadcast(rawLog types.Log, decodedLog interface{}) Broadcast {
	return &broadcast{
		latestBlockNumber: 0,
		latestBlockHash:   common.Hash{},
		decodedLog:        decodedLog,
		rawLog:            rawLog,
		jobID:             NewJobIdV1(models.NilJobID),
	}
}

type JobIdSelect struct {
	JobIDV1 models.JobID
	JobIDV2 int32
	IsV2    bool
}

func NewJobIdV1(id models.JobID) JobIdSelect {
	return JobIdSelect{
		JobIDV1: id,
	}
}
func NewJobIdV2(id int32) JobIdSelect {
	return JobIdSelect{
		JobIDV2: id,
		IsV2:    true,
	}
}
func NewJobIdFromListener(listener Listener) JobIdSelect {
	if listener.IsV2Job() {
		return NewJobIdV2(listener.JobIDV2())
	}
	return NewJobIdV1(listener.JobID())
}

func (j JobIdSelect) String() string {
	jobId := j.JobIDV1.String()
	if j.IsV2 {
		jobId = fmt.Sprintf("%v", j.JobIDV2)
	}
	return jobId
}
