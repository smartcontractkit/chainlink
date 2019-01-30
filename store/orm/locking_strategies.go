package orm

import (
	"fmt"
	"time"
)

// NewLockingStrategy returns the appropriate dialect for the
// passed connection string.
func NewLockingStrategy(dialect DialectName, path string) (LockingStrategy, error) {
	switch dialect {
	case DialectPostgres:
		return UnlockedLockingStrategy{}, nil
	case DialectSqlite:
		return UnlockedLockingStrategy{}, nil
	}

	return nil, fmt.Errorf("unable to create locking strategy for dialect %s and path %s", dialect, path)
}

// LockingStrategy employs the locking and unlocking of an underlying
// resource, usually a file or db.
type LockingStrategy interface {
	Lock(timeout time.Duration) error
	Unlock()
}

// UnlockedLockingStrategy is a strategy that's always unlocked.
type UnlockedLockingStrategy struct{}

// Lock returns immediately and assumes is always unlocked.
func (s UnlockedLockingStrategy) Lock(timeout time.Duration) error { return nil }

// Unlock is a noop.
func (s UnlockedLockingStrategy) Unlock() {}
