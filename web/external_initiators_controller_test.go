package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalInitiatorsController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators", nil)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 201)
}

func TestExternalInitiatorsController_Delete(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Delete("/v2/external_initiators")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode)
}

func TestExternalInitiatorsController_DeleteNotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	err := app.GetStore().CreateExternalInitiator(&models.ExternalInitiator{
		AccessKey: "abracadabra",
	})
	require.NoError(t, err)

	client := app.NewHTTPClient()

	resp, cleanup := client.Delete("/v2/external_initiators/abracadabra")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 204)
}
