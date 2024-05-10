package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/jmoiron/sqlx"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	valuespb "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

// `DBStore` is a postgres-backed
// data store that persists workflow progress.
type DBStore struct {
	db    sqlutil.DataSource
	clock clockwork.Clock
}

// `workflowExecutionRow` describes a row
// of the `workflow_executions` table
type workflowExecutionRow struct {
	ID         string
	WorkflowID *string
	Status     string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	FinishedAt *time.Time
}

// `workflowStepRow` describes a row
// of the `workflow_steps` table
type workflowStepRow struct {
	ID                  uint
	WorkflowExecutionID string `db:"workflow_execution_id"`
	Ref                 string
	Status              string
	Inputs              []byte
	OutputErr           *string    `db:"output_err"`
	OutputValue         []byte     `db:"output_value"`
	UpdatedAt           *time.Time `db:"updated_at"`
}

// `UpdateStatus` updates the status of the given workflow execution
func (d *DBStore) UpdateStatus(ctx context.Context, executionID string, status string) error {
	sql := `UPDATE workflow_executions SET status = $1, updated_at = $2 WHERE id = $3`

	// If we're completing the workflow execution, let's also set a finished_at timestamp.
	if status != StatusStarted {
		sql = "UPDATE workflow_executions SET status = $1, updated_at = $2, finished_at = $2 WHERE id = $3"
	}
	_, err := d.db.ExecContext(ctx, sql, status, d.clock.Now(), executionID)
	return err
}

// `UpsertStep` updates the given step. This will correspond to an insert, or an update
// depending on whether a step with the ref already exists.
func (d *DBStore) UpsertStep(ctx context.Context, stepState *WorkflowExecutionStep) (WorkflowExecution, error) {
	step, err := stateToStep(stepState)
	if err != nil {
		return WorkflowExecution{}, err
	}

	err = d.upsertSteps(ctx, []workflowStepRow{step})
	if err != nil {
		return WorkflowExecution{}, err
	}

	return d.Get(ctx, step.WorkflowExecutionID)
}

// `Get` fetches the ExecutionState from the database.
func (d *DBStore) Get(ctx context.Context, executionID string) (WorkflowExecution, error) {
	wex := &workflowExecutionRow{}
	err := d.db.GetContext(ctx, wex, `SELECT * FROM workflow_executions WHERE id = $1`, executionID)
	if err != nil {
		return WorkflowExecution{}, err
	}

	ws := []workflowStepRow{}
	err = d.db.SelectContext(ctx, &ws, `SELECT * FROM workflow_steps WHERE workflow_execution_id = $1`, wex.ID)
	if err != nil {
		return WorkflowExecution{}, err
	}

	refToStep := map[string]*WorkflowExecutionStep{}
	for _, s := range ws {
		ss, err := stepToState(s)
		if err != nil {
			return WorkflowExecution{}, err
		}

		refToStep[s.Ref] = ss
	}

	var workflowID string
	if wex.WorkflowID != nil {
		workflowID = *wex.WorkflowID
	}

	es := WorkflowExecution{
		ExecutionID: wex.ID,
		WorkflowID:  workflowID,
		Status:      wex.Status,
		Steps:       refToStep,
		CreatedAt:   wex.CreatedAt,
		UpdatedAt:   wex.UpdatedAt,
		FinishedAt:  wex.FinishedAt,
	}
	return es, nil
}

func stepToState(step workflowStepRow) (*WorkflowExecutionStep, error) {
	var inputs *values.Map
	if len(step.Inputs) > 0 {
		vmProto := &valuespb.Map{}
		err := proto.Unmarshal(step.Inputs, vmProto)
		if err != nil {
			return nil, err
		}

		inputs = values.FromMapValueProto(vmProto)
	}

	var (
		outputErr error
		outputs   values.Value
	)

	if step.OutputErr != nil {
		outputErr = errors.New(*step.OutputErr)
	}

	if len(step.OutputValue) != 0 {
		vProto := &valuespb.Value{}
		err := proto.Unmarshal(step.OutputValue, vProto)
		if err != nil {
			return nil, err
		}

		outputs = values.FromProto(vProto)
	}

	var so *StepOutput
	if outputErr != nil || outputs != nil {
		so = &StepOutput{
			Err:   outputErr,
			Value: outputs,
		}
	}

	return &WorkflowExecutionStep{
		ExecutionID: step.WorkflowExecutionID,
		Ref:         step.Ref,
		Status:      step.Status,
		Inputs:      inputs,
		Outputs:     so,
	}, nil
}

func stateToStep(state *WorkflowExecutionStep) (workflowStepRow, error) {
	var inpb []byte
	if state.Inputs != nil {
		p := values.Proto(state.Inputs).GetMapValue()
		ib, err := proto.Marshal(p)
		if err != nil {
			return workflowStepRow{}, err
		}
		inpb = ib
	}

	wsr := workflowStepRow{
		WorkflowExecutionID: state.ExecutionID,
		Ref:                 state.Ref,
		Status:              state.Status,
		Inputs:              inpb,
	}

	if state.Outputs == nil {
		return wsr, nil
	}

	if state.Outputs.Value != nil {
		p := values.Proto(state.Outputs.Value)
		ob, err := proto.Marshal(p)
		if err != nil {
			return workflowStepRow{}, err
		}

		wsr.OutputValue = ob
	}

	if state.Outputs.Err != nil {
		errs := state.Outputs.Err.Error()
		wsr.OutputErr = &errs
	}
	return wsr, nil
}

// `Add` creates the relevant workflow_execution and workflow_step entries
// to persist the passed in ExecutionState.
func (d *DBStore) Add(ctx context.Context, state *WorkflowExecution) error {
	return d.transact(ctx, func(db *DBStore) error {
		var wid *string
		if state.WorkflowID != "" {
			wid = &state.WorkflowID
		}
		wex := &workflowExecutionRow{
			ID:         state.ExecutionID,
			WorkflowID: wid,
			Status:     state.Status,
		}
		err := db.insertWorkflowExecution(ctx, wex)
		if err != nil {
			return fmt.Errorf("could not insert workflow execution %s: %w", state.ExecutionID, err)
		}

		ws := []workflowStepRow{}
		for _, step := range state.Steps {
			step, err := stateToStep(step)
			if err != nil {
				return err
			}
			ws = append(ws, step)
		}
		if len(ws) > 0 {
			return db.upsertSteps(ctx, ws)
		}
		return nil
	})
}

func (d *DBStore) upsertSteps(ctx context.Context, steps []workflowStepRow) error {
	for _, s := range steps {
		now := d.clock.Now()
		s.UpdatedAt = &now
	}

	sql := `
	INSERT INTO
	workflow_steps(workflow_execution_id, ref, status, inputs, output_err, output_value, updated_at)
	VALUES (:workflow_execution_id, :ref, :status, :inputs, :output_err, :output_value, :updated_at)
	ON CONFLICT ON CONSTRAINT uniq_workflow_execution_id_ref
	DO UPDATE SET
		workflow_execution_id = EXCLUDED.workflow_execution_id,
		ref = EXCLUDED.ref,
		status = EXCLUDED.status,
		inputs = EXCLUDED.inputs,
		output_err = EXCLUDED.output_err,
		output_value = EXCLUDED.output_value,
		updated_at = EXCLUDED.updated_at;
	`
	stmt, args, err := sqlx.Named(sql, steps)
	if err != nil {
		return err
	}
	stmt = d.db.Rebind(stmt)
	_, err = d.db.ExecContext(ctx, stmt, args...)
	return err
}

func (d *DBStore) insertWorkflowExecution(ctx context.Context, execution *workflowExecutionRow) error {
	sql := `
	INSERT INTO
	workflow_executions(id, workflow_id, status, created_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err := d.db.ExecContext(ctx, sql, execution.ID, execution.WorkflowID, execution.Status, d.clock.Now())
	return err
}

func (d *DBStore) transact(ctx context.Context, fn func(*DBStore) error) error {
	return sqlutil.Transact(
		ctx,
		func(ds sqlutil.DataSource) *DBStore {
			return &DBStore{db: ds, clock: d.clock}
		},
		d.db,
		nil,
		fn,
	)
}

func (d *DBStore) GetUnfinished(ctx context.Context, offset, limit int) ([]WorkflowExecution, error) {
	sql := `
	SELECT
		workflow_steps.workflow_execution_id AS ws_workflow_execution_id,
		workflow_steps.ref AS ws_ref,
		workflow_steps.status AS ws_status,
		workflow_steps.inputs AS ws_inputs,
		workflow_steps.output_err AS ws_output_err,
		workflow_steps.output_value AS ws_output_value,
		workflow_steps.updated_at AS ws_updated_at,
		workflow_executions.id AS we_id,
		workflow_executions.workflow_id AS we_workflow_id,
		workflow_executions.status AS we_status,
		workflow_executions.created_at AS we_created_at,
		workflow_executions.updated_at AS we_updated_at,
		workflow_executions.finished_at AS we_finished_at
	FROM workflow_executions
	JOIN workflow_steps
	ON  workflow_steps.workflow_execution_id = workflow_executions.id
	WHERE workflow_executions.status = $1
	ORDER BY workflow_executions.created_at DESC
	LIMIT $2
	OFFSET $3
	`
	joinRecords := []struct {
		// WorkflowExecutionStep fields
		WSWorkflowExecutionID string     `db:"ws_workflow_execution_id"`
		WSRef                 string     `db:"ws_ref"`
		WSStatus              string     `db:"ws_status"`
		WSInputs              []byte     `db:"ws_inputs"`
		WSOutputErr           *string    `db:"ws_output_err"`
		WSOutputValue         []byte     `db:"ws_output_value"`
		WSUpdatedAt           *time.Time `db:"ws_updated_at"`

		// WorkflowExecution fields
		WEID         string     `db:"we_id"`
		WEWorkflowID *string    `db:"we_workflow_id"`
		WEStatus     string     `db:"we_status"`
		WECreatedAt  *time.Time `db:"we_created_at"`
		WEUpdatedAt  *time.Time `db:"we_updated_at"`
		WEFinishedAt *time.Time `db:"we_finished_at"`
	}{}
	err := d.db.SelectContext(ctx, &joinRecords, sql, StatusStarted, limit, offset)
	if err != nil {
		return []WorkflowExecution{}, err
	}

	idToExecutionState := map[string]*WorkflowExecution{}
	for _, jr := range joinRecords {
		var wid string
		if jr.WEWorkflowID != nil {
			wid = *jr.WEWorkflowID
		}
		if _, ok := idToExecutionState[jr.WEID]; !ok {
			idToExecutionState[jr.WEID] = &WorkflowExecution{
				ExecutionID: jr.WEID,
				WorkflowID:  wid,
				Status:      jr.WEStatus,
				Steps:       map[string]*WorkflowExecutionStep{},
				CreatedAt:   jr.WECreatedAt,
				UpdatedAt:   jr.WEUpdatedAt,
				FinishedAt:  jr.WEFinishedAt,
			}
		}

		state, err := stepToState(workflowStepRow{
			WorkflowExecutionID: jr.WSWorkflowExecutionID,
			Ref:                 jr.WSRef,
			OutputErr:           jr.WSOutputErr,
			OutputValue:         jr.WSOutputValue,
			Inputs:              jr.WSInputs,
			Status:              jr.WSStatus,
			UpdatedAt:           jr.WSUpdatedAt,
		})
		if err != nil {
			return nil, err
		}

		es := idToExecutionState[jr.WEID]
		es.Steps[state.Ref] = state
	}

	states := []WorkflowExecution{}
	for _, s := range idToExecutionState {
		states = append(states, *s)
	}

	return states, nil
}

func NewDBStore(ds sqlutil.DataSource, clock clockwork.Clock) *DBStore {
	return &DBStore{db: ds, clock: clock}
}
