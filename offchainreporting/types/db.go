package types

import (
	"database/sql/driver"
	"time"

	"github.com/pkg/errors"
)

type Database interface {
	ReadState(configDigest ConfigDigest) (*PersistentState, error)
	WriteState(configDigest ConfigDigest, state PersistentState) error

	ReadConfig() (*ContractConfig, error)
	WriteConfig(config ContractConfig) error

	StorePendingTransmission(PendingTransmissionKey, PendingTransmission) error
	PendingTransmissionsWithConfigDigest(ConfigDigest) (map[PendingTransmissionKey]PendingTransmission, error)
	DeletePendingTransmission(PendingTransmissionKey) error
	DeletePendingTransmissionsOlderThan(time.Time) error
}

type PendingTransmissionKey struct {
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
	HighestReceivedEpoch []uint32
}

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

func (c ConfigDigest) Value() (driver.Value, error) {
	return c[:], nil
}
