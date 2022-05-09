package pipeline

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateSpec(pipeline Pipeline, maxTaskTimeout models.Interval, qopts ...pg.QOpt) (int32, error)
	CreateRun(run *Run, qopts ...pg.QOpt) (err error)
	InsertRun(run *Run, qopts ...pg.QOpt) error
	DeleteRun(id int64) error
	StoreRun(run *Run, qopts ...pg.QOpt) (restart bool, err error)
	UpdateTaskRunResult(taskID uuid.UUID, result Result) (run Run, start bool, err error)
	InsertFinishedRun(run *Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) (err error)

	// InsertFinishedRuns inserts all the given runs into the database.
	// If saveSuccessfulTaskRuns is false, only errored runs are saved.
	InsertFinishedRuns(run []*Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) (err error)

	DeleteRunsOlderThan(context.Context, time.Duration) error
	FindRun(id int64) (Run, error)
	GetAllRuns() ([]Run, error)
	GetUnfinishedRuns(context.Context, time.Time, func(run Run) error) error
	GetQ() pg.Q
}

type orm struct {
	q    pg.Q
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *orm {
	return &orm{pg.NewQ(db, lggr, cfg), lggr}
}

func (o *orm) CreateSpec(pipeline Pipeline, maxTaskDuration models.Interval, qopts ...pg.QOpt) (id int32, err error) {
	q := o.q.WithOpts(qopts...)
	sql := `INSERT INTO pipeline_specs (dot_dag_source, max_task_duration, created_at)
	VALUES ($1, $2, NOW())
	RETURNING id;`
	err = q.Get(&id, sql, pipeline.Source, maxTaskDuration)
	return id, errors.WithStack(err)
}

func (o *orm) CreateRun(run *Run, qopts ...pg.QOpt) (err error) {
	if run.CreatedAt.IsZero() {
		return errors.New("run.CreatedAt must be set")
	}

	q := o.q.WithOpts(qopts...)
	err = q.Transaction(func(tx pg.Queryer) error {
		if e := o.InsertRun(run, pg.WithQueryer(tx)); e != nil {
			return errors.Wrap(e, "error inserting pipeline_run")
		}

		// Now create pipeline_task_runs if any
		if len(run.PipelineTaskRuns) == 0 {
			return nil
		}

		// update the ID key everywhere
		for i := range run.PipelineTaskRuns {
			run.PipelineTaskRuns[i].PipelineRunID = run.ID
		}

		sql := `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at);`
		_, err = tx.NamedExec(sql, run.PipelineTaskRuns)
		return err
	})

	return errors.Wrap(err, "CreateRun failed")
}

// InsertRun inserts a run into the database
func (o *orm) InsertRun(run *Run, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	sql := `INSERT INTO pipeline_runs (pipeline_spec_id, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
		VALUES (:pipeline_spec_id, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state)
		RETURNING *;`
	return q.GetNamed(sql, run, run)
}

// StoreRun will persist a partially executed run before suspending, or finish a run.
// If `restart` is true, then new task run data is available and the run should be resumed immediately.
func (o *orm) StoreRun(run *Run, qopts ...pg.QOpt) (restart bool, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.Transaction(func(tx pg.Queryer) error {
		finished := run.FinishedAt.Valid
		if !finished {
			// Lock the current run. This prevents races with /v2/resume
			sql := `SELECT id FROM pipeline_runs WHERE id = $1 FOR UPDATE;`
			if _, err = tx.Exec(sql, run.ID); err != nil {
				return errors.Wrap(err, "StoreRun")
			}

			taskRuns := []TaskRun{}
			// Reload task runs, we want to check for any changes while the run was ongoing
			if err = sqlx.Select(tx, &taskRuns, `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`, run.ID); err != nil {
				return errors.Wrap(err, "StoreRun")
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
				return nil
			}

			// Suspend the run
			run.State = RunStatusSuspended
			if _, err = sqlx.NamedExec(tx, `UPDATE pipeline_runs SET state = :state WHERE id = :id`, run); err != nil {
				return errors.Wrap(err, "StoreRun")
			}
		} else {
			// Simply finish the run, no need to do any sort of locking
			if run.Outputs.Val == nil || len(run.FatalErrors) == 0 {
				return errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.FatalErrors)
			}
			sql := `UPDATE pipeline_runs SET state = :state, finished_at = :finished_at, all_errors= :all_errors, fatal_errors= :fatal_errors, outputs = :outputs WHERE id = :id`
			if _, err = sqlx.NamedExec(tx, sql, run); err != nil {
				return errors.Wrap(err, "StoreRun")
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
			return errors.Wrap(err, "StoreRun")
		}
		taskRuns := []TaskRun{}
		if err = sqlx.StructScan(rows, &taskRuns); err != nil {
			return errors.Wrap(err, "StoreRun")
		}
		// replace with new task run data
		run.PipelineTaskRuns = taskRuns
		return nil
	})
	return
}

// DeleteRun cleans up a run that failed and is marked failEarly (should leave no trace of the run)
func (o *orm) DeleteRun(id int64) error {
	// NOTE: this will cascade and wipe pipeline_task_runs too
	_, err := o.q.Exec(`DELETE FROM pipeline_runs WHERE id = $1`, id)
	return err
}

func (o *orm) UpdateTaskRunResult(taskID uuid.UUID, result Result) (run Run, start bool, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
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

// InsertFinishedRuns inserts all the given runs into the database.
func (o *orm) InsertFinishedRuns(runs []*Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.Transaction(func(tx pg.Queryer) error {
		pipelineRunsQuery := `
INSERT INTO pipeline_runs 
	(pipeline_spec_id, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
VALUES 
	(:pipeline_spec_id, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state) 
RETURNING id
	`
		rows, errQ := tx.NamedQuery(pipelineRunsQuery, runs)
		if errQ != nil {
			return errors.Wrap(errQ, "inserting finished pipeline runs")
		}

		var runIDs []int64
		for rows.Next() {
			var runID int64
			if errS := rows.Scan(&runID); errS != nil {
				return errors.Wrap(errS, "scanning pipeline runs id row")
			}
			runIDs = append(runIDs, runID)
		}

		for i, run := range runs {
			for j := range run.PipelineTaskRuns {
				run.PipelineTaskRuns[j].PipelineRunID = runIDs[i]
			}
		}

		pipelineTaskRunsQuery := `
INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at);
	`
		var pipelineTaskRuns []TaskRun
		for _, run := range runs {
			if !saveSuccessfulTaskRuns && !run.HasErrors() {
				continue
			}
			pipelineTaskRuns = append(pipelineTaskRuns, run.PipelineTaskRuns...)
		}

		_, errE := tx.NamedExec(pipelineTaskRunsQuery, pipelineTaskRuns)
		return errors.Wrap(errE, "insert pipeline task runs")
	})
	return errors.Wrap(err, "InsertFinishedRuns failed")
}

func (o *orm) checkFinishedRun(run *Run, saveSuccessfulTaskRuns bool) error {
	if run.CreatedAt.IsZero() {
		return errors.New("run.CreatedAt must be set")
	}
	if run.FinishedAt.IsZero() {
		return errors.New("run.FinishedAt must be set")
	}
	if run.Outputs.Val == nil || len(run.FatalErrors) == 0 {
		return errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, Errors: %#v", run.Outputs.Val, run.FatalErrors)
	}
	if len(run.PipelineTaskRuns) == 0 && (saveSuccessfulTaskRuns || run.HasErrors()) {
		return errors.New("must provide task run results")
	}
	return nil
}

// InsertFinishedRun inserts the given run into the database.
// If saveSuccessfulTaskRuns = false, we only save errored runs.
// That way if the job is run frequently (such as OCR) we avoid saving a large number of successful task runs
// which do not provide much value.
func (o *orm) InsertFinishedRun(run *Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) (err error) {
	if err = o.checkFinishedRun(run, saveSuccessfulTaskRuns); err != nil {
		return err
	}

	q := o.q.WithOpts(qopts...)
	err = q.Transaction(func(tx pg.Queryer) error {
		sql := `INSERT INTO pipeline_runs (pipeline_spec_id, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
		VALUES (:pipeline_spec_id, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state)
		RETURNING id;`

		query, args, e := tx.BindNamed(sql, run)
		if e != nil {
			return errors.Wrap(e, "failed to bind")
		}

		if err = tx.QueryRowx(query, args...).Scan(&run.ID); err != nil {
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
		_, err = tx.NamedExec(sql, run.PipelineTaskRuns)
		return errors.Wrap(err, "failed to insert pipeline_task_runs")
	})
	return errors.Wrap(err, "InsertFinishedRun failed")
}

// DeleteRunsOlderThan deletes all pipeline_runs that have been finished for a certain threshold to free DB space
func (o *orm) DeleteRunsOlderThan(ctx context.Context, threshold time.Duration) error {
	// Addede 1 minute timeout to account for big databases
	q := o.q.WithOpts(pg.WithParentCtx(ctx), pg.WithLongQueryTimeout())

	err := pg.Batch(func(_, limit uint) (count uint, err error) {
		result, cancel, err := q.ExecQIter(`
WITH batched_pipeline_runs AS (
	SELECT * FROM pipeline_runs
	WHERE finished_at < ($1)
	ORDER BY finished_at ASC
	LIMIT $2
)
DELETE FROM pipeline_runs
USING batched_pipeline_runs
WHERE pipeline_runs.id = batched_pipeline_runs.id`,
			time.Now().Add(-threshold),
			limit,
		)
		defer cancel()
		if err != nil {
			return count, errors.Wrap(err, "DeleteRunsOlderThan failed to delete old pipeline_runs")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return count, errors.Wrap(err, "DeleteRunsOlderThan failed to get rows affected")
		}

		return uint(rowsAffected), err
	})
	if err != nil {
		return errors.Wrap(err, "DeleteRunsOlderThan failed")
	}

	return nil
}

func (o *orm) FindRun(id int64) (r Run, err error) {
	var runs []Run
	err = o.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Select(&runs, `SELECT * from pipeline_runs WHERE id = $1 LIMIT 1`, id); err != nil {
			return errors.Wrap(err, "failed to load runs")
		}
		return loadAssociations(tx, runs)
	})
	if len(runs) == 0 {
		return r, sql.ErrNoRows
	}
	return runs[0], err
}

func (o *orm) GetAllRuns() (runs []Run, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&runs, `SELECT * from pipeline_runs ORDER BY created_at ASC, id ASC`)
		if err != nil {
			return errors.Wrap(err, "failed to load runs")
		}

		return loadAssociations(tx, runs)
	})
	return runs, err
}

func (o *orm) GetUnfinishedRuns(ctx context.Context, now time.Time, fn func(run Run) error) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	return pg.Batch(func(offset, limit uint) (count uint, err error) {
		var runs []Run

		err = q.Transaction(func(tx pg.Queryer) error {
			err = tx.Select(&runs, `SELECT * from pipeline_runs WHERE state = $1 AND created_at < $2 ORDER BY created_at ASC, id ASC OFFSET $3 LIMIT $4`, RunStatusRunning, now, offset, limit)
			if err != nil {
				return errors.Wrap(err, "failed to load runs")
			}

			err = loadAssociations(tx, runs)
			if err != nil {
				return err
			}

			for _, run := range runs {
				if err = fn(run); err != nil {
					return err
				}
			}
			return nil
		})

		return uint(len(runs)), err
	})
}

// loads PipelineSpec and PipelineTaskRuns for Runs in exactly 2 queries
func loadAssociations(q pg.Queryer, runs []Run) error {
	if len(runs) == 0 {
		return nil
	}
	var specs []Spec
	pipelineSpecIDM := make(map[int32]Spec)
	var pipelineSpecIDs []int32 // keyed by pipelineSpecID
	pipelineRunIDs := make([]int64, len(runs))
	for i, run := range runs {
		pipelineRunIDs[i] = run.ID
		if _, exists := pipelineSpecIDM[run.PipelineSpecID]; !exists {
			pipelineSpecIDs = append(pipelineSpecIDs, run.PipelineSpecID)
			pipelineSpecIDM[run.PipelineSpecID] = Spec{}
		}
	}
	if err := q.Select(&specs, `SELECT * FROM pipeline_specs WHERE id = ANY($1)`, pipelineSpecIDs); err != nil {
		return errors.Wrap(err, "failed to postload pipeline_specs for runs")
	}
	for _, spec := range specs {
		pipelineSpecIDM[spec.ID] = spec
	}

	var taskRuns []TaskRun
	taskRunPRIDM := make(map[int64][]TaskRun, len(runs)) // keyed by pipelineRunID
	if err := q.Select(&taskRuns, `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = ANY($1) ORDER BY created_at ASC, id ASC`, pipelineRunIDs); err != nil {
		return errors.Wrap(err, "failed to postload pipeline_task_runs for runs")
	}
	for _, taskRun := range taskRuns {
		taskRunPRIDM[taskRun.PipelineRunID] = append(taskRunPRIDM[taskRun.PipelineRunID], taskRun)
	}

	for i, run := range runs {
		runs[i].PipelineSpec = pipelineSpecIDM[run.PipelineSpecID]
		runs[i].PipelineTaskRuns = taskRunPRIDM[run.ID]
	}

	return nil
}

func (o *orm) GetQ() pg.Q {
	return o.q
}
