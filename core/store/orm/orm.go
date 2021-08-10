package orm

import (
	"crypto/subtle"
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
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	// ErrReleaseLockFailed  is returned when releasing the advisory lock fails.
	ErrReleaseLockFailed = errors.New("advisory lock release failed")
	// ErrOptimisticUpdateConflict is returned when a record update failed
	// because another update occurred while the model was in memory and the
	// differences must be reconciled.
	ErrOptimisticUpdateConflict = errors.New("conflict while updating record")
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

func (orm *ORM) MustSQLDB() *sql.DB {
	d, err := orm.DB.DB()
	if err != nil {
		panic(err)
	}
	return d
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

// Unscoped returns a new instance of this ORM that includes soft deleted items.
func (orm *ORM) Unscoped() *ORM {
	return &ORM{
		DB:              orm.DB.Unscoped(),
		lockingStrategy: orm.lockingStrategy,
	}
}

// FindBridge looks up a Bridge by its Name.
func (orm *ORM) FindBridge(name models.TaskType) (bt models.BridgeType, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return bt, err
	}
	return bt, orm.DB.First(&bt, "name = ?", name.String()).Error
}

// ExternalInitiatorsSorted returns many ExternalInitiators sorted by Name from the store adhering
// to the passed parameters.
func (orm *ORM) ExternalInitiatorsSorted(offset int, limit int) ([]models.ExternalInitiator, int, error) {
	count, err := orm.CountOf(&models.ExternalInitiator{})
	if err != nil {
		return nil, 0, err
	}

	var exis []models.ExternalInitiator
	err = orm.getRecords(&exis, "name asc", offset, limit)
	return exis, count, err
}

// CreateExternalInitiator inserts a new external initiator
func (orm *ORM) CreateExternalInitiator(externalInitiator *models.ExternalInitiator) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	err := orm.DB.Create(externalInitiator).Error
	return err
}

// DeleteExternalInitiator removes an external initiator
func (orm *ORM) DeleteExternalInitiator(name string) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	err := orm.DB.Exec("DELETE FROM external_initiators WHERE name = ?", name).Error
	return err
}

// FindExternalInitiator finds an external initiator given an authentication request
func (orm *ORM) FindExternalInitiator(
	eia *auth.Token,
) (*models.ExternalInitiator, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, err
	}
	initiator := &models.ExternalInitiator{}
	err := orm.DB.Where("access_key = ?", eia.AccessKey).First(initiator).Error
	if err != nil {
		return nil, errors.Wrap(err, "error finding external initiator")
	}

	return initiator, nil
}

// FindExternalInitiatorByName finds an external initiator given an authentication request
func (orm *ORM) FindExternalInitiatorByName(iname string) (exi models.ExternalInitiator, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return exi, err
	}
	return exi, orm.DB.First(&exi, "lower(name) = lower(?)", iname).Error
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
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, err
	}
	ethTxAttempt := &bulletprooftxmanager.EthTxAttempt{}
	if err := orm.DB.Preload("EthTx").First(ethTxAttempt, "hash = ?", hash).Error; err != nil {
		return nil, errors.Wrap(err, "FindEthTxAttempt First(ethTxAttempt) failed")
	}
	return ethTxAttempt, nil
}

// FindUser will return the one API user, or an error.
func (orm *ORM) FindUser() (models.User, error) {
	return findUser(orm.DB)
}

func findUser(db *gorm.DB) (user models.User, err error) {
	return user, db.Preload(clause.Associations).Order("created_at desc").First(&user).Error
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (orm *ORM) AuthorizedUserWithSession(sessionID string, sessionDuration time.Duration) (models.User, error) {
	if len(sessionID) == 0 {
		return models.User{}, errors.New("Session ID cannot be empty")
	}

	var session models.Session
	err := orm.DB.First(&session, "id = ?", sessionID).Error
	if err != nil {
		return models.User{}, err
	}
	now := time.Now()
	if session.LastUsed.Add(sessionDuration).Before(now) {
		return models.User{}, errors.New("Session has expired")
	}
	session.LastUsed = now
	if err := orm.DB.Save(&session).Error; err != nil {
		return models.User{}, err
	}
	return orm.FindUser()
}

// DeleteUser will delete the API User in the db.
func (orm *ORM) DeleteUser() error {
	return postgres.GormTransactionWithDefaultContext(orm.DB, func(dbtx *gorm.DB) error {
		user, err := findUser(dbtx)
		if err != nil {
			return err
		}

		if err = dbtx.Delete(&user).Error; err != nil {
			return err
		}

		return dbtx.Exec("DELETE FROM sessions").Error
	})
}

// DeleteUserSession will erase the session ID for the sole API User.
func (orm *ORM) DeleteUserSession(sessionID string) error {
	return orm.DB.Delete(models.Session{ID: sessionID}).Error
}

// DeleteBridgeType removes the bridge type
func (orm *ORM) DeleteBridgeType(bt *models.BridgeType) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.Delete(bt).Error
}

// CreateSession will check the password in the SessionRequest against
// the hashed API User password in the db.
func (orm *ORM) CreateSession(sr models.SessionRequest) (string, error) {
	user, err := orm.FindUser()
	if err != nil {
		return "", err
	}

	if !constantTimeEmailCompare(sr.Email, user.Email) {
		return "", errors.New("Invalid email")
	}

	if utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		session := models.NewSession()
		return session.ID, orm.DB.Save(&session).Error
	}
	return "", errors.New("Invalid password")
}

const constantTimeEmailLength = 256

func constantTimeEmailCompare(left, right string) bool {
	length := utils.MaxInt(constantTimeEmailLength, len(left), len(right))
	leftBytes := make([]byte, length)
	rightBytes := make([]byte, length)
	copy(leftBytes, left)
	copy(rightBytes, right)
	return subtle.ConstantTimeCompare(leftBytes, rightBytes) == 1
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (orm *ORM) ClearNonCurrentSessions(sessionID string) error {
	return orm.DB.Delete(&models.Session{}, "id != ?", sessionID).Error
}

// BridgeTypes returns bridge types ordered by name filtered limited by the
// passed params.
func (orm *ORM) BridgeTypes(offset int, limit int) ([]models.BridgeType, int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, 0, err
	}
	count, err := orm.CountOf(&models.BridgeType{})
	if err != nil {
		return nil, 0, err
	}

	var bridges []models.BridgeType
	err = orm.getRecords(&bridges, "name asc", offset, limit)
	return bridges, count, err
}

// SaveUser saves the user.
func (orm *ORM) SaveUser(user *models.User) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.Save(user).Error
}

// CreateBridgeType saves the bridge type.
func (orm *ORM) CreateBridgeType(bt *models.BridgeType) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.Create(bt).Error
}

// UpdateBridgeType updates the bridge type.
func (orm *ORM) UpdateBridgeType(bt *models.BridgeType, btr *models.BridgeTypeRequest) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	bt.URL = btr.URL
	bt.Confirmations = btr.Confirmations
	bt.MinimumContractPayment = btr.MinimumContractPayment
	return orm.DB.Save(bt).Error
}

func (orm *ORM) CountOf(t interface{}) (int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return 0, err
	}
	var count int64
	return int(count), orm.DB.Model(t).Count(&count).Error
}

func (orm *ORM) getRecords(collection interface{}, order string, offset, limit int) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.
		Preload(clause.Associations).
		Order(order).Limit(limit).Offset(offset).
		Find(collection).Error
}

func (orm *ORM) RawDBWithAdvisoryLock(fn func(*gorm.DB) error) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return fn(orm.DB)
}

// Connection manages all of the possible database connection setup and config.
type Connection struct {
	name               dialects.DialectName
	uri                string
	dialect            dialects.DialectName
	locking            bool
	advisoryLockID     int64
	lockRetryInterval  time.Duration
	transactionWrapped bool
	maxOpenConns       int
	maxIdleConns       int
}

// NewConnection returns a Connection which holds all of the configuration
// necessary for managing the database connection.
func NewConnection(dialect dialects.DialectName, uri string, advisoryLockID int64, lockRetryInterval time.Duration, maxOpenConns, maxIdleConns int) (Connection, error) {
	switch dialect {
	case dialects.Postgres:
		return Connection{
			advisoryLockID:     advisoryLockID,
			dialect:            dialects.Postgres,
			locking:            true,
			name:               dialect,
			transactionWrapped: false,
			uri:                uri,
			lockRetryInterval:  lockRetryInterval,
			maxOpenConns:       maxOpenConns,
			maxIdleConns:       maxIdleConns,
		}, nil
	case dialects.PostgresWithoutLock:
		return Connection{
			advisoryLockID:     advisoryLockID,
			dialect:            dialects.Postgres,
			locking:            false,
			name:               dialect,
			transactionWrapped: false,
			uri:                uri,
			maxOpenConns:       maxOpenConns,
			maxIdleConns:       maxIdleConns,
		}, nil
	case dialects.TransactionWrappedPostgres:
		return Connection{
			advisoryLockID:     advisoryLockID,
			dialect:            dialects.TransactionWrappedPostgres,
			locking:            true,
			name:               dialect,
			transactionWrapped: true,
			uri:                uri,
			lockRetryInterval:  lockRetryInterval,
			maxOpenConns:       maxOpenConns,
			maxIdleConns:       maxIdleConns,
		}, nil
	}
	return Connection{}, errors.Errorf("%s is not a valid dialect type", dialect)
}

func (ct *Connection) initializeDatabase() (*gorm.DB, error) {
	originalURI := ct.uri
	if ct.transactionWrapped {
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

// SortType defines the different sort orders available.
type SortType int

const (
	// Ascending is the sort order going up, i.e. 1,2,3.
	Ascending SortType = iota
	// Descending is the sort order going down, i.e. 3,2,1.
	Descending
)

func (s SortType) String() string {
	orderStr := "asc"
	if s == Descending {
		orderStr = "desc"
	}
	return orderStr
}
