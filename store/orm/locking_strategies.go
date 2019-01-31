package orm

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/jinzhu/gorm"
	"go.uber.org/multierr"
)

// NewLockingStrategy returns the locking strategy for a particular dialect
// to ensure exlusive access to the orm.
func NewLockingStrategy(dialect DialectName, path string) (LockingStrategy, error) {
	switch dialect {
	case DialectPostgres:
		return NewPostgresLockingStrategy(dialect, path)
	case DialectSqlite:
		return NewFileLockingStrategy(dialect, path)
	}

	return nil, fmt.Errorf("unable to create locking strategy for dialect %s and path %s", dialect, path)
}

// LockingStrategy employs the locking and unlocking of an underlying
// resource for exclusive access, usually a file or database.
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

// FileLockingStrategy uses a file lock on disk to ensure exclusive access.
type FileLockingStrategy struct {
	path     string
	fileLock *flock.Flock
	m        *sync.Mutex
}

// NewFileLockingStrategy creates a new instance of FileLockingStrategy
// at the passed path.
func NewFileLockingStrategy(_ DialectName, dbpath string) (LockingStrategy, error) {
	uri, err := url.Parse(dbpath)
	if err != nil {
		return nil, multierr.Append(errors.New("unable to create file locking strategy"), err)
	}
	dbpath = uri.Path
	directory := filepath.Dir(dbpath)
	lockPath := filepath.Join(directory, "chainlink.lock")
	return &FileLockingStrategy{
		path:     lockPath,
		fileLock: flock.New(lockPath),
		m:        &sync.Mutex{},
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
	db          *gorm.DB
	dialectName DialectName
	path        string
	m           *sync.Mutex
}

// NewPostgresLockingStrategy returns a new instance of the PostgresLockingStrategy.
func NewPostgresLockingStrategy(dialectName DialectName, path string) (LockingStrategy, error) {
	return &PostgresLockingStrategy{
		m:           &sync.Mutex{},
		dialectName: dialectName,
		path:        path,
	}, nil
}

const postgresAdvisoryLockID int64 = 1027321974924625846

// Lock uses a blocking postgres advisory lock that times out at the passed
// timeout.
func (s *PostgresLockingStrategy) Lock(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	db, err := gorm.Open(string(s.dialectName), s.path)
	if err != nil {
		return err
	}
	s.db = db

	locked := make(chan struct{})
	go func() {
		err = s.db.Exec("SELECT pg_advisory_lock(?);", postgresAdvisoryLockID).Error
		close(locked)
	}()
	select {
	case <-locked:
	case <-normalizedTimeout(timeout):
		return fmt.Errorf("postgres advisory locking strategy timed out for advisory lock ID %v", postgresAdvisoryLockID)
	}
	return err
}

// Unlock unlocks the locked postgres advisory lock.
func (s *PostgresLockingStrategy) Unlock() error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
