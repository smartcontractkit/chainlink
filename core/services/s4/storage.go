package s4

import (
	"context"

	"github.com/jonboulle/clockwork"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Constraints specifies the global storage constraints.
type Constraints struct {
	MaxPayloadSizeBytes    uint   `json:"maxPayloadSizeBytes"`
	MaxSlotsPerUser        uint   `json:"maxSlotsPerUser"`
	MaxExpirationLengthSec uint64 `json:"maxExpirationLengthSec"`
}

// Key identifies a versioned user record.
type Key struct {
	// Address is a user address
	Address common.Address
	// SlotId is a slot number
	SlotId uint
	// Version is a data version
	Version uint64
}

// Record represents a user record persisted by S4.
type Record struct {
	// Arbitrary user data
	Payload []byte
	// Expiration timestamp assigned by user (unix time in milliseconds)
	Expiration int64
}

// Metadata is the internal S4 data associated with a Record
type Metadata struct {
	// Confirmed turns true once consensus is reached.
	Confirmed bool
	// Signature contains the original user signature.
	Signature []byte
}

// Storage represents S4 storage access interface.
// All functions are thread-safe.
type Storage interface {
	// Constraints returns a copy of Constraints struct specified during service creation.
	// The implementation is thread-safe.
	Constraints() Constraints

	// Get returns a copy of record (with metadata) associated with the specified key.
	// The returned Record & Metadata are always a copy.
	Get(ctx context.Context, key *Key) (*Record, *Metadata, error)

	// Put creates (or updates) a record identified by the specified key.
	// For signature calculation see envelope.go
	Put(ctx context.Context, key *Key, record *Record, signature []byte) error

	// List returns a snapshot for the specified address.
	// Slots having no data are not returned.
	List(ctx context.Context, address common.Address) ([]*SnapshotRow, error)
}

type storage struct {
	lggr       logger.Logger
	contraints Constraints
	orm        ORM
	clock      clockwork.Clock
}

var _ Storage = (*storage)(nil)

func NewStorage(lggr logger.Logger, contraints Constraints, orm ORM, clock clockwork.Clock) Storage {
	return &storage{
		lggr:       lggr.Named("S4Storage"),
		contraints: contraints,
		orm:        orm,
		clock:      clock,
	}
}

func (s *storage) Constraints() Constraints {
	return s.contraints
}

func (s *storage) Get(ctx context.Context, key *Key) (*Record, *Metadata, error) {
	if key.SlotId >= s.contraints.MaxSlotsPerUser {
		return nil, nil, ErrSlotIdTooBig
	}

	bigAddress := big.New(key.Address.Big())
	row, err := s.orm.Get(ctx, bigAddress, key.SlotId)
	if err != nil {
		return nil, nil, err
	}

	if row.Version != key.Version || row.Expiration <= s.clock.Now().UnixMilli() {
		return nil, nil, ErrNotFound
	}

	record := &Record{
		Payload:    make([]byte, len(row.Payload)),
		Expiration: row.Expiration,
	}
	copy(record.Payload, row.Payload)

	metadata := &Metadata{
		Confirmed: row.Confirmed,
		Signature: make([]byte, len(row.Signature)),
	}
	copy(metadata.Signature, row.Signature)

	return record, metadata, nil
}

func (s *storage) List(ctx context.Context, address common.Address) ([]*SnapshotRow, error) {
	bigAddress := big.New(address.Big())
	sar, err := NewSingleAddressRange(bigAddress)
	if err != nil {
		return nil, err
	}
	return s.orm.GetSnapshot(ctx, sar)
}

func (s *storage) Put(ctx context.Context, key *Key, record *Record, signature []byte) error {
	if key.SlotId >= s.contraints.MaxSlotsPerUser {
		return ErrSlotIdTooBig
	}
	if len(record.Payload) > int(s.contraints.MaxPayloadSizeBytes) {
		return ErrPayloadTooBig
	}
	now := s.clock.Now().UnixMilli()
	if now > record.Expiration {
		return ErrPastExpiration
	}
	if record.Expiration-now > int64(s.contraints.MaxExpirationLengthSec)*1000 {
		return ErrExpirationTooLong
	}

	envelope := NewEnvelopeFromRecord(key, record)
	signer, err := envelope.GetSignerAddress(signature)
	if err != nil || signer != key.Address {
		return ErrWrongSignature
	}

	row := &Row{
		Address:    big.New(key.Address.Big()),
		SlotId:     key.SlotId,
		Payload:    make([]byte, len(record.Payload)),
		Version:    key.Version,
		Expiration: record.Expiration,
		Confirmed:  false,
		Signature:  make([]byte, len(signature)),
	}
	copy(row.Payload, record.Payload)
	copy(row.Signature, signature)

	return s.orm.Update(ctx, row)
}
