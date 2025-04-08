package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateExternalInitiator(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := bridges.NewORM(db)

	url := cltest.WebURL(t, "https://a.web.url")

	//  Add duplicate
	exi := bridges.ExternalInitiator{
		Name: "duplicate",
		URL:  &url,
	}

	assert.NoError(t, orm.CreateExternalInitiator(testutils.Context(t), &exi))

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"basic", `{"name":"bitcoin","url":"https://test.url"}`, false},
		{"basic w/ underscore", `{"name":"bit_coin","url":"https://test.url"}`, false},
		{"basic w/ underscore in url", `{"name":"bitcoin","url":"https://chainlink_bit-coin_1.url"}`, false},
		{"missing url", `{"name":"missing_url"}`, false},
		{"duplicate name", `{"name":"duplicate","url":"https://test.url"}`, true},
		{"invalid name characters", `{"name":"<invalid>","url":"https://test.url"}`, true},
		{"missing name", `{"url":"https://test.url"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			var exr bridges.ExternalInitiatorRequest

			assert.NoError(t, json.Unmarshal([]byte(test.input), &exr))
			result := web.ValidateExternalInitiator(ctx, &exr, orm)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}

func TestExternalInitiatorsController_Index(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithConfig(t,
		configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		}))
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	db := app.GetDB()
	borm := bridges.NewORM(db)

	eiFoo := cltest.MustInsertExternalInitiatorWithOpts(t, borm, cltest.ExternalInitiatorOpts{
		NamePrefix:    "foo",
		URL:           cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingToken: "outgoing_token",
	})
	eiBar := cltest.MustInsertExternalInitiatorWithOpts(t, borm, cltest.ExternalInitiatorOpts{NamePrefix: "bar"})

	resp, cleanup := client.Get("/v2/external_initiators?size=x")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)

	resp, cleanup = client.Get("/v2/external_initiators?size=1")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	body := cltest.ParseResponseBody(t, resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	require.NoError(t, err)
	require.Equal(t, 2, metaCount)

	var links jsonapi.Links
	var eis []presenters.ExternalInitiatorResource
	err = web.ParsePaginatedResponse(body, &eis, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, eis, 1)
	assert.Equal(t, strconv.FormatInt(eiBar.ID, 10), eis[0].ID)
	assert.Equal(t, eiBar.Name, eis[0].Name)
	assert.Nil(t, eis[0].URL)
	assert.Equal(t, eiBar.AccessKey, eis[0].AccessKey)
	assert.Equal(t, eiBar.OutgoingToken, eis[0].OutgoingToken)

	resp, cleanup = client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	eis = []presenters.ExternalInitiatorResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &eis, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])

	assert.Len(t, eis, 1)
	assert.Equal(t, strconv.FormatInt(eiFoo.ID, 10), eis[0].ID)
	assert.Equal(t, eiFoo.Name, eis[0].Name)
	assert.Equal(t, eiFoo.URL.String(), eis[0].URL.String())
	assert.Equal(t, eiFoo.AccessKey, eis[0].AccessKey)
	assert.Equal(t, eiFoo.OutgoingToken, eis[0].OutgoingToken)
}

func TestExternalInitiatorsController_Create_success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithConfig(t,
		configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		}))
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"name":"bitcoin","url":"http://without.a.name"}`),
	)
	t.Cleanup(cleanup)
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

	app := cltest.NewApplicationWithConfig(t,
		configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		}))
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"name":"no-url"}`),
	)
	t.Cleanup(cleanup)
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

	app := cltest.NewApplicationWithConfig(t,
		configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		}))
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Post("/v2/external_initiators",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}

func TestExternalInitiatorsController_Delete(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithConfig(t,
		configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		}))
	require.NoError(t, app.Start(testutils.Context(t)))

	exi := bridges.ExternalInitiator{
		Name: "abracadabra",
	}
	err := app.BridgeORM().CreateExternalInitiator(testutils.Context(t), &exi)
	require.NoError(t, err)

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Delete("/v2/external_initiators/" + exi.Name)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusNoContent)
}

func TestExternalInitiatorsController_DeleteNotFound(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithConfig(t,
		configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		}))
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)

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
		t.Run(test.Name, func(t *testing.T) {
			resp, cleanup := client.Delete(test.URL)
			t.Cleanup(cleanup)
			assert.Equal(t, http.StatusText(http.StatusNotFound), http.StatusText(resp.StatusCode))
		})
	}
}
