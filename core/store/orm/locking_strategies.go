package orm

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"go.uber.org/multierr"
)

// NewLockingStrategy returns the locking strategy for a particular dialect
// to ensure exlusive access to the orm.
func NewLockingStrategy(ct Connection) (LockingStrategy, error) {
	switch ct.name {
	case dialects.Postgres, dialects.PostgresWithoutLock, dialects.TransactionWrappedPostgres:
		return NewPostgresLockingStrategy(ct)
	}

	return nil, fmt.Errorf("unable to create locking strategy for dialect %s and path %s", ct.dialect, ct.uri)
}

// LockingStrategy employs the locking and unlocking of an underlying
// resource for exclusive access, usually a file or database.
type LockingStrategy interface {
	Lock(timeout models.Duration) error
	Unlock(timeout models.Duration) error
}

// PostgresLockingStrategy uses a postgres advisory lock to ensure exclusive
// access.
type PostgresLockingStrategy struct {
	db     *sql.DB
	conn   *sql.Conn
	m      *sync.Mutex
	config Connection
}

// NewPostgresLockingStrategy returns a new instance of the PostgresLockingStrategy.
func NewPostgresLockingStrategy(ct Connection) (LockingStrategy, error) {
	return &PostgresLockingStrategy{
		config: ct,
		m:      &sync.Mutex{},
	}, nil
}

// Lock uses a blocking postgres advisory lock that times out at the passed
// timeout.
func (s *PostgresLockingStrategy) Lock(timeout models.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	ctx := context.Background()
	if !timeout.IsInstant() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout.Duration())
		defer cancel()
	}

	if s.conn == nil {
		uri, err := url.Parse(s.config.uri)
		if err != nil {
			return err
		}
		static.SetConsumerName(uri, "PostgresLockingStrategy")
		db, err := sql.Open(string(dialects.Postgres), uri.String())
		if err != nil {
			return err
		}
		s.db = db

		// `database/sql`.DB does opaque connection pooling, but PG advisory locks are per-connection
		conn, err := db.Conn(ctx)
		if err != nil {
			return err
		}

		s.conn = conn
	}

	if s.config.locking {
		err := s.waitForLock(ctx)
		if err != nil {
			return errors.Wrapf(ErrNoAdvisoryLock, "postgres advisory locking strategy failed on .Lock, timeout set to %v: %v, lock ID: %v", displayTimeout(timeout), err, s.config.advisoryLockID)
		}
	}
	return nil
}

func (s *PostgresLockingStrategy) waitForLock(ctx context.Context) error {
	ticker := time.NewTicker(s.config.lockRetryInterval)
	defer ticker.Stop()
	retryCount := 0
	for {
		rows, err := s.conn.QueryContext(ctx, "SELECT pg_try_advisory_lock($1)", s.config.advisoryLockID)
		if err != nil {
			return err
		}
		var gotLock bool
		for rows.Next() {
			err := rows.Scan(&gotLock)
			if err != nil {
				return multierr.Combine(err, rows.Close())
			}
		}
		if err := rows.Close(); err != nil {
			return err
		}
		if gotLock {
			return nil
		}

		select {
		case <-ticker.C:
			retryCount++
			logRetry(retryCount)
			continue
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "timeout expired while waiting for lock")
		}
	}
}

// logRetry logs messages at
// 1
// 2
// 4
// 8
// 16
// 32
/// ... etc, then every 1000
func logRetry(count int) {
	if count == 1 {
		logger.Infow("Could not get lock, retrying...", "failCount", count)
	} else if count%1000 == 0 || count&(count-1) == 0 {
		logger.Infow("Still waiting for lock...", "failCount", count)
	}
}

// Unlock unlocks the locked postgres advisory lock.
func (s *PostgresLockingStrategy) Unlock(timeout models.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.conn == nil {
		return nil
	}

	connErr := s.conn.Close()
	if connErr == sql.ErrConnDone {
		connErr = nil
	}
	dbErr := s.db.Close()
	if dbErr == sql.ErrConnDone {
		dbErr = nil
	}

	s.db = nil
	s.conn = nil

	return multierr.Combine(
		connErr,
		dbErr,
	)
}
