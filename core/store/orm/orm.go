package orm

import (
	"bytes"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"gorm.io/gorm/clause"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"

	"gorm.io/datatypes"
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

// UpsertNodeVersion inserts a new NodeVersion
func (orm *ORM) UpsertNodeVersion(version models.NodeVersion) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}

	return orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&version).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// FindLatestNodeVersion looks up the latest node version
func (orm *ORM) FindLatestNodeVersion() (*models.NodeVersion, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, err
	}
	var nodeVersion models.NodeVersion
	err := orm.DB.Order("created_at DESC").First(&nodeVersion).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil && strings.Contains(err.Error(), "relation \"node_versions\" does not exist") {
		logger.Default.Debug("Failed to find any node version in the DB, the node_versions table does not exist yet.")
		return nil, nil
	}
	return &nodeVersion, err
}

// FindBridge looks up a Bridge by its Name.
func (orm *ORM) FindBridge(name models.TaskType) (bt models.BridgeType, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return bt, err
	}
	return bt, orm.DB.First(&bt, "name = ?", name.String()).Error
}

// FindBridgesByNames finds multiple bridges by their names.
func (orm *ORM) FindBridgesByNames(names []string) ([]models.BridgeType, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, err
	}
	var bt []models.BridgeType
	if err := orm.DB.Where("name IN (?)", names).Find(&bt).Error; err != nil {
		return nil, err
	}
	if len(bt) != len(names) {
		return nil, errors.New("bridge names don't exist or duplicates present")
	}
	return bt, nil
}

// PendingBridgeType returns the bridge type of the current pending task,
// or error if not pending bridge.
func (orm *ORM) PendingBridgeType(jr models.JobRun) (bt models.BridgeType, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return bt, err
	}
	nextTask := jr.NextTaskRun()
	if nextTask == nil {
		return models.BridgeType{}, errors.New("Cannot find the pending bridge type of a job run with no unfinished tasks")
	}
	return orm.FindBridge(nextTask.TaskSpec.Type)
}

// FindJob looks up a JobSpec by its ID.
func (orm *ORM) FindJobSpec(id models.JobID) (job models.JobSpec, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return job, err
	}
	return job, orm.preloadJobs().First(&job, "id = ?", id).Error
}

func (orm *ORM) FindJobSpecUnscoped(id models.JobID) (job models.JobSpec, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return job, err
	}
	return job, orm.preloadJobs().First(&job, "id = ?", id).Error
}

// FindJobWithErrors looks up a Job by its ID and preloads JobSpecErrors.
func (orm *ORM) FindJobWithErrors(id models.JobID) (models.JobSpec, error) {
	var job models.JobSpec
	err := orm.
		preloadJobs().
		Preload("Errors", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Order("id asc")
		}).
		First(&job, "id = ?", id).Error
	return job, err
}

// FindInitiator returns the single initiator defined by the passed ID.
func (orm *ORM) FindInitiator(ID int64) (initr models.Initiator, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return initr, err
	}
	return initr, orm.DB.
		Preload(clause.Associations).
		First(&initr, "id = ?", ID).Error
}

func (orm *ORM) preloadJobs() *gorm.DB {
	return orm.DB.
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
	return orm.DB.
		Preload("Initiator", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("RunRequest").
		Preload("TaskRuns", func(db *gorm.DB) *gorm.DB {
			return preloadTaskRuns(db).Order("task_spec_id asc")
		}).
		Preload("Result")
}

func (orm *ORM) preloadJobRunsUnscoped() *gorm.DB {
	return preloadJobRunsUnscoped(orm.DB)
}

func preloadJobRunsUnscoped(db *gorm.DB) *gorm.DB {
	return db.Unscoped().
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
func (orm *ORM) FindJobRun(id uuid.UUID) (jr models.JobRun, err error) {
	if err = orm.MustEnsureAdvisoryLock(); err != nil {
		return jr, err
	}
	err = orm.preloadJobRuns().First(&jr, "id = ?", id).Error
	return jr, err
}

func (orm *ORM) FindJobRunIncludingArchived(id uuid.UUID) (jr models.JobRun, err error) {
	if err = orm.MustEnsureAdvisoryLock(); err != nil {
		return jr, err
	}
	err = orm.preloadJobRunsUnscoped().First(&jr, "id = ?", id).Error
	return jr, err
}

func (orm *ORM) Transaction(fc func(tx *gorm.DB) error) (err error) {
	return orm.convenientTransaction(fc)
}

// convenientTransaction handles setup and teardown for a gorm database
// transaction, handing off the database transaction to the callback parameter.
// Encourages the use of transactions for gorm calls that translate
// into multiple sql calls, i.e. orm.SaveJobRun(run), which are better suited
// in a database transaction.
func (orm *ORM) convenientTransaction(callback func(*gorm.DB) error) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return postgres.GormTransactionWithDefaultContext(orm.DB, callback)
}

// SaveJobRun updates UpdatedAt for a JobRun and updates its status, finished
// at and run results.
// It auto-inserts run results if none are existing already for the given
// jobrun/taskruns.
func (orm *ORM) SaveJobRun(run *models.JobRun) error {
	if run.ID == uuid.Nil {
		return errors.New("job run did not have an ID set")
	}
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	err := postgres.GormTransactionWithDefaultContext(orm.DB, func(dbtx *gorm.DB) error {
		result := dbtx.Exec(`
UPDATE job_runs SET "status"=?, "finished_at"=?, "updated_at"=NOW(), "creation_height"=?, "observed_height"=?, "payment"=?
WHERE updated_at = ? AND "id" = ?`,
			run.Status, run.FinishedAt, run.CreationHeight, run.ObservedHeight, run.Payment,
			run.UpdatedAt, run.ID,
		)
		if result.Error != nil {
			return errors.Wrap(result.Error, "failed to update job run")
		}
		if result.RowsAffected == 0 {
			return ErrOptimisticUpdateConflict
		}

		// NOTE: uRunResultStrs/uRunResultArgs handles updating the run_results
		// for both the job_run and task_runs, if any
		uRunResultStrs := []string{}
		uRunResultArgs := []interface{}{}
		if run.Result.ID > 0 {
			uRunResultStrs = append(uRunResultStrs, "(?::bigint, ?::jsonb, ?::text)")
			uRunResultArgs = append(uRunResultArgs, run.ResultID, run.Result.Data, run.Result.ErrorMessage)
		} else {
			if run.ResultID.Valid {
				return errors.Errorf("got valid ResultID %v for job run %s but Result.ID was 0; expected JobRun.Result to be preloaded", run.ResultID.Int64, run.ID)
			}
			// Insert the run result for the job run
			err := dbtx.Exec(`
WITH run_result AS (
	INSERT INTO run_results (data, error_message, created_at, updated_at) VALUES(?, ?, NOW(), NOW()) RETURNING id
)
UPDATE job_runs
SET result_id = run_result.id
FROM run_result
WHERE job_runs.id = ?
`, run.Result.Data, run.Result.ErrorMessage, run.ID).Error
			if err != nil {
				return errors.Wrap(err, "error inserting run result for job_run")
			}
		}

		uTaskRunStrs := []string{}
		uTaskRunArgs := []interface{}{}
		for _, tr := range run.TaskRuns {
			uTaskRunStrs = append(uTaskRunStrs, "(?::uuid, ?::run_status)")
			uTaskRunArgs = append(uTaskRunArgs, tr.ID, tr.Status)
			if tr.Result.ID > 0 {
				uRunResultStrs = append(uRunResultStrs, "(?::bigint, ?::jsonb, ?::text)")
				uRunResultArgs = append(uRunResultArgs, tr.ResultID, tr.Result.Data, tr.Result.ErrorMessage)
			} else {
				if tr.ResultID.Valid {
					return errors.Errorf("got valid ResultID %v for task run %s but Result.ID was 0; expected TaskRun.Result to be preloaded", tr.ResultID.Int64, tr.ID)
				}
				// Insert the run result for the task run
				res := dbtx.Exec(`
WITH run_result AS (
	INSERT INTO run_results (data, error_message, created_at, updated_at) VALUES(?, ?, NOW(), NOW()) RETURNING id
)
UPDATE task_runs
SET result_id = run_result.id
FROM run_result
WHERE task_runs.id = ?
`, tr.Result.Data, tr.Result.ErrorMessage, tr.ID)
				if res.Error != nil {
					return errors.Wrap(res.Error, "error inserting run result for job_run")
				}
				if res.RowsAffected == 0 {
					return errors.Errorf("failed to insert run_result; task run with id %v was missing", tr.ID)
				}
			}
		}

		if len(uTaskRunStrs) > 0 {
			updateTaskRunsSQL := `
UPDATE task_runs SET status=updates.status, updated_at=NOW()
FROM (VALUES
%s
) AS updates(id, status)
WHERE task_runs.id = updates.id
`
			/* #nosec G201 */
			stmt := fmt.Sprintf(updateTaskRunsSQL, strings.Join(uTaskRunStrs, ","))
			err := dbtx.Exec(stmt, uTaskRunArgs...).Error
			if err != nil {
				return errors.Wrap(err, "failed to update task runs")
			}
		}

		if len(uRunResultStrs) > 0 {
			updateRunResultsSQL := `
UPDATE run_results SET data=updates.data, error_message=updates.error_message, updated_at=NOW()
FROM (VALUES
%s
) AS updates(id, data, error_message)
WHERE run_results.id = updates.id
`
			/* #nosec G201 */
			stmt := fmt.Sprintf(updateRunResultsSQL, strings.Join(uRunResultStrs, ","))
			err := dbtx.Exec(stmt, uRunResultArgs...).Error
			if err != nil {
				return errors.Wrap(err, "failed to update task run results")
			}
		}

		// Preload the lot again, it's the easiest way to make sure our
		// returned model is in sync with the database
		err := preloadJobRunsUnscoped(dbtx).First(run, "id = ?", run.ID).Error
		if err != nil {
			return errors.Wrap(err, "failed to reload job run after update")
		}

		// Insert sync_event
		return errors.Wrap(synchronization.InsertSyncEventForJobRun(dbtx, run), "failed to insert sync_event for updated run")
	})

	return errors.Wrap(err, "SaveJobRun failed")
}

// CreateJobRun inserts a new JobRun
func (orm *ORM) CreateJobRun(run *models.JobRun) error {
	return orm.convenientTransaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&run.RunRequest).Error
		if err != nil {
			return err
		}

		err = tx.Create(&run.Result).Error
		if err != nil {
			return err
		}

		run.RunRequestID = clnull.Int64From(run.RunRequest.ID)
		run.ResultID = clnull.Int64From(run.Result.ID)
		err = tx.Create(run).Error
		if err != nil {
			return err
		}

		for i := range run.TaskRuns {
			err = tx.Create(&run.TaskRuns[i].Result).Error
			if err != nil {
				return err
			}

			run.TaskRuns[i].JobRunID = run.ID
			run.TaskRuns[i].ResultID = clnull.Int64From(run.TaskRuns[i].Result.ID)
			err = tx.Create(&run.TaskRuns[i]).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// LinkEarnedFor shows the total link earnings for a job
func (orm *ORM) LinkEarnedFor(spec *models.JobSpec) (*assets.Link, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, err
	}
	var earned *assets.Link
	query := orm.DB.Table("job_runs").
		Joins("JOIN job_specs ON job_runs.job_spec_id = job_specs.id").
		Where("job_specs.id = ? AND job_runs.status = ? AND job_runs.finished_at IS NOT NULL", spec.ID, models.RunStatusCompleted)

	query = query.Select("SUM(payment)")

	err := query.Row().Scan(&earned)
	if err != nil {
		return nil, errors.Wrap(err, "error obtaining link earned from job_runs")
	}
	return earned, nil
}

// UpsertErrorFor upserts a JobSpecError record, incrementing the occurrences counter by 1
// if the record is found
func (orm *ORM) UpsertErrorFor(jobID models.JobID, description string) {
	jse := models.NewJobSpecError(jobID, description)
	err := orm.DB.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "job_spec_id"}, {Name: "description"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"occurrences": gorm.Expr("job_spec_errors.occurrences + 1"),
				"updated_at":  gorm.Expr("excluded.updated_at"),
			}),
		}).
		Create(&jse).
		Error

	logger.ErrorIf(err, fmt.Sprintf("Unable to create JobSpecError: %v", err))
}

// FindJobSpecError looks for a JobSpecError record with the given jobID and description
func (orm *ORM) FindJobSpecError(jobID models.JobID, description string) (*models.JobSpecError, error) {
	jobSpecErr := &models.JobSpecError{}
	err := orm.DB.
		Where("job_spec_id = ? AND description = ?", jobID, description).
		First(&jobSpecErr).Error
	return jobSpecErr, err
}

// DeleteJobSpecError removes a JobSpecError
func (orm *ORM) DeleteJobSpecError(ID int64) error {
	result := orm.DB.Exec("DELETE FROM job_spec_errors WHERE id = ?", ID)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

// FindServiceAgreement looks up a ServiceAgreement by its ID.
func (orm *ORM) FindServiceAgreement(id string) (sa models.ServiceAgreement, err error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return sa, err
	}
	return sa, orm.DB.Preload(clause.Associations).First(&sa, "id = ?", id).Error
}

// Jobs fetches all jobs.
func (orm *ORM) Jobs(cb func(*models.JobSpec) bool, initrTypes ...string) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return postgres.Batch(func(offset, limit uint) (uint, error) {
		scope := orm.DB.Order("job_specs.id asc").Limit(int(limit)).Offset(int(offset))
		if len(initrTypes) > 0 {
			scope = scope.Where("initiators.type IN (?)", initrTypes)
			scope = scope.Joins("JOIN initiators ON job_specs.id = initiators.job_spec_id::uuid")
		}
		var ids []string
		err := scope.Table("job_specs").Distinct("job_specs.id").Pluck("job_specs.id", &ids).Error
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
func (orm *ORM) JobRunsFor(jobSpecID models.JobID, limit ...int) ([]models.JobRun, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, err
	}
	runs := []models.JobRun{}
	var lim int
	if len(limit) == 0 {
		lim = 100
	} else if len(limit) >= 1 {
		lim = limit[0]
	}
	err := orm.preloadJobRuns().
		Limit(lim).
		Where("job_spec_id = ?", jobSpecID.UUID()).
		Order("created_at desc").
		Find(&runs).Error
	return runs, err
}

// JobRunsCountFor returns the current number of runs for the job
func (orm *ORM) JobRunsCountFor(jobSpecID models.JobID) (int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return 0, err
	}
	var count int64
	err := orm.DB.
		Model(&models.JobRun{}).
		Where("job_spec_id = ?", jobSpecID).
		Count(&count).Error
	return int(count), err
}

func (orm *ORM) JobRunsCountForGivenStatus(jobSpecID models.JobID, status string) (int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return 0, err
	}
	var count int64
	err := orm.DB.
		Model(&models.JobRun{}).
		Where("job_spec_id = ?", jobSpecID).
		Where("status = ?", status).
		Count(&count).Error
	return int(count), err
}

// CreateJob saves a job to the database and adds IDs to associated tables.
func (orm *ORM) CreateJob(job *models.JobSpec) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return orm.createJob(orm.DB, job)
	})
}

func (orm *ORM) createJob(tx *gorm.DB, job *models.JobSpec) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}

	err := tx.Create(job).Error
	if err != nil {
		return nil
	}

	for i := range job.Tasks {
		job.Tasks[i].JobSpecID = job.ID
		err := tx.Create(&job.Tasks[i]).Error
		if err != nil {
			return err
		}
	}

	for i := range job.Initiators {
		job.Initiators[i].JobSpecID = job.ID
		err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&job.Initiators[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// ArchiveJob soft deletes the job, job_runs and its initiator.
// It is idempotent, subsequent runs will do nothing and return no error
func (orm *ORM) ArchiveJob(ID models.JobID) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		return multierr.Combine(
			dbtx.Exec("UPDATE initiators SET deleted_at = NOW() WHERE job_spec_id = ? AND deleted_at IS NULL", ID).Error,
			dbtx.Exec("UPDATE task_specs SET deleted_at = NOW() WHERE job_spec_id = ? AND deleted_at IS NULL", ID).Error,
			dbtx.Exec("UPDATE job_runs SET deleted_at = NOW() WHERE job_spec_id = ? AND deleted_at IS NULL", ID).Error,
			dbtx.Exec("UPDATE job_specs SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", ID).Error,
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

		err = dbtx.Create(&sa.Encumbrance).Error
		if err != nil {
			return errors.Wrap(err, "Failed to create Encumberance for SA")
		}

		sa.EncumbranceID = sa.Encumbrance.ID
		return dbtx.Create(sa).Error
	})
}

// UnscopedJobRunsWithStatus passes all JobRuns to a callback, one by one,
// including those that were soft deleted.
func (orm *ORM) UnscopedJobRunsWithStatus(cb func(*models.JobRun), statuses ...models.RunStatus) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	var runIDs []string
	err := orm.DB.Unscoped().
		Table("job_runs").
		Where("status IN (?)", statuses).
		Order("created_at asc").
		Pluck("id", &runIDs).Error
	if err != nil {
		return errors.Wrap(err, "finding job ids")
	}

	return postgres.Batch(func(offset, limit uint) (uint, error) {
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
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return false, err
	}
	db := orm.DB
	var taskSpec models.TaskSpec
	err := db.Where("type = ?", taskTypeName).First(&taskSpec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "error looking for job of type %s", taskTypeName)
	}
	return true, nil
}

func (orm *ORM) FindJobIDsWithBridge(bridgeName string) ([]models.JobID, error) {
	// Non-FM jobs specify bridges in task specs.
	var bridgeJobIDs []models.JobID
	err := orm.DB.Raw(`SELECT job_spec_id FROM task_specs WHERE type = ? AND deleted_at IS NULL`, bridgeName).Find(&bridgeJobIDs).Error
	if err != nil {
		return nil, err
	}
	var fmInitiators []models.Initiator
	err = orm.DB.Raw(`SELECT * FROM initiators WHERE type = ? AND deleted_at IS NULL`, models.InitiatorFluxMonitor).Find(&fmInitiators).Error
	if err != nil {
		return nil, err
	}
	// FM jobs specify bridges in the initiator.
	for _, fmi := range fmInitiators {
		if fmi.Feeds.IsArray() && len(fmi.Feeds.Array()) == 1 {
			bn := fmi.Feeds.Array()[0].Get("bridge")
			if bn.String() == bridgeName {
				bridgeJobIDs = append(bridgeJobIDs, fmi.JobSpecID)
			}
		}
	}
	return bridgeJobIDs, nil
}

// IdempotentInsertEthTaskRunTx creates both eth_task_run_transaction and eth_tx in one hit
// It can be called multiple times without error as long as the outcome would have resulted in the same database state
func (orm *ORM) IdempotentInsertEthTaskRunTx(meta models.EthTxMeta, fromAddress common.Address, toAddress common.Address, encodedPayload []byte, gasLimit uint64) error {
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	etx := bulletprooftxmanager.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          assets.NewEthValue(0),
		GasLimit:       gasLimit,
		State:          bulletprooftxmanager.EthTxUnstarted,
		Meta:           datatypes.JSON(metaBytes),
	}
	ethTaskRunTransaction := bulletprooftxmanager.EthTaskRunTx{
		TaskRunID: meta.TaskRunID,
	}
	err = orm.DB.Transaction(func(dbtx *gorm.DB) error {
		if err = dbtx.Create(&etx).Error; err != nil {
			return err
		}
		ethTaskRunTransaction.EthTxID = etx.ID
		return dbtx.Create(&ethTaskRunTransaction).Error
	})
	switch v := err.(type) {
	case *pgconn.PgError:
		if v.ConstraintName == "idx_eth_task_run_txes_task_run_id" {
			savedRecord, e := orm.FindEthTaskRunTxByTaskRunID(meta.TaskRunID)
			if e != nil {
				return e
			}
			t := savedRecord.EthTx
			if t.ToAddress != toAddress || !bytes.Equal(t.EncodedPayload, encodedPayload) {
				return fmt.Errorf(
					"transaction already exists for task run ID %s but it has different parameters\n"+
						"New parameters: toAddress: %s, encodedPayload: 0x%s"+
						"Existing record has: toAddress: %s, encodedPayload: 0x%s",
					meta.TaskRunID.String(),
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

// FindEthTaskRunTxByTaskRunID finds the EthTaskRunTx with its EthTxes and EthTxAttempts preloaded
func (orm *ORM) FindEthTaskRunTxByTaskRunID(taskRunID uuid.UUID) (*bulletprooftxmanager.EthTaskRunTx, error) {
	etrt := &bulletprooftxmanager.EthTaskRunTx{}
	err := orm.DB.Preload("EthTx").First(etrt, "task_run_id = ?", &taskRunID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return etrt, err
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

// MarkRan will set Ran to true for a given initiator
func (orm *ORM) MarkRan(i models.Initiator, ran bool) error {
	return orm.convenientTransaction(func(dbtx *gorm.DB) error {
		var newi models.Initiator
		if err := dbtx.Select("ran").First(&newi, "ID = ?", i.ID).Error; err != nil {
			return err
		}

		if ran && newi.Ran {
			return fmt.Errorf("initiator %v for job spec %s has already been run", i.ID, i.JobSpecID.String())
		}

		return dbtx.Model(i).UpdateColumn("ran", ran).Error
	})
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

// JobsSorted returns many JobSpecs sorted by CreatedAt from the store adhering
// to the passed parameters.
func (orm *ORM) JobsSorted(sort SortType, offset int, limit int) ([]models.JobSpec, int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, 0, err
	}
	count, err := orm.CountOf(&models.JobSpec{})
	if err != nil {
		return nil, 0, err
	}

	var jobs []models.JobSpec
	order := fmt.Sprintf("created_at %s", sort.String())
	err = orm.getRecords(&jobs, order, offset, limit)
	return jobs, count, err
}

// JobRunsSorted returns job runs ordered and filtered by the passed params.
func (orm *ORM) JobRunsSorted(sort SortType, offset int, limit int) ([]models.JobRun, int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, 0, err
	}
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
func (orm *ORM) JobRunsSortedFor(id models.JobID, order SortType, offset int, limit int) ([]models.JobRun, int, int, int, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return nil, 0, 0, 0, err
	}
	count, err := orm.JobRunsCountFor(id)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	completedCount, err := orm.JobRunsCountForGivenStatus(id, "completed")
	if err != nil {
		return nil, 0, 0, 0, err
	}
	erroredCount, err := orm.JobRunsCountForGivenStatus(id, "errored")
	if err != nil {
		return nil, 0, 0, 0, err
	}
	var runs []models.JobRun
	err = orm.preloadJobRuns().
		Where("job_spec_id = ?", id).
		Order(fmt.Sprintf("created_at %s", order.String())).
		Limit(limit).
		Offset(offset).
		Find(&runs).Error
	return runs, count, completedCount, erroredCount, err
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

// CreateInitiator saves the initiator.
func (orm *ORM) CreateInitiator(initr *models.Initiator) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.Create(initr).Error
}

// FindOrCreateFluxMonitorRoundStats find the round stats record for a given oracle on a given round, or creates
// it if no record exists
func (orm *ORM) FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32) (stats models.FluxMonitorRoundStats, err error) {
	if err = orm.MustEnsureAdvisoryLock(); err != nil {
		return stats, err
	}
	err = orm.DB.FirstOrCreate(&stats, models.FluxMonitorRoundStats{Aggregator: aggregator, RoundID: roundID}).Error
	return stats, err
}

// DeleteFluxMonitorRoundsBackThrough deletes all the RoundStat records for a given oracle address
// starting from the most recent round back through the given round
func (orm *ORM) DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.Exec(`
        DELETE FROM flux_monitor_round_stats
        WHERE aggregator = ?
          AND round_id >= ?
    `, aggregator, roundID).Error
}

// MostRecentFluxMonitorRoundID finds roundID of the most recent round that the provided oracle
// address submitted to
func (orm *ORM) MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error) {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return 0, err
	}
	var stats models.FluxMonitorRoundStats
	err := orm.DB.Order("round_id DESC").First(&stats, "aggregator = ?", aggregator).Error
	if err != nil {
		return 0, err
	}
	return stats.RoundID, nil
}

// UpdateFluxMonitorRoundStats trys to create a RoundStat record for the given oracle
// at the given round. If one already exists, it increments the num_submissions column.
func (orm *ORM) UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, jobRunID uuid.UUID) error {
	if err := orm.MustEnsureAdvisoryLock(); err != nil {
		return err
	}
	return orm.DB.Exec(`
        INSERT INTO flux_monitor_round_stats (
            aggregator, round_id, job_run_id, num_new_round_logs, num_submissions
        ) VALUES (
            ?, ?, ?, 0, 1
        ) ON CONFLICT (aggregator, round_id)
        DO UPDATE SET
					num_submissions = flux_monitor_round_stats.num_submissions + 1,
					job_run_id = EXCLUDED.job_run_id
    `, aggregator, roundID, jobRunID).Error
}

const removeUnstartedJobRunsQuery = `
WITH deleted_job_runs AS (
	DELETE FROM job_runs as jr
	USING
		task_runs as tr
	WHERE
		jr.status = 'unstarted' AND
		jr.id = tr.job_run_id
	RETURNING
		jr.result_id as job_run_result_id,
		tr.result_id as task_run_result_id,
		jr.run_request_id
),
deleted_job_run_results AS (
	DELETE FROM run_results WHERE id IN (SELECT job_run_result_id FROM deleted_job_runs)
),
deleted_task_run_results AS (
	DELETE FROM run_results WHERE id IN (SELECT task_run_result_id FROM deleted_job_runs)
),
preserved_job_runs AS (
	SELECT run_request_id
	FROM job_runs
	WHERE run_request_id NOT IN (SELECT run_request_id FROM deleted_job_runs)
)
DELETE FROM run_requests
WHERE
	id IN (SELECT run_request_id FROM deleted_job_runs) AND
	id NOT IN (SELECT run_request_id FROM preserved_job_runs)
`

const removeUnstartedTransactionsQuery = `
WITH deleted_job_runs AS (
	DELETE FROM job_runs AS jr
	USING
		task_runs AS tr,
		eth_txes AS tx,
		eth_task_run_txes AS tx_tr
	WHERE
		tx.state = 'unstarted' AND
		tx_tr.eth_tx_id = tx.id AND
		tx_tr.task_run_id = tr.id AND
		jr.id = tr.job_run_id
	RETURNING
		jr.result_id as job_run_result_id,
		tr.result_id as task_run_result_id,
		jr.run_request_id
),
deleted_job_run_results AS (
	DELETE FROM run_results WHERE id IN (SELECT job_run_result_id FROM deleted_job_runs)
),
deleted_task_run_results AS (
	DELETE FROM run_results WHERE id IN (SELECT task_run_result_id FROM deleted_job_runs)
),
preserved_job_runs AS (
	SELECT run_request_id
	FROM job_runs
	WHERE run_request_id NOT IN (SELECT run_request_id FROM deleted_job_runs)
)
DELETE FROM run_requests
WHERE
	id IN (SELECT run_request_id FROM deleted_job_runs) AND
	id NOT IN (SELECT run_request_id FROM preserved_job_runs)
`

func (orm *ORM) RemoveUnstartedTransactions() error {
	return orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(removeUnstartedTransactionsQuery).Error; err != nil {
			return err
		}
		return tx.Exec(removeUnstartedJobRunsQuery).Error
	})
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
