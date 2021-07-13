package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var (
	ErrNoSuchBridge = errors.New("no such bridge exists")
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateSpec(ctx context.Context, tx *gorm.DB, pipeline Pipeline, maxTaskTimeout models.Interval) (int32, error)
	CreateRun(db *gorm.DB, run *Run) (err error)
	StoreRun(db *sql.DB, run *Run) (restart bool, err error)
	UpdateTaskRunResult(db *sql.DB, taskID uuid.UUID, result interface{}) (run Run, start bool, err error)
	InsertFinishedRun(db *gorm.DB, run Run, trrs []TaskRunResult, saveSuccessfulTaskRuns bool) (runID int64, err error)
	DeleteRunsOlderThan(threshold time.Duration) error
	FindRun(id int64) (Run, error)
	GetAllRuns() ([]Run, error)
	GetUnfinishedRuns(now time.Time, fn func(run Run) error) error
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

func (o *orm) CreateRun(db *gorm.DB, run *Run) (err error) {
	if run.CreatedAt.IsZero() {
		return errors.New("run.CreatedAt must be set")
	}
	if err = db.Create(run).Error; err != nil {
		return errors.Wrap(err, "error inserting pipeline_run")
	}
	return err
}

// StoreRun will persist a partially executed run before suspending, or finish a run.
// If `restart` is true, then new task run data is available and the run should be resumed immediately.
func (o *orm) StoreRun(db *sql.DB, run *Run) (restart bool, err error) {
	finished := run.FinishedAt.Valid
	err = postgres.SqlxTransaction(context.Background(), db, func(tx *sqlx.Tx) error {
		if !finished {
			// Lock the current run. This prevents races with /v2/resume
			sql := `SELECT id FROM pipeline_runs WHERE id = $1 FOR UPDATE;`
			if _, err = tx.Exec(sql, run.ID); err != nil {
				return err
			}

			taskRuns := []TaskRun{}
			// Reload task runs, we want to check for any changes while the run was ongoing
			if err = tx.Select(&taskRuns, `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`, run.ID); err != nil {
				return err
			}

			// Construct a temporary run so we can use r.ByDotID
			tempRun := Run{PipelineTaskRuns: taskRuns}

			// Diff with current state, if updated, swap run.PipelineTaskRuns and early return with restart = true
			for i, tr := range run.PipelineTaskRuns {
				if !tr.IsPending() {
					continue
				}

				// Look for new data
				if taskRun := tempRun.ByDotID(tr.DotID); taskRun != nil {
					// Swap in the latest state
					run.PipelineTaskRuns[i] = *taskRun
					restart = true
				}
			}

			if restart {
				return nil
			}

			// Suspend the run
			run.State = RunStatusSuspended
			if _, err = tx.NamedExec(`UPDATE pipeline_runs SET state = :state`, run); err != nil {
				return err
			}
		} else {
			// Simply finish the run, no need to do any sort of locking
			if run.Outputs.Val == nil || len(run.Errors) == 0 {
				return errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.Errors)
			}
			sql := `UPDATE pipeline_runs SET state = :state, finished_at = :finished_at, errors= :errors, outputs = :outputs WHERE id = :id`
			if _, err = tx.NamedExec(sql, run); err != nil {
				return err
			}
		}

		sql := `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at)
		ON CONFLICT (pipeline_run_id, dot_id) DO UPDATE SET
		output = EXCLUDED.output, error = EXCLUDED.error, finished_at = EXCLUDED.finished_at
		RETURNING *;
		`

		// NOTE: can't use Select() to auto scan because we're using NamedQuery,
		// sqlx.Named + Select is possible but it's about the same amount of code
		var rows *sqlx.Rows
		rows, err = tx.NamedQuery(sql, run.PipelineTaskRuns)
		if err != nil {
			return err
		}
		taskRuns := []TaskRun{}
		if err = sqlx.StructScan(rows, &taskRuns); err != nil {
			return err
		}
		// replace with new task run data
		run.PipelineTaskRuns = taskRuns
		return nil
	})
	return restart, err
}

func (o *orm) UpdateTaskRunResult(db *sql.DB, taskID uuid.UUID, result interface{}) (run Run, start bool, err error) {
	err = postgres.SqlxTransaction(context.Background(), db, func(tx *sqlx.Tx) error {
		sql := `
		SELECT pipeline_runs.*, pipeline_specs.dot_dag_source "pipeline_spec.dot_dag_source"
		FROM pipeline_runs
		JOIN pipeline_task_runs ON (pipeline_task_runs.pipeline_run_id = pipeline_runs.id)
		JOIN pipeline_specs ON (pipeline_specs.id = pipeline_runs.pipeline_spec_id)
		WHERE pipeline_task_runs.id = $1 AND pipeline_runs.state in ('running', 'suspended')
		FOR UPDATE`
		if err = tx.Get(&run, sql, taskID); err != nil {
			return err
		}

		// Update the task with result
		sql = `UPDATE pipeline_task_runs SET output = $2, finished_at = $3 WHERE id = $1`
		if _, err = tx.Exec(sql, taskID, JSONSerializable{Val: result}, time.Now()); err != nil {
			return err
		}

		if run.State == RunStatusSuspended {
			start = true
			run.State = RunStatusRunning

			// We're going to restart the run, so set it back to "in progress"
			sql = `UPDATE pipeline_runs SET state = $2 WHERE id = $1`
			if _, err = tx.Exec(sql, run.ID, run.State); err != nil {
				return err
			}

			// NOTE: can't join and preload in a single query unless explicitly listing all the struct fields...
			// https://snippets.aktagon.com/snippets/757-how-to-join-two-tables-with-jmoiron-sqlx
			sql = `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`
			return tx.Select(&run.PipelineTaskRuns, sql, run.ID)
		}

		return nil
	})
	return run, start, err
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

		if !saveSuccessfulTaskRuns && !run.HasErrors() {
			return nil
		}

		sql := `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES %s
		`
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, trr := range trrs {
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?)")
			valueArgs = append(valueArgs, run.ID, trr.ID, trr.Task.Type(), trr.Task.OutputIndex(), trr.Result.OutputDB(), trr.Result.ErrorDB(), trr.Task.DotID(), trr.CreatedAt, trr.FinishedAt)
		}

		/* #nosec G201 */
		stmt := fmt.Sprintf(sql, strings.Join(valueStrings, ","))
		return tx.Exec(stmt, valueArgs...).Error
	})
	return run.ID, err
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

func (o *orm) GetUnfinishedRuns(now time.Time, fn func(run Run) error) error {
	return postgres.Batch(func(offset, limit uint) (count uint, err error) {
		var runs []Run

		err = o.db.
			Preload("PipelineSpec").
			Preload("PipelineTaskRuns", func(db *gorm.DB) *gorm.DB {
				return db.
					Order("created_at ASC, id ASC")
			}).
			Where(`state = ? AND created_at < ?`, RunStatusRunning, now).
			Order("created_at ASC, id ASC").
			Limit(int(limit)).
			Offset(int(offset)).
			Find(&runs).Error

		if err != nil {
			return 0, err
		}

		for _, run := range runs {
			if err := fn(run); err != nil {
				return 0, err
			}
		}

		return uint(len(runs)), nil
	})
}

func (o *orm) DB() *gorm.DB {
	return o.db
}
