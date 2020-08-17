package web_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalInitiatorsController_Create_success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"name":"bitcoin","url":"http://without.a.name"}`),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusCreated)
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

func TestExternalInitiatorsController_Create_without_URL(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"name":"no-url"}`),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 201)
	ei := &presenters.ExternalInitiatorAuthentication{}
	err := cltest.ParseJSONAPIResponse(t, resp, ei)
	require.NoError(t, err)

	assert.Equal(t, "no-url", ei.Name)
	assert.Equal(t, "", ei.URL.String())
	assert.NotEmpty(t, ei.AccessKey)
	assert.NotEmpty(t, ei.Secret)
	assert.NotEmpty(t, ei.OutgoingToken)
	assert.NotEmpty(t, ei.OutgoingSecret)
}

func TestExternalInitiatorsController_Create_invalid(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}

func TestExternalInitiatorsController_Delete(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	exi := models.ExternalInitiator{
		Name: "abracadabra",
	}
	err := app.GetStore().CreateExternalInitiator(&exi)
	require.NoError(t, err)

	client := app.NewHTTPClient()

	resp, cleanup := client.Delete("/v2/external_initiators/" + exi.Name)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusNoContent)
}

func TestExternalInitiatorsController_DeleteNotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	tests := []struct {
		Name string
		URL  string
	}{
		{
			Name: "No external initiator specified",
			URL:  "/v2/external_initiators",
		},
		{
			Name: "Unknown initiator",
			URL:  "/v2/external_initiators/not-exist",
		},
	}

	for _, test := range tests {
		t.Log(test.Name)
		resp, cleanup := client.Delete(test.URL)
		defer cleanup()
		assert.Equal(t, http.StatusText(http.StatusNotFound), http.StatusText(resp.StatusCode))
	}
}
