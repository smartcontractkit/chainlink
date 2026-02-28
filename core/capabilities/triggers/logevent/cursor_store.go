package logevent

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// CursorRecord represents a persisted trigger cursor position.
type CursorRecord struct {
	TriggerID   string    `db:"trigger_id"`
	ChainID     string    `db:"chain_id"`
	Cursor      string    `db:"cursor_value"`
	BlockNumber uint64    `db:"block_number"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// CursorStore persists trigger cursor positions across restarts.
type CursorStore interface {
	// Get retrieves the last persisted cursor for a trigger.
	Get(ctx context.Context, triggerID string) (*CursorRecord, error)
	// Save persists the cursor position for a trigger.
	Save(ctx context.Context, triggerID string, chainID string, cursor string, blockNumber uint64) error
	// Delete removes the cursor for a trigger.
	Delete(ctx context.Context, triggerID string) error
}

// DBCursorStore is a Postgres-backed CursorStore.
type DBCursorStore struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

// NewDBCursorStore creates a new Postgres-backed cursor store.
func NewDBCursorStore(ds sqlutil.DataSource, lggr logger.Logger) *DBCursorStore {
	return &DBCursorStore{
		ds:   ds,
		lggr: logger.Named(lggr, "CursorStore"),
	}
}

func (s *DBCursorStore) Get(ctx context.Context, triggerID string) (*CursorRecord, error) {
	var record CursorRecord
	err := s.ds.GetContext(ctx, &record,
		`SELECT trigger_id, chain_id, cursor_value, block_number, updated_at
		 FROM evm_trigger_cursors WHERE trigger_id = $1`, triggerID)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *DBCursorStore) Save(ctx context.Context, triggerID string, chainID string, cursor string, blockNumber uint64) error {
	_, err := s.ds.ExecContext(ctx,
		`INSERT INTO evm_trigger_cursors (trigger_id, chain_id, cursor_value, block_number, updated_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (trigger_id) DO UPDATE SET
			cursor_value = EXCLUDED.cursor_value,
			block_number = EXCLUDED.block_number,
			updated_at = NOW()`,
		triggerID, chainID, cursor, blockNumber)
	return err
}

func (s *DBCursorStore) Delete(ctx context.Context, triggerID string) error {
	_, err := s.ds.ExecContext(ctx,
		`DELETE FROM evm_trigger_cursors WHERE trigger_id = $1`, triggerID)
	return err
}

// InMemoryCursorStore is a simple in-memory implementation for testing.
type InMemoryCursorStore struct {
	cursors map[string]*CursorRecord
}

func NewInMemoryCursorStore() *InMemoryCursorStore {
	return &InMemoryCursorStore{
		cursors: make(map[string]*CursorRecord),
	}
}

func (s *InMemoryCursorStore) Get(_ context.Context, triggerID string) (*CursorRecord, error) {
	record, ok := s.cursors[triggerID]
	if !ok {
		return nil, nil
	}
	return record, nil
}

func (s *InMemoryCursorStore) Save(_ context.Context, triggerID string, chainID string, cursor string, blockNumber uint64) error {
	s.cursors[triggerID] = &CursorRecord{
		TriggerID:   triggerID,
		ChainID:     chainID,
		Cursor:      cursor,
		BlockNumber: blockNumber,
		UpdatedAt:   time.Now(),
	}
	return nil
}

func (s *InMemoryCursorStore) Delete(_ context.Context, triggerID string) error {
	delete(s.cursors, triggerID)
	return nil
}
