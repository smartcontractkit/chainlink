package s4

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// Row represents a data row persisted by ORM.
type Row struct {
	Address    string
	SlotId     uint
	Payload    []byte
	Version    uint64
	Expiration int64
	Confirmed  bool
	Signature  []byte
	UpdatedAt  int64
}

var (
	MinAddress *big.Int = common.HexToAddress("0x0").Big()
	MaxAddress *big.Int = common.HexToAddress("0x" + strings.Repeat("ff", common.AddressLength)).Big()
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

// ORM represents S4 persistence layer.
// All functions are thread-safe.
type ORM interface {
	// Get reads a row for the given address and slotId combination.
	// If such row does not exist, ErrNotFound is returned.
	// There is no filter on Expiration.
	Get(address common.Address, slotId uint, qopts ...pg.QOpt) (*Row, error)

	// Update inserts or updates the row identified by (Address, SlotId) pair.
	// When updating, the new row must have version greater than the existing,
	// otherwise ErrVersionTooLow is returned.
	// UpdatedAt field value is ignored.
	Update(row *Row, qopts ...pg.QOpt) error

	// DeleteExpired deletes any entries having Expiration < now().
	DeleteExpired(qopts ...pg.QOpt) error

	// GetSnapshot selects all rows ordered by UpdatedAt.
	// Min/Max addresses are specifying the scope.
	// To get a full snapshot, use MinAddress and MaxAddress variables.
	GetSnapshot(minAddress, maxAddress *big.Int, qopts ...pg.QOpt) ([]*Row, error)
}

func (r Row) Clone() *Row {
	clone := Row{
		Address:    r.Address,
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
