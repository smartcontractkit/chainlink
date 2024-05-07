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
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func randomID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func Test_StoreDB(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	store := &DBStore{db: db, clock: clockwork.NewFakeClock()}

	id := randomID()
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": {
				ExecutionID: id,
				Ref:         "step1",
				Status:      "completed",
			},
			"step2": {
				ExecutionID: id,
				Ref:         "step2",
				Status:      "started",
			},
		},
		ExecutionID: id,
		Status:      "started",
	}

	err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	gotEs, err := store.Get(tests.Context(t), es.ExecutionID)
	// Zero out the created at timestamp; this isn't present on `es`
	// but is added by the db store.
	gotEs.CreatedAt = nil
	require.NoError(t, err)
	assert.Equal(t, es, gotEs)
}

func Test_StoreDB_DuplicateEntry(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	store := &DBStore{db: db, clock: clockwork.NewFakeClock()}

	id := randomID()
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": {
				ExecutionID: id,
				Ref:         "step1",
				Status:      "completed",
			},
			"step2": {
				ExecutionID: id,
				Ref:         "step2",
				Status:      "started",
			},
		},
		ExecutionID: id,
		Status:      "started",
	}

	err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	err = store.Add(tests.Context(t), &es)
	assert.ErrorContains(t, err, "duplicate key value violates")
}

func Test_StoreDB_UpdateStatus(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	store := &DBStore{db: db, clock: clockwork.NewFakeClock()}

	id := randomID()
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": {
				ExecutionID: id,
				Ref:         "step1",
				Status:      "completed",
			},
			"step2": {
				ExecutionID: id,
				Ref:         "step2",
				Status:      "started",
			},
		},
		ExecutionID: id,
		Status:      "started",
	}

	err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	completedStatus := "completed"
	err = store.UpdateStatus(tests.Context(t), es.ExecutionID, "completed")
	require.NoError(t, err)

	gotEs, err := store.Get(tests.Context(t), es.ExecutionID)
	require.NoError(t, err)

	assert.Equal(t, gotEs.Status, completedStatus)
}

func Test_StoreDB_UpdateStep(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	store := &DBStore{db: db, clock: clockwork.NewFakeClock()}

	id := randomID()
	stepOne := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step1",
		Status:      "completed",
	}
	stepTwo := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step2",
		Status:      "started",
	}
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": stepOne,
			"step2": stepTwo,
		},
		ExecutionID: id,
		Status:      "started",
	}

	err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	stepOne.Status = "completed"
	nm, err := values.NewMap(map[string]any{"hello": "world"})
	require.NoError(t, err)

	stepOne.Inputs = nm
	stepOne.Outputs = &StepOutput{Err: errors.New("some error")}

	es, err = store.UpsertStep(tests.Context(t), stepOne)
	require.NoError(t, err)

	gotStep := es.Steps[stepOne.Ref]
	assert.Equal(t, stepOne, gotStep)

	stepTwo.Outputs = &StepOutput{Value: nm}
	es, err = store.UpsertStep(tests.Context(t), stepTwo)
	require.NoError(t, err)

	gotStep = es.Steps[stepTwo.Ref]
	assert.Equal(t, stepTwo, gotStep)
}

func Test_StoreDB_GetUnfinishedSteps(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	store := &DBStore{db: db, clock: clockwork.NewFakeClock()}

	id := randomID()
	stepOne := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step1",
		Status:      "completed",
	}
	stepTwo := &WorkflowExecutionStep{
		ExecutionID: id,
		Ref:         "step2",
		Status:      "started",
	}
	es := WorkflowExecution{
		Steps: map[string]*WorkflowExecutionStep{
			"step1": stepOne,
			"step2": stepTwo,
		},
		ExecutionID: id,
		Status:      "started",
	}

	err := store.Add(tests.Context(t), &es)
	require.NoError(t, err)

	id = randomID()
	esTwo := WorkflowExecution{
		ExecutionID: id,
		Status:      "completed",
		Steps:       map[string]*WorkflowExecutionStep{},
	}
	err = store.Add(tests.Context(t), &esTwo)
	require.NoError(t, err)

	states, err := store.GetUnfinished(tests.Context(t), 0, 100)
	require.NoError(t, err)

	assert.Len(t, states, 1)
	// Zero out the completedAt timestamp
	states[0].CreatedAt = nil
	assert.Equal(t, es, states[0])
}
