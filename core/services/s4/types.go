package s4

import "errors"

type RecordState int

const (
	NewRecordState RecordState = iota
	ConfirmedRecordState
	ExpiredRecordState
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrRecordExpired     = errors.New("record expired")
	ErrWrongSignature    = errors.New("wrong signature")
	ErrTooBigSlotId      = errors.New("too big slot id")
	ErrTooBigPayload     = errors.New("too big payload")
	ErrOlderVersion      = errors.New("older version")
	ErrPastExpiration    = errors.New("past expiration")
	ErrServiceNotStarted = errors.New("service not started")
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
	State             RecordState
	HighestExpiration int64
	Signature         []byte
}

func (r Record) Clone() Record {
	clone := Record{
		Payload:    make([]byte, len(r.Payload)),
		Version:    r.Version,
		Expiration: r.Expiration,
	}
	copy(clone.Payload, r.Payload)
	return clone
}

func (m Metadata) Clone() Metadata {
	clone := Metadata{
		Signature:         make([]byte, len(m.Signature)),
		State:             m.State,
		HighestExpiration: m.HighestExpiration,
	}
	copy(clone.Signature, m.Signature)
	return clone
}
