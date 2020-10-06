package job

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	ListenForNewJobs() (*utils.PostgresEventListener, error)
	ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error)
	CreateJob(jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error
	DeleteJob(ctx context.Context, id int32) error
	Close() error
}

type orm struct {
	db           *gorm.DB
	uri          string
	advisoryLock *utils.PostgresAdvisoryLock
	pipelineORM  pipeline.ORM
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, uri string, pipelineORM pipeline.ORM) *orm {
	return &orm{db, uri, &utils.PostgresAdvisoryLock{URI: uri}, pipelineORM}
}

func (o *orm) Close() error {
	return o.advisoryLock.Close()
}

const (
	postgresChannelJobCreated = "insert_on_jobs"
)

func (o *orm) ListenForNewJobs() (*utils.PostgresEventListener, error) {
	listener := &utils.PostgresEventListener{
		URI:                  o.uri,
		Event:                postgresChannelJobCreated,
		MinReconnectInterval: 1 * time.Second,
		MaxReconnectDuration: 1 * time.Minute,
	}
	err := listener.Start()
	if err != nil {
		return nil, errors.Wrap(err, "could not start postgres event listener")
	}
	return listener, nil
}

// ClaimUnclaimedJobs returns all currently unlocked jobs, with the lock taken
func (o *orm) ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error) {
	var unclaimedJobs []models.JobSpecV2
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		var maybeJobs []models.JobSpecV2
		err := o.db.
			Preload("OffchainreportingOracleSpec").
			Find(&maybeJobs).Error
		if err != nil {
			return errors.Wrap(err, "ClaimUnclaimedJobs failed to load jobs")
		}

		for _, job := range maybeJobs {
			err = o.advisoryLock.TryLock(ctx, utils.AdvisoryLockClassID_JobSpawner, job.ID)
			if err == nil {
				unclaimedJobs = append(unclaimedJobs, job)
			}
		}
		return nil
	})
	return unclaimedJobs, errors.Wrap(err, "Job Spawner ORM could not load unclaimed job specs")
}

func (o *orm) CreateJob(jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error {
	if taskDAG.HasCycles() {
		return errors.New("task DAG has cycles, which are not permitted")
	}
	return utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		pipelineSpecID, err := o.pipelineORM.CreateSpec(taskDAG)
		if err != nil {
			return err
		}
		jobSpec.PipelineSpecID = pipelineSpecID

		err = tx.Create(jobSpec).Error
		if err != nil && err.Error() != "sql: no rows in result set" {
			return err
		}
		return nil
	})
}

func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	return utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		// If we can take the advisory lock, that means either we own this job or
		// nobody does.  That gives us permission to delete the job.  Note that we
		// have to unlock twice at the end (as we already have it).
		err := o.advisoryLock.TryLock(ctx, utils.AdvisoryLockClassID_JobSpawner, id)
		if err != nil {
			return err
		}
		defer o.advisoryLock.Unlock(ctx, utils.AdvisoryLockClassID_JobSpawner, id)
		defer o.advisoryLock.Unlock(ctx, utils.AdvisoryLockClassID_JobSpawner, id)
		return o.db.Exec(`
            WITH deleted_jobs AS (
                DELETE FROM jobs WHERE id = ? RETURNING offchainreporting_oracle_spec_id
            )
            DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT id FROM deleted_jobs)
        `, id).Error
	})
}
