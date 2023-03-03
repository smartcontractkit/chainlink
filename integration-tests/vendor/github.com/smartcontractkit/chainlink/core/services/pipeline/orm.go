package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/sqlx"
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

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	services.ServiceCtx
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
	utils.StartStopOnce
	q                 pg.Q
	lggr              logger.Logger
	maxSuccessfulRuns uint64
	// jobID => count
	pm   sync.Map
	wg   sync.WaitGroup
	ctx  context.Context
	cncl context.CancelFunc
}

var _ ORM = (*orm)(nil)

type ORMConfig interface {
	pg.QConfig
	JobPipelineMaxSuccessfulRuns() uint64
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg ORMConfig) *orm {
	ctx, cancel := context.WithCancel(context.Background())
	return &orm{
		utils.StartStopOnce{},
		pg.NewQ(db, lggr, cfg),
		lggr.Named("PipelineORM"),
		cfg.JobPipelineMaxSuccessfulRuns(),
		sync.Map{},
		sync.WaitGroup{},
		ctx,
		cancel,
	}
}

func (o *orm) Start(_ context.Context) error {
	return o.StartOnce("pipeline.ORM", func() error {
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
	return o.StopOnce("pipeline.ORM", func() error {
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
	if run.Status() == RunStatusCompleted {
		defer o.Prune(o.q, run.PipelineSpecID)
	}
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
			defer o.Prune(tx, run.PipelineSpecID)
			// Simply finish the run, no need to do any sort of locking
			if run.Outputs.Val == nil || len(run.FatalErrors)+len(run.AllErrors) == 0 {
				return errors.Errorf("run must have both Outputs and Errors, got Outputs: %#v, FatalErrors: %#v, AllErrors: %#v", run.Outputs.Val, run.FatalErrors, run.AllErrors)
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
	if result.OutputDB().Valid && result.ErrorDB().Valid {
		panic("run result must specify either output or error, not both")
	}
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

			sql = `UPDATE pipeline_runs SET state = $2 WHERE id = $1`
			if _, err = tx.Exec(sql, run.ID, run.State); err != nil {
				return errors.Wrap(err, "UpdateTaskRunResult")
			}
		}

		return loadAssociations(tx, []*Run{&run})
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

		pipelineSpecIDm := make(map[int32]struct{})
		for i, run := range runs {
			pipelineSpecIDm[run.PipelineSpecID] = struct{}{}
			for j := range run.PipelineTaskRuns {
				run.PipelineTaskRuns[j].PipelineRunID = runIDs[i]
			}
		}

		defer func() {
			for pipelineSpecID := range pipelineSpecIDm {
				o.Prune(tx, pipelineSpecID)
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
func (o *orm) InsertFinishedRun(run *Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) (err error) {
	if err = o.checkFinishedRun(run, saveSuccessfulTaskRuns); err != nil {
		return err
	}

	if o.maxSuccessfulRuns == 0 {
		// optimisation: avoid persisting if we oughtn't to save any
		return nil
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

		defer o.Prune(tx, run.PipelineSpecID)
		sql = `
		INSERT INTO pipeline_task_runs (pipeline_run_id, id, type, index, output, error, dot_id, created_at, finished_at)
		VALUES (:pipeline_run_id, :id, :type, :index, :output, :error, :dot_id, :created_at, :finished_at);`
		_, err = tx.NamedExec(sql, run.PipelineTaskRuns)
		return errors.Wrap(err, "failed to insert pipeline_task_runs")
	})
	return errors.Wrap(err, "InsertFinishedRun failed")
}

// DeleteRunsOlderThan deletes all pipeline_runs that have been finished for a certain threshold to free DB space
// Caller is expected to set timeout on calling context.
func (o *orm) DeleteRunsOlderThan(ctx context.Context, threshold time.Duration) error {
	start := time.Now()

	q := o.q.WithOpts(pg.WithParentCtxInheritTimeout(ctx))

	queryThreshold := start.Add(-threshold)

	rowsDeleted := int64(0)

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
			queryThreshold,
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

	err = q.ExecQ("VACUUM ANALYZE pipeline_runs")
	if err != nil {
		o.lggr.Warnw("DeleteRunsOlderThan successfully deleted old pipeline_runs rows, but failed to run VACUUM ANALYZE", "err", err)
		return nil
	}

	return nil
}

func (o *orm) FindRun(id int64) (r Run, err error) {
	var runs []*Run
	err = o.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Select(&runs, `SELECT * from pipeline_runs WHERE id = $1 LIMIT 1`, id); err != nil {
			return errors.Wrap(err, "failed to load runs")
		}
		return loadAssociations(tx, runs)
	})
	if len(runs) == 0 {
		return r, sql.ErrNoRows
	}
	return *runs[0], err
}

func (o *orm) GetAllRuns() (runs []Run, err error) {
	var runsPtrs []*Run
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&runsPtrs, `SELECT * from pipeline_runs ORDER BY created_at ASC, id ASC`)
		if err != nil {
			return errors.Wrap(err, "failed to load runs")
		}

		return loadAssociations(tx, runsPtrs)
	})
	runs = make([]Run, len(runsPtrs))
	for i, runPtr := range runsPtrs {
		runs[i] = *runPtr
	}
	return runs, err
}

func (o *orm) GetUnfinishedRuns(ctx context.Context, now time.Time, fn func(run Run) error) error {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	return pg.Batch(func(offset, limit uint) (count uint, err error) {
		var runs []*Run

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
func loadAssociations(q pg.Queryer, runs []*Run) error {
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
	if err := q.Select(&specs, `SELECT ps.id, ps.dot_dag_source, ps.created_at, ps.max_task_duration, coalesce(jobs.id, 0) "job_id", coalesce(jobs.name, '') "job_name", coalesce(jobs.type, '') "job_type" FROM pipeline_specs ps LEFT OUTER JOIN jobs ON jobs.pipeline_spec_id=ps.id WHERE ps.id = ANY($1)`, pipelineSpecIDs); err != nil {
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

func (o *orm) loadCount(pipelineSpecID int32) *atomic.Uint64 {
	// fast path; avoids allocation
	actual, exists := o.pm.Load(pipelineSpecID)
	if exists {
		return actual.(*atomic.Uint64)
	}
	// "slow" path
	actual, _ = o.pm.LoadOrStore(pipelineSpecID, new(atomic.Uint64))
	return actual.(*atomic.Uint64)
}

// Runs will be pruned async on a sampled basis if maxSuccessfulRuns is set to
// this value or higher
const syncLimit = 1000

// Prune attempts to keep the pipeline_runs table capped close to the
// maxSuccessfulRuns length for each pipeline_spec_id.
//
// It does this synchronously for small values and async/sampled for large
// values.
//
// Note this does not guarantee the pipeline_runs table is kept to exactly the
// max length, rather that it doesn't excessively larger than it.
func (o *orm) Prune(tx pg.Queryer, pipelineSpecID int32) {
	if pipelineSpecID == 0 {
		o.lggr.Panic("expected a non-zero pipeline spec ID")
	}
	// For small maxSuccessfulRuns its fast enough to prune every time
	if o.maxSuccessfulRuns < syncLimit {
		o.execPrune(tx, pipelineSpecID)
		return
	}
	// for large maxSuccessfulRuns we do it async on a sampled basis
	every := o.maxSuccessfulRuns / 20 // it can get up to 5% larger than maxSuccessfulRuns before a prune
	cnt := o.loadCount(pipelineSpecID)
	val := cnt.Add(1)
	if val%every == 0 {
		ok := o.IfStarted(func() {
			o.wg.Add(1)
			go func() {
				o.lggr.Debugw("Pruning runs", "pipelineSpecID", pipelineSpecID, "count", val, "every", every, "maxSuccessfulRuns", o.maxSuccessfulRuns)
				defer o.wg.Done()
				// Must not use tx here since it's async and the transaction
				// could be stale
				o.execPrune(o.q.WithOpts(pg.WithLongQueryTimeout()), pipelineSpecID)
			}()
		})
		if !ok {
			o.lggr.Warnw("Cannot prune: ORM is not running", "pipelineSpecID", pipelineSpecID)
			return
		}
	}
}

func (o *orm) execPrune(q pg.Queryer, pipelineSpecID int32) {
	res, err := q.ExecContext(o.ctx, `DELETE FROM pipeline_runs WHERE pipeline_spec_id = $1 AND state = $2 AND id NOT IN (
SELECT id FROM pipeline_runs
WHERE pipeline_spec_id = $1 AND state = $2
ORDER BY id DESC
LIMIT $3
)`, pipelineSpecID, RunStatusCompleted, o.maxSuccessfulRuns)
	if err != nil {
		o.lggr.Errorw("Failed to prune runs", "err", err, "pipelineSpecID", pipelineSpecID)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		o.lggr.Errorw("Failed to get RowsAffected while pruning runs", "err", err, "pipelineSpecID", pipelineSpecID)
		return
	}
	if rowsAffected == 0 {
		// check the spec still exists and garbage collect if necessary
		var exists bool
		if err := q.GetContext(o.ctx, &exists, `SELECT EXISTS(SELECT * FROM pipeline_specs WHERE id = $1)`, pipelineSpecID); err != nil {
			o.lggr.Errorw("Failed check existence of pipeline_spec while pruning runs", "err", err, "pipelineSpecID", pipelineSpecID)
			return
		}
		if !exists {
			o.lggr.Debugw("Pipeline spec no longer exists, removing prune count", "pipelineSpecID", pipelineSpecID)
			o.pm.Delete(pipelineSpecID)
		}
	} else if o.maxSuccessfulRuns < syncLimit {
		o.lggr.Tracew("Pruned runs", "rowsAffected", rowsAffected, "pipelineSpecID", pipelineSpecID)
	} else {
		o.lggr.Debugw("Pruned runs", "rowsAffected", rowsAffected, "pipelineSpecID", pipelineSpecID)
	}
}
