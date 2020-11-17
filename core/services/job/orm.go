package job

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

var (
	ErrViolatesForeignKeyConstraint = errors.New("violates foreign key constraint")
)

type ORM interface {
	ListenForNewJobs() (postgres.Subscription, error)
	ListenForDeletedJobs() (postgres.Subscription, error)
	ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error)
	CreateJob(ctx context.Context, jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error
	DeleteJob(ctx context.Context, id int32) error
	RecordError(ctx context.Context, jobID int32, description string)
	UnclaimJob(ctx context.Context, id int32) error
	CheckForDeletedJobs(ctx context.Context) (deletedJobIDs []int32, err error)
	Close() error
}

type orm struct {
	db                  *gorm.DB
	config              Config
	advisoryLocker      postgres.AdvisoryLocker
	advisoryLockClassID int32
	pipelineORM         pipeline.ORM
	eventBroadcaster    postgres.EventBroadcaster
	claimedJobs         map[int32]models.JobSpecV2
	claimedJobsMu       *sync.RWMutex
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, config Config, pipelineORM pipeline.ORM, eventBroadcaster postgres.EventBroadcaster, advisoryLocker postgres.AdvisoryLocker) *orm {
	return &orm{
		db:                  db,
		config:              config,
		advisoryLocker:      advisoryLocker,
		advisoryLockClassID: postgres.AdvisoryLockClassID_JobSpawner,
		pipelineORM:         pipelineORM,
		eventBroadcaster:    eventBroadcaster,
		claimedJobs:         make(map[int32]models.JobSpecV2),
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
func (o *orm) ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error) {
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
                FROM (SELECT id FROM jobs WHERE id != ANY(?) OFFSET 0) not_claimed_by_us
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

	var newlyClaimedJobs []models.JobSpecV2
	err := o.db.
		Joins(join, args...).
		Preload("OffchainreportingOracleSpec").
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

func (o *orm) CreateJob(ctx context.Context, jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error {
	if taskDAG.HasCycles() {
		return errors.New("task DAG has cycles, which are not permitted")
	}

	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	return postgres.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		pipelineSpecID, err := o.pipelineORM.CreateSpec(ctx, taskDAG)
		if err != nil {
			return errors.Wrap(err, "failed to create pipeline spec")
		}
		jobSpec.PipelineSpecID = pipelineSpecID

		err = tx.Create(jobSpec).Error
		return errors.Wrap(err, "failed to create job")
	})
}

// DeleteJob removes a job that is claimed by this orm
func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	err := o.db.Exec(`
            WITH deleted_jobs AS (
            	DELETE FROM jobs WHERE id = $1 RETURNING offchainreporting_oracle_spec_id, pipeline_spec_id
            ),
            deleted_oracle_specs AS (
				DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT offchainreporting_oracle_spec_id FROM deleted_jobs)
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

	rows, err := o.db.DB().QueryContext(ctx, `SELECT id FROM jobs WHERE id = ANY($1)`, pq.Array(claimedJobIDs))
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
	pse := models.JobSpecErrorV2{JobID: jobID, Description: description, Occurrences: 1}
	err := o.db.
		Set(
			"gorm:insert_option",
			`ON CONFLICT (job_id, description)
			DO UPDATE SET occurrences = job_spec_errors_v2.occurrences + 1, updated_at = excluded.updated_at`,
		).
		Create(&pse).
		Error
	// Noop if the job has been deleted.
	if err != nil && strings.Contains(err.Error(), ErrViolatesForeignKeyConstraint.Error()) {
		return
	}
	logger.ErrorIf(err, fmt.Sprintf("error creating JobSpecErrorV2 %v", description))
}
