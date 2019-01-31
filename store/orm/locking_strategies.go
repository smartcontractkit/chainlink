package orm

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

// NewLockingStrategy returns the appropriate dialect for the
// passed connection string.
func NewLockingStrategy(dialect DialectName, path string) (LockingStrategy, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	path = uri.Path
	switch dialect {
	case DialectPostgres:
		return UnlockedLockingStrategy{}, nil
	case DialectSqlite:
		directory := filepath.Dir(path)
		return NewFileLockingStrategy(filepath.Join(directory, "chainlink.lock")), nil
	}

	return nil, fmt.Errorf("unable to create locking strategy for dialect %s and path %s", dialect, path)
}

// LockingStrategy employs the locking and unlocking of an underlying
// resource, usually a file or db.
type LockingStrategy interface {
	Lock(timeout time.Duration) error
	Unlock() error
}

// UnlockedLockingStrategy is a strategy that's always unlocked.
type UnlockedLockingStrategy struct{}

// Lock returns immediately and assumes is always unlocked.
func (s UnlockedLockingStrategy) Lock(timeout time.Duration) error { return nil }

// Unlock is a noop.
func (s UnlockedLockingStrategy) Unlock() error { return nil }

// FileLockingStrategy uses the write lock for a file on disk to ensure
// exclusive operations.
type FileLockingStrategy struct {
	path     string
	fileLock *flock.Flock
	m        *sync.Mutex
}

// NewFileLockingStrategy creates a new instance of FileLockingStrategy
// at the passed path.
func NewFileLockingStrategy(path string) LockingStrategy {
	return FileLockingStrategy{
		path:     path,
		fileLock: flock.New(path),
		m:        &sync.Mutex{},
	}
}

// Lock returns immediately and assumes is always unlocked.
func (s FileLockingStrategy) Lock(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	var err error
	locked := make(chan struct{})
	go func() {
		err = s.fileLock.Lock()
		close(locked)
	}()
	select {
	case <-locked:
	case <-time.After(timeout):
		return fmt.Errorf("file locking strategy timed out for %s", s.path)
	}
	return err
}

// Unlock is a noop.
func (s FileLockingStrategy) Unlock() error {
	s.m.Lock()
	defer s.m.Unlock()
	return s.fileLock.Unlock()
}
