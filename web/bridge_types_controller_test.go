package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBridgeTypesController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	setupJobSpecsControllerIndex(app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs")
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestBridgeTypesController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt, err := setupBridgeControllerIndex(app)
	assert.NoError(t, err)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/specs?size=x")
	cltest.AssertServerResponse(t, resp, 422)

	resp = cltest.BasicAuthGet(app.Server.URL + "/v2/bridge_types?size=1")
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
	assert.Equal(t, bt[0].DefaultConfirmations, bridges[0].DefaultConfirmations, "should have the same DefaultConfirmations")

	resp = cltest.BasicAuthGet(app.Server.URL + links["next"].Href)
	cltest.AssertServerResponse(t, resp, 200)

	bridges = []models.BridgeType{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &bridges, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])
	assert.Len(t, bridges, 1)
	assert.Equal(t, bt[1].Name, bridges[0].Name, "should have the same Name")
	assert.Equal(t, bt[1].URL.String(), bridges[0].URL.String(), "should have the same URL")
	assert.Equal(t, bt[1].DefaultConfirmations, bridges[0].DefaultConfirmations, "should have the same DefaultConfirmations")
}

func setupBridgeControllerIndex(app *cltest.TestApplication) ([]*models.BridgeType, error) {

	bt1 := &models.BridgeType{Name: "testingbridges1",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	err := app.AddAdapter(bt1)
	if err != nil {
		return nil, err
	}

	bt2 := &models.BridgeType{Name: "testingbridges2",
		URL:                  cltest.WebURL("https://testing.com/tari"),
		DefaultConfirmations: 0}
	err = app.AddAdapter(bt2)

	return []*models.BridgeType{bt1, bt2}, err
}

func TestBridgeTypesController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/create_random_number_bridge_type.json")),
	)
	cltest.AssertServerResponse(t, resp, 200)
	btName := cltest.ParseCommonJSON(resp.Body).Name

	bt := &models.BridgeType{}
	assert.Nil(t, app.Store.One("Name", btName, bt))
	assert.Equal(t, "randomnumber", bt.Name)
	assert.Equal(t, uint64(10), bt.DefaultConfirmations)
	assert.Equal(t, "https://example.com/randomNumber", bt.URL.String())
}

func TestBridgeController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := &models.BridgeType{Name: "testingbridges1",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	err := app.AddAdapter(bt)
	assert.NoError(t, err)

	resp := cltest.BasicAuthGet(app.Server.URL + "/v2/bridge_types/" + bt.Name)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respBridge presenters.BridgeType
	json.Unmarshal(cltest.ParseResponseBody(resp), &respBridge)
	assert.Equal(t, respBridge.Name, bt.Name, "should have the same schedule")
	assert.Equal(t, respBridge.URL.String(), bt.URL.String(), "should have the same URL")
	assert.Equal(t, respBridge.DefaultConfirmations, bt.DefaultConfirmations, "should have the same DefaultConfirmations")

	resp = cltest.BasicAuthGet(app.Server.URL + "/v2/bridge_types/nosuchbridge")
	assert.Equal(t, 404, resp.StatusCode, "Response should be 404")
}

func TestBridgeController_RemoveOne(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	resp := cltest.BasicAuthDelete(app.Server.URL+"/v2/bridge_types/testingbridges1",
		"application/json",
		nil)
	assert.Equal(t, 404, resp.StatusCode, "Response should be 404")

	bt := &models.BridgeType{Name: "testingbridges2",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	err := app.AddAdapter(bt)
	assert.NoError(t, err)

	resp = cltest.BasicAuthDelete(app.Server.URL+"/v2/bridge_types/"+bt.Name,
		"application/json",
		nil)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	resp = cltest.BasicAuthGet(app.Server.URL + "/v2/bridge_types/testingbridges2")
	assert.Equal(t, 404, resp.StatusCode, "Response should be 404")
}

func TestBridgeController_RemoveMany(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	var bridges []models.BridgeType

	for i := 0; i < 8; i++ {
		bt := models.BridgeType{Name: fmt.Sprintf("testbridge%v", i),
			URL:                  cltest.WebURL(fmt.Sprintf("https://testing.com/bridges%v", i%2)),
			DefaultConfirmations: uint64(i % 5)}
		bridges = append(bridges, bt)
		err := app.AddAdapter(&bt)
		assert.NoError(t, err)
	}

	cases := []struct {
		name        string
		jsonInput   string
		searchCheck models.BridgeTypeCleaner
		resp        int
		nilAssert   bool
		expectedLen int
	}{
		{"value not found, no removals",
			`{"defaultConfirmations":8}`,
			models.BridgeTypeCleaner{bridges[3], ""},
			500,
			false,
			1,
		},
		{"single value removal",
			`{"defaultConfirmations":3}`,
			models.BridgeTypeCleaner{bridges[3], ""},
			200,
			true,
			0,
		},
		{"multiple value removal",
			`{"url":"https://testing.com/bridges1"}`,
			models.BridgeTypeCleaner{bridges[1], ""},
			200,
			true,
			0,
		},
		{"regex based multiple value removal",
			`{"name":"^test.+[0]"}`,
			models.BridgeTypeCleaner{bridges[0], ""},
			200,
			true,
			0,
		},
		{"empty input, remove all values",
			`{}`,
			models.BridgeTypeCleaner{bridges[5], ""},
			200,
			true,
			0,
		},
	}

	for _, test := range cases {
		resp := cltest.BasicAuthDelete(app.Server.URL+"/v2/bridge_types",
			"application/json",
			bytes.NewBufferString(test.jsonInput))
		assert.Equal(t, test.resp, resp.StatusCode, "Response should be successful")
		query, err := app.Store.AdvancedBridgeSearch(test.searchCheck)
		assert.Equal(t, test.nilAssert, err != nil)
		assert.Equal(t, test.expectedLen, len(query))
	}
}

func TestBridgeController_AdvancedSearch(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	bt := &models.BridgeType{Name: "testingbridges1",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	err := app.AddAdapter(bt)
	assert.NoError(t, err)

	bt = &models.BridgeType{Name: "testingbridges2",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	err = app.AddAdapter(bt)
	assert.NoError(t, err)

	cases := []struct {
		name        string
		search      models.BridgeTypeCleaner
		expectedLen int
		nilAssert   bool
	}{
		{"Search by URL",
			models.BridgeTypeCleaner{models.BridgeType{URL: cltest.WebURL("https://testing.com/bridges")}, ""},
			2,
			false,
		},
		{"Search by Name",
			models.BridgeTypeCleaner{models.BridgeType{Name: "^test"}, ""},
			2,
			false,
		},
		{"Not found",
			models.BridgeTypeCleaner{models.BridgeType{Name: "notabridge"}, ""},
			0,
			true,
		},
	}

	for _, test := range cases {
		query, err := app.Store.AdvancedBridgeSearch(test.search)
		assert.Equal(t, test.nilAssert, err != nil)
		assert.Equal(t, test.expectedLen, len(query))
	}

}

func TestBridgeTypesController_Create_AdapterExistsError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/existing_core_adapter.json")),
	)
	cltest.AssertServerResponse(t, resp, 400)
}

func TestBridgeTypesController_Create_BindJSONError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBufferString("}"),
	)
	cltest.AssertServerResponse(t, resp, 500)
}

func TestBridgeTypesController_Create_DatabaseError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBufferString(`{"url":"http://without.a.name"}`),
	)
	cltest.AssertServerResponse(t, resp, 400)
}
