package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

var (
	ErrNoSuchBridge = errors.New("no such bridge exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateSpec(ctx context.Context, tx *gorm.DB, taskDAG TaskDAG, maxTaskTimeout models.Interval) (int32, error)
	InsertFinishedRun(db *gorm.DB, run Run, trrs []TaskRunResult, saveSuccessfulTaskRuns bool) (runID int64, err error)
	DeleteRunsOlderThan(threshold time.Duration) error
	FindBridge(name models.TaskType) (models.BridgeType, error)
	FindRun(id int64) (Run, error)
	DB() *gorm.DB
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

// If saveSuccessfulTaskRuns = false, we only save errored runs.
// That way if the job is run frequently (such as OCR) we avoid saving a large number of successful task runs
// which do not provide much value.
func (o *orm) InsertFinishedRun(db *gorm.DB, run Run, trrs []TaskRunResult, saveSuccessfulTaskRuns bool) (runID int64, err error) {
	if run.CreatedAt.IsZero() {
		return 0, errors.New("run.CreatedAt must be set")
	}
	if run.FinishedAt.IsZero() {
		return 0, errors.New("run.FinishedAt must be set")
	}
	if run.Outputs.Val == nil || len(run.Errors) == 0 {
		return 0, errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.Errors)
	}
	if len(trrs) == 0 && saveSuccessfulTaskRuns {
		return 0, errors.New("must provide task run results")
	}

	err = postgres.GormTransactionWithoutContext(db, func(tx *gorm.DB) error {
		if err = tx.Create(&run).Error; err != nil {
			return errors.Wrap(err, "error inserting finished pipeline_run")
		}

		runID = run.ID
		if !saveSuccessfulTaskRuns && !run.HasErrors() {
			return nil
		}

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
		return tx.Exec(stmt, valueArgs...).Error
	})
	return runID, err
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
	var run = Run{ID: id}
	err := o.db.Preload("PipelineTaskRuns").First(&run).Error
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
