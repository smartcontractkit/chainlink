package job

import (
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	// "github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

type RunnerORM interface {
	JobSpec(id models.ID) (JobSpec, error)
	CreateJobRun(jobRun *JobRun) error
	LockFirstIncompleteTaskRunWithCompletedParents() (TaskRun, func(), error)
	MarkTaskRunCompleted(taskRunID uint64, output interface{}, err error) error
	OutputTaskRunsForTaskRun(taskRunID uint64) ([]TaskRun, error)
	AllTaskRunsCompleted(jobRunID uint64) (bool, error)
	MarkJobRunCompleted(jobRunID uint64)
}

type JobRun struct {
	ID          uint64 `gorm:"primary_key;auto_increment;not null"`
	JobSpecID   *models.ID
	JobSpecType string

	TaskRuns []TaskRun
}

type TaskRun struct {
	ID       uint64 `gorm:"primary_key;auto_increment;not null"`
	JobRunID uint64

	Output    *JSONSerializable `gorm:"type:jsonb"`
	Error     string
	Completed bool `gorm:"not null;default:false"`

	Task          *TaskDBRow `json:"-"`
	InputTaskRuns []TaskRun  `json:"-" gorm:"-"`
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

func (r *runner) CreateJobRun(id models.ID) error {
	jobSpec, err := r.orm.JobSpec(id)
	if err != nil {
		return err
	}

	jobRun := &JobRun{
		JobSpecID:   &id,
		JobSpecType: string(jobSpec.JobType()),
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

func (r *runner) processIncompleteTaskRuns() {
	for {
		// SELECT * FROM task_runs t
		// LEFT JOIN task_runs_join_table AS join ON t.id = join.child_id
		// LEFT JOIN task_runs AS parent ON join.parent_id = parent.id
		// WHERE t.id = ? AND parent.completed = true
		taskRun, unlock, err := r.orm.LockFirstIncompleteTaskRunWithCompletedParents()
		// if errors.Cause(err) == Err404 {
		// 	// All task runs complete
		// 	break
		// } else
		if err != nil {
			logger.Errorf("error fetching task runs: %v", err)
			return
		}
		defer unlock()

		inputs := make([]Result, len(taskRun.InputTaskRuns))
		for i, parent := range taskRun.InputTaskRuns {
			inputs[i] = Result{
				Value: parent.Output.Value,
				Error: errors.New(parent.Error),
			}
		}

		output, err := taskRun.Task.Task().Run(inputs)
		if err != nil {
			logger.Errorf("error in task run %v:", err)
		}

		r.orm.MarkTaskRunCompleted(taskRun.ID, output, err)

		// If this task has no children, it's an output task.
		// If it's an output task, it might be the last remaining output task.
		// If there's a chance that it's the last output task, we need to check
		//     for job run completion.
		// If we find that the job run has completed, we need to update the
		//     job run in the DB to reflect this.
		outputTaskRuns, err := r.orm.OutputTaskRunsForTaskRun(taskRun.ID)
		if err != nil {
			logger.Errorw("error fetching output task runs",
				"jobRunID", taskRun.JobRunID,
				"taskRunID", taskRun.ID,
				"error", err,
			)
			return
		}
		if len(outputTaskRuns) == 0 {
			// SELECT j.num_outputs as num_outputs, COUNT(t.id) as num_finished_task_runs
			//   FROM job_runs AS j
			//   LEFT JOIN task_runs AS t ON j.id = t.job_run_id
			//   WHERE j.id = ? AND t.completed = true
			jobRunCompleted, err := r.orm.AllTaskRunsCompleted(taskRun.JobRunID)
			if err != nil {
				logger.Errorw("error checking job completion status",
					"taskRunID", taskRun.ID,
					"error", err,
				)
				return
			}

			if jobRunCompleted {
				r.orm.MarkJobRunCompleted(taskRun.JobRunID)
			}
		}
	}
}
