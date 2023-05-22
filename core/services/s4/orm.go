package s4

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// Row represents a data row persisted by ORM.
type Row struct {
	Payload    []byte
	Version    uint64
	Expiration int64
	Confirmed  bool
	Signature  []byte
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

// ORM represents S4 persistence layer.
// All functions are thread-safe.
type ORM interface {
	// Get reads a row for the given address/slotId combination.
	// If a row is missing, ErrNotFound is returned.
	// Returned row is a clone, safe to modify.
	Get(address common.Address, slotId uint, qopts ...pg.QOpt) (*Row, error)

	// Put inserts (or updates) a row identified by the specified address and slotId.
	// No validation is applied for signature, version, etc.
	// Implementation clones the given Row when persisting.
	Upsert(address common.Address, slotId uint, row *Row, qopts ...pg.QOpt) error

	// DeleteExpired deletes any entries having HighestExpiration < now().
	// The function can be called by OCR plugin on every round.
	// (shall be cheap for postgres implementation, given the right columns are indexed).
	DeleteExpired(qopts ...pg.QOpt) error
}

func (r Row) Clone() *Row {
	clone := Row{
		Payload:    make([]byte, len(r.Payload)),
		Version:    r.Version,
		Expiration: r.Expiration,
		Confirmed:  r.Confirmed,
		Signature:  make([]byte, len(r.Signature)),
	}
	copy(clone.Payload, r.Payload)
	copy(clone.Signature, r.Signature)
	return &clone
}
