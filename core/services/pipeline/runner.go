package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	// Runner checks the DB for incomplete TaskRuns and runs them.  For a
	// TaskRun to be eligible to be run, its parent/input tasks must already
	// all be complete.
	Runner interface {
		Start()
		Stop()
		CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
		AwaitRun(ctx context.Context, runID int64) error
		ResultsForRun(ctx context.Context, runID int64) ([]Result, error)
	}

	runner struct {
		orm                             ORM
		config                          Config
		processIncompleteTaskRunsWorker utils.SleeperTask
		runReaperWorker                 utils.SleeperTask

		utils.StartStopOnce
		chStop chan struct{}
		chDone chan struct{}
	}
)

var (
	promPipelineTaskExecutionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_task_execution_time",
		Help: "How long each pipeline task took to execute",
	},
		[]string{"pipeline_spec_id", "task_type"},
	)
)

func NewRunner(orm ORM, config Config) *runner {
	r := &runner{
		orm:    orm,
		config: config,
		chStop: make(chan struct{}),
		chDone: make(chan struct{}),
	}
	r.processIncompleteTaskRunsWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.processIncompleteTaskRuns),
	)
	r.runReaperWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.runReaper),
	)
	return r
}

func (r *runner) Start() {
	if !r.OkayToStart() {
		logger.Error("Pipeline runner has already been started")
		return
	}
	go r.runLoop()
}

func (r *runner) Stop() {
	if !r.OkayToStop() {
		logger.Error("Pipeline runner has already been stopped")
		return
	}

	close(r.chStop)
	<-r.chDone
}

func (r *runner) destroy() {
	err := r.processIncompleteTaskRunsWorker.Stop()
	if err != nil {
		logger.Error(err)
	}
	err = r.runReaperWorker.Stop()
	if err != nil {
		logger.Error(err)
	}
}

func (r *runner) runLoop() {
	defer close(r.chDone)
	defer r.destroy()

	var newRunEvents <-chan postgres.Event
	newRunsSubscription, err := r.orm.ListenForNewRuns()
	if err != nil {
		logger.Error("Pipeline runner could not subscribe to new run events, falling back to polling")
	} else {
		defer newRunsSubscription.Close()
		newRunEvents = newRunsSubscription.Events()
	}

	dbPollTicker := time.NewTicker(utils.WithJitter(r.config.TriggerFallbackDBPollInterval()))
	defer dbPollTicker.Stop()

	runReaperTicker := time.NewTicker(r.config.JobPipelineReaperInterval())
	defer runReaperTicker.Stop()

	for {
		select {
		case <-r.chStop:
			return
		case <-newRunEvents:
			r.processIncompleteTaskRunsWorker.WakeUp()
		case <-dbPollTicker.C:
			r.processIncompleteTaskRunsWorker.WakeUp()
		case <-runReaperTicker.C:
			r.runReaperWorker.WakeUp()
		}
	}
}

func (r *runner) CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error) {
	runID, err := r.orm.CreateRun(ctx, jobID, meta)
	if err != nil {
		logger.Errorw("Error creating new pipeline run", "jobID", jobID, "error", err)
		return 0, err
	}
	logger.Infow("Pipeline run created", "jobID", jobID, "runID", runID)
	return runID, nil
}

func (r *runner) AwaitRun(ctx context.Context, runID int64) error {
	ctx, cancel := utils.CombinedContext(r.chStop, ctx)
	defer cancel()
	return r.orm.AwaitRun(ctx, runID)
}

func (r *runner) ResultsForRun(ctx context.Context, runID int64) ([]Result, error) {
	ctx, cancel := utils.CombinedContext(r.chStop, ctx)
	defer cancel()
	return r.orm.ResultsForRun(ctx, runID)
}

// NOTE: This could potentially run on a different machine in the cluster than
// the one that originally added the task runs.
func (r *runner) processIncompleteTaskRuns() {
	threads := int(r.config.JobPipelineParallelism())

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-r.chStop:
					return
				default:
				}

				anyRemaining, err := r.processTaskRun()
				if err != nil {
					logger.Errorf("Error processing incomplete task runs: %v", err)
					return
				} else if !anyRemaining {
					return
				}
			}
		}()
	}
	wg.Wait()
}

func (r *runner) processTaskRun() (anyRemaining bool, err error) {
	ctx, cancel := utils.CombinedContext(r.chStop, r.config.JobPipelineMaxTaskDuration())
	defer cancel()

	return r.orm.ProcessNextUnclaimedTaskRun(ctx, func(ctx context.Context, txdb *gorm.DB, jobID int32, taskRun TaskRun, predecessors []TaskRun) Result {
		loggerFields := []interface{}{
			"jobID", jobID,
			"taskName", taskRun.PipelineTaskSpec.DotID,
			"taskID", taskRun.PipelineTaskSpecID,
			"runID", taskRun.PipelineRunID,
			"taskRunID", taskRun.ID,
		}

		start := time.Now()

		logger.Infow("Running pipeline task", loggerFields...)

		inputs := make([]Result, len(predecessors))
		for i, predecessor := range predecessors {
			inputs[i] = predecessor.Result()
		}

		task, err := UnmarshalTaskFromMap(
			taskRun.PipelineTaskSpec.Type,
			taskRun.PipelineTaskSpec.JSON.Val,
			taskRun.PipelineTaskSpec.DotID,
			r.config,
			txdb,
		)
		if err != nil {
			logger.Errorw("Pipeline task run could not be unmarshaled", append(loggerFields, "error", err)...)
			return Result{Error: err}
		}
		var job models.JobSpecV2
		err = txdb.Find(&job, "id = ?", jobID).Error
		if err != nil {
			logger.Errorw("unexpected error could not find job by ID", append(loggerFields, "error", err)...)
			return Result{Error: err}
		}

		// Order of precedence for task timeout:
		// - Specific task timeout (task.TaskTimeout)
		// - Job level task timeout (job.MaxTaskDuration)
		// - Node level task timeout (JobPipelineMaxTaskDuration)
		taskTimeout, isSet := task.TaskTimeout()
		if isSet {
			ctx, cancel = utils.CombinedContext(r.chStop, taskTimeout)
			defer cancel()
		} else if job.MaxTaskDuration != models.Interval(time.Duration(0)) {
			ctx, cancel = utils.CombinedContext(r.chStop, time.Duration(job.MaxTaskDuration))
			defer cancel()
		}

		result := task.Run(ctx, taskRun, inputs)
		if _, is := result.Error.(FinalErrors); !is && result.Error != nil {
			logger.Errorw("Pipeline task run errored", append(loggerFields, "error", result.Error)...)
		} else {
			f := append(loggerFields, "result", result.Value)
			switch v := result.Value.(type) {
			case []byte:
				f = append(f, "resultString", fmt.Sprintf("%q", v))
				f = append(f, "resultHex", fmt.Sprintf("%x", v))
			}
			logger.Infow("Pipeline task completed", f...)
		}

		elapsed := time.Since(start)
		promPipelineTaskExecutionTime.WithLabelValues(string(taskRun.PipelineTaskSpec.PipelineSpecID), string(taskRun.PipelineTaskSpec.Type)).Set(float64(elapsed))

		return result
	})
}

func (r *runner) runReaper() {
	err := r.orm.DeleteRunsOlderThan(r.config.JobPipelineReaperThreshold())
	if err != nil {
		logger.Errorw("Pipeline run reaper failed", "error", err)
	}
}
