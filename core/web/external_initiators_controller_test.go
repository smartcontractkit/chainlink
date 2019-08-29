package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalInitiatorsController_Create_success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"name":"bitcoin","url":"http://without.a.name"}`),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 201)
	ei := &presenters.ExternalInitiatorAuthentication{}
	err := cltest.ParseJSONAPIResponse(t, resp, ei)
	require.NoError(t, err)

	assert.Equal(t, "bitcoin", ei.Name)
	assert.Equal(t, "http://without.a.name", ei.URL.String())
	assert.NotEmpty(t, ei.AccessKey)
	assert.NotEmpty(t, ei.Secret)
	assert.NotEmpty(t, ei.OutgoingToken)
	assert.NotEmpty(t, ei.OutgoingSecret)
}

func TestExternalInitiatorsController_Create_invalid(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 400)
}

func TestExternalInitiatorsController_Delete(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Delete("/v2/external_initiators")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode)
}

func TestExternalInitiatorsController_DeleteNotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
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
