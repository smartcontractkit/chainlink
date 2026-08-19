package artifacts

import (
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

func TestWorkflowArtifactsORM_GetAndUpdate(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	giveURL := "https://example.com/" + uuid.New().String()[:8]
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
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	owner := "owner-" + uuid.New().String()[:8]
	name := "name-" + uuid.New().String()[:8]
	cid := "cid-" + uuid.New().String()[:8]

	t.Run("inserts new spec", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
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

	t.Run("updates existing spec", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
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
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	owner := "owner-" + uuid.New().String()[:8]
	name := "name-" + uuid.New().String()[:8]
	cid := "cid-" + uuid.New().String()[:8]

	t.Run("deletes a workflow spec", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
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

	t.Run("fails if no workflow spec exists", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		err := orm.DeleteWorkflowSpec(ctx, owner, name)
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func Test_GetWorkflowSpec(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	owner := "owner-" + uuid.New().String()[:8]
	name := "name-" + uuid.New().String()[:8]
	cid := "cid-" + uuid.New().String()[:8]

	t.Run("gets a workflow spec", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
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

	t.Run("fails if no workflow spec exists", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		dbSpec, err := orm.GetWorkflowSpec(ctx, owner, name)
		require.Error(t, err)
		require.Nil(t, dbSpec)
	})
}

func Test_GetWorkflowSpecByID(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	owner := "owner-" + uuid.New().String()[:8]
	name := "name-" + uuid.New().String()[:8]
	cid := "cid-" + uuid.New().String()[:8]

	t.Run("gets a workflow spec by ID", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		id, err := orm.UpsertWorkflowSpec(ctx, spec)
		require.NoError(t, err)
		require.NotZero(t, id)

		dbSpec, err := orm.GetWorkflowSpecByID(ctx, spec.WorkflowID)
		require.NoError(t, err)
		require.Equal(t, spec.Workflow, dbSpec.Workflow)

		err = orm.DeleteWorkflowSpec(ctx, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
	})

	t.Run("fails if no workflow spec exists", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		dbSpec, err := orm.GetWorkflowSpecByID(ctx, cid)
		require.Error(t, err)
		require.Nil(t, dbSpec)
	})
}

func Test_GetContentsByWorkflowID(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "wf-" + uuid.New().String()[:8]
	workflowOwner := "owner-" + uuid.New().String()[:8]
	workflowName := "name-" + uuid.New().String()[:8]

	// workflow_id is missing
	_, _, err := orm.GetContentsByWorkflowID(ctx, "doesnt-exist")
	require.ErrorContains(t, err, "no rows in result set")

	// secrets_id is nil; should return EmptySecrets
	_, err = orm.UpsertWorkflowSpec(ctx, &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	_, _, err = orm.GetContentsByWorkflowID(ctx, workflowID)
	require.ErrorIs(t, err, ErrEmptySecrets)

	// retrieves the artifact if provided
	giveURL := "https://example.com/" + uuid.New().String()[:8]
	giveBytes, err := crypto.Keccak256([]byte(giveURL))
	require.NoError(t, err)
	giveHash := hex.EncodeToString(giveBytes)
	giveContent := "some contents"

	secretsID, err := orm.Create(ctx, giveURL, giveHash, giveContent)
	require.NoError(t, err)

	// Use a new workflow ID to avoid conflicts
	workflowID2 := "wf-" + uuid.New().String()[:8]
	_, err = orm.UpsertWorkflowSpec(ctx, &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID2,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName + "-2",
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)
	_, err = orm.GetWorkflowSpec(ctx, workflowOwner, workflowName+"-2")
	require.NoError(t, err)

	gotHash, gotContent, err := orm.GetContentsByWorkflowID(ctx, workflowID2)
	require.NoError(t, err)
	assert.Equal(t, giveHash, gotHash)
	assert.Equal(t, giveContent, gotContent)
}

func Test_GetContentsByWorkflowID_SecretsProvidedButEmpty(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	// workflow_id is missing
	_, _, err := orm.GetContentsByWorkflowID(ctx, "doesnt-exist")
	require.ErrorContains(t, err, "no rows in result set")

	// secrets_id is nil; should return EmptySecrets
	workflowID := "wf-" + uuid.New().String()[:8]
	workflowOwner := "owner-" + uuid.New().String()[:8]
	workflowName := "name-" + uuid.New().String()[:8]
	giveURL := "https://example.com/" + uuid.New().String()[:8]
	giveBytes, err := crypto.Keccak256([]byte(giveURL))
	require.NoError(t, err)
	giveHash := hex.EncodeToString(giveBytes)
	giveContent := ""
	_, err = orm.UpsertWorkflowSpecWithSecrets(ctx, &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	}, giveURL, giveHash, giveContent)
	require.NoError(t, err)

	_, _, err = orm.GetContentsByWorkflowID(ctx, workflowID)
	require.ErrorIs(t, err, ErrEmptySecrets)
}

func Test_UpsertWorkflowSpecWithSecrets(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	owner := "owner-" + uuid.New().String()[:8]
	name := "name-" + uuid.New().String()[:8]
	cid := "cid-" + uuid.New().String()[:8]

	t.Run("inserts new spec and new secrets", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		giveURL := "https://example.com/" + uuid.New().String()[:8]
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)
		giveHash := hex.EncodeToString(giveBytes)
		giveContent := "some contents"

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		_, err = orm.UpsertWorkflowSpecWithSecrets(ctx, spec, giveURL, giveHash, giveContent)
		require.NoError(t, err)

		// Verify the record exists in the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs WHERE workflow_owner = $1 AND workflow_name = $2`, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
		require.Equal(t, spec.Workflow, dbSpec.Workflow)

		// Verify the secrets exists in the database
		contents, err := orm.GetContents(ctx, giveURL)
		require.NoError(t, err)
		require.Equal(t, giveContent, contents)
	})

	t.Run("updates existing spec and secrets", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		giveURL := "https://example.com/" + uuid.New().String()[:8]
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)
		giveHash := hex.EncodeToString(giveBytes)
		giveContent := "some contents"

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		_, err = orm.UpsertWorkflowSpecWithSecrets(ctx, spec, giveURL, giveHash, giveContent)
		require.NoError(t, err)

		// Update the status
		spec.Status = job.WorkflowSpecStatusPaused

		_, err = orm.UpsertWorkflowSpecWithSecrets(ctx, spec, giveURL, giveHash, "new contents")
		require.NoError(t, err)

		// Verify the record is updated in the database
		var dbSpec job.WorkflowSpec
		err = db.Get(&dbSpec, `SELECT * FROM workflow_specs WHERE workflow_owner = $1 AND workflow_name = $2`, spec.WorkflowOwner, spec.WorkflowName)
		require.NoError(t, err)
		require.Equal(t, spec.Config, dbSpec.Config)

		// Verify the secrets is updated in the database
		contents, err := orm.GetContents(ctx, giveURL)
		require.NoError(t, err)
		require.Equal(t, "new contents", contents)
	})

	t.Run("updates existing spec and secrets if spec has executions", func(t *testing.T) { //nolint:paralleltest // subtests share database setup
		giveURL := "https://example.com/" + uuid.New().String()[:8]
		giveBytes, err := crypto.Keccak256([]byte(giveURL))
		require.NoError(t, err)
		giveHash := hex.EncodeToString(giveBytes)
		giveContent := "some contents"

		spec := &job.WorkflowSpec{
			Workflow:      "test_workflow",
			Config:        "test_config",
			WorkflowID:    cid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        job.WorkflowSpecStatusActive,
			BinaryURL:     "http://example.com/binary/" + cid,
			ConfigURL:     "http://example.com/config/" + cid,
			CreatedAt:     time.Now(),
			SpecType:      job.WASMFile,
		}

		_, err = orm.UpsertWorkflowSpecWithSecrets(ctx, spec, giveURL, giveHash, giveContent)
		require.NoError(t, err)

		_, err = db.ExecContext(
			ctx,
			`INSERT INTO workflow_executions (id, workflow_id, status, created_at) VALUES ($1, $2, $3, $4)`,
			uuid.New().String(),
			cid,
			"started",
			time.Now(),
		)
		require.NoError(t, err)

		// Update the status
		spec.WorkflowID = "cid-456-" + uuid.New().String()[:8]

		_, err = orm.UpsertWorkflowSpecWithSecrets(ctx, spec, giveURL, giveHash, "new contents")
		require.NoError(t, err)
	})
}
