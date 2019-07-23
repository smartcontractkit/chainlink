package orm

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/multierr"
)

// NewLockingStrategy returns the locking strategy for a particular dialect
// to ensure exlusive access to the orm.
func NewLockingStrategy(dialect DialectName, dbpath string) (LockingStrategy, error) {
	switch dialect {
	case DialectPostgres:
		return NewPostgresLockingStrategy(dbpath)
	case DialectSqlite:
		return NewFileLockingStrategy(dbpath)
	}

	return nil, fmt.Errorf("unable to create locking strategy for dialect %s and path %s", dialect, dbpath)
}

// LockingStrategy employs the locking and unlocking of an underlying
// resource for exclusive access, usually a file or database.
type LockingStrategy interface {
	Lock(timeout time.Duration) error
	Unlock() error
}

// FileLockingStrategy uses a file lock on disk to ensure exclusive access.
type FileLockingStrategy struct {
	path     string
	fileLock *flock.Flock
	m        *deadlock.Mutex
}

// NewFileLockingStrategy creates a new instance of FileLockingStrategy
// at the passed path.
func NewFileLockingStrategy(dbpath string) (LockingStrategy, error) {
	directory := filepath.Dir(dbpath)
	lockPath := filepath.Join(directory, "chainlink.lock")
	return &FileLockingStrategy{
		path:     lockPath,
		fileLock: flock.New(lockPath),
		m:        &deadlock.Mutex{},
	}, nil
}

// Lock returns immediately and assumes is always unlocked.
func (s *FileLockingStrategy) Lock(timeout time.Duration) error {
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
	case <-normalizedTimeout(timeout):
		return fmt.Errorf("file locking strategy timed out for %s", s.path)
	}
	return err
}

func normalizedTimeout(timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return make(chan time.Time) // never time out
	}
	return time.After(timeout)
}

// Unlock is a noop.
func (s *FileLockingStrategy) Unlock() error {
	s.m.Lock()
	defer s.m.Unlock()
	return s.fileLock.Unlock()
}

// PostgresLockingStrategy uses a postgres advisory lock to ensure exclusive
// access.
type PostgresLockingStrategy struct {
	db     *sql.DB
	path   string
	m      *deadlock.Mutex
	locked bool
}

// NewPostgresLockingStrategy returns a new instance of the PostgresLockingStrategy.
func NewPostgresLockingStrategy(path string) (LockingStrategy, error) {
	return &PostgresLockingStrategy{
		m:    &deadlock.Mutex{},
		path: path,
	}, nil
}

const postgresAdvisoryLockID int64 = 1027321974924625846

// Lock uses a blocking postgres advisory lock that times out at the passed
// timeout.
func (s *PostgresLockingStrategy) Lock(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.locked {
		return nil
	}

	db, err := sql.Open(string(DialectPostgres), s.path)
	if err != nil {
		return err
	}
	s.db = db

	ctx := context.Background()
	if timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	_, err = s.db.ExecContext(ctx, "SELECT pg_advisory_lock($1)", postgresAdvisoryLockID)
	if err != nil {
		return multierr.Append(
			fmt.Errorf("postgres advisory locking strategy failed, timeout set to %v", displayTimeout(timeout)),
			err)
	}
	s.locked = true
	return nil
}

// Unlock unlocks the locked postgres advisory lock.
func (s *PostgresLockingStrategy) Unlock() error {
	s.m.Lock()
	defer s.m.Unlock()

	s.locked = false
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
