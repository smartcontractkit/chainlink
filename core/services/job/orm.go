package job

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/sqlx"

	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNoSuchPeerID             = errors.New("no such peer id exists")
	ErrNoSuchKeyBundle          = errors.New("no such key bundle exists")
	ErrNoSuchTransmitterAddress = errors.New("no such transmitter address exists")
	ErrNoSuchPublicKey          = errors.New("no such public key exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

var (
	ErrViolatesForeignKeyConstraint = errors.New("violates foreign key constraint")
)

type ORM interface {
	ListenForNewJobs() (postgres.Subscription, error)
	ListenForDeletedJobs() (postgres.Subscription, error)
	ClaimUnclaimedJobs(ctx context.Context) ([]Job, error)
	CreateJob(ctx context.Context, jobSpec *Job, pipeline pipeline.Pipeline) (Job, error)
	JobsV2(offset, limit int) ([]Job, int, error)
	FindJobTx(id int32) (Job, error)
	FindJob(ctx context.Context, id int32) (Job, error)
	FindJobIDsWithBridge(name string) ([]int32, error)
	DeleteJob(ctx context.Context, id int32) error
	RecordError(ctx context.Context, jobID int32, description string)
	DismissError(ctx context.Context, errorID int32) error
	UnclaimJob(ctx context.Context, id int32) error
	CheckForDeletedJobs(ctx context.Context) (deletedJobIDs []int32, err error)
	Close() error
	PipelineRuns(offset, size int) ([]pipeline.Run, int, error)
	PipelineRunsByJobID(jobID int32, offset, size int) ([]pipeline.Run, int, error)
}

type orm struct {
	db                  *gorm.DB
	config              Config
	keyStore            keystore.Master
	advisoryLocker      postgres.AdvisoryLocker
	advisoryLockClassID int32
	pipelineORM         pipeline.ORM
	eventBroadcaster    postgres.EventBroadcaster
	claimedJobs         map[int32]Job
	claimedJobsMu       *sync.RWMutex
}

var _ ORM = (*orm)(nil)

func NewORM(
	db *gorm.DB,
	cfg Config,
	pipelineORM pipeline.ORM,
	eventBroadcaster postgres.EventBroadcaster,
	advisoryLocker postgres.AdvisoryLocker,
	keyStore keystore.Master, // needed to validation key properties on new job creation
) *orm {
	return &orm{
		db:                  db,
		config:              cfg,
		keyStore:            keyStore,
		advisoryLocker:      advisoryLocker,
		advisoryLockClassID: postgres.AdvisoryLockClassID_JobSpawner,
		pipelineORM:         pipelineORM,
		eventBroadcaster:    eventBroadcaster,
		claimedJobs:         make(map[int32]Job),
		claimedJobsMu:       new(sync.RWMutex),
	}
}

func PreloadAllJobTypes(db *gorm.DB) *gorm.DB {
	return db.
		Preload("PipelineSpec").
		Preload("FluxMonitorSpec").
		Preload("DirectRequestSpec").
		Preload("OffchainreportingOracleSpec").
		Preload("KeeperSpec").
		Preload("PipelineSpec").
		Preload("CronSpec").
		Preload("WebhookSpec").
		Preload("VRFSpec")
}

func (o *orm) Close() error {
	return nil
}

func (o *orm) ListenForNewJobs() (postgres.Subscription, error) {
	return o.eventBroadcaster.Subscribe(postgres.ChannelJobCreated, "")
}

func (o *orm) ListenForDeletedJobs() (postgres.Subscription, error) {
	return o.eventBroadcaster.Subscribe(postgres.ChannelJobDeleted, "")
}

// ClaimUnclaimedJobs locks all currently unlocked jobs and returns all jobs locked by this process
func (o *orm) ClaimUnclaimedJobs(ctx context.Context) ([]Job, error) {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	claimedJobIDs := o.claimedJobIDs()

	var join string
	var args []interface{}
	if len(claimedJobIDs) > 0 {
		// NOTE: OFFSET 0 is a postgres trick that doesn't change the result,
		// but prevents the optimiser from trying to pull the where condition
		// up out of the subquery
		join = `
            INNER JOIN (
                SELECT not_claimed_by_us.id, pg_try_advisory_lock(?::integer, not_claimed_by_us.id) AS locked
                FROM (SELECT id FROM jobs WHERE NOT (id = ANY(?)) OFFSET 0) not_claimed_by_us
            ) claimed_jobs ON jobs.id = claimed_jobs.id AND claimed_jobs.locked
        `
		args = []interface{}{o.advisoryLockClassID, pq.Array(claimedJobIDs)}
	} else {
		join = `
            INNER JOIN (
                SELECT not_claimed_by_us.id, pg_try_advisory_lock(?::integer, not_claimed_by_us.id) AS locked
                FROM jobs not_claimed_by_us
            ) claimed_jobs ON jobs.id = claimed_jobs.id AND claimed_jobs.locked
        `
		args = []interface{}{o.advisoryLockClassID}
	}

	var newlyClaimedJobs []Job
	err := postgres.GormTransactionWithDefaultContext(o.db, func(tx *gorm.DB) error {
		err := PreloadAllJobTypes(tx.
			Joins(join, args...)).
			Find(&newlyClaimedJobs).Error
		if err != nil {
			return err
		}

		for i := range newlyClaimedJobs {
			o.claimedJobs[newlyClaimedJobs[i].ID] = newlyClaimedJobs[i]
		}
		return nil
	})
	return newlyClaimedJobs, errors.Wrap(err, "Job Spawner ORM could not load unclaimed job specs")
}

func (o *orm) claimedJobIDs() (ids []int32) {
	ids = []int32{}
	for _, job := range o.claimedJobs {
		ids = append(ids, job.ID)
	}
	return
}

// CreateJob creates the job and it's associated spec record.
//
// NOTE: This is not wrapped in a db transaction so if you call this, you should
// use postgres.TransactionManager to create the transaction in the context.
// Expects an unmarshaled job spec as the jobSpec argument i.e. output from ValidatedXX.
// Returns a fully populated Job.
func (o *orm) CreateJob(ctx context.Context, jobSpec *Job, p pipeline.Pipeline) (Job, error) {
	var jb Job
	for _, task := range p.Tasks {
		if task.Type() == pipeline.TaskTypeBridge {
			// Bridge must exist
			name := task.(*pipeline.BridgeTask).Name
			bt := models.BridgeType{}
			if err := o.db.First(&bt, "name = ?", name).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return jb, errors.Wrap(pipeline.ErrNoSuchBridge, name)
				}
				return jb, err
			}
		}
	}

	tx := postgres.TxFromContext(ctx, o.db)

	// Autogenerate a job ID if not specified
	if jobSpec.ExternalJobID == (uuid.UUID{}) {
		jobSpec.ExternalJobID = uuid.NewV4()
	}

	switch jobSpec.Type {
	case DirectRequest:
		err := tx.Create(&jobSpec.DirectRequestSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create DirectRequestSpec for jobSpec")
		}
		jobSpec.DirectRequestSpecID = &jobSpec.DirectRequestSpec.ID
	case FluxMonitor:
		err := tx.Create(&jobSpec.FluxMonitorSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create FluxMonitorSpec for jobSpec")
		}
		jobSpec.FluxMonitorSpecID = &jobSpec.FluxMonitorSpec.ID
	case OffchainReporting:
		if jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID != nil {
			_, err := o.keyStore.OCR().Get(jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID.String())
			if err != nil {
				return jb, errors.Wrapf(ErrNoSuchKeyBundle, "%v", jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
			}
		}
		if jobSpec.OffchainreportingOracleSpec.P2PPeerID != nil {
			_, err := o.keyStore.P2P().Get(jobSpec.OffchainreportingOracleSpec.P2PPeerID.Raw())
			if err != nil {
				return jb, errors.Wrapf(ErrNoSuchPeerID, "%v", jobSpec.OffchainreportingOracleSpec.P2PPeerID)
			}
		}
		if jobSpec.OffchainreportingOracleSpec.TransmitterAddress != nil {
			_, err := o.keyStore.Eth().Get(jobSpec.OffchainreportingOracleSpec.TransmitterAddress.Hex())
			if err != nil {
				return jb, errors.Wrapf(ErrNoSuchTransmitterAddress, "%v", jobSpec.OffchainreportingOracleSpec.TransmitterAddress)
			}
		}

		err := tx.Create(&jobSpec.OffchainreportingOracleSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create OffchainreportingOracleSpec for jobSpec")
		}
		jobSpec.OffchainreportingOracleSpecID = &jobSpec.OffchainreportingOracleSpec.ID
	case Keeper:
		err := tx.Create(&jobSpec.KeeperSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create KeeperSpec for jobSpec")
		}
		jobSpec.KeeperSpecID = &jobSpec.KeeperSpec.ID
	case Cron:
		err := tx.Create(&jobSpec.CronSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create CronSpec for jobSpec")
		}
		jobSpec.CronSpecID = &jobSpec.CronSpec.ID
	case VRF:
		err := tx.Create(&jobSpec.VRFSpec).Error
		pqErr, ok := err.(*pgconn.PgError)
		if err != nil && ok && pqErr.Code == "23503" {
			if pqErr.ConstraintName == "vrf_specs_public_key_fkey" {
				return jb, errors.Wrapf(ErrNoSuchPublicKey, "%s", jobSpec.VRFSpec.PublicKey.String())
			}
		}
		if err != nil {
			return jb, errors.Wrap(err, "failed to create VRFSpec for jobSpec")
		}
		jobSpec.VRFSpecID = &jobSpec.VRFSpec.ID
	case Webhook:
		err := tx.Create(&jobSpec.WebhookSpec).Error
		if err != nil {
			return jb, errors.Wrap(err, "failed to create WebhookSpec for jobSpec")
		}
		jobSpec.WebhookSpecID = &jobSpec.WebhookSpec.ID
		for i, eiWS := range jobSpec.WebhookSpec.ExternalInitiatorWebhookSpecs {
			jobSpec.WebhookSpec.ExternalInitiatorWebhookSpecs[i].WebhookSpecID = jobSpec.WebhookSpec.ID
			err := tx.Create(&jobSpec.WebhookSpec.ExternalInitiatorWebhookSpecs[i]).Error
			if err != nil {
				return jb, errors.Wrapf(err, "failed to create ExternalInitiatorWebhookSpec for WebhookSpec: %#v", eiWS)
			}
		}
	default:
		logger.Fatalf("Unsupported jobSpec.Type: %v", jobSpec.Type)
	}

	pipelineSpecID, err := o.pipelineORM.CreateSpec(ctx, tx, p, jobSpec.MaxTaskDuration)
	if err != nil {
		return jb, errors.Wrap(err, "failed to create pipeline spec")
	}
	jobSpec.PipelineSpecID = pipelineSpecID
	err = tx.Create(jobSpec).Error
	if err != nil {
		return jb, errors.Wrap(err, "failed to create job")
	}

	return o.FindJob(ctx, jobSpec.ID)
}

// DeleteJob removes a job that is claimed by this orm
func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	err := o.db.Exec(`
		WITH deleted_jobs AS (
			DELETE FROM jobs WHERE id = ? RETURNING
				pipeline_spec_id,
				offchainreporting_oracle_spec_id,
				keeper_spec_id,
				cron_spec_id,
				flux_monitor_spec_id,
				vrf_spec_id,
				webhook_spec_id,
				direct_request_spec_id
		),
		deleted_oracle_specs AS (
			DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT offchainreporting_oracle_spec_id FROM deleted_jobs)
		),
		deleted_keeper_specs AS (
			DELETE FROM keeper_specs WHERE id IN (SELECT keeper_spec_id FROM deleted_jobs)
		),
		deleted_cron_specs AS (
			DELETE FROM cron_specs WHERE id IN (SELECT cron_spec_id FROM deleted_jobs)
		),
		deleted_fm_specs AS (
			DELETE FROM flux_monitor_specs WHERE id IN (SELECT flux_monitor_spec_id FROM deleted_jobs)
		),
		deleted_vrf_specs AS (
			DELETE FROM vrf_specs WHERE id IN (SELECT vrf_spec_id FROM deleted_jobs)
		),
		deleted_webhook_specs AS (
			DELETE FROM webhook_specs WHERE id IN (SELECT webhook_spec_id FROM deleted_jobs)
		),
		deleted_dr_specs AS (
			DELETE FROM direct_request_specs WHERE id IN (SELECT direct_request_spec_id FROM deleted_jobs)
		)
		DELETE FROM pipeline_specs WHERE id IN (SELECT pipeline_spec_id FROM deleted_jobs)
	`, id).Error
	if err != nil {
		return errors.Wrap(err, "DeleteJob failed to delete job")
	}

	if err := o.unclaimJob(ctx, id); err != nil {
		return errors.Wrap(err, "DeleteJob failed to unclaim job")
	}

	return nil
}

func (o *orm) CheckForDeletedJobs(ctx context.Context) (deletedJobIDs []int32, err error) {
	o.claimedJobsMu.RLock()
	defer o.claimedJobsMu.RUnlock()
	var claimedJobIDs = o.claimedJobIDs()

	rows, err := o.db.Raw(`SELECT id FROM jobs WHERE id = ANY(?)`, pq.Array(claimedJobIDs)).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "could not query for jobs")
	}
	defer logger.ErrorIfCalling(rows.Close)

	foundJobs := make(map[int32]struct{})
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "could not scan row")
		}
		foundJobs[id] = struct{}{}
	}

	var deletedClaimedJobs []int32

	for _, claimedID := range claimedJobIDs {
		if _, ok := foundJobs[claimedID]; !ok {
			deletedClaimedJobs = append(deletedClaimedJobs, claimedID)
		}
	}

	return deletedClaimedJobs, nil
}

func (o *orm) UnclaimJob(ctx context.Context, id int32) error {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()
	return o.unclaimJob(ctx, id)
}

func (o *orm) unclaimJob(ctx context.Context, id int32) error {
	if _, ok := o.claimedJobs[id]; ok {
		delete(o.claimedJobs, id)
		return errors.Wrap(o.advisoryLocker.Unlock(ctx, o.advisoryLockClassID, id), "DeleteJob failed to unlock job")
	}
	return nil
}

func (o *orm) RecordError(ctx context.Context, jobID int32, description string) {
	pse := SpecError{JobID: jobID, Description: description, Occurrences: 1}
	err := o.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "job_id"}, {Name: "description"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"occurrences": gorm.Expr("job_spec_errors.occurrences + 1"),
				"updated_at":  gorm.Expr("excluded.updated_at"),
			}),
		}).
		Create(&pse).
		Error
	// Noop if the job has been deleted.
	if err != nil && strings.Contains(err.Error(), ErrViolatesForeignKeyConstraint.Error()) {
		return
	}
	logger.ErrorIf(err, fmt.Sprintf("error creating SpecError %v", description))
}

func (o *orm) DismissError(ctx context.Context, ID int32) error {
	result := o.db.Exec("DELETE FROM job_spec_errors WHERE id = ?", ID)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (o *orm) JobsV2(offset, limit int) ([]Job, int, error) {
	var count int64
	var jobs []Job
	err := postgres.GormTransactionWithDefaultContext(o.db, func(tx *gorm.DB) error {
		err := tx.
			Model(Job{}).
			Count(&count).
			Error

		if err != nil {
			return err
		}

		err = PreloadAllJobTypes(tx).
			Preload("JobSpecErrors").
			Limit(limit).
			Offset(offset).
			Order("id ASC").
			Find(&jobs).
			Error
		if err != nil {
			return err
		}
		for i := range jobs {
			if jobs[i].OffchainreportingOracleSpec != nil {
				jobs[i].OffchainreportingOracleSpec = LoadDynamicConfigVars(o.config, *jobs[i].OffchainreportingOracleSpec)
			}
		}
		return nil
	})
	return jobs, int(count), err
}

type OCRSpecConfig interface {
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRObservationTimeout() time.Duration
}

func LoadDynamicConfigVars(cfg OCRSpecConfig, os OffchainReportingOracleSpec) *OffchainReportingOracleSpec {
	// Load dynamic variables
	if os.ObservationTimeout == 0 {
		os.ObservationTimeout = models.Interval(cfg.OCRObservationTimeout())
	}
	if os.BlockchainTimeout == 0 {
		os.BlockchainTimeout = models.Interval(cfg.OCRBlockchainTimeout())
	}
	if os.ContractConfigTrackerSubscribeInterval == 0 {
		os.ContractConfigTrackerSubscribeInterval = models.Interval(cfg.OCRContractSubscribeInterval())
	}
	if os.ContractConfigTrackerPollInterval == 0 {
		os.ContractConfigTrackerPollInterval = models.Interval(cfg.OCRContractPollInterval())
	}
	if os.ContractConfigConfirmations == 0 {
		os.ContractConfigConfirmations = cfg.OCRContractConfirmations()
	}
	return &os
}

func (o *orm) FindJobTx(id int32) (Job, error) {
	var jb Job
	var err error
	txm := postgres.NewGormTransactionManager(o.db)
	err = txm.Transact(func(ctx context.Context) error {
		jb, err = o.FindJob(ctx, id)
		return err
	})
	return jb, err
}

// FindJob returns job by ID
func (o *orm) FindJob(ctx context.Context, id int32) (Job, error) {
	var jb Job
	tx := postgres.TxFromContext(ctx, o.db)
	err := PreloadAllJobTypes(tx).
		Preload("JobSpecErrors").
		First(&jb, "jobs.id = ?", id).
		Error
	if err != nil {
		return jb, err
	}

	if jb.OffchainreportingOracleSpec != nil {
		jb.OffchainreportingOracleSpec = LoadDynamicConfigVars(o.config, *jb.OffchainreportingOracleSpec)
	}
	return jb, nil
}

func (o *orm) FindJobIDsWithBridge(name string) ([]int32, error) {
	var jobs []Job
	err := o.db.Preload("PipelineSpec").Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	var jids []int32
	for _, job := range jobs {
		p, err := pipeline.Parse(job.PipelineSpec.DotDagSource)
		if err != nil {
			return nil, err
		}
		for _, task := range p.Tasks {
			if task.Type() == pipeline.TaskTypeBridge {
				if task.(*pipeline.BridgeTask).Name == name {
					jids = append(jids, job.ID)
				}
			}
		}
	}
	return jids, nil
}

// Preload PipelineSpec.JobID for each Run
func (o *orm) preloadJobIDs(runs []pipeline.Run) error {
	db := postgres.UnwrapGormDB(o.db)

	ids := make([]int32, 0, len(runs))
	for _, run := range runs {
		ids = append(ids, run.PipelineSpecID)
	}

	// construct a WHERE IN query
	sql := `SELECT id, pipeline_spec_id FROM jobs WHERE pipeline_spec_id IN (?);`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	query = db.Rebind(query)
	var results []struct {
		ID             int32
		PipelineSpecID int32
	}
	if err := db.Select(&results, query, args...); err != nil {
		return err
	}

	// fill in fields
	for i := range runs {
		for _, result := range results {
			if result.PipelineSpecID == runs[i].PipelineSpecID {
				runs[i].PipelineSpec.JobID = result.ID
			}
		}
	}

	return nil
}

// PipelineRunsByJobID returns all pipeline runs
func (o *orm) PipelineRuns(offset, size int) ([]pipeline.Run, int, error) {
	var pipelineRuns []pipeline.Run
	var count int64
	err := o.db.
		Model(pipeline.Run{}).
		Count(&count).
		Error

	if err != nil {
		return pipelineRuns, 0, err
	}

	err = o.db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		}).
		Limit(size).
		Offset(offset).
		Order("created_at DESC, id DESC").
		Find(&pipelineRuns).
		Error

	if err != nil {
		return pipelineRuns, int(count), err
	}

	err = o.preloadJobIDs(pipelineRuns)

	return pipelineRuns, int(count), err
}

// PipelineRunsByJobID returns pipeline runs for a job
func (o *orm) PipelineRunsByJobID(jobID int32, offset, size int) ([]pipeline.Run, int, error) {
	var pipelineRuns []pipeline.Run
	var count int64
	err := o.db.
		Model(pipeline.Run{}).
		Joins("INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id").
		Where("jobs.id = ?", jobID).
		Count(&count).
		Error

	if err != nil {
		return pipelineRuns, 0, err
	}

	err = o.db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		}).
		Joins("INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id").
		Where("jobs.id = ?", jobID).
		Limit(size).
		Offset(offset).
		Order("created_at DESC, id DESC").
		Find(&pipelineRuns).
		Error

	// can skip preloadJobIDs since we already know the jobID
	for i := range pipelineRuns {
		pipelineRuns[i].PipelineSpec.JobID = jobID
	}

	return pipelineRuns, int(count), err
}
