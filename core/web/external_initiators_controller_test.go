package web_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalInitiatorsController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationEthereumDisabled(t)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	db := app.GetStore().DB

	eiFoo := cltest.MustInsertExternalInitiatorWithOpts(t, db, cltest.ExternalInitiatorOpts{
		NamePrefix:    "foo",
		URL:           cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingToken: "outgoing_token",
	})
	eiBar := cltest.MustInsertExternalInitiatorWithOpts(t, db, cltest.ExternalInitiatorOpts{NamePrefix: "bar"})

	resp, cleanup := client.Get("/v2/external_initiators?size=x")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)

	resp, cleanup = client.Get("/v2/external_initiators?size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 2, metaCount)

	var links jsonapi.Links
	eis := []presenters.ExternalInitiatorResource{}
	err = web.ParsePaginatedResponse(body, &eis, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, eis, 1)
	assert.Equal(t, fmt.Sprintf("%d", eiBar.ID), eis[0].ID)
	assert.Equal(t, eiBar.Name, eis[0].Name)
	assert.Nil(t, eis[0].URL)
	assert.Equal(t, eiBar.AccessKey, eis[0].AccessKey)
	assert.Equal(t, eiBar.OutgoingToken, eis[0].OutgoingToken)

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	eis = []presenters.ExternalInitiatorResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &eis, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])

	assert.Len(t, eis, 1)
	assert.Equal(t, fmt.Sprintf("%d", eiFoo.ID), eis[0].ID)
	assert.Equal(t, eiFoo.Name, eis[0].Name)
	assert.Equal(t, eiFoo.URL.String(), eis[0].URL.String())
	assert.Equal(t, eiFoo.AccessKey, eis[0].AccessKey)
	assert.Equal(t, eiFoo.OutgoingToken, eis[0].OutgoingToken)
}

func TestExternalInitiatorsController_Create_success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationEthereumDisabled(t)
	t.Cleanup(cleanup)
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

	app, cleanup := cltest.NewApplicationEthereumDisabled(t)
	t.Cleanup(cleanup)
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

	app, cleanup := cltest.NewApplicationEthereumDisabled(t)
	t.Cleanup(cleanup)
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

	app, cleanup := cltest.NewApplicationEthereumDisabled(t)
	t.Cleanup(cleanup)
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

	app, cleanup := cltest.NewApplicationEthereumDisabled(t)
	t.Cleanup(cleanup)
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
