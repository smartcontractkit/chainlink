package pg

import (
	"context"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

// LockedDB bounds DB connection and DB locks.
type LockedDB interface {
	Open(ctx context.Context) error
	Close() error
	DB() *sqlx.DB
}

type LockedDBConfig interface {
	ConnectionConfig
	URL() url.URL
	DefaultQueryTimeout() time.Duration
	DriverName() string
}

type lockedDb struct {
	services.StateMachine
	appID         uuid.UUID
	cfg           LockedDBConfig
	lockCfg       config.Lock
	lggr          logger.Logger
	db            *sqlx.DB
	leaseLock     LeaseLock
	statsReporter *StatsReporter
}

// NewLockedDB creates a new instance of LockedDB.
func NewLockedDB(appID uuid.UUID, cfg LockedDBConfig, lockCfg config.Lock, lggr logger.Logger) LockedDB {
	return &lockedDb{
		appID:   appID,
		cfg:     cfg,
		lockCfg: lockCfg,
		lggr:    lggr.Named("LockedDB"),
	}
}

// OpenUnlockedDB just opens DB connection, without any DB locks.
// This should be used carefully, when we know we don't need any locks.
// Currently this is used by RebroadcastTransactions command only.
func OpenUnlockedDB(ctx context.Context, appID uuid.UUID, cfg LockedDBConfig) (db *sqlx.DB, err error) {
	return openDB(ctx, appID, cfg)
}

// Open function connects to DB and acquires DB locks based on configuration.
// If any of the steps fails or ctx is cancelled, it reverts everything.
// This is a blocking function and it may execute long due to DB locks acquisition.
func (l *lockedDb) Open(ctx context.Context) (err error) {
	return l.StartOnce("LockedDB", func() error {
		// Step 1: open DB connection
		l.db, err = openDB(ctx, l.appID, l.cfg)
		if err != nil {
			// l.db will be nil in case of error
			return errors.Wrap(err, "failed to open db")
		}
		revert := func() {
			// Let Open() return the actual error, while l.Close() error is just logged.
			if err2 := l.close(); err2 != nil {
				l.lggr.Errorf("failed to cleanup LockedDB: %v", err2)
			}
		}

		// Step 2: start the stat reporter
		l.statsReporter = NewStatsReporter(l.db.Stats, l.lggr)
		l.statsReporter.Start()

		// Step 3: acquire DB locks
		lockingMode := l.lockCfg.LockingMode()
		l.lggr.Debugf("Using database locking mode: %s", lockingMode)

		// Take the lease before any other DB operations
		switch lockingMode {
		case "lease":
			cfg := LeaseLockConfig{
				DefaultQueryTimeout:  l.cfg.DefaultQueryTimeout(),
				LeaseDuration:        l.lockCfg.LeaseDuration(),
				LeaseRefreshInterval: l.lockCfg.LeaseRefreshInterval(),
			}
			l.leaseLock = NewLeaseLock(l.db, l.appID, l.lggr, cfg)
			if err = l.leaseLock.TakeAndHold(ctx); err != nil {
				defer revert()
				return errors.Wrap(err, "failed to take initial lease on database")
			}
		}

		return nil
	})
}

// Close function releases DB locks (if acquired by Open) and closes DB connection.
// Closing of a closed LockedDB instance has no effect.
func (l *lockedDb) Close() error {
	err := l.StopOnce("LockedDB", func() error {
		return l.close()
	})
	if !errors.Is(err, services.ErrAlreadyStopped) {
		return err
	}
	return nil
}

func (l *lockedDb) close() error {
	defer func() {
		l.db = nil
		l.leaseLock = nil
		l.statsReporter = nil
	}()

	// Step 0: stop the stat reporter
	if l.statsReporter != nil {
		l.statsReporter.Stop()
	}

	// Step 1: release DB locks
	if l.leaseLock != nil {
		l.leaseLock.Release()
	}

	// Step 2: close DB connection
	if l.db != nil {
		return l.db.Close()
	}

	return nil
}

// DB returns DB connection if Opened successfully, or nil.
func (l *lockedDb) DB() *sqlx.DB {
	return l.db
}

func openDB(ctx context.Context, appID uuid.UUID, cfg LockedDBConfig) (db *sqlx.DB, err error) {
	uri := cfg.URL()
	static.SetConsumerName(&uri, "App", &appID)
	db, err = NewConnection(ctx, uri.String(), cfg.DriverName(), cfg)
	return
}
