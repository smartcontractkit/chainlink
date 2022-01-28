package pg

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// AdvisoryLock is an interface for postgresql advisory locks.
type AdvisoryLock interface {
	TakeAndHold(ctx context.Context) error
	Release()
}

// advisoryLock implements the Locker interface.
type advisoryLock struct {
	id            int64
	db            *sqlx.DB
	conn          *sqlx.Conn
	checkInterval time.Duration
	logger        logger.Logger
	stop          func()
	wgReleased    sync.WaitGroup
}

// NewAdvisoryLock returns an advisoryLocker
func NewAdvisoryLock(db *sqlx.DB, id int64, lggr logger.Logger, checkInterval time.Duration) AdvisoryLock {
	return &advisoryLock{id, db, nil, checkInterval, lggr.Named("AdvisoryLock").With("advisoryLockID", id), func() {}, sync.WaitGroup{}}
}

// TakeAndHold will block and wait indefinitely until it can get its first lock or ctx is cancelled.
// Use Release() function to release the acquired lock.
// NOT THREAD SAFE
func (l *advisoryLock) TakeAndHold(ctx context.Context) (err error) {
	l.logger.Debug("Taking initial advisory lock...")
	retryCount := 0

	lctx, cancel := context.WithCancel(context.Background())
	l.stop = cancel

	for {
		var gotLock bool
		var err error

		err = func() error {
			qctx, cancel := DefaultQueryCtxWithParent(lctx)
			defer cancel()
			if l.conn == nil {
				if err = l.checkoutConn(qctx); err != nil {
					return errors.Wrap(err, "advisory lock failed to checkout initial connection")
				}
			}
			gotLock, err = l.getLock(qctx)
			if errors.Is(err, sql.ErrConnDone) {
				l.logger.Warnw("DB connection was unexpectedly closed; checking out a new one", "err", err)
				l.conn = nil
				return err
			}
			return nil
		}()

		if errors.Is(err, sql.ErrConnDone) {
			continue
		} else if err != nil {
			err = errors.Wrap(err, "failed to get advisory lock")
			if l.conn != nil {
				err = multierr.Combine(err, l.conn.Close())
			}
			return err
		}
		if gotLock {
			break
		}
		l.logRetry(retryCount)
		retryCount++
		err = func() error {
			select {
			case <-ctx.Done():
				return errors.New("stopped by parent context")
			case <-lctx.Done():
				return errors.New("stopped by Release()")
			case <-time.After(utils.WithJitter(l.checkInterval)):
				return nil
			}
		}()
		if err != nil {
			if l.conn != nil {
				err = multierr.Combine(err, l.conn.Close())
			}
			return err
		}
	}

	l.logger.Debug("Got advisory lock")
	l.wgReleased.Add(1)
	go l.loop(lctx)
	return nil
}

// Release requests the lock to release and blocks until it gets released.
// Calling Release for a released lock has no effect.
func (l *advisoryLock) Release() {
	l.stop()
	l.wgReleased.Wait()
}

// advisory lock must use one specific connection
func (l *advisoryLock) checkoutConn(ctx context.Context) (err error) {
	newConn, err := l.db.Connx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed checking out connection from pool")
	}
	l.conn = newConn
	return nil
}

// getLock obtains exclusive session level advisory lock if available.
// It will either obtain the lock and return true, or return false if the lock cannot be acquired immediately.
func (l *advisoryLock) getLock(ctx context.Context) (locked bool, err error) {
	l.logger.Trace("Taking advisory lock")
	sqlQuery := "SELECT pg_try_advisory_lock($1)"
	err = l.conn.QueryRowContext(ctx, sqlQuery, l.id).Scan(&locked)
	return locked, errors.WithStack(err)
}

func (l *advisoryLock) logRetry(count int) {
	if count%1000 == 0 || (count < 1000 && count&(count-1) == 0) {
		l.logger.Infow("Another application holds the advisory lock, waiting...", "tryCount", count)
	}
}

const checkAdvisoryLockStmt = `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE locktype = 'advisory' AND pid = pg_backend_pid() AND (classid::bigint << 32) | objid::bigint = $1)`

func (l *advisoryLock) loop(ctx context.Context) {
	defer l.wgReleased.Done()

	ticker := time.NewTicker(utils.WithJitter(l.checkInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			qctx, cancel := DefaultQueryCtx()
			err := multierr.Combine(
				utils.JustError(l.conn.ExecContext(qctx, `SELECT pg_advisory_unlock($1)`, l.id)),
				l.conn.Close(),
			)
			cancel()
			if err != nil {
				l.logger.Warnw("Error trying to unlock advisory lock on shutdown", "err", err)
			}
			return
		case <-ticker.C:
			var gotLock bool

			qctx, cancel := DefaultQueryCtxWithParent(ctx)
			l.logger.Trace("Checking advisory lock")
			err := l.conn.QueryRowContext(qctx, checkAdvisoryLockStmt, l.id).Scan(&gotLock)
			if errors.Is(err, sql.ErrConnDone) {
				// conn was closed and advisory lock lost, try to check out a new connection and re-lock
				l.logger.Warnw("DB connection was unexpectedly closed; checking out a new one", "err", err)
				if err = l.checkoutConn(qctx); err != nil {
					l.logger.Warnw("Error trying to checkout connection", "err", err)
				}
				gotLock, err = l.getLock(qctx)
			}
			cancel()
			if err != nil {
				l.logger.Errorw("Error while checking advisory lock", "err", err)
			} else if !gotLock {
				l.logger.Fatal("Another node has taken the advisory lock, exiting")
			}
			ticker.Reset(utils.WithJitter(l.checkInterval))
		}
	}
}
