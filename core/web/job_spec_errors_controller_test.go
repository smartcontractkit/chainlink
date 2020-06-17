package web_test

import (
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobSpecErrorsController_Delete(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	j := cltest.NewJob()
	app.Store.CreateJob(&j)

	description := "job spec error description"

	err := app.Store.UpsertErrorFor(j.ID, description)
	assert.NoError(t, err)

	jse, err := app.Store.FindJobSpecError(j.ID, description)
	assert.NoError(t, err)

	resp, cleanup := client.Delete("/v2/job_spec_errors/" + jse.ID.String())
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusNoContent)

	_, err = app.Store.FindJobSpecError(jse.JobSpecID, jse.Description)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestJobSpecErrorsController_Delete_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Delete("/v2/job_spec_errors/190AE4CE-40B6-4D60-A3DA-061C5ACD32D0")
	defer cleanup()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be not found")
}

func TestJobSpecErrorsController_Delete_InvalidUuid(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/specs/garbage")
	defer cleanup()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "Response should be unprocessable entity")
}

func TestJobSpecErrorsController_Delete_Unauthenticated(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	defer cleanup()

	resp, err := http.Get(app.Server.URL + "/v2/specs/garbage")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response should be forbidden")
}
