package pg

import (
	"context"
	"net/url"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

// LockedDB bounds DB connection and DB locks.
type LockedDB interface {
	Open(ctx context.Context) error
	Close() error
	DB() *sqlx.DB
}

type LockedDBConfig interface {
	ConnectionConfig
	AppID() uuid.UUID
	DatabaseLockingMode() string
	DatabaseURL() url.URL
	DatabaseDefaultQueryTimeout() time.Duration
	LeaseLockDuration() time.Duration
	LeaseLockRefreshInterval() time.Duration
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	MigrateDatabase() bool
}

type lockedDb struct {
	cfg           LockedDBConfig
	lggr          logger.Logger
	db            *sqlx.DB
	leaseLock     LeaseLock
	statsReporter *StatsReporter
}

// NewLockedDB creates a new instance of LockedDB.
func NewLockedDB(cfg LockedDBConfig, lggr logger.Logger) LockedDB {
	return &lockedDb{
		cfg:  cfg,
		lggr: lggr.Named("LockedDB"),
	}
}

// OpenUnlockedDB just opens DB connection, without any DB locks.
// This should be used carefully, when we know we don't need any locks.
// Currently this is used by RebroadcastTransactions command only.
func OpenUnlockedDB(cfg LockedDBConfig) (db *sqlx.DB, err error) {
	return openDB(cfg)
}

// Open function connects to DB and acquires DB locks based on configuration.
// If any of the steps fails or ctx is cancelled, it reverts everything.
// This is a blocking function and it may execute long due to DB locks acquisition.
// NOT THREAD SAFE
func (l *lockedDb) Open(ctx context.Context) (err error) {
	// If Open succeeded previously, db will not be nil
	if l.db != nil {
		l.lggr.Panic("calling Open() twice")
	}

	// Step 1: open DB connection
	l.db, err = openDB(l.cfg)
	if err != nil {
		// l.db will be nil in case of error
		return errors.Wrap(err, "failed to open db")
	}
	revert := func() {
		// Let Open() return the actual error, while l.Close() error is just logged.
		if err2 := l.Close(); err2 != nil {
			l.lggr.Errorf("failed to cleanup LockedDB: %v", err2)
		}
	}

	// Step 2: start the stat reporter
	l.statsReporter = NewStatsReporter(l.db.Stats, l.lggr)
	l.statsReporter.Start(ctx)

	// Step 3: acquire DB locks
	lockingMode := l.cfg.DatabaseLockingMode()
	l.lggr.Debugf("Using database locking mode: %s", lockingMode)

	// Take the lease before any other DB operations
	switch lockingMode {
	case "lease":
		l.leaseLock = NewLeaseLock(l.db, l.cfg.AppID(), l.lggr, l.cfg)
		if err = l.leaseLock.TakeAndHold(ctx); err != nil {
			defer revert()
			return errors.Wrap(err, "failed to take initial lease on database")
		}
	}

	return
}

// Close function releases DB locks (if acquired by Open) and closes DB connection.
// Closing of a closed LockedDB instance has no effect.
// NOT THREAD SAFE
func (l *lockedDb) Close() error {
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
func (l lockedDb) DB() *sqlx.DB {
	return l.db
}

func openDB(cfg LockedDBConfig) (db *sqlx.DB, err error) {
	uri := cfg.DatabaseURL()
	appid := cfg.AppID()
	static.SetConsumerName(&uri, "App", &appid)
	dialect := cfg.GetDatabaseDialectConfiguredOrDefault()
	db, err = NewConnection(uri.String(), dialect, cfg)
	return
}
