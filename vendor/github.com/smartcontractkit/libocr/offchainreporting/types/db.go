package types

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/pkg/errors"
)

// Database persistently stores information on-disk.
// All its functions should be thread-safe.
type Database interface {
	ReadState(ctx context.Context, configDigest ConfigDigest) (*PersistentState, error)
	WriteState(ctx context.Context, configDigest ConfigDigest, state PersistentState) error

	ReadConfig(ctx context.Context) (*ContractConfig, error)
	WriteConfig(ctx context.Context, config ContractConfig) error

	StorePendingTransmission(context.Context, ReportTimestamp, PendingTransmission) error
	PendingTransmissionsWithConfigDigest(context.Context, ConfigDigest) (map[ReportTimestamp]PendingTransmission, error)
	DeletePendingTransmission(context.Context, ReportTimestamp) error
	DeletePendingTransmissionsOlderThan(context.Context, time.Time) error
}

type ReportTimestamp struct {
	ConfigDigest ConfigDigest
	Epoch        uint32
	Round        uint8
}

type PendingTransmission struct {
	Time             time.Time
	Median           Observation
	SerializedReport []byte
	Rs               [][32]byte
	Ss               [][32]byte
	Vs               [32]byte
}

type PersistentState struct {
	Epoch                uint32
	HighestSentEpoch     uint32
	HighestReceivedEpoch []uint32 // length: at most MaxOracles
}

func (ps PersistentState) Equal(ps2 PersistentState) bool {
	if ps.Epoch != ps2.Epoch {
		return false
	}
	if ps.HighestSentEpoch != ps2.HighestSentEpoch {
		return false
	}
	if len(ps.HighestReceivedEpoch) != len(ps2.HighestReceivedEpoch) {
		return false
	}
	for i := 0; i < len(ps.HighestReceivedEpoch); i++ {
		if ps.HighestReceivedEpoch[i] != ps2.HighestReceivedEpoch[i] {
			return false
		}
	}
	return true
}

//
// database/sql/driver interface functions for ConfigDigest
//

// Scan complies with sql Scanner interface
func (c *ConfigDigest) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("unable to convert %v of type %T to ConfigDigest", value, value)
	}
	if len(b) != 16 {
		return errors.Errorf("unable to convert blob 0x%x of length %v to ConfigDigest", b, len(b))
	}
	copy(c[:], b)
	return nil
}

// Value returns this instance serialized for database storage.
func (c ConfigDigest) Value() (driver.Value, error) {
	return c[:], nil
}
