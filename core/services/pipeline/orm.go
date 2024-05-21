package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// KeepersObservationSource is the same for all keeper jobs and it is not persisted in DB
const KeepersObservationSource = `
    encode_check_upkeep_tx      [type=ethabiencode
                                 abi="checkUpkeep(uint256 id, address from)"
                                 data="{\"id\":$(jobSpec.upkeepID),\"from\":$(jobSpec.effectiveKeeperAddress)}"]
    check_upkeep_tx             [type=ethcall
                                 failEarly=true
                                 extractRevertReason=true
                                 evmChainID="$(jobSpec.evmChainID)"
                                 contract="$(jobSpec.contractAddress)"
                                 gasUnlimited=true
                                 gasPrice="$(jobSpec.gasPrice)"
                                 gasTipCap="$(jobSpec.gasTipCap)"
                                 gasFeeCap="$(jobSpec.gasFeeCap)"
                                 data="$(encode_check_upkeep_tx)"]
    decode_check_upkeep_tx      [type=ethabidecode
                                 abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
    calculate_perform_data_len  [type=length
                                 input="$(decode_check_upkeep_tx.performData)"]
    perform_data_lessthan_limit [type=lessthan
                                 left="$(calculate_perform_data_len)"
                                 right="$(jobSpec.maxPerformDataSize)"]
    check_perform_data_limit    [type=conditional
                                 failEarly=true
                                 data="$(perform_data_lessthan_limit)"]
    encode_perform_upkeep_tx    [type=ethabiencode
                                 abi="performUpkeep(uint256 id, bytes calldata performData)"
                                 data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
    simulate_perform_upkeep_tx  [type=ethcall
                                 extractRevertReason=true
                                 evmChainID="$(jobSpec.evmChainID)"
                                 contract="$(jobSpec.contractAddress)"
                                 from="$(jobSpec.effectiveKeeperAddress)"
                                 gasUnlimited=true
                                 data="$(encode_perform_upkeep_tx)"]
    decode_check_perform_tx     [type=ethabidecode
                                 abi="bool success"]
    check_success            	[type=conditional
                                 failEarly=true
                                 data="$(decode_check_perform_tx.success)"]
    perform_upkeep_tx        	[type=ethtx
                                 minConfirmations=0
                                 to="$(jobSpec.contractAddress)"
                                 from="[$(jobSpec.fromAddress)]"
                                 evmChainID="$(jobSpec.evmChainID)"
                                 data="$(encode_perform_upkeep_tx)"
                                 gasLimit="$(jobSpec.performUpkeepGasLimit)"
                                 txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
    encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> calculate_perform_data_len -> perform_data_lessthan_limit -> check_perform_data_limit -> encode_perform_upkeep_tx -> simulate_perform_upkeep_tx -> decode_check_perform_tx -> check_success -> perform_upkeep_tx
`

type CreateDataSource interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	services.Service

	CreateSpec(ctx context.Context, pipeline Pipeline, maxTaskTimeout models.Interval) (int32, error)
	CreateRun(ctx context.Context, run *Run) (err error)
	InsertRun(ctx context.Context, run *Run) error
	DeleteRun(ctx context.Context, id int64) error
	StoreRun(ctx context.Context, run *Run) (restart bool, err error)
	UpdateTaskRunResult(ctx context.Context, taskID uuid.UUID, result Result) (run Run, start bool, err error)
	InsertFinishedRun(ctx context.Context, run *Run, saveSuccessfulTaskRuns bool) (err error)
	InsertFinishedRunWithSpec(ctx context.Context, run *Run, saveSuccessfulTaskRuns bool) (err error)

	// InsertFinishedRuns inserts all the given runs into the database.
	// If saveSuccessfulTaskRuns is false, only errored runs are saved.
	InsertFinishedRuns(ctx context.Context, run []*Run, saveSuccessfulTaskRuns bool) (err error)

	DeleteRunsOlderThan(context.Context, time.Duration) error
	FindRun(ctx context.Context, id int64) (Run, error)
	GetAllRuns(ctx context.Context) ([]Run, error)
	GetUnfinishedRuns(context.Context, time.Time, func(run Run) error) error

	DataSource() sqlutil.DataSource
	WithDataSource(sqlutil.DataSource) ORM
	Transact(context.Context, func(ORM) error) error
}

type orm struct {
	services.StateMachine
	ds                sqlutil.DataSource
	lggr              logger.Logger
	maxSuccessfulRuns uint64
	// jobID => count
	pm   sync.Map
	wg   sync.WaitGroup
	ctx  context.Context
	cncl context.CancelFunc
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger, jobPipelineMaxSuccessfulRuns uint64) *orm {
	ctx, cancel := context.WithCancel(context.Background())
	return &orm{
		ds:                ds,
		lggr:              lggr.Named("PipelineORM"),
		maxSuccessfulRuns: jobPipelineMaxSuccessfulRuns,
		ctx:               ctx,
		cncl:              cancel,
	}
}

func (o *orm) Start(_ context.Context) error {
	return o.StartOnce("PipelineORM", func() error {
		var msg string
		if o.maxSuccessfulRuns == 0 {
			msg = "Pipeline runs saving is disabled for all jobs: MaxSuccessfulRuns=0"
		} else {
			msg = fmt.Sprintf("Pipeline runs will be pruned above per-job limit of MaxSuccessfulRuns=%d", o.maxSuccessfulRuns)
		}
		o.lggr.Info(msg)
		return nil
	})
}

func (o *orm) Close() error {
	return o.StopOnce("PipelineORM", func() error {
		o.cncl()
		o.wg.Wait()
		return nil
	})
}

func (o *orm) Name() string {
	return o.lggr.Name()
}

func (o *orm) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *orm) Transact(ctx context.Context, fn func(ORM) error) error {
	return sqlutil.Transact(ctx, func(tx sqlutil.DataSource) ORM {
		return o.withDataSource(tx)
	}, o.ds, nil, func(tx ORM) error {
		if err := tx.Start(ctx); err != nil {
			return fmt.Errorf("failed to start tx orm: %w", err)
		}
		defer func() {
			if err := tx.Close(); err != nil {
				o.lggr.Warnw("Error closing temporary transactional ORM", "err", err)
			}
		}()
		return fn(tx)
	})
}

func (o *orm) DataSource() sqlutil.DataSource { return o.ds }

func (o *orm) WithDataSource(ds sqlutil.DataSource) ORM { return o.withDataSource(ds) }

func (o *orm) withDataSource(ds sqlutil.DataSource) *orm {
	ctx, cancel := context.WithCancel(context.Background())
	return &orm{
		ds:                ds,
		lggr:              o.lggr,
		maxSuccessfulRuns: o.maxSuccessfulRuns,
		ctx:               ctx,
		cncl:              cancel,
	}
}

func (o *orm) transact(ctx context.Context, fn func(*orm) error) error {
	return sqlutil.Transact(ctx, o.withDataSource, o.ds, nil, fn)
}

func (o *orm) CreateSpec(ctx context.Context, pipeline Pipeline, maxTaskDuration models.Interval) (id int32, err error) {
	sql := `INSERT INTO pipeline_specs (dot_dag_source, max_task_duration, created_at)
	VALUES ($1, $2, NOW())
	RETURNING id;`
	err = o.ds.GetContext(ctx, &id, sql, pipeline.Source, maxTaskDuration)
	return id, errors.WithStack(err)
}

func (o *orm) CreateRun(ctx context.Context, run *Run) (err error) {
	if run.CreatedAt.IsZero() {
		return errors.New("run.CreatedAt must be set")
	}

	err = o.transact(ctx, func(tx *orm) error {
		if e := tx.InsertRun(ctx, run); e != nil {
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

		sql := `INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at);`
		_, err = tx.ds.NamedExecContext(ctx, sql, run.PipelineTaskRuns)
		return err
	})

	return errors.Wrap(err, "CreateRun failed")
}

// InsertRun inserts a run into the database
func (o *orm) InsertRun(ctx context.Context, run *Run) error {
	if run.Status() == RunStatusCompleted {
		defer o.prune(o.ds, run.PruningKey)
	}
	query, args, err := o.ds.BindNamed(`INSERT INTO pipeline_runs (pipeline_spec_id, pruning_key, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
		VALUES (:pipeline_spec_id, :pruning_key, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state)
		RETURNING *;`, run)
	if err != nil {
		return fmt.Errorf("error binding arg: %w", err)
	}
	return o.ds.GetContext(ctx, run, query, args...)
}

// StoreRun will persist a partially executed run before suspending, or finish a run.
// If `restart` is true, then new task run data is available and the run should be resumed immediately.
func (o *orm) StoreRun(ctx context.Context, run *Run) (restart bool, err error) {
	err = o.transact(ctx, func(tx *orm) error {
		finished := run.FinishedAt.Valid
		if !finished {
			// Lock the current run. This prevents races with /v2/resume
			sql := `SELECT id FROM pipeline_runs WHERE id = $1 FOR UPDATE;`
			if _, err = tx.ds.ExecContext(ctx, sql, run.ID); err != nil {
				return fmt.Errorf("failed to select pipeline run %d: %w", run.ID, err)
			}

			taskRuns := []TaskRun{}
			// Reload task runs, we want to check for any changes while the run was ongoing
			if err = tx.ds.SelectContext(ctx, &taskRuns, `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`, run.ID); err != nil {
				return fmt.Errorf("failed to select piepline task run %d: %w", run.ID, err)
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
			if _, err = tx.ds.NamedExecContext(ctx, `UPDATE pipeline_runs SET state = :state WHERE id = :id`, run); err != nil {
				return fmt.Errorf("failed to update pipeline run %d to %s: %w", run.ID, run.State, err)
			}
		} else {
			defer o.prune(tx.ds, run.PruningKey)
			// Simply finish the run, no need to do any sort of locking
			if run.Outputs.Val == nil || len(run.FatalErrors)+len(run.AllErrors) == 0 {
				return fmt.Errorf("run must have both Outputs and Errors, got Outputs: %#v, FatalErrors: %#v, AllErrors: %#v", run.Outputs.Val, run.FatalErrors, run.AllErrors)
			}
			sql := `UPDATE pipeline_runs SET state = :state, finished_at = :finished_at, all_errors= :all_errors, fatal_errors= :fatal_errors, outputs = :outputs WHERE id = :id`
			if _, err = tx.ds.NamedExecContext(ctx, sql, run); err != nil {
				return fmt.Errorf("failed to update pipeline run %d: %w", run.ID, err)
			}
		}

		sql := `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at)
		ON CONFLICT (pipeline_run_id, dot_id) DO UPDATE SET
		output = EXCLUDED.output, error = EXCLUDED.error, finished_at = EXCLUDED.finished_at
		RETURNING *;
		`

		taskRuns := []TaskRun{}
		query, args, err := tx.ds.BindNamed(sql, run.PipelineTaskRuns)
		if err != nil {
			return fmt.Errorf("failed to prepare named query: %w", err)
		}
		err = tx.ds.SelectContext(ctx, &taskRuns, query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert pipeline task runs: %w", err)
		}
		run.PipelineTaskRuns = taskRuns
		return nil
	})
	return
}

// DeleteRun cleans up a run that failed and is marked failEarly (should leave no trace of the run)
func (o *orm) DeleteRun(ctx context.Context, id int64) error {
	// NOTE: this will cascade and wipe pipeline_task_runs too
	_, err := o.ds.ExecContext(ctx, `DELETE FROM pipeline_runs WHERE id = $1`, id)
	return err
}

func (o *orm) UpdateTaskRunResult(ctx context.Context, taskID uuid.UUID, result Result) (run Run, start bool, err error) {
	if result.OutputDB().Valid && result.ErrorDB().Valid {
		panic("run result must specify either output or error, not both")
	}
	err = o.transact(ctx, func(tx *orm) error {
		sql := `
		SELECT pipeline_runs.*, pipeline_specs.dot_dag_source "pipeline_spec.dot_dag_source", job_pipeline_specs.job_id "job_id"
		FROM pipeline_runs
		JOIN pipeline_task_runs ON (pipeline_task_runs.pipeline_run_id = pipeline_runs.id)
		JOIN pipeline_specs ON (pipeline_specs.id = pipeline_runs.pipeline_spec_id)
		JOIN job_pipeline_specs ON (job_pipeline_specs.pipeline_spec_id = pipeline_specs.id)
		WHERE pipeline_task_runs.id = $1 AND pipeline_runs.state in ('running', 'suspended')
		FOR UPDATE`
		if err = tx.ds.GetContext(ctx, &run, sql, taskID); err != nil {
			return fmt.Errorf("failed to find pipeline run for task ID %s: %w", taskID.String(), err)
		}

		// Update the task with result
		sql = `UPDATE pipeline_task_runs SET output = $2, error = $3, finished_at = $4 WHERE id = $1`
		if _, err = tx.ds.ExecContext(ctx, sql, taskID, result.OutputDB(), result.ErrorDB(), time.Now()); err != nil {
			return fmt.Errorf("failed to update pipeline task run: %w", err)
		}

		if run.State == RunStatusSuspended {
			start = true
			run.State = RunStatusRunning

			sql = `UPDATE pipeline_runs SET state = $2 WHERE id = $1`
			if _, err = tx.ds.ExecContext(ctx, sql, run.ID, run.State); err != nil {
				return fmt.Errorf("failed to update pipeline run state: %w", err)
			}
		}

		return loadAssociations(ctx, tx.ds, []*Run{&run})
	})

	return run, start, err
}

// InsertFinishedRuns inserts all the given runs into the database.
func (o *orm) InsertFinishedRuns(ctx context.Context, runs []*Run, saveSuccessfulTaskRuns bool) error {
	err := o.transact(ctx, func(tx *orm) error {
		pipelineRunsQuery := `
INSERT INTO pipeline_runs 
	(pipeline_spec_id, pruning_key, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
VALUES 
	(:pipeline_spec_id, :pruning_key, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state) 
RETURNING id
	`

		var runIDs []int64
		err := sqlutil.NamedQueryContext(ctx, tx.ds, pipelineRunsQuery, runs, func(row sqlutil.RowScanner) error {
			var runID int64
			if errS := row.Scan(&runID); errS != nil {
				return errors.Wrap(errS, "scanning pipeline runs id row")
			}
			runIDs = append(runIDs, runID)
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "inserting finished pipeline runs")
		}

		pruningKeysm := make(map[int32]struct{})
		for i, run := range runs {
			pruningKeysm[run.PruningKey] = struct{}{}
			for j := range run.PipelineTaskRuns {
				run.PipelineTaskRuns[j].PipelineRunID = runIDs[i]
			}
		}

		defer func() {
			for pruningKey := range pruningKeysm {
				o.prune(tx.ds, pruningKey)
			}
		}()

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

		_, errE := tx.ds.NamedExecContext(ctx, pipelineTaskRunsQuery, pipelineTaskRuns)
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
	if run.Outputs.Val == nil || len(run.FatalErrors)+len(run.AllErrors) == 0 {
		return errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, FatalErrors: %#v, AllErrors: %#v", run.Outputs.Val, run.FatalErrors, run.AllErrors)
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
func (o *orm) InsertFinishedRun(ctx context.Context, run *Run, saveSuccessfulTaskRuns bool) (err error) {
	if err = o.checkFinishedRun(run, saveSuccessfulTaskRuns); err != nil {
		return err
	}

	if o.maxSuccessfulRuns == 0 {
		// optimisation: avoid persisting if we oughtn't to save any
		return nil
	}

	err = o.insertFinishedRun(ctx, run, saveSuccessfulTaskRuns)
	return errors.Wrap(err, "InsertFinishedRun failed")
}

// InsertFinishedRunWithSpec works like InsertFinishedRun but also inserts the pipeline spec.
func (o *orm) InsertFinishedRunWithSpec(ctx context.Context, run *Run, saveSuccessfulTaskRuns bool) (err error) {
	if err = o.checkFinishedRun(run, saveSuccessfulTaskRuns); err != nil {
		return err
	}

	if o.maxSuccessfulRuns == 0 {
		// optimisation: avoid persisting if we oughtn't to save any
		return nil
	}

	err = o.transact(ctx, func(tx *orm) error {
		sqlStmt1 := `INSERT INTO pipeline_specs (dot_dag_source, max_task_duration, created_at)
	VALUES ($1, $2, NOW())
	RETURNING id;`
		err = tx.ds.GetContext(ctx, &run.PipelineSpecID, sqlStmt1, run.PipelineSpec.DotDagSource, run.PipelineSpec.MaxTaskDuration)
		if err != nil {
			return errors.Wrap(err, "failed to insert pipeline_specs")
		}
		// This `job_pipeline_specs` record won't be primary since when this method is called, the job already exists, so it will have primary record.
		sqlStmt2 := `INSERT INTO job_pipeline_specs (job_id, pipeline_spec_id, is_primary) VALUES ($1, $2, false)`
		_, err = tx.ds.ExecContext(ctx, sqlStmt2, run.JobID, run.PipelineSpecID)
		if err != nil {
			return errors.Wrap(err, "failed to insert job_pipeline_specs")
		}
		return tx.insertFinishedRun(ctx, run, saveSuccessfulTaskRuns)
	})
	return errors.Wrap(err, "InsertFinishedRun failed")
}

func (o *orm) insertFinishedRun(ctx context.Context, run *Run, saveSuccessfulTaskRuns bool) error {
	sql := `INSERT INTO pipeline_runs (pipeline_spec_id, pruning_key, meta, all_errors, fatal_errors, inputs, outputs, created_at, finished_at, state)
		VALUES (:pipeline_spec_id, :pruning_key, :meta, :all_errors, :fatal_errors, :inputs, :outputs, :created_at, :finished_at, :state)
		RETURNING id;`

	query, args, err := o.ds.BindNamed(sql, run)
	if err != nil {
		return errors.Wrap(err, "failed to bind")
	}

	if err = o.ds.QueryRowxContext(ctx, query, args...).Scan(&run.ID); err != nil {
		return errors.Wrap(err, "error inserting finished pipeline_run")
	}

	// update the ID key everywhere
	for i := range run.PipelineTaskRuns {
		run.PipelineTaskRuns[i].PipelineRunID = run.ID
	}

	if !saveSuccessfulTaskRuns && !run.HasErrors() {
		return nil
	}

	defer o.prune(o.ds, run.PruningKey)
	sql = `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at);`
	_, err = o.ds.NamedExecContext(ctx, sql, run.PipelineTaskRuns)
	return errors.Wrap(err, "failed to insert pipeline_task_runs")
}

// DeleteRunsOlderThan deletes all pipeline_runs that have been finished for a certain threshold to free DB space
// Caller is expected to set timeout on calling context.
func (o *orm) DeleteRunsOlderThan(ctx context.Context, threshold time.Duration) error {
	start := time.Now()

	queryThreshold := start.Add(-threshold)

	rowsDeleted := int64(0)

	err := pg.Batch(func(_, limit uint) (count uint, err error) {
		result, err := o.ds.ExecContext(ctx, `
WITH batched_pipeline_runs AS (
	SELECT * FROM pipeline_runs
	WHERE finished_at < ($1)
	ORDER BY finished_at ASC
	LIMIT $2
)
DELETE FROM pipeline_runs
USING batched_pipeline_runs
WHERE pipeline_runs.id = batched_pipeline_runs.id`,
			queryThreshold,
			limit,
		)
		if err != nil {
			return count, errors.Wrap(err, "DeleteRunsOlderThan failed to delete old pipeline_runs")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return count, errors.Wrap(err, "DeleteRunsOlderThan failed to get rows affected")
		}
		rowsDeleted += rowsAffected

		return uint(rowsAffected), err
	})
	if err != nil {
		return errors.Wrap(err, "DeleteRunsOlderThan failed")
	}

	deleteTS := time.Now()

	o.lggr.Debugw("pipeline_runs reaper DELETE query completed", "rowsDeleted", rowsDeleted, "duration", deleteTS.Sub(start))
	defer func(start time.Time) {
		o.lggr.Debugw("pipeline_runs reaper VACUUM ANALYZE query completed", "duration", time.Since(start))
	}(deleteTS)

	_, err = o.ds.ExecContext(ctx, "VACUUM ANALYZE pipeline_runs")
	if err != nil {
		o.lggr.Warnw("DeleteRunsOlderThan successfully deleted old pipeline_runs rows, but failed to run VACUUM ANALYZE", "err", err)
		return nil
	}

	return nil
}

func (o *orm) FindRun(ctx context.Context, id int64) (r Run, err error) {
	var runs []*Run
	err = o.transact(ctx, func(tx *orm) error {
		if err = tx.ds.SelectContext(ctx, &runs, `SELECT * from pipeline_runs WHERE id = $1 LIMIT 1`, id); err != nil {
			return errors.Wrap(err, "failed to load runs")
		}
		return loadAssociations(ctx, tx.ds, runs)
	})
	if len(runs) == 0 {
		return r, sql.ErrNoRows
	}
	return *runs[0], err
}

func (o *orm) GetAllRuns(ctx context.Context) (runs []Run, err error) {
	var runsPtrs []*Run
	err = o.transact(ctx, func(tx *orm) error {
		err = tx.ds.SelectContext(ctx, &runsPtrs, `SELECT * from pipeline_runs ORDER BY created_at ASC, id ASC`)
		if err != nil {
			return errors.Wrap(err, "failed to load runs")
		}

		return loadAssociations(ctx, tx.ds, runsPtrs)
	})
	runs = make([]Run, len(runsPtrs))
	for i, runPtr := range runsPtrs {
		runs[i] = *runPtr
	}
	return runs, err
}

func (o *orm) GetUnfinishedRuns(ctx context.Context, now time.Time, fn func(run Run) error) error {
	return pg.Batch(func(offset, limit uint) (count uint, err error) {
		var runs []*Run

		err = o.transact(ctx, func(tx *orm) error {
			err = tx.ds.SelectContext(ctx, &runs, `SELECT * from pipeline_runs WHERE state = $1 AND created_at < $2 ORDER BY created_at ASC, id ASC OFFSET $3 LIMIT $4`, RunStatusRunning, now, offset, limit)
			if err != nil {
				return errors.Wrap(err, "failed to load runs")
			}

			err = loadAssociations(ctx, tx.ds, runs)
			if err != nil {
				return err
			}

			for _, run := range runs {
				if err = fn(*run); err != nil {
					return err
				}
			}
			return nil
		})

		return uint(len(runs)), err
	})
}

// loads PipelineSpec and PipelineTaskRuns for Runs in exactly 2 queries
func loadAssociations(ctx context.Context, ds sqlutil.DataSource, runs []*Run) error {
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
	sqlQuery := `SELECT
			ps.id,
			ps.dot_dag_source,
			ps.created_at,
			ps.max_task_duration,
			coalesce(jobs.id, 0) "job_id",
			coalesce(jobs.name, '') "job_name",
			coalesce(jobs.type, '') "job_type"
		FROM pipeline_specs ps
		LEFT JOIN job_pipeline_specs jps ON jps.pipeline_spec_id=ps.id
		LEFT JOIN jobs ON jobs.id=jps.job_id
		WHERE ps.id = ANY($1)`
	if err := ds.SelectContext(ctx, &specs, sqlQuery, pipelineSpecIDs); err != nil {
		return errors.Wrap(err, "failed to postload pipeline_specs for runs")
	}
	for _, spec := range specs {
		if spec.JobType == "keeper" {
			spec.DotDagSource = KeepersObservationSource
		}
		pipelineSpecIDM[spec.ID] = spec
	}

	var taskRuns []TaskRun
	taskRunPRIDM := make(map[int64][]TaskRun, len(runs)) // keyed by pipelineRunID
	if err := ds.SelectContext(ctx, &taskRuns, `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = ANY($1) ORDER BY created_at ASC, id ASC`, pipelineRunIDs); err != nil {
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

func (o *orm) loadCount(jobID int32) *atomic.Uint64 {
	// fast path; avoids allocation
	actual, exists := o.pm.Load(jobID)
	if exists {
		return actual.(*atomic.Uint64)
	}
	// "slow" path
	actual, _ = o.pm.LoadOrStore(jobID, new(atomic.Uint64))
	return actual.(*atomic.Uint64)
}

// Runs will be pruned async on a sampled basis if maxSuccessfulRuns is set to
// this value or higher
const syncLimit = 1000

// prune attempts to keep the pipeline_runs table capped close to the
// maxSuccessfulRuns length for each job_id.
//
// It does this synchronously for small values and async/sampled for large
// values.
//
// Note this does not guarantee the pipeline_runs table is kept to exactly the
// max length, rather that it doesn't excessively larger than it.
func (o *orm) prune(tx sqlutil.DataSource, jobID int32) {
	if jobID == 0 {
		o.lggr.Panic("expected a non-zero job ID")
	}
	// For small maxSuccessfulRuns its fast enough to prune every time
	if o.maxSuccessfulRuns < syncLimit {
		o.withDataSource(tx).execPrune(o.ctx, jobID)
		return
	}
	// for large maxSuccessfulRuns we do it async on a sampled basis
	every := o.maxSuccessfulRuns / 20 // it can get up to 5% larger than maxSuccessfulRuns before a prune
	cnt := o.loadCount(jobID)
	val := cnt.Add(1)
	if val%every == 0 {
		ok := o.IfStarted(func() {
			o.wg.Add(1)
			go func() {
				o.lggr.Debugw("Pruning runs", "jobID", jobID, "count", val, "every", every, "maxSuccessfulRuns", o.maxSuccessfulRuns)
				defer o.wg.Done()
				ctx, cancel := context.WithTimeout(sqlutil.WithoutDefaultTimeout(o.ctx), time.Minute)
				defer cancel()

				// Must not use tx here since it could be stale by the time we execute async.
				o.execPrune(ctx, jobID)
			}()
		})
		if !ok {
			o.lggr.Warnw("Cannot prune: ORM is not running", "jobID", jobID)
			return
		}
	}
}

func (o *orm) execPrune(ctx context.Context, jobID int32) {
	res, err := o.ds.ExecContext(o.ctx, `DELETE FROM pipeline_runs WHERE pruning_key = $1 AND state = $2 AND id NOT IN (
SELECT id FROM pipeline_runs
WHERE pruning_key = $1 AND state = $2
ORDER BY id DESC
LIMIT $3
)`, jobID, RunStatusCompleted, o.maxSuccessfulRuns)
	if err != nil {
		o.lggr.Errorw("Failed to prune runs", "err", err, "jobID", jobID)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		o.lggr.Errorw("Failed to get RowsAffected while pruning runs", "err", err, "jobID", jobID)
		return
	}
	if rowsAffected == 0 {
		// check the spec still exists and garbage collect if necessary
		var exists bool
		if err := o.ds.GetContext(ctx, &exists, `SELECT EXISTS(SELECT ps.* FROM pipeline_specs ps JOIN job_pipeline_specs jps ON (ps.id=jps.pipeline_spec_id) WHERE jps.job_id = $1)`, jobID); err != nil {
			o.lggr.Errorw("Failed check existence of pipeline_spec while pruning runs", "err", err, "jobID", jobID)
			return
		}
		if !exists {
			o.lggr.Debugw("Pipeline spec no longer exists, removing prune count", "jobID", jobID)
			o.pm.Delete(jobID)
		}
	} else if o.maxSuccessfulRuns < syncLimit {
		o.lggr.Tracew("Pruned runs", "rowsAffected", rowsAffected, "jobID", jobID)
	} else {
		o.lggr.Debugw("Pruned runs", "rowsAffected", rowsAffected, "jobID", jobID)
	}
}
