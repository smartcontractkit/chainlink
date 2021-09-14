package orm

import (
	"database/sql"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"gorm.io/gorm/clause"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"go.uber.org/multierr"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	// We've specified a later version in go.mod than is currently used by gorm
	// to get this fix in https://github.com/jackc/pgx/pull/975.
	// As soon as pgx releases a 4.12 and gorm [https://github.com/go-gorm/postgres/blob/master/go.mod#L6]
	// bumps their version to 4.12, we can remove this.
	_ "github.com/jackc/pgx/v4"
)

var (
	// ErrorNotFound is returned when finding a single value fails.
	ErrorNotFound = gorm.ErrRecordNotFound
	// ErrNoAdvisoryLock is returned when an advisory lock can't be acquired.
	ErrNoAdvisoryLock = errors.New("can't acquire advisory lock")
)

// ORM contains the database object used by Chainlink.
type ORM struct {
	DB                  *gorm.DB
	lockingStrategy     LockingStrategy
	advisoryLockTimeout models.Duration
	closeOnce           sync.Once
	shutdownSignal      gracefulpanic.Signal
}

// NewORM initializes the orm with the configured uri
func NewORM(uri string, timeout models.Duration, shutdownSignal gracefulpanic.Signal, dialect dialects.DialectName, advisoryLockID int64, lockRetryInterval time.Duration, maxOpenConns, maxIdleConns int) (*ORM, error) {
	ct, err := NewConnection(dialect, uri, advisoryLockID, lockRetryInterval, maxOpenConns, maxIdleConns)
	if err != nil {
		return nil, err
	}
	// Locking strategy for transaction wrapped postgres must use original URI
	lockingStrategy, err := NewLockingStrategy(ct)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create ORM lock")
	}

	orm := &ORM{
		lockingStrategy:     lockingStrategy,
		advisoryLockTimeout: timeout,
		shutdownSignal:      shutdownSignal,
	}

	db, err := ct.initializeDatabase()
	if err != nil {
		return nil, errors.Wrap(err, "unable to init DB")
	}
	orm.DB = db

	return orm, nil
}

// MustEnsureAdvisoryLock sends a shutdown signal to the ORM if it an advisory
// lock cannot be acquired.
func (orm *ORM) MustEnsureAdvisoryLock() error {
	err := orm.lockingStrategy.Lock(orm.advisoryLockTimeout)
	if err != nil {
		logger.Errorf("unable to lock ORM: %v", err)
		orm.shutdownSignal.Panic()
		return err
	}
	return nil
}

func displayTimeout(timeout models.Duration) string {
	if timeout.IsInstant() {
		return "indefinite"
	}
	return timeout.String()
}

// SetLogging turns on SQL statement logging
func (orm *ORM) SetLogging(enabled bool) {
	orm.DB.Logger = NewOrmLogWrapper(logger.Default, enabled, time.Second)
}

// Close closes the underlying database connection.
func (orm *ORM) Close() error {
	var err error
	db, _ := orm.DB.DB()
	orm.closeOnce.Do(func() {
		err = multierr.Combine(
			db.Close(),
			orm.lockingStrategy.Unlock(orm.advisoryLockTimeout),
		)
	})
	return err
}

// EthTransactionsWithAttempts returns all eth transactions with at least one attempt
// limited by passed parameters. Attempts are sorted by created_at.
func (orm *ORM) EthTransactionsWithAttempts(offset, limit int) ([]bulletprooftxmanager.EthTx, int, error) {
	ethTXIDs := orm.DB.
		Select("DISTINCT eth_tx_id").
		Table("eth_tx_attempts")

	var count int64
	err := orm.DB.
		Table("eth_txes").
		Where("id IN (?)", ethTXIDs).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var txs []bulletprooftxmanager.EthTx
	err = orm.DB.
		Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Where("id IN (?)", ethTXIDs).
		Order("id desc").Limit(limit).Offset(offset).
		Find(&txs).Error

	return txs, int(count), err
}

// EthTxAttempts returns the last tx attempts sorted by created_at descending.
func (orm *ORM) EthTxAttempts(offset, limit int) ([]bulletprooftxmanager.EthTxAttempt, int, error) {
	count, err := orm.CountOf(&bulletprooftxmanager.EthTxAttempt{})
	if err != nil {
		return nil, 0, err
	}

	var attempts []bulletprooftxmanager.EthTxAttempt
	err = orm.DB.
		Preload("EthTx").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&attempts).Error

	return attempts, count, err
}

// FindEthTxAttempt returns an individual EthTxAttempt
func (orm *ORM) FindEthTxAttempt(hash common.Hash) (*bulletprooftxmanager.EthTxAttempt, error) {
	ethTxAttempt := &bulletprooftxmanager.EthTxAttempt{}
	if err := orm.DB.Preload("EthTx").First(ethTxAttempt, "hash = ?", hash).Error; err != nil {
		return nil, errors.Wrap(err, "FindEthTxAttempt First(ethTxAttempt) failed")
	}
	return ethTxAttempt, nil
}

func (orm *ORM) CountOf(t interface{}) (int, error) {
	var count int64
	return int(count), orm.DB.Model(t).Count(&count).Error
}

func (orm *ORM) RawDBWithAdvisoryLock(fn func(*gorm.DB) error) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return fn(orm.DB)
}

// Connection manages all of the possible database connection setup and config.
type Connection struct {
	name              dialects.DialectName
	uri               string
	dialect           dialects.DialectName
	locking           bool
	advisoryLockID    int64
	lockRetryInterval time.Duration
	maxOpenConns      int
	maxIdleConns      int
}

// NewConnection returns a Connection which holds all of the configuration
// necessary for managing the database connection.
func NewConnection(dialect dialects.DialectName, uri string, advisoryLockID int64, lockRetryInterval time.Duration, maxOpenConns, maxIdleConns int) (Connection, error) {
	ct := Connection{
		advisoryLockID: advisoryLockID,
		locking:        true,
		name:           dialect,
		uri:            uri,
		maxOpenConns:   maxOpenConns,
		maxIdleConns:   maxIdleConns,
	}
	switch dialect {
	case dialects.Postgres:
		ct.dialect = dialects.Postgres
		ct.locking = true
		ct.lockRetryInterval = lockRetryInterval
	case dialects.PostgresWithoutLock:
		ct.dialect = dialects.Postgres
		ct.locking = false
	case dialects.TransactionWrappedPostgres:
		ct.dialect = dialects.TransactionWrappedPostgres
		ct.locking = true
		ct.lockRetryInterval = lockRetryInterval
	default:
		return Connection{}, errors.Errorf("%s is not a valid dialect type", dialect)
	}
	return ct, nil
}

func (ct *Connection) initializeDatabase() (*gorm.DB, error) {
	originalURI := ct.uri
	if ct.dialect == dialects.TransactionWrappedPostgres {
		// Dbtx uses the uri as a unique identifier for each transaction. Each ORM
		// should be encapsulated in it's own transaction, and thus needs its own
		// unique id.
		//
		// We can happily throw away the original uri here because if we are using
		// txdb it should have already been set at the point where we called
		// txdb.Register
		ct.uri = uuid.NewV4().String()
	} else {
		uri, err := url.Parse(ct.uri)
		if err != nil {
			return nil, err
		}
		static.SetConsumerName(uri, "ORM")
		ct.uri = uri.String()
	}

	newLogger := NewOrmLogWrapper(logger.Default, false, time.Second)

	// Use the underlying connection with the unique uri for txdb.
	d, err := sql.Open(string(ct.dialect), ct.uri)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, errors.Errorf("unable to open %s received a nil connection", ct.uri)
	}
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn: d,
		DSN:  originalURI,
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %s for gorm DB conn %v", ct.uri, d)
	}
	db = db.Omit(clause.Associations).Session(&gorm.Session{})

	if err = db.Exec(`SET TIME ZONE 'UTC'`).Error; err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(ct.maxOpenConns)
	d.SetMaxIdleConns(ct.maxIdleConns)

	return db, nil
}
