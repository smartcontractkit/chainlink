package orm

import (
	"crypto/subtle"
	"encoding"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"chainlink/core/assets"
	"chainlink/core/auth"
	"chainlink/core/gracefulpanic"
	"chainlink/core/logger"
	"chainlink/core/store/dbutil"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // http://doc.gorm.io/database.html#connecting-to-a-database
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// BatchSize is the safe number of records to cache during Batch calls for
// SQLite without causing load problems.
const BatchSize = 100

var (
	// ErrorNotFound is returned when finding a single value fails.
	ErrorNotFound = gorm.ErrRecordNotFound
	ErrorConflict = errors.New("record already exists")
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
	db                  *gorm.DB
	lockingStrategy     LockingStrategy
	advisoryLockTimeout time.Duration
	dialectName         DialectName
	closeOnce           sync.Once
}

var (
	ErrNoAdvisoryLock    = errors.New("can't acquire advisory lock")
	ErrReleaseLockFailed = errors.New("advisory lock release failed")
)

// mapError tries to coerce the error into package defined errors.
func mapError(err error) error {
	err = errors.Cause(err)
	if v, ok := err.(sqlite3.Error); ok && v.Code == sqlite3.ErrConstraint {
		return ErrorConflict
	} else if v, ok := err.(*pq.Error); ok && v.Code.Class() == "23" {
		return ErrorConflict
	}
	return err
}

// NewORM initializes a new database file at the configured uri.
func NewORM(uri string, timeout time.Duration) (*ORM, error) {
	dialect, err := DeduceDialect(uri)
	if err != nil {
		return nil, err
	}

	lockingStrategy, err := NewLockingStrategy(dialect, uri)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create ORM lock")
	}

	logger.Infof("Locking %v for exclusive access with %v timeout", dialect, displayTimeout(timeout))

	orm := &ORM{
		lockingStrategy:     lockingStrategy,
		advisoryLockTimeout: timeout,
		dialectName:         dialect,
	}
	orm.MustEnsureAdvisoryLock()

	db, err := initializeDatabase(string(dialect), uri)
	if err != nil {
		return nil, errors.Wrap(err, "unable to init DB")
	}

	orm.db = db

	return orm, nil
}

func (orm *ORM) MustEnsureAdvisoryLock() {
	if orm.dialectName != DialectPostgres {
		return
	}
	err := orm.lockingStrategy.Lock(orm.advisoryLockTimeout)
	if err != nil {
		logger.Errorf("unable to lock ORM: %v", err)
		gracefulpanic.Panic()
	}
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
		return nil, errors.Wrapf(err, "unable to open %s for gorm DB", path)
	}

	db.SetLogger(newOrmLogWrapper(logger.GetLogger()))

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

// SetLogging turns on SQL statement logging
func (orm *ORM) SetLogging(enabled bool) {
	orm.db.LogMode(enabled)
}

// Close closes the underlying database connection.
func (orm *ORM) Close() error {
	var err error
	orm.closeOnce.Do(func() {
		err = multierr.Combine(
			orm.db.Close(),
			orm.lockingStrategy.Unlock(orm.advisoryLockTimeout),
		)
	})
	return err
}

// Unscoped returns a new instance of this ORM that includes soft deleted items.
func (orm *ORM) Unscoped() *ORM {
	return &ORM{
		db:              orm.db.Unscoped(),
		lockingStrategy: orm.lockingStrategy,
	}
}

// Where fetches multiple objects with "Find".
func (orm *ORM) Where(field string, value interface{}, instance interface{}) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Where(fmt.Sprintf("%v = ?", field), value).Find(instance).Error
}

// FindBridge looks up a Bridge by its Name.
func (orm *ORM) FindBridge(name models.TaskType) (models.BridgeType, error) {
	orm.MustEnsureAdvisoryLock()
	var bt models.BridgeType
	return bt, orm.db.First(&bt, "name = ?", name.String()).Error
}

// FindBridgesByNames finds multiple bridges by their names.
func (orm *ORM) FindBridgesByNames(names []string) ([]models.BridgeType, error) {
	orm.MustEnsureAdvisoryLock()
	var bt []models.BridgeType
	return bt, orm.db.Where("name IN (?)", names).Find(&bt).Error
}

// PendingBridgeType returns the bridge type of the current pending task,
// or error if not pending bridge.
func (orm *ORM) PendingBridgeType(jr models.JobRun) (models.BridgeType, error) {
	orm.MustEnsureAdvisoryLock()
	nextTask := jr.NextTaskRun()
	if nextTask == nil {
		return models.BridgeType{}, errors.New("Cannot find the pending bridge type of a job run with no unfinished tasks")
	}
	return orm.FindBridge(nextTask.TaskSpec.Type)
}

// FindJob looks up a Job by its ID.
func (orm *ORM) FindJob(id *models.ID) (models.JobSpec, error) {
	orm.MustEnsureAdvisoryLock()
	var job models.JobSpec
	return job, orm.preloadJobs().First(&job, "id = ?", id).Error
}

// FindInitiator returns the single initiator defined by the passed ID.
func (orm *ORM) FindInitiator(ID uint) (models.Initiator, error) {
	orm.MustEnsureAdvisoryLock()
	initr := models.Initiator{}
	return initr, orm.db.
		Set("gorm:auto_preload", true).
		First(&initr, "id = ?", ID).Error
}

func (orm *ORM) preloadJobs() *gorm.DB {
	return orm.db.
		Preload("Initiators", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Order(`"id" asc`)
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
	return orm.db.
		Preload("Initiator", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("RunRequest").
		Preload("TaskRuns", func(db *gorm.DB) *gorm.DB {
			return preloadTaskRuns(db).Order("task_spec_id asc")
		}).
		Preload("Result")
}

// FindJobRun looks up a JobRun by its ID.
func (orm *ORM) FindJobRun(id *models.ID) (models.JobRun, error) {
	orm.MustEnsureAdvisoryLock()
	var jr models.JobRun
	err := orm.preloadJobRuns().First(&jr, "id = ?", id).Error
	return jr, err
}

// AllSyncEvents returns all sync events
func (orm *ORM) AllSyncEvents(cb func(*models.SyncEvent) error) error {
	orm.MustEnsureAdvisoryLock()
	return Batch(BatchSize, func(offset, limit uint) (uint, error) {
		var events []models.SyncEvent
		err := orm.db.
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
	orm.MustEnsureAdvisoryLock()
	dbtx := orm.db.Begin()
	if dbtx.Error != nil {
		return dbtx.Error
	}
	defer dbtx.Rollback()

	err := callback(dbtx)
	if err != nil {
		return err
	}

	return dbtx.Commit().Error
}

// OptimisticUpdateConflictError is returned when a record update failed
// because another update occurred while the model was in memory and the
// differences must be reconciled.
var OptimisticUpdateConflictError = errors.New("conflict while updating record")

// SaveJobRun updates UpdatedAt for a JobRun and saves it
func (orm *ORM) SaveJobRun(run *models.JobRun) error {
	orm.MustEnsureAdvisoryLock()
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		result := dbtx.Unscoped().
			Model(run).
			Where("updated_at = ?", run.UpdatedAt).
			Omit("deleted_at").
			Save(run)
		if result.RowsAffected == 0 {
			return OptimisticUpdateConflictError
		}
		return result.Error
	})
}

// CreateJobRun inserts a new JobRun
func (orm *ORM) CreateJobRun(run *models.JobRun) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Create(run).Error
}

// LinkEarnedFor shows the total link earnings for a job
func (orm *ORM) LinkEarnedFor(spec *models.JobSpec) (*assets.Link, error) {
	orm.MustEnsureAdvisoryLock()
	var earned *assets.Link
	query := orm.db.Table("job_runs").
		Joins("JOIN job_specs ON job_runs.job_spec_id = job_specs.id").
		Where("job_specs.id = ? AND job_runs.status = ? AND job_runs.finished_at IS NOT NULL", spec.ID, models.RunStatusCompleted)

	if dbutil.IsPostgres(orm.db) {
		query = query.Select("SUM(payment)")
	} else {
		query = query.Select("CAST(SUM(CAST(SUBSTR(payment, 1, 10) as BIGINT)) as varchar(255))")
	}

	err := query.Row().Scan(&earned)
	if err != nil {
		return nil, errors.Wrap(err, "error obtaining link earned from job_runs")
	}
	return earned, nil
}

// CreateExternalInitiator inserts a new external initiator
func (orm *ORM) CreateExternalInitiator(externalInitiator *models.ExternalInitiator) error {
	orm.MustEnsureAdvisoryLock()
	err := orm.db.Create(externalInitiator).Error
	return mapError(err)
}

// DeleteExternalInitiator removes an external initiator
func (orm *ORM) DeleteExternalInitiator(name string) error {
	orm.MustEnsureAdvisoryLock()
	err := orm.db.Delete(&models.ExternalInitiator{Name: name}).Error
	return mapError(err)
}

// FindExternalInitiator finds an external initiator given an authentication request
func (orm *ORM) FindExternalInitiator(
	eia *auth.Token,
) (*models.ExternalInitiator, error) {
	orm.MustEnsureAdvisoryLock()
	initiator := &models.ExternalInitiator{}
	err := orm.db.Where("access_key = ?", eia.AccessKey).Find(initiator).Error
	if err != nil {
		return nil, errors.Wrap(err, "error finding external initiator")
	}

	return initiator, nil
}

// FindExternalInitiatorByName finds an external initiator given an authentication request
func (orm *ORM) FindExternalInitiatorByName(iname string) (models.ExternalInitiator, error) {
	orm.MustEnsureAdvisoryLock()
	var exi models.ExternalInitiator
	return exi, orm.db.First(&exi, "lower(name) = lower(?)", iname).Error
}

// FindServiceAgreement looks up a ServiceAgreement by its ID.
func (orm *ORM) FindServiceAgreement(id string) (models.ServiceAgreement, error) {
	orm.MustEnsureAdvisoryLock()
	var sa models.ServiceAgreement
	return sa, orm.db.Set("gorm:auto_preload", true).First(&sa, "id = ?", id).Error
}

// Jobs fetches all jobs.
func (orm *ORM) Jobs(cb func(*models.JobSpec) bool, initrTypes ...string) error {
	orm.MustEnsureAdvisoryLock()
	return Batch(BatchSize, func(offset, limit uint) (uint, error) {
		scope := orm.db.Limit(limit).Offset(offset)
		if len(initrTypes) > 0 {
			scope = scope.Where("initiators.type IN (?)", initrTypes)
			if dbutil.IsPostgres(orm.db) {
				scope = scope.Joins("JOIN initiators ON job_specs.id = initiators.job_spec_id::uuid")
			} else {
				scope = scope.Joins("JOIN initiators ON job_specs.id = initiators.job_spec_id")
			}
		}
		var ids []string
		err := scope.Table("job_specs").Pluck("job_specs.id", &ids).Error
		if err != nil {
			return 0, err
		}

		if len(ids) == 0 {
			return 0, nil
		}

		jobs := []models.JobSpec{}
		err = orm.preloadJobs().Unscoped().Find(&jobs, "id IN (?)", ids).Error
		if err != nil {
			return 0, err
		}
		for _, j := range jobs {
			temp := j
			if temp.DeletedAt.Valid {
				continue
			}
			if !cb(&temp) {
				return 0, nil
			}
		}

		return uint(len(jobs)), nil
	})
}

// JobRunsFor fetches all JobRuns with a given Job ID,
// sorted by their created at time.
func (orm *ORM) JobRunsFor(jobSpecID *models.ID, limit ...int) ([]models.JobRun, error) {
	orm.MustEnsureAdvisoryLock()
	runs := []models.JobRun{}
	var lim int
	if len(limit) == 0 {
		lim = 100
	} else if len(limit) >= 1 {
		lim = limit[0]
	}
	err := orm.preloadJobRuns().
		Limit(lim).
		Where("job_spec_id = ?", jobSpecID).
		Order("created_at desc").
		Find(&runs).Error
	return runs, err
}

// JobRunsCountFor returns the current number of runs for the job
func (orm *ORM) JobRunsCountFor(jobSpecID *models.ID) (int, error) {
	orm.MustEnsureAdvisoryLock()
	var count int
	err := orm.db.
		Model(&models.JobRun{}).
		Where("job_spec_id = ?", jobSpecID).
		Count(&count).Error
	return count, err
}

// Sessions returns all sessions limited by the parameters.
func (orm *ORM) Sessions(offset, limit int) ([]models.Session, error) {
	orm.MustEnsureAdvisoryLock()
	var sessions []models.Session
	err := orm.db.
		Set("gorm:auto_preload", true).
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

// GetConfigValue returns the value for a named configuration entry
func (orm *ORM) GetConfigValue(field string, value encoding.TextUnmarshaler) error {
	orm.MustEnsureAdvisoryLock()
	name := EnvVarName(field)
	config := models.Configuration{}
	if err := orm.db.First(&config, "name = ?", name).Error; err != nil {
		return err
	}
	return value.UnmarshalText([]byte(config.Value))
}

// SetConfigValue returns the value for a named configuration entry
func (orm *ORM) SetConfigValue(field string, value encoding.TextMarshaler) error {
	orm.MustEnsureAdvisoryLock()
	name := EnvVarName(field)
	textValue, err := value.MarshalText()
	if err != nil {
		return err
	}
	return orm.db.Where(models.Configuration{Name: name}).
		Assign(models.Configuration{Name: name, Value: string(textValue)}).
		FirstOrCreate(&models.Configuration{}).Error
}

// CreateJob saves a job to the database and adds IDs to associated tables.
func (orm *ORM) CreateJob(job *models.JobSpec) error {
	orm.MustEnsureAdvisoryLock()
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return orm.createJob(dbtx, job)
	})
}

func (orm *ORM) createJob(tx *gorm.DB, job *models.JobSpec) error {
	orm.MustEnsureAdvisoryLock()
	for i := range job.Initiators {
		job.Initiators[i].JobSpecID = job.ID
	}

	return tx.Create(job).Error
}

// ArchiveJob soft deletes the job and its associated job runs.
func (orm *ORM) ArchiveJob(ID *models.ID) error {
	orm.MustEnsureAdvisoryLock()
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
	orm.MustEnsureAdvisoryLock()
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
	orm.MustEnsureAdvisoryLock()
	var runIDs []string
	err := orm.db.Unscoped().
		Table("job_runs").
		Where("status IN (?)", statuses).
		Order("created_at asc").
		Pluck("ID", &runIDs).Error
	if err != nil {
		return errors.Wrap(err, "finding job ids")
	}

	return Batch(BatchSize, func(offset, limit uint) (uint, error) {
		batchIDs := runIDs[offset:utils.MinUint(limit, uint(len(runIDs)))]
		var runs []models.JobRun
		err := orm.Unscoped().
			preloadJobRuns().
			Order("job_runs.created_at asc").
			Find(&runs, "job_runs.id IN (?)", batchIDs).Error
		if err != nil {
			return 0, errors.Wrap(err, "error fetching job run batch")
		}

		for _, run := range runs {
			r := run
			cb(&r)
		}
		return uint(len(batchIDs)), nil
	})
}

// AnyJobWithType returns true if there is at least one job associated with
// the type name specified and false otherwise
func (orm *ORM) AnyJobWithType(taskTypeName string) (bool, error) {
	orm.MustEnsureAdvisoryLock()
	db := orm.db
	var taskSpec models.TaskSpec
	rval := db.Where("type = ?", taskTypeName).First(&taskSpec)
	found := !rval.RecordNotFound()
	return found, ignoreRecordNotFound(rval)
}

// CreateTx returns a transaction by its surrogate key, if it exists, or
// creates it
func (orm *ORM) CreateTx(tx *models.Tx) (*models.Tx, error) {
	orm.MustEnsureAdvisoryLock()

	err := orm.convenientTransaction(func(dbtx *gorm.DB) error {
		var query *gorm.DB
		if tx.SurrogateID.Valid {
			query = dbtx.First(&models.Tx{}, "surrogate_id = ?", tx.SurrogateID.ValueOrZero())
		} else {
			query = dbtx.First(&models.Tx{}, "hash = ?", tx.Hash)
		}

		ids := []uint64{}
		err := query.Pluck("id", &ids).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "CreateTx#First failed")
		}

		if err == gorm.ErrRecordNotFound {
			return dbtx.Create(tx).Error
		}
		tx.ID = ids[0]
		return dbtx.Save(tx).Error
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// AddTxAttempt attaches a new attempt to a Tx, after the attempt has been sent to the chain
func (orm *ORM) AddTxAttempt(tx *models.Tx, newTxAttempt *models.Tx) (*models.TxAttempt, error) {
	orm.MustEnsureAdvisoryLock()

	tx.From = newTxAttempt.From
	tx.Nonce = newTxAttempt.Nonce
	tx.GasPrice = newTxAttempt.GasPrice
	tx.Hash = newTxAttempt.Hash
	tx.SentAt = newTxAttempt.SentAt
	tx.SignedRawTx = newTxAttempt.SignedRawTx
	txAttempt := &models.TxAttempt{
		Hash:        newTxAttempt.Hash,
		GasPrice:    newTxAttempt.GasPrice,
		SentAt:      newTxAttempt.SentAt,
		SignedRawTx: newTxAttempt.SignedRawTx,
	}
	tx.Attempts = append(tx.Attempts, txAttempt)

	return txAttempt, orm.db.Save(tx).Error
}

// MarkTxSafe updates the database for the given transaction and attempt to
// show that the transaction has not just been confirmed,
// but has met the minimum number of outgoing confirmations to be deemed
// safely written on the blockchain.
func (orm *ORM) MarkTxSafe(tx *models.Tx, txAttempt *models.TxAttempt) error {
	orm.MustEnsureAdvisoryLock()
	txAttempt.Confirmed = true
	tx.Hash = txAttempt.Hash
	tx.GasPrice = txAttempt.GasPrice
	tx.Confirmed = txAttempt.Confirmed
	tx.SentAt = txAttempt.SentAt
	tx.SignedRawTx = txAttempt.SignedRawTx
	return orm.db.Save(tx).Error
}

func preloadAttempts(dbtx *gorm.DB) *gorm.DB {
	return dbtx.
		Preload("Attempts", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at asc")
		})
}

// FindTx returns the specific transaction for the passed ID.
func (orm *ORM) FindTx(ID uint64) (*models.Tx, error) {
	orm.MustEnsureAdvisoryLock()
	tx := &models.Tx{}
	err := preloadAttempts(orm.db).First(tx, "id = ?", ID).Error
	return tx, err
}

// FindTxByAttempt returns the specific transaction attempt with the hash.
func (orm *ORM) FindTxByAttempt(hash common.Hash) (*models.Tx, *models.TxAttempt, error) {
	orm.MustEnsureAdvisoryLock()
	txAttempt := &models.TxAttempt{}
	if err := orm.db.First(txAttempt, "hash = ?", hash).Error; err != nil {
		return nil, nil, err
	}
	tx, err := orm.FindTx(txAttempt.TxID)
	if err != nil {
		return nil, nil, err
	}
	return tx, txAttempt, nil
}

// FindTxAttempt returns an individual TxAttempt
func (orm *ORM) FindTxAttempt(hash common.Hash) (*models.TxAttempt, error) {
	orm.MustEnsureAdvisoryLock()
	txAttempt := &models.TxAttempt{}
	if err := orm.db.Preload("Tx").First(txAttempt, "hash = ?", hash).Error; err != nil {
		return nil, errors.Wrap(err, "FindTxByAttempt First(txAttempt) failed")
	}
	return txAttempt, nil
}

// GetLastNonce retrieves the last known nonce in the database for an account
func (orm *ORM) GetLastNonce(address common.Address) (uint64, error) {
	orm.MustEnsureAdvisoryLock()
	var transaction models.Tx
	rval := orm.db.Order("nonce desc").Where(`"from" = ?`, address).First(&transaction)
	return transaction.Nonce, ignoreRecordNotFound(rval)
}

// MarkRan will set Ran to true for a given initiator
func (orm *ORM) MarkRan(i *models.Initiator, ran bool) error {
	orm.MustEnsureAdvisoryLock()
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		var newi models.Initiator
		if err := dbtx.Select("ran").First(&newi, "ID = ?", i.ID).Error; err != nil {
			return err
		}

		if ran && newi.Ran {
			return fmt.Errorf("Initiator %v for job spec %s has already been run", i.ID, i.JobSpecID.String())
		}

		if err := dbtx.Model(i).UpdateColumn("ran", ran).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindUser will return the one API user, or an error.
func (orm *ORM) FindUser() (models.User, error) {
	orm.MustEnsureAdvisoryLock()
	user := models.User{}
	err := orm.db.
		Set("gorm:auto_preload", true).
		Order("created_at desc").
		First(&user).Error
	return user, err
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (orm *ORM) AuthorizedUserWithSession(sessionID string, sessionDuration time.Duration) (models.User, error) {
	orm.MustEnsureAdvisoryLock()
	if len(sessionID) == 0 {
		return models.User{}, errors.New("Session ID cannot be empty")
	}

	var session models.Session
	err := orm.db.First(&session, "id = ?", sessionID).Error
	if err != nil {
		return models.User{}, err
	}
	now := time.Now()
	if session.LastUsed.Add(sessionDuration).Before(now) {
		return models.User{}, errors.New("Session has expired")
	}
	session.LastUsed = now
	if err := orm.db.Save(&session).Error; err != nil {
		return models.User{}, err
	}
	return orm.FindUser()
}

// DeleteUser will delete the API User in the db.
func (orm *ORM) DeleteUser() (models.User, error) {
	orm.MustEnsureAdvisoryLock()
	user, err := orm.FindUser()
	if err != nil {
		return user, err
	}

	return user, orm.convenientTransaction(func(dbtx *gorm.DB) error {
		if err := dbtx.Delete(&user).Error; err != nil {
			return err
		}

		if err := dbtx.Delete(models.Session{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteUserSession will erase the session ID for the sole API User.
func (orm *ORM) DeleteUserSession(sessionID string) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Where("id = ?", sessionID).Delete(models.Session{}).Error
}

// DeleteBridgeType removes the bridge type
func (orm *ORM) DeleteBridgeType(bt *models.BridgeType) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Delete(bt).Error
}

// CreateSession will check the password in the SessionRequest against
// the hashed API User password in the db.
func (orm *ORM) CreateSession(sr models.SessionRequest) (string, error) {
	orm.MustEnsureAdvisoryLock()
	user, err := orm.FindUser()
	if err != nil {
		return "", err
	}

	if !constantTimeEmailCompare(sr.Email, user.Email) {
		return "", errors.New("Invalid email")
	}

	if utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		session := models.NewSession()
		return session.ID, orm.db.Save(&session).Error
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
	orm.MustEnsureAdvisoryLock()
	return orm.db.Delete(models.Session{}).Error
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (orm *ORM) ClearNonCurrentSessions(sessionID string) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Where("id <> ?", sessionID).Delete(models.Session{}).Error
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
	orm.MustEnsureAdvisoryLock()
	count, err := orm.CountOf(&models.JobSpec{})
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
	orm.MustEnsureAdvisoryLock()
	txs := []models.Tx{}
	return txs, preloadAttempts(orm.db).Find(&txs, `"from" = ?`, from).Error
}

// Transactions returns all transactions limited by passed parameters.
func (orm *ORM) Transactions(offset, limit int) ([]models.Tx, int, error) {
	orm.MustEnsureAdvisoryLock()
	count, err := orm.CountOf(&models.Tx{})
	if err != nil {
		return nil, 0, err
	}

	var txs []models.Tx
	err = orm.getRecords(&txs, "id desc", offset, limit)
	return txs, count, err
}

// TxAttempts returns the last tx attempts sorted by sent at descending.
func (orm *ORM) TxAttempts(offset, limit int) ([]models.TxAttempt, int, error) {
	orm.MustEnsureAdvisoryLock()
	count, err := orm.CountOf(&models.TxAttempt{})
	if err != nil {
		return nil, 0, err
	}

	var attempts []models.TxAttempt
	err = orm.getRecords(&attempts, "sent_at desc", offset, limit)
	return attempts, count, err
}

// UnconfirmedTxAttempts returns all TxAttempts for which the associated Tx is still unconfirmed.
func (orm *ORM) UnconfirmedTxAttempts() ([]models.TxAttempt, error) {
	orm.MustEnsureAdvisoryLock()
	var items []models.TxAttempt

	err := orm.db.
		Preload("Tx").
		Joins("inner join txes on txes.id = tx_attempts.tx_id").
		Where("txes.confirmed = ?", false).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, err
}

// JobRunsSorted returns job runs ordered and filtered by the passed params.
func (orm *ORM) JobRunsSorted(sort SortType, offset int, limit int) ([]models.JobRun, int, error) {
	orm.MustEnsureAdvisoryLock()
	count, err := orm.CountOf(&models.JobRun{})
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
func (orm *ORM) JobRunsSortedFor(id *models.ID, order SortType, offset int, limit int) ([]models.JobRun, int, error) {
	orm.MustEnsureAdvisoryLock()
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
	orm.MustEnsureAdvisoryLock()
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
	orm.MustEnsureAdvisoryLock()
	return orm.db.Save(user).Error
}

// SaveSession saves the session.
func (orm *ORM) SaveSession(session *models.Session) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Save(session).Error
}

// SaveTx saves the Ethereum Transaction.
func (orm *ORM) SaveTx(tx *models.Tx) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Save(tx).Error
}

// CreateBridgeType saves the bridge type.
func (orm *ORM) CreateBridgeType(bt *models.BridgeType) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Create(bt).Error
}

// UpdateBridgeType updates the bridge type.
func (orm *ORM) UpdateBridgeType(bt *models.BridgeType, btr *models.BridgeTypeRequest) error {
	orm.MustEnsureAdvisoryLock()
	bt.URL = btr.URL
	bt.Confirmations = btr.Confirmations
	bt.MinimumContractPayment = btr.MinimumContractPayment
	return orm.db.Save(bt).Error
}

// CreateInitiator saves the initiator.
func (orm *ORM) CreateInitiator(initr *models.Initiator) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Create(initr).Error
}

// CreateHead creates a head record that tracks which block heads we've observed in the HeadTracker
func (orm *ORM) CreateHead(n *models.Head) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Create(n).Error
}

// FirstHead returns the oldest persisted head entry.
func (orm *ORM) FirstHead() (*models.Head, error) {
	orm.MustEnsureAdvisoryLock()
	number := &models.Head{}
	err := orm.db.Order("number asc").First(number).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return number, err
}

// LastHead returns the most recently persisted head entry.
func (orm *ORM) LastHead() (*models.Head, error) {
	orm.MustEnsureAdvisoryLock()
	number := &models.Head{}
	err := orm.db.Order("number desc").First(number).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return number, err
}

// DeleteStaleSessions deletes all sessions before the passed time.
func (orm *ORM) DeleteStaleSessions(before time.Time) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Where("last_used < ?", before).Delete(models.Session{}).Error
}

// DeleteTransaction deletes a transaction an all of its attempts.
func (orm *ORM) DeleteTransaction(ethtx *models.Tx) error {
	orm.MustEnsureAdvisoryLock()
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
	orm.MustEnsureAdvisoryLock()
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		// NOTE: SQLite doesn't support compound delete statements, so delete run
		// results for job_runs ...
		err := dbtx.Exec(`
			DELETE
			FROM run_results
			WHERE run_results.id IN (SELECT result_id
															FROM job_runs
															WHERE status IN (?) AND updated_at < ?)`,
			bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
		if err != nil {
			return errors.Wrap(err, "error deleting JobRun's RunResults")
		}

		// and run_requests
		err = dbtx.Exec(`
			DELETE
			FROM run_requests
			WHERE run_requests.id IN (SELECT run_request_id
															FROM job_runs
															WHERE status IN (?) AND updated_at < ?)`,
			bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
		if err != nil {
			return errors.Wrap(err, "error deleting JobRun's RunRequests")
		}

		// and then task runs using a join in the subquery
		err = dbtx.Exec(`
			DELETE
			FROM run_results
			WHERE run_results.id IN (SELECT task_runs.result_id
															FROM task_runs
															INNER JOIN job_runs ON
																task_runs.job_run_id = job_runs.id
															WHERE job_runs.status IN (?) AND job_runs.updated_at < ?)`,
			bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
		if err != nil {
			return errors.Wrap(err, "error deleting TaskRuns's RunResults")
		}

		err = dbtx.
			Where("status IN (?)", bulkQuery.Status.ToStrings()).
			Where("updated_at < ?", bulkQuery.UpdatedBefore).
			Unscoped().
			Delete(&[]models.JobRun{}).
			Error
		if err != nil {
			return errors.Wrap(err, "error deleting JobRuns")
		}

		return nil
	})
}

// Keys returns all keys stored in the orm.
func (orm *ORM) Keys() ([]*models.Key, error) {
	orm.MustEnsureAdvisoryLock()
	var keys []*models.Key
	return keys, orm.db.Find(&keys).Error
}

// FirstOrCreateKey returns the first key found or creates a new one in the orm.
func (orm *ORM) FirstOrCreateKey(k *models.Key) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.FirstOrCreate(k).Error
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

func (orm *ORM) CountOf(t interface{}) (int, error) {
	orm.MustEnsureAdvisoryLock()
	var count int
	return count, orm.db.Model(t).Count(&count).Error
}

func (orm *ORM) getRecords(collection interface{}, order string, offset, limit int) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.
		Set("gorm:auto_preload", true).
		Order(order).Limit(limit).Offset(offset).
		Find(collection).Error
}

func (orm *ORM) RawDB(fn func(*gorm.DB) error) error {
	orm.MustEnsureAdvisoryLock()
	return fn(orm.db)
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
