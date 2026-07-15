package v2

import (
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
)

// MaxModulePanics is how many times a workflow module may panic the host before
// it is quarantined and no longer executed. Each panic rethrows and restarts the
// process, so this bounds a deterministic-panic reboot loop to at most
// MaxModulePanics restarts before the workflow is disabled and the node stays up.
const MaxModulePanics = 3

type ModulePanicStore interface {
	Store(ctx context.Context, key string, val []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

func modulePanicKey(workflowID string) string {
	return "module_panic_count:" + workflowID
}

func recordModulePanic(ctx context.Context, s ModulePanicStore, workflowID string) (uint64, error) {
	count, err := modulePanicCount(ctx, s, workflowID)
	if err != nil {
		return 0, err
	}
	count++

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], count)
	if err := s.Store(ctx, modulePanicKey(workflowID), buf[:]); err != nil {
		return 0, fmt.Errorf("store module panic count: %w", err)
	}
	return count, nil
}

func modulePanicCount(ctx context.Context, s ModulePanicStore, workflowID string) (uint64, error) {
	b, err := s.Get(ctx, modulePanicKey(workflowID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("get module panic count: %w", err)
	}
	if len(b) < 8 {
		return 0, nil
	}
	return binary.BigEndian.Uint64(b), nil
}

// ClearModulePanic removes a workflow's panic counter so it is no longer
// quarantined. Call it when the workflow's artifacts change (redeploy) or the
// workflow is deleted. A missing entry is not an error.
func ClearModulePanic(ctx context.Context, s ModulePanicStore, workflowID string) error {
	if err := s.Delete(ctx, modulePanicKey(workflowID)); err != nil {
		return fmt.Errorf("clear module panic count: %w", err)
	}
	return nil
}
