package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func randomID() string {
	return random(32)
}

func random(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func newTestDBStore(t *testing.T) *DBStore {
	db := pgtest.NewSqlxDB(t)
	return &DBStore{db: db, lggr: logger.TestLogger(t), clock: clockwork.NewFakeClock()}
}

func Test_StoreDB(t *testing.T) {
	store := newTestDBStore(t)

	id := randomID()
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": {
				ExecutionID: id,
				Ref:         "step1",
				Status:      StatusCompleted,
			},
			"step2": {
				ExecutionID: id,
				Ref:         "step2",
				Status:      StatusStarted,
			},
		},
		ExecutionID: id,
		Status:      StatusStarted,
	}

	_, err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	gotEs, err := store.Get(tests.Context(t), es.ExecutionID)
	// Zero out the created at timestamp; this isn't present on `es`
	// but is added by the db store.
	gotEs.CreatedAt = nil
	require.NoError(t, err)
	assert.Equal(t, es, gotEs)
}

func Test_StoreDB_DuplicateEntry(t *testing.T) {
	store := newTestDBStore(t)

	id := randomID()
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": {
				ExecutionID: id,
				Ref:         "step1",
				Status:      StatusCompleted,
			},
			"step2": {
				ExecutionID: id,
				Ref:         "step2",
				Status:      StatusStarted,
			},
		},
		ExecutionID: id,
		Status:      StatusStarted,
	}

	_, err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	_, err = store.Add(tests.Context(t), &es)
	assert.ErrorContains(t, err, "duplicate key value violates")
}

func Test_StoreDB_UpdateStatus(t *testing.T) {
	store := newTestDBStore(t)

	id := randomID()
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": {
				ExecutionID: id,
				Ref:         "step1",
				Status:      StatusCompleted,
			},
			"step2": {
				ExecutionID: id,
				Ref:         "step2",
				Status:      StatusStarted,
			},
		},
		ExecutionID: id,
		Status:      StatusStarted,
	}

	_, err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	completedStatus := StatusCompleted
	err = store.UpdateStatus(tests.Context(t), es.ExecutionID, StatusCompleted)
	require.NoError(t, err)

	gotEs, err := store.Get(tests.Context(t), es.ExecutionID)
	require.NoError(t, err)

	assert.Equal(t, gotEs.Status, completedStatus)
}

func Test_StoreDB_UpdateStep(t *testing.T) {
	store := newTestDBStore(t)

	id := randomID()
	stepOne := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step1",
		Status:      StatusCompleted,
	}
	stepTwo := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step2",
		Status:      StatusStarted,
	}
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": stepOne,
			"step2": stepTwo,
		},
		ExecutionID: id,
		Status:      StatusStarted,
	}

	_, err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	stepOne.Status = StatusCompleted
	nm, err := values.NewMap(map[string]any{"hello": "world"})
	require.NoError(t, err)

	stepOne.Inputs = nm
	stepOne.Outputs = StepOutput{Err: errors.New("some error")}

	es, err = store.UpsertStep(tests.Context(t), stepOne)
	require.NoError(t, err)

	gotStep := es.Steps[stepOne.Ref]
	assert.Equal(t, stepOne, gotStep)

	stepTwo.Outputs = StepOutput{Value: nm}
	es, err = store.UpsertStep(tests.Context(t), stepTwo)
	require.NoError(t, err)

	gotStep = es.Steps[stepTwo.Ref]
	assert.Equal(t, stepTwo, gotStep)
}

func Test_StoreDB_WorkflowStatus(t *testing.T) {
	store := newTestDBStore(t)

	for s := range ValidStatuses {
		id := randomID()
		stepOne := &WorkflowExecutionStep{
			ExecutionID: id,
			Ref:         "step1",
			Status:      StatusCompleted,
		}
		stepTwo := &WorkflowExecutionStep{
			ExecutionID: id,
			Ref:         "step2",
			Status:      StatusStarted,
		}
		es := WorkflowExecution{
			Steps: map[string]*WorkflowExecutionStep{
				"step1": stepOne,
				"step2": stepTwo,
			},
			ExecutionID: id,
			Status:      s,
		}

		_, err := store.Add(tests.Context(t), &es)
		require.NoError(t, err)
	}
}

func Test_StoreDB_WorkflowStepStatus(t *testing.T) {
	store := newTestDBStore(t)

	id := randomID()
	stepOne := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step1",
		Status:      StatusCompleted,
	}
	stepTwo := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step2",
		Status:      StatusStarted,
	}
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": stepOne,
			"step2": stepTwo,
		},
		ExecutionID: id,
		Status:      StatusStarted,
	}

	_, err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	for s := range ValidStatuses {
		stepOne.Status = s
		_, err := store.UpsertStep(tests.Context(t), stepOne)
		require.NoError(t, err)
	}
}

func createWorkflow(t *testing.T, store *DBStore, id string) {
	sql := `INSERT INTO workflow_specs (workflow, workflow_id, workflow_owner, workflow_name, created_at, updated_at)
	VALUES (:workflow, :workflow_id, :workflow_owner, :workflow_name, NOW(), NOW())`
	var wfSpec job.WorkflowSpec
	wfSpec.Workflow = ""
	wfSpec.WorkflowID = id
	wfSpec.WorkflowOwner = random(20)
	wfSpec.WorkflowName = random(20)
	_, err := store.db.NamedExecContext(tests.Context(t), sql, wfSpec)
	require.NoError(t, err)
}

func Test_StoreDB_GetUnfinishedSteps(t *testing.T) {
	store := newTestDBStore(t)

	id := randomID()
	wid := randomID()
	createWorkflow(t, store, wid)
	stepOne := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step1",
		Status:      StatusCompleted,
	}
	stepTwo := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step2",
		Status:      StatusStarted,
	}
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": stepOne,
			"step2": stepTwo,
		},
		ExecutionID: id,
		WorkflowID:  wid,
		Status:      StatusStarted,
	}

	_, err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	id2 := randomID()
	wid2 := randomID()
	createWorkflow(t, store, wid2)
	es2 := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": stepOne,
			"step2": stepTwo,
		},
		ExecutionID: id2,
		WorkflowID:  wid2,
		Status:      StatusStarted,
	}

	_, err = store.Add(tests.Context(t), &es2)
	require.NoError(t, err)

	id = randomID()
	esTwo := WorkflowExecution{
		ExecutionID: id,
		Status:      StatusCompleted,
		Steps:       map[string]*WorkflowExecutionStep{},
	}
	_, err = store.Add(tests.Context(t), &esTwo)
	require.NoError(t, err)

	states, err := store.GetUnfinished(tests.Context(t), wid, 0, 100)
	require.NoError(t, err)

	assert.Len(t, states, 1)
	// Zero out the completedAt timestamp
	states[0].CreatedAt = nil
	assert.Equal(t, es, states[0])
}
