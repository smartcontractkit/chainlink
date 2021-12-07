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
	TakeAndHold() error
	Release()
}

// advisoryLock implements the Locker interface.
type advisoryLock struct {
	id            int64
	db            *sqlx.DB
	conn          *sqlx.Conn
	checkInterval time.Duration
	logger        logger.Logger

	chStop chan struct{}
	wg     sync.WaitGroup
}

// NewAdvisoryLock returns an advisoryLocker
func NewAdvisoryLock(db *sqlx.DB, id int64, lggr logger.Logger, checkInterval time.Duration) AdvisoryLock {
	return &advisoryLock{id, db, nil, checkInterval, lggr.Named("AdvisoryLock").With("advisoryLockID", id), make(chan struct{}), sync.WaitGroup{}}
}

// TakeAndHold will block and wait indefinitely until it can get its first lock
// NOT THREAD SAFE
func (l *advisoryLock) TakeAndHold() (err error) {
	l.logger.Debug("Taking initial advisory lock...")
	retryCount := 0

	ctxStop, cancel := utils.ContextFromChan(l.chStop)
	defer cancel()

	for {
		var gotLock bool
		var err error

		err = func() error {
			ctx, cancel := DefaultQueryCtxWithParent(ctxStop)
			defer cancel()
			if l.conn == nil {
				if err = l.checkoutConn(ctx); err != nil {
					return errors.Wrap(err, "advisory lock failed to checkout initial connection")
				}
			}
			gotLock, err = l.getLock(ctx)
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
		select {
		case <-l.chStop:
			err = errors.New("stopped")
			if l.conn != nil {
				err = multierr.Combine(err, l.conn.Close())
			}
			return err
		case <-time.After(utils.WithJitter(l.checkInterval)):
		}
	}

	l.logger.Debug("Got advisory lock")
	l.wg.Add(1)
	go l.loop()
	return nil
}

// advisory lock must use one specific connection
func (l *advisoryLock) checkoutConn(ctx context.Context) (err error) {
	l.conn, err = l.db.Connx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed checking out connection from pool")
	}
	return nil
}

// getLock obtains exclusive session level advisory lock if available.
// It will either obtain the lock and return true, or return false if the lock cannot be acquired immediately.
func (l *advisoryLock) getLock(ctx context.Context) (locked bool, err error) {
	l.logger.Trace("Taking advisory lock")
	sqlQuery := "SELECT pg_try_advisory_lock($1)"
	err = l.conn.QueryRowContext(ctx, sqlQuery, l.id).Scan(&locked)
	return locked, err
}

func (l *advisoryLock) logRetry(count int) {
	if count%1000 == 0 || (count < 1000 && count&(count-1) == 0) {
		l.logger.Infow("Another application holds the advisory lock, waiting...", "tryCount", count)
	}
}

const checkAdvisoryLockStmt = `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE locktype = 'advisory' AND pid = pg_backend_pid() AND (classid::bigint << 32) | objid::bigint = $1)`

func (l *advisoryLock) loop() {
	defer l.wg.Done()

	ticker := time.NewTicker(utils.WithJitter(l.checkInterval))
	defer ticker.Stop()

	for {
		select {
		case <-l.chStop:
			ctx, cancel := DefaultQueryCtx()
			err := multierr.Combine(
				utils.JustError(l.conn.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, l.id)),
				l.conn.Close(),
			)
			cancel()
			if err != nil {
				l.logger.Warnw("Error trying to unlock advisory lock on shutdown", "err", err)
			}
			return
		case <-ticker.C:
			var gotLock bool

			ctx, cancel := utils.ContextFromChan(l.chStop)
			ctx, cancel2 := context.WithTimeout(ctx, l.checkInterval)
			l.logger.Trace("Checking advisory lock")
			err := l.conn.QueryRowContext(ctx, checkAdvisoryLockStmt, l.id).Scan(&gotLock)
			if errors.Is(err, sql.ErrConnDone) {
				// conn was closed and advisory lock lost, try to check out a new connection and re-lock
				l.logger.Warnw("DB connection was unexpectedly closed; checking out a new one", "err", err)
				if err = l.checkoutConn(ctx); err != nil {
					l.logger.Warnw("Error trying to checkout connection", "err", err)
				}
				gotLock, err = l.getLock(ctx)
			}
			cancel2()
			cancel()
			if err != nil {
				l.logger.Errorw("Error while checking advisory lock", "err", err)
			} else if !gotLock {
				l.logger.Fatal("Another node has taken the advisory lock, exiting")
			}
		}
	}
}

// Unlock releases the lock and DB connection.
// Should only be called once, and ends the lifecycle of the lock object
func (l *advisoryLock) Release() {
	close(l.chStop)
	l.wg.Wait()
}
