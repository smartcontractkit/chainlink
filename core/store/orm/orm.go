package orm

import (
	"crypto/subtle"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // http://doc.gorm.io/database.html#connecting-to-a-database
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
)

var (
	// ErrorNotFound is returned when finding a single value fails.
	ErrorNotFound = gorm.ErrRecordNotFound
)

// DialectName is a compiler enforced type used that maps to gorm's dialect
// names.
type DialectName string

const (
	// DialectPostgres represents the postgres dialect.
	DialectPostgres DialectName = "postgres"
	// DialectSqlite represents the sqlite dialect.
	DialectSqlite = "sqlite3"
)

// ORM contains the database object used by Chainlink.
type ORM struct {
	DB              *gorm.DB
	lockingStrategy LockingStrategy
	dialectName     DialectName
}

// NewORM initializes a new database file at the configured uri.
func NewORM(uri string, timeout time.Duration) (*ORM, error) {
	dialect, err := DeduceDialect(uri)
	if err != nil {
		return nil, err
	}

	lockingStrategy, err := NewLockingStrategy(dialect, uri)
	if err != nil {
		return nil, fmt.Errorf("unable to create ORM lock: %+v", err)
	}

	logger.Infof("Locking %v for exclusive access with %v timeout", dialect, displayTimeout(timeout))
	err = lockingStrategy.Lock(timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to lock ORM: %+v", err)
	}

	db, err := initializeDatabase(string(dialect), uri)
	if err != nil {
		return nil, fmt.Errorf("unable to init DB: %+v", err)
	}

	orm := &ORM{
		DB:              db,
		lockingStrategy: lockingStrategy,
		dialectName:     dialect,
	}

	return orm, nil
}

func displayTimeout(timeout time.Duration) string {
	if timeout == 0 {
		return "indefinite"
	}
	return timeout.String()
}

func initializeDatabase(dialect, path string) (*gorm.DB, error) {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s for gorm DB: %+v", path, err)
	}

	db.SetLogger(ormLogWrapper{logger.GetLogger()})

	if err := dbutil.SetTimezone(db); err != nil {
		return nil, err
	}

	if err := dbutil.SetSqlitePragmas(db); err != nil {
		return nil, err
	}

	if err := dbutil.LimitSqliteOpenConnections(db); err != nil {
		return nil, err
	}

	return db, nil
}

// DeduceDialect returns the appropriate dialect for the passed connection string.
func DeduceDialect(path string) (DialectName, error) {
	url, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	scheme := strings.ToLower(url.Scheme)
	switch scheme {
	case "postgresql", "postgres":
		return DialectPostgres, nil
	case "file", "":
		if len(strings.Split(url.Path, " ")) > 1 {
			return "", errors.New("error deducing ORM dialect, no spaces allowed, please use a postgres URL or file path")
		}
		return DialectSqlite, nil
	case "sqlite3", "sqlite":
		return "", fmt.Errorf("do not have full support for the sqlite URL, please use file:// instead for path %s", path)
	}

	return DialectSqlite, nil
}

func ignoreRecordNotFound(db *gorm.DB) error {
	var merr error
	for _, e := range db.GetErrors() {
		if e != gorm.ErrRecordNotFound {
			merr = multierr.Append(merr, e)
		}
	}
	return merr
}

func (orm *ORM) DialectName() DialectName {
	return orm.dialectName
}

// EnableLogging turns on SQL statement logging
func (orm *ORM) EnableLogging() {
	orm.DB.LogMode(true)
}

// Close closes the underlying database connection.
func (orm *ORM) Close() error {
	return multierr.Append(
		orm.DB.Close(),
		orm.lockingStrategy.Unlock(),
	)
}

// Unscoped returns a new instance of this ORM that includes soft deleted items.
func (orm *ORM) Unscoped() *ORM {
	return &ORM{
		DB:              orm.DB.Unscoped(),
		lockingStrategy: orm.lockingStrategy,
	}
}

// Where fetches multiple objects with "Find".
func (orm *ORM) Where(field string, value interface{}, instance interface{}) error {
	return orm.DB.Where(fmt.Sprintf("%v = ?", field), value).Find(instance).Error
}

// FindBridge looks up a Bridge by its Name.
func (orm *ORM) FindBridge(name models.TaskType) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, orm.DB.First(&bt, "name = ?", name.String()).Error
}

// PendingBridgeType returns the bridge type of the current pending task,
// or error if not pending bridge.
func (orm *ORM) PendingBridgeType(jr models.JobRun) (models.BridgeType, error) {
	nextTask := jr.NextTaskRun()
	if nextTask == nil {
		return models.BridgeType{}, errors.New("Cannot find the pending bridge type of a job run with no unfinished tasks")
	}
	return orm.FindBridge(nextTask.TaskSpec.Type)
}

// FindJob looks up a Job by its ID.
func (orm *ORM) FindJob(id string) (models.JobSpec, error) {
	var job models.JobSpec
	return job, orm.preloadJobs().First(&job, "id = ?", id).Error
}

// FindInitiator returns the single initiator defined by the passed ID.
func (orm *ORM) FindInitiator(ID uint) (models.Initiator, error) {
	initr := models.Initiator{}
	return initr, orm.DB.
		Set("gorm:auto_preload", true).
		First(&initr, "id = ?", ID).Error
}

func (orm *ORM) preloadJobs() *gorm.DB {
	return orm.DB.
		Preload("Initiators", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Order("\"id\" asc")
		}).
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Order("id asc")
		})
}

func preloadTaskRuns(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Result").
		Preload("TaskSpec", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		})
}

func (orm *ORM) preloadJobRuns() *gorm.DB {
	return orm.DB.
		Preload("Initiator", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("RunRequest").
		Preload("Overrides").
		Preload("TaskRuns", func(db *gorm.DB) *gorm.DB {
			return preloadTaskRuns(db).Order("task_spec_id asc")
		}).
		Preload("Result")
}

// FindJobRun looks up a JobRun by its ID.
func (orm *ORM) FindJobRun(id string) (models.JobRun, error) {
	var jr models.JobRun
	err := orm.preloadJobRuns().First(&jr, "id = ?", id).Error
	return jr, err
}

// AllSyncEvents returns all sync events
func (orm *ORM) AllSyncEvents(cb func(*models.SyncEvent) error) error {
	return Batch(1000, func(offset, limit uint) (uint, error) {
		var events []models.SyncEvent
		err := orm.DB.
			Limit(limit).
			Offset(offset).
			Order("id, created_at asc").
			Find(&events).Error
		if err != nil {
			return 0, err
		}

		for _, event := range events {
			err = cb(&event)
			if err != nil {
				return 0, err
			}
		}

		return uint(len(events)), err
	})
}

// convenientTransaction handles setup and teardown for a gorm database
// transaction, handing off the database transaction to the callback parameter.
// Encourages the use of transactions for gorm calls that translate
// into multiple sql calls, i.e. orm.SaveJobRun(run), which are better suited
// in a database transaction.
// Improves efficiency in sqlite by preventing autocommit on each line, instead
// Batch committing at the end of the transaction.
func (orm *ORM) convenientTransaction(callback func(*gorm.DB) error) error {
	dbtx := orm.DB.Begin()
	if dbtx.Error != nil {
		return dbtx.Error
	}

	err := callback(dbtx)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	return dbtx.Commit().Error
}

// SaveJobRun updates UpdatedAt for a JobRun and saves it
func (orm *ORM) SaveJobRun(run *models.JobRun) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return dbtx.Unscoped().Omit("deleted_at").Save(run).Error
	})
}

// CreateJobRun inserts a new JobRun
func (orm *ORM) CreateJobRun(run *models.JobRun) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return updateDeletedAtFromJobSpec(run, dbtx.Create(run)).Error
	})
}

// CreateExternalInitiator inserts a new external initiator
func (orm *ORM) CreateExternalInitiator(externalInitiator *models.ExternalInitiator) error {
	return orm.DB.Create(externalInitiator).Error
}

// DeleteExternalInitiator removes an external initiator
func (orm *ORM) DeleteExternalInitiator(accessKey string) error {
	return orm.DB.
		Where("access_key = ?", accessKey).
		Delete(&models.ExternalInitiator{}).
		Error
}

// FindExternalInitiator finds an external initiator given an authentication request
func (orm *ORM) FindExternalInitiator(eia *models.ExternalInitiatorAuthentication) (*models.ExternalInitiator, error) {
	initiator := &models.ExternalInitiator{}
	err := orm.DB.Where("access_key = ?", eia.AccessKey).Find(initiator).Error
	if err != nil {
		return nil, errors.Wrap(err, "error finding external initiator")
	}

	return initiator, nil
}

// updateDeletedAtFromJobSpec will update a runs deleted_at from its parent
// job spec.
// This is of particular importance in the edge case when a runlog starts a run
// at the same time a job is archived. Said run will never retain a deleted_at,
// unless this is invoked.
func updateDeletedAtFromJobSpec(run *models.JobRun, db *gorm.DB) *gorm.DB {
	return db.Exec(fmt.Sprintf(`
		UPDATE job_runs SET deleted_at = (
			SELECT job_specs.deleted_at FROM job_specs WHERE job_specs.ID = '%s')
		WHERE job_runs.job_spec_id = '%s'
	`, run.JobSpecID, run.JobSpecID))
}

// FindServiceAgreement looks up a ServiceAgreement by its ID.
func (orm *ORM) FindServiceAgreement(id string) (models.ServiceAgreement, error) {
	var sa models.ServiceAgreement
	return sa, orm.DB.Set("gorm:auto_preload", true).First(&sa, "id = ?", id).Error
}

// Jobs fetches all jobs.
func (orm *ORM) Jobs(cb func(models.JobSpec) bool) error {
	return Batch(1000, func(offset, limit uint) (uint, error) {
		jobs := []models.JobSpec{}
		err := orm.preloadJobs().
			Limit(limit).
			Offset(offset).
			Find(&jobs).Error
		if err != nil {
			return 0, err
		}

		for _, j := range jobs {
			if !cb(j) {
				return 0, nil
			}
		}

		return uint(len(jobs)), nil
	})
}

// JobRunsFor fetches all JobRuns with a given Job ID,
// sorted by their created at time.
func (orm *ORM) JobRunsFor(jobSpecID string) ([]models.JobRun, error) {
	runs := []models.JobRun{}
	err := orm.preloadJobRuns().
		Where("job_spec_id = ?", jobSpecID).
		Order("created_at desc").
		Find(&runs).Error
	return runs, err
}

// JobRunsCountFor returns the current number of runs for the job
func (orm *ORM) JobRunsCountFor(jobSpecID string) (int, error) {
	var count int
	err := orm.DB.
		Model(&models.JobRun{}).
		Where("job_spec_id = ?", jobSpecID).
		Count(&count).Error
	return count, err
}

// Sessions returns all sessions limited by the parameters.
func (orm *ORM) Sessions(offset, limit int) ([]models.Session, error) {
	var sessions []models.Session
	err := orm.DB.
		Set("gorm:auto_preload", true).
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

// CreateJob saves a job to the database and adds IDs to associated tables.
func (orm *ORM) CreateJob(job *models.JobSpec) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return orm.createJob(dbtx, job)
	})
}

func (orm *ORM) createJob(tx *gorm.DB, job *models.JobSpec) error {
	for i := range job.Initiators {
		job.Initiators[i].JobSpecID = job.ID
	}

	return tx.Create(job).Error
}

// Archived returns whether or not a job has been archived.
func (orm *ORM) Archived(ID string) bool {
	j, err := orm.Unscoped().FindJob(ID)
	if err != nil {
		return false
	}
	return j.DeletedAt.Valid
}

// ArchiveJob soft deletes the job and its associated job runs.
func (orm *ORM) ArchiveJob(ID string) error {
	j, err := orm.FindJob(ID)
	if err != nil {
		return err
	}

	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return multierr.Combine(
			dbtx.Where("job_spec_id = ?", ID).Delete(&models.Initiator{}).Error,
			dbtx.Where("job_spec_id = ?", ID).Delete(&models.TaskSpec{}).Error,
			dbtx.Where("job_spec_id = ?", ID).Delete(&models.JobRun{}).Error,
			dbtx.Delete(&j).Error,
		)
	})
}

// CreateServiceAgreement saves a Service Agreement, its JobSpec and its
// associations to the database.
func (orm *ORM) CreateServiceAgreement(sa *models.ServiceAgreement) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		err := orm.createJob(dbtx, &sa.JobSpec)
		if err != nil {
			return errors.Wrap(err, "Failed to create job for SA")
		}

		return dbtx.Create(sa).Error
	})
}

// UnscopedJobRunsWithStatus passes all JobRuns to a callback, one by one,
// including those that were soft deleted.
func (orm *ORM) UnscopedJobRunsWithStatus(cb func(*models.JobRun), statuses ...models.RunStatus) error {
	var runIDs []string
	err := orm.DB.Unscoped().
		Table("job_runs").
		Where("status IN (?)", statuses).
		Order("created_at asc").
		Pluck("ID", &runIDs).Error
	if err != nil {
		return fmt.Errorf("error finding job ids %v", err)
	}

	for _, id := range runIDs {
		var run models.JobRun
		err := orm.Unscoped().preloadJobRuns().First(&run, "id = ?", id).Error
		if err != nil {
			return fmt.Errorf("error finding job run %v", err)
		}

		cb(&run)
	}

	return nil
}

// AnyJobWithType returns true if there is at least one job associated with
// the type name specified and false otherwise
func (orm *ORM) AnyJobWithType(taskTypeName string) (bool, error) {
	db := orm.DB
	var taskSpec models.TaskSpec
	rval := db.Where("type = ?", taskTypeName).First(&taskSpec)
	found := !rval.RecordNotFound()
	return found, ignoreRecordNotFound(rval)
}

// CreateTxAndAttempt persists a TX and its first attempt
func (orm *ORM) CreateTx(
	tx *models.Tx,
) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return dbtx.Create(tx).Error
	})
}

// MarkTxSafe updates the database for the given transaction and attempt to
// show that the transaction has not just been confirmed,
// but has met the minimum number of outgoing confirmations to be deemed
// safely written on the blockchain.
func (orm *ORM) MarkTxSafe(tx *models.Tx, txat *models.TxAttempt) error {
	txat.Confirmed = true
	tx.AssignTxAttempt(txat)

	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return dbtx.Save(tx).Save(txat).Error
	})
}

// FindTx returns the specific transaction for the passed ID.
func (orm *ORM) FindTx(ID uint64) (*models.Tx, error) {
	tx := &models.Tx{}
	err := orm.DB.Set("gorm:auto_preload", true).First(tx, "id = ?", ID).Error
	return tx, err
}

// FindTxByAttempt returns the specific transaction attempt with the hash.
func (orm *ORM) FindTxByAttempt(hash common.Hash) (*models.Tx, error) {
	txat := &models.TxAttempt{}
	if err := orm.DB.Set("gorm:auto_preload", true).First(txat, "hash = ?", hash).Error; err != nil {
		return nil, err
	}
	tx, err := orm.FindTx(txat.TxID)
	if err != nil {
		return nil, err
	}
	tx.Hash = txat.Hash
	tx.GasPrice = txat.GasPrice
	tx.Confirmed = txat.Confirmed
	tx.Hex = txat.Hex
	tx.SentAt = txat.SentAt
	return tx, nil
}

// TxAttemptsFor returns the Transaction Attempts (TxAttempt) for a
// given Transaction ID (TxID).
func (orm *ORM) TxAttemptsFor(id uint64) ([]models.TxAttempt, error) {
	attempts := []models.TxAttempt{}
	err := orm.DB.
		Order("created_at asc").
		Where("tx_id = ?", id).
		Find(&attempts).Error
	return attempts, err
}

// AddTxAttempt creates a new transaction attempt and stores it
// in the database.
func (orm *ORM) AddTxAttempt(
	tx *models.Tx,
	etx *types.Transaction,
	blkNum uint64,
) (*models.TxAttempt, error) {
	hex, err := utils.EncodeTxToHex(etx)
	if err != nil {
		return nil, err
	}
	attempt := &models.TxAttempt{
		Hash:     etx.Hash(),
		GasPrice: models.NewBig(etx.GasPrice()),
		Hex:      hex,
		TxID:     tx.ID,
		SentAt:   blkNum,
	}
	if !tx.Confirmed {
		tx.AssignTxAttempt(attempt)
	}

	err = orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return dbtx.Save(tx).Save(attempt).Error
	})
	return attempt, err
}

// GetLastNonce retrieves the last known nonce in the database for an account
func (orm *ORM) GetLastNonce(address common.Address) (uint64, error) {
	var transaction models.Tx
	rval := orm.DB.Order("nonce desc").Where("\"from\" = ?", address).First(&transaction)
	return transaction.Nonce, ignoreRecordNotFound(rval)
}

// MarkRan will set Ran to true for a given initiator
func (orm *ORM) MarkRan(i *models.Initiator, ran bool) error {
	tx := orm.DB.Begin()
	defer tx.Rollback()

	var newi models.Initiator
	if err := tx.Select("ran").First(&newi, "ID = ?", i.ID).Error; err != nil {
		return err
	}

	if ran && newi.Ran {
		return fmt.Errorf("Initiator %v for job spec %s has already been run", i.ID, i.JobSpecID)
	}

	if err := tx.Model(i).UpdateColumn("ran", ran).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// FindUser will return the one API user, or an error.
func (orm *ORM) FindUser() (models.User, error) {
	user := models.User{}
	err := orm.DB.
		Set("gorm:auto_preload", true).
		Order("created_at desc").
		First(&user).Error
	return user, err
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
func (orm *ORM) DeleteUser() (models.User, error) {
	user, err := orm.FindUser()
	if err != nil {
		return user, err
	}

	tx := orm.DB.Begin()
	defer tx.Rollback()

	if err := tx.Delete(&user).Error; err != nil {
		return user, err
	}

	if err := tx.Delete(models.Session{}).Error; err != nil {
		return user, err
	}

	return user, tx.Commit().Error
}

// DeleteUserSession will erase the session ID for the sole API User.
func (orm *ORM) DeleteUserSession(sessionID string) error {
	return orm.DB.Where("id = ?", sessionID).Delete(models.Session{}).Error
}

// DeleteBridgeType removes the bridge type
func (orm *ORM) DeleteBridgeType(bt *models.BridgeType) error {
	return orm.DB.Delete(bt).Error
}

// DeleteJobRun deletes the job run and corresponding task runs.
func (orm *ORM) DeleteJobRun(ID string) error {
	tx := orm.DB.Begin()
	defer tx.Rollback()
	if err := tx.Where("id = ?", ID).Delete(models.JobRun{}).Error; err != nil {
		return err
	}

	if err := tx.Where("job_run_id = ?", ID).Delete(models.TaskRun{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
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

// ClearSessions removes all sessions.
func (orm *ORM) ClearSessions() error {
	return orm.DB.Delete(models.Session{}).Error
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (orm *ORM) ClearNonCurrentSessions(sessionID string) error {
	return orm.DB.Where("id <> ?", sessionID).Delete(models.Session{}).Error
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

// JobsSorted returns many JobSpecs sorted by CreatedAt from the store adhering
// to the passed parameters.
func (orm *ORM) JobsSorted(sort SortType, offset int, limit int) ([]models.JobSpec, int, error) {
	count, err := orm.countOf(&models.JobSpec{})
	if err != nil {
		return nil, 0, err
	}

	var jobs []models.JobSpec
	order := fmt.Sprintf("created_at %s", sort.String())
	err = orm.getRecords(&jobs, order, offset, limit)
	return jobs, count, err
}

// TxFrom returns all transactions from a particular address.
func (orm *ORM) TxFrom(from common.Address) ([]models.Tx, error) {
	txs := []models.Tx{}
	return txs, orm.DB.
		Set("gorm:auto_preload", true).
		Find(&txs, "\"from\" = ?", from).Error
}

// Transactions returns all transactions limited by passed parameters.
func (orm *ORM) Transactions(offset, limit int) ([]models.Tx, int, error) {
	count, err := orm.countOf(&models.Tx{})
	if err != nil {
		return nil, 0, err
	}

	var txs []models.Tx
	err = orm.getRecords(&txs, "id desc", offset, limit)
	return txs, count, err
}

// TxAttempts returns the last tx attempts sorted by sent at descending.
func (orm *ORM) TxAttempts(offset, limit int) ([]models.TxAttempt, int, error) {
	count, err := orm.countOf(&models.TxAttempt{})
	if err != nil {
		return nil, 0, err
	}

	var attempts []models.TxAttempt
	err = orm.getRecords(&attempts, "sent_at desc", offset, limit)
	return attempts, count, err
}

// JobRunsSorted returns job runs ordered and filtered by the passed params.
func (orm *ORM) JobRunsSorted(sort SortType, offset int, limit int) ([]models.JobRun, int, error) {
	count, err := orm.countOf(&models.JobRun{})
	if err != nil {
		return nil, 0, err
	}

	var runs []models.JobRun
	order := fmt.Sprintf("created_at %s", sort.String())
	err = orm.getRecords(&runs, order, offset, limit)
	return runs, count, err
}

// JobRunsSortedFor returns job runs for a specific job spec ordered and
// filtered by the passed params.
func (orm *ORM) JobRunsSortedFor(id string, order SortType, offset int, limit int) ([]models.JobRun, int, error) {
	count, err := orm.JobRunsCountFor(id)
	if err != nil {
		return nil, 0, err
	}

	var runs []models.JobRun
	err = orm.preloadJobRuns().
		Where("job_spec_id = ?", id).
		Order(fmt.Sprintf("created_at %s", order.String())).
		Limit(limit).
		Offset(offset).
		Find(&runs).Error
	return runs, count, err
}

// BridgeTypes returns bridge types ordered by name filtered limited by the
// passed params.
func (orm *ORM) BridgeTypes(offset int, limit int) ([]models.BridgeType, int, error) {
	count, err := orm.countOf(&models.BridgeType{})
	if err != nil {
		return nil, 0, err
	}

	var bridges []models.BridgeType
	err = orm.getRecords(&bridges, "name asc", offset, limit)
	return bridges, count, err
}

// SaveUser saves the user.
func (orm *ORM) SaveUser(user *models.User) error {
	return orm.DB.Save(user).Error
}

// SaveSession saves the session.
func (orm *ORM) SaveSession(session *models.Session) error {
	return orm.DB.Save(session).Error
}

// CreateBridgeType saves the bridge type.
func (orm *ORM) CreateBridgeType(bt *models.BridgeType) error {
	return orm.DB.Create(bt).Error
}

// UpdateBridgeType updates the bridge type.
func (orm *ORM) UpdateBridgeType(bt *models.BridgeType, btr *models.BridgeTypeRequest) error {
	bt.URL = btr.URL
	bt.Confirmations = btr.Confirmations
	bt.MinimumContractPayment = btr.MinimumContractPayment
	return orm.DB.Save(bt).Error
}

// SaveTx saves the transaction.
func (orm *ORM) SaveTx(tx *models.Tx) error {
	return orm.DB.Save(tx).Error
}

// CreateInitiator saves the initiator.
func (orm *ORM) CreateInitiator(initr *models.Initiator) error {
	return orm.DB.Create(initr).Error
}

// SaveHead saves the indexable block number related to head tracker.
func (orm *ORM) SaveHead(n *models.Head) error {
	return orm.DB.Save(n).Error
}

// LastHead returns the most recently persisted head entry.
func (orm *ORM) LastHead() (*models.Head, error) {
	number := &models.Head{}
	err := orm.DB.Order("number desc").First(number).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return number, err
}

// DeleteStaleSessions deletes all sessions before the passed time.
func (orm *ORM) DeleteStaleSessions(before time.Time) error {
	return orm.DB.Where("last_used < ?", before).Delete(models.Session{}).Error
}

// DeleteTransaction deletes a transaction an all of its attempts.
func (orm *ORM) DeleteTransaction(ethtx *models.Tx) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		err := dbtx.Where("id = ?", ethtx.ID).Delete(models.Tx{}).Error
		err = multierr.Append(err, dbtx.Where("tx_id = ?", ethtx.ID).Delete(models.TxAttempt{}).Error)
		return err
	})
}

// BulkDeleteRuns removes JobRuns and their related records: TaskRuns and
// RunResults.
//
// TaskRuns are removed by ON DELETE CASCADE when the JobRuns are
// deleted, but RunResults are not using foreign keys because multiple foreign
// keys on a record creates an ambiguity with gorm.
func (orm *ORM) BulkDeleteRuns(bulkQuery *models.BulkDeleteRunRequest) error {
	tx := orm.DB.Begin()
	defer tx.Rollback()
	if tx.Error != nil {
		return tx.Error
	}

	// NOTE: SQLite doesn't support compound delete statements, so delete run
	// results for job_runs ...
	err := tx.Exec(`
		DELETE
		FROM run_results
		WHERE run_results.id IN (SELECT result_id
													   FROM job_runs
														 WHERE status IN (?) AND updated_at < ?)`,
		bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
	if err != nil {
		return fmt.Errorf("error deleting JobRun's RunResults: %v", err)
	}

	// and run_requests
	err = tx.Exec(`
		DELETE
		FROM run_requests
		WHERE run_requests.id IN (SELECT run_request_id
													   FROM job_runs
														 WHERE status IN (?) AND updated_at < ?)`,
		bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
	if err != nil {
		return fmt.Errorf("error deleting JobRun's RunRequests: %v", err)
	}

	// and then task runs using a join in the subquery
	err = tx.Exec(`
		DELETE
		FROM run_results
		WHERE run_results.id IN (SELECT task_runs.result_id
													   FROM task_runs
														 INNER JOIN job_runs ON
															 task_runs.job_run_id = job_runs.id
														 WHERE job_runs.status IN (?) AND job_runs.updated_at < ?)`,
		bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
	if err != nil {
		return fmt.Errorf("error deleting TaskRuns's RunResults: %v", err)
	}

	err = tx.
		Where("status IN (?)", bulkQuery.Status.ToStrings()).
		Where("updated_at < ?", bulkQuery.UpdatedBefore).
		Unscoped().
		Delete(&[]models.JobRun{}).
		Error
	if err != nil {
		return err
	}
	return tx.Commit().Error
}

// Keys returns all keys stored in the orm.
func (orm *ORM) Keys() ([]*models.Key, error) {
	var keys []*models.Key
	return keys, orm.DB.Find(&keys).Error
}

// FirstOrCreateKey returns the first key found or creates a new one in the orm.
func (orm *ORM) FirstOrCreateKey(k *models.Key) error {
	return orm.DB.FirstOrCreate(k).Error
}

// ClobberDiskKeyStoreWithDBKeys writes all keys stored in the orm to
// the keys folder on disk, deleting anything there prior.
func (orm *ORM) ClobberDiskKeyStoreWithDBKeys(keysDir string) error {
	if err := os.RemoveAll(keysDir); err != nil {
		return err
	}

	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return err
	}

	keys, err := orm.Keys()
	if err != nil {
		return err
	}

	var merr error
	for _, k := range keys {
		merr = multierr.Append(
			k.WriteToDisk(filepath.Join(keysDir, fmt.Sprintf("%s.json", k.Address.String()))),
			merr)
	}
	return merr
}

func (orm *ORM) countOf(t interface{}) (int, error) {
	var count int
	return count, orm.DB.Model(t).Count(&count).Error
}

func (orm *ORM) getRecords(collection interface{}, order string, offset, limit int) error {
	return orm.DB.
		Set("gorm:auto_preload", true).
		Order(order).Limit(limit).Offset(offset).
		Find(collection).Error
}

// Batch is an iterator _like_ for batches of records
func Batch(chunkSize uint, cb func(offset, limit uint) (uint, error)) error {
	offset := uint(0)
	limit := uint(1000)

	for {
		count, err := cb(offset, limit)
		if err != nil {
			return err
		}

		if count < limit {
			return nil
		}

		offset += limit
	}
}
