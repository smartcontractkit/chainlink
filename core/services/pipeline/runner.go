package pipeline

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	// "github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

type Runner interface {
	Start()
	Stop()
	CreatePipelineRun(pipelineSpecID models.Sha256Hash) error
}

type runner struct {
	// processTasks services.SleeperTask
	orm    RunnerORM
	chStop chan struct{}
	chDone chan struct{}
}

// FIXME: This interface probably needs rethinking
type RunnerORM interface {
	LoadPipelineSpec(id int64) (PipelineSpec, error)
	CreatePipelineRun(id int64) error
	NextTaskRunForExecution(func(*gorm.DB, PipelineTaskRun) error) error
}

func NewRunner(orm RunnerORM) *runner {
	r := &runner{
		orm:    orm,
		chStop: make(chan struct{}),
		chDone: make(chan struct{}),
	}
	// r.processTasks = services.NewSleeperTask(services.SleeperTaskFuncWorker(r.processIncompleteTaskRuns))
	return r
}

func (r *runner) Start() {
	go func() {
		defer close(r.chDone)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.chStop:
				return
			case <-ticker.C:
				r.processIncompleteTaskRuns()
			}
		}
	}()
}

func (r *runner) Stop() {
	close(r.chStop)
	<-r.chDone
}

func (r *runner) CreatePipelineRun(pipelineSpecID models.Sha256Hash) error {
	spec, err := r.orm.LoadPipelineSpec(pipelineSpecID)
	if err != nil {
		return err
	}

	run := &PipelineRun{
		PipelineSpecID: &id,
	}

	err = r.orm.CreatePipelineRun(run)
	if err != nil {
		return err
	}

	// r.processTasks.WakeUp()
	return nil
}

type Result struct {
	Value interface{}
	Error error
}

// NOTE: This could potentially run on another machine in the cluster
func (r *runner) processIncompleteTaskRuns() error {
	for {
		var jobRunID int64
		r.orm.NextTaskRunForExecution(func(ptRun PipelineTaskRun, predecessors []PipelineTaskRun) error {
			jobRunID = ptRun.PipelineRunID

			inputs := make([]Result, len(predecessors))
			for i, predecessor := range predecessors {
				inputs[i] = Result{
					Value: predecessor.Output.Value,
					Error: predecessor.ResultError(),
				}
			}

			output, err := ptRun.PipelineTaskSpec.Task.Run(inputs)
			if err != nil {
				logger.Errorf("error in task run %v:", err)
			}

			err = r.finishTaskRun(tx, taskRun.ID, output, err)
			if err != nil {
				return errors.Wrap(err, "could not mark task run completed")
			}
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// All task runs complete
			break
		} else if err != nil {
			return err
		}

		// XXX: An example of a completion notification could be as follows
		err = r.orm.NotifyCompletion(jobRunID)

		if err != nil {
			return err
		}
	}
}

func (r *runner) finishTaskRun(tx *gorm.DB, ptRunID int64, output *JSONSerializable, resultErr error) {
	timeNow := time.Now()
	err := tx.Exec(`UPDATE pipeline_task_runs SET output = ?, err = ?, finished_at = ? WHERE id = ?`, output, resultErr, timeNow, ptRunID).Error
	return errors.Wrap(err, "could not mark pipeline_task_run as finished")
}

type runnerORM struct {
	orm *orm.ORM
}

func (r *runnerORM) LoadSpec(id models.ID) (j Spec, err error) {
	err = r.orm.First(&j, id)
	return j, err
}

const pipelineTaskRunInsertSql = `
INSERT INTO pipeline_task_runs (
	pipeline_run_id, pipeline_task_spec_id, created_at
)
SELECT ? AS pipeline_run_id, id AS pipeline_task_spec_id, NOW() AS created_at
FROM pipeline_task_specs
WHERE pipeline_spec_id = ?
`

func (r *runnerORM) CreatePipelineRun(prun PipelineRun) error {
	return r.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&prun).Error
		if err != nil {
			return errors.Wrap(err, "could not create pipeline run")
		}

		err = tx.Exec(pipelineTaskRunInsertSql, prun.ID, prun.PipelineSpecID).Error
		return errors.Wrap(err, "could not create pipeline task runs")
	})
}

func (r *runnerORM) NextTaskRunForExecution(f func(tx *gorm.DB, ptRun PipelineTaskRun, predecessors []PipelineTaskRun) error) error {
	return r.orm.Transaction(func(tx *gorm.DB) error {
		var ptRun PipelineTaskRun
		var predecessors []PipelineTaskRun

		// NOTE: This could conceivably be done in pure SQL and made marginally more efficient by preloading with joins

		// Find the next unlocked, unfinished pipeline_task_run with no uncompleted predecessors
		err := tx.Table("pipeline_task_runs AS successor").
			Set("gorm:query_option", "FOR UPDATE OF successor SKIP LOCKED").
			Joins("LEFT JOIN pipeline_task_runs AS predecessors ON successor.id = predecessor.successor_id AND predecessors.finished_at IS NULL").
			Where("predecessors.id IS NULL").
			Where("successor.finished_at IS NULL").
			Preload("PipelineTaskSpec").
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

		return f(tx, ptRun, predecessors)
	})
}

func (r *runnerORM) MarkTaskRunCompleted(tx *gorm.DB, taskRunID uint64, output sql.Scanner, err error) error {
	return tx.Exec(`UPDATE q_task_runs SET completed = true, output = ?, error = ? WHERE id = ?`, output, err, taskRunID).Error
}

func (r *runnerORM) NotifyCompletion(jobRunID int64) error {
	return r.orm.DB.Exec(`
	$$
	BEGIN
		IF (SELECT bool_and(q_task_runs.error IS NOT NULL OR q_task_runs.output IS NOT NULL) FROM q_job_runs JOIN q_task_runs ON q_task_runs.q_job_run_id = q_job_runs.id WHERE q_job_runs.id = $1)
			PERFORM pg_notify('q_job_run_completed', $1::text);
		END IF;
	END;
	$$ LANGUAGE plpgsql;
	)`, jobRunID).Error
}
