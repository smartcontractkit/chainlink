package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkBridgeTypesController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	setupJobSpecsControllerIndex(app)
	client := app.NewHTTPClient()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/specs")
		defer cleanup()
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestBridgeTypesController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	bt, err := setupBridgeControllerIndex(app)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/specs?size=x")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)

	resp, cleanup = client.Get("/v2/bridge_types?size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	var links jsonapi.Links
	bridges := []models.BridgeType{}

	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &bridges, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, bridges, 1)
	assert.Equal(t, bt[0].Name, bridges[0].Name, "should have the same Name")
	assert.Equal(t, bt[0].URL.String(), bridges[0].URL.String(), "should have the same URL")
	assert.Equal(t, bt[0].Confirmations, bridges[0].Confirmations, "should have the same Confirmations")

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	bridges = []models.BridgeType{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &bridges, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])
	assert.Len(t, bridges, 1)
	assert.Equal(t, bt[1].Name, bridges[0].Name, "should have the same Name")
	assert.Equal(t, bt[1].URL.String(), bridges[0].URL.String(), "should have the same URL")
	assert.Equal(t, bt[1].Confirmations, bridges[0].Confirmations, "should have the same Confirmations")
}

func setupBridgeControllerIndex(app *cltest.TestApplication) ([]*models.BridgeType, error) {

	bt1 := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL("https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := app.AddAdapter(bt1)
	if err != nil {
		return nil, err
	}

	bt2 := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges2"),
		URL:           cltest.WebURL("https://testing.com/tari"),
		Confirmations: 0,
	}
	err = app.AddAdapter(bt2)

	return []*models.BridgeType{bt1, bt2}, err
}

func TestBridgeTypesController_Create_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBuffer(cltest.MustReadFile(t, "testdata/create_random_number_bridge_type.json")),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	respJSON := cltest.ParseJSON(resp.Body)
	btName := respJSON.Get("data.attributes.name").String()

	assert.NotEmpty(t, respJSON.Get("data.attributes.incomingToken").String())
	assert.NotEmpty(t, respJSON.Get("data.attributes.outgoingToken").String())

	bt, err := app.Store.FindBridge(btName)
	assert.NoError(t, err)
	assert.Equal(t, "randomnumber", bt.Name.String())
	assert.Equal(t, uint64(10), bt.Confirmations)
	assert.Equal(t, "https://example.com/randomNumber", bt.URL.String())
	assert.Equal(t, *assets.NewLink(100), bt.MinimumContractPayment)
	assert.NotEmpty(t, bt.IncomingToken)
	assert.NotEmpty(t, bt.OutgoingToken)
}

func TestBridgeTypesController_Update_Success(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := &models.BridgeType{
		Name: models.MustNewTaskType("bridgea"),
		URL:  cltest.WebURL("http://mybridge"),
	}
	assert.NoError(t, app.AddAdapter(bt))

	ud := bytes.NewBuffer([]byte(`{"url":"http://yourbridge"}`))
	resp, cleanup := client.Patch("/v2/bridge_types/bridgea", ud)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	ubt, err := app.Store.FindBridge(bt.Name.String())
	assert.NoError(t, err)
	assert.Equal(t, cltest.WebURL("http://yourbridge"), ubt.URL)
}

func TestBridgeController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	bt := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL("https://testing.com/bridges"),
		Confirmations: 0,
	}
	assert.NoError(t, app.AddAdapter(bt))

	resp, cleanup := client.Get("/v2/bridge_types/" + bt.Name.String())
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respBridge presenters.BridgeType
	require.NoError(t, cltest.ParseJSONAPIResponse(resp, &respBridge))
	assert.Equal(t, respBridge.Name, bt.Name, "should have the same schedule")
	assert.Equal(t, respBridge.URL.String(), bt.URL.String(), "should have the same URL")
	assert.Equal(t, respBridge.Confirmations, bt.Confirmations, "should have the same Confirmations")

	resp, cleanup = client.Get("/v2/bridge_types/nosuchbridge")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be 404")
}

func TestBridgeController_Destroy(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()
	resp, cleanup := client.Delete("/v2/bridge_types/testingbridges1")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be 404")

	bridgeJSON := cltest.MustReadFile(t, "testdata/create_random_number_bridge_type.json")
	var bt models.BridgeType
	err := json.Unmarshal(bridgeJSON, &bt)
	assert.NoError(t, err)
	assert.NoError(t, app.AddAdapter(&bt))

	resp, cleanup = client.Delete("/v2/bridge_types/" + bt.Name.String())
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	resp, cleanup = client.Get("/v2/bridge_types/" + bt.Name.String())
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be 404")

	assert.NoError(t, app.AddAdapter(&bt))

	js := cltest.NewJobWithWebInitiator()
	js.Tasks = []models.TaskSpec{models.TaskSpec{Type: bt.Name}}
	assert.NoError(t, app.Store.CreateJob(&js))

	resp, cleanup = client.Delete("/v2/bridge_types/" + bt.Name.String())
	defer cleanup()
	assert.Equal(t, 409, resp.StatusCode, "Response should be 409")
}

func TestBridgeTypesController_Create_AdapterExistsError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBuffer(cltest.MustReadFile(t, "testdata/existing_core_adapter.json")),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 400)
}

func TestBridgeTypesController_Create_BindJSONError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString("}"),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 500)
}

func TestBridgeTypesController_Create_DatabaseError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 400)
}
