package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkBridgeTypesController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication(b)
	b.Cleanup(cleanup)
	setupJobSpecsControllerIndex(app)
	client := app.NewHTTPClient()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/specs")
		defer cleanup()
		assert.Equal(b, http.StatusOK, resp.StatusCode, "Response should be successful")
	}
}

func TestBridgeTypesController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	bt, err := setupBridgeControllerIndex(t, app.Store)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/specs?size=x")
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

func setupBridgeControllerIndex(t testing.TB, store *store.Store) ([]*models.BridgeType, error) {
	bt1 := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := store.CreateBridgeType(bt1)
	if err != nil {
		return nil, err
	}

	bt2 := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges2"),
		URL:           cltest.WebURL(t, "https://testing.com/tari"),
		Confirmations: 0,
	}
	err = store.CreateBridgeType(bt2)
	return []*models.BridgeType{bt1, bt2}, err
}

func TestBridgeTypesController_Create_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

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

	bt, err := app.Store.FindBridge(models.MustNewTaskType(btName))
	assert.NoError(t, err)
	assert.Equal(t, "randomnumber", bt.Name.String())
	assert.Equal(t, uint32(10), bt.Confirmations)
	assert.Equal(t, "https://example.com/randomNumber", bt.URL.String())
	assert.Equal(t, assets.NewLink(100), bt.MinimumContractPayment)
	assert.NotEmpty(t, bt.OutgoingToken)
}

func TestBridgeTypesController_Update_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	bt := &models.BridgeType{
		Name: models.MustNewTaskType("BRidgea"),
		URL:  cltest.WebURL(t, "http://mybridge"),
	}
	require.NoError(t, app.GetStore().CreateBridgeType(bt))

	ud := bytes.NewBuffer([]byte(`{"name": "BRidgea","url":"http://yourbridge"}`))
	resp, cleanup := client.Patch("/v2/bridge_types/bridgea", ud)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	ubt, err := app.Store.FindBridge(bt.Name)
	assert.NoError(t, err)
	assert.Equal(t, cltest.WebURL(t, "http://yourbridge"), ubt.URL)
}

func TestBridgeController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	bt := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	require.NoError(t, app.GetStore().CreateBridgeType(bt))

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

func TestBridgeController_Destroy(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	resp, cleanup := client.Delete("/v2/bridge_types/testingbridges1")
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be 404")

	bridgeJSON := cltest.MustReadFile(t, "../testdata/apiresponses/create_random_number_bridge_type.json")
	var bt models.BridgeType
	err := json.Unmarshal(bridgeJSON, &bt)
	assert.NoError(t, err)
	require.NoError(t, app.GetStore().CreateBridgeType(&bt))

	resp, cleanup = client.Delete("/v2/bridge_types/" + bt.Name.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be successful")

	resp, cleanup = client.Get("/v2/bridge_types/" + bt.Name.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be 404")

	require.NoError(t, app.GetStore().CreateBridgeType(&bt))

	// Create an FM and DR job using the bridge.
	js := cltest.NewJobWithWebInitiator()
	js.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.NoError(t, app.Store.CreateJob(&js))

	jsFM := cltest.NewJobWithFluxMonitorInitiatorWithBridge(bt.Name.String())
	jsFM.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	assert.NoError(t, app.Store.CreateJob(&jsFM))

	resp, cleanup = client.Delete("/v2/bridge_types/" + bt.Name.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Response should be 409")

	require.NoError(t, app.Store.ArchiveJob(js.ID))

	// Still fails because FM job using bridge.
	resp, cleanup = client.Delete("/v2/bridge_types/" + bt.Name.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Response should be 409")

	require.NoError(t, app.Store.ArchiveJob(jsFM.ID))

	// Succeeds because FM job is archived.
	resp, cleanup = client.Delete("/v2/bridge_types/" + bt.Name.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be 200")
}

func TestBridgeTypesController_Create_AdapterExistsError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBuffer(cltest.MustReadFile(t, "../testdata/apiresponses/existing_core_adapter.json")),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}

func TestBridgeTypesController_Create_BindJSONError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString("}"),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusUnprocessableEntity)
}

func TestBridgeTypesController_Create_DatabaseError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusBadRequest)
}
