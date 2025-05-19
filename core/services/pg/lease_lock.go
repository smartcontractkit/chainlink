package pg

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// LeaseLock handles taking an exclusive lease on database access. This is not
// enforced by any database primitives, but rather voluntarily respected by
// other instances of the Chainlink application.
//
// Chainlink is designed to run as a single instance. Running multiple
// instances of Chainlink on a single database at the same time is not
// supported and likely to lead to strange errors and possibly even data
// integrity failures.
//
// With that being said, a common use case is to run multiple Chainlink
// instances in failover mode. The first instance will take some kind of lock
// on the database and subsequent instances will wait trying to take this lock
// in case the first instance disappears or dies.
//
// Traditionally Chainlink has used an advisory lock to manage this. However,
// advisory locks come with several problems, notably: - Postgres does not
// really like it when you hold locks open for a very long time (hours/days).
// It hampers certain internal cleanup tasks and is explicitly discouraged by
// the postgres maintainers - The advisory lock can silently disappear on
// postgres upgrade - Advisory locks do not play nicely with pooling tools such
// as pgbouncer - If the application crashes, the advisory lock can be left
// hanging around for a while (sometimes hours) and can require manual
// intervention to remove it
//
// For this reason, we now use a database leaseLock instead, which works as
// such: - Have one row in a database which is updated periodically with the
// client ID - CL node A will run a background process on start that updates
// this e.g. once per second - CL node B will spinlock, checking periodically
// to see if the update got too old. If it goes more than, say, 5s without
// updating, it assumes that node A is dead and takes over. Now CL node B is
// the owner of the row and it updates this every second - If CL node A comes
// back somehow, it will go to take out a lease and realise that the database
// has been leased to another process, so it will panic and quit immediately
type LeaseLock interface {
	TakeAndHold(ctx context.Context) error
	ClientID() uuid.UUID
	Release()
}

type LeaseLockConfig struct {
	DefaultQueryTimeout  time.Duration
	LeaseDuration        time.Duration
	LeaseRefreshInterval time.Duration
}

var _ LeaseLock = &leaseLock{}

type leaseLock struct {
	id         uuid.UUID
	db         *sqlx.DB
	conn       *sqlx.Conn
	cfg        LeaseLockConfig
	logger     logger.Logger
	stop       func()
	wgReleased sync.WaitGroup
}

// NewLeaseLock creates a "leaseLock" - an entity that tries to take an exclusive lease on the database
func NewLeaseLock(db *sqlx.DB, appID uuid.UUID, lggr logger.Logger, cfg LeaseLockConfig) LeaseLock {
	if cfg.LeaseRefreshInterval > cfg.LeaseDuration/2 {
		panic("refresh interval must be <= half the lease duration")
	}
	return &leaseLock{appID, db, nil, cfg, lggr.Named("LeaseLock").With("appID", appID), func() {}, sync.WaitGroup{}}
}

// TakeAndHold will block and wait indefinitely until it can get its first lock or ctx is cancelled.
// Release() function must be used to release the acquired lock.
// NOT THREAD SAFE
func (l *leaseLock) TakeAndHold(ctx context.Context) (err error) {
	l.logger.Debug("Taking initial lease...")
	retryCount := 0
	isInitial := true

	for {
		var gotLease bool
		var err error

		err = func() error {
			qctx, cancel := context.WithTimeout(ctx, l.cfg.DefaultQueryTimeout)
			defer cancel()
			if l.conn == nil {
				if err = l.checkoutConn(qctx); err != nil {
					return errors.Wrap(err, "lease lock failed to checkout initial connection")
				}
			}
			gotLease, err = l.getLease(qctx, isInitial)
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
			err = errors.Wrap(err, "failed to get lease lock")
			if l.conn != nil {
				err = multierr.Combine(err, l.conn.Close())
			}
			return err
		}
		if gotLease {
			break
		}
		isInitial = false
		l.logRetry(retryCount)
		retryCount++
		select {
		case <-ctx.Done():
			err = errors.New("stopped")
			if l.conn != nil {
				err = multierr.Combine(err, l.conn.Close())
			}
			return err
		case <-time.After(utils.WithJitter(l.cfg.LeaseRefreshInterval)):
		}
	}
	l.logger.Debug("Got exclusive lease on database")

	lctx, cancel := context.WithCancel(context.Background())
	l.stop = cancel

	l.wgReleased.Add(1)
	// Once the lock is acquired, Release() method must be used to release the lock (hence different context).
	// This is done on purpose: Release() method has exclusive control on releasing the lock.
	go l.loop(lctx)

	return nil
}

// Release requests the lock to release and blocks until it gets released.
// Calling Release for a released lock has no effect.
func (l *leaseLock) Release() {
	l.stop()
	l.wgReleased.Wait()
}

// checkout dedicated connection for lease lock to bypass any DB contention
func (l *leaseLock) checkoutConn(ctx context.Context) (err error) {
	newConn, err := l.db.Connx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed checking out connection from pool")
	}
	l.conn = newConn
	if err = l.setInitialTimeouts(ctx); err != nil {
		return multierr.Combine(
			errors.Wrap(err, "failed to set initial timeouts"),
			l.conn.Close(),
		)
	}
	return nil
}

func (l *leaseLock) setInitialTimeouts(ctx context.Context) error {
	// Set short timeouts to prevent some kind of pathological situation
	// occurring where we get stuck waiting for the table lock, or hang during
	// the transaction - we do not want to leave rows locked if this process is
	// dead
	ms := l.cfg.LeaseDuration.Milliseconds()
	return multierr.Combine(
		utils.JustError(l.conn.ExecContext(ctx, fmt.Sprintf(`SET SESSION lock_timeout = %d`, ms))),
		utils.JustError(l.conn.ExecContext(ctx, fmt.Sprintf(`SET SESSION idle_in_transaction_session_timeout = %d`, ms))),
	)
}

func (l *leaseLock) logRetry(count int) {
	if count%1000 == 0 || (count < 1000 && count&(count-1) == 0) {
		l.logger.Infow("Another application is currently holding the database lease (or a previous instance exited uncleanly), waiting for lease to expire...", "tryCount", count)
	}
}

func (l *leaseLock) loop(ctx context.Context) {
	defer l.wgReleased.Done()

	refresh := time.NewTicker(l.cfg.LeaseRefreshInterval)
	defer refresh.Stop()

	for {
		select {
		case <-ctx.Done():
			qctx, cancel := context.WithTimeout(context.Background(), l.cfg.DefaultQueryTimeout)
			err := multierr.Combine(
				utils.JustError(l.conn.ExecContext(qctx, `UPDATE lease_lock SET expires_at=NOW() WHERE client_id = $1 AND expires_at > NOW()`, l.id)),
				l.conn.Close(),
			)
			cancel()
			if err != nil {
				l.logger.Warnw("Error trying to release lease on cancelled ctx", "err", err)
			}
			return
		case <-refresh.C:
			qctx, cancel := context.WithTimeout(ctx, l.cfg.LeaseDuration)
			gotLease, err := l.getLease(qctx, false)
			if errors.Is(err, sql.ErrConnDone) {
				l.logger.Warnw("DB connection was unexpectedly closed; checking out a new one", "err", err)
				if err = l.checkoutConn(ctx); err != nil {
					l.logger.Warnw("Error trying to refresh connection", "err", err)
				}
				gotLease, err = l.getLease(ctx, false)
			}
			cancel()
			if err != nil {
				l.logger.Errorw("Error trying to refresh database lease", "err", err)
			} else if !gotLease {
				if err := l.db.Close(); err != nil {
					l.logger.Errorw("Failed to close DB", "err", err)
				}
				l.logger.Fatal("Another node has taken the lease, exiting immediately")
			}
		}
	}
}

// initialSQL is necessary because the application attempts to take the lease
// lock BEFORE running migrations
var initialSQL = []string{
	`CREATE TABLE IF NOT EXISTS lease_lock (client_id uuid NOT NULL, expires_at timestamptz NOT NULL)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS only_one_lease_lock ON lease_lock ((client_id IS NOT NULL))`,
}

// GetLease tries to get a lease from the DB
// If successful, returns true
// If the lease is currently held by someone else, returns false
// If some other error occurred, returns the error
func (l *leaseLock) getLease(ctx context.Context, isInitial bool) (gotLease bool, err error) {
	l.logger.Trace("Refreshing database lease")
	leaseDuration := fmt.Sprintf("%f seconds", l.cfg.LeaseDuration.Seconds())

	// NOTE: Uses database time for all calculations since it's conceivable
	// that node local times might be skewed compared to each other
	err = sqlutil.TransactConn(ctx, func(ds sqlutil.DataSource) sqlutil.DataSource {
		return ds
	}, l.conn, nil, func(tx sqlutil.DataSource) error {
		if isInitial {
			for _, query := range initialSQL {
				if _, err = tx.ExecContext(ctx, query); err != nil {
					return errors.Wrap(err, "failed to create initial lease_lock table")
				}
			}
		}

		// Upsert the lease_lock, only overwriting an existing one if the existing one has expired
		var res sql.Result
		res, err = tx.ExecContext(ctx, `
INSERT INTO lease_lock (client_id, expires_at) VALUES ($1, NOW()+$2::interval) ON CONFLICT ((client_id IS NOT NULL)) DO UPDATE SET
client_id = EXCLUDED.client_id,
expires_at = EXCLUDED.expires_at
WHERE
lease_lock.client_id = $1
OR
lease_lock.expires_at < NOW()
`, l.id, leaseDuration)
		if err != nil {
			return errors.Wrap(err, "failed to upsert lease_lock")
		}
		var rowsAffected int64
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "failed to get RowsAffected for lease lock upsert")
		}
		if rowsAffected > 0 {
			gotLease = true
		}
		return nil
	})
	return gotLease, errors.Wrap(err, "leaseLock#GetLease failed")
}

func (l *leaseLock) ClientID() uuid.UUID {
	return l.id
}
