package job

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/lib/pq"

	"gorm.io/gorm/clause"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	storm "github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	ErrNoSuchPeerID             = errors.New("no such peer id exists")
	ErrNoSuchKeyBundle          = errors.New("no such key bundle exists")
	ErrNoSuchTransmitterAddress = errors.New("no such transmitter address exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

var (
	ErrViolatesForeignKeyConstraint = errors.New("violates foreign key constraint")
)

type ORM interface {
	ListenForNewJobs() (postgres.Subscription, error)
	ListenForDeletedJobs() (postgres.Subscription, error)
	ClaimUnclaimedJobs(ctx context.Context) ([]Job, error)
	CreateJob(ctx context.Context, jobSpec *Job, taskDAG pipeline.TaskDAG) error
	JobsV2() ([]Job, error)
	FindJob(id int32) (Job, error)
	FindJobIDsWithBridge(name string) ([]int32, error)
	DeleteJob(ctx context.Context, id int32) error
	RecordError(ctx context.Context, jobID int32, description string)
	UnclaimJob(ctx context.Context, id int32) error
	CheckForDeletedJobs(ctx context.Context) (deletedJobIDs []int32, err error)
	Close() error
	PipelineRunsByJobID(jobID int32, offset, size int) ([]pipeline.Run, int, error)
}

type orm struct {
	db                  *gorm.DB
	config              *storm.Config
	advisoryLocker      postgres.AdvisoryLocker
	advisoryLockClassID int32
	pipelineORM         pipeline.ORM
	eventBroadcaster    postgres.EventBroadcaster
	claimedJobs         map[int32]Job
	claimedJobsMu       *sync.RWMutex
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, config *storm.Config, pipelineORM pipeline.ORM, eventBroadcaster postgres.EventBroadcaster, advisoryLocker postgres.AdvisoryLocker) *orm {
	return &orm{
		db:                  db,
		config:              config,
		advisoryLocker:      advisoryLocker,
		advisoryLockClassID: postgres.AdvisoryLockClassID_JobSpawner,
		pipelineORM:         pipelineORM,
		eventBroadcaster:    eventBroadcaster,
		claimedJobs:         make(map[int32]Job),
		claimedJobsMu:       new(sync.RWMutex),
	}
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
	err := o.db.
		Joins(join, args...).
		Preload("FluxMonitorSpec").
		Preload("DirectRequestSpec").
		Preload("OffchainreportingOracleSpec").
		Preload("KeeperSpec").
		Preload("PipelineSpec").
		Find(&newlyClaimedJobs).Error
	if err != nil {
		return nil, errors.Wrap(err, "ClaimUnclaimedJobs failed to load jobs")
	}

	for _, job := range newlyClaimedJobs {
		o.claimedJobs[job.ID] = job
	}

	return newlyClaimedJobs, errors.Wrap(err, "Job Spawner ORM could not load unclaimed job specs")
}

func (o *orm) claimedJobIDs() (ids []int32) {
	ids = []int32{}
	for _, job := range o.claimedJobs {
		ids = append(ids, job.ID)
	}
	return
}

func (o *orm) CreateJob(ctx context.Context, jobSpec *Job, taskDAG pipeline.TaskDAG) error {
	if taskDAG.HasCycles() {
		return errors.New("task DAG has cycles, which are not permitted")
	}
	tasks, err := taskDAG.TasksInDependencyOrder()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		if task.Type() == pipeline.TaskTypeBridge {
			// Bridge must exist
			name := task.(*pipeline.BridgeTask).Name
			bt := models.BridgeType{}
			if err := o.db.First(&bt, "name = ?", name).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.Wrap(pipeline.ErrNoSuchBridge, name)
				}
				return err
			}
		}
	}

	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	return postgres.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		pipelineSpecID, err := o.pipelineORM.CreateSpec(ctx, tx, taskDAG, jobSpec.MaxTaskDuration)
		if err != nil {
			return errors.Wrap(err, "failed to create pipeline spec")
		}
		jobSpec.PipelineSpecID = pipelineSpecID

		if jobSpec.DirectRequestSpec != nil {
			err = tx.FirstOrCreate(&jobSpec.DirectRequestSpec).Error
			if err != nil {
				return errors.Wrap(err, "error creating direct request spec")
			}
			jobSpec.DirectRequestSpecID = &jobSpec.DirectRequestSpec.ID
		}

		err = tx.Omit("DirectRequestSpec").Create(jobSpec).Error
		pqErr, ok := err.(*pgconn.PgError)
		if err != nil && ok && pqErr.Code == "23503" {
			if pqErr.ConstraintName == "offchainreporting_oracle_specs_p2p_peer_id_fkey" {
				return errors.Wrapf(ErrNoSuchPeerID, "%v", jobSpec.OffchainreportingOracleSpec.P2PPeerID)
			}
			if jobSpec.OffchainreportingOracleSpec != nil && !jobSpec.OffchainreportingOracleSpec.IsBootstrapPeer {
				if pqErr.ConstraintName == "offchainreporting_oracle_specs_transmitter_address_fkey" {
					return errors.Wrapf(ErrNoSuchTransmitterAddress, "%v", jobSpec.OffchainreportingOracleSpec.TransmitterAddress)
				}
				if pqErr.ConstraintName == "offchainreporting_oracle_specs_encrypted_ocr_key_bundle_id_fkey" {
					return errors.Wrapf(ErrNoSuchKeyBundle, "%v", jobSpec.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
				}
			}
		}
		return errors.Wrap(err, "failed to create job")
	})
}

// DeleteJob removes a job that is claimed by this orm
func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	err := o.db.Exec(`
			WITH deleted_jobs AS (
				DELETE FROM jobs WHERE id = ? RETURNING offchainreporting_oracle_spec_id, pipeline_spec_id, keeper_spec_id
			),
			deleted_oracle_specs AS (
				DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT offchainreporting_oracle_spec_id FROM deleted_jobs)
			),
			deleted_keeper_specs AS (
				DELETE FROM keeper_specs WHERE id IN (SELECT keeper_spec_id FROM deleted_jobs)
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
	var claimedJobIDs []int32 = o.claimedJobIDs()

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
	err := o.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "job_id"}, {Name: "description"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"occurrences": gorm.Expr("job_spec_errors_v2.occurrences + 1"),
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

// OffChainReportingJobs returns job specs
func (o *orm) JobsV2() ([]Job, error) {
	var jobs []Job
	err := o.db.
		Preload("PipelineSpec").
		Preload("OffchainreportingOracleSpec").
		Preload("DirectRequestSpec").
		Preload("FluxMonitorSpec").
		Preload("JobSpecErrors").
		Preload("KeeperSpec").
		Find(&jobs).
		Error
	for i := range jobs {
		if jobs[i].OffchainreportingOracleSpec != nil {
			jobs[i].OffchainreportingOracleSpec = loadDynamicConfigVars(o.config, *jobs[i].OffchainreportingOracleSpec)
		}
	}
	return jobs, err
}

func loadDynamicConfigVars(cfg *storm.Config, os OffchainReportingOracleSpec) *OffchainReportingOracleSpec {
	// Load dynamic variables
	return &OffchainReportingOracleSpec{
		IDEmbed: IDEmbed{
			os.ID,
		},
		ContractAddress:                        os.ContractAddress,
		P2PPeerID:                              os.P2PPeerID,
		P2PBootstrapPeers:                      os.P2PBootstrapPeers,
		IsBootstrapPeer:                        os.IsBootstrapPeer,
		EncryptedOCRKeyBundleID:                os.EncryptedOCRKeyBundleID,
		TransmitterAddress:                     os.TransmitterAddress,
		ObservationTimeout:                     models.Interval(cfg.OCRObservationTimeout(time.Duration(os.ObservationTimeout))),
		BlockchainTimeout:                      models.Interval(cfg.OCRBlockchainTimeout(time.Duration(os.BlockchainTimeout))),
		ContractConfigTrackerSubscribeInterval: models.Interval(cfg.OCRContractSubscribeInterval(time.Duration(os.ContractConfigTrackerSubscribeInterval))),
		ContractConfigTrackerPollInterval:      models.Interval(cfg.OCRContractPollInterval(time.Duration(os.ContractConfigTrackerPollInterval))),
		ContractConfigConfirmations:            cfg.OCRContractConfirmations(os.ContractConfigConfirmations),
		CreatedAt:                              os.CreatedAt,
		UpdatedAt:                              os.UpdatedAt,
	}
}

// FindJob returns job by ID
func (o *orm) FindJob(id int32) (Job, error) {
	var job Job
	err := o.db.
		Preload("PipelineSpec").
		Preload("OffchainreportingOracleSpec").
		Preload("FluxMonitorSpec").
		Preload("DirectRequestSpec").
		Preload("JobSpecErrors").
		Preload("KeeperSpec").
		First(&job, "jobs.id = ?", id).
		Error
	if job.OffchainreportingOracleSpec != nil {
		job.OffchainreportingOracleSpec = loadDynamicConfigVars(o.config, *job.OffchainreportingOracleSpec)
	}
	return job, err
}

func (o *orm) FindJobIDsWithBridge(name string) ([]int32, error) {
	var jobs []Job
	err := o.db.Preload("PipelineSpec").Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	var jids []int32
	for _, job := range jobs {
		d := pipeline.TaskDAG{}
		err = d.UnmarshalText([]byte(job.PipelineSpec.DotDagSource))
		if err != nil {
			return nil, err
		}
		tasks, err := d.TasksInDependencyOrder()
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			if task.Type() == pipeline.TaskTypeBridge {
				if task.(*pipeline.BridgeTask).Name == name {
					jids = append(jids, job.ID)
				}
			}
		}
	}
	return jids, nil
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
				Where(`pipeline_task_runs.type != 'result'`).
				Order("created_at ASC, id ASC")
		}).
		Joins("INNER JOIN jobs ON pipeline_runs.pipeline_spec_id = jobs.pipeline_spec_id").
		Where("jobs.id = ?", jobID).
		Limit(size).
		Offset(offset).
		Order("created_at DESC, id DESC").
		Find(&pipelineRuns).
		Error

	return pipelineRuns, int(count), err
}
