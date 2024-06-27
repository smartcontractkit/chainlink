package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestPipelineJobSpecErrorsController_Delete_2(t *testing.T) {
	app, client, _, jID, _, _ := setupJobSpecsControllerTestsWithJobs(t)

	description := "job spec error description"

	require.NoError(t, app.JobORM().RecordError(jID, description))

	// FindJob -> find error
	j, err := app.JobORM().FindJob(testutils.Context(t), jID)
	require.NoError(t, err)
	t.Log(j.JobSpecErrors)
	require.GreaterOrEqual(t, len(j.JobSpecErrors), 1) // second 'got nil head' error may have occured also
	var id int64 = -1
	for i := range j.JobSpecErrors {
		jse := j.JobSpecErrors[i]
		if jse.Description == description {
			id = jse.ID
			break
		}
	}
	require.NotEqual(t, -1, id, "error not found")

	resp, cleanup := client.Delete(fmt.Sprintf("/v2/pipeline/job_spec_errors/%v", id))
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusNoContent)

	// FindJob -> error is gone
	j, err = app.JobORM().FindJob(testutils.Context(t), j.ID)
	require.NoError(t, err)
	for i := range j.JobSpecErrors {
		jse := j.JobSpecErrors[i]
		require.NotEqual(t, id, jse.ID)
	}
}

func TestPipelineJobSpecErrorsController_Delete_NotFound(t *testing.T) {
	_, client, _, _, _, _ := setupJobSpecsControllerTestsWithJobs(t)

	resp, cleanup := client.Delete("/v2/pipeline/job_spec_errors/1")
	defer cleanup()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be not found")
}
