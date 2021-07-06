package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

var (
	ErrNoSuchBridge = errors.New("no such bridge exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateSpec(ctx context.Context, tx *gorm.DB, pipeline Pipeline, maxTaskTimeout models.Interval) (int32, error)
	InsertFinishedRun(db *gorm.DB, run Run, trrs []TaskRunResult, saveSuccessfulTaskRuns bool) (runID int64, err error)
	DeleteRunsOlderThan(threshold time.Duration) error
	FindRun(id int64) (Run, error)
	GetAllRuns() ([]Run, error)
	DB() *gorm.DB
}

type orm struct {
	db *gorm.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

// The tx argument must be an already started transaction.
func (o *orm) CreateSpec(ctx context.Context, tx *gorm.DB, pipeline Pipeline, maxTaskDuration models.Interval) (int32, error) {
	spec := Spec{
		DotDagSource:    pipeline.Source,
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
	if len(trrs) == 0 && (saveSuccessfulTaskRuns || run.HasErrors()) {
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
	err := o.db.Exec(
		`DELETE FROM pipeline_runs WHERE finished_at < ?`, time.Now().Add(-threshold),
	).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *orm) FindRun(id int64) (Run, error) {
	var run = Run{ID: id}
	err := o.db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		}).First(&run).Error
	return run, err
}

func (o *orm) GetAllRuns() ([]Run, error) {
	var runs []Run
	err := o.db.
		Preload("PipelineSpec").
		Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at ASC, id ASC")
		}).Find(&runs).Error
	return runs, err
}

func (o *orm) DB() *gorm.DB {
	return o.db
}
