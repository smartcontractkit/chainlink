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
	CreateSpec(taskDAG TaskDAG) (int32, error)
	CreateRun(jobID int32) (int64, error)
	WithNextUnclaimedTaskRun(f func(taskRun TaskRun, predecessors []TaskRun) Result) (bool, error)
	AwaitRun(ctx context.Context, runID int64) error
	// NotifyCompletion(pipelineRunID int64) error
	ResultsForRun(runID int64) ([]Result, error)

	FindBridge(name models.TaskType) (models.BridgeType, error)
}

type orm struct {
	db *gorm.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

func (o *orm) CreateSpec(taskDAG TaskDAG) (int32, error) {
	var specID int32
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		now := time.Now()

		spec := Spec{
			DotDagSource: taskDAG.DOTSource,
			CreatedAt:    now,
		}
		err := tx.Create(&spec).Error
		if err != nil {
			return err
		}
		specID = spec.ID

		tasks, err := taskDAG.TasksInDependencyOrder()
		if err != nil {
			return err
		}

		taskSpecIDs := make(map[Task]int32)
		for _, task := range tasks {
			var successorID null.Int
			if task.OutputTask() != nil {
				successor := task.OutputTask()
				successorID = null.IntFrom(int64(taskSpecIDs[successor]))
			}

			taskSpec := TaskSpec{
				DotID:          task.DotID(),
				PipelineSpecID: spec.ID,
				Type:           task.Type(),
				JSON:           JSONSerializable{task},
				Index:          task.OutputIndex(),
				SuccessorID:    successorID,
				CreatedAt:      now,
			}
			err = tx.Create(&taskSpec).Error
			if err != nil {
				return err
			}

			taskSpecIDs[task] = taskSpec.ID
		}
		return nil
	})
	return specID, err
}

// CreateRun adds a Run record to the DB, and one TaskRun
// per TaskSpec associated with the given Spec.  Processing of the
// TaskRuns is maximally parallelized across all of the Chainlink nodes in the
// cluster.
func (o *orm) CreateRun(jobID int32) (int64, error) {
	var runID int64
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		// Fetch the spec ID from the job
		var job models.JobSpecV2
		err := o.db.Where("id = ?", jobID).First(&job).Error
		if err != nil {
			return err
		}

		// Ensure the spec exists
		var spec Spec
		err = o.db.Where("id = ?", job.PipelineSpecID).First(&spec).Error
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
                pipeline_run_id, pipeline_task_spec_id, index, dot_id, created_at
            )
            SELECT ? AS pipeline_run_id, id AS pipeline_task_spec_id, index, dot_id, NOW() AS created_at
            FROM pipeline_task_specs
            WHERE pipeline_spec_id = ?
        `, run.ID, run.PipelineSpecID).Error
		return errors.Wrap(err, "could not create pipeline task runs")
	})
	return runID, err
}

// WithNextUnclaimedTaskRun chooses any arbitrary incomplete TaskRun from the DB
// whose parent TaskRuns have already been processed.
func (o *orm) WithNextUnclaimedTaskRun(fn func(ptRun TaskRun, predecessors []TaskRun) Result) (bool, error) {
	var done bool
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		var ptRun TaskRun
		var predecessors []TaskRun

		// Find the next unlocked, unfinished pipeline_task_run with no uncompleted predecessors
		err := tx.Raw(`
            SELECT * from pipeline_task_runs WHERE id IN (
                SELECT pipeline_task_runs.id FROM pipeline_task_runs
                    INNER JOIN pipeline_task_specs ON pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id
                    LEFT JOIN pipeline_task_specs AS predecessor_specs ON predecessor_specs.successor_id = pipeline_task_specs.id
                    LEFT JOIN pipeline_task_runs AS predecessor_unfinished_runs ON predecessor_specs.id = predecessor_unfinished_runs.pipeline_task_spec_id
                          AND pipeline_task_runs.pipeline_run_id = predecessor_unfinished_runs.pipeline_run_id
                WHERE pipeline_task_runs.finished_at IS NULL
                GROUP BY (pipeline_task_runs.id)
                HAVING (
                    bool_and(predecessor_unfinished_runs.finished_at IS NOT NULL)
                    OR
                    count(predecessor_unfinished_runs.id) = 0
                )
            )
            LIMIT 1
            FOR UPDATE OF pipeline_task_runs SKIP LOCKED;
        `).Scan(&ptRun).Error
		if gorm.IsRecordNotFoundError(err) {
			done = true
			return nil
		} else if err != nil {
			return errors.Wrap(err, "error finding next task run")
		}

		// Fill in the TaskSpec, as gorm can't seem to .Preload() it in the query above
		err = tx.Where("id = ?", ptRun.PipelineTaskSpecID).First(&ptRun.PipelineTaskSpec).Error
		if gorm.IsRecordNotFoundError(err) {
			// done = true
			return err
		} else if err != nil {
			return errors.Wrap(err, "error finding next task run's spec")
		}

		// Find all the predecessors
		err = tx.Raw(`
                SELECT pipeline_task_runs.* FROM pipeline_task_runs
                LEFT JOIN pipeline_task_specs AS specs_right ON specs_right.id = pipeline_task_runs.pipeline_task_spec_id
                LEFT JOIN pipeline_task_specs AS specs_left ON specs_left.id = specs_right.successor_id
                LEFT JOIN pipeline_task_runs AS successors ON successors.pipeline_task_spec_id = specs_left.id
                      AND successors.pipeline_run_id = pipeline_task_runs.pipeline_run_id
                WHERE successors.id = ?`, ptRun.ID).Find(&predecessors).Error
		if err != nil {
			return errors.Wrap(err, "error finding task run predecessors")
		}

		// Call the callback and convert its output to a format appropriate for the DB
		result := fn(ptRun, predecessors)

		var out *JSONSerializable
		var errString null.String
		if result.Value != nil {
			out = &JSONSerializable{Val: result.Value}
		} else if result.Error != nil {
			errString = null.StringFrom(err.Error())
		}

		// Update the task run record with the output and error
		err = tx.Exec(`
            UPDATE pipeline_task_runs
            SET output = ?, error = ?, finished_at = ?
            WHERE id = ?
        `, out, errString, time.Now(), ptRun.ID).Error
		if err != nil {
			return errors.Wrap(err, "could not mark pipeline_task_run as finished")
		}
		return nil
	})
	return done, err
}

func (o *orm) AwaitRun(ctx context.Context, runID int64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		done, err := runFinished(o.db, runID)
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

func runFinished(tx *gorm.DB, runID int64) (bool, error) {
	var done struct{ Done bool }
	err := tx.Raw(`
        SELECT bool_and(pipeline_task_runs.error IS NOT NULL OR pipeline_task_runs.output IS NOT NULL) AS done
        FROM pipeline_runs
        JOIN pipeline_task_runs ON pipeline_task_runs.pipeline_run_id = pipeline_runs.id
        WHERE pipeline_runs.id = $1
    `, runID).Scan(&done).Error
	return done.Done, err
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
	var results []Result
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) error {
		done, err := runFinished(tx, runID)
		if err != nil {
			return err
		} else if !done {
			return errors.New("can't fetch run results, run is still in progress")
		}

		var taskRuns []TaskRun
		err = o.db.
			Joins("LEFT JOIN pipeline_task_specs ON pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id").
			Where("pipeline_run_id = ?", runID).
			Where("error IS NOT NULL OR output IS NOT NULL").
			Where("pipeline_task_specs.successor_id IS NULL").
			Order("index ASC").
			Find(&taskRuns).
			Error
		if err != nil {
			return err
		}

		results = make([]Result, len(taskRuns))
		for i, taskRun := range taskRuns {
			results[i] = taskRun.Result()
		}
		return nil
	})
	return results, err
}

func (o *orm) FindBridge(name models.TaskType) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, o.db.First(&bt, "name = ?", name.String()).Error
}
