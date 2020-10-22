package pipeline

import (
	"context"
	"encoding/json"
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
	CreateSpec(ctx context.Context, taskDAG TaskDAG) (int32, error)
	CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
	ProcessNextUnclaimedTaskRun(ctx context.Context, fn ProcessTaskRunFunc) (bool, error)
	ListenForNewRuns() (*utils.PostgresEventListener, error)
	AwaitRun(ctx context.Context, runID int64) error
	RunFinished(runID int64) (bool, error)
	ResultsForRun(ctx context.Context, runID int64) ([]Result, error)
	DeleteRunsOlderThan(threshold time.Duration) error

	FindBridge(name models.TaskType) (models.BridgeType, error)
}

type orm struct {
	db     *gorm.DB
	config Config
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, config Config) *orm {
	return &orm{db, config}
}

func (o *orm) CreateSpec(ctx context.Context, taskDAG TaskDAG) (int32, error) {
	var specID int32

	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	err := utils.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		// Create the pipeline spec
		spec := Spec{
			DotDagSource: taskDAG.DOTSource,
		}
		err := tx.Create(&spec).Error
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

		// Create the final result task that collects the answers from the pipeline's
		// outputs.  This is a Postgres-related performance optimization.
		resultTask := ResultTask{BaseTask{dotID: ResultTaskDotID}}
		for _, task := range tasks {
			if task.DotID() == ResultTaskDotID {
				return errors.Errorf("%v is a reserved keyword and cannot be used in job specs", ResultTaskDotID)
			}
			if task.OutputTask() == nil {
				task.SetOutputTask(&resultTask)
			}
		}
		tasks = append([]Task{&resultTask}, tasks...)

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
	return specID, errors.WithStack(err)
}

// CreateRun adds a Run record to the DB, and one TaskRun
// per TaskSpec associated with the given Spec.  Processing of the
// TaskRuns is maximally parallelized across all of the Chainlink nodes in the
// cluster.
func (o *orm) CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error) {
	var runID int64

	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	err := utils.GormTransaction(ctx, o.db, func(tx *gorm.DB) (err error) {
		// Create the job run
		run := Run{}

		err = tx.Raw(`
            INSERT INTO pipeline_runs (pipeline_spec_id, meta, created_at)
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
            	pipeline_run_id, pipeline_task_spec_id, type, index, created_at
            )
            SELECT ? AS pipeline_run_id, id AS pipeline_task_spec_id, type, index, NOW() AS created_at
            FROM pipeline_task_specs
            WHERE pipeline_spec_id = ?`, run.ID, run.PipelineSpecID).Error
		return errors.Wrap(err, "could not create pipeline task runs")
	})
	return runID, errors.WithStack(err)
}

type ProcessTaskRunFunc func(jobID int32, ptRun TaskRun, predecessors []TaskRun) Result

// ProcessNextUnclaimedTaskRun chooses any arbitrary incomplete TaskRun from the DB
// whose parent TaskRuns have already been processed.
func (o *orm) ProcessNextUnclaimedTaskRun(ctx context.Context, fn ProcessTaskRunFunc) (anyRemaining bool, err error) {
	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		err = o.processNextUnclaimedTaskRun(ctx, fn)
		// "Record not found" errors mean that we're done with all unclaimed
		// task runs.
		if utils.IsRecordNotFound(err) {
			anyRemaining = false
			retry = false
			err = nil
		} else if err != nil {
			retry = true
			err = errors.Wrap(err, "Pipeline runner could not process task run")
			logger.Error(err)

		} else {
			anyRemaining = true
			retry = false
		}
		return
	})
	return anyRemaining, errors.WithStack(err)
}

func (o *orm) processNextUnclaimedTaskRun(ctx context.Context, fn ProcessTaskRunFunc) error {
	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	err := utils.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		var ptRun TaskRun
		var predecessors []TaskRun

		// NOTE: Manual loads below can probably be replaced with Joins in
		// gormv2.
		//
		// Further optimisations (condensing into fewer queries) are
		// probably possible if this turns out to be a hot path

		// Find (and lock) the next unlocked, unfinished pipeline_task_run with no uncompleted predecessors
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
		if err != nil {
			return errors.Wrap(err, "error finding next task run")
		}

		// Load the TaskSpec
		if err = tx.Where("id = ?", ptRun.PipelineTaskSpecID).First(&ptRun.PipelineTaskSpec).Error; err != nil {
			return errors.Wrap(err, "error finding next task run's spec")
		}

		// Load the PipelineRun
		if err = tx.Where("id = ?", ptRun.PipelineRunID).First(&ptRun.PipelineRun).Error; err != nil {
			return errors.Wrap(err, "error finding next task run's pipeline.Run")
		}

		// Find all the predecessor task runs
		err = tx.
			Preload("PipelineTaskSpec").
			Raw(`
                SELECT pipeline_task_runs.* FROM pipeline_task_runs
                INNER JOIN pipeline_task_specs AS specs_right ON specs_right.id = pipeline_task_runs.pipeline_task_spec_id
                LEFT JOIN pipeline_task_specs AS specs_left ON specs_left.id = specs_right.successor_id
                LEFT JOIN pipeline_task_runs AS successors ON successors.pipeline_task_spec_id = specs_left.id
                      AND successors.pipeline_run_id = pipeline_task_runs.pipeline_run_id
                WHERE successors.id = ?
                ORDER BY pipeline_task_runs.index ASC
            `, ptRun.ID).
			Find(&predecessors).Error
		if err != nil {
			return errors.Wrap(err, "error finding task run predecessors")
		}

		// Get the job ID
		var job struct{ ID int32 }
		err = tx.Raw(`
            SELECT jobs.id FROM pipeline_task_runs
            INNER JOIN pipeline_task_specs ON pipeline_task_specs.id = pipeline_task_runs.pipeline_task_spec_id
            INNER JOIN jobs ON jobs.pipeline_spec_id = pipeline_task_specs.pipeline_spec_id
            WHERE pipeline_task_runs.id = ?
    		LIMIT 1
        `, ptRun.ID).Scan(&job).Error
		// TODO: Needs test, what happens if it can't find the job?!
		if err != nil {
			return errors.Wrap(err, "error finding job ID")
		}

		// Call the callback
		result := fn(job.ID, ptRun, predecessors)

		// Update the task run record with the output and error
		var out interface{}
		var errString null.String
		if result.Value != nil {
			out = &JSONSerializable{Val: result.Value}
		}
		if finalErrors, is := result.Error.(FinalErrors); is {
			errString = null.StringFrom(finalErrors.Error())
		} else if result.Error != nil {
			logger.Errorw("Error in pipeline task", "error", result.Error)
			errString = null.StringFrom(result.Error.Error())
		}
		err = tx.Exec(`UPDATE pipeline_task_runs SET output = ?, error = ?, finished_at = ? WHERE id = ?`, out, errString, time.Now(), ptRun.ID).Error
		if err != nil {
			return errors.Wrap(err, "could not mark pipeline_task_run as finished")
		}

		if ptRun.PipelineTaskSpec.IsFinalPipelineOutput() {
			err = tx.Exec(`UPDATE pipeline_runs SET finished_at = ? WHERE id = ?`, time.Now(), ptRun.PipelineTaskSpecID).Error
			if err != nil {
				return errors.Wrap(err, "could not mark pipeline_run as finished")
			}
			// Emit a Postgres notification if this is the final `ResultTask`
			err = tx.Exec(`SELECT pg_notify('pipeline_run_completed', ?::text);`, ptRun.PipelineRunID).Error
			if err != nil {
				return errors.Wrap(err, "could not notify pipeline_run_completed")
			}
			logger.Infow("Pipeline run completed", "runID", ptRun.PipelineRunID)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "while processing task run")
	}

	return nil
}

const (
	postgresChannelRunStarted   = "pipeline_run_started"
	postgresChannelRunCompleted = "pipeline_run_completed"
)

func (o *orm) ListenForNewRuns() (*utils.PostgresEventListener, error) {
	listener := &utils.PostgresEventListener{
		URI:                  o.config.DatabaseURL(),
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
	// This goroutine polls the DB at a set interval
	chPoll := make(chan error)
	chDone := make(chan struct{})
	defer close(chDone)
	go func() {
		var err error
		for {
			select {
			case <-chDone:
				return
			case <-ctx.Done():
				return
			default:
			}

			var done bool
			done, err = o.RunFinished(runID)
			if err != nil || done {
				break
			}
			time.Sleep(1 * time.Second)
		}

		select {
		case chPoll <- err:
		case <-chDone:
		case <-ctx.Done():
		}
	}()

	// This listener subscribes to the Postgres event informing us of a completed pipeline run
	listener := utils.PostgresEventListener{
		URI:                  o.config.DatabaseURL(),
		Event:                postgresChannelRunCompleted,
		PayloadFilter:        fmt.Sprintf("%v", runID),
		MinReconnectInterval: 1 * time.Second,
		MaxReconnectDuration: 1 * time.Minute,
	}
	err := listener.Start()
	if err != nil {
		return err
	}
	defer logger.ErrorIfCalling(listener.Stop)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chPoll:
		return err
	case <-listener.Events():
		return nil
	}
}

func (o *orm) ResultsForRun(ctx context.Context, runID int64) ([]Result, error) {
	// TODO(sam): I think this can be optimised by condensing it down into one query
	// See: https://www.pivotaltracker.com/story/show/175288635
	done, err := o.RunFinished(runID)
	if err != nil {
		return nil, err
	} else if !done {
		return nil, errors.New("can't fetch run results, run is still in progress")
	}

	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	var results []Result
	err = utils.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		var resultTaskRun TaskRun
		err = o.db.
			Preload("PipelineTaskSpec").
			Joins("INNER JOIN pipeline_task_specs ON pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id").
			Where("pipeline_run_id = ?", runID).
			Where("finished_at IS NOT NULL").
			Where("pipeline_task_specs.successor_id IS NULL").
			Where("pipeline_task_specs.dot_id = ?", ResultTaskDotID).
			First(&resultTaskRun).
			Error
		if err != nil {
			return errors.Wrapf(err, "Pipeline runner could not fetch pipeline results (runID: %v)", runID)
		}

		var values []interface{}
		var errs FinalErrors
		if resultTaskRun.Output != nil && resultTaskRun.Output.Val != nil {
			vals, is := resultTaskRun.Output.Val.([]interface{})
			if !is {
				return errors.Errorf("Pipeline runner invariant violation: result task run's output must be []interface{}, got %T", resultTaskRun.Output.Val)
			}
			values = vals
		}
		if !resultTaskRun.Error.IsZero() {
			err = json.Unmarshal([]byte(resultTaskRun.Error.ValueOrZero()), &errs)
			if err != nil {
				return errors.Errorf("Pipeline runner invariant violation: result task run's errors must be []error, got %v", resultTaskRun.Error.ValueOrZero())
			}
		}
		if len(values) != len(errs) {
			return errors.Errorf("Pipeline runner invariant violation: result task run must have equal numbers of outputs and errors (got %v and %v)", len(values), len(errs))
		}
		results = make([]Result, len(values))
		for i := range values {
			results[i].Value = values[i]
			if !errs[i].IsZero() {
				results[i].Error = errors.New(errs[i].ValueOrZero())
			}
		}
		return nil
	})
	return results, err
}

func (o *orm) RunFinished(runID int64) (bool, error) {
	var done struct{ Done bool }
	err := o.db.Raw(`
        SELECT finished_at IS NOT NULL AS done
        FROM pipeline_task_runs
        INNER JOIN pipeline_task_specs ON pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id
        WHERE pipeline_task_runs.pipeline_run_id = ? AND pipeline_task_specs.successor_id IS NULL
		LIMIT 1
    `, runID).Scan(&done).Error
	return done.Done, errors.Wrapf(err, "could not determine if run is finished (run ID: %v)", runID)
}

func (o *orm) DeleteRunsOlderThan(threshold time.Duration) error {
	err := o.db.Exec(`DELETE FROM pipeline_runs WHERE finished_at < ?`, time.Now().Add(-threshold)).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *orm) FindBridge(name models.TaskType) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, o.db.First(&bt, "name = ?", name.String()).Error
}
