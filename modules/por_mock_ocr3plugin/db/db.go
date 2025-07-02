package db

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3types.Database = &FakeDatabase{}

type FakeDatabase struct {
}

func NewFakeFatabase() *FakeDatabase {
	return &FakeDatabase{}
}

// In case the key is not found, nil should be returned.
func (db *FakeDatabase) ReadProtocolState(ctx context.Context, configDigest types.ConfigDigest, key string) ([]byte, error) {
	return nil, nil
}

// Writing with a nil value is the same as deleting.
func (db *FakeDatabase) WriteProtocolState(ctx context.Context, configDigest types.ConfigDigest, key string, value []byte) error {
	return nil
}

// ReadConfig returns the stored configuration (or nil if not set).
func (db *FakeDatabase) ReadConfig(ctx context.Context) (*types.ContractConfig, error) {
	return nil, nil
}

// WriteConfig stores the given configuration.
func (db *FakeDatabase) WriteConfig(ctx context.Context, config types.ContractConfig) error {
	return nil
}

// ReadState retrieves the persistent state for a given configuration digest.
// Returns nil if no state exists.
func (db *FakeDatabase) ReadState(ctx context.Context, configDigest types.ConfigDigest) (*types.PersistentState, error) {
	return nil, nil
}

// WriteState stores the persistent state for the given configuration digest.
func (db *FakeDatabase) WriteState(ctx context.Context, configDigest types.ConfigDigest, state types.PersistentState) error {
	return nil
}

// StorePendingTransmission stores a pending transmission associated with the
// ReportTimestamp’s configuration digest.
func (db *FakeDatabase) StorePendingTransmission(ctx context.Context, ts types.ReportTimestamp, pt types.PendingTransmission) error {
	return nil
}

// PendingTransmissionsWithConfigDigest returns all pending transmissions for a given configDigest.
func (db *FakeDatabase) PendingTransmissionsWithConfigDigest(ctx context.Context, configDigest types.ConfigDigest) (map[types.ReportTimestamp]types.PendingTransmission, error) {
	return nil, nil
}

// DeletePendingTransmission removes the pending transmission identified by the ReportTimestamp.
func (db *FakeDatabase) DeletePendingTransmission(ctx context.Context, ts types.ReportTimestamp) error {
	return nil
}

// DeletePendingTransmissionsOlderThan removes any pending transmissions whose
// associated transmission time is older than the specified time.
func (db *FakeDatabase) DeletePendingTransmissionsOlderThan(ctx context.Context, t time.Time) error {
	return nil
}
