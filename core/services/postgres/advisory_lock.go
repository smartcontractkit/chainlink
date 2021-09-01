package postgres

import (
	"context"
	"database/sql"
	"net/url"
	"sync"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// NOTE: All advisory lock class IDs used by the Chainlink application MUST be
// kept here to avoid accidental re-use
const (
	AdvisoryLockClassID_EthBroadcaster int32 = 0
	AdvisoryLockClassID_JobSpawner     int32 = 1
	AdvisoryLockClassID_EthConfirmer   int32 = 2

	// ORM takes lock on 1027321974924625846 which splits into ClassID 239192036, ObjID 2840971190
	AdvisoryLockClassID_ORM int32 = 239192036

	AdvisoryLockObjectID_EthConfirmer int32 = 0
)

//go:generate mockery --name AdvisoryLocker --output ../../internal/mocks/ --case=underscore
type (
	postgresAdvisoryLock struct {
		URI  string
		conn *sql.Conn
		db   *sql.DB
		mu   *sync.Mutex
	}

	AdvisoryLocker interface {
		Unlock(ctx context.Context, classID int32, objectID int32) error
		WithAdvisoryLock(ctx context.Context, classID int32, objectID int32, f func() error) error
		Close() error
	}
)

func NewAdvisoryLock(uri url.URL) AdvisoryLocker {
	static.SetConsumerName(&uri, "AdvisoryLocker")
	return &postgresAdvisoryLock{
		URI: uri.String(),
		mu:  &sync.Mutex{},
	}
}

func (lock *postgresAdvisoryLock) Close() error {
	lock.mu.Lock()
	defer lock.mu.Unlock()

	var connErr, dbErr error

	if lock.conn != nil {
		connErr = lock.conn.Close()
		if connErr == sql.ErrConnDone {
			connErr = nil
		}
	}
	if lock.db != nil {
		dbErr = lock.db.Close()
		if dbErr == sql.ErrConnDone {
			dbErr = nil
		}
	}

	lock.db = nil
	lock.conn = nil

	return multierr.Combine(connErr, dbErr)
}

func (lock *postgresAdvisoryLock) tryLock(ctx context.Context, classID int32, objectID int32) (err error) {
	lock.mu.Lock()
	defer lock.mu.Unlock()
	defer utils.WrapIfError(&err, "TryAdvisoryLock failed")

	if lock.conn == nil {
		db, err2 := sql.Open(string(dialects.Postgres), lock.URI)
		if err2 != nil {
			return err2
		}
		lock.db = db

		// `database/sql`.DB does opaque connection pooling, but PG advisory locks are per-connection
		conn, err2 := db.Conn(ctx)
		if err2 != nil {
			logger.ErrorIfCalling(lock.db.Close)
			lock.db = nil
			return err2
		}
		lock.conn = conn
	}

	gotLock := false
	if err = lock.conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1, $2)", classID, objectID).Scan(&gotLock); err != nil {
		return err
	}
	if gotLock {
		return nil
	}
	return errors.Errorf("could not get advisory lock for classID, objectID %v, %v", classID, objectID)
}

func (lock *postgresAdvisoryLock) Unlock(ctx context.Context, classID int32, objectID int32) error {
	lock.mu.Lock()
	defer lock.mu.Unlock()

	if lock.conn == nil {
		return nil
	}
	_, err := lock.conn.ExecContext(ctx, "SELECT pg_advisory_unlock($1, $2)", classID, objectID)
	return errors.Wrap(err, "AdvisoryUnlock failed")
}

func (lock *postgresAdvisoryLock) WithAdvisoryLock(ctx context.Context, classID int32, objectID int32, f func() error) error {
	err := lock.tryLock(ctx, classID, objectID)
	if err != nil {
		return errors.Wrapf(err, "could not get advisory lock for classID, objectID %v, %v", classID, objectID)
	}
	defer logger.ErrorIfCalling(func() error { return lock.Unlock(ctx, classID, objectID) })
	return f()
}

var _ AdvisoryLocker = &NullAdvisoryLocker{}

func NewNullAdvisoryLocker() *NullAdvisoryLocker {
	return &NullAdvisoryLocker{}
}

type NullAdvisoryLocker struct {
	mu     sync.Mutex
	closed bool
}

func (n *NullAdvisoryLocker) Close() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.closed {
		panic("already closed")
	}
	n.closed = true
	return nil
}

func (*NullAdvisoryLocker) Unlock(ctx context.Context, classID int32, objectID int32) error {
	return nil
}

func (*NullAdvisoryLocker) WithAdvisoryLock(ctx context.Context, classID int32, objectID int32, f func() error) error {
	return f()
}
