package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateSpec(taskDAG TaskDAG) (int32, error)
	CreateRun(jobID int32, meta map[string]interface{}) (int64, error)
	ProcessNextUnclaimedTaskRun(f func(jobID int32, taskRun TaskRun, predecessors []TaskRun) Result) (bool, error)
	ListenForNewRuns() (*utils.PostgresEventListener, error)
	AwaitRun(ctx context.Context, runID int64) error
	ResultsForRun(runID int64) ([]Result, error)

	FindBridge(name models.TaskType) (models.BridgeType, error)
}

type orm struct {
	db  *gorm.DB
	uri string
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, uri string) *orm {
	return &orm{db, uri}
}

func (o *orm) CreateSpec(taskDAG TaskDAG) (int32, error) {
	var specID int32
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) (err error) {
		// Create the pipeline spec
		spec := Spec{
			DotDagSource: taskDAG.DOTSource,
		}
		err = tx.Create(&spec).Error
		if err != nil {
			return err
		}
		specID = spec.ID

		// Create the pipeline task specs in dependency order so
		// that we know what the successor ID for each task is
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
func (o *orm) CreateRun(jobID int32, meta map[string]interface{}) (int64, error) {
	var runID int64
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) (err error) {
		// Create the job run
		run := Run{}

		err = tx.Raw(`
INSERT INTO job_runs (pipeline_spec_id, meta, created_at)
SELECT pipeline_spec_id, $1, NOW()
FROM jobs WHERE id = $2
RETURNING *`, JSONSerializable{Val: meta}, jobID).Scan(&run).Error
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
WHERE pipeline_spec_id = ?`, run.ID, run.PipelineSpecID).Error
		return errors.Wrap(err, "could not create pipeline task runs")
	})
	return runID, err
}

// ProcessNextUnclaimedTaskRun chooses any arbitrary incomplete TaskRun from the DB
// whose parent TaskRuns have already been processed.
func (o *orm) ProcessNextUnclaimedTaskRun(fn func(jobID int32, ptRun TaskRun, predecessors []TaskRun) Result) (done bool, err error) {
	var ptRun TaskRun

	err = utils.GormTransaction(o.db, func(tx *gorm.DB) (err error) {
		var predecessors []TaskRun

		// NOTE: Manual loads below can probably be replaced with Joins in
		// gormv2.
		//
		// Further optimisations (condensing into fewer queries) are
		// probably possible if this turns out to be a hot path

		// Find (and lock) the next unlocked, unfinished pipeline_task_run with no uncompleted predecessors
		err = tx.Raw(`
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

		// Load the TaskSpec
		if err := tx.Where("id = ?", ptRun.PipelineTaskSpecID).First(&ptRun.PipelineTaskSpec).Error; err != nil {
			return errors.Wrap(err, "error finding next task run's spec")
		}

		// Load the PipelineRun
		if err := tx.Where("id = ?", ptRun.PipelineRunID).First(&ptRun.PipelineRun).Error; err != nil {
			return errors.Wrap(err, "error finding next task run's pipeline.Run")
		}

		// Find all the predecessor task runs
		err = tx.Raw(`
                SELECT pipeline_task_runs.* FROM pipeline_task_runs
                INNER JOIN pipeline_task_specs AS specs_right ON specs_right.id = pipeline_task_runs.pipeline_task_spec_id
                LEFT JOIN pipeline_task_specs AS specs_left ON specs_left.id = specs_right.successor_id
                LEFT JOIN pipeline_task_runs AS successors ON successors.pipeline_task_spec_id = specs_left.id
                      AND successors.pipeline_run_id = pipeline_task_runs.pipeline_run_id
                WHERE successors.id = ?`, ptRun.ID).Find(&predecessors).Error
		if err != nil {
			return errors.Wrap(err, "error finding task run predecessors")
		}

		// Get the job ID
		var job struct{ ID int32 }
		err = tx.Raw(`
            SELECT jobs.id FROM pipeline_task_runs
            LEFT JOIN pipeline_task_specs ON pipeline_task_specs.id = pipeline_task_runs.pipeline_task_spec_id
            LEFT JOIN jobs ON jobs.pipeline_spec_id = pipeline_task_specs.pipeline_spec_id
            WHERE pipeline_task_runs.id = ?
        `, ptRun.ID).Scan(&job).Error
		if err != nil {
			return errors.Wrap(err, "error finding job ID")
		}

		// Call the callback
		result := fn(job.ID, ptRun, predecessors)

		// Update the task run record with the output and error
		var out *JSONSerializable
		var errString null.String
		if result.Value != nil {
			out = &JSONSerializable{Val: result.Value}
		} else if result.Error != nil {
			logger.Errorw("Error in pipeline task", "error", result.Error)
			errString = null.StringFrom(result.Error.Error())
		}
		if err := tx.Exec(`UPDATE pipeline_task_runs SET output = ?, error = ?, finished_at = ? WHERE id = ?`, out, errString, time.Now(), ptRun.ID).Error; err != nil {
			return errors.Wrap(err, "could not mark pipeline_task_run as finished")
		}

		// FIXME: Can there be more than one terminal task? If so this won't work
		if ptRun.PipelineTaskSpec.SuccessorID.IsZero() {
			// No more successors means this was the last task in the run
			err = tx.Exec(`
				NOTIFY pipeline_run_completed, ?::text
			`, ptRun.ID).Error
			if err != nil {
				return errors.Wrap(err, "could not notify pipeline_run_completed")
			}
		}
		return nil
	})

	return done, err
}

const (
	postgresChannelRunStarted   = "pipeline_run_started"
	postgresChannelRunCompleted = "pipeline_run_completed"
)

func (o *orm) ListenForNewRuns() (*utils.PostgresEventListener, error) {
	listener := &utils.PostgresEventListener{
		URI:                  o.uri,
		Event:                postgresChannelRunStarted,
		MinReconnectInterval: 1 * time.Second,
		MaxReconnectDuration: 1 * time.Minute,
	}
	err := listener.Start()
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// AwaitRun waits until a run has completed (either successfully or with errors)
// and then returns.  It uses two distinct methods to determine when a job run
// has completed:
//    1) periodic polling
//    2) Postgres notifications
func (o *orm) AwaitRun(ctx context.Context, runID int64) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// This goroutine polls the DB at a set interval
	chPoll := make(chan error)
	go func() {
		var err error
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var done bool
			done, err = runFinished(o.db, runID)
			if err != nil || done {
				break
			}
			time.Sleep(1 * time.Second)
		}

		select {
		case chPoll <- err:
		case <-ctx.Done():
		}
	}()

	// This listener subscribes to the Postgres event informing us of a completed pipeline run
	listener := utils.PostgresEventListener{
		URI:                  o.uri,
		Event:                postgresChannelRunCompleted,
		PayloadFilter:        fmt.Sprintf("%v", runID),
		MinReconnectInterval: 1 * time.Second,
		MaxReconnectDuration: 1 * time.Minute,
	}
	err := listener.Start()
	if err != nil {
		return err
	}
	defer listener.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chPoll:
		return err
	case <-listener.Events():
		return nil
	}
}

func (o *orm) ResultsForRun(runID int64) ([]Result, error) {
	var results []Result
	err := utils.GormTransaction(o.db, func(tx *gorm.DB) (err error) {
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
			Where("finished_at IS NOT NULL").
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

func runFinished(tx *gorm.DB, runID int64) (bool, error) {
	var done struct{ Done bool }
	err := tx.Raw(`
        SELECT bool_and(finished_at) AS done
        FROM pipeline_task_runs
        WHERE pipeline_run_id = $1
    `, runID).Scan(&done).Error
	return done.Done, err
}

func (o *orm) FindBridge(name models.TaskType) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, o.db.First(&bt, "name = ?", name.String()).Error
}
