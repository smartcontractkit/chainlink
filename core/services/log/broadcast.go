package log

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Broadcast --output ./mocks/ --case=underscore --structname Broadcast --filename broadcast.go

type (
	// The Broadcast type wraps a models.Log but provides additional functionality
	// for determining whether or not the log has been consumed and for marking
	// the log as consumed
	Broadcast interface {
		DecodedLog() interface{}
		RawLog() types.Log
		SetDecodedLog(interface{})
		WasAlreadyConsumed() (bool, error)
		MarkConsumed() error
	}

	broadcast struct {
		orm        ORM
		decodedLog interface{}
		rawLog     types.Log
		jobID      models.JobID
		jobIDV2    int32
		isV2       bool
	}
)

func (b *broadcast) DecodedLog() interface{} {
	return b.decodedLog
}

func (b *broadcast) RawLog() types.Log {
	return b.rawLog
}

func (b *broadcast) SetDecodedLog(newLog interface{}) {
	b.decodedLog = newLog
}

// WasAlreadyConsumed reports whether the given consumer had already consumed the given log
func (b *broadcast) WasAlreadyConsumed() (bool, error) {
	return b.orm.WasBroadcastConsumed(b.rawLog.BlockHash, b.rawLog.Index, b.JobID())
}

// MarkConsumed marks the log as having been successfully consumed by the subscriber
func (b *broadcast) MarkConsumed() error {
	return b.orm.MarkBroadcastConsumed(b.rawLog.BlockHash, b.rawLog.Index, b.JobID())
}

func (b broadcast) JobID() interface{} {
	if b.isV2 {
		return b.jobIDV2
	}
	return b.jobID
}
