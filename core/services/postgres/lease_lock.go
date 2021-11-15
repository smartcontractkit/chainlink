package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
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
	TakeAndHold() error
	Release()
	ClientID() uuid.UUID
}

var _ LeaseLock = &leaseLock{}

type leaseLock struct {
	// TODO: Use a "master" application parent ctx that is cancelled on stop?
	// https://app.shortcut.com/chainlinklabs/story/20770/application-should-have-base-context-that-cancels-on-stop
	id              uuid.UUID
	db              *sqlx.DB
	refreshInterval time.Duration
	leaseDuration   time.Duration
	logger          logger.Logger

	chStop chan struct{}
	wg     sync.WaitGroup
}

// NewLeaseLock creates a "leaseLock" - an entity that tries to take an exclusive lease on the database
func NewLeaseLock(db *sqlx.DB, appID uuid.UUID, lggr logger.Logger, refreshInterval, leaseDuration time.Duration) LeaseLock {
	if refreshInterval > leaseDuration/2 {
		panic("refresh interval must be <= half the lease duration")
	}
	return &leaseLock{appID, db, refreshInterval, leaseDuration, lggr.Named("LeaseLock").With("appID", appID), make(chan struct{}), sync.WaitGroup{}}
}

// TakeAndHold will block and wait indefinitely until it can get its first lock
func (l *leaseLock) TakeAndHold() error {
	l.logger.Debug("Taking initial lease...")
	retryCount := 0
	isInitial := true
	for {
		ctx, cancel := utils.ContextFromChan(l.chStop)
		ctx, cancel2 := DefaultQueryCtxWithParent(ctx)
		gotLease, err := l.getLease(ctx, isInitial)
		cancel2()
		cancel()
		if err != nil {
			return errors.Wrap(err, "failed to get lock")
		}
		if gotLease {
			break
		}
		isInitial = false
		l.logRetry(retryCount)
		retryCount++
		select {
		case <-l.chStop:
			return errors.New("stopped")
		case <-time.After(l.refreshInterval):
		}
	}
	l.logger.Debug("Got exclusive lease on database")
	l.wg.Add(1)
	go l.loop()
	return nil
}

func (l *leaseLock) logRetry(count int) {
	if count%1000 == 0 || count&(count-1) == 0 {
		l.logger.Infow("Another application holds the database lease, waiting...", "failCount", count+1)
	}
}

func (l *leaseLock) Release() {
	close(l.chStop)
	l.wg.Wait()
}

func (l *leaseLock) loop() {
	defer l.wg.Done()

	ticker := time.NewTicker(l.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-l.chStop:
			_, err := l.db.Exec(`UPDATE lease_lock SET expires_at=NOW() WHERE client_id = $1 AND expires_at > NOW()`, l.id)
			if err != nil {
				l.logger.Warn("Error trying to release lease on shutdown", "err", err)
			}
			return
		case <-ticker.C:
			ctx, cancel := utils.ContextFromChan(l.chStop)
			ctx, cancel2 := context.WithTimeout(ctx, l.refreshInterval)
			gotLease, err := l.getLease(ctx, false)
			cancel2()
			cancel()
			if err != nil {
				l.logger.Errorw("Error trying to refresh database lease", "err", err)
			} else if !gotLease {
				panic("another node has taken the lease")
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
	leaseDuration := fmt.Sprintf("%f seconds", l.leaseDuration.Seconds())

	// Set short timeouts to prevent some kind of pathological situation
	// occurring where we get stuck waiting for the table lock, or hang during
	// the transaction - we do not want to leave rows locked if this process is
	// dead
	opts := TxOptions{LockTimeout: l.refreshInterval, IdleInTxSessionTimeout: l.refreshInterval}

	// NOTE: Uses database time for all calculations since it's conceivable
	// that node local times might be skewed compared to each other
	err = SqlxTransaction(ctx, l.db, l.logger, func(tx Queryer) error {
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
	}, opts)
	return gotLease, errors.Wrap(err, "leaseLock#GetLease failed")
}

func (l *leaseLock) ClientID() uuid.UUID {
	return l.id
}
