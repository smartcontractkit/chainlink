package s4

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Row represents a data row persisted by ORM.
type Row struct {
	Address    *utils.Big
	SlotId     uint
	Payload    []byte
	Version    uint64
	Expiration int64
	Confirmed  bool
	Signature  []byte
	UpdatedAt  int64
}

// VersionRow used by GetVersions function.
type VersionRow struct {
	Address *utils.Big
	SlotId  uint
	Version uint64
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

// ORM represents S4 persistence layer.
// All functions are thread-safe.
type ORM interface {
	// Get reads a row for the given address and slotId combination.
	// If such row does not exist, ErrNotFound is returned.
	// There is no filter on Expiration.
	Get(address *utils.Big, slotId uint, qopts ...pg.QOpt) (*Row, error)

	// Update inserts or updates the row identified by (Address, SlotId) pair.
	// When updating, the new row must have greater or equal version,
	// otherwise ErrVersionTooLow is returned.
	// UpdatedAt field value is ignored.
	Update(row *Row, qopts ...pg.QOpt) error

	// DeleteExpired deletes any entries having Expiration < now().
	DeleteExpired(qopts ...pg.QOpt) error

	// GetVersions selects all non-expired row versions for the given addresses range.
	// For the full address range, use NewFullAddressRange().
	GetVersions(addressRange *AddressRange, qopts ...pg.QOpt) ([]*VersionRow, error)

	// GetUnconfirmedRows selects all non-expired, non-confirmed rows ordered by UpdatedAt.
	// The number of returned rows is limited to the given limit.
	GetUnconfirmedRows(limit uint, qopts ...pg.QOpt) ([]*Row, error)
}

func (r Row) Clone() *Row {
	clone := Row{
		Address:    utils.NewBig(r.Address.ToInt()),
		SlotId:     r.SlotId,
		Payload:    make([]byte, len(r.Payload)),
		Version:    r.Version,
		Expiration: r.Expiration,
		Confirmed:  r.Confirmed,
		Signature:  make([]byte, len(r.Signature)),
		UpdatedAt:  r.UpdatedAt,
	}
	copy(clone.Payload, r.Payload)
	copy(clone.Signature, r.Signature)
	return &clone
}
