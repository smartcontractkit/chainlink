package s4

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// Row represents a data row persisted by ORM.
type Row struct {
	Address    *big.Big
	SlotId     uint
	Payload    []byte
	Version    uint64
	Expiration int64
	Confirmed  bool
	Signature  []byte
}

// SnapshotRow(s) are returned by GetSnapshot function.
type SnapshotRow struct {
	Address     *big.Big
	SlotId      uint
	Version     uint64
	Expiration  int64
	Confirmed   bool
	PayloadSize uint64
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

// ORM represents S4 persistence layer.
// All functions are thread-safe.
type ORM interface {
	// Get reads a row for the given address and slotId combination.
	// If such row does not exist, ErrNotFound is returned.
	// There is no filter on Expiration.
	Get(ctx context.Context, address *big.Big, slotId uint) (*Row, error)

	// Update inserts or updates the row identified by (Address, SlotId) pair.
	// When updating, the new row must have greater or equal version,
	// otherwise ErrVersionTooLow is returned.
	// UpdatedAt field value is ignored.
	Update(ctx context.Context, row *Row) error

	// DeleteExpired deletes any entries having Expiration < utcNow,
	// up to the given limit.
	// Returns the number of deleted rows.
	DeleteExpired(ctx context.Context, limit uint, utcNow time.Time) (int64, error)

	// GetSnapshot selects all non-expired row versions for the given addresses range.
	// For the full address range, use NewFullAddressRange().
	GetSnapshot(ctx context.Context, addressRange *AddressRange) ([]*SnapshotRow, error)

	// GetUnconfirmedRows selects all non-expired, non-confirmed rows ordered by UpdatedAt.
	// The number of returned rows is limited to the given limit.
	GetUnconfirmedRows(ctx context.Context, limit uint) ([]*Row, error)
}

func (r Row) Clone() *Row {
	clone := Row{
		Address:    big.New(r.Address.ToInt()),
		SlotId:     r.SlotId,
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
