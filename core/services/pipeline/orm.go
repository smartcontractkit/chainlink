package pipeline

import (
	"context"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
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
	CreateRun(db postgres.Queryer, run *Run) (err error)
	DeleteRun(id int64) error
	StoreRun(db postgres.Queryer, run *Run) (restart bool, err error)
	UpdateTaskRunResult(taskID uuid.UUID, result Result) (run Run, start bool, err error)
	InsertFinishedRun(db postgres.Queryer, run Run, saveSuccessfulTaskRuns bool) (runID int64, err error)
	DeleteRunsOlderThan(context.Context, time.Duration) error
	FindRun(id int64) (Run, error)
	GetAllRuns() ([]Run, error)
	GetUnfinishedRuns(context.Context, time.Time, func(run Run) error) error
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

func (o *orm) CreateRun(db postgres.Queryer, run *Run) (err error) {
	if run.CreatedAt.IsZero() {
		return errors.New("run.CreatedAt must be set")
	}

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = postgres.SqlxTransaction(ctx, db, func(tx *sqlx.Tx) error {
		sql := `INSERT INTO pipeline_runs (pipeline_spec_id, meta, inputs, created_at, state)
		VALUES (:pipeline_spec_id, :meta, :inputs, :created_at, :state)
		RETURNING id`

		query, args, e := tx.BindNamed(sql, run)
		if e != nil {
			return err
		}
		if err = tx.Get(run, query, args...); err != nil {
			return errors.Wrap(err, "error inserting pipeline_run")
		}

		// Now create pipeline_task_runs if any
		if len(run.PipelineTaskRuns) == 0 {
			return nil
		}

		// update the ID key everywhere
		for i := range run.PipelineTaskRuns {
			run.PipelineTaskRuns[i].PipelineRunID = run.ID
		}

		sql = `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at);`
		_, err = tx.NamedExecContext(ctx, sql, run.PipelineTaskRuns)
		return err
	})

	return err
}

// StoreRun will persist a partially executed run before suspending, or finish a run.
// If `restart` is true, then new task run data is available and the run should be resumed immediately.
func (o *orm) StoreRun(tx postgres.Queryer, run *Run) (restart bool, err error) {
	finished := run.FinishedAt.Valid
	if !finished {
		// Lock the current run. This prevents races with /v2/resume
		sql := `SELECT id FROM pipeline_runs WHERE id = $1 FOR UPDATE;`
		if _, err = tx.Exec(sql, run.ID); err != nil {
			return restart, errors.Wrap(err, "StoreRun")
		}

		taskRuns := []TaskRun{}
		// Reload task runs, we want to check for any changes while the run was ongoing
		if err = sqlx.Select(tx, &taskRuns, `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`, run.ID); err != nil {
			return restart, errors.Wrap(err, "StoreRun")
		}

		// Construct a temporary run so we can use r.ByDotID
		tempRun := Run{PipelineTaskRuns: taskRuns}

		// Diff with current state, if updated, swap run.PipelineTaskRuns and early return with restart = true
		for i, tr := range run.PipelineTaskRuns {
			if !tr.IsPending() {
				continue
			}

			// Look for new data
			if taskRun := tempRun.ByDotID(tr.DotID); taskRun != nil && !taskRun.IsPending() {
				// Swap in the latest state
				run.PipelineTaskRuns[i] = *taskRun
				restart = true
			}
		}

		if restart {
			return restart, nil
		}

		// Suspend the run
		run.State = RunStatusSuspended
		if _, err = sqlx.NamedExec(tx, `UPDATE pipeline_runs SET state = :state WHERE id = :id`, run); err != nil {
			return false, errors.Wrap(err, "StoreRun")
		}
	} else {
		// Simply finish the run, no need to do any sort of locking
		if run.Outputs.Val == nil || len(run.FatalErrors) == 0 {
			return false, errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.FatalErrors)
		}
		sql := `UPDATE pipeline_runs SET state = :state, finished_at = :finished_at, all_errors= :all_errors, fatal_errors= :fatal_errors, outputs = :outputs WHERE id = :id`
		if _, err = sqlx.NamedExec(tx, sql, run); err != nil {
			return false, errors.Wrap(err, "StoreRun")
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
	rows, err = sqlx.NamedQuery(tx, sql, run.PipelineTaskRuns)
	if err != nil {
		return false, errors.Wrap(err, "StoreRun")
	}
	taskRuns := []TaskRun{}
	if err = sqlx.StructScan(rows, &taskRuns); err != nil {
		return false, errors.Wrap(err, "StoreRun")
	}
	// replace with new task run data
	run.PipelineTaskRuns = taskRuns
	return false, nil
}

// Used for cleaning up a run that failed and is marked failEarly (should leave no trace of the run)
func (o *orm) DeleteRun(id int64) error {
	db := postgres.UnwrapGormDB(o.db)
	// NOTE: this will cascade and wipe pipeline_task_runs too
	_, err := db.Exec(`DELETE FROM pipeline_runs WHERE id = $1`, id)
	return err
}

func (o *orm) UpdateTaskRunResult(taskID uuid.UUID, result Result) (run Run, start bool, err error) {
	err = postgres.SqlxTransaction(context.Background(), postgres.UnwrapGormDB(o.db), func(tx *sqlx.Tx) error {
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
		sql = `UPDATE pipeline_task_runs SET output = $2, error = $3, finished_at = $4 WHERE id = $1`
		if _, err = tx.Exec(sql, taskID, result.OutputDB(), result.ErrorDB(), time.Now()); err != nil {
			return errors.Wrap(err, "UpdateTaskRunResult")
		}

		if run.State == RunStatusSuspended {
			start = true
			run.State = RunStatusRunning

			// We're going to restart the run, so set it back to "in progress"
			sql = `UPDATE pipeline_runs SET state = $2 WHERE id = $1`
			if _, err = tx.Exec(sql, run.ID, run.State); err != nil {
				return errors.Wrap(err, "UpdateTaskRunResult")
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
func (o *orm) InsertFinishedRun(db postgres.Queryer, run Run, saveSuccessfulTaskRuns bool) (runID int64, err error) {
	if run.CreatedAt.IsZero() {
		return 0, errors.New("run.CreatedAt must be set")
	}
	if run.FinishedAt.IsZero() {
		return 0, errors.New("run.FinishedAt must be set")
	}
	if run.Outputs.Val == nil || len(run.FatalErrors) == 0 {
		return 0, errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.FatalErrors)
	}
	if len(run.PipelineTaskRuns) == 0 && (saveSuccessfulTaskRuns || run.HasErrors()) {
		return 0, errors.New("must provide task run results")
	}

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = postgres.SqlxTransaction(ctx, db, func(tx *sqlx.Tx) error {
		sql := `INSERT INTO pipeline_runs (pipeline_spec_id, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
		VALUES (:pipeline_spec_id, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state)
		RETURNING *;`

		query, args, e := tx.BindNamed(sql, run)
		if e != nil {
			return err
		}

		if err = tx.GetContext(ctx, &run, query, args...); err != nil {
			return errors.Wrap(err, "error inserting finished pipeline_run")
		}

		// update the ID key everywhere
		for i := range run.PipelineTaskRuns {
			run.PipelineTaskRuns[i].PipelineRunID = run.ID
		}

		if !saveSuccessfulTaskRuns && !run.HasErrors() {
			return nil
		}

		sql = `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at);`
		_, err = tx.NamedExecContext(ctx, sql, run.PipelineTaskRuns)
		return err
	})
	return run.ID, err
}

func (o *orm) DeleteRunsOlderThan(ctx context.Context, threshold time.Duration) error {
	return o.db.
		WithContext(ctx).
		Exec(
			`DELETE FROM pipeline_runs WHERE finished_at < ?`, time.Now().Add(-threshold),
		).Error
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

func (o *orm) GetUnfinishedRuns(ctx context.Context, now time.Time, fn func(run Run) error) error {
	return postgres.Batch(func(offset, limit uint) (count uint, err error) {
		var runs []Run

		err = o.db.
			WithContext(ctx).
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

		if ctx.Err() != nil {
			return 0, nil
		} else if err != nil {
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
