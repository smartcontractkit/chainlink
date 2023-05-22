package s4

import (
	"context"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	"github.com/ethereum/go-ethereum/common"
)

// Constraints specifies the global storage constraints.
type Constraints struct {
	MaxPayloadSizeBytes uint
	MaxSlotsPerUser     uint
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

//go:generate mockery --quiet --name Storage --output ./mocks/ --case=underscore

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
}

type storage struct {
	lggr       logger.Logger
	contraints Constraints
	orm        ORM
}

var _ Storage = (*storage)(nil)

func NewStorage(lggr logger.Logger, contraints Constraints, orm ORM) Storage {
	return &storage{
		lggr:       lggr.Named("s4_storage"),
		contraints: contraints,
		orm:        orm,
	}
}

func (s *storage) Constraints() Constraints {
	return s.contraints
}

func (s *storage) Get(ctx context.Context, key *Key) (*Record, *Metadata, error) {
	if key.SlotId >= s.contraints.MaxSlotsPerUser {
		return nil, nil, ErrSlotIdTooBig
	}

	row, err := s.orm.Get(key.Address, key.SlotId, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, nil, err
	}

	if row.Expiration <= time.Now().UnixMilli() {
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

func (s *storage) Put(ctx context.Context, key *Key, record *Record, signature []byte) error {
	if key.SlotId >= s.contraints.MaxSlotsPerUser {
		return ErrSlotIdTooBig
	}
	if len(record.Payload) > int(s.contraints.MaxPayloadSizeBytes) {
		return ErrPayloadTooBig
	}
	if time.Now().UnixMilli() > record.Expiration {
		return ErrPastExpiration
	}

	envelope := NewEnvelopeFromRecord(key, record)
	signer, err := envelope.GetSignerAddress(signature)
	if err != nil || signer != key.Address {
		return ErrWrongSignature
	}

	row, err := s.orm.Get(key.Address, key.SlotId, pg.WithParentCtx(ctx))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if row != nil && key.Version <= row.Version {
		return ErrVersionTooLow
	}

	row = &Row{
		Payload:    make([]byte, len(record.Payload)),
		Version:    key.Version,
		Expiration: record.Expiration,
		Confirmed:  false,
		Signature:  make([]byte, len(signature)),
	}
	copy(row.Payload, record.Payload)
	copy(row.Signature, signature)

	return s.orm.Upsert(key.Address, key.SlotId, row, pg.WithParentCtx(ctx))
}
