package pipeline

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	// Runner
	CreateRun(jobSpecID int32) (int64, error)
	WithNextUnclaimedTaskRun(f func(ptRun TaskRun, predecessors []TaskRun) Result) error
	AwaitRun(ctx context.Context, runID int64) error
	// NotifyCompletion(pipelineRunID int64) error
	ResultsForRun(runID int64) ([]Result, error)

	// BridgeTask
	FindBridge(name models.TaskType) (models.BridgeType, error)
}

type orm struct {
	db *gorm.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

// CreateRun adds a Run record to the DB, and one TaskRun
// per TaskSpec associated with the given Spec.  Processing of the
// TaskRuns is maximally parallelized across all of the Chainlink nodes in the
// cluster.
func (o *orm) CreateRun(jobSpecID int32) (int64, error) {
	var runID int64
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		// Fetch the spec ID
		var spec Spec
		err := o.db.
			Where("job_spec_id = ?", jobSpecID).
			First(&spec).
			Error
		if err != nil {
			return err
		}

		// Create the job run
		run := Run{
			PipelineSpecID: spec.ID,
			CreatedAt:      time.Now(),
		}

		err = tx.Create(&run).Error
		if err != nil {
			return errors.Wrap(err, "could not create pipeline run")
		}

		runID = run.ID

		// Create the task runs
		err = tx.Exec(`
            INSERT INTO pipeline_task_runs (
                pipeline_run_id, pipeline_task_spec_id, created_at
            )
            SELECT ? AS pipeline_run_id, id AS pipeline_task_spec_id, NOW() AS created_at
            FROM pipeline_task_specs
            WHERE pipeline_spec_id = ?
        `, run.ID, run.PipelineSpecID).Error
		return errors.Wrap(err, "could not create pipeline task runs")
	})
	return runID, err
}

// WithNextUnclaimedTaskRun chooses any arbitrary incomplete TaskRun from the DB
// whose parent TaskRuns have already been processed.
func (o *orm) WithNextUnclaimedTaskRun(fn func(ptRun TaskRun, predecessors []TaskRun) Result) error {
	return utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		var ptRun TaskRun
		var predecessors []TaskRun

		// NOTE: This could conceivably be done in pure SQL and made marginally more efficient by preloading with joins

		// Find the next unlocked, unfinished pipeline_task_run with no uncompleted predecessors
		err := tx.Table("pipeline_task_runs AS successor").
			Set("gorm:query_option", "FOR UPDATE OF successor SKIP LOCKED"). // Row-level lock
			Joins("LEFT JOIN pipeline_task_runs AS predecessors ON successor.id = predecessor.successor_id AND predecessors.finished_at IS NULL").
			Where("predecessors.id IS NULL").
			Where("successor.finished_at IS NULL").
			Preload("TaskSpec").
			Order("id ASC").
			First(&ptRun).
			Error
		if err != nil {
			return errors.Wrap(err, "error finding next task run")
		}

		// Find all the predecessors
		err = tx.Where("successor_id = ?", ptRun.ID).Find(&predecessors).Error
		if err != nil {
			return errors.Wrap(err, "error finding task run predecessors")
		}

		// Call the callback and convert its output to a format appropriate for the DB
		result := fn(ptRun, predecessors)
		var out *JSONSerializable
		var errString null.String
		if result.Value != nil {
			out = &JSONSerializable{Value: result.Value}
		} else if result.Error != nil {
			errString = null.StringFrom(err.Error())
		}

		// Update the task run record with the output and error
		err = tx.Exec(`
            UPDATE pipeline_task_runs
            SET output = ?, error = ?, finished_at = ?
            WHERE id = ?
        `, out, errString, time.Now(), ptRun.ID).Error

		return errors.Wrap(err, "could not mark pipeline_task_run as finished")
	})
}

func (o *orm) AwaitRun(ctx context.Context, runID int64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var done bool
		err := o.db.Raw(`
            SELECT bool_and(pipeline_task_runs.error IS NOT NULL OR pipeline_task_runs.output IS NOT NULL)
            FROM pipeline_job_runs
            JOIN pipeline_task_runs ON pipeline_task_runs.pipeline_job_run_id = pipeline_job_runs.id
            WHERE pipeline_job_runs.id = $1
        `, runID).Scan(&done).Error
		if err != nil {
			// TODO: log error
			time.Sleep(1 * time.Second)
			continue
		}
		if done {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

// func (o *orm) NotifyCompletion(runID int64) error {
// 	return o.db.Exec(`
//     $$
//     BEGIN
//         IF (SELECT bool_and(pipeline_task_runs.error IS NOT NULL OR pipeline_task_runs.output IS NOT NULL) FROM pipeline_job_runs JOIN pipeline_task_runs ON pipeline_task_runs.pipeline_job_run_id = pipeline_job_runs.id WHERE pipeline_job_runs.id = $1)
//             PERFORM pg_notify('pipeline_job_run_completed', $1::text);
//         END IF;
//     END;
//     $$ LANGUAGE plpgsql;
//     )`, runID).Error
// }

func (o *orm) ResultsForRun(runID int64) ([]Result, error) {
	var taskRuns []TaskRun
	err := o.db.
		Where("pipeline_job_run_id = ?", runID).
		Where("error IS NOT NULL OR output IS NOT NULL").
		Where("successor_id IS NULL").
		Find(&taskRuns).
		Error
	if err != nil {
		return nil, err
	}

	results := make([]Result, len(taskRuns))
	for i, taskRun := range taskRuns {
		results[i] = taskRun.Result()
	}
	return results, nil
}

func (o *orm) FindBridge(name models.TaskType) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, o.db.First(&bt, "name = ?", name.String()).Error
}
