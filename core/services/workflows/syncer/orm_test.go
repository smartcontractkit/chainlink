package syncer

import (
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowArtifactsORM_GetAndUpdate(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	giveURL := "https://example.com"
	giveBytes, err := crypto.Keccak256([]byte(giveURL))
	require.NoError(t, err)
	giveHash := hex.EncodeToString(giveBytes)
	giveContent := "some contents"

	gotID, err := orm.Create(ctx, giveURL, giveHash, giveContent)
	require.NoError(t, err)

	url, err := orm.GetSecretsURLByID(ctx, gotID)
	require.NoError(t, err)
	assert.Equal(t, giveURL, url)

	contents, err := orm.GetContents(ctx, giveURL)
	require.NoError(t, err)
	assert.Equal(t, "some contents", contents)

	contents, err = orm.GetContentsByHash(ctx, giveHash)
	require.NoError(t, err)
	assert.Equal(t, "some contents", contents)

	_, err = orm.Update(ctx, giveHash, "new contents")
	require.NoError(t, err)

	contents, err = orm.GetContents(ctx, giveURL)
	require.NoError(t, err)
	assert.Equal(t, "new contents", contents)

	contents, err = orm.GetContentsByHash(ctx, giveHash)
	require.NoError(t, err)
	assert.Equal(t, "new contents", contents)
}

func Test_UpsertWorkflowSpec(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	t.Run("inserts new spec", func(t *testing.T) {
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
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
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs WHERE workflow_owner = $1 AND workflow_name = $2`, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
		require.Equal(t, spec.Workflow, dbSpec.Workflow)
	})

	t.Run("updates existing spec", func(t *testing.T) {
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
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
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs WHERE workflow_owner = $1 AND workflow_name = $2`, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
		require.Equal(t, spec.Config, dbSpec.Config)
		require.Equal(t, spec.Status, dbSpec.Status)
	})
}

func Test_DeleteWorkflowSpec(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	t.Run("deletes a workflow spec", func(t *testing.T) {
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		id, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)
		require.NotZero(t, id)

		err = orm.DeleteWorkflowSpec(ctx, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)

		// Verify the record is deleted from the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs WHERE id = $1`, id)
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("fails if no workflow spec exists", func(t *testing.T) {
		err := orm.DeleteWorkflowSpec(ctx, "owner-123", "Test Workflow")
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func Test_GetWorkflowSpec(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	t.Run("gets a workflow spec", func(t *testing.T) {
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    "cid-123",
			WorkflowOwner: "owner-123",
			WorkflowName:  "Test Workflow",
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary",
			ConfigURL:     "http://example.com/config",
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		id, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)
		require.NotZero(t, id)

		dbSpec, err := orm.GetWorkflowSpec(ctx, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
		require.Equal(t, spec.Workflow, dbSpec.Workflow)

		err = orm.DeleteWorkflowSpec(ctx, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
	})

	t.Run("fails if no workflow spec exists", func(t *testing.T) {
		dbSpec, err := orm.GetWorkflowSpec(ctx, "owner-123", "Test Workflow")
		require.Error(t, err)
		require.Nil(t, dbSpec)
	})
}
