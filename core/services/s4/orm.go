package s4

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	ErrEntryNotFound = errors.New("entry not found")
)

// Entry represents a data row persisted by ORM.
type Entry struct {
	Payload           []byte
	Version           int64
	Expiration        int64
	Confirmed         bool
	HighestExpiration int64
	Signature         []byte
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

// ORM represents S4 persistence layer.
// All functions are thread-safe.
type ORM interface {
	// Get reads an entry for the given address/slotId combination.
	// If an entry is missing, ErrEntryNotFound is returned.
	// Returned entry is a clone, safe to modify.
	Get(address common.Address, slotId int, qopts ...pg.QOpt) (*Entry, error)

	// Put inserts (or updates) an entry identified by the specified address and slotId.
	// No validation is applied for signature, version, etc.
	// Implementation clones the given Entry when persisting.
	Upsert(address common.Address, slotId int, entry *Entry, qopts ...pg.QOpt) error

	// DeleteExpired deletes any entries having HighestExpiration < now().
	// The function can be called by OCR plugin on every round.
	// (shall be cheap for postgres implementation, given the right columns are indexed).
	DeleteExpired(qopts ...pg.QOpt) error
}

func (e Entry) Clone() *Entry {
	clone := Entry{
		Payload:           make([]byte, len(e.Payload)),
		Version:           e.Version,
		Expiration:        e.Expiration,
		Confirmed:         e.Confirmed,
		HighestExpiration: e.HighestExpiration,
		Signature:         make([]byte, len(e.Signature)),
	}
	copy(clone.Payload, e.Payload)
	copy(clone.Signature, e.Signature)
	return &clone
}
