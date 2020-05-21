package orm

import (
	"bytes"
	"crypto/subtle"
	"database/sql"
	"encoding"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
	"github.com/pkg/errors"
	"go.uber.org/multierr"
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
	db                  *gorm.DB
	lockingStrategy     LockingStrategy
	advisoryLockTimeout models.Duration
	closeOnce           sync.Once
	shutdownSignal      gracefulpanic.Signal
}

// NewORM initializes a new database file at the configured uri.
func NewORM(uri string, timeout models.Duration, shutdownSignal gracefulpanic.Signal, dialect DialectName, advisoryLockID int64) (*ORM, error) {
	ct, err := NewConnection(dialect, uri, advisoryLockID)
	if err != nil {
		return nil, err
	}
	// Locking strategy for transaction wrapped postgres must use original URI
	lockingStrategy, err := NewLockingStrategy(ct)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create ORM lock")
	}

	logger.Infof("Locking %v for exclusive access with %v timeout", ct.name, displayTimeout(timeout))
	orm := &ORM{
		lockingStrategy:     lockingStrategy,
		advisoryLockTimeout: timeout,
		shutdownSignal:      shutdownSignal,
	}
	orm.MustEnsureAdvisoryLock()

	db, err := ct.initializeDatabase()
	if err != nil {
		return nil, errors.Wrap(err, "unable to init DB")
	}
	orm.db = db

	return orm, nil
}

// MustEnsureAdvisoryLock sends a shutdown signal to the ORM if it an advisory
// lock cannot be acquired.
func (orm *ORM) MustEnsureAdvisoryLock() {
	err := orm.lockingStrategy.Lock(orm.advisoryLockTimeout)
	if err != nil {
		logger.Errorf("unable to lock ORM: %v", err)
		orm.shutdownSignal.Panic()
	}
}

func displayTimeout(timeout models.Duration) string {
	if timeout.IsInstant() {
		return "indefinite"
	}
	return timeout.String()
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
	if err := orm.db.Where("name IN (?)", names).Find(&bt).Error; err != nil {
		return nil, err
	}
	if len(bt) != len(names) {
		return nil, errors.New("bridge names don't exist or duplicates present")
	}
	return bt, nil
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
func (orm *ORM) FindInitiator(ID int64) (models.Initiator, error) {
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

// NOTE: Copied verbatim from gorm master
// Transaction start a transaction as a block,
// return error will rollback, otherwise to commit.
func (orm *ORM) Transaction(fc func(tx *gorm.DB) error) (err error) {
	tx := orm.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
			tx.Rollback()
			return
		}
	}()

	err = fc(tx)

	if err == nil {
		err = tx.Commit().Error
	}

	// Makesure rollback when Block error or Commit error
	if err != nil {
		tx.Rollback()
	}
	return
}

// convenientTransaction handles setup and teardown for a gorm database
// transaction, handing off the database transaction to the callback parameter.
// Encourages the use of transactions for gorm calls that translate
// into multiple sql calls, i.e. orm.SaveJobRun(run), which are better suited
// in a database transaction.
func (orm *ORM) convenientTransaction(callback func(*gorm.DB) error) error {
	orm.MustEnsureAdvisoryLock()
	return orm.Transaction(callback)
}

// SaveJobRun updates UpdatedAt for a JobRun and saves it
func (orm *ORM) SaveJobRun(run *models.JobRun) error {
	orm.MustEnsureAdvisoryLock()
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		result := dbtx.Unscoped().
			Model(run).
			Where("updated_at = ?", run.UpdatedAt).
			Omit("deleted_at").
			Save(run)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrOptimisticUpdateConflict
		}
		return nil
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
	return err
}

// DeleteExternalInitiator removes an external initiator
func (orm *ORM) DeleteExternalInitiator(name string) error {
	orm.MustEnsureAdvisoryLock()
	err := orm.db.Delete(&models.ExternalInitiator{Name: name}).Error
	return err
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

// ArchiveJob soft deletes the job, job_runs and its initiator.
func (orm *ORM) ArchiveJob(ID *models.ID) error {
	orm.MustEnsureAdvisoryLock()
	j, err := orm.FindJob(ID)
	if err != nil {
		return err
	}

	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return multierr.Combine(
			dbtx.Exec("UPDATE initiators SET deleted_at = NOW() WHERE job_spec_id = ?", ID).Error,
			dbtx.Exec("UPDATE task_specs SET deleted_at = NOW() WHERE job_spec_id = ?", ID).Error,
			dbtx.Exec("UPDATE job_runs SET deleted_at = NOW() WHERE job_spec_id = ?", ID).Error,
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

// IdempotentInsertEthTaskRunTransaction creates both eth_task_run_transaction and eth_transaction in one hit
// It can be called multiple times without error as long as the outcome would have resulted in the same database state
func (orm *ORM) IdempotentInsertEthTaskRunTransaction(taskRunID models.ID, fromAddress common.Address, toAddress common.Address, encodedPayload []byte, gasLimit uint64) error {
	ethTransaction := models.EthTransaction{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          assets.NewEthValue(0),
		GasLimit:       gasLimit,
	}
	ethTaskRunTransaction := models.EthTaskRunTransaction{
		TaskRunID: taskRunID.UUID(),
	}
	err := orm.Transaction(func(dbtx *gorm.DB) error {
		if err := dbtx.Save(&ethTransaction).Error; err != nil {
			return err
		}
		ethTaskRunTransaction.EthTransactionID = ethTransaction.ID
		if err := dbtx.Create(&ethTaskRunTransaction).Error; err != nil {
			return err
		}
		return nil
	})
	switch v := err.(type) {
	case *pq.Error:
		if v.Constraint == "idx_eth_task_run_transactions_task_run_id" {
			savedRecord, err := orm.FindEthTaskRunTransactionByTaskRunID(taskRunID)
			if err != nil {
				return err
			}
			t := savedRecord.EthTransaction
			if t.ToAddress != toAddress || !bytes.Equal(t.EncodedPayload, encodedPayload) {
				return fmt.Errorf(
					"transaction already exists for task run ID %s but it has different parameters\n"+
						"New parameters: toAddress: %s, encodedPayload: 0x%s"+
						"Existing record has: toAddress: %s, encodedPayload: 0x%s",
					taskRunID.String(),
					toAddress.String(), hex.EncodeToString(encodedPayload),
					t.ToAddress.String(), hex.EncodeToString(t.EncodedPayload),
				)
			}
			return nil
		}
		return err
	default:
		return err
	}
}

// FindEthTaskRunTransactionByTaskRunID finds the EthTaskRunTransaction with its EthTransaction preloaded
func (orm *ORM) FindEthTaskRunTransactionByTaskRunID(taskRunID models.ID) (models.EthTaskRunTransaction, error) {
	etrt := models.EthTaskRunTransaction{}
	err := orm.db.Preload("EthTransaction").First(&etrt, "task_run_id = ?", &taskRunID).Error
	return etrt, err
}

// FindEthTransactionWithAttempts finds the EthTransaction with its attempts preloaded
func (orm *ORM) FindEthTransactionWithAttempts(ethTransactionID int64) (models.EthTransaction, error) {
	etx := models.EthTransaction{}
	err := orm.db.Preload("EthTransactionAttempts").First(&etx, "id = ?", &ethTransactionID).Error
	return etx, err
}

// CreateTx finds and overwrites a transaction by its surrogate key, if it exists, or
// creates it
func (orm *ORM) CreateTx(tx *models.Tx) (*models.Tx, error) {
	orm.MustEnsureAdvisoryLock()

	err := orm.convenientTransaction(func(dbtx *gorm.DB) error {
		var query *gorm.DB
		foundTx := models.Tx{}
		if tx.SurrogateID.Valid {
			query = dbtx.First(&foundTx, "surrogate_id = ?", tx.SurrogateID.ValueOrZero())
		} else {
			query = dbtx.First(&foundTx, "hash = ?", tx.Hash)
		}
		err := query.Error

		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return errors.Wrap(err, "CreateTx#First failed")
		}

		if gorm.IsRecordNotFoundError(err) {
			return dbtx.Create(tx).Error
		}

		tx.ID = foundTx.ID
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

// FindAllTxsInNonceRange returns an array of transactions matching the inclusive range between beginningNonce and endingNonce
func (orm *ORM) FindAllTxsInNonceRange(beginningNonce uint, endingNonce uint) ([]models.Tx, error) {
	orm.MustEnsureAdvisoryLock()
	var txs []models.Tx
	err := orm.db.Order("nonce ASC, sent_at ASC").Where(`nonce BETWEEN ? AND ?`, beginningNonce, endingNonce).Find(&txs).Error
	return txs, err
}

// FindTxsBySenderAndRecipient returns an array of transactions sent by `sender` to `recipient`
func (orm *ORM) FindTxsBySenderAndRecipient(sender, recipient common.Address, offset, limit uint) ([]models.Tx, error) {
	orm.MustEnsureAdvisoryLock()
	var txs []models.Tx
	err := orm.db.
		Where(`"from" = ? AND "to" = ?`, sender, recipient).
		Order("nonce DESC").
		Offset(offset).
		Limit(limit).
		Find(&txs).Error
	return txs, err
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
			return fmt.Errorf("initiator %v for job spec %s has already been run", i.ID, i.JobSpecID.String())
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
	if initr.JobSpecID == nil {
		// NOTE: This hangs forever if we don't check this here and the
		// supplied initiator does not have a JobSpecID set.
		// I do not know why. Seems to be something going wrong inside gorm
		logger.Error("cannot create initiator without job spec ID")
		return errors.New("requires job spec ID")
	}
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
// RunResults and RunRequests are pointed at by JobRuns so we must use two CTEs
// to remove both parents in one hit.
//
// TaskRuns are removed by ON DELETE CASCADE when the JobRuns and RunResults
// are deleted.
func (orm *ORM) BulkDeleteRuns(bulkQuery *models.BulkDeleteRunRequest) error {
	orm.MustEnsureAdvisoryLock()
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		err := dbtx.Exec(`
			WITH deleted_job_runs AS (
				DELETE FROM job_runs WHERE status IN (?) AND updated_at < ? RETURNING result_id, run_request_id
			),
			deleted_run_results AS (
				DELETE FROM run_results WHERE id IN (SELECT result_id FROM deleted_job_runs)
			)
			DELETE FROM run_requests WHERE id IN (SELECT run_request_id FROM deleted_job_runs)`,
			bulkQuery.Status.ToStrings(), bulkQuery.UpdatedBefore).Error
		if err != nil {
			return errors.Wrap(err, "error deleting JobRuns")
		}

		return nil
	})
}

// Keys returns all of the keys recorded in the database.
func (orm *ORM) Keys() ([]*models.Key, error) {
	orm.MustEnsureAdvisoryLock()
	var keys []*models.Key
	return keys, orm.db.Find(&keys).Order("created_at ASC, address ASC").Error
}

// DeleteKey deletes a key whose address matches the supplied bytes.
func (orm *ORM) DeleteKey(address []byte) error {
	return orm.db.Exec("DELETE FROM keys WHERE address = ?", address).Error
}

// FirstOrCreateKey returns the first key found or creates a new one in the orm.
func (orm *ORM) FirstOrCreateKey(k *models.Key) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.FirstOrCreate(k).Error
}

// FirstOrCreateEncryptedSecretKey returns the first key found or creates a new one in the orm.
func (orm *ORM) FirstOrCreateEncryptedSecretVRFKey(k *models.EncryptedSecretVRFKey) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.FirstOrCreate(k).Error
}

// DeleteEncryptedSecretKey deletes k from the encrypted keys table, or errors
func (orm *ORM) DeleteEncryptedSecretVRFKey(k *models.EncryptedSecretVRFKey) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Delete(k).Error
}

// FindEncryptedSecretKeys retrieves matches to where from the encrypted keys table, or errors
func (orm *ORM) FindEncryptedSecretVRFKeys(where ...models.EncryptedSecretVRFKey) (
	retrieved []*models.EncryptedSecretVRFKey, err error) {
	orm.MustEnsureAdvisoryLock()
	var anonWhere []interface{} // Find needs "where" contents coerced to interface{}
	for _, constraint := range where {
		anonWhere = append(anonWhere, &constraint)
	}
	return retrieved, orm.db.Find(&retrieved, anonWhere...).Error
}

// SaveLogCursor saves the log cursor.
func (orm *ORM) SaveLogCursor(logCursor *models.LogCursor) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Save(logCursor).Error
}

// FindLogCursor will find the given log cursor.
func (orm *ORM) FindLogCursor(name string) (models.LogCursor, error) {
	orm.MustEnsureAdvisoryLock()
	lc := models.LogCursor{}
	err := orm.db.
		Where("name = ?", name).
		First(&lc).Error
	return lc, err
}

// HasConsumedLog reports whether the given consumer had already consumed the given log
func (orm *ORM) HasConsumedLog(rawLog eth.RawLog, JobID *models.ID) (bool, error) {
	lc := models.LogConsumption{
		BlockHash: rawLog.GetBlockHash(),
		LogIndex:  rawLog.GetIndex(),
		JobID:     JobID,
	}
	return orm.LogConsumptionExists(&lc)
}

// LogConsumptionExists reports whether a given LogConsumption record already exists
func (orm *ORM) LogConsumptionExists(lc *models.LogConsumption) (bool, error) {
	subQuery := "SELECT id FROM log_consumptions " +
		"WHERE block_hash=$1 " +
		"AND log_index=$2 " +
		"AND job_id=$3"
	query := "SELECT exists (" + subQuery + ")"

	var exists bool
	err := orm.db.DB().
		QueryRow(query, lc.BlockHash, lc.LogIndex, lc.JobID).
		Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

// CreateLogConsumption creates a new LogConsumption record
func (orm *ORM) CreateLogConsumption(lc *models.LogConsumption) error {
	orm.MustEnsureAdvisoryLock()
	return orm.db.Create(lc).Error
}

// FindLogConsumer finds the consuming job of a particular LogConsumption record
func (orm *ORM) FindLogConsumer(lc *models.LogConsumption) (models.JobSpec, error) {
	orm.MustEnsureAdvisoryLock()
	return orm.FindJob(lc.JobID)
}

// ClobberDiskKeyStoreWithDBKeys writes all keys stored in the orm to
// the keys folder on disk, deleting anything there prior.
func (orm *ORM) ClobberDiskKeyStoreWithDBKeys(keysDir string) error {
	if err := os.RemoveAll(keysDir); err != nil {
		return err
	}

	if err := utils.EnsureDirAndPerms(keysDir, 0700); err != nil {
		return err
	}

	keys, err := orm.Keys()
	if err != nil {
		return err
	}

	var merr error
	for _, k := range keys {
		merr = multierr.Append(
			k.WriteToDisk(filepath.Join(keysDir, keyFileName(k.Address, k.CreatedAt))),
			merr)
	}
	return merr
}

// Copied directly from geth - see: https://github.com/ethereum/go-ethereum/blob/32d35c9c088463efac49aeb0f3e6d48cfb373a40/accounts/keystore/key.go#L217
func keyFileName(keyAddr models.EIP55Address, createdAt time.Time) string {
	return fmt.Sprintf("UTC--%s--%s", toISO8601(createdAt), keyAddr[2:])
}

// Copied directly from geth - see: https://github.com/ethereum/go-ethereum/blob/32d35c9c088463efac49aeb0f3e6d48cfb373a40/accounts/keystore/key.go#L217
func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
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

// DialectName is a compiler enforced type used that maps to gorm's dialect
// names.
type DialectName string

const (
	// DialectPostgres represents the postgres dialect.
	DialectPostgres DialectName = "postgres"
	// DialectTransactionWrappedPostgres is useful for tests.
	// When the connection is opened, it starts a transaction and all
	// operations performed on the DB will be within that transaction.
	//
	// HACK: This must be the string 'cloudsqlpostgres' because of an absolutely
	// horrible design in gorm. We need gorm to enable postgres-specific
	// features for the txdb driver, but it can only do that if the dialect is
	// called "postgres" or "cloudsqlpostgres".
	//
	// Since "postgres" is already taken, "cloudsqlpostgres" is our only
	// remaining option
	//
	// See: https://github.com/jinzhu/gorm/blob/master/dialect_postgres.go#L15
	DialectTransactionWrappedPostgres DialectName = "cloudsqlpostgres"
	// DialectPostgresWithoutLock represents the postgres dialect but it does not
	// wait for a lock to connect. Intended to be used for read only access.
	DialectPostgresWithoutLock DialectName = "postgresWithoutLock"
)

// Connection manages all of the possible database connection setup and config.
type Connection struct {
	name               DialectName
	uri                string
	dialect            DialectName
	locking            bool
	advisoryLockID     int64
	transactionWrapped bool
}

// NewConnection returns a Connection which holds all of the configuration
// necessary for managing the database connection.
func NewConnection(dialect DialectName, uri string, advisoryLockID int64) (Connection, error) {
	switch dialect {
	case DialectPostgres:
		return Connection{
			advisoryLockID:     advisoryLockID,
			dialect:            DialectPostgres,
			locking:            true,
			name:               dialect,
			transactionWrapped: false,
			uri:                uri,
		}, nil
	case DialectPostgresWithoutLock:
		return Connection{
			advisoryLockID:     advisoryLockID,
			dialect:            DialectPostgres,
			locking:            false,
			name:               dialect,
			transactionWrapped: false,
			uri:                uri,
		}, nil
	case DialectTransactionWrappedPostgres:
		return Connection{
			advisoryLockID:     advisoryLockID,
			dialect:            DialectTransactionWrappedPostgres,
			locking:            true,
			name:               dialect,
			transactionWrapped: true,
			uri:                uri,
		}, nil
	}
	return Connection{}, errors.Errorf("%s is not a valid dialect type", dialect)
}

func (ct Connection) initializeDatabase() (*gorm.DB, error) {
	if ct.transactionWrapped {
		// Dbtx uses the uri as a unique identifier for each transaction. Each ORM
		// should be encapsulated in it's own transaction, and thus needs its own
		// unique id.
		//
		// We can happily throw away the original uri here because if we are using
		// txdb it should have already been set at the point where we called
		// txdb.Register
		ct.uri = models.NewID().String()
	}

	db, err := gorm.Open(string(ct.dialect), ct.uri)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %s for gorm DB", ct.uri)
	}

	db.SetLogger(newOrmLogWrapper(logger.GetLogger()))

	if err := dbutil.SetTimezone(db); err != nil {
		return nil, err
	}

	if ct.transactionWrapped {
		// Required to prevent phantom reads in overlapping tests
		err := db.Exec(`SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE`).Error
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// BatchSize is the safe number of records to cache during Batch calls for
// SQLite without causing load problems.
// NOTE: Now we no longer support SQLite, perhaps this can be tuned?
const BatchSize = 100

// NOTE: Skips advisory lock, use caution!
func (orm *ORM) GetRawDB() *gorm.DB {
	return orm.db
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
