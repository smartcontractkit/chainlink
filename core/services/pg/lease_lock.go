package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"
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
}

var _ LeaseLock = &leaseLock{}

type leaseLock struct {
	id              uuid.UUID
	db              *sqlx.DB
	conn            *sqlx.Conn
	refreshInterval time.Duration
	leaseDuration   time.Duration
	logger          logger.Logger
}

// NewLeaseLock creates a "leaseLock" - an entity that tries to take an exclusive lease on the database
func NewLeaseLock(db *sqlx.DB, appID uuid.UUID, lggr logger.Logger, refreshInterval, leaseDuration time.Duration) LeaseLock {
	if refreshInterval > leaseDuration/2 {
		panic("refresh interval must be <= half the lease duration")
	}
	return &leaseLock{appID, db, nil, refreshInterval, leaseDuration, lggr.Named("LeaseLock").With("appID", appID)}
}

// TakeAndHold will block and wait indefinitely until it can get its first lock
// NOT THREAD SAFE
func (l *leaseLock) TakeAndHold(ctx context.Context) (err error) {
	l.logger.Debug("Taking initial lease...")
	retryCount := 0
	isInitial := true

	for {
		var gotLease bool
		var err error

		err = func() error {
			qCtx, cancel := DefaultQueryCtxWithParent(ctx)
			defer cancel()
			if l.conn == nil {
				if err = l.checkoutConn(qCtx); err != nil {
					return errors.Wrap(err, "lease lock failed to checkout initial connection")
				}
			}
			gotLease, err = l.getLease(qCtx, isInitial)
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
			return multierr.Combine(
				errors.New("application shutdown"),
				l.conn.Close(),
			)
		case <-time.After(utils.WithJitter(l.refreshInterval)):
		}
	}
	l.logger.Debug("Got exclusive lease on database")
	go l.loop(ctx)
	return nil
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
	return multierr.Combine(
		utils.JustError(l.conn.ExecContext(ctx, fmt.Sprintf(`SET SESSION lock_timeout = %d`, l.leaseDuration.Milliseconds()))),
		utils.JustError(l.conn.ExecContext(ctx, fmt.Sprintf(`SET SESSION idle_in_transaction_session_timeout = %d`, l.leaseDuration.Milliseconds()))),
	)
}

func (l *leaseLock) logRetry(count int) {
	if count%1000 == 0 || (count < 1000 && count&(count-1) == 0) {
		l.logger.Infow("Another application is currently holding the database lease (or a previous instance exited uncleanly), waiting for lease to expire...", "tryCount", count)
	}
}

func (l *leaseLock) loop(ctx context.Context) {
	ticker := time.NewTicker(l.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			qCtx, cancel := DefaultQueryCtx()
			err := multierr.Combine(
				utils.JustError(l.conn.ExecContext(qCtx, `UPDATE lease_lock SET expires_at=NOW() WHERE client_id = $1 AND expires_at > NOW()`, l.id)),
				l.conn.Close(),
			)
			cancel()
			if err != nil {
				l.logger.Warnw("Error trying to release lease on shutdown", "err", err)
			}
			return
		case <-ticker.C:
			qCtx, cancel := context.WithTimeout(context.Background(), l.leaseDuration)
			gotLease, err := l.getLease(qCtx, false)
			if errors.Is(err, sql.ErrConnDone) {
				l.logger.Warnw("DB connection was unexpectedly closed; checking out a new one", "err", err)
				if err = l.checkoutConn(qCtx); err != nil {
					l.logger.Warnw("Error trying to refresh connection", "err", err)
				}
				gotLease, err = l.getLease(qCtx, false)
			}
			cancel()
			if err != nil {
				l.logger.Errorw("Error trying to refresh database lease", "err", err)
			} else if !gotLease {
				l.logger.Fatal("Another node has taken the lease, exiting")
			}
		}
	}
}

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
	leaseDuration := fmt.Sprintf("%f seconds", l.leaseDuration.Seconds())

	// NOTE: Uses database time for all calculations since it's conceivable
	// that node local times might be skewed compared to each other
	err = sqlxTransactionQ(ctx, l.conn, l.logger, func(tx Queryer) error {
		if isInitial {
			for _, query := range initialSQL {
				if _, err = tx.Exec(query); err != nil {
					return errors.Wrap(err, "failed to create initial lease_lock table")
				}
			}
		}
		if _, err = tx.Exec(`LOCK TABLE lease_lock`); err != nil {
			return errors.Wrap(err, "failed to lock lease_lock table")
		}
		var count int
		err = tx.Get(&count, `SELECT count(*) FROM lease_lock`)
		if count == 0 {
			// first time anybody claimed a lock on this table
			_, err = tx.Exec(`INSERT INTO lease_lock (client_id, expires_at) VALUES ($1, NOW()+$2::interval)`, l.id, leaseDuration)
			gotLease = true
			return errors.Wrap(err, "failed to create initial lease_lock")
		} else if count > 1 {
			return errors.Errorf("expected only one row in lease_lock, got %d", count)
		}
		var res sql.Result
		res, err = tx.Exec(`
UPDATE lease_lock
SET client_id = $1, expires_at = NOW()+$2::interval
WHERE (
	lease_lock.client_id = $1
	OR
	lease_lock.expires_at < NOW()
)`, l.id, leaseDuration)
		if err != nil {
			return errors.Wrap(err, "failed to update lease_lock")
		}
		var rowsAffected int64
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "failed to update lease_lock (RowsAffected)")
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
