package pool

import (
	"math"

	"github.com/smartcontractkit/libocr/commontypes"
)

type Pool[T any] struct {
	maxItemsPerSender int

	completedSeqNr uint64
	entries        map[uint64]map[commontypes.OracleID]*Entry[T]
	count          map[commontypes.OracleID]int
}

func NewPool[T any](maxItemsPerSender int) *Pool[T] {
	return &Pool[T]{
		maxItemsPerSender,

		0,
		map[uint64]map[commontypes.OracleID]*Entry[T]{},
		map[commontypes.OracleID]int{},
	}
}

type Entry[T any] struct {
	Item     T
	Verified *bool
}

func (p *Pool[M]) ReapCompleted(completedSeqNr uint64) {
	for seqNr, messagesByOracleID := range p.entries {
		if seqNr > completedSeqNr {
			continue
		}
		for sender := range messagesByOracleID {
			p.count[sender] -= 1
			delete(messagesByOracleID, sender)
		}
		delete(p.entries, seqNr)
	}
	p.completedSeqNr = completedSeqNr
}

type PutResult string

const (
	PutResultOK               PutResult = "ok"
	PutResultDuplicate        PutResult = "duplicate"
	PutResultFull             PutResult = "pool is full for sender"
	PutResultAlreadyCompleted PutResult = "seqNr too low"
)

func (p *Pool[T]) Put(seqNr uint64, sender commontypes.OracleID, item T) PutResult {
	if seqNr <= p.completedSeqNr {
		return PutResultAlreadyCompleted
	}
	if p.maxItemsPerSender <= p.count[sender] {
		return PutResultFull
	}
	if p.entries[seqNr] == nil {
		p.entries[seqNr] = map[commontypes.OracleID]*Entry[T]{}
	}
	if p.entries[seqNr][sender] != nil {
		return PutResultDuplicate
	}
	p.entries[seqNr][sender] = &Entry[T]{item, nil}
	p.count[sender]++
	return PutResultOK
}

func (p *Pool[M]) StoreVerified(seqNr uint64, sender commontypes.OracleID, verified bool) {
	if p.entries[seqNr] == nil {
		return
	}
	if p.entries[seqNr][sender] == nil {
		return
	}
	if p.entries[seqNr][sender].Verified != nil {
		// only store first verification result
		return
	}
	p.entries[seqNr][sender].Verified = &verified
}

func (p *Pool[M]) Entries(seqNr uint64) map[commontypes.OracleID]*Entry[M] {
	return p.entries[seqNr]
}

func (p *Pool[M]) EntriesWithMinSeqNr() map[commontypes.OracleID]*Entry[M] {
	if len(p.entries) == 0 {
		return nil
	}
	minSeqNr := uint64(math.MaxUint64)
	for seqNr := range p.entries {
		if seqNr < minSeqNr {
			minSeqNr = seqNr
		}
	}
	return p.Entries(minSeqNr)
}
