package v2

import (
	"database/sql"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/stretchr/testify/require"
)

func Test_UpsertWorkflowSpec(t *testing.T) {
	t.Parallel()
	t.Run("inserts new spec", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			WorkflowTag:   "workflowTag",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		_, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)

		// Verify the record exists in the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs_v2 WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_tag = $3`, spec.WorkflowOwner, spec.WorkflowName, spec.WorkflowTag)
		require.NoError(t, err)
		require.Equal(t, spec.Workflow, dbSpec.Workflow)
	})

	t.Run("updates existing spec", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			WorkflowTag:   "workflowTag",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		_, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)

		// Update the status
		spec.Status = job.WorkflowSpecStatusPaused

		_, err = orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)

		// Verify the record is updated in the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs_v2 WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_tag = $3`, spec.WorkflowOwner, spec.WorkflowName, spec.WorkflowTag)
		require.NoError(t, err)
		require.Equal(t, spec.Config, dbSpec.Config)
		require.Equal(t, spec.Status, dbSpec.Status)
	})

	t.Run("sets CreatedAt to current time if zero value is provided", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-zero-time",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			WorkflowTag:   "workflowTag",
			Status:        job.WorkflowSpecStatusActive,
			SpecType:      job.WASMFile,
			// Intentionally leaving CreatedAt as zero
		}

		startTime := time.Now().UTC()
		_, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)

		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs_v2 WHERE workflow_id = $1`, spec.WorkflowID)
		require.NoError(t, err)

		// Assert CreatedAt was populated correctly in the DB
		require.False(t, dbSpec.CreatedAt.IsZero(), "CreatedAt should not be the zero time")
		require.True(t, dbSpec.CreatedAt.After(startTime) || dbSpec.CreatedAt.Equal(startTime), "CreatedAt should be >= start time")
		require.True(t, dbSpec.CreatedAt.Before(time.Now().UTC()), "CreatedAt should be <= current time")
	})

	t.Run("workflow is unique by workflow ID", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		WFID1 := "cid-123"
		WFID2 := "cid-456"
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    WFID1,
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			WorkflowTag:   "workflowTag",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		_, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)

		// Verify the record exists in the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs_v2 WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_id = $3 AND workflow_tag = $4`, spec.WorkflowOwner, spec.WorkflowName, WFID1, spec.WorkflowTag)
		require.NoError(t, err)
		require.Equal(t, WFID1, dbSpec.WorkflowID)

		// Create another entry with a different ID
		spec.WorkflowID = WFID2
		_, err = orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)

		// Verify the original record is still there
		var dbSpec2 job.WorkflowSpec
		err = db.Get(&dbSpec2, `SELECT * FROM workflow_specs_v2 WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_id = $3 AND workflow_tag = $4`, spec.WorkflowOwner, spec.WorkflowName, WFID1, spec.WorkflowTag)
		require.NoError(t, err)
		require.Equal(t, WFID1, dbSpec2.WorkflowID)

		// Verify the new record is there
		var dbSpec3 job.WorkflowSpec
		err = db.Get(&dbSpec3, `SELECT * FROM workflow_specs_v2 WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_id = $3 AND workflow_tag = $4`, spec.WorkflowOwner, spec.WorkflowName, WFID2, spec.WorkflowTag)
		require.NoError(t, err)
		require.Equal(t, WFID2, dbSpec3.WorkflowID)
	})
}

func Test_DeleteWorkflowSpec(t *testing.T) {
	t.Parallel()
	t.Run("deletes a workflow spec by ID", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			WorkflowTag:   "workflowTag",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		id, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)
		require.NotZero(t, id)

		err = orm.DeleteWorkflowSpec(ctx, spec.WorkflowID)
		require.NoError(t, err)

		// Verify the record is deleted from the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs_v2 WHERE id = $1`, id)
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("fails if no workflow spec exists", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		err := orm.DeleteWorkflowSpec(ctx, "non-existent-workflow-id")
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func Test_DeleteWorkflowSpecs(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	specs := []struct {
		id   string
		name string
	}{
		{"wf-1", "Workflow 1"},
		{"wf-2", "Workflow 2"},
		{"wf-3", "Workflow 3"},
	}
	for _, s := range specs {
		_, err := orm.UpsertWorkflowSpec(ctx, &job.WorkflowSpec{
			Workflow:      "binary",
			Config:        "config",
			WorkflowID:    s.id,
			WorkflowOwner: "owner",
			WorkflowName:  s.name,
			Status:        job.WorkflowSpecStatusActive,
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		})
		require.NoError(t, err)
	}

	err := orm.DeleteWorkflowSpecs(ctx, []string{"wf-1", "wf-2"})
	require.NoError(t, err)

	_, err = orm.GetWorkflowSpec(ctx, "wf-1")
	require.ErrorIs(t, err, sql.ErrNoRows)
	_, err = orm.GetWorkflowSpec(ctx, "wf-2")
	require.ErrorIs(t, err, sql.ErrNoRows)

	got, err := orm.GetWorkflowSpec(ctx, "wf-3")
	require.NoError(t, err)
	require.Equal(t, "Workflow 3", got.WorkflowName)

	// empty slice is a no-op
	require.NoError(t, orm.DeleteWorkflowSpecs(ctx, []string{}))
}

func Test_GetWorkflowSpec(t *testing.T) {
	t.Parallel()
	t.Run("gets a workflow spec by ID", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			WorkflowTag:   "workflowTag",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		id, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)
		require.NotZero(t, id)

		dbSpec, err := orm.GetWorkflowSpec(ctx, spec.WorkflowID)
		require.NoError(t, err)
		require.Equal(t, spec.Workflow, dbSpec.Workflow)

		err = orm.DeleteWorkflowSpec(ctx, spec.WorkflowID)
		require.NoError(t, err)
	})

	t.Run("fails if no workflow spec exists", func(t *testing.T) {
		t.Parallel()
		db := pgtest.NewSqlxDB(t)
		ctx := t.Context()
		lggr := logger.TestLogger(t)
		orm := &orm{ds: db, lggr: lggr}

		dbSpec, err := orm.GetWorkflowSpec(ctx, "inexistent-workflow-id")
		require.Error(t, err)
		require.Nil(t, dbSpec)
	})
}
