package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

var (
	ErrNoSuchBridge = errors.New("no such bridge exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateSpec(ctx context.Context, db *gorm.DB, taskDAG TaskDAG, maxTaskTimeout models.Interval) (int32, error)
	InsertFinishedRunWithResults(ctx context.Context, run Run, trrs []TaskRunResult) (runID int64, err error)
	DeleteRunsOlderThan(threshold time.Duration) error

	FindBridge(name models.TaskType) (models.BridgeType, error)
	FindRun(id int64) (Run, error)
	DB() *gorm.DB

	// Note below methods are not currently used to process runs.
	CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
	AwaitRun(ctx context.Context, runID int64) error
	ProcessNextUnfinishedRun(ctx context.Context, fn ProcessRunFunc) (bool, error)
	ListenForNewRuns() (postgres.Subscription, error)
	RunFinished(runID int64) (bool, error)
	ResultsForRun(ctx context.Context, runID int64) ([]Result, error)
}

type orm struct {
	db               *gorm.DB
	config           Config
	eventBroadcaster postgres.EventBroadcaster
}

var _ ORM = (*orm)(nil)

var (
	promPipelineRunErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pipeline_run_errors",
		Help: "Number of errors for each pipeline spec",
	},
		[]string{"job_id", "job_name"},
	)
	promPipelineRunTotalTimeToCompletion = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_run_total_time_to_completion",
		Help: "How long each pipeline run took to finish (from the moment it was created)",
	},
		[]string{"job_id", "job_name"},
	)
	promPipelineTasksTotalFinished = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pipeline_tasks_total_finished",
		Help: "The total number of pipeline tasks which have finished",
	},
		[]string{"job_id", "job_name", "task_type", "status"},
	)
)

func NewORM(db *gorm.DB, config Config, eventBroadcaster postgres.EventBroadcaster) *orm {
	return &orm{db, config, eventBroadcaster}
}

// The tx argument must be an already started transaction.
func (o *orm) CreateSpec(ctx context.Context, tx *gorm.DB, taskDAG TaskDAG, maxTaskDuration models.Interval) (int32, error) {
	spec := Spec{
		DotDagSource:    taskDAG.DOTSource,
		MaxTaskDuration: maxTaskDuration,
	}
	err := tx.Create(&spec).Error
	if err != nil {
		return 0, err
	}
	return spec.ID, errors.WithStack(err)
}

// CreateRun adds a Run record to the DB, and one TaskRun
// per TaskSpec associated with the given Spec.  Processing of the
// TaskRuns is maximally parallelized across all of the Chainlink nodes in the
// cluster.
func (o *orm) CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (runID int64, err error) {
	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	err = postgres.GormTransaction(ctx, o.db, func(tx *gorm.DB) (err error) {
		// Create the job run
		run := Run{}

		err = tx.Raw(`
            INSERT INTO pipeline_runs (pipeline_spec_id, meta, created_at)
            SELECT pipeline_spec_id, ?, NOW()
            FROM jobs WHERE id = ?
            RETURNING *`, JSONSerializable{Val: meta}, jobID).Scan(&run).Error
		if run.ID == 0 {
			return errors.Errorf("no job found with id %v (most likely it was deleted)", jobID)
		} else if err != nil {
			return errors.Wrap(err, "could not create pipeline run")
		}

		if err = tx.Preload("PipelineSpec").First(&run).Error; err != nil {
			return err
		}
		d := TaskDAG{}
		if err = d.UnmarshalText([]byte(run.PipelineSpec.DotDagSource)); err != nil {
			return err
		}

		var trs []TaskRun
		tasks, err := d.TasksInDependencyOrder()
		if err != nil {
			return err
		}
		for _, ts := range tasks {
			trs = append(trs, TaskRun{
				Type:          ts.Type(),
				PipelineRunID: run.ID,
				Index:         ts.OutputIndex(),
				DotID:         ts.DotID(),
			})
		}
		runID = run.ID
		if len(trs) > 0 {
			return tx.Create(&trs).Error
		}
		return nil
	})
	return runID, errors.WithStack(err)
}

type ProcessRunFunc func(ctx context.Context, txdb *gorm.DB, spec Spec, meta JSONSerializable, l logger.Logger) (TaskRunResults, bool, error)

func (o *orm) ProcessNextUnfinishedRun(ctx context.Context, fn ProcessRunFunc) (bool, error) {
	// Passed in context cancels on (chStop || JobPipelineMaxTaskDuration)
	txContext, cancel := context.WithTimeout(context.Background(), o.config.DatabaseMaximumTxDuration())
	defer cancel()
	var pRun Run

	err := postgres.GormTransaction(txContext, o.db, func(tx *gorm.DB) error {
		err := tx.
			Preload("PipelineSpec").
			Preload("PipelineTaskRuns").
			Where("pipeline_runs.finished_at IS NULL").
			Order("id ASC").
			Clauses(clause.Locking{
				Strength: "UPDATE",
				Options:  "SKIP LOCKED",
			}).
			First(&pRun).Error
		if err != nil {
			return errors.Wrap(err, "error finding unfinished run")
		}
		logger.Infow("Pipeline run started", "runID", pRun.ID)

		trrs, _, err := fn(ctx, tx, pRun.PipelineSpec, pRun.Meta, *logger.Default)
		if err != nil {
			return errors.Wrap(err, "error calling ProcessRunFunc")
		}

		// Populate the task run result IDs by matching the dot
		// IDs.
		for i, trr := range trrs {
			for _, tr := range pRun.PipelineTaskRuns {
				if trr.Task.DotID() == tr.DotID {
					trrs[i].ID = tr.ID
				}
			}
		}

		if err = o.updateTaskRuns(tx, trrs); err != nil {
			return errors.Wrap(err, "could not update task runs")
		}

		if err = o.UpdatePipelineRun(tx, &pRun, trrs.FinalResult()); err != nil {
			return errors.Wrap(err, "could not mark pipeline_run as finished")
		}

		err = o.eventBroadcaster.NotifyInsideGormTx(tx, postgres.ChannelRunCompleted, fmt.Sprintf("%v", pRun.ID))
		if err != nil {
			return errors.Wrap(err, "could not notify pipeline_run_completed")
		}

		elapsed := time.Since(pRun.CreatedAt)
		promPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", pRun.PipelineSpec.JobID), pRun.PipelineSpec.JobName).Set(float64(elapsed))

		if pRun.HasErrors() {
			promPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", pRun.PipelineSpec.JobID), pRun.PipelineSpec.JobName).Inc()
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.Wrap(err, "while processing run")
	}
	logger.Infow("Pipeline run completed", "runID", pRun.ID)
	return true, nil
}

// updateTaskRuns updates multiple task runs in one query
func (o *orm) updateTaskRuns(db *gorm.DB, trrs TaskRunResults) error {
	sql := `
UPDATE pipeline_task_runs AS ptr SET
output = updates.output,
error = updates.error,
finished_at = updates.finished_at
FROM (VALUES
%s
) AS updates(id, output, error, finished_at)
WHERE ptr.id = updates.id
`
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, trr := range trrs {
		valueStrings = append(valueStrings, "(?::bigint, ?::jsonb, ?::text, ?::timestamptz)")
		valueArgs = append(valueArgs, trr.ID, trr.Result.OutputDB(), trr.Result.ErrorDB(), trr.FinishedAt)
	}

	/* #nosec G201 */
	stmt := fmt.Sprintf(sql, strings.Join(valueStrings, ","))
	return db.Exec(stmt, valueArgs...).Error
}

func (o *orm) UpdatePipelineRun(db *gorm.DB, run *Run, result FinalResult) error {
	return db.Raw(`
		UPDATE pipeline_runs SET finished_at = ?, outputs = ?, errors = ?
		WHERE id = ?
		RETURNING *
		`, time.Now(), result.OutputsDB(), result.ErrorsDB(), run.ID).
		Scan(run).Error
}

func (o *orm) ListenForNewRuns() (postgres.Subscription, error) {
	return o.eventBroadcaster.Subscribe(postgres.ChannelRunStarted, "")
}

func (o *orm) InsertFinishedRunWithResults(ctx context.Context, run Run, trrs []TaskRunResult) (runID int64, err error) {
	if run.CreatedAt.IsZero() {
		return 0, errors.New("run.CreatedAt must be set")
	}
	if run.FinishedAt.IsZero() {
		return 0, errors.New("run.FinishedAt must be set")
	}
	if run.Outputs.Val == nil || len(run.Errors) == 0 {
		return 0, errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.Errors)
	}

	err = postgres.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		if err = tx.Create(&run).Error; err != nil {
			return errors.Wrap(err, "error inserting finished pipeline_run")
		}

		runID = run.ID
		sql := `
		INSERT INTO pipeline_task_runs (pipeline_run_id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES %s
		`
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, trr := range trrs {
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?)")
			valueArgs = append(valueArgs, run.ID, trr.Task.Type(), trr.Task.OutputIndex(), trr.Result.OutputDB(), trr.Result.ErrorDB(), trr.Task.DotID(), trr.CreatedAt, trr.FinishedAt)
		}

		/* #nosec G201 */
		stmt := fmt.Sprintf(sql, strings.Join(valueStrings, ","))
		err = tx.Exec(stmt, valueArgs...).Error
		return errors.Wrap(err, "error inserting finished pipeline_task_runs")
	})

	return runID, err
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
	sub, err := o.eventBroadcaster.Subscribe(postgres.ChannelRunCompleted, fmt.Sprintf("%d", runID))
	if err != nil {
		return err
	}
	defer sub.Close()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chPoll:
		return err
	case <-sub.Events():
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
	err = postgres.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		var run Run
		err = tx.Raw(`
			SELECT * FROM pipeline_runs 
				WHERE id = ? 
				AND finished_at IS NOT NULL
`, runID).Scan(&run).
			Error
		if err != nil {
			return errors.Wrapf(err, "Pipeline runner could not fetch pipeline results (runID: %v)", runID)
		}

		var values []interface{}
		if !run.Outputs.Null {
			vals, is := run.Outputs.Val.([]interface{})
			if !is {
				return errors.Errorf("Pipeline runner invariant violation: result task run's output must be []interface{}, got %T", run.Outputs.Val)
			}
			values = vals
		}

		if len(values) != len(run.Errors) {
			return errors.Errorf("Pipeline runner invariant violation: result task run must have equal numbers of outputs and errors (got %v and %v)", len(values), len(run.Errors))
		}
		results = make([]Result, len(values))
		for i := range values {
			results[i].Value = values[i]
			if !run.Errors[i].IsZero() {
				results[i].Error = errors.New(run.Errors[i].String)
			}
		}
		return nil
	})
	return results, err
}

func (o *orm) RunFinished(runID int64) (bool, error) {
	var tr Run
	err := o.db.Where("id = ?", runID).First(&tr).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.Wrapf(err, "could not determine if run is finished (run ID: %v)", runID)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.Errorf("run not found - could not determine if run is finished (run ID: %v)", runID)
	}
	return tr.FinishedAt != nil, nil
}

func (o *orm) DeleteRunsOlderThan(threshold time.Duration) error {
	err := o.db.Exec(`DELETE FROM pipeline_runs WHERE finished_at < ?`, time.Now().Add(-threshold)).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *orm) FindBridge(name models.TaskType) (models.BridgeType, error) {
	return FindBridge(o.db, name)
}

func (o *orm) FindRun(id int64) (Run, error) {
	var run Run
	err := o.db.Raw(`
        SELECT * 
        FROM pipeline_runs
		WHERE id = ?
    `, id).Scan(&run).Error

	return run, err
}

// FindBridge find a bridge using the given database
func FindBridge(db *gorm.DB, name models.TaskType) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, errors.Wrapf(db.First(&bt, "name = ?", name.String()).Error, "could not find bridge with name '%s'", name)
}

func (o *orm) DB() *gorm.DB {
	return o.db
}
