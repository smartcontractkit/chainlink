package web_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateBridgeType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		request     bridges.BridgeTypeRequest
		want        error
	}{
		{
			"no adapter name",
			bridges.BridgeTypeRequest{
				URL: cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("No name specified"),
		},
		{
			"invalid adapter name",
			bridges.BridgeTypeRequest{
				Name: "invalid/adapter",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("task type validation: name invalid/adapter contains invalid characters"),
		},
		{
			"invalid with blank url",
			bridges.BridgeTypeRequest{
				Name: "validadaptername",
				URL:  cltest.WebURL(t, ""),
			},
			models.NewJSONAPIErrorsWith("URL must be present"),
		},
		{
			"valid url",
			bridges.BridgeTypeRequest{
				Name: "adapterwithvalidurl",
				URL:  cltest.WebURL(t, "//denergy"),
			},
			nil,
		},
		{
			"valid docker url",
			bridges.BridgeTypeRequest{
				Name: "adapterwithdockerurl",
				URL:  cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
			},
			nil,
		},
		{
			"valid MinimumContractPayment positive",
			bridges.BridgeTypeRequest{
				Name:                   "adapterwithdockerurl",
				URL:                    cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
				MinimumContractPayment: assets.NewLinkFromJuels(1),
			},
			nil,
		},
		{
			"invalid MinimumContractPayment negative",
			bridges.BridgeTypeRequest{
				Name:                   "adapterwithdockerurl",
				URL:                    cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
				MinimumContractPayment: assets.NewLinkFromJuels(-1),
			},
			models.NewJSONAPIErrorsWith("MinimumContractPayment must be positive"),
		},
		{
			"existing core adapter (no longer fails since core adapters no longer exist)",
			bridges.BridgeTypeRequest{
				Name: "ethtx",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			nil,
		},
		{
			"new external adapter",
			bridges.BridgeTypeRequest{
				Name: "gdaxprice",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			nil,
		}}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := web.ValidateBridgeType(&test.request)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateBridgeNotExist(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(true)
	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)

	// Create a duplicate
	bt := bridges.BridgeType{}
	bt.Name = bridges.MustParseBridgeName("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, orm.CreateBridgeType(&bt))

	newBridge := bridges.BridgeTypeRequest{
		Name: "solargridreporting",
	}
	expected := models.NewJSONAPIErrorsWith("Bridge Type solargridreporting already exists")
	result := web.ValidateBridgeTypeNotExist(&newBridge, orm)
	assert.Equal(t, expected, result)
}

func BenchmarkBridgeTypesController_Index(b *testing.B) {
	app := cltest.NewApplication(b)
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/bridge_types")
		defer cleanup()
		assert.Equal(b, http.StatusOK, resp.StatusCode, "Response should be successful")
	}
}

func TestBridgeTypesController_Index(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	bt, err := setupBridgeControllerIndex(t, app.BridgeORM())
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/bridge_types?size=x")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)

	resp, cleanup = client.Get("/v2/bridge_types?size=1")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	resources := []presenters.BridgeResource{}

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &resources, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, resources, 1)
	assert.Equal(t, bt[0].Name.String(), resources[0].Name, "should have the same Name")
	assert.Equal(t, bt[0].URL.String(), resources[0].URL, "should have the same URL")
	assert.Equal(t, bt[0].Confirmations, resources[0].Confirmations, "should have the same Confirmations")

	resp, cleanup = client.Get(links["next"].Href)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resources = []presenters.BridgeResource{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(t, resp), &resources, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])
	assert.Len(t, resources, 1)
	assert.Equal(t, bt[1].Name.String(), resources[0].Name, "should have the same Name")
	assert.Equal(t, bt[1].URL.String(), resources[0].URL, "should have the same URL")
	assert.Equal(t, bt[1].Confirmations, resources[0].Confirmations, "should have the same Confirmations")
}

// cannot randomize bridge names here since they are ordered by name on the API
// leading in random order for assertion...
func setupBridgeControllerIndex(t testing.TB, orm bridges.ORM) ([]*bridges.BridgeType, error) {
	bt1 := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName("indexbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := orm.CreateBridgeType(bt1)
	if err != nil {
		return nil, err
	}

	bt2 := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName("indexbridges2"),
		URL:           cltest.WebURL(t, "https://testing.com/tari"),
		Confirmations: 0,
	}
	err = orm.CreateBridgeType(bt2)
	return []*bridges.BridgeType{bt1, bt2}, err
}

func TestBridgeTypesController_Create_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBuffer(cltest.MustReadFile(t, "../testdata/apiresponses/create_random_number_bridge_type.json")),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	respJSON := cltest.ParseJSON(t, resp.Body)
	btName := respJSON.Get("data.attributes.name").String()

	assert.NotEmpty(t, respJSON.Get("data.attributes.incomingToken").String())
	assert.NotEmpty(t, respJSON.Get("data.attributes.outgoingToken").String())

	bt, err := app.BridgeORM().FindBridge(bridges.MustParseBridgeName(btName))
	assert.NoError(t, err)
	assert.Equal(t, "randomnumber", bt.Name.String())
	assert.Equal(t, uint32(10), bt.Confirmations)
	assert.Equal(t, "https://example.com/randomNumber", bt.URL.String())
	assert.Equal(t, assets.NewLinkFromJuels(100), bt.MinimumContractPayment)
	assert.NotEmpty(t, bt.OutgoingToken)
}

func TestBridgeTypesController_Update_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	bridgeName := testutils.RandomizeName("BRidgea")
	bt := &bridges.BridgeType{
		Name: bridges.MustParseBridgeName(bridgeName),
		URL:  cltest.WebURL(t, "http://mybridge"),
	}
	require.NoError(t, app.BridgeORM().CreateBridgeType(bt))

	body := fmt.Sprintf(`{"name": "%s","url":"http://yourbridge"}`, bridgeName)
	ud := bytes.NewBuffer([]byte(body))
	resp, cleanup := client.Patch("/v2/bridge_types/"+bridgeName, ud)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	ubt, err := app.BridgeORM().FindBridge(bt.Name)
	assert.NoError(t, err)
	assert.Equal(t, cltest.WebURL(t, "http://yourbridge"), ubt.URL)
}

func TestBridgeController_Show(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	bt := &bridges.BridgeType{
		Name:          bridges.MustParseBridgeName(testutils.RandomizeName("showbridge")),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	require.NoError(t, app.BridgeORM().CreateBridgeType(bt))

	resp, cleanup := client.Get("/v2/bridge_types/" + bt.Name.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be successful")

	var resource presenters.BridgeResource
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &resource))
	assert.Equal(t, bt.Name.String(), resource.Name, "should have the same name")
	assert.Equal(t, bt.URL.String(), resource.URL, "should have the same URL")
	assert.Equal(t, bt.Confirmations, resource.Confirmations, "should have the same Confirmations")

	resp, cleanup = client.Get("/v2/bridge_types/nosuchbridge")
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be 404")
}

func TestBridgeTypesController_Create_AdapterExistsError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBuffer(cltest.MustReadFile(t, "../testdata/apiresponses/existing_core_adapter.json")),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}

func TestBridgeTypesController_Create_BindJSONError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString("}"),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestBridgeTypesController_Create_DatabaseError(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}
