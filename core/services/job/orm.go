package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	ListenForNewJobs() (*utils.PostgresEventListener, error)
	ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error)
	CreateJob(jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error
	DeleteJob(ctx context.Context, id int32) error
	Close() error
}

type orm struct {
	db            *gorm.DB
	uri           string
	advisoryLock  *utils.PostgresAdvisoryLock
	pipelineORM   pipeline.ORM
	claimedJobs   []models.JobSpecV2
	claimedJobsMu *sync.Mutex
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, config Config, pipelineORM pipeline.ORM) *orm {
	return &orm{
		db:            db,
		uri:           config.DatabaseURL(),
		advisoryLock:  utils.NewPostgresAdvisoryLock(config.DatabaseURL()),
		pipelineORM:   pipelineORM,
		claimedJobs:   make([]models.JobSpecV2, 0),
		claimedJobsMu: &sync.Mutex{},
	}
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

// ClaimUnclaimedJobs locks all currently unlocked jobs and returns all jobs locked by this process
func (o *orm) ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error) {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	var ids []int32
	err := o.db.Raw(`SELECT id FROM jobs`).Scan(&ids).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("111 ~>", ids)

	var jobs []models.JobSpecV2
	err = o.db.Find(&jobs).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("222 ~>", jobs)

	err = o.db.Raw(`SELECT id FROM jobs WHERE id NOT IN (1, 2, 3)`).Scan(&ids).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("111 ~>", ids)

	var newlyClaimedJobs []models.JobSpecV2
	err = o.db.
		// NOTE: OFFSET 0 is a postgres trick that doesn't change the result,
		// but prevents the optimiser from trying to pull the where condition
		// up out of the subquery
		Joins(`
			INNER JOIN (
				SELECT not_claimed_by_us.id, pg_try_advisory_lock(?::integer, not_claimed_by_us.id) AS locked
				FROM (
					SELECT id FROM jobs WHERE id != ANY(ARRAY[9, 10]) OFFSET 0
				) not_claimed_by_us
			) claimed_jobs ON jobs.id = claimed_jobs.id AND claimed_jobs.locked
			`, utils.AdvisoryLockClassID_JobSpawner, pq.Array(o.claimedJobIDs())).
		Preload("OffchainreportingOracleSpec").
		Find(&newlyClaimedJobs).Error
	if err != nil {
		return nil, errors.Wrap(err, "ClaimUnclaimedJobs failed to load jobs")
	}

	o.claimedJobs = append(o.claimedJobs, newlyClaimedJobs...)

	return newlyClaimedJobs, errors.Wrap(err, "Job Spawner ORM could not load unclaimed job specs")
}

func (o *orm) claimedJobIDs() (ids []int32) {
	for _, job := range o.claimedJobs {
		ids = append(ids, job.ID)
	}
	return
}

func (o *orm) CreateJob(jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error {
	if taskDAG.HasCycles() {
		return errors.New("task DAG has cycles, which are not permitted")
	}
	return utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		pipelineSpecID, err := o.pipelineORM.CreateSpec(taskDAG)
		if err != nil {
			return errors.Wrap(err, "failed to create pipeline spec")
		}
		jobSpec.PipelineSpecID = pipelineSpecID

		err = tx.Create(jobSpec).Error
		return errors.Wrap(err, "failed to create job")
	})
}

// DeleteJob removes a job that is claimed by this orm
// NOTE: It may be nice to extend this in future so it can delete any job and other nodes handle it gracefully
func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	for i, job := range o.claimedJobs {
		if job.ID == id {
			if _, err := o.db.DB().ExecContext(ctx, `
                WITH deleted_jobs AS (
                	DELETE FROM jobs WHERE id = $1 RETURNING offchainreporting_oracle_spec_id
                )
                DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT offchainreporting_oracle_spec_id FROM deleted_jobs)
			`, id); err != nil {
				return errors.Wrap(err, "DeleteJob failed to delete job")
			}
			if err := o.advisoryLock.Unlock(ctx, utils.AdvisoryLockClassID_JobSpawner, id); err != nil {
				return errors.Wrap(err, "DeleteJob failed to unlock job")
			}
			// Delete the current job from the claimedJobs list
			o.claimedJobs[i] = o.claimedJobs[len(o.claimedJobs)-1] // Copy last element to current position
			o.claimedJobs = o.claimedJobs[:len(o.claimedJobs)-1]   // Truncate slice.
			return nil
		}
	}

	return errors.New("cannot delete job that is not claimed by this orm")
}
