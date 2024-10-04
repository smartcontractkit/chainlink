package generic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type memoryDb struct {
	// The ID is used for logging and error messages
	// A single standard capabilities spec can instantiate multiple oracles
	// TODO: NewOracle should take a unique identifier for the oracle
	specID               int32
	lggr                 logger.SugaredLogger
	config               *ocrtypes.ContractConfig
	states               map[ocrtypes.ConfigDigest]*ocrtypes.PersistentState
	pendingTransmissions map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission
	protocolStates       map[ocrtypes.ConfigDigest]map[string][]byte
}

var (
	_ ocrtypes.Database = &memoryDb{}
)

// NewDB returns a new DB scoped to this instanceID
func NewMemoryDB(specID int32, lggr logger.Logger) *memoryDb {
	return &memoryDb{
		specID:               specID,
		lggr:                 logger.Sugared(lggr.Named("OracleFactoryMemoryDb")),
		states:               make(map[ocrtypes.ConfigDigest]*ocrtypes.PersistentState),
		pendingTransmissions: make(map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission),
		protocolStates:       make(map[ocrtypes.ConfigDigest]map[string][]byte),
	}
}

func (md *memoryDb) ReadState(ctx context.Context, cd ocrtypes.ConfigDigest) (ps *ocrtypes.PersistentState, err error) {
	ps, ok := md.states[cd]
	if !ok {
		return nil, errors.Errorf("state not found for standard capabilities spec ID %d, config digest %s", md.specID, cd)
	}

	return ps, nil
}

func (md *memoryDb) WriteState(ctx context.Context, cd ocrtypes.ConfigDigest, state ocrtypes.PersistentState) error {
	md.states[cd] = &state
	return nil
}

func (md *memoryDb) ReadConfig(ctx context.Context) (c *ocrtypes.ContractConfig, err error) {
	if md.config == nil {
		// Returning nil, nil because this is a cache miss
		return nil, nil
	}
	return md.config, nil
}

func (md *memoryDb) WriteConfig(ctx context.Context, c ocrtypes.ContractConfig) error {
	md.config = &c

	cBytes, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "MemoryDB: WriteConfig failed to marshal config")
	}

	md.lggr.Debugw("MemoryDB: WriteConfig", "ocrtypes.ContractConfig", string(cBytes))

	return nil
}

func (md *memoryDb) StorePendingTransmission(ctx context.Context, t ocrtypes.ReportTimestamp, tx ocrtypes.PendingTransmission) error {
	md.pendingTransmissions[t] = tx
	return nil
}

func (md *memoryDb) PendingTransmissionsWithConfigDigest(ctx context.Context, cd ocrtypes.ConfigDigest) (map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission, error) {
	m := make(map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission)
	for k, v := range md.pendingTransmissions {
		if k.ConfigDigest == cd {
			m[k] = v
		}
	}

	return m, nil
}

func (md *memoryDb) DeletePendingTransmission(ctx context.Context, t ocrtypes.ReportTimestamp) error {
	delete(md.pendingTransmissions, t)
	return nil
}

func (md *memoryDb) DeletePendingTransmissionsOlderThan(ctx context.Context, t time.Time) error {
	for k, v := range md.pendingTransmissions {
		if v.Time.Before(t) {
			delete(md.pendingTransmissions, k)
		}
	}

	return nil
}

func (md *memoryDb) ReadProtocolState(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	key string,
) ([]byte, error) {
	value, ok := md.protocolStates[configDigest][key]
	if !ok {
		// Previously implementation returned nil if the state is not found
		// TODO: Should this return nil, nil?
		return nil, nil
	}
	return value, nil
}

func (md *memoryDb) WriteProtocolState(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	key string,
	value []byte,
) error {
	if value == nil {
		delete(md.protocolStates[configDigest], key)
	} else {
		if md.protocolStates[configDigest] == nil {
			md.protocolStates[configDigest] = make(map[string][]byte)
		}
		md.protocolStates[configDigest][key] = value
	}
	return nil
}
