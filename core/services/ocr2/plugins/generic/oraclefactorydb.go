package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type oracleFactoryDb struct {
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
	_ ocrtypes.Database = &oracleFactoryDb{}
)

// NewDB returns a new DB scoped to this instanceID
func OracleFactoryDB(specID int32, lggr logger.Logger) *oracleFactoryDb {
	return &oracleFactoryDb{
		specID:               specID,
		lggr:                 logger.Sugared(lggr.Named("OracleFactoryMemoryDb")),
		states:               make(map[ocrtypes.ConfigDigest]*ocrtypes.PersistentState),
		pendingTransmissions: make(map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission),
		protocolStates:       make(map[ocrtypes.ConfigDigest]map[string][]byte),
	}
}

func (ofdb *oracleFactoryDb) ReadState(ctx context.Context, cd ocrtypes.ConfigDigest) (ps *ocrtypes.PersistentState, err error) {
	ps, ok := ofdb.states[cd]
	if !ok {
		return nil, fmt.Errorf("state not found for standard capabilities spec ID %d, config digest %s", ofdb.specID, cd)
	}

	return ps, nil
}

func (ofdb *oracleFactoryDb) WriteState(ctx context.Context, cd ocrtypes.ConfigDigest, state ocrtypes.PersistentState) error {
	ofdb.states[cd] = &state
	return nil
}

func (ofdb *oracleFactoryDb) ReadConfig(ctx context.Context) (c *ocrtypes.ContractConfig, err error) {
	if ofdb.config == nil {
		// Returning nil, nil because this is a cache miss
		return nil, nil
	}
	return ofdb.config, nil
}

func (ofdb *oracleFactoryDb) WriteConfig(ctx context.Context, c ocrtypes.ContractConfig) error {
	ofdb.config = &c

	cBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("MemoryDB: WriteConfig failed to marshal config: %w", err)
	}

	ofdb.lggr.Debugw("MemoryDB: WriteConfig", "ocrtypes.ContractConfig", string(cBytes))

	return nil
}

func (ofdb *oracleFactoryDb) StorePendingTransmission(ctx context.Context, t ocrtypes.ReportTimestamp, tx ocrtypes.PendingTransmission) error {
	ofdb.pendingTransmissions[t] = tx
	return nil
}

func (ofdb *oracleFactoryDb) PendingTransmissionsWithConfigDigest(ctx context.Context, cd ocrtypes.ConfigDigest) (map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission, error) {
	m := make(map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission)
	for k, v := range ofdb.pendingTransmissions {
		if k.ConfigDigest == cd {
			m[k] = v
		}
	}

	return m, nil
}

func (ofdb *oracleFactoryDb) DeletePendingTransmission(ctx context.Context, t ocrtypes.ReportTimestamp) error {
	delete(ofdb.pendingTransmissions, t)
	return nil
}

func (ofdb *oracleFactoryDb) DeletePendingTransmissionsOlderThan(ctx context.Context, t time.Time) error {
	for k, v := range ofdb.pendingTransmissions {
		if v.Time.Before(t) {
			delete(ofdb.pendingTransmissions, k)
		}
	}

	return nil
}

func (ofdb *oracleFactoryDb) ReadProtocolState(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	key string,
) ([]byte, error) {
	value, ok := ofdb.protocolStates[configDigest][key]
	if !ok {
		// Previously implementation returned nil if the state is not found
		return nil, nil
	}
	return value, nil
}

func (ofdb *oracleFactoryDb) WriteProtocolState(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	key string,
	value []byte,
) error {
	if value == nil {
		delete(ofdb.protocolStates[configDigest], key)
	} else {
		if ofdb.protocolStates[configDigest] == nil {
			ofdb.protocolStates[configDigest] = make(map[string][]byte)
		}
		ofdb.protocolStates[configDigest][key] = value
	}
	return nil
}
