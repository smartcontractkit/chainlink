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
	CreateJobRun(id models.ID) error
}

type runner struct {
	// processTasks services.SleeperTask
	orm    RunnerORM
	chStop chan struct{}
	chDone chan struct{}
}

// FIXME: This interface probably needs rethinking
type RunnerORM interface {
	LoadSpec(id int64) (Spec, error)
	CreatePipelineRun(id int64) error
	NextTaskRunForExecution(func(*gorm.DB, TaskRun) error) error
	MarkTaskRunCompleted(tx *gorm.DB, taskRunID int64, output sql.Scanner, err error) error
}

type PipelineSpec struct {
	ID               int64
	PipelineSpecType string
}

type PipelineRun struct {
	ID           int64 `gorm:"primary_key;auto_increment;not null"`
	PipelineSpec PipelineSpec

	TaskRuns []TaskRun
}

type TaskRun struct {
	ID       int64 `gorm:"primary_key;auto_increment;not null"`
	JobRunID int64

	Output *JSONSerializable `gorm:"type:jsonb"`
	Error  string

	Task          Task      `json:"-"`
	InputTaskRuns []TaskRun `json:"-" gorm:"many2many:q_task_run_edges;joinForeignKey:q_child_id;JoinReferences:q_parent_id"`
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

func (r *runner) CreatePipelineRun(id int64) error {
	spec, err := r.orm.LoadSpec(id)
	if err != nil {
		return err
	}

	run := &PipelineRun{
		PipelineSpecID:   &id,
		PipelineSpecType: string(spec.Type()),
	}

	// for _, task := range jobSpec.Tasks() {
	// 	jobRun.TaskRuns = append(jobRun.TaskRuns, TaskRun{
	// 		TaskID:   task.TaskID(),
	// 		TaskType: task.TaskType(),
	// 	})
	// }

	err = r.orm.CreateJobRun(jobRun)
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
		r.orm.NextTaskRunForExecution(func(taskRun TaskRun) error {
			jobRunID = taskRun.JobRunID

			inputs := make([]Result, len(taskRun.InputTaskRuns))
			for i, parent := range taskRun.InputTaskRuns {
				inputs[i] = Result{
					Value: parent.Output.Value,
					Error: parent.Error,
				}
			}

			output, err := taskRun.Task.Run(inputs)
			if err != nil {
				logger.Errorf("error in task run %v:", err)
			}

			err = r.orm.MarkTaskRunCompleted(tx, taskRun.ID, output, err)
			if err != nil {
				return errors.Wrap(err, "could not mark task run completed")
			}
		})
		if errors.Cause(err) == gorm.ErrRecordNotFound {
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

type runnerORM struct {
	orm *orm.ORM
}

func (r *runnerORM) LoadSpec(id models.ID) (j Spec, err error) {
	err = r.orm.First(&j, id)
	return j, err
}

func (r *runnerORM) CreatePipelineRun(pr *PipelineRun) error {
	return r.orm.Create(pr)
}

func (r *runnerORM) NextTaskRunForExecution(f func(*gorm.DB, TaskRun) error) error {
	return r.orm.Transaction(func(tx *gorm.DB) error {
		var taskRun TaskRun
		// NOTE: Convert to join preloads with gormv2
		err := tx.Table("q_task_runs AS child").
			Set("gorm:query_option", "FOR UPDATE OF child SKIP LOCKED").
			Joins("LEFT JOIN q_task_run_edges AS edge ON child.id = edge.child_id").
			Joins("LEFT JOIN q_task_runs AS parent ON parent.id = edge.parent_id").
			Where("parent.id IS NULL OR parent.completed = true").
			Preload("Task").
			Preload("InputTaskRuns").
			Order("id ASC").
			First(&taskRun)
		if err != nil {
			return errors.Wrap(err, "error finding next task run")
		}
		return f(tx, taskRun)
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
