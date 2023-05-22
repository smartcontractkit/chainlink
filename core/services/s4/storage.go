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

// Record represents a user record persisted by S4
type Record struct {
	// Arbitrary user data
	Payload []byte
	// Version attribute assigned by user (unix timestamp is recommended)
	Version uint64
	// Expiration timestamp assigned by user (unix time in milliseconds)
	Expiration int64
}

// Metadata is the internal S4 data associated with a Record
type Metadata struct {
	Confirmed bool
	Signature []byte
}

//go:generate mockery --quiet --name Storage --output ./mocks/ --case=underscore

// Storage represents S4 storage access interface.
// All functions are thread-safe.
type Storage interface {
	// Constraints returns a copy of Constraints struct specified during service creation.
	// The implementation is thread-safe.
	Constraints() Constraints

	// Get returns a copy of record (with metadata) associated with the specified address and slotId.
	// The returned Record & Metadata are always a copy.
	Get(ctx context.Context, address common.Address, slotId uint) (*Record, *Metadata, error)

	// Put creates (or updates) a record identified by the specified address and slotId.
	// For signature calculation see envelope.go
	Put(ctx context.Context, address common.Address, slotId uint, record *Record, signature []byte) error
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

func (s *storage) Get(ctx context.Context, address common.Address, slotId uint) (*Record, *Metadata, error) {
	if slotId >= s.contraints.MaxSlotsPerUser {
		return nil, nil, ErrSlotIdTooBig
	}

	row, err := s.orm.Get(address, slotId, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, nil, err
	}

	if row.Expiration <= time.Now().UnixMilli() {
		return nil, nil, ErrNotFound
	}

	record := &Record{
		Payload:    make([]byte, len(row.Payload)),
		Version:    row.Version,
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

func (s *storage) Put(ctx context.Context, address common.Address, slotId uint, record *Record, signature []byte) error {
	if slotId >= s.contraints.MaxSlotsPerUser {
		return ErrSlotIdTooBig
	}
	if len(record.Payload) > int(s.contraints.MaxPayloadSizeBytes) {
		return ErrPayloadTooBig
	}
	if time.Now().UnixMilli() > record.Expiration {
		return ErrPastExpiration
	}

	envelope := NewEnvelopeFromRecord(address, slotId, record)
	signer, err := envelope.GetSignerAddress(signature)
	if err != nil || signer != address {
		return ErrWrongSignature
	}

	row, err := s.orm.Get(address, slotId, pg.WithParentCtx(ctx))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if row != nil && record.Version <= row.Version {
		return ErrVersionTooLow
	}

	row = &Row{
		Payload:    make([]byte, len(record.Payload)),
		Version:    record.Version,
		Expiration: record.Expiration,
		Confirmed:  false,
		Signature:  make([]byte, len(signature)),
	}
	copy(row.Payload, record.Payload)
	copy(row.Signature, signature)

	return s.orm.Upsert(address, slotId, row, pg.WithParentCtx(ctx))
}
