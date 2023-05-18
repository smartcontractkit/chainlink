package s4

import (
	"context"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	"github.com/ethereum/go-ethereum/common"
)

type RecordState int

var (
	ErrRecordExpired  = errors.New("record expired")
	ErrWrongSignature = errors.New("wrong signature")
	ErrTooBigSlotId   = errors.New("too big slot id")
	ErrTooBigPayload  = errors.New("too big payload")
	ErrPastExpiration = errors.New("past expiration")
	ErrOlderVersion   = errors.New("older version")
)

// Constraints specifies the global storage constraints.
type Constraints struct {
	MaxPayloadSizeBytes int
	MaxSlotsPerUser     int
}

// Record represents a user record persisted by S4
type Record struct {
	// Arbitrary user data
	Payload []byte
	// Version attribute assigned by user
	Version int64
	// Expiration timestamp assigned by user (milliseconds)
	Expiration int64
}

// Metadata is the internal S4 data associated with a Record
type Metadata struct {
	Confirmed         bool
	HighestExpiration int64
	Signature         []byte
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
	Get(ctx context.Context, address common.Address, slotId int) (*Record, *Metadata, error)

	// Put creates (or updates) a record identified by the specified address and slotId.
	// For signature calculation see envelope.go
	Put(ctx context.Context, address common.Address, slotId int, record *Record, signature []byte) error
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

func (s *storage) Get(ctx context.Context, address common.Address, slotId int) (*Record, *Metadata, error) {
	if slotId >= s.contraints.MaxSlotsPerUser {
		return nil, nil, ErrTooBigSlotId
	}

	entry, err := s.orm.Get(address, slotId, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, nil, err
	}

	if entry.Expiration <= time.Now().UnixMilli() {
		return nil, nil, ErrRecordExpired
	}

	record := &Record{
		Payload:    make([]byte, len(entry.Payload)),
		Version:    entry.Version,
		Expiration: entry.Expiration,
	}
	copy(record.Payload, entry.Payload)

	metadata := &Metadata{
		Confirmed:         entry.Confirmed,
		HighestExpiration: entry.HighestExpiration,
		Signature:         make([]byte, len(entry.Signature)),
	}
	copy(metadata.Signature, entry.Signature)

	return record, metadata, nil
}

func (s *storage) Put(ctx context.Context, address common.Address, slotId int, record *Record, signature []byte) error {
	if slotId >= s.contraints.MaxSlotsPerUser {
		return ErrTooBigSlotId
	}
	if len(record.Payload) > s.contraints.MaxPayloadSizeBytes {
		return ErrTooBigPayload
	}
	if time.Now().UnixMilli() > record.Expiration {
		return ErrPastExpiration
	}

	envelope := NewEnvelopeFromRecord(address, slotId, record)
	signer, err := envelope.GetSignerAddress(signature)
	if err != nil || signer != address {
		return ErrWrongSignature
	}

	entry, err := s.orm.Get(address, slotId, pg.WithParentCtx(ctx))
	if err != nil && !errors.Is(err, ErrEntryNotFound) {
		return err
	}

	highestExpiration := record.Expiration
	if entry != nil {
		highestExpiration = entry.HighestExpiration
		if highestExpiration < record.Expiration {
			highestExpiration = record.Expiration
		}
		if record.Version <= entry.Version {
			return ErrOlderVersion
		}
	}

	entry = &Entry{
		Payload:           make([]byte, len(record.Payload)),
		Version:           record.Version,
		Expiration:        record.Expiration,
		Confirmed:         false,
		HighestExpiration: highestExpiration,
		Signature:         make([]byte, len(signature)),
	}
	copy(entry.Payload, record.Payload)
	copy(entry.Signature, signature)

	return s.orm.Upsert(address, slotId, entry, pg.WithParentCtx(ctx))
}
